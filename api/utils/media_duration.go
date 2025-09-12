package utils

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// AudioDuration returns duration of an audio file.
// Supported formats: MP3, WAV (auto-detected by header)
func AudioDuration(path string) (time.Duration, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// Peek first 12 bytes to detect format
	head := make([]byte, 12)
	n, err := io.ReadFull(f, head)
	if err != nil {
		return 0, err
	}
	if n < 12 {
		return 0, errors.New("file too small")
	}

	// WAV: RIFF....WAVE
	if string(head[0:4]) == "RIFF" && string(head[8:12]) == "WAVE" {
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return 0, err
		}
		return wavDuration(f)
	}

	// MP3 can start with ID3 or frame sync 0xFFEx
	if string(head[0:3]) == "ID3" || (head[0] == 0xFF && (head[1]&0xE0) == 0xE0) {
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return 0, err
		}
		return mp3Duration(f)
	}

	return 0, errors.New("unsupported audio format")
}

// AudioDurationFromURL downloads the url to a temp file and returns duration.
func AudioDurationFromURL(url string) (time.Duration, error) {
	path, err := fetchURLToTemp(url, 30*time.Second)
	if err != nil {
		return 0, err
	}
	defer os.Remove(path)
	return AudioDuration(path)
}

// VideoDurationMP4 returns duration of an MP4 file (MOV/MP4 base media).
func VideoDurationMP4(path string) (time.Duration, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return mp4Duration(f)
}

// VideoDurationMP4FromURL downloads the url to a temp file and returns duration.
func VideoDurationMP4FromURL(url string) (time.Duration, error) {
	path, err := fetchURLToTemp(url, 30*time.Second)
	if err != nil {
		return 0, err
	}
	defer os.Remove(path)
	return VideoDurationMP4(path)
}

// ---------------------- helpers ----------------------

func fetchURLToTemp(url string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status: %d", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "media-*")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		path := tmp.Name()
		_ = os.Remove(path)
		return "", err
	}
	return tmp.Name(), nil
}

// ---------------------- WAV ----------------------

func wavDuration(r io.ReadSeeker) (time.Duration, error) {
	// RIFF header already checked outside if needed. We parse chunks to get fmt and data.
	// WAV little-endian
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	// Read RIFF header (12 bytes)
	head := make([]byte, 12)
	if _, err := io.ReadFull(r, head); err != nil {
		return 0, err
	}
	if string(head[0:4]) != "RIFF" || string(head[8:12]) != "WAVE" {
		return 0, errors.New("invalid wav header")
	}

	var sampleRate uint32
	var numChans uint16
	var bitsPerSample uint16
	var byteRate uint32
	var dataSize uint32

	for {
		chunkHdr := make([]byte, 8)
		if _, err := io.ReadFull(r, chunkHdr); err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			return 0, err
		}
		ckID := string(chunkHdr[0:4])
		ckSize := binary.LittleEndian.Uint32(chunkHdr[4:8])

		switch ckID {
		case "fmt ":
			fmtData := make([]byte, ckSize)
			if _, err := io.ReadFull(r, fmtData); err != nil {
				return 0, err
			}
			// audioFormat := binary.LittleEndian.Uint16(fmtData[0:2]) // 1 = PCM
			numChans = binary.LittleEndian.Uint16(fmtData[2:4])
			sampleRate = binary.LittleEndian.Uint32(fmtData[4:8])
			byteRate = binary.LittleEndian.Uint32(fmtData[8:12])
			// blockAlign := binary.LittleEndian.Uint16(fmtData[12:14])
			if len(fmtData) >= 16 {
				bitsPerSample = binary.LittleEndian.Uint16(fmtData[14:16])
			}
		case "data":
			dataSize = ckSize
			// Skip data content
			if _, err := r.Seek(int64(ckSize), io.SeekCurrent); err != nil {
				return 0, err
			}
		default:
			// Skip other chunks
			if _, err := r.Seek(int64(ckSize), io.SeekCurrent); err != nil {
				return 0, err
			}
		}
		// Chunks are word-aligned (pad byte if odd size)
		if ckSize%2 == 1 {
			if _, err := r.Seek(1, io.SeekCurrent); err != nil { // skip pad
				return 0, err
			}
		}
	}

	if sampleRate == 0 || numChans == 0 {
		return 0, errors.New("invalid wav fmt")
	}

	var durationSeconds float64
	if byteRate != 0 {
		durationSeconds = float64(dataSize) / float64(byteRate)
	} else {
		bytesPerSec := float64(sampleRate) * float64(numChans) * float64(bitsPerSample) / 8.0
		if bytesPerSec == 0 {
			return 0, errors.New("invalid wav parameters")
		}
		durationSeconds = float64(dataSize) / bytesPerSec
	}
	return time.Duration(durationSeconds * float64(time.Second)), nil
}

// ---------------------- MP3 ----------------------

func mp3Duration(r io.ReadSeeker) (time.Duration, error) {
	// Strategy:
	// 1) Skip ID3v2 header if present.
	// 2) Try read first frame and detect XING/Info or VBRI to get total frames and duration.
	// 3) If VBR headers not present, fall back to CBR estimation: (audioDataBytes * 8) / bitrate.

	// File size
	fi, err := fileSizeFromSeeker(r)
	if err != nil {
		return 0, err
	}

	// Skip ID3v2
	var id3v2Size int64
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}
	id3v2Size, err = skipID3v2(r)
	if err != nil {
		return 0, err
	}

	// Remember audio start offset
	startOffset, _ := r.Seek(0, io.SeekCurrent)

	// Read first frame header (search sync)
	off, fh, err := findNextMP3Frame(r)
	if err != nil {
		return 0, err
	}
	if _, err := r.Seek(off, io.SeekStart); err != nil {
		return 0, err
	}

	// Check for XING/Info header in first frame (for VBR)
	totalFrames, sampleRate, samplesPerFrame, bitrateKbps, vbrFound, err := parseFirstFrameForVBR(r, fh)
	if err != nil {
		return 0, err
	}

	if vbrFound && totalFrames > 0 && sampleRate > 0 && samplesPerFrame > 0 {
		seconds := (float64(totalFrames) * float64(samplesPerFrame)) / float64(sampleRate)
		return time.Duration(seconds * float64(time.Second)), nil
	}

	// Fall back to CBR estimate using bitrate and data size (excluding ID3v2 and ID3v1)
	// Detect ID3v1 at end (128 bytes TAG)
	var id3v1Size int64
	if fi >= 128 {
		if _, err := r.Seek(fi-128, io.SeekStart); err == nil {
			buf := make([]byte, 3)
			if _, err := io.ReadFull(r, buf); err == nil {
				if string(buf) == "TAG" {
					id3v1Size = 128
				}
			}
		}
	}

	audioBytes := fi - id3v2Size - id3v1Size - startOffset
	if audioBytes <= 0 || bitrateKbps == 0 {
		return 0, errors.New("unable to estimate mp3 duration")
	}
	// bitrateKbps in kbps, bytes -> bits
	seconds := float64(audioBytes*8) / float64(bitrateKbps*1000)
	return time.Duration(seconds * float64(time.Second)), nil
}

type mp3FrameHeader struct {
	Version     int // 1: MPEG1, 2: MPEG2, 25: MPEG2.5
	Layer       int // 1,2,3
	BitrateKbps int
	SampleRate  int
	Padding     int
	ChannelMode int // 0:Stereo,1:Joint,2:Dual,3:Mono
}

func findNextMP3Frame(r io.ReadSeeker) (int64, mp3FrameHeader, error) {
	var hdr mp3FrameHeader
	// Start from current pos and scan up to 64KB
	start, _ := r.Seek(0, io.SeekCurrent)
	limit := int64(64 * 1024)
	buf := make([]byte, limit)
	n, err := r.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, hdr, err
	}
	for i := 0; i+4 <= n; i++ {
		if buf[i] == 0xFF && (buf[i+1]&0xE0) == 0xE0 { // sync
			if h, ok := parseMP3Header(buf[i : i+4]); ok {
				offset := start + int64(i)
				return offset, h, nil
			}
		}
	}
	return 0, hdr, errors.New("mp3 frame not found")
}

func parseMP3Header(b []byte) (mp3FrameHeader, bool) {
	var h mp3FrameHeader
	if len(b) < 4 {
		return h, false
	}
	if b[0] != 0xFF || (b[1]&0xE0) != 0xE0 {
		return h, false
	}
	versionBits := (b[1] >> 3) & 0x03
	layerBits := (b[1] >> 1) & 0x03
	bitrateBits := (b[2] >> 4) & 0x0F
	sampleRateBits := (b[2] >> 2) & 0x03
	paddingBit := (b[2] >> 1) & 0x01
	channelMode := (b[3] >> 6) & 0x03

	var version int
	switch versionBits {
	case 0x00:
		version = 25 // MPEG 2.5
	case 0x02:
		version = 2 // MPEG 2
	case 0x03:
		version = 1 // MPEG 1
	default:
		return h, false
	}

	var layer int
	switch layerBits {
	case 0x01:
		layer = 3
	case 0x02:
		layer = 2
	case 0x03:
		layer = 1
	default:
		return h, false
	}

	br := mp3BitrateKbps(version, layer, int(bitrateBits))
	if br == 0 {
		return h, false
	}
	sr := mp3SampleRate(version, int(sampleRateBits))
	if sr == 0 {
		return h, false
	}

	h = mp3FrameHeader{
		Version:     version,
		Layer:       layer,
		BitrateKbps: br,
		SampleRate:  sr,
		Padding:     int(paddingBit),
		ChannelMode: int(channelMode),
	}
	return h, true
}

func mp3BitrateKbps(version, layer, index int) int {
	// index: 1..14 valid; 0,15 invalid
	if index <= 0 || index == 15 {
		return 0
	}
	// Tables per ISO/IEC 11172-3/13818-3 (common subset)
	var tbl [15]int
	if layer == 1 { // Layer I
		if version == 1 { // MPEG1
			tbl = [15]int{0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448}
		} else { // MPEG2/2.5
			tbl = [15]int{0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256}
		}
	} else if layer == 2 { // Layer II
		if version == 1 {
			tbl = [15]int{0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384}
		} else {
			tbl = [15]int{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160}
		}
	} else { // Layer III
		if version == 1 {
			tbl = [15]int{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320}
		} else {
			tbl = [15]int{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160}
		}
	}
	return tbl[index]
}

func mp3SampleRate(version, index int) int {
	if index == 3 {
		return 0
	}
	// base table for MPEG1
	base := [3]int{44100, 48000, 32000}
	sr := base[index]
	if version == 2 { // MPEG2
		sr /= 2
	} else if version == 25 { // MPEG2.5
		sr /= 4
	}
	return sr
}

func samplesPerMP3Frame(version, layer int) int {
	switch layer {
	case 1:
		return 384
	case 2:
		return 1152
	case 3:
		if version == 1 {
			return 1152
		}
		return 576 // MPEG2/2.5 Layer III
	default:
		return 0
	}
}

func parseFirstFrameForVBR(r io.ReadSeeker, fh mp3FrameHeader) (totalFrames uint32, sampleRate int, samplesPerFrame int, bitrateKbps int, vbrFound bool, err error) {
	// After the 4-byte header, possible side info and then XING/Info
	if _, err = r.Seek(0, io.SeekCurrent); err != nil {
		return
	}
	// Re-read header
	hdr := make([]byte, 4)
	if _, err = io.ReadFull(r, hdr); err != nil {
		return
	}

	// side info size depends on MPEG version and channel mode (for Layer III)
	sideInfoSize := 0
	if fh.Layer == 3 { // Layer III
		if fh.Version == 1 { // MPEG1
			if fh.ChannelMode == 3 { // mono
				sideInfoSize = 17
			} else {
				sideInfoSize = 32
			}
		} else { // MPEG2/2.5
			if fh.ChannelMode == 3 {
				sideInfoSize = 9
			} else {
				sideInfoSize = 17
			}
		}
	}

	// Read next up to 120 bytes to search for XING/Info or VBRI
	buf := make([]byte, sideInfoSize+120)
	if _, err = io.ReadFull(r, buf); err != nil {
		// If short, still try within available
		if !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
			return
		}
	}

	// Search XING/Info signature
	sigs := [][]byte{[]byte("Xing"), []byte("Info")}
	for _, sig := range sigs {
		idx := indexOf(buf, sig)
		if idx >= 0 {
			// flags after signature (4 bytes), then if frames flag set, 4 bytes frames
			if len(buf) >= idx+4+4 {
				flags := binary.BigEndian.Uint32(buf[idx+4 : idx+8])
				var frames uint32
				if (flags & 0x01) != 0 { // frames present
					if len(buf) >= idx+8+4 {
						frames = binary.BigEndian.Uint32(buf[idx+8 : idx+12])
					}
				}
				if frames > 0 {
					vbrFound = true
					totalFrames = frames
					sampleRate = fh.SampleRate
					samplesPerFrame = samplesPerMP3Frame(fh.Version, fh.Layer)
					bitrateKbps = fh.BitrateKbps
					return
				}
			}
		}
	}

	// Check VBRI (usually at 32 bytes after header for MPEG1 Layer III)
	if len(buf) >= 4 {
		idx := indexOf(buf, []byte("VBRI"))
		if idx >= 0 {
			if len(buf) >= idx+4+2+2+4+4 {
				// VBRI layout: 'VBRI'(4) + version(2) + delay(2) + quality(2?) varies; but at offset 10 comes bytes: bytes (4), frames (4)
				// Some docs: offset 10: bytes, offset 14: frames (big-endian)
				bytesOffset := idx + 10
				framesOffset := idx + 14
				if len(buf) >= framesOffset+4 {
					frames := binary.BigEndian.Uint32(buf[framesOffset : framesOffset+4])
					if frames > 0 {
						vbrFound = true
						totalFrames = frames
						sampleRate = fh.SampleRate
						samplesPerFrame = samplesPerMP3Frame(fh.Version, fh.Layer)
						bitrateKbps = fh.BitrateKbps
						_ = bytesOffset // not used
						return
					}
				}
			}
		}
	}

	// No VBR header. Provide header info for CBR fallback
	sampleRate = fh.SampleRate
	samplesPerFrame = samplesPerMP3Frame(fh.Version, fh.Layer)
	bitrateKbps = fh.BitrateKbps
	return
}

func indexOf(haystack []byte, needle []byte) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		match := true
		for j := 0; j < len(needle); j++ {
			if haystack[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func skipID3v2(r io.ReadSeeker) (int64, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}
	head := make([]byte, 10)
	if _, err := io.ReadFull(r, head); err != nil {
		return 0, nil // no header
	}
	if string(head[0:3]) != "ID3" {
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			return 0, err
		}
		return 0, nil
	}
	// size: 4 synchsafe bytes
	sz := int64((int(head[6]&0x7F) << 21) | (int(head[7]&0x7F) << 14) | (int(head[8]&0x7F) << 7) | int(head[9]&0x7F))
	// total header size = 10 + sz (+ footer 10 if flag set)
	footer := int64(0)
	if (head[5] & 0x10) != 0 { // footer present
		footer = 10
	}
	total := 10 + sz + footer
	if _, err := r.Seek(total, io.SeekStart); err != nil {
		return 0, err
	}
	return total, nil
}

func fileSizeFromSeeker(r io.ReadSeeker) (int64, error) {
	cur, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	end, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	if _, err := r.Seek(cur, io.SeekStart); err != nil {
		return 0, err
	}
	return end, nil
}

// ---------------------- MP4 ----------------------

type mp4BoxHeader struct {
	Size uint64
	Type [4]byte
}

func readBoxHeader(r io.ReadSeeker) (mp4BoxHeader, error) {
	var h mp4BoxHeader
	buf := make([]byte, 8)
	if _, err := io.ReadFull(r, buf); err != nil {
		return h, err
	}
	sz := binary.BigEndian.Uint32(buf[0:4])
	copy(h.Type[:], buf[4:8])
	if sz == 1 {
		// 64-bit size follows
		ext := make([]byte, 8)
		if _, err := io.ReadFull(r, ext); err != nil {
			return h, err
		}
		h.Size = binary.BigEndian.Uint64(ext)
	} else {
		h.Size = uint64(sz)
	}
	return h, nil
}

func skipBox(r io.ReadSeeker, boxSize uint64, alreadyRead int64) error {
	toSkip := int64(boxSize) - alreadyRead
	if toSkip < 0 {
		return fmt.Errorf("invalid box size")
	}
	_, err := r.Seek(toSkip, io.SeekCurrent)
	return err
}

func mp4Duration(r io.ReadSeeker) (time.Duration, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}
	// Find moov box
	var moovStart int64
	var moovSize uint64
	for {
		pos, _ := r.Seek(0, io.SeekCurrent)
		h, err := readBoxHeader(r)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			return 0, err
		}
		if string(h.Type[:]) == "moov" {
			moovStart = pos
			moovSize = h.Size
			break
		}
		if h.Size < 8 {
			return 0, errors.New("invalid mp4 box size")
		}
		if err := skipBox(r, h.Size, 8); err != nil {
			return 0, err
		}
	}
	if moovStart == 0 && moovSize == 0 {
		return 0, errors.New("moov not found")
	}
	// Parse inside moov for video trak mdhd, else mvhd
	if _, err := r.Seek(moovStart+8, io.SeekStart); err != nil { // skip moov header
		return 0, err
	}
	end := moovStart + int64(moovSize)
	var movieTimescale uint32
	var movieDuration uint64
	var foundVideoMdhd bool
	var mdhdTimescale uint32
	var mdhdDuration uint64

	for {
		pos, _ := r.Seek(0, io.SeekCurrent)
		if pos >= end {
			break
		}
		h, err := readBoxHeader(r)
		if err != nil {
			return 0, err
		}
		switch string(h.Type[:]) {
		case "mvhd":
			// movie header
			ver := make([]byte, 1)
			if _, err := io.ReadFull(r, ver); err != nil {
				return 0, err
			}
			if _, err := r.Seek(3, io.SeekCurrent); err != nil { // flags
				return 0, err
			}
			if ver[0] == 1 {
				// version 1: 64-bit duration
				buf := make([]byte, 8+8+4+8) // ctime(8) mtime(8) timescale(4) duration(8)
				if _, err := io.ReadFull(r, buf); err != nil {
					return 0, err
				}
				movieTimescale = binary.BigEndian.Uint32(buf[16:20])
				movieDuration = binary.BigEndian.Uint64(buf[20:28])
			} else {
				buf := make([]byte, 4+4+4+4) // ctime mtime timescale duration
				if _, err := io.ReadFull(r, buf); err != nil {
					return 0, err
				}
				movieTimescale = binary.BigEndian.Uint32(buf[8:12])
				movieDuration = uint64(binary.BigEndian.Uint32(buf[12:16]))
			}
			// skip rest of mvhd
			read := int64(1 + 3)
			if ver[0] == 1 {
				read += int64(8 + 8 + 4 + 8)
			} else {
				read += int64(4 + 4 + 4 + 4)
			}
			if err := skipBox(r, h.Size, 8+read); err != nil {
				return 0, err
			}
		case "trak":
			// parse trak for hdlr and mdhd
			tEnd := int64(0)
			if h.Size < 8 {
				return 0, errors.New("invalid trak size")
			}
			tEnd = pos + int64(h.Size)
			var isVideo bool
			var tMdhdTimescale uint32
			var tMdhdDuration uint64
			for {
				cpos, _ := r.Seek(0, io.SeekCurrent)
				if cpos >= tEnd {
					break
				}
				ch, err := readBoxHeader(r)
				if err != nil {
					return 0, err
				}
				switch string(ch.Type[:]) {
				case "mdia":
					mEnd := cpos + int64(ch.Size)
					for {
						mpos, _ := r.Seek(0, io.SeekCurrent)
						if mpos >= mEnd {
							break
						}
						mh, err := readBoxHeader(r)
						if err != nil {
							return 0, err
						}
						switch string(mh.Type[:]) {
						case "hdlr":
							// skip version+flags (4), pre_defined(4)
							if _, err := r.Seek(8, io.SeekCurrent); err != nil {
								return 0, err
							}
							handler := make([]byte, 4)
							if _, err := io.ReadFull(r, handler); err != nil {
								return 0, err
							}
							if string(handler) == "vide" {
								isVideo = true
							}
							if err := skipBox(r, mh.Size, 8+8+4); err != nil { // header + skipped + read handler
								return 0, err
							}
						case "mdhd":
							ver := make([]byte, 1)
							if _, err := io.ReadFull(r, ver); err != nil {
								return 0, err
							}
							if _, err := r.Seek(3, io.SeekCurrent); err != nil { // flags
								return 0, err
							}
							if ver[0] == 1 {
								buf := make([]byte, 8+8+4+8)
								if _, err := io.ReadFull(r, buf); err != nil {
									return 0, err
								}
								tMdhdTimescale = binary.BigEndian.Uint32(buf[16:20])
								tMdhdDuration = binary.BigEndian.Uint64(buf[20:28])
							} else {
								buf := make([]byte, 4+4+4+4)
								if _, err := io.ReadFull(r, buf); err != nil {
									return 0, err
								}
								tMdhdTimescale = binary.BigEndian.Uint32(buf[8:12])
								tMdhdDuration = uint64(binary.BigEndian.Uint32(buf[12:16]))
							}
							if err := skipBox(r, mh.Size, 8+1+3+int64(lenVersionPayload(ver[0]))); err != nil {
								return 0, err
							}
						default:
							if err := skipBox(r, mh.Size, 8); err != nil {
								return 0, err
							}
						}
					}
				default:
					if err := skipBox(r, ch.Size, 8); err != nil {
						return 0, err
					}
				}
			}
			if isVideo && tMdhdTimescale != 0 && tMdhdDuration != 0 {
				foundVideoMdhd = true
				mdhdTimescale = tMdhdTimescale
				mdhdDuration = tMdhdDuration
			}
			// Skip remaining of trak if any
			if _, err := r.Seek(tEnd, io.SeekStart); err != nil {
				return 0, err
			}
		default:
			if err := skipBox(r, h.Size, 8); err != nil {
				return 0, err
			}
		}
	}

	if foundVideoMdhd && mdhdTimescale != 0 {
		sec := float64(mdhdDuration) / float64(mdhdTimescale)
		return time.Duration(sec * float64(time.Second)), nil
	}
	if movieTimescale != 0 {
		sec := float64(movieDuration) / float64(movieTimescale)
		return time.Duration(sec * float64(time.Second)), nil
	}
	return 0, errors.New("failed to read mp4 duration")
}

func lenVersionPayload(ver byte) int {
	if ver == 1 {
		return 8 + 8 + 4 + 8
	}
	return 4 + 4 + 4 + 4
}

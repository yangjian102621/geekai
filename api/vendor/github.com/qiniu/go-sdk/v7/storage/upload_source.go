package storage

import (
	"errors"
	"io"
	"os"
)

type UploadSource interface {
	Size() int64
	Rewindable() bool
	Rewind() error
}

func NewUploadSourceReader(reader io.Reader, size int64) (UploadSource, error) {
	return &uploadSourceReader{
		reader: reader,
		size:   size,
	}, nil
}

type uploadSourceReader struct {
	reader io.Reader
	size   int64
}

func (u *uploadSourceReader) Rewindable() bool {
	return false
}

func (u *uploadSourceReader) Rewind() error {
	return errors.New("resource not support rewind")
}

func (u *uploadSourceReader) Size() int64 {
	return u.size
}

func NewUploadSourceReaderAt(reader io.ReaderAt, size int64) (UploadSource, error) {
	if size <= 0 {
		return nil, errors.New("source size must be set")
	}

	return &uploadSourceReaderAt{
		reader: reader,
		size:   size,
	}, nil
}

type uploadSourceReaderAt struct {
	reader io.ReaderAt
	size   int64
}

func (u *uploadSourceReaderAt) Rewindable() bool {
	return true
}

func (u *uploadSourceReaderAt) Rewind() error {
	return nil
}

func (u *uploadSourceReaderAt) Size() int64 {
	return u.size
}

func NewUploadSourceFile(filePath string) (UploadSource, error) {
	if fileInfo, err := os.Stat(filePath); err != nil {
		return nil, err
	} else {
		return &uploadSourceFile{
			fileInfo: fileInfo,
			filePath: filePath,
		}, nil
	}
}

type uploadSourceFile struct {
	filePath string
	fileInfo os.FileInfo
}

func (u *uploadSourceFile) Rewindable() bool {
	return true
}

func (u *uploadSourceFile) Rewind() error {
	return nil
}

func (u *uploadSourceFile) Size() int64 {
	return u.fileInfo.Size()
}

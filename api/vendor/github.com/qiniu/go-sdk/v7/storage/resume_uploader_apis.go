package storage

import (
	"context"
	"encoding/base64"
	"github.com/qiniu/go-sdk/v7/internal/hostprovider"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/conf"
)

type resumeUploaderAPIs struct {
	Client *client.Client
	Cfg    *Config
}

// BlkputRet 表示分片上传每个片上传完毕的返回值
type BlkputRet struct {
	Ctx        string `json:"ctx"`
	Checksum   string `json:"checksum"`
	Crc32      uint32 `json:"crc32"`
	Offset     uint32 `json:"offset"`
	Host       string `json:"host"`
	ExpiredAt  int64  `json:"expired_at"`
	chunkSize  int
	fileOffset int64
	blkIdx     int
}

func (p *resumeUploaderAPIs) mkBlk(ctx context.Context, upToken, upHost string, ret *BlkputRet, blockSize int64, body io.Reader, size int64) error {
	reqUrl := upHost + "/mkblk/" + strconv.FormatInt(blockSize, 10)

	return p.Client.CallWith64(ctx, ret, "POST", reqUrl, makeHeadersForUpload(upToken), body, size)
}

func (p *resumeUploaderAPIs) bput(ctx context.Context, upToken string, ret *BlkputRet, body io.Reader, size int64) error {
	reqUrl := ret.Host + "/bput/" + ret.Ctx + "/" + strconv.FormatUint(uint64(ret.Offset), 10)

	return p.Client.CallWith64(ctx, ret, "POST", reqUrl, makeHeadersForUpload(upToken), body, size)
}

// RputExtra 表示分片上传额外可以指定的参数
type RputExtra struct {
	Recorder           Recorder          // 可选。上传进度记录
	Params             map[string]string // 可选。用户自定义参数，以"x:"开头，而且值不能为空，否则忽略
	UpHost             string
	MimeType           string                                        // 可选。
	ChunkSize          int                                           // 可选。每次上传的Chunk大小
	TryTimes           int                                           // 可选。尝试次数
	HostFreezeDuration time.Duration                                 // 可选。主备域名冻结时间（默认：600s），当一个域名请求失败（单个域名会被重试 TryTimes 次），会被冻结一段时间，使用备用域名进行重试，在冻结时间内，域名不能被使用，当一个操作中所有域名竣备冻结操作不在进行重试，返回最后一次操作的错误。
	Progresses         []BlkputRet                                   // 可选。上传进度
	Notify             func(blkIdx int, blkSize int, ret *BlkputRet) // 可选。进度提示（注意多个block是并行传输的）
	NotifyErr          func(blkIdx int, blkSize int, err error)
}

func (extra *RputExtra) init() {
	if extra.ChunkSize == 0 {
		extra.ChunkSize = settings.ChunkSize
	}
	if extra.TryTimes == 0 {
		extra.TryTimes = settings.TryTimes
	}
	if extra.HostFreezeDuration <= 0 {
		extra.HostFreezeDuration = 10 * 60 * time.Second
	}
	if extra.Notify == nil {
		extra.Notify = func(blkIdx, blkSize int, ret *BlkputRet) {}
	}
	if extra.NotifyErr == nil {
		extra.NotifyErr = func(blkIdx, blkSize int, err error) {}
	}
}

func (extra *RputExtra) getUpHost(useHttps bool) string {
	return hostAddSchemeIfNeeded(useHttps, extra.UpHost)
}

func (p *resumeUploaderAPIs) mkfile(ctx context.Context, upToken, upHost string, ret interface{}, key string, hasKey bool, fsize int64, extra *RputExtra) (err error) {
	url := upHost + "/mkfile/" + strconv.FormatInt(fsize, 10)
	if extra == nil {
		extra = &RputExtra{}
	}
	if extra.MimeType != "" {
		url += "/mimeType/" + encode(extra.MimeType)
	}
	if hasKey {
		url += "/key/" + encode(key)
	}
	for k, v := range extra.Params {
		if (strings.HasPrefix(k, "x:") || strings.HasPrefix(k, "x-qn-meta-")) && v != "" {
			url += "/" + k + "/" + encode(v)
		}
	}
	ctxs := make([]string, len(extra.Progresses))
	for i, progress := range extra.Progresses {
		ctxs[i] = progress.Ctx
	}
	buf := strings.Join(ctxs, ",")
	return p.Client.CallWith(ctx, ret, "POST", url, makeHeadersForUpload(upToken), strings.NewReader(buf), len(buf))
}

// InitPartsRet 表示分片上传 v2 初始化完毕的返回值
type InitPartsRet struct {
	UploadID string `json:"uploadId"`
	ExpireAt int64  `json:"expireAt"`
}

func (p *resumeUploaderAPIs) initParts(ctx context.Context, upToken, upHost, bucket, key string, hasKey bool, ret *InitPartsRet) error {
	reqUrl := upHost + "/buckets/" + bucket + "/objects/" + encodeV2(key, hasKey) + "/uploads"

	return p.Client.CallWith(ctx, ret, "POST", reqUrl, makeHeadersForUploadEx(upToken, ""), nil, 0)
}

// UploadPartsRet 表示分片上传 v2 每个片上传完毕的返回值
type UploadPartsRet struct {
	Etag string `json:"etag"`
	MD5  string `json:"md5"`
}

func (p *resumeUploaderAPIs) uploadParts(ctx context.Context, upToken, upHost, bucket, key string, hasKey bool, uploadId string, partNumber int64, partMD5 string, ret *UploadPartsRet, body io.Reader, size int64) error {
	reqUrl := upHost + "/buckets/" + bucket + "/objects/" + encodeV2(key, hasKey) + "/uploads/" + uploadId + "/" + strconv.FormatInt(partNumber, 10)

	return p.Client.CallWith64(ctx, ret, "PUT", reqUrl, makeHeadersForUploadPart(upToken, partMD5), body, size)
}

type UploadPartInfo struct {
	Etag       string `json:"etag"`
	PartNumber int64  `json:"partNumber"`
	partSize   int
	fileOffset int64
}

// RputV2Extra 表示分片上传 v2 额外可以指定的参数
type RputV2Extra struct {
	Recorder           Recorder          // 可选。上传进度记录
	Metadata           map[string]string // 可选。用户自定义文件 metadata 信息
	CustomVars         map[string]string // 可选。用户自定义参数，以"x:"开头，而且值不能为空，否则忽略
	UpHost             string
	MimeType           string                                      // 可选。
	PartSize           int64                                       // 可选。每次上传的块大小
	TryTimes           int                                         // 可选。尝试次数
	HostFreezeDuration time.Duration                               // 可选。主备域名冻结时间（默认：600s），当一个域名请求失败（单个域名会被重试 TryTimes 次），会被冻结一段时间，使用备用域名进行重试，在冻结时间内，域名不能被使用，当一个操作中所有域名竣备冻结操作不在进行重试，返回最后一次操作的错误。
	Progresses         []UploadPartInfo                            // 上传进度
	Notify             func(partNumber int64, ret *UploadPartsRet) // 可选。进度提示（注意多个block是并行传输的）
	NotifyErr          func(partNumber int64, err error)
}

func (extra *RputV2Extra) init() {
	if extra.PartSize == 0 {
		extra.PartSize = settings.PartSize
	}
	if extra.TryTimes == 0 {
		extra.TryTimes = settings.TryTimes
	}
	if extra.HostFreezeDuration <= 0 {
		extra.HostFreezeDuration = 10 * 60 * time.Second
	}
	if extra.Notify == nil {
		extra.Notify = func(partNumber int64, ret *UploadPartsRet) {}
	}
	if extra.NotifyErr == nil {
		extra.NotifyErr = func(partNumber int64, err error) {}
	}
}

func (extra *RputV2Extra) getUpHost(useHttps bool) string {
	return hostAddSchemeIfNeeded(useHttps, extra.UpHost)
}

func hostAddSchemeIfNeeded(useHttps bool, host string) string {
	if host == "" {
		return ""
	} else if strings.Contains(host, "://") {
		return host
	} else {
		return endpoint(useHttps, host)
	}
}

func (p *resumeUploaderAPIs) completeParts(ctx context.Context, upToken, upHost string, ret interface{}, bucket, key string, hasKey bool, uploadId string, extra *RputV2Extra) (err error) {
	type CompletePartBody struct {
		Parts      []UploadPartInfo  `json:"parts"`
		MimeType   string            `json:"mimeType,omitempty"`
		Metadata   map[string]string `json:"metadata,omitempty"`
		CustomVars map[string]string `json:"customVars,omitempty"`
	}
	if extra == nil {
		extra = &RputV2Extra{}
	}
	completePartBody := CompletePartBody{
		Parts:      extra.Progresses,
		MimeType:   extra.MimeType,
		Metadata:   extra.Metadata,
		CustomVars: make(map[string]string),
	}
	for k, v := range extra.CustomVars {
		if strings.HasPrefix(k, "x:") && v != "" {
			completePartBody.CustomVars[k] = v
		}
	}

	reqUrl := upHost + "/buckets/" + bucket + "/objects/" + encodeV2(key, hasKey) + "/uploads/" + uploadId

	return p.Client.CallWithJson(ctx, ret, "POST", reqUrl, makeHeadersForUploadEx(upToken, conf.CONTENT_TYPE_JSON), &completePartBody)
}

func (p *resumeUploaderAPIs) upHost(ak, bucket string) (upHost string, err error) {
	return getUpHost(p.Cfg, 0, 0, ak, bucket)
}

func (p *resumeUploaderAPIs) upHostProvider(ak, bucket string, retryMax int, hostFreezeDuration time.Duration) (hostProvider hostprovider.HostProvider, err error) {
	return getUpHostProvider(p.Cfg, retryMax, hostFreezeDuration, ak, bucket)
}

func makeHeadersForUpload(upToken string) http.Header {
	return makeHeadersForUploadEx(upToken, conf.CONTENT_TYPE_OCTET)
}

func makeHeadersForUploadPart(upToken, partMD5 string) http.Header {
	headers := makeHeadersForUpload(upToken)
	headers.Add("Content-MD5", partMD5)
	return headers
}

func makeHeadersForUploadEx(upToken, contentType string) http.Header {
	headers := http.Header{}
	if contentType != "" {
		headers.Add("Content-Type", contentType)
	}
	headers.Add("Authorization", "UpToken "+upToken)
	return headers
}

func encode(raw string) string {
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func encodeV2(key string, hasKey bool) string {
	if !hasKey {
		return "~"
	} else {
		return encode(key)
	}
}

type blkputRets []BlkputRet

func (rets blkputRets) Len() int {
	return len(rets)
}

func (rets blkputRets) Less(i, j int) bool {
	return rets[i].blkIdx < rets[j].blkIdx
}

func (rets blkputRets) Swap(i, j int) {
	rets[i], rets[j] = rets[j], rets[i]
}

type uploadPartInfos []UploadPartInfo

func (infos uploadPartInfos) Len() int {
	return len(infos)
}

func (infos uploadPartInfos) Less(i, j int) bool {
	return infos[i].PartNumber < infos[j].PartNumber
}

func (infos uploadPartInfos) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}

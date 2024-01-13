package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/internal/hostprovider"
	"hash"
	"hash/crc32"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/client"
)

// PutExtra 为表单上传的额外可选项
type PutExtra struct {
	// 可选，用户自定义参数，必须以 "x:" 开头。若不以x:开头，则忽略。
	Params map[string]string

	UpHost string

	TryTimes int // 可选。尝试次数

	// 主备域名冻结时间（默认：600s），当一个域名请求失败（单个域名会被重试 TryTimes 次），会被冻结一段时间，使用备用域名进行重试，在冻结时间内，域名不能被使用，当一个操作中所有域名竣备冻结操作不在进行重试，返回最后一次操作的错误。
	HostFreezeDuration time.Duration

	// 可选，当为 "" 时候，服务端自动判断。
	MimeType string

	// 上传事件：进度通知。这个事件的回调函数应该尽可能快地结束。
	OnProgress func(fsize, uploaded int64)
}

func (extra *PutExtra) init() {
	if extra.TryTimes == 0 {
		extra.TryTimes = settings.TryTimes
	}
	if extra.HostFreezeDuration <= 0 {
		extra.HostFreezeDuration = 10 * 60 * time.Second
	}
}

func (extra *PutExtra) getUpHost(useHttps bool) string {
	return hostAddSchemeIfNeeded(useHttps, extra.UpHost)
}

// PutRet 为七牛标准的上传回复内容。
// 如果使用了上传回调或者自定义了returnBody，那么需要根据实际情况，自己自定义一个返回值结构体
type PutRet struct {
	Hash         string `json:"hash"`
	PersistentID string `json:"persistentId"`
	Key          string `json:"key"`
}

// FormUploader 表示一个表单上传的对象
type FormUploader struct {
	Client *client.Client
	Cfg    *Config
}

// NewFormUploader 用来构建一个表单上传的对象
func NewFormUploader(cfg *Config) *FormUploader {
	if cfg == nil {
		cfg = &Config{}
	}

	return &FormUploader{
		Client: &client.DefaultClient,
		Cfg:    cfg,
	}
}

// NewFormUploaderEx 用来构建一个表单上传的对象
func NewFormUploaderEx(cfg *Config, clt *client.Client) *FormUploader {
	if cfg == nil {
		cfg = &Config{}
	}

	if clt == nil {
		clt = &client.DefaultClient
	}

	return &FormUploader{
		Client: clt,
		Cfg:    cfg,
	}
}

// PutFile 用来以表单方式上传一个文件，和 Put 不同的只是一个通过提供文件路径来访问文件内容，一个通过 io.Reader 来访问。
//
// ctx       是请求的上下文。
// ret       是上传成功后返回的数据。如果 uptoken 中没有设置 callbackUrl 或 returnBody，那么返回的数据结构是 PutRet 结构。
// uptoken   是由业务服务器颁发的上传凭证。
// key       是要上传的文件访问路径。比如："foo/bar.jpg"。注意我们建议 key 不要以 '/' 开头。另外，key 为空字符串是合法的。
// localFile 是要上传的文件的本地路径。
// extra     是上传的一些可选项，可以指定为nil。详细见 PutExtra 结构的描述。
func (p *FormUploader) PutFile(
	ctx context.Context, ret interface{}, uptoken, key, localFile string, extra *PutExtra) (err error) {
	return p.putFile(ctx, ret, uptoken, key, true, localFile, extra)
}

// PutFileWithoutKey 用来以表单方式上传一个文件。不指定文件上传后保存的key的情况下，文件命名方式首先看看
// uptoken 中是否设置了 saveKey，如果设置了 saveKey，那么按 saveKey 要求的规则生成 key，否则自动以文件的 hash 做 key。
// 和 Put 不同的只是一个通过提供文件路径来访问文件内容，一个通过 io.Reader 来访问。
//
// ctx       是请求的上下文。
// ret       是上传成功后返回的数据。如果 uptoken 中没有设置 CallbackUrl 或 ReturnBody，那么返回的数据结构是 PutRet 结构。
// uptoken   是由业务服务器颁发的上传凭证。
// localFile 是要上传的文件的本地路径。
// extra     是上传的一些可选项。可以指定为nil。详细见 PutExtra 结构的描述。
func (p *FormUploader) PutFileWithoutKey(
	ctx context.Context, ret interface{}, uptoken, localFile string, extra *PutExtra) (err error) {
	return p.putFile(ctx, ret, uptoken, "", false, localFile, extra)
}

func (p *FormUploader) putFile(
	ctx context.Context, ret interface{}, upToken string,
	key string, hasKey bool, localFile string, extra *PutExtra) (err error) {

	f, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}
	fsize := fi.Size()

	return p.put(ctx, ret, upToken, key, hasKey, f, fsize, extra, filepath.Base(localFile))
}

// Put 用来以表单方式上传一个文件。
//
// ctx     是请求的上下文。
// ret     是上传成功后返回的数据。如果 uptoken 中没有设置 callbackUrl 或 returnBody，那么返回的数据结构是 PutRet 结构。
// uptoken 是由业务服务器颁发的上传凭证。
// key     是要上传的文件访问路径。比如："foo/bar.jpg"。注意我们建议 key 不要以 '/' 开头。另外，key 为空字符串是合法的。
// data    是文件内容的访问接口（io.Reader）。
// fsize   是要上传的文件大小。
// extra   是上传的一些可选项。可以指定为nil。详细见 PutExtra 结构的描述。
func (p *FormUploader) Put(
	ctx context.Context, ret interface{}, uptoken, key string, data io.Reader, size int64, extra *PutExtra) (err error) {
	err = p.put(ctx, ret, uptoken, key, true, data, size, extra, path.Base(key))
	return
}

// PutWithoutKey 用来以表单方式上传一个文件。不指定文件上传后保存的key的情况下，文件命名方式首先看看 uptoken 中是否设置了 saveKey，
// 如果设置了 saveKey，那么按 saveKey 要求的规则生成 key，否则自动以文件的 hash 做 key。
//
// ctx     是请求的上下文。
// ret     是上传成功后返回的数据。如果 uptoken 中没有设置 CallbackUrl 或 ReturnBody，那么返回的数据结构是 PutRet 结构。
// uptoken 是由业务服务器颁发的上传凭证。
// data    是文件内容的访问接口（io.Reader）。
// fsize   是要上传的文件大小。
// extra   是上传的一些可选项。详细见 PutExtra 结构的描述。
func (p *FormUploader) PutWithoutKey(
	ctx context.Context, ret interface{}, uptoken string, data io.Reader, size int64, extra *PutExtra) (err error) {
	err = p.put(ctx, ret, uptoken, "", false, data, size, extra, "filename")
	return err
}

func (p *FormUploader) put(
	ctx context.Context, ret interface{}, upToken string,
	key string, hasKey bool, data io.Reader, size int64, extra *PutExtra, fileName string) error {

	if extra == nil {
		extra = &PutExtra{}
	}
	extra.init()

	seekableData, ok := data.(io.ReadSeeker)
	if !ok {
		dataBytes, rErr := ioutil.ReadAll(data)
		if rErr != nil {
			return rErr
		}
		if size <= 0 {
			size = int64(len(dataBytes))
		}
		seekableData = bytes.NewReader(dataBytes)
	}

	return p.putSeekableData(ctx, ret, upToken, key, hasKey, seekableData, size, extra, fileName)
}

func (p *FormUploader) putSeekableData(ctx context.Context, ret interface{}, upToken string,
	key string, hasKey bool, data io.ReadSeeker, dataSize int64, extra *PutExtra, fileName string) error {

	formFieldBuff := new(bytes.Buffer)
	formWriter := multipart.NewWriter(formFieldBuff)
	// 写入表单头、token、key、fileName 等信息
	if wErr := writeMultipart(formWriter, upToken, key, hasKey, extra, fileName); wErr != nil {
		return wErr
	}

	// 计算文件 crc32
	crc32Hash := crc32.NewIEEE()
	if _, cErr := io.Copy(crc32Hash, data); cErr != nil {
		return cErr
	}
	crcReader := newCrc32Reader(formWriter.Boundary(), crc32Hash)
	crcBytes, rErr := ioutil.ReadAll(crcReader)
	if rErr != nil {
		return rErr
	}
	crcReader = nil

	// 表单写入文件 crc32
	if _, wErr := formFieldBuff.Write(crcBytes); wErr != nil {
		return wErr
	}
	crcBytes = nil

	formHead := make(textproto.MIMEHeader)
	formHead.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`,
		escapeQuotes(fileName)))
	if extra.MimeType != "" {
		formHead.Set("Content-Type", extra.MimeType)
	}
	if _, cErr := formWriter.CreatePart(formHead); cErr != nil {
		return cErr
	}
	formHead = nil

	// 表单 Fields
	formFieldData := formFieldBuff.Bytes()
	formFieldBuff = nil

	// 表单最后一行
	formEndLine := []byte(fmt.Sprintf("\r\n--%s--\r\n", formWriter.Boundary()))

	// 不再重新构造 formBody ，避免内存峰值问题
	var formBodyLen int64 = -1
	if dataSize >= 0 {
		formBodyLen = int64(len(formFieldData)) + dataSize + int64(len(formEndLine))
	}

	progress := newUploadProgress(extra.OnProgress)
	getBodyReader := func() (io.Reader, error) {
		if _, err := data.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		var formReader = io.MultiReader(bytes.NewReader(formFieldData), data, bytes.NewReader(formEndLine))
		if extra.OnProgress != nil {
			formReader = &readerWithProgress{reader: formReader, fsize: formBodyLen, onProgress: progress.onProgress}
		}
		return formReader, nil
	}
	getBodyReadCloser := func() (io.ReadCloser, error) {
		reader, err := getBodyReader()
		if err != nil {
			return nil, err
		}
		return ioutil.NopCloser(reader), nil
	}

	var err error
	var hostProvider hostprovider.HostProvider = nil
	if extra.UpHost != "" {
		hostProvider = hostprovider.NewWithHosts([]string{extra.getUpHost(p.Cfg.UseHTTPS)})
	} else {
		hostProvider, err = p.getUpHostProviderFromUploadToken(upToken, extra)
		if err != nil {
			return err
		}
	}

	// 上传
	contentType := formWriter.FormDataContentType()
	headers := http.Header{}
	headers.Add("Content-Type", contentType)
	err = doUploadAction(hostProvider, extra.TryTimes, extra.HostFreezeDuration, func(host string) error {
		reader, gErr := getBodyReader()
		if gErr != nil {
			return gErr
		}

		return p.Client.CallWithBodyGetter(ctx, ret, "POST", host, headers, reader, getBodyReadCloser, formBodyLen)
	})
	if err != nil {
		return err
	}
	if extra.OnProgress != nil {
		extra.OnProgress(formBodyLen, formBodyLen)
	}

	return nil
}

func (p *FormUploader) getUpHostProviderFromUploadToken(upToken string, extra *PutExtra) (hostprovider.HostProvider, error) {
	ak, bucket, err := getAkBucketFromUploadToken(upToken)
	if err != nil {
		return nil, err
	}
	return getUpHostProvider(p.Cfg, extra.TryTimes, extra.HostFreezeDuration, ak, bucket)
}

type crc32Reader struct {
	h                hash.Hash32
	boundary         string
	r                io.Reader
	inited           bool
	nlDashBoundaryNl string
	header           string
	crc32PadLen      int64
}

func newCrc32Reader(boundary string, h hash.Hash32) *crc32Reader {
	nlDashBoundaryNl := fmt.Sprintf("\r\n--%s\r\n", boundary)
	header := `Content-Disposition: form-data; name="crc32"` + "\r\n\r\n"
	return &crc32Reader{
		h:                h,
		boundary:         boundary,
		nlDashBoundaryNl: nlDashBoundaryNl,
		header:           header,
		crc32PadLen:      10,
	}
}

func (r *crc32Reader) Read(p []byte) (int, error) {
	if !r.inited {
		crc32Sum := r.h.Sum32()
		crc32Line := r.nlDashBoundaryNl + r.header + fmt.Sprintf("%010d", crc32Sum) //padding crc32 results to 10 digits
		r.r = strings.NewReader(crc32Line)
		r.inited = true
	}
	return r.r.Read(p)
}

func (r crc32Reader) length() (length int64) {
	return int64(len(r.nlDashBoundaryNl+r.header)) + r.crc32PadLen
}

func (p *FormUploader) UpHost(ak, bucket string) (upHost string, err error) {
	return getUpHost(p.Cfg, 0, 0, ak, bucket)
}

type readerWithProgress struct {
	reader     io.Reader
	uploaded   int64
	fsize      int64
	onProgress func(fsize, uploaded int64)
}

func (p *readerWithProgress) Read(b []byte) (n int, err error) {
	if p.uploaded > 0 {
		p.onProgress(p.fsize, p.uploaded)
	}

	n, err = p.reader.Read(b)
	p.uploaded += int64(n)
	if p.fsize > 0 && p.uploaded > p.fsize {
		p.uploaded = p.fsize
	}
	return
}

func writeMultipart(writer *multipart.Writer, uptoken, key string, hasKey bool,
	extra *PutExtra, fileName string) (err error) {

	//token
	if err = writer.WriteField("token", uptoken); err != nil {
		return
	}

	//key
	if hasKey {
		if err = writer.WriteField("key", key); err != nil {
			return
		}
	}

	//extra.Params
	if extra.Params != nil {
		for k, v := range extra.Params {
			if (strings.HasPrefix(k, "x:") || strings.HasPrefix(k, "x-qn-meta-")) && v != "" {
				err = writer.WriteField(k, v)
				if err != nil {
					return
				}
			}
		}
	}

	return err
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

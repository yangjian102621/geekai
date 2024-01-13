package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/qiniu/go-sdk/v7/internal/hostprovider"
	"hash/crc32"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/client"
)

// Base64Uploader 表示一个Base64上传对象
type Base64Uploader struct {
	client *client.Client
	cfg    *Config
}

// NewBase64Uploader 用来构建一个Base64上传的对象
func NewBase64Uploader(cfg *Config) *Base64Uploader {
	if cfg == nil {
		cfg = &Config{}
	}

	return &Base64Uploader{
		client: &client.DefaultClient,
		cfg:    cfg,
	}
}

// NewBase64UploaderEx 用来构建一个Base64上传的对象
func NewBase64UploaderEx(cfg *Config, clt *client.Client) *Base64Uploader {
	if cfg == nil {
		cfg = &Config{}
	}

	if clt == nil {
		clt = &client.DefaultClient
	}

	return &Base64Uploader{
		client: clt,
		cfg:    cfg,
	}
}

// Base64PutExtra 为Base64上传的额外可选项
type Base64PutExtra struct {
	// 可选，用户自定义参数，必须以 "x:" 开头。若不以x:开头，则忽略。
	Params map[string]string

	// 可选，当为 "" 时候，服务端自动判断。
	MimeType string

	TryTimes int // 可选。尝试次数

	// 主备域名冻结时间（默认：600s），当一个域名请求失败（单个域名会被重试 TryTimes 次），会被冻结一段时间，使用备用域名进行重试，在冻结时间内，域名不能被使用，当一个操作中所有域名竣备冻结操作不在进行重试，返回最后一次操作的错误。
	HostFreezeDuration time.Duration
}

func (r *Base64PutExtra) init() {
	if r.TryTimes == 0 {
		r.TryTimes = settings.TryTimes
	}
	if r.HostFreezeDuration <= 0 {
		r.HostFreezeDuration = 10 * 60 * time.Second
	}
}

// Put 用来以Base64方式上传一个文件
//
// ctx        是请求的上下文。
// ret        是上传成功后返回的数据。如果 uptoken 中没有设置 callbackUrl 或 returnBody，那么返回的数据结构是 PutRet 结构。
// uptoken    是由业务服务器颁发的上传凭证。
// key        是要上传的文件访问路径。比如："foo/bar.jpg"。注意我们建议 key 不要以 '/' 开头。另外，key 为空字符串是合法的。
// base64Data 是要上传的Base64数据，一般为图片数据的Base64编码字符串
// extra      是上传的一些可选项，可以指定为nil。详细见 Base64PutExtra 结构的描述。
func (p *Base64Uploader) Put(
	ctx context.Context, ret interface{}, uptoken, key string, base64Data []byte, extra *Base64PutExtra) (err error) {
	return p.put(ctx, ret, uptoken, key, true, base64Data, extra)
}

// PutWithoutKey 用来以Base64方式上传一个文件，保存的文件名以文件的内容hash作为文件名
func (p *Base64Uploader) PutWithoutKey(
	ctx context.Context, ret interface{}, uptoken string, base64Data []byte, extra *Base64PutExtra) (err error) {
	return p.put(ctx, ret, uptoken, "", false, base64Data, extra)
}

func (p *Base64Uploader) put(
	ctx context.Context, ret interface{}, uptoken, key string, hasKey bool, base64Data []byte, extra *Base64PutExtra) (err error) {

	//set default extra
	if extra == nil {
		extra = &Base64PutExtra{}
	}
	extra.init()

	//calc crc32
	h := crc32.NewIEEE()
	rawReader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(base64Data))
	fsize, decodeErr := io.Copy(h, rawReader)
	if decodeErr != nil {
		err = fmt.Errorf("invalid base64 data, %s", decodeErr.Error())
		return
	}
	fCrc32 := h.Sum32()

	postPath := bytes.NewBufferString("/putb64")
	//add fsize
	postPath.WriteString("/")
	postPath.WriteString(strconv.Itoa(int(fsize)))

	//add key
	if hasKey {
		postPath.WriteString("/key/")
		postPath.WriteString(base64.URLEncoding.EncodeToString([]byte(key)))
	}
	//add mimeType
	if extra.MimeType != "" {
		postPath.WriteString("/mimeType/")
		postPath.WriteString(base64.URLEncoding.EncodeToString([]byte(extra.MimeType)))
	}

	//add crc32
	postPath.WriteString("/crc32/")
	postPath.WriteString(fmt.Sprintf("%d", fCrc32))

	//add extra params
	if len(extra.Params) > 0 {
		for k, v := range extra.Params {
			if strings.HasPrefix(k, "x:") && v != "" {
				postPath.WriteString("/")
				postPath.WriteString(k)
				postPath.WriteString("/")
				postPath.WriteString(base64.URLEncoding.EncodeToString([]byte(v)))
			}
		}
	}

	//get up host
	ak, bucket, gErr := getAkBucketFromUploadToken(uptoken)
	if gErr != nil {
		err = gErr
		return
	}

	var upHostProvider hostprovider.HostProvider
	upHostProvider, err = p.upHostProvider(ak, bucket, extra)
	if err != nil {
		return
	}

	return doUploadAction(upHostProvider, extra.TryTimes, extra.HostFreezeDuration, func(host string) error {
		postURL := fmt.Sprintf("%s%s", host, postPath.String())
		headers := http.Header{}
		headers.Add("Content-Type", "application/octet-stream")
		headers.Add("Authorization", "UpToken "+uptoken)

		return p.client.CallWith(ctx, ret, "POST", postURL, headers, bytes.NewReader(base64Data), len(base64Data))
	})
}

func (p *Base64Uploader) upHostProvider(ak, bucket string, extra *Base64PutExtra) (hostProvider hostprovider.HostProvider, err error) {
	return getUpHostProvider(p.cfg, extra.TryTimes, extra.HostFreezeDuration, ak, bucket)
}

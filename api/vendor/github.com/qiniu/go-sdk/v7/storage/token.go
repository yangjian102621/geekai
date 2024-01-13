package storage

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
)

// PutPolicy 表示文件上传的上传策略，参考 https://developer.qiniu.com/kodo/manual/1206/put-policy
type PutPolicy struct {

	// 指定上传的目标资源空间 Bucket 和资源键 Key（最大为 750 字节）。有三种格式：
	// <bucket>，表示允许用户上传文件到指定的 bucket。在这种格式下文件只能新增（分片上传 v1 版 需要指定 insertOnly 为 1 才是新增，否则也为覆盖上传），若已存在同名资源（且文件内容/etag不一致），上传会失败；若已存在资源的内容/etag一致，则上传会返回成功。
	// <bucket>:<key>，表示只允许用户上传指定 key 的文件。在这种格式下文件默认允许修改，若已存在同名资源则会被覆盖。如果只希望上传指定 key 的文件，并且不允许修改，那么可以将下面的 insertOnly 属性值设为 1。
	// <bucket>:<keyPrefix>，表示只允许用户上传指定以 keyPrefix 为前缀的文件，当且仅当 isPrefixalScope 字段为 1 时生效，isPrefixalScope 为 1 时无法覆盖上传。
	Scope string `json:"scope"`

	// 若为 1，表示允许用户上传以 scope 的 keyPrefix 为前缀的文件。
	IsPrefixalScope int `json:"isPrefixalScope,omitempty"`

	// 上传凭证有效截止时间。Unix时间戳，单位为秒。该截止时间为上传完成后，在七牛空间生成文件的校验时间，而非上传的开始时间，
	// 一般建议设置为上传开始时间 + 3600s，用户可根据具体的业务场景对凭证截止时间进行调整。
	Expires uint64 `json:"deadline"`

	// 若非0, 即使Scope为 Bucket:Key 的形式也是insert only
	InsertOnly uint16 `json:"insertOnly,omitempty"`

	// 唯一属主标识。特殊场景下非常有用，例如根据 App-Client 标识给图片或视频打水印。
	EndUser string `json:"endUser,omitempty"`

	// Web 端文件上传成功后，浏览器执行 303 跳转的 URL。通常用于表单上传。
	// 文件上传成功后会跳转到 <returnUrl>?upload_ret=<queryString>，<queryString>包含 returnBody 内容。
	// 如不设置 returnUrl，则直接将 returnBody 的内容返回给客户端。
	ReturnURL string `json:"returnUrl,omitempty"`

	// 上传成功后，自定义七牛云最终返回給上传端（在指定 returnUrl 时是携带在跳转路径参数中）的数据。支持魔法变量和自定义变量。
	// returnBody 要求是合法的 JSON 文本。
	// 例如 {“key”: $(key), “hash”: $(etag), “w”: $(imageInfo.width), “h”: $(imageInfo.height)}。
	ReturnBody string `json:"returnBody,omitempty"`

	// 上传成功后，七牛云向业务服务器发送 POST 请求的 URL。必须是公网上可以正常进行 POST 请求并能响应 HTTP/1.1 200 OK 的有效 URL。
	// 另外，为了给客户端有一致的体验，我们要求 callbackUrl 返回包 Content-Type 为 “application/json”，即返回的内容必须是合法的
	// JSON 文本。出于高可用的考虑，本字段允许设置多个 callbackUrl（用英文符号 ; 分隔），在前一个 callbackUrl 请求失败的时候会依次
	// 重试下一个 callbackUrl。一个典型例子是：http://<ip1>/callback;http://<ip2>/callback，并同时指定下面的 callbackHost 字段。
	// 在 callbackUrl 中使用 ip 的好处是减少对 dns 解析的依赖，可改善回调的性能和稳定性。指定 callbackUrl，必须指定 callbackbody，
	// 且值不能为空。
	CallbackURL string `json:"callbackUrl,omitempty"`

	// 上传成功后，七牛云向业务服务器发送回调通知时的 Host 值。与 callbackUrl 配合使用，仅当设置了 callbackUrl 时才有效。
	CallbackHost string `json:"callbackHost,omitempty"`

	// 上传成功后，七牛云向业务服务器发送 Content-Type: application/x-www-form-urlencoded 的 POST 请求。业务服务器可以通过直接读取
	// 请求的 query 来获得该字段，支持魔法变量和自定义变量。callbackBody 要求是合法的 url query string。
	// 例如key=$(key)&hash=$(etag)&w=$(imageInfo.width)&h=$(imageInfo.height)。如果callbackBodyType指定为application/json，
	// 则callbackBody应为json格式，例如:{“key”:"$(key)",“hash”:"$(etag)",“w”:"$(imageInfo.width)",“h”:"$(imageInfo.height)"}。
	CallbackBody string `json:"callbackBody,omitempty"`

	// 上传成功后，七牛云向业务服务器发送回调通知 callbackBody 的 Content-Type。默认为 application/x-www-form-urlencoded，也可设置
	// 为 application/json。
	CallbackBodyType string `json:"callbackBodyType,omitempty"`

	// 资源上传成功后触发执行的预转持久化处理指令列表。fileType=2或3（上传归档存储或深度归档存储文件）时，不支持使用该参数。支持魔法变量和自
	// 定义变量。每个指令是一个 API 规格字符串，多个指令用;分隔。请参阅persistenOps详解与示例。同时添加 persistentPipeline 字段，使用专
	// 用队列处理，请参阅persistentPipeline。
	PersistentOps string `json:"persistentOps,omitempty"`

	// 接收持久化处理结果通知的 URL。必须是公网上可以正常进行 POST 请求并能响应 HTTP/1.1 200 OK 的有效 URL。该 URL 获取的内容和持久化处
	// 理状态查询的处理结果一致。发送 body 格式是 Content-Type 为 application/json 的 POST 请求，需要按照读取流的形式读取请求的 body
	// 才能获取。
	PersistentNotifyURL string `json:"persistentNotifyUrl,omitempty"`

	// 转码队列名。资源上传成功后，触发转码时指定独立的队列进行转码。为空则表示使用公用队列，处理速度比较慢。建议使用专用队列。
	PersistentPipeline string `json:"persistentPipeline,omitempty"`

	// saveKey 的优先级设置。为 true 时，saveKey不能为空，会忽略客户端指定的key，强制使用saveKey进行文件命名。参数不设置时，
	// 默认值为false
	ForceSaveKey bool `json:"forceSaveKey,omitempty"` //

	// 自定义资源名。支持魔法变量和自定义变量。forceSaveKey 为false时，这个字段仅当用户上传的时候没有主动指定 key 时起作用；
	// forceSaveKey 为true时，将强制按这个字段的格式命名。
	SaveKey string `json:"saveKey,omitempty"`

	// 限定上传文件大小最小值，单位Byte。小于限制上传文件大小的最小值会被判为上传失败，返回 403 状态码
	FsizeMin int64 `json:"fsizeMin,omitempty"`

	// 限定上传文件大小最大值，单位Byte。超过限制上传文件大小的最大值会被判为上传失败，返回 413 状态码。
	FsizeLimit int64 `json:"fsizeLimit,omitempty"`

	// 开启 MimeType 侦测功能，并按照下述规则进行侦测；如不能侦测出正确的值，会默认使用 application/octet-stream 。
	// 设为非 0 值，则忽略上传端传递的文件 MimeType 信息，并按如下顺序侦测 MimeType 值：
	// 1. 侦测内容； 2. 检查文件扩展名； 3. 检查 Key 扩展名。
	// 默认设为 0 值，如上传端指定了 MimeType 则直接使用该值，否则按如下顺序侦测 MimeType 值：
	// 1. 检查文件扩展名； 2. 检查 Key 扩展名； 3. 侦测内容。
	DetectMime uint8 `json:"detectMime,omitempty"`

	// 限定用户上传的文件类型。指定本字段值，七牛服务器会侦测文件内容以判断 MimeType，再用判断值跟指定值进行匹配，匹配成功则允许上传，匹配失败则返回 403 状态码。示例：
	// image/* 表示只允许上传图片类型
	// image/jpeg;image/png 表示只允许上传 jpg 和 png 类型的图片
	// !application/json;text/plain 表示禁止上传 json 文本和纯文本。注意最前面的感叹号！
	MimeLimit string `json:"mimeLimit,omitempty"`

	// 资源的存储类型，0表示标准存储，1 表示低频存储，2 表示归档存储，3 表示深度归档存储。
	FileType int `json:"fileType,omitempty"`

	CallbackFetchKey uint8 `json:"callbackFetchKey,omitempty"`

	DeleteAfterDays int `json:"deleteAfterDays,omitempty"`
}

// UploadToken 方法用来进行上传凭证的生成
// 该方法生成的过期时间是现对于现在的时间
func (p *PutPolicy) UploadToken(cred *auth.Credentials) string {
	return p.uploadToken(cred)
}

func (p PutPolicy) uploadToken(cred *auth.Credentials) (token string) {
	if p.Expires == 0 {
		p.Expires = 3600 // 默认一小时过期
	}
	p.Expires += uint64(time.Now().Unix())
	putPolicyJSON, _ := json.Marshal(p)
	token = cred.SignWithData(putPolicyJSON)
	return
}

func getAkBucketFromUploadToken(token string) (ak, bucket string, err error) {
	items := strings.Split(token, ":")
	// KODO-11919
	if len(items) == 5 && items[0] == "" {
		items = items[2:]
	} else if len(items) != 3 {
		err = errors.New("invalid upload token, format error")
		return
	}

	ak = items[0]
	policyBytes, dErr := base64.URLEncoding.DecodeString(items[2])
	if dErr != nil {
		err = errors.New("invalid upload token, invalid put policy")
		return
	}

	putPolicy := PutPolicy{}
	uErr := json.Unmarshal(policyBytes, &putPolicy)
	if uErr != nil {
		err = errors.New("invalid upload token, invalid put policy")
		return
	}

	bucket = strings.Split(putPolicy.Scope, ":")[0]
	return
}

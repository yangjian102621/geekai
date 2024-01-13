// package storage 提供了用户存储配置(uc)方面的功能, 定义了UC API 的返回结构体类型
package storage

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/internal/clientv2"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
)

// BucketSummary 存储空间信息
type BucketSummary struct {
	// 存储空间名字
	Name string     `json:"name"`
	Info BucketInfo `json:"info"`
}

// BucketInfo 存储空间的详细信息
type BucketInfo struct {
	// 镜像回源地址， 接口返回的多个地址以；分割
	Source string `json:"source"`

	// 镜像回源的时候请求头中的HOST
	Host string `json:"host"`

	// 镜像回源地址过期时间(秒数)， 现在这个功能没有实现，因此这个字段现在是没有意义的
	Expires int `json:"expires"`

	// 是否开启了原图保护
	Protected int `json:"protected"`

	// 是否是私有空间
	Private int `json:"private"`

	// 如果NoIndexPage是false表示开启了空间根目录index.html
	// 如果是true, 表示没有开启
	// 开启了根目录下的index.html, 文件将会被作为默认首页展示
	NoIndexPage int `json:"no_index_page"`

	// 在规定的时效内使客户端缓存更新的效果
	MaxAge int `json:"max_age"`

	// 图片样式分隔符， 接口返回的可能有多个
	Separator string `json:"separator"`

	// 图片样式， map中的key表示图片样式命令名字
	// map中的value表示图片样式命令的内容
	Styles map[string]string `json:"styles"`

	// 防盗链模式
	// 1 - 表示设置了防盗链的referer白名单
	// 2 - 表示设置了防盗链的referer黑名单
	AntiLeechMode int `json:"anti_leech_mode"`

	// 使用token签名进行防盗链
	// 0 - 表示关闭
	// 1 - 表示开启
	TokenAntiLeechMode int `json:"token_anti_leech"`

	// 防盗链referer白名单列表
	ReferWl []string `json:"refer_wl"`

	// 防盗链referer黑名单列表
	ReferBl []string `json:"refer_bl"`

	// 在源站支持的情况下开启源站的Referer防盗链
	EnableSource bool `json:"source_enabled"`

	// 是否允许空的referer访问
	NoRefer bool `json:"no_refer"`

	// 用于防盗链token的生成
	MacKey string `json:"mac_key"`

	// 用于防盗链token的生成
	MacKey2 string `json:"mac_key2"`

	// 存储区域， 兼容保留
	Zone string

	// 存储区域
	Region string

	// 空间备注信息
	Remark string

	// 空间创建时间
	Ctime time.Time
}

// ReferAntiLeechConfig 是用户存储空间的Refer防盗链配置
type ReferAntiLeechConfig struct {
	// 防盗链模式， 0 - 关闭Refer防盗链, 1 - 开启Referer白名单，2 - 开启Referer黑名单
	Mode int

	// 是否允许空的referer访问
	AllowEmptyReferer bool

	// Pattern 匹配HTTP Referer头, 当模式是1或者2的时候有效
	// Mode为1的时候表示允许Referer符合该Pattern的HTTP请求访问
	// Mode为2的时候表示禁止Referer符合该Pattern的HTTP请求访问
	// 当前允许的匹配字符串格式分为三种:
	// 一种为空主机头域名, 比如 foo.com; 一种是泛域名, 比如 *.bar.com;
	// 一种是完全通配符, 即一个 *;
	// 多个规则之间用;隔开, 比如: foo.com;*.bar.com;sub.foo.com;*.sub.bar.com
	Pattern string

	// 是否开启源站的防盗链， 默认为0， 只开启CDN防盗链， 当设置为1的时候
	// 在源站支持的情况下开启源站的Referer防盗链
	EnableSource bool
}

// SetMode 设置referer防盗链模式
func (r *ReferAntiLeechConfig) SetMode(mode int) *ReferAntiLeechConfig {
	if mode != 0 && mode != 1 && mode != 2 {
		panic("Referer anti_leech_mode must be in [0, 1, 2]")
	}
	r.Mode = mode
	return r
}

// SetEmptyReferer 设置是否允许空Referer访问
func (r *ReferAntiLeechConfig) SetEmptyReferer(enable bool) *ReferAntiLeechConfig {
	r.AllowEmptyReferer = enable
	return r
}

// SetPattern 设置匹配Referer的模式
func (r *ReferAntiLeechConfig) SetPattern(pattern string) *ReferAntiLeechConfig {
	if pattern == "" {
		panic("Empty pattern is not allowed")
	}

	r.Pattern = pattern
	return r
}

// AddDomainPattern 添加pattern到Pattern字段
// 假入Pattern值为"*.qiniu.com"， 使用AddDomainPattern("*.baidu.com")后
// r.Pattern的值为"*.qiniu.com;*.baidu.com"
func (r *ReferAntiLeechConfig) AddDomainPattern(pattern string) *ReferAntiLeechConfig {
	if strings.HasSuffix(r.Pattern, ";") {
		r.Pattern = strings.TrimRight(r.Pattern, ";")
	}
	r.Pattern = strings.Join([]string{r.Pattern, pattern}, ";")
	return r
}

// SetEnableSource 设置是否开启源站的防盗链
func (r *ReferAntiLeechConfig) SetEnableSource(enable bool) *ReferAntiLeechConfig {
	r.EnableSource = enable
	return r
}

// AsQueryString 编码成query参数格式
func (r *ReferAntiLeechConfig) AsQueryString() string {
	params := make(url.Values, 4)
	params.Add("mode", strconv.Itoa(r.Mode))
	params.Add("pattern", r.Pattern)
	if r.AllowEmptyReferer {
		params.Add("norefer", "1")
	} else {
		params.Add("norefer", "0")
	}
	if r.EnableSource {
		params.Add("source_enabled", "1")
	} else {
		params.Add("source_enabled", "0")
	}
	return params.Encode()
}

// ProtectedOn 返回true or false
// 如果开启了原图保护，返回true, 否则false
func (b *BucketInfo) ProtectedOn() bool {
	return b.Protected == 1
}

// IsPrivate  返回布尔值
// 如果是私有空间， 返回 true, 否则返回false
func (b *BucketInfo) IsPrivate() bool {
	return b.Private == 1
}

// ImageSources 返回多个镜像回源地址的列表
func (b *BucketInfo) ImageSources() (srcs []string) {
	srcs = strings.Split(b.Source, ";")
	return
}

// IndexPageOn 返回空间是否开启了根目录下的index.html
func (b *BucketInfo) IndexPageOn() bool {
	return b.NoIndexPage == 0
}

// Separators 返回分隔符列表
func (b *BucketInfo) Separators() (ret []rune) {
	for _, r := range b.Separator {
		ret = append(ret, r)
	}
	return
}

// WhiteListSet 是否设置了防盗链白名单
func (b *BucketInfo) WhiteListSet() bool {
	return b.AntiLeechMode == 1
}

// BlackListSet 是否设置了防盗链黑名单
func (b *BucketInfo) BlackListSet() bool {
	return b.AntiLeechMode == 2
}

// TokenAntiLeechModeOn 返回是否使用token签名防盗链开启了
func (b *BucketInfo) TokenAntiLeechModeOn() bool {
	return b.TokenAntiLeechMode == 1
}

// GetBucketInfo 返回BucketInfo结构
func (m *BucketManager) GetBucketInfo(bucketName string) (bucketInfo BucketInfo, err error) {
	reqURL := fmt.Sprintf("%s/v2/bucketInfo?bucket=%s", getUcHost(m.Cfg.UseHTTPS), bucketName)
	_, err = clientv2.DoAndDecodeJsonResponse(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	}, &bucketInfo)
	return bucketInfo, err
}

// SetRemark 设置空间备注信息
func (m *BucketManager) SetRemark(bucketName, remark string) (err error) {
	reqURL := fmt.Sprintf("%s/buckets/%s?remark", getUcHost(m.Cfg.UseHTTPS), bucketName)
	body := struct {
		Remark string `json:"remark"`
	}{Remark: remark}
	_, err = clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPut,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: clientv2.RequestBodyCreatorOfJson(body),
	})
	return err
}

// BucketInfosForRegion 获取指定区域的该用户的所有bucketInfo信息
func (m *BucketManager) BucketInfosInRegion(region RegionID, statistics bool) (bucketInfos []BucketSummary, err error) {
	reqURL := fmt.Sprintf("%s/v2/bucketInfos?region=%s&fs=%t", getUcHost(m.Cfg.UseHTTPS), string(region), statistics)
	_, err = clientv2.DoAndDecodeJsonResponse(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	}, &bucketInfos)
	return bucketInfos, err
}

// SetReferAntiLeechMode 配置存储空间referer防盗链模式
func (m *BucketManager) SetReferAntiLeechMode(bucketName string, refererAntiLeechConfig *ReferAntiLeechConfig) (err error) {
	reqURL := fmt.Sprintf("%s/referAntiLeech?bucket=%s&%s", getUcHost(m.Cfg.UseHTTPS), bucketName, refererAntiLeechConfig.AsQueryString())
	err = m.Client.CredentialedCall(context.Background(), m.Mac, auth.TokenQiniu, nil, "POST", reqURL, nil)
	return
}

// BucketLifeCycleRule 定义了关于七牛存储空间关于生命周期的一些配置，规则。
// 比如存储空间中文件可以设置多少天后删除，多少天后转低频存储等等
type BucketLifeCycleRule struct {
	// 规则名称， 在设置的bucket中规则名称需要是唯一的
	// 同时长度小于50， 不能为空
	// 由字母，数字和下划线组成
	Name string `json:"name"`

	// 以该前缀开头的文件应用此规则
	Prefix string `json:"prefix"`

	// 指定存储空间内的文件多少天后删除
	// 0 - 不删除
	// > 0 表示多少天后删除
	DeleteAfterDays int `json:"delete_after_days"`

	// 在多少天后转低频存储
	// 0  - 表示不转低频
	// < 0 表示上传的文件立即使用低频存储
	// > 0 表示多少天后转低频存储
	ToLineAfterDays int `json:"to_line_after_days"`

	// 指定文件上传多少天后转归档存储。
	// 0 表示不转归档存储，
	// < 0 表示上传的文件立即变归档存储
	// > 0 表示多少天后转归档存储
	ToArchiveAfterDays int `json:"to_archive_after_days"`

	// 指定文件上传多少天后转深度归档存储。
	// < 0 表示上传的文件立即变深度归档存储
	// > 0 表示多少天后转深度归档存储
	ToDeepArchiveAfterDays int `json:"to_deep_archive_after_days"`
}

// SetBucketLifeCycleRule 设置存储空间内文件的生命周期规则
func (m *BucketManager) AddBucketLifeCycleRule(bucketName string, lifeCycleRule *BucketLifeCycleRule) (err error) {
	params := make(map[string][]string)

	// 没有检查参数的合法性，交给服务端检查
	params["bucket"] = []string{bucketName}
	params["name"] = []string{lifeCycleRule.Name}
	params["prefix"] = []string{lifeCycleRule.Prefix}
	params["delete_after_days"] = []string{strconv.Itoa(lifeCycleRule.DeleteAfterDays)}
	params["to_ia_after_days"] = []string{strconv.Itoa(lifeCycleRule.ToLineAfterDays)}
	params["to_archive_after_days"] = []string{strconv.Itoa(lifeCycleRule.ToArchiveAfterDays)}
	params["to_deep_archive_after_days"] = []string{strconv.Itoa(lifeCycleRule.ToDeepArchiveAfterDays)}

	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/rules/add"
	_, err = clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: clientv2.RequestBodyCreatorForm(params),
	})
	return err
}

// DelBucketLifeCycleRule 删除特定存储空间上设定的规则
func (m *BucketManager) DelBucketLifeCycleRule(bucketName, ruleName string) (err error) {
	params := make(map[string][]string)

	params["bucket"] = []string{bucketName}
	params["name"] = []string{ruleName}

	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/rules/delete"
	_, err = clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: clientv2.RequestBodyCreatorForm(params),
	})
	return err
}

// UpdateBucketLifeCycleRule 更新特定存储空间上的生命周期规则
func (m *BucketManager) UpdateBucketLifeCycleRule(bucketName string, rule *BucketLifeCycleRule) (err error) {
	params := make(map[string][]string)

	params["bucket"] = []string{bucketName}
	params["name"] = []string{rule.Name}
	params["prefix"] = []string{rule.Prefix}
	params["delete_after_days"] = []string{strconv.Itoa(rule.DeleteAfterDays)}
	params["to_line_after_days"] = []string{strconv.Itoa(rule.ToLineAfterDays)}
	params["to_archive_after_days"] = []string{strconv.Itoa(rule.ToArchiveAfterDays)}
	params["to_deep_archive_after_days"] = []string{strconv.Itoa(rule.ToDeepArchiveAfterDays)}

	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/rules/update"
	_, err = clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: clientv2.RequestBodyCreatorForm(params),
	})
	return err
}

// GetBucketLifeCycleRule 获取指定空间上设置的生命周期规则
func (m *BucketManager) GetBucketLifeCycleRule(bucketName string) (rules []BucketLifeCycleRule, err error) {
	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/rules/get?bucket=" + bucketName
	_, err = clientv2.DoAndDecodeJsonResponse(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodGet,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	}, &rules)
	return rules, err
}

// BucketEnvent 定义了存储空间发生事件时候的通知规则
// 比如调用了存储的"delete"删除接口删除文件， 这个是一个事件；
// 当这个事件发生的时候， 我们要对哪些文件，做什么处理，是否要作回调，
// 都可以通过这个结构体配置
type BucketEventRule struct {
	// 规则名字
	Name string `json:"name"`

	// 匹配文件前缀
	Prefix string `json:"prefix"`

	// 匹配文件后缀
	Suffix string `json:"suffix"`

	// 事件类型
	// put,mkfile,delete,copy,move,append,disable,enable,deleteMarkerCreate
	Event []string `json:"event"`

	// 回调通知地址， 可以指定多个
	CallbackURL []string `json:"callback_urls"`

	// 用户的AccessKey， 可选， 设置的话会对通知请求用对应的ak、sk进行签名
	AccessKey string `json:"access_key"`

	// 回调通知的请求HOST, 可选
	Host string `json:"host"`
}

// Params 返回一个hash结构
func (r *BucketEventRule) Params(bucket string) map[string][]string {
	params := make(map[string][]string)

	params["bucket"] = []string{bucket}
	params["name"] = []string{r.Name}
	if r.Prefix != "" {
		params["prefix"] = []string{r.Prefix}
	}
	if r.Suffix != "" {
		params["suffix"] = []string{r.Suffix}
	}
	params["event"] = r.Event
	params["callbackURL"] = r.CallbackURL
	if r.AccessKey != "" {
		params["access_key"] = []string{r.AccessKey}
	}
	if r.Host != "" {
		params["host"] = []string{r.Host}
	}
	return params
}

// AddBucketEvent 增加存储空间事件通知规则
func (m *BucketManager) AddBucketEvent(bucket string, rule *BucketEventRule) (err error) {
	params := rule.Params(bucket)
	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/events/add"
	_, err = clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: clientv2.RequestBodyCreatorForm(params),
	})
	return err
}

// DelBucketEvent 删除指定存储空间的通知事件规则
func (m *BucketManager) DelBucketEvent(bucket, ruleName string) (err error) {
	params := make(map[string][]string)
	params["bucket"] = []string{bucket}
	params["name"] = []string{ruleName}

	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/events/delete"
	_, err = clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: clientv2.RequestBodyCreatorForm(params),
	})
	return err
}

// UpdateBucketEnvent 更新指定存储空间的事件通知规则
func (m *BucketManager) UpdateBucketEnvent(bucket string, rule *BucketEventRule) (err error) {
	params := rule.Params(bucket)
	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/events/update"
	_, err = clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: clientv2.RequestBodyCreatorForm(params),
	})
	return err
}

// GetBucketEvent 获取指定存储空间的事件通知规则
func (m *BucketManager) GetBucketEvent(bucket string) (rule []BucketEventRule, err error) {
	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/events/get?bucket=" + bucket
	_, err = clientv2.DoAndDecodeJsonResponse(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodGet,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	}, &rule)
	return rule, err
}

// CorsRule 是关于存储的跨域规则
// 最多允许设置10条跨域规则
// 对于同一个域名如果设置了多条规则，那么按顺序使用第一条匹配的规则去生成返回值。
// 对于简单跨域请求，只匹配 Origin；
// 对于预检请求， 需要匹配 Origin、AllowedMethod、AllowedHeader；

// 如果没有设置任何corsRules，那么默认允许所有的跨域请求
// 参考： https://www.w3.org/TR/cors/
type CorsRule struct {

	// allowed_orgin: 允许的域名。必填；支持通配符*；*表示全部匹配；只有第一个*生效；需要设置"Scheme"；大小写敏感。例如
	//   规则：http://*.abc.*.com 请求："http://test.abc.test.com" 结果：不通过
	//   规则："http://abc.com" 请求："https://abc.com"/"abc.com" 结果：不通过
	//   规则："abc.com" 请求："http://abc.com" 结果：不通过
	AllowedOrigin []string `json:"allowed_origin"`

	// allowed_method: 允许的方法。必填；不支持通配符；大小写不敏感；
	AllowedMethod []string `json:"allowed_method"`

	// allowed_header: 允许的header。选填；支持通配符*，但只能是单独的*，表示允许全部header，其他*不生效；空则不允许任何header；大小写不敏感；
	AllowedHeader []string `json:"allowed_header"`

	// 暴露的header。选填；不支持通配符；X-Log, X-Reqid是默认会暴露的两个header；其他的header如果没有设置，则不会暴露；大小写不敏感；
	ExposedHeader []string `json:"exposed_header"`

	// max_age: 结果可以缓存的时间。选填；空则不缓存
	MaxAge int64 `json:"max_age"`
}

// AddCorsRules 设置指定存储空间的跨域规则
func (m *BucketManager) AddCorsRules(bucket string, corsRules []CorsRule) (err error) {
	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/corsRules/set/" + bucket
	_, err = clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: clientv2.RequestBodyCreatorOfJson(corsRules),
	})
	return err
}

// GetCorsRules 获取指定存储空间的跨域规则
func (m *BucketManager) GetCorsRules(bucket string) (corsRules []CorsRule, err error) {
	reqURL := getUcHost(m.Cfg.UseHTTPS) + "/corsRules/get/" + bucket
	_, err = clientv2.DoAndDecodeJsonResponse(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodGet,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	}, &corsRules)
	return corsRules, err
}

// BucketQuota 七牛存储空间的配额信息
type BucketQuota struct {
	// 如果HTTP请求没有发送该参数或者发送的参数是0，表示不更改当前配置
	// 如果是-1， 表示取消限额
	// 一下两个参数都使用于这个逻辑

	// 空间存储量配额信息
	Size int64

	// 空间文件数配置信息
	Count int64
}

// SetBucketQuota 设置存储空间的配额限制
// 配额限制主要是两块， 空间存储量的限制和空间文件数限制
func (m *BucketManager) SetBucketQuota(bucket string, size, count int64) (err error) {
	reqURL := fmt.Sprintf("%s/setbucketquota/%s/size/%d/count/%d", getUcHost(m.Cfg.UseHTTPS), bucket, size, count)
	_, err = clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	})
	return err
}

// GetBucketQuota 获取存储空间的配额信息
func (m *BucketManager) GetBucketQuota(bucket string) (quota BucketQuota, err error) {
	reqURL := fmt.Sprintf("%s/getbucketquota/%s", getUcHost(m.Cfg.UseHTTPS), bucket)
	_, err = clientv2.DoAndDecodeJsonResponse(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	}, &quota)
	return quota, err
}

// SetBucketAccessStyle 可以用来开启或关闭制定存储空间的原图保护
// mode - 1 ==> 开启原图保护
// mode - 0 ==> 关闭原图保护
func (m *BucketManager) SetBucketAccessStyle(bucket string, mode int) error {
	reqURL := fmt.Sprintf("%s/accessMode/%s/mode/%d", getUcHost(m.Cfg.UseHTTPS), bucket, mode)
	_, err := clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	})
	return err
}

// TurnOffBucketProtected 开启指定存储空间的原图保护
func (m *BucketManager) TurnOnBucketProtected(bucket string) error {
	return m.SetBucketAccessStyle(bucket, 1)
}

// TurnOffBucketProtected 关闭指定空间的原图保护
func (m *BucketManager) TurnOffBucketProtected(bucket string) error {
	return m.SetBucketAccessStyle(bucket, 0)
}

// SetBucketMaxAge 设置指定存储空间的MaxAge响应头
// maxAge <= 0时，表示使用默认值31536000
func (m *BucketManager) SetBucketMaxAge(bucket string, maxAge int64) error {
	reqURL := fmt.Sprintf("%s/maxAge?bucket=%s&maxAge=%d", getUcHost(m.Cfg.UseHTTPS), bucket, maxAge)
	_, err := clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	})
	return err
}

// SetBucketAccessMode 设置指定空间的私有属性

// mode - 1 表示设置空间为私有空间， 私有空间访问需要鉴权
// mode - 0 表示设置空间为公开空间
func (m *BucketManager) SetBucketAccessMode(bucket string, mode int) error {
	reqURL := fmt.Sprintf("%s/private?bucket=%s&private=%d", getUcHost(m.Cfg.UseHTTPS), bucket, mode)
	_, err := clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	})
	return err
}

// MakeBucketPublic 设置空间为公有空间
func (m *BucketManager) MakeBucketPublic(bucket string) error {
	return m.SetBucketAccessMode(bucket, 0)
}

// MakeBucketPrivate 设置空间为私有空间
func (m *BucketManager) MakeBucketPrivate(bucket string) error {
	return m.SetBucketAccessMode(bucket, 1)
}

// TurnOnIndexPage 设置默认首页
func (m *BucketManager) TurnOnIndexPage(bucket string) error {
	return m.setIndexPage(bucket, 0)
}

// TurnOnIndexPage 关闭默认首页
func (m *BucketManager) TurnOffIndexPage(bucket string) error {
	return m.setIndexPage(bucket, 1)
}

func (m *BucketManager) setIndexPage(bucket string, noIndexPage int) error {
	reqURL := fmt.Sprintf("%s/noIndexPage?bucket=%s&noIndexPage=%d", getUcHost(m.Cfg.UseHTTPS), bucket, noIndexPage)
	_, err := clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     nil,
		Method:      clientv2.RequestMethodPost,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	})
	return err
}

// BucketTagging 为 Bucket 设置标签
type BucketTagging struct {
	Tags []BucketTag `json:"Tags"`
}

type BucketTag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// SetTagging 设置 Bucket 标签

// SetTagging 该方法为覆盖所有 Bucket 上之前设置的标签，标签 Key 最大 64 字节，Value 最大 128 字节，均不能为空，且区分大小写
// Key 不能以 kodo 为前缀，Key 和 Value 的字符只能为：字母，数字，空格，+，-，=，.，_，:，/，@，不能支持中文
func (m *BucketManager) SetTagging(bucket string, tags map[string]string) error {
	tagging := BucketTagging{Tags: make([]BucketTag, 0, len(tags))}
	for key, value := range tags {
		tagging.Tags = append(tagging.Tags, BucketTag{Key: key, Value: value})
	}

	reqURL := fmt.Sprintf("%s/bucketTagging?bucket=%s", getUcHost(m.Cfg.UseHTTPS), bucket)
	_, err := clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     context.Background(),
		Method:      clientv2.RequestMethodPut,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: clientv2.RequestBodyCreatorOfJson(tagging),
	})
	return err
}

// ClearTagging 清空 Bucket 标签
func (m *BucketManager) ClearTagging(bucket string) error {
	reqURL := fmt.Sprintf("%s/bucketTagging?bucket=%s", getUcHost(m.Cfg.UseHTTPS), bucket)
	_, err := clientv2.Do(m.getUCClient(), clientv2.RequestParams{
		Context:     context.Background(),
		Method:      clientv2.RequestMethodDelete,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	})
	return err
}

// GetTagging 获取 Bucket 标签
func (m *BucketManager) GetTagging(bucket string) (tags map[string]string, err error) {
	var tagging BucketTagging
	reqURL := fmt.Sprintf("%s/bucketTagging?bucket=%s", getUcHost(m.Cfg.UseHTTPS), bucket)
	_, err = clientv2.DoAndDecodeJsonResponse(m.getUCClient(), clientv2.RequestParams{
		Context:     context.Background(),
		Method:      clientv2.RequestMethodGet,
		Url:         reqURL,
		Header:      nil,
		BodyCreator: nil,
	}, &tagging)
	if err != nil {
		return
	}
	tags = make(map[string]string, len(tagging.Tags))
	for _, tag := range tagging.Tags {
		tags[tag.Key] = tag.Value
	}
	return
}

func (m *BucketManager) getUCClient() clientv2.Client {
	return getUCClient(ucClientConfig{
		IsUcQueryApi:       false,
		RetryMax:           m.options.RetryMax,
		HostFreezeDuration: m.options.HostFreezeDuration,
		Client:             m.Client,
	}, m.Mac)
}

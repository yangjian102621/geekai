package storage

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/internal/clientv2"
	"github.com/qiniu/go-sdk/v7/internal/hostprovider"
	"strings"
	"time"
)

// 存储所在的地区，例如华东，华南，华北
// 每个存储区域可能有多个机房信息，每个机房可能有多个上传入口
type Region struct {
	// 上传入口
	SrcUpHosts []string `json:"src_up,omitempty"`

	// 加速上传入口
	CdnUpHosts []string `json:"cdn_up,omitempty"`

	// 获取文件信息入口
	RsHost string `json:"rs,omitempty"`

	// bucket列举入口
	RsfHost string `json:"rsf,omitempty"`

	ApiHost string `json:"api,omitempty"`

	// 存储io 入口
	IovipHost string `json:"io,omitempty"`

	// 源站下载入口
	IoSrcHost string `json:"io_src,omitempty"`
}

type RegionID string

// GetDefaultReion 根据RegionID获取对应的Region信息
func GetRegionByID(regionID RegionID) (Region, bool) {
	if r, ok := regionMap[regionID]; ok {
		return r, ok
	}
	return Region{}, false
}

func (r *Region) String() string {
	str := ""
	str += fmt.Sprintf("SrcUpHosts: %v\n", r.SrcUpHosts)
	str += fmt.Sprintf("CdnUpHosts: %v\n", r.CdnUpHosts)
	str += fmt.Sprintf("IovipHost: %s\n", r.IovipHost)
	str += fmt.Sprintf("RsHost: %s\n", r.RsHost)
	str += fmt.Sprintf("RsfHost: %s\n", r.RsfHost)
	str += fmt.Sprintf("ApiHost: %s\n", r.ApiHost)
	return str
}

func endpoint(useHttps bool, host string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}

	if strings.HasPrefix(host, "http://") ||
		strings.HasPrefix(host, "https://") {
		return host
	}

	scheme := "http://"
	if useHttps {
		scheme = "https://"
	}
	return fmt.Sprintf("%s%s", scheme, host)
}

// 获取rsfHost
func (r *Region) GetRsfHost(useHttps bool) string {
	return endpoint(useHttps, r.RsfHost)
}

// 获取io host
func (r *Region) GetIoHost(useHttps bool) string {
	return endpoint(useHttps, r.IovipHost)
}

// 获取RsHost
func (r *Region) GetRsHost(useHttps bool) string {
	return endpoint(useHttps, r.RsHost)
}

// 获取api host
func (r *Region) GetApiHost(useHttps bool) string {
	return endpoint(useHttps, r.ApiHost)
}

var (
	// regionHuadong 表示华东机房
	regionHuadong = Region{
		SrcUpHosts: []string{
			"up.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload.qiniup.com",
		},
		RsHost:    "rs.qbox.me",
		RsfHost:   "rsf.qbox.me",
		ApiHost:   "api.qiniu.com",
		IovipHost: "iovip.qbox.me",
	}

	// regionHuadongZhejiang2 表示华东-浙江2
	regionHuadongZhejiang2 = Region{
		SrcUpHosts: []string{
			"up-cn-east-2.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-cn-east-2.qiniup.com",
		},
		RsHost:    "rs-cn-east-2.qiniuapi.com",
		RsfHost:   "rsf-cn-east-2.qiniuapi.com",
		ApiHost:   "api-cn-east-2.qiniuapi.com",
		IovipHost: "iovip-cn-east-2.qiniuio.com",
	}

	// regionHuabei 表示华北机房
	regionHuabei = Region{
		SrcUpHosts: []string{
			"up-z1.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-z1.qiniup.com",
		},
		RsHost:    "rs-z1.qbox.me",
		RsfHost:   "rsf-z1.qbox.me",
		ApiHost:   "api-z1.qiniuapi.com",
		IovipHost: "iovip-z1.qbox.me",
	}

	// regionHuanan 表示华南机房
	regionHuanan = Region{
		SrcUpHosts: []string{
			"up-z2.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-z2.qiniup.com",
		},
		RsHost:    "rs-z2.qbox.me",
		RsfHost:   "rsf-z2.qbox.me",
		ApiHost:   "api-z2.qiniuapi.com",
		IovipHost: "iovip-z2.qbox.me",
	}

	// regionNorthAmerica 表示北美机房
	regionNorthAmerica = Region{
		SrcUpHosts: []string{
			"up-na0.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-na0.qiniup.com",
		},
		RsHost:    "rs-na0.qbox.me",
		RsfHost:   "rsf-na0.qbox.me",
		ApiHost:   "api-na0.qiniuapi.com",
		IovipHost: "iovip-na0.qbox.me",
	}

	// regionSingapore 表示新加坡机房
	regionSingapore = Region{
		SrcUpHosts: []string{
			"up-as0.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload-as0.qiniup.com",
		},
		RsHost:    "rs-as0.qbox.me",
		RsfHost:   "rsf-as0.qbox.me",
		ApiHost:   "api-as0.qiniuapi.com",
		IovipHost: "iovip-as0.qbox.me",
	}
)

const (
	// region code
	RIDHuadong          = RegionID("z0")
	RIDHuadongZheJiang2 = RegionID("cn-east-2")
	RIDHuabei           = RegionID("z1")
	RIDHuanan           = RegionID("z2")
	RIDNorthAmerica     = RegionID("na0")
	RIDSingapore        = RegionID("as0")
)

// regionMap 是RegionID到具体的Region的映射
var regionMap = map[RegionID]Region{
	RIDHuadong:          regionHuadong,
	RIDHuadongZheJiang2: regionHuadongZhejiang2,
	RIDHuanan:           regionHuanan,
	RIDHuabei:           regionHuabei,
	RIDSingapore:        regionSingapore,
	RIDNorthAmerica:     regionNorthAmerica,
}

const (
	defaultApiHost = "api.qiniu.com"
	defaultUcHost0 = "kodo-config.qiniuapi.com"
	defaultUcHost1 = "uc.qbox.me"
)

// UcHost 为查询空间相关域名的 API 服务地址
// 设置 UcHost 时，如果不指定 scheme 默认会使用 https
// Deprecated 使用 SetUcHosts 替换
var UcHost = ""

// 公有云包括 defaultApiHost，非 uc query api 使用时需要移除 defaultApiHost
// 用户配置时，不能配置 api 域名
var ucHosts = []string{defaultUcHost0, defaultUcHost1, defaultApiHost}

// SetUcHost
// Deprecated 使用 SetUcHosts 替换
func SetUcHost(host string, useHttps bool) {
	if len(host) == 0 {
		return
	}
	host = endpoint(useHttps, host)
	ucHosts = []string{host}
}

// SetUcHosts 配置多个 UC 域名
func SetUcHosts(hosts ...string) {
	var newHosts []string
	for _, host := range hosts {
		if len(host) > 0 {
			newHosts = append(newHosts, host)
		}
	}
	ucHosts = newHosts
}

func getUcHost(useHttps bool) string {
	// 兼容老版本，优先使用 UcHost
	host := ""
	if len(UcHost) > 0 {
		host = UcHost
	} else if len(ucHosts) > 0 {
		host = ucHosts[0]
	}
	return endpoint(useHttps, host)
}

// 不带 scheme
func getUcBackupHosts() []string {
	var hosts []string
	if len(UcHost) > 0 {
		hosts = append(hosts, removeHostScheme(UcHost))
	}

	for _, host := range ucHosts {
		if len(host) > 0 {
			hosts = append(hosts, removeHostScheme(host))
		}
	}

	hosts = removeRepeatStringItem(hosts)
	return hosts
}

// GetRegion 用来根据ak和bucket来获取空间相关的机房信息
// 延用 v2, v2 结构和 v4 结构不同且暂不可替代
// Deprecated 使用 GetRegionWithOptions 替换
func GetRegion(ak, bucket string) (*Region, error) {
	return GetRegionWithOptions(ak, bucket, DefaultUCApiOptions())
}

// GetRegionWithOptions 用来根据ak和bucket来获取空间相关的机房信息
func GetRegionWithOptions(ak, bucket string, options UCApiOptions) (*Region, error) {
	return getRegionByV2(ak, bucket, options)
}

// 使用 v4
func getRegionGroup(ak, bucket string) (*RegionGroup, error) {
	return getRegionByV4(ak, bucket, DefaultUCApiOptions())
}

func getRegionGroupWithOptions(ak, bucket string, options UCApiOptions) (*RegionGroup, error) {
	return getRegionByV4(ak, bucket, options)
}

type RegionInfo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func SetRegionCachePath(newPath string) {
	setRegionV2CachePath(newPath)
	setRegionV4CachePath(newPath)
}

// GetRegionsInfo Deprecated and use GetRegionsInfoWithOptions instead
// Deprecated
func GetRegionsInfo(mac *auth.Credentials) ([]RegionInfo, error) {
	return GetRegionsInfoWithOptions(mac, DefaultUCApiOptions())
}

func GetRegionsInfoWithOptions(mac *auth.Credentials, options UCApiOptions) ([]RegionInfo, error) {
	var regions struct {
		Regions []RegionInfo `json:"regions"`
	}

	reqUrl := getUcHost(options.UseHttps) + "/regions"
	c := getUCClient(ucClientConfig{
		IsUcQueryApi:       false,
		RetryMax:           options.RetryMax,
		HostFreezeDuration: options.HostFreezeDuration,
	}, mac)
	_, qErr := clientv2.DoAndDecodeJsonResponse(c, clientv2.RequestParams{
		Context:     context.Background(),
		Method:      clientv2.RequestMethodGet,
		Url:         reqUrl,
		Header:      nil,
		BodyCreator: nil,
	}, &regions)
	if qErr != nil {
		return nil, fmt.Errorf("query region error, %s", qErr.Error())
	} else {
		return regions.Regions, nil
	}
}

type ucClientConfig struct {
	// 非 uc query api 需要去除默认域名 defaultApiHost
	IsUcQueryApi bool

	// 单域名重试次数
	RetryMax int

	// 主备域名冻结时间（默认：600s），当一个域名请求失败（单个域名会被重试 TryTimes 次），会被冻结一段时间，使用备用域名进行重试，在冻结时间内，域名不能被使用，当一个操作中所有域名竣备冻结操作不在进行重试，返回最后一次操作的错误。
	HostFreezeDuration time.Duration

	Client *client.Client
}

func getUCClient(config ucClientConfig, mac *auth.Credentials) clientv2.Client {
	allHosts := getUcBackupHosts()
	var hosts []string = nil
	if !config.IsUcQueryApi {
		// 非 uc query api 去除 defaultApiHost
		for _, host := range allHosts {
			if host != defaultApiHost {
				hosts = append(hosts, host)
			}
		}
	} else {
		hosts = allHosts
	}

	is := []clientv2.Interceptor{
		clientv2.NewHostsRetryInterceptor(clientv2.HostsRetryConfig{
			RetryConfig: clientv2.RetryConfig{
				RetryMax:      len(hosts),
				RetryInterval: nil,
				ShouldRetry:   nil,
			},
			ShouldFreezeHost:   nil,
			HostFreezeDuration: 0,
			HostProvider:       hostprovider.NewWithHosts(hosts),
		}),
		clientv2.NewSimpleRetryInterceptor(clientv2.RetryConfig{
			RetryMax:      config.RetryMax,
			RetryInterval: nil,
			ShouldRetry:   nil,
		}),
	}

	if mac != nil {
		is = append(is, clientv2.NewAuthInterceptor(clientv2.AuthConfig{
			Credentials: *mac,
			TokenType:   auth.TokenQiniu,
		}))
	}

	return clientv2.NewClient(clientv2.NewClientWithClientV1(config.Client), is...)
}

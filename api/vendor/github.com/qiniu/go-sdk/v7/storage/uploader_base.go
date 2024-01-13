package storage

import (
	"github.com/qiniu/go-sdk/v7/internal/hostprovider"
	"time"
)

// retryMax: 为 0，使用默认值，每个域名只请求一次
// hostFreezeDuration: 为 0，使用默认值：50ms ~ 100ms
func getUpHost(config *Config, retryMax int, hostFreezeDuration time.Duration, ak, bucket string) (upHost string, err error) {
	region := config.GetRegion()
	if region == nil {
		if region, err = GetRegionWithOptions(ak, bucket, UCApiOptions{
			RetryMax:           retryMax,
			HostFreezeDuration: hostFreezeDuration,
		}); err != nil {
			return "", err
		}
	}

	host := region.SrcUpHosts[0]
	if config.UseCdnDomains {
		host = region.CdnUpHosts[0]
	}

	upHost = endpoint(config.UseHTTPS, host)
	return
}

// retryMax: 为 0，使用默认值，每个域名只请求一次
// hostFreezeDuration: 为 0，使用默认值：50ms ~ 100ms
func getUpHostProvider(config *Config, retryMax int, hostFreezeDuration time.Duration, ak, bucket string) (hostprovider.HostProvider, error) {
	region := config.GetRegion()
	var err error
	if region == nil {
		if region, err = GetRegionWithOptions(ak, bucket, UCApiOptions{
			RetryMax:           retryMax,
			HostFreezeDuration: hostFreezeDuration,
		}); err != nil {
			return nil, err
		}
	}

	hosts := make([]string, 0, 0)
	if config.UseCdnDomains && len(region.CdnUpHosts) > 0 {
		hosts = append(hosts, region.CdnUpHosts...)
	} else if len(region.SrcUpHosts) > 0 {
		hosts = append(hosts, region.SrcUpHosts...)
	}

	for i := 0; i < len(hosts); i++ {
		hosts[i] = endpoint(config.UseHTTPS, hosts[i])
	}

	return hostprovider.NewWithHosts(hosts), nil
}

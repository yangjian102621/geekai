package storage

import (
	"context"
	api "github.com/qiniu/go-sdk/v7"
	"github.com/qiniu/go-sdk/v7/internal/clientv2"
	"github.com/qiniu/go-sdk/v7/internal/hostprovider"
	"strings"
	"time"
)

// ParsePutTime 提供了将PutTime转换为 time.Time 的功能
func ParsePutTime(putTime int64) (t time.Time) {
	t = time.Unix(0, putTime*100)
	return
}

// IsContextExpired 检查分片上传的ctx是否过期，提前一天让它过期
// 因为我们认为如果断点继续上传的话，最长需要1天时间
func IsContextExpired(blkPut BlkputRet) bool {
	if blkPut.Ctx == "" {
		return false
	}
	target := time.Unix(blkPut.ExpiredAt, 0).AddDate(0, 0, -1)
	now := time.Now()
	return now.After(target)
}

func isUploadContextExpired(expiredAt int64) bool {
	target := time.Unix(expiredAt, 0).Add(-2 * time.Hour)
	return time.Now().After(target)
}

func shouldUploadAgain(err error) bool {
	if err == nil {
		return false
	}

	if isCancelErr(err) {
		return false
	}

	switch t := err.(type) {
	case *ErrorInfo:
		// 4xx 不重试
		return t.Code < 400 || t.Code > 499
	case *api.QError:
		return t.Code == ErrMaxUpRetry
	default:
		// 网络异常可重试
		return shouldUploadRetryWithOtherHost(err)
	}
}

func isContextExpiredError(err error) bool {
	if err == nil {
		return false
	}

	errInfo, ok := err.(*ErrorInfo)
	if !ok {
		return false
	}

	return errInfo.Code == 701 || (errInfo.Code == 612 && strings.Contains(errInfo.Error(), "no such uploadId"))
}

func shouldUploadRetryWithOtherHost(err error) bool {
	return clientv2.IsErrorRetryable(err)
}

func doUploadAction(hostProvider hostprovider.HostProvider, retryMax int, freezeDuration time.Duration, action func(host string) error) error {
	for {
		host, err := hostProvider.Provider()
		if err != nil {
			return api.NewError(ErrMaxUpRetry, err.Error())
		}

		for i := 0; ; i++ {
			err = action(host)

			// 请求成功
			if err == nil {
				return nil
			}

			// 不可重试错误
			if !shouldUploadRetryWithOtherHost(err) {
				return err
			}

			// 超过重试次数退出
			if i >= retryMax {
				break
			}
		}

		// 单个 host 失败，冻结此 host，换其他 host
		_ = hostProvider.Freeze(host, err, freezeDuration)
	}
}

func isCancelErr(err error) bool {
	if err == context.Canceled {
		return true
	}

	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "context canceled")
}

func removeRepeatStringItem(slc []string) []string {
	var result []string
	tempMap := map[string]uint8{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func removeHostScheme(host string) string {
	host = strings.TrimLeft(host, "http://")
	host = strings.TrimLeft(host, "https://")
	return host
}

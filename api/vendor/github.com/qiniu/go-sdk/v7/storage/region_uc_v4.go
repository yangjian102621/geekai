package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qiniu/go-sdk/v7/internal/clientv2"
	"golang.org/x/sync/singleflight"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ucQueryV4Ret struct {
	Hosts []ucQueryV4Region `json:"hosts"`
}

type ucQueryV4Region struct {
	RegionId string          `json:"region"`
	TTL      int             `json:"ttl"`
	Io       ucQueryV4Server `json:"io"`
	IoSrc    ucQueryV4Server `json:"io_src"`
	Up       ucQueryV4Server `json:"up"`
	Rs       ucQueryV4Server `json:"rs"`
	Rsf      ucQueryV4Server `json:"rsf"`
	Api      ucQueryV4Server `json:"api"`
}

type ucQueryV4Server struct {
	Domains []string `json:"domains"`
	Old     []string `json:"old"`
}

func (s *ucQueryV4Server) getOneServer() string {
	if len(s.Domains) > 0 {
		return s.Domains[0]
	}
	if len(s.Old) > 0 {
		return s.Old[0]
	}
	return ""
}

var ucQueryV4Group singleflight.Group

type regionV4CacheValue struct {
	Regions  []*Region `json:"regions"`
	Deadline time.Time `json:"deadline"`
}

func (r *regionV4CacheValue) getRegions() []*Region {
	if r == nil {
		return nil
	}
	return r.Regions
}

type regionV4CacheMap map[string]regionV4CacheValue

const regionV4CacheFileName = "query_v4_00.cache.json"

var (
	regionV4CachePath     = filepath.Join(os.TempDir(), "qiniu-golang-sdk", regionV4CacheFileName)
	regionV4Cache         sync.Map
	regionV4CacheLock     sync.RWMutex
	regionV4CacheSyncLock sync.Mutex
	regionV4CacheLoaded   bool = false
)

func setRegionV4CachePath(newPath string) {
	cacheDir := filepath.Dir(newPath)
	if len(cacheDir) == 0 {
		return
	}

	regionV4CacheLock.Lock()
	defer regionV4CacheLock.Unlock()

	regionV4CachePath = filepath.Join(cacheDir, regionV4CacheFileName)
	regionV4CacheLoaded = false
}

func loadRegionV4Cache() {
	cacheFile, err := os.Open(regionV4CachePath)
	if err != nil {
		return
	}
	defer cacheFile.Close()

	var cacheMap regionV4CacheMap
	if err = json.NewDecoder(cacheFile).Decode(&cacheMap); err != nil {
		return
	}
	for cacheKey, cacheValue := range cacheMap {
		regionV4Cache.Store(cacheKey, cacheValue)
	}
}

func storeRegionV4Cache() {
	err := os.MkdirAll(filepath.Dir(regionV4CachePath), 0700)
	if err != nil {
		return
	}

	cacheFile, err := os.OpenFile(regionV4CachePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer cacheFile.Close()

	cacheMap := make(regionV4CacheMap)
	regionV4Cache.Range(func(cacheKey, cacheValue interface{}) bool {
		cacheMap[cacheKey.(string)] = cacheValue.(regionV4CacheValue)
		return true
	})
	if err = json.NewEncoder(cacheFile).Encode(cacheMap); err != nil {
		return
	}
}

func getRegionByV4(ak, bucket string, options UCApiOptions) (*RegionGroup, error) {
	regionV4CacheLock.RLock()
	if regionV4CacheLoaded {
		regionV4CacheLock.RUnlock()
	} else {
		regionV4CacheLock.RUnlock()
		func() {
			regionV4CacheLock.Lock()
			defer regionV4CacheLock.Unlock()

			if !regionV4CacheLoaded {
				loadRegionV4Cache()
				regionV4CacheLoaded = true
			}
		}()
	}

	regionID := fmt.Sprintf("%s:%s", ak, bucket)
	//check from cache
	if v, ok := regionV4Cache.Load(regionID); ok && time.Now().Before(v.(regionV4CacheValue).Deadline) {
		cacheValue, _ := v.(regionV4CacheValue)
		return NewRegionGroup(cacheValue.getRegions()...), nil
	}

	newRegion, err, _ := ucQueryV4Group.Do(regionID, func() (interface{}, error) {
		reqURL := fmt.Sprintf("%s/v4/query?ak=%s&bucket=%s", getUcHost(options.UseHttps), ak, bucket)

		var ret ucQueryV4Ret
		c := getUCClient(ucClientConfig{
			IsUcQueryApi:       true,
			RetryMax:           options.RetryMax,
			HostFreezeDuration: options.HostFreezeDuration,
		}, nil)
		_, err := clientv2.DoAndDecodeJsonResponse(c, clientv2.RequestParams{
			Context:     context.Background(),
			Method:      clientv2.RequestMethodGet,
			Url:         reqURL,
			Header:      nil,
			BodyCreator: nil,
		}, &ret)
		if err != nil {
			return nil, fmt.Errorf("query region error, %s", err.Error())
		}

		ttl := math.MaxInt32
		regions := make([]*Region, 0, 0)
		for _, host := range ret.Hosts {
			if ttl > host.TTL {
				ttl = host.TTL
			}
			regions = append(regions, &Region{
				SrcUpHosts: host.Up.Domains,
				CdnUpHosts: host.Up.Domains,
				RsHost:     host.Rs.getOneServer(),
				RsfHost:    host.Rsf.getOneServer(),
				ApiHost:    host.Api.getOneServer(),
				IovipHost:  host.Io.getOneServer(),
				IoSrcHost:  host.IoSrc.getOneServer(),
			})
		}

		regionV4Cache.Store(regionID, regionV4CacheValue{
			Regions:  regions,
			Deadline: time.Now().Add(time.Duration(ttl) * time.Second),
		})

		regionV4CacheSyncLock.Lock()
		defer regionV4CacheSyncLock.Unlock()

		storeRegionV4Cache()

		return NewRegionGroup(regions...), nil
	})

	if newRegion == nil {
		return nil, err
	}

	return newRegion.(*RegionGroup), err
}

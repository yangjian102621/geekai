package storage

import "time"

type (
	// BucketV4 查询条件
	BucketV4Input struct {
		// 指定区域 ID，如果传入空字符串，则查询所有区域
		Region string
		// 最多获取的空间数，如果传入 0，则查询 20 个空间
		Limit uint64
		// 获取下一页的标记
		Marker string
	}

	// BucketV4 返回的空间信息
	BucketsV4Output struct {
		// 下页开始的 Marker
		NextMarker string `json:"next_marker"`
		// 列举是否被阶段，如果为 true，则表示还有下一页
		IsTruncated bool `json:"is_truncated"`
		// 空间列表
		Buckets []BucketV4Output `json:"buckets"`
	}

	// BucketV4 返回的空间信息
	BucketV4Output struct {
		// 空间名称
		Name string `json:"name"`
		// 空间区域 ID
		Region string `json:"region"`
		// 空间是否私有
		Private bool `json:"private"`
		// 空间创建时间
		Ctime time.Time `json:"ctime"`
	}
)

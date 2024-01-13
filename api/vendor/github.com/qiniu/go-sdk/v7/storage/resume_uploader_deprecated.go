package storage

import (
	"errors"
	"io"

	"github.com/qiniu/go-sdk/v7"
)

// 分片上传过程中可能遇到的错误
var (
	ErrInvalidPutProgress = errors.New("invalid put progress")
	ErrPutFailed          = errors.New("resumable put failed")
	ErrBadToken           = errors.New("invalid token")
)

const (
	// 上传一个分片失败
	ErrUploadChunkFailed = "ErrUploadChunkFailed"
	// 取消了分片的上传
	ErrChunkUpCanceled = "ErrChunkUpCanceled"
)

// BlockCount 用来计算文件的分块数量
func BlockCount(fsize int64) int {
	return int((fsize + blockMask) >> blockBits)
}

// 上传进度过期错误
const (
	InvalidCtx = 701 // UP: 无效的上下文(bput)，可能情况：Ctx非法或者已经被淘汰（太久未使用）
)

// ChunkPutRetryer 上传分片失败时候重试接口
type ChunkPutRetryer interface {
	Retry(ck *Chunk)
}

// Chunk表示要上传的数据块, 该片的大小不能大于4M
// 上传块的过程： 1. 调用接口在七牛后端创建块 2. 上传数据到该块
// 详细可以参考 https://developer.qiniu.com/kodo/api/1286/mkblk
type Chunk struct {
	// 要上传的块数据
	Body io.ReadSeeker

	// 该片数据所属的块的大小
	BlkSize int

	// 将大的数据切成4M一块，每个块都有一个index
	// 该片数据所属的块的index
	Index int

	// 上传该片发生的错误
	Err error

	// 是否调用了mkblk接口在后台创建了该片所属的块
	Created bool

	// 上传块的返回值
	Ret *BlkputRet
}

// ShouldRetry 是否需要重新上传
func (b *Chunk) ShouldRetry() bool {
	return b.Err != nil
}

// BlockLength 返回实际要上传的数据的长度
func (b *Chunk) ChunkLength() (int, error) {
	n, err := api.SeekerLen(b.Body)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// ResetBody 重置Body到开头
func (b *Chunk) ResetBody() error {
	_, err := b.Body.Seek(0, io.SeekStart)
	return err
}

// Reset 重置Body和Err
func (b *Chunk) Reset() error {
	err := b.ResetBody()
	if err != nil {
		return err
	}
	b.Err = nil
	return nil
}

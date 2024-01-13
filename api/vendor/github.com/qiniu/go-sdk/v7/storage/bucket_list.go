package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"net/url"
	"strconv"
)

// ListItem 为文件列举的返回值
type ListItem struct {

	// 资源名
	Key string `json:"key"`

	// 上传时间，单位：100纳秒，其值去掉低七位即为Unix时间戳。
	PutTime int64 `json:"putTime"`

	// 文件的HASH值，使用hash值算法计算。
	Hash string `json:"hash"`

	// 资源内容的大小，单位：字节。
	Fsize int64 `json:"fsize"`

	// 资源的 MIME 类型。
	MimeType string `json:"mimeType"`

	/**
	 * 文件上传时设置的endUser
	 */
	EndUser string `json:"endUser"`

	/**
	 * 资源的存储类型
	 * 0 表示标准存储
	 * 1 表示低频存储
	 * 2 表示归档存储
	 * 3 表示深度归档存储
	 */
	Type int `json:"type"`

	/**
	 * 文件的存储状态，即禁用状态和启用状态间的的互相转换，请参考：文件状态。
	 * 0 表示启用
	 * 1 表示禁用
	 */
	Status int `json:"status"`

	/**
	 * 文件的 md5 值
	 */
	Md5 string `json:"md5"`
}

// 接口可能返回空的记录
func (l *ListItem) IsEmpty() (empty bool) {
	if l == nil {
		return true
	}

	return l.Key == "" && l.Hash == "" && l.Fsize == 0 && l.PutTime == 0
}

func (l *ListItem) String() string {
	str := ""
	str += fmt.Sprintf("Hash:     %s\n", l.Hash)
	str += fmt.Sprintf("Fsize:    %d\n", l.Fsize)
	str += fmt.Sprintf("PutTime:  %d\n", l.PutTime)
	str += fmt.Sprintf("MimeType: %s\n", l.MimeType)
	str += fmt.Sprintf("Type:     %d\n", l.Type)
	str += fmt.Sprintf("EndUser:  %s\n", l.EndUser)
	return str
}

type ListFilesRet struct {
	Marker         string     `json:"marker"`
	Items          []ListItem `json:"items"`
	CommonPrefixes []string   `json:"commonPrefixes"`
}

// ListFiles 用来获取空间文件列表，可以根据需要指定文件的前缀 prefix，文件的目录 delimiter，循环列举的时候下次
// 列举的位置 marker，以及每次返回的文件的最大数量limit，其中limit最大为1000。
func (m *BucketManager) ListFiles(bucket, prefix, delimiter, marker string,
	limit int) (entries []ListItem, commonPrefixes []string, nextMarker string, hasNext bool, err error) {

	ret, hNext, e := m.ListFilesWithContext(context.Background(), bucket,
		ListInputOptionsPrefix(prefix),
		ListInputOptionsDelimiter(delimiter),
		ListInputOptionsMarker(marker),
		ListInputOptionsLimit(limit))
	if e != nil {
		return nil, nil, "", false, e
	}
	return ret.Items, ret.CommonPrefixes, ret.Marker, hNext, nil
}

type listInputOptions struct {
	prefix    string
	delimiter string
	marker    string
	limit     int
}

type ListInputOption func(options *listInputOptions)

func ListInputOptionsPrefix(prefix string) ListInputOption {
	return func(input *listInputOptions) {
		input.prefix = prefix
	}
}

func ListInputOptionsDelimiter(delimiter string) ListInputOption {
	return func(input *listInputOptions) {
		input.delimiter = delimiter
	}
}

func ListInputOptionsMarker(marker string) ListInputOption {
	return func(input *listInputOptions) {
		input.marker = marker
	}
}

func ListInputOptionsLimit(limit int) ListInputOption {
	return func(input *listInputOptions) {
		input.limit = limit
	}
}

// ListFilesWithContext
//
//	 @Description: 用来获取空间文件列表，可以根据需要指定文件的列举条件
//	 @receiver m BucketManager
//	 @param ctx context
//	 @param bucket 列举的 bucket
//	 @param options 列举的可选条件
//					列举条件-需要列举 Key 的前缀：ListInputOptionsPrefix(prefix)
//					列举条件-文件的目录分隔符：ListInputOptionsDelimiter(delimiter)
//					列举条件-下次列举的位置：ListInputOptionsMarker(marker)
//					列举条件-每次返回的文件的最大数量：ListInputOptionsLimit(limit) 范围：1~1000
//	 @return ret 列举的对象数据
//	 @return hasNext 是否还有数据未被列举
//	 @return err 列举时的错误信息
func (m *BucketManager) ListFilesWithContext(ctx context.Context, bucket string, options ...ListInputOption) (ret *ListFilesRet, hasNext bool, err error) {
	if len(bucket) == 0 {
		return nil, false, errors.New("bucket can't empty")
	}

	inputOptions := listInputOptions{}
	for _, option := range options {
		option(&inputOptions)
	}

	if inputOptions.limit <= 0 || inputOptions.limit > 1000 {
		return nil, false, errors.New("invalid list limit, only allow [1, 1000]")
	}

	ctx = auth.WithCredentialsType(ctx, m.Mac, auth.TokenQiniu)
	host, reqErr := m.RsfReqHost(bucket)
	if reqErr != nil {
		return nil, false, reqErr
	}

	ret = &ListFilesRet{}
	reqURL := fmt.Sprintf("%s%s", host, uriListFiles(bucket, inputOptions.prefix, inputOptions.delimiter, inputOptions.marker, inputOptions.limit))
	err = m.Client.CredentialedCall(ctx, m.Mac, auth.TokenQiniu, ret, "POST", reqURL, nil)
	if err != nil {
		return nil, false, err
	}

	return ret, len(ret.Marker) > 0, nil
}

type listFilesRet2 struct {
	Marker string   `json:"marker"`
	Item   ListItem `json:"item"`
	Dir    string   `json:"dir"`
}

// ListBucket 用来获取空间文件列表，可以根据需要指定文件的前缀 prefix，文件的目录 delimiter，流式返回每条数据。
// Deprecated
func (m *BucketManager) ListBucket(bucket, prefix, delimiter, marker string) (retCh chan listFilesRet2, err error) {
	return m.ListBucketContext(context.Background(), bucket, prefix, delimiter, marker)
}

// ListBucketContext 用来获取空间文件列表，可以根据需要指定文件的前缀 prefix，文件的目录 delimiter，流式返回每条数据。
// 接受的context可以用来取消列举操作
// Deprecated
func (m *BucketManager) ListBucketContext(ctx context.Context, bucket, prefix, delimiter, marker string) (retCh chan listFilesRet2, err error) {

	ret, _, lErr := m.ListFilesWithContext(ctx, bucket,
		ListInputOptionsLimit(250),
		ListInputOptionsPrefix(prefix),
		ListInputOptionsDelimiter(delimiter),
		ListInputOptionsMarker(marker))
	if lErr != nil {
		return nil, lErr
	}

	count := len(ret.CommonPrefixes) + len(ret.Items)
	retCh = make(chan listFilesRet2, count)
	defer close(retCh)

	if len(ret.CommonPrefixes) > 0 {
		for _, commonPrefix := range ret.CommonPrefixes {
			retCh <- listFilesRet2{
				Marker: ret.Marker,
				Item:   ListItem{},
				Dir:    commonPrefix,
			}
		}
	}

	if len(ret.Items) > 0 {
		for _, item := range ret.Items {
			retCh <- listFilesRet2{
				Marker: ret.Marker,
				Item:   item,
				Dir:    "",
			}
		}
	}

	return retCh, err
}

func uriListFiles(bucket, prefix, delimiter, marker string, limit int) string {
	query := make(url.Values)
	query.Add("bucket", bucket)
	if prefix != "" {
		query.Add("prefix", prefix)
	}
	if delimiter != "" {
		query.Add("delimiter", delimiter)
	}
	if marker != "" {
		query.Add("marker", marker)
	}
	if limit > 0 {
		query.Add("limit", strconv.FormatInt(int64(limit), 10))
	}
	return fmt.Sprintf("/list?%s", query.Encode())
}

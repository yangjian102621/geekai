# Changelog
## 7.17.1
* 优化
  * 调整在获取 Bucket 所在区域服务域名时的主备域名顺序
* 调整
  * 移除内置的亚太-首尔区域

## 7.17.0
* 优化
  * 对象存储，UC 服务相关请求支持主备重试

## 7.16.0
* 新增
    * 对象存储，`BucketManager` `BucketsV4` 获取该用户的指定区域内的空间信息，注意该 API 以分页形式返回 Bucket 列表
    * 对象存储，`BucketManager` `SetRemark` 设置空间的备注信息
* 更新
    * 对象存储，`BucketManager` `GetBucketInfo` 返回空间的备注信息

## 7.15.0
* 优化
    * 表单上传 Put 方法的 PutExtra 参数，支持传入 nil
    * Bucket 镜像源/配额方法内部请求使用 UC 域名
    * `BucketManager` `ListBucket` 和 `ListBucketContext` 方法内部接口由 /v2/list 调用调整为 /list
* 新增
    * `BucketManager` 新增批量方法 `BatchWithContext`
    * `BucketManager` 增加 Bucket 列举方法 `ListFileWithContext`

## 7.14.0
* Go SDK  对七牛回调请求的鉴权验证函数支持 Qiniu 签名
* UploadMananger 支持双活区域
* 优化分片上传 ctx 超时检测
* 表单 / 分片 v1 / 分片 v2 上传支持主备域名重试逻辑。
* 表单上传支持配置重试次数
* 处理 Config 配置固定 Region， 但不生效的问题(配置 Zone 生效)
* 移除雾存储区域：华东一区
* 增加亚太-首尔区域固定 Region

## 7.13.0
* 对象存储，管理类 API 发送请求时增加 [X-Qiniu-Date](https://developer.qiniu.com/kodo/3924/common-request-headers) （生成请求的时间） header

## 7.12.1
* 对象存储，补充 Stat API 可查看对象元数据信息

## 7.12.0
* 对象存储，新增支持 [深度归档存储类型](https://developer.qiniu.com/kodo/3956/kodo-category#deep_archive)
* 对象存储，全面支持 Qiniu 签名

## 7.11.1
* 优化容器环境下 pod 当前内存工作集 (working set) 使用量

## 7.11.0
* 新增直播云服务端管理能力，包括：直播空间管理、域名管理、推拉流地址管理、直播流管理和统计数据查询 API


## 7.10.1
* 优化了分片上传内存占用
* 修复部分已知问题

## 7.10.0
* 增加了 PutWithoutKeyAndSize API，在上传时支持可不指定 size 和 key
* 修复了 已知 UcQuery 解析问题
## 7.9.8
* 补充了查询 object 元信息返回部分字段

## 7.9.7
* 修复了表单上传 FormUploader 在内部重试情况下的已知问题

## 7.9.6
* 在需要指定存储服务 host 情况下兼容了只配置域名和同时指定域名和访问 protocol 的问题

## 7.9.5
优化几个已知小问题
* 支持指定空间管理域名，默认是公有云地址
* 支持下载 escape 编码文件
* 优化对一些错误情况处理

## 7.9.4
* 兼容不同格式上传token

## 7.9.3
* 修复在复用上传token时，过期时间叠加问题

## 7.9.2
* UploadPartInfo 结构体公开使用，可用于定制分片上传过程
* 保持兼容支持上传API extra.UpHost参数

## 7.9.1
* 修复buckets api 已知问题

## 7.9.0
* 从 github.com/qiniu/api.v7 迁移至 github.com/qiniu/go-sdk

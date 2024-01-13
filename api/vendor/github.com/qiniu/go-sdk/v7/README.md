github.com/qiniu/go-sdk
===============

[![LICENSE](https://img.shields.io/github/license/qiniu/go-sdk.svg)](https://github.com/qiniu/go-sdk/blob/master/LICENSE)
[![Build Status](https://github.com/qiniu/go-sdk/workflows/Run%20Test%20Cases/badge.svg)](https://github.com/qiniu/go-sdk/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/go-sdk)](https://goreportcard.com/report/github.com/qiniu/go-sdk)
[![GitHub release](https://img.shields.io/github/v/tag/qiniu/go-sdk.svg?label=release)](https://github.com/qiniu/go-sdk/releases)
[![codecov](https://codecov.io/gh/qiniu/go-sdk/branch/master/graph/badge.svg)](https://codecov.io/gh/qiniu/go-sdk)
[![GoDoc](https://godoc.org/github.com/qiniu/go-sdk?status.svg)](https://godoc.org/github.com/qiniu/go-sdk)

[![Qiniu Logo](http://open.qiniudn.com/logo.png)](http://qiniu.com/)

# 下载

## 使用 Go mod【推荐】

在您的项目中的 `go.mod` 文件内添加这行代码

```
require github.com/qiniu/go-sdk/v7 v7.17.1
```

并且在项目中使用 `"github.com/qiniu/go-sdk/v7"` 引用 Qiniu Go SDK。

例如

```go
import (
    "github.com/qiniu/go-sdk/v7/auth"
    "github.com/qiniu/go-sdk/v7/storage"
)
```

# Golang 版本需求

需要 go1.10 或者 1.10 以上

#  文档

[七牛SDK文档站](https://developer.qiniu.com/kodo/sdk/1238/go) 或者 [项目WIKI](https://github.com/qiniu/go-sdk/wiki)

# 示例

[参考代码](https://github.com/qiniu/go-sdk/tree/master/examples)

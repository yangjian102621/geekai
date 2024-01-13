package conf

import (
	"os"
	"strings"
)

const Version = "7.17.1"

const (
	CONTENT_TYPE_JSON      = "application/json"
	CONTENT_TYPE_FORM      = "application/x-www-form-urlencoded"
	CONTENT_TYPE_OCTET     = "application/octet-stream"
	CONTENT_TYPE_MULTIPART = "multipart/form-data"

	disableQiniuTimestampSignatureEnvKey = "DISABLE_QINIU_TIMESTAMP_SIGNATURE"
)

func IsDisableQiniuTimestampSignature() bool {
	value := os.Getenv(disableQiniuTimestampSignatureEnvKey)
	value = strings.ToLower(value)
	return value == "true" || value == "yes" || value == "y" || value == "1"
}

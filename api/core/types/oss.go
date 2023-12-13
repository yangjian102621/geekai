package types

type OSSConfig struct {
	Active string
	Local  LocalStorageConfig
	Minio  MiniOssConfig
	QiNiu  QiNiuOssConfig
	AliYun AliYunOssConfig
}
type MiniOssConfig struct {
	Endpoint     string
	AccessKey    string
	AccessSecret string
	Bucket       string
	SubDir       string
	UseSSL       bool
	Domain       string
}

type QiNiuOssConfig struct {
	Zone         string
	AccessKey    string
	AccessSecret string
	Bucket       string
	SubDir       string
	Domain       string
}

type AliYunOssConfig struct {
	Endpoint     string
	AccessKey    string
	AccessSecret string
	Bucket       string
	SubDir       string
	Domain       string
}

type LocalStorageConfig struct {
	BasePath string
	BaseURL  string
}

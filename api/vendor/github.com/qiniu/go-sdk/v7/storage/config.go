package storage

// Config 为文件上传，资源管理等配置
type Config struct {
	//兼容保留
	Zone *Region //空间所在的存储区域

	Region *Region

	// 如果设置的Host本身是以http://开头的，又设置了该字段为true，那么优先使用该字段，使用https协议
	// 同理如果该字段为false, 但是设置的host以https开头，那么使用http协议通信
	UseHTTPS      bool   //是否使用https域名
	UseCdnDomains bool   //是否使用cdn加速域名
	CentralRsHost string //中心机房的RsHost，用于list bucket

	// 兼容保留
	RsHost  string
	RsfHost string
	UpHost  string
	ApiHost string
	IoHost  string
}

// reqHost 返回一个Host链接
// 主要用于Config 中Host的获取，Region优先级最高， Zone次之， 最后才使用设置的Host信息
// topHost是优先级最高的， host次之，如果都没有，使用默认的defaultHost
func reqHost(useHttps bool, topHost, host, defaultHost string) (endp string) {
	if topHost != "" {
		return topHost
	}
	if host == "" {
		host = defaultHost
	}
	return endpoint(useHttps, host)
}

// 获取RsHost
// 优先使用Zone中的Host信息，如果Zone中的host信息没有配置，那么使用Config中的Host信息
func (c *Config) RsReqHost() string {
	rzHost := c.hostFromRegion("rs")
	return reqHost(c.UseHTTPS, rzHost, c.RsHost, DefaultRsHost)
}

// GetRegion返回一个Region指针
// 默认返回最新的Region， 如果该字段没有，那么返回兼容保留的Zone, 如果都为nil, 就返回nil
func (c *Config) GetRegion() *Region {
	if c.Region != nil {
		return c.Region
	}
	if c.Zone != nil {
		return c.Zone
	}
	return nil
}

func (c *Config) hostFromRegion(typ string) string {
	region := c.GetRegion()
	if region != nil {
		switch typ {
		case "rs":
			return region.GetRsHost(c.UseHTTPS)
		case "rsf":
			return region.GetRsfHost(c.UseHTTPS)
		case "api":
			return region.GetApiHost(c.UseHTTPS)
		case "io":
			return region.GetIoHost(c.UseHTTPS)
		}
	}
	return ""
}

// 获取rsfHost
// 优先使用Zone中的Host信息，如果Zone中的host信息没有配置，那么使用Config中的Host信息
func (c *Config) RsfReqHost() string {
	rsHost := c.hostFromRegion("rsf")
	return reqHost(c.UseHTTPS, rsHost, c.RsfHost, DefaultRsfHost)
}

// 获取apiHost
// 优先使用Zone中的Host信息，如果Zone中的host信息没有配置，那么使用Config中的Host信息
func (c *Config) ApiReqHost() string {
	rzHost := c.hostFromRegion("api")
	return reqHost(c.UseHTTPS, rzHost, c.ApiHost, DefaultAPIHost)
}

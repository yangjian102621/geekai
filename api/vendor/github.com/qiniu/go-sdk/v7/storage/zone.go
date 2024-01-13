package storage

// Zone 是Region的别名
// 兼容保留
type Zone = Region

// GetZone 用来根据ak和bucket来获取空间相关的机房信息
// 新版本使用GetRegion, 这个函数用来保持兼容
func GetZone(ak, bucket string) (zone *Zone, err error) {
	return GetRegion(ak, bucket)
}

var (
	// 华东机房
	// 兼容保留
	ZoneHuadong, _ = GetRegionByID(RIDHuadong)

	// 华北机房
	// 兼容保留
	ZoneHuabei, _ = GetRegionByID(RIDHuabei)

	// 华南机房
	// 兼容保留
	ZoneHuanan, _ = GetRegionByID(RIDHuanan)

	// 北美机房
	// 兼容保留
	ZoneBeimei, _ = GetRegionByID(RIDNorthAmerica)

	// 新加坡机房
	// 兼容保留
	ZoneXinjiapo, _ = GetRegionByID(RIDSingapore)

	// 华东浙江 2 区
	ZoneHuadongZheJiang2, _ = GetRegionByID(RIDHuadongZheJiang2)

	// 兼容保留
	Zone_z0 = ZoneHuadong
	// 兼容保留
	Zone_z1 = ZoneHuabei
	// 兼容保留
	Zone_z2 = ZoneHuanan
	// 兼容保留
	Zone_na0 = ZoneBeimei
	// 兼容保留
	Zone_as0 = ZoneXinjiapo
)

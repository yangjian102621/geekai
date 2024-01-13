package storage

import "sync"

type RegionGroup struct {
	locker             sync.RWMutex
	currentRegionIndex int
	regions            []*Region
}

func NewRegionGroup(region ...*Region) *RegionGroup {
	return &RegionGroup{
		locker:             sync.RWMutex{},
		currentRegionIndex: 0,
		regions:            region,
	}
}

func (g *RegionGroup) GetRegion() *Region {
	g.locker.RLock()
	defer g.locker.RUnlock()
	if g.currentRegionIndex >= len(g.regions) {
		return nil
	}
	return g.regions[g.currentRegionIndex]
}

func (g *RegionGroup) CouldSwitchRegion() bool {
	g.locker.RLock()
	defer g.locker.RUnlock()
	return len(g.regions) > (g.currentRegionIndex + 1)
}

func (g *RegionGroup) SwitchRegion() bool {
	g.locker.Lock()
	defer g.locker.Unlock()
	if len(g.regions) <= (g.currentRegionIndex + 1) {
		return false
	}
	g.currentRegionIndex++
	return true
}

func (g *RegionGroup) clone() *RegionGroup {
	g.locker.RLock()
	defer g.locker.RUnlock()
	return &RegionGroup{
		locker:             sync.RWMutex{},
		currentRegionIndex: g.currentRegionIndex,
		regions:            g.regions,
	}
}

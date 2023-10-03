package core

type RegionType int

const (
	RegionTypeNone RegionType = iota
	RegionTypeSelection
	RegionTypeCursor
	RegionTypeDirty
	RegionTypeHighlight
)

type Range struct {
	Start int64
	End   int64
}

type IndexedRange struct {
	Index int
	Range
}

type Region struct {
	Type RegionType
	Range
}

// GetActiveRegions returns the list of regions that are active at the given
// position. Regions are assumed to be sorted by position.
func GetActiveRegions(regions []Region, pos int64) []Region {
	active := make([]Region, 0)
	for _, r := range regions {
		if pos >= r.Start && pos <= r.End {
			active = append(active, r)
		}
		if r.Start > pos {
			break
		}
	}
	return active
}

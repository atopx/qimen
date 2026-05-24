package plate

import (
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/internal/tables"
)

// FindStem searches a StemPlate for a particular stem. When excludeCenter
// is true, palace 5 is skipped. Returns (palace, true) when found.
func FindStem(p *StemPlate, stem almanac.Stem, excludeCenter bool) (uint8, bool) {
	stemIdx := stem.Index()
	for palace := uint8(1); palace <= 9; palace++ {
		if excludeCenter && palace == 5 {
			continue
		}
		if v, ok := p.Get(palace); ok && v.Index() == stemIdx {
			return palace, true
		}
	}
	return 0, false
}

// FindHourStem locates the palace of the current 时辰 heavenly stem.
//
// 甲 stem is "hidden" so it is mapped to zhiFuPalace. Other stems are
// looked up in the StemPlate; if found in center (palace 5), the
// classical 寄宫 rule maps it to 2 (坤).
func FindHourStem(p *StemPlate, stem almanac.Stem, zhiFuPalace uint8) uint8 {
	if stem.Index() == 0 {
		return zhiFuPalace
	}
	if palace, ok := FindStem(p, stem, true); ok {
		return palace
	}
	if c, ok := p.Get(5); ok && c.Index() == stem.Index() {
		return 2
	}
	return zhiFuPalace
}

// FindDoor searches a DoorPlate for a particular door.
func FindDoor(p *DoorPlate, door enum.Door) (uint8, bool) {
	for palace := uint8(1); palace <= 9; palace++ {
		if v, ok := p.Get(palace); ok && v == door {
			return palace, true
		}
	}
	return 0, false
}

// IsStemInTomb reports whether a stem lands in its 入墓 palace.
func IsStemInTomb(stem almanac.Stem, palace uint8) bool {
	idx := stem.Index()
	if idx < 0 || idx >= 10 {
		return false
	}
	return tables.StemTombPalace[idx] == palace
}

// TenXunStartBranchIndex returns the starting 地支 index of the 旬
// containing the given hour pillar.
func TenXunStartBranchIndex(hour almanac.Cycle) int {
	tenIdx := hour.Ten().Index()
	if tenIdx > 5 {
		tenIdx = 5
	}
	if tenIdx < 0 {
		tenIdx = 0
	}
	return int(tables.TenXunStartBranch[tenIdx])
}

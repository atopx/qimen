package plate

import (
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/internal/tables"
)

// BuildEarth populates 地盘 (三奇六仪) from the local 局 and 阴/阳 遁.
func BuildEarth(yy almanac.YinYang, ju uint8) StemPlate {
	var p StemPlate
	palace := ju
	for _, stemIdx := range tables.SanQiLiuYi {
		p.Set(palace, almanac.StemOf(int(stemIdx)))
		palace = StepPalace(palace, yy)
	}
	return p
}

// BuildHeaven populates 天盘 from 地盘 + 值符原宫 + 时干落宫.
func BuildHeaven(earth *StemPlate, yy almanac.YinYang, zhiFuOriginalPalace, hourPalace uint8) StemPlate {
	var p StemPlate
	zhiFuIdx := int(tables.LuoShuIndex[zhiFuOriginalPalace])
	hourIdx := int(tables.LuoShuIndex[hourPalace])
	var steps int
	if yy == almanac.Yang {
		steps = (hourIdx + 8 - zhiFuIdx) % 8
	} else {
		steps = (zhiFuIdx + 8 - hourIdx) % 8
	}
	for i, palace := range tables.LuoShuOrder {
		if stem, ok := earth.Get(palace); ok {
			var targetIdx int
			if yy == almanac.Yang {
				targetIdx = (i + steps) % 8
			} else {
				targetIdx = (i + 8 - steps) % 8
			}
			p.Set(tables.LuoShuOrder[targetIdx], stem)
		}
	}
	if center, ok := earth.Get(5); ok {
		p.Set(5, center)
	}
	return p
}

// BuildStar populates 九星盘. 禽 falls into 2 with the merged 禽芮 marker.
func BuildStar(zhiFuOriginalPalace, hourPalace uint8) StarPlate {
	var p StarPlate
	zhiFuIdx := int(tables.LuoShuIndex[zhiFuOriginalPalace])
	hourIdx := int(tables.LuoShuIndex[hourPalace])
	steps := (hourIdx + 8 - zhiFuIdx) % 8
	for i, originalPalace := range tables.LuoShuOrder {
		targetPalace := tables.LuoShuOrder[(i+steps)%8]
		// 禽 has no native palace; falls into palace 2's slot as 禽芮.
		var star enum.Star
		if originalPalace == 2 {
			star = enum.StarQinRui
		} else {
			star = enum.StarOfPalace(originalPalace)
		}
		p.Set(targetPalace, star)
	}
	return p
}

// BuildDoor populates 八门盘 and returns the 值使落宫 alongside it.
//
// zhiShiOriginalPalace is the home palace of the 值使 door (= 值符原宫);
// the door advances from there by the number of 时辰 elapsed since the
// 旬首, which also determines the rotation of the other seven doors.
func BuildDoor(yy almanac.YinYang, zhiShiOriginalPalace uint8, hour almanac.Cycle) (DoorPlate, uint8) {
	var p DoorPlate
	xunStartBranchIdx := hour.Ten().FirstBranch().Index()
	hourBranchIdx := hour.Branch().Index()
	branchSteps := (hourBranchIdx + 12 - xunStartBranchIdx) % 12
	zhiShiPalace := MoveBy(zhiShiOriginalPalace, branchSteps, yy)
	zhiShiOriginIdx := int(tables.LuoShuIndex[zhiShiOriginalPalace])
	zhiShiIdx := int(tables.LuoShuIndex[zhiShiPalace])
	var steps int
	if yy == almanac.Yang {
		steps = (zhiShiIdx + 8 - zhiShiOriginIdx) % 8
	} else {
		steps = (zhiShiOriginIdx + 8 - zhiShiIdx) % 8
	}
	for i, originalPalace := range tables.LuoShuOrder {
		var targetIdx int
		if yy == almanac.Yang {
			targetIdx = (i + steps) % 8
		} else {
			targetIdx = (i + 8 - steps) % 8
		}
		// LuoShuOrder excludes center palace 5; every entry has a door.
		p.Set(tables.LuoShuOrder[targetIdx], enum.DoorOfPalace(originalPalace))
	}
	return p, zhiShiPalace
}

// BuildGod populates 九神盘: the gods follow the LuoShu ring from the
// 值符落宫, forward in 阳遁 and backward in 阴遁.
func BuildGod(yy almanac.YinYang, zhiFuPalace uint8) GodPlate {
	var p GodPlate
	startPalace := zhiFuPalace
	if startPalace == 5 {
		startPalace = 2
	}
	start := int(tables.LuoShuIndex[startPalace])
	for i, god := range tables.GodsOrder {
		var pos int
		if yy == almanac.Yang {
			pos = (start + i) % 8
		} else {
			pos = (start - i + 8) % 8
		}
		p.Set(tables.LuoShuOrder[pos], god)
	}
	return p
}

// BuildHidden populates 暗干盘 (从 8 宫起).
func BuildHidden(yy almanac.YinYang) StemPlate {
	var p StemPlate
	palace := uint8(8)
	for _, stemIdx := range tables.SanQiLiuYi {
		p.Set(palace, almanac.StemOf(int(stemIdx)))
		palace = StepPalace(palace, yy)
	}
	return p
}

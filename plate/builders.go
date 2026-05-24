package plate

import (
	"errors"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/internal/tables"
)

// ErrUnsupportedStyle indicates a chart style that is not yet implemented.
//
// Style validation is a Chart-entry-point concern; the BuildXxx
// primitives assume StyleRotate semantics and have no error return.
// Currently only enum.StyleRotate is supported; StyleFly is reserved.
var ErrUnsupportedStyle = errors.New("plate: unsupported chart style")

// BuildEarth populates 地盘 (三奇六仪) from the local 局 and 阴/阳 遁.
func BuildEarth(yy almanac.YinYang, ju uint8) *StemPlate {
	p := New[almanac.Stem]()
	palace := ju
	for _, stemIdx := range tables.SanQiLiuYi {
		p.Set(palace, almanac.StemOf(int(stemIdx)))
		palace = StepPalace(palace, yy)
	}
	return p
}

// BuildHeaven populates 天盘 from 地盘 + 值符落宫 + 时干落宫.
func BuildHeaven(earth *StemPlate, yy almanac.YinYang, zhiFuPalace, hourPalace uint8) *StemPlate {
	p := New[almanac.Stem]()
	zhiFuIdx := int(tables.LuoShuIndex[zhiFuPalace])
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
func BuildStar(zhiFuPalace, hourPalace uint8) *StarPlate {
	p := New[enum.Star]()
	zhiFuIdx := int(tables.LuoShuIndex[zhiFuPalace])
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

// BuildDoor populates 八门盘.
func BuildDoor(yy almanac.YinYang, zhiFuPalace uint8, hour almanac.Cycle) *DoorPlate {
	p := New[enum.Door]()
	xunStartBranchIdx := TenXunStartBranchIndex(hour)
	hourBranchIdx := hour.Branch().Index()
	branchSteps := (hourBranchIdx + 12 - xunStartBranchIdx) % 12
	zhiShiPalace := MoveBy(zhiFuPalace, branchSteps, yy)
	zhiFuIdx := int(tables.LuoShuIndex[zhiFuPalace])
	zhiShiIdx := int(tables.LuoShuIndex[zhiShiPalace])
	var steps int
	if yy == almanac.Yang {
		steps = (zhiShiIdx + 8 - zhiFuIdx) % 8
	} else {
		steps = (zhiFuIdx + 8 - zhiShiIdx) % 8
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
	return p
}

// BuildGod populates 九神盘.
func BuildGod(yy almanac.YinYang, zhiFuPalace uint8) *GodPlate {
	p := New[enum.God]()
	var order [8]uint8
	if yy == almanac.Yang {
		order = tables.GodYangOrder
	} else {
		order = tables.GodYinOrder
	}
	startPalace := zhiFuPalace
	if startPalace == 5 {
		startPalace = 2
	}
	startIndex := 0
	for i, v := range order {
		if v == startPalace {
			startIndex = i
			break
		}
	}
	for i, god := range tables.GodsOrder {
		p.Set(order[(startIndex+i)%8], god)
	}
	return p
}

// BuildHidden populates 暗干盘 (从 8 宫起).
func BuildHidden(yy almanac.YinYang) *StemPlate {
	p := New[almanac.Stem]()
	palace := uint8(8)
	for _, stemIdx := range tables.SanQiLiuYi {
		p.Set(palace, almanac.StemOf(int(stemIdx)))
		palace = StepPalace(palace, yy)
	}
	return p
}

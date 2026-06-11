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

// RotateStems rotates the 地盘 stems rigidly around the LuoShu ring by
// the arc from palace `from` to palace `to` (the center stem stays in
// place). Rotation is direction-free: 阴/阳 遁 only decide where the
// duty lands, not how the ring turns. This single primitive lays both
// rotated stem plates of a chart:
//
//   - 天盘: the stems follow the 值符 (from 值符原宫 to 时干落宫);
//   - 暗干: the stems follow the 值使 (from 值使门本位宫 to 值使落宫).
//
// Precondition: from, to ∈ [1, 9] \ {5} (center is projected to 坤 2
// by the caller).
func RotateStems(earth *StemPlate, from, to uint8) StemPlate {
	var p StemPlate
	shift := (int(tables.LuoShuIndex[to]) - int(tables.LuoShuIndex[from]) + 8) % 8
	for i, palace := range tables.LuoShuOrder {
		if stem, ok := earth.Get(palace); ok {
			p.Set(tables.LuoShuOrder[(i+shift)%8], stem)
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

// BuildDoor populates 八门盘 and returns the (real, possibly center)
// 值使落宫 alongside it.
//
// The 值使 door marches from the real 值符原宫 (marchFrom, which may be
// the center palace 5) by the number of 时辰 elapsed since the 旬首;
// the eight doors then rotate rigidly by the arc from the door's home
// palace (zhiShiHome, the 寄坤 projection of the origin) to the landing
// palace (likewise projected when it is the center).
func BuildDoor(yy almanac.YinYang, zhiShiHome, marchFrom uint8, hour almanac.Cycle) (DoorPlate, uint8) {
	var p DoorPlate
	xunStartBranchIdx := hour.Ten().FirstBranch().Index()
	hourBranchIdx := hour.Branch().Index()
	branchSteps := (hourBranchIdx + 12 - xunStartBranchIdx) % 12
	zhiShiPalace := MoveBy(marchFrom, branchSteps, yy)
	landEff := zhiShiPalace
	if landEff == 5 {
		landEff = 2
	}
	shift := (int(tables.LuoShuIndex[landEff]) - int(tables.LuoShuIndex[zhiShiHome]) + 8) % 8
	for i, originalPalace := range tables.LuoShuOrder {
		// LuoShuOrder excludes center palace 5; every entry has a door.
		p.Set(tables.LuoShuOrder[(i+shift)%8], enum.DoorOfPalace(originalPalace))
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

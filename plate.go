package qimen

import (
	"github.com/6tail/tyme4go/tyme"
)

// ===================== 常量表 =====================

// SanQiLiuYi 三奇六仪填充顺序: 戊己庚辛壬癸丁丙乙 (对应天干索引 4..=8, 9, 3, 2, 1)。
var SanQiLiuYi = [9]uint8{4, 5, 6, 7, 8, 9, 3, 2, 1}

// LuoShuOrder 九宫八方排列, 坎宫起, 顺八卦次序 (跳过中宫): 坎坤震巽乾兑艮离。
var LuoShuOrder = [8]uint8{1, 8, 3, 4, 9, 2, 7, 6}

// LuoShuIndex 由宫位号反查其在 LuoShuOrder 中的位置 (越界返回 0)。
var LuoShuIndex = [10]uint8{0, 0, 5, 2, 3, 0, 7, 6, 1, 4}

// GodYangOrder 阳遁九神排列。
var GodYangOrder = [8]uint8{1, 8, 3, 4, 9, 2, 7, 6}

// GodYinOrder 阴遁九神排列。
var GodYinOrder = [8]uint8{1, 6, 7, 2, 9, 4, 3, 8}

// GodsOrder 九神排列顺序 (依次: 值符、腾蛇、太阴、六合、白虎、玄武、九地、九天)。
var GodsOrder = [8]QimenGod{
	QimenGodZhiFu,
	QimenGodTengShe,
	QimenGodTaiYin,
	QimenGodLiuHe,
	QimenGodBaiHu,
	QimenGodXuanWu,
	QimenGodJiuDi,
	QimenGodJiuTian,
}

// Grid 三行三列九宫展示: [巽离坤; 震中兑; 艮坎乾] (按行从左到右)。
var Grid = [3][3]uint8{{4, 9, 2}, {3, 5, 7}, {8, 1, 6}}

// PalaceNames 宫名 (索引 0 占位为空, 1..=9 对应实际宫位)。
var PalaceNames = [10]string{"", "坎", "坤", "震", "巽", "中", "乾", "兑", "艮", "离"}

// TermJu 节气索引 (0..24) → [上元局, 中元局, 下元局]。
//
// 索引对应 tyme4go SolarTermNames: 0冬至, 1小寒, ..., 23大雪。
var TermJu = [24][3]uint8{
	{1, 7, 4}, // 0  冬至
	{2, 8, 5}, // 1  小寒
	{3, 9, 6}, // 2  大寒
	{8, 5, 2}, // 3  立春
	{9, 6, 3}, // 4  雨水
	{1, 7, 4}, // 5  惊蛰
	{3, 9, 6}, // 6  春分
	{4, 1, 7}, // 7  清明
	{5, 2, 8}, // 8  谷雨
	{4, 1, 7}, // 9  立夏
	{5, 2, 8}, // 10 小满
	{6, 3, 9}, // 11 芒种
	{9, 3, 6}, // 12 夏至
	{8, 2, 5}, // 13 小暑
	{7, 1, 4}, // 14 大暑
	{2, 5, 8}, // 15 立秋
	{1, 4, 7}, // 16 处暑
	{9, 3, 6}, // 17 白露
	{7, 1, 4}, // 18 秋分
	{6, 9, 3}, // 19 寒露
	{5, 8, 2}, // 20 霜降
	{6, 9, 3}, // 21 立冬
	{5, 8, 2}, // 22 小雪
	{4, 7, 1}, // 23 大雪
}

// TenXunShou 旬首 (按 tyme4go TenNames 索引):
// 甲子→戊, 甲戌→己, 甲申→庚, 甲午→辛, 甲辰→壬, 甲寅→癸。
var TenXunShou = [6]uint8{4, 5, 6, 7, 8, 9}

// TenXunStartBranch 旬首起始地支 (按旬索引):
// 甲子→子, 甲戌→戌, 甲申→申, 甲午→午, 甲辰→辰, 甲寅→寅。
var TenXunStartBranch = [6]uint8{0, 10, 8, 6, 4, 2}

// TenKongBranches 旬空亡地支对 (按旬索引):
// 甲子→戌亥, 甲戌→申酉, 甲申→午未, 甲午→辰巳, 甲辰→寅卯, 甲寅→子丑。
var TenKongBranches = [6][2]uint8{
	{10, 11},
	{8, 9},
	{6, 7},
	{4, 5},
	{2, 3},
	{0, 1},
}

// BranchToPalace 地支 (0..12) → 落宫: 子1, 丑寅8, 卯3, 辰巳4, 午9, 未申2, 酉7, 戌亥6。
var BranchToPalace = [12]uint8{1, 8, 8, 3, 4, 4, 9, 2, 2, 7, 6, 6}

// StemTombPalace 天干 (0..10) → 入墓宫位: 甲癸→2, 乙丙戊→6, 丁己庚→8, 辛壬→4。
var StemTombPalace = [10]uint8{
	2, // 甲
	6, // 乙
	6, // 丙
	8, // 丁
	6, // 戊
	8, // 己
	8, // 庚
	4, // 辛
	4, // 壬
	2, // 癸
}

// ===================== Plate[T] =====================

// Plate 九宫盘: 索引 0..=8 对应宫位 1..=9 (palace - 1)。
//
// 中宫无值用 nil 表示 (星/门/神盘的中宫为空; 地/天/暗干盘的中宫有值)。
type Plate[T any] struct {
	cells [9]*T
}

// NewPlate 构造一个空 Plate。
func NewPlate[T any]() *Plate[T] { return &Plate[T]{} }

// Set 在指定宫位设置值; palace 必须在 [1, 9] 区间内, 否则 panic。
func (p *Plate[T]) Set(palace uint8, v T) {
	if palace < 1 || palace > 9 {
		panic("palace must be in [1, 9]")
	}
	val := v
	p.cells[palace-1] = &val
}

// Get 取指定宫位的值指针; 越界或未设置返回 nil。
func (p *Plate[T]) Get(palace uint8) *T {
	if palace < 1 || palace > 9 {
		return nil
	}
	return p.cells[palace-1]
}

// ===================== 步进辅助 =====================

// stepPalace 按阴阳遁次序前进一步 (跳过 5 中宫)。
func stepPalace(palace uint8, yy tyme.YinYang) uint8 {
	if yy == tyme.YANG {
		if palace == 9 {
			return 1
		}
		return palace + 1
	}
	if palace == 1 {
		return 9
	}
	return palace - 1
}

// movePalaceBySteps 按阴阳遁连续前进 steps 步, 最终若落 5 宫则改寄 2 宫。
func movePalaceBySteps(palace uint8, steps int, yy tyme.YinYang) uint8 {
	target := palace
	for i := 0; i < steps; i++ {
		target = stepPalace(target, yy)
	}
	if target == 5 {
		return 2
	}
	return target
}

// areOppositePalaces 判断两个宫位是否对冲 (1↔9, 2↔8, 3↔7, 4↔6)。
func areOppositePalaces(a, b uint8) bool {
	switch {
	case a == 1 && b == 9, a == 9 && b == 1:
		return true
	case a == 2 && b == 8, a == 8 && b == 2:
		return true
	case a == 3 && b == 7, a == 7 && b == 3:
		return true
	case a == 4 && b == 6, a == 6 && b == 4:
		return true
	}
	return false
}

// palaceBranchIndices 宫位 → 该宫所统辖的地支索引列表 (按子丑..亥的索引值)。
//
// 中宫 5 返回空切片。
func palaceBranchIndices(palace uint8) []uint8 {
	switch palace {
	case 1:
		return []uint8{0}
	case 2:
		return []uint8{7, 8}
	case 3:
		return []uint8{3}
	case 4:
		return []uint8{4, 5}
	case 6:
		return []uint8{10, 11}
	case 7:
		return []uint8{9}
	case 8:
		return []uint8{1, 2}
	case 9:
		return []uint8{6}
	}
	return nil
}

// isStemInTomb 判断天干是否落入自己的墓库宫。
func isStemInTomb(stem tyme.HeavenStem, palace uint8) bool {
	idx := stem.GetIndex()
	if idx < 0 || idx >= 10 {
		return false
	}
	return StemTombPalace[idx] == palace
}

// tenXunStartBranchIndex 由旬索引取旬首起始地支索引。
func tenXunStartBranchIndex(hour tyme.SixtyCycle) int {
	tenIdx := hour.GetTen().GetIndex()
	if tenIdx > 5 {
		tenIdx = 5
	}
	if tenIdx < 0 {
		tenIdx = 0
	}
	return int(TenXunStartBranch[tenIdx])
}

// ===================== 盘内查找 =====================

// findStemPalace 在天干盘中查找指定天干所在的宫位 (含/不含中宫)。
func findStemPalace(plate *Plate[tyme.HeavenStem], stem tyme.HeavenStem, excludeCenter bool) *uint8 {
	stemIdx := stem.GetIndex()
	for palace := uint8(1); palace <= 9; palace++ {
		if excludeCenter && palace == 5 {
			continue
		}
		v := plate.Get(palace)
		if v != nil && v.GetIndex() == stemIdx {
			p := palace
			return &p
		}
	}
	return nil
}

// findHourStemPalace 求当前时辰天干所在地盘宫位。
//
// 甲干寄旬首符使所在宫; 否则按地盘排找; 若中宫含则寄 2 宫。
func findHourStemPalace(plate *Plate[tyme.HeavenStem], stem tyme.HeavenStem, zhiFuPalace uint8) uint8 {
	if stem.GetIndex() == 0 {
		return zhiFuPalace
	}
	if p := findStemPalace(plate, stem, true); p != nil {
		return *p
	}
	if c := plate.Get(5); c != nil && c.GetIndex() == stem.GetIndex() {
		return 2
	}
	return zhiFuPalace
}

// findDoorPalace 在门盘中查找指定门所在的宫位。
func findDoorPalace(plate *Plate[QimenDoor], door QimenDoor) *uint8 {
	for palace := uint8(1); palace <= 9; palace++ {
		v := plate.Get(palace)
		if v != nil && *v == door {
			p := palace
			return &p
		}
	}
	return nil
}

// ===================== 六个盘构造器 =====================

// buildEarthPlate 构造地盘天干 (三奇六仪)。
func buildEarthPlate(yy tyme.YinYang, ju uint8) *Plate[tyme.HeavenStem] {
	p := NewPlate[tyme.HeavenStem]()
	palace := ju
	for _, stemIdx := range SanQiLiuYi {
		p.Set(palace, tyme.HeavenStem{}.FromIndex(int(stemIdx)))
		palace = stepPalace(palace, yy)
	}
	return p
}

// buildHeavenPlate 构造天盘天干。
//
// 需提前算出 zhiFuPalace (旬首落宫) 与 hourPalace (时干落宫) 复用。
func buildHeavenPlate(earth *Plate[tyme.HeavenStem], yy tyme.YinYang, zhiFuPalace, hourPalace uint8) *Plate[tyme.HeavenStem] {
	p := NewPlate[tyme.HeavenStem]()
	zhiFuIdx := int(LuoShuIndex[zhiFuPalace])
	hourIdx := int(LuoShuIndex[hourPalace])
	var steps int
	if yy == tyme.YANG {
		steps = (hourIdx + 8 - zhiFuIdx) % 8
	} else {
		steps = (zhiFuIdx + 8 - hourIdx) % 8
	}
	for i, palace := range LuoShuOrder {
		if stem := earth.Get(palace); stem != nil {
			var targetIdx int
			if yy == tyme.YANG {
				targetIdx = (i + steps) % 8
			} else {
				targetIdx = (i + 8 - steps) % 8
			}
			p.Set(LuoShuOrder[targetIdx], *stem)
		}
	}
	if center := earth.Get(5); center != nil {
		p.Set(5, *center)
	}
	return p
}

// buildStarPlate 构造九星盘。
func buildStarPlate(zhiFuPalace, hourPalace uint8) *Plate[QimenStar] {
	p := NewPlate[QimenStar]()
	zhiFuIdx := int(LuoShuIndex[zhiFuPalace])
	hourIdx := int(LuoShuIndex[hourPalace])
	steps := (hourIdx + 8 - zhiFuIdx) % 8
	for i, originalPalace := range LuoShuOrder {
		targetPalace := LuoShuOrder[(i+steps)%8]
		var star QimenStar
		if originalPalace == 2 {
			star = QimenStarQinRui
		} else if s := QimenStarFromPalace(originalPalace); s != nil {
			star = *s
		} else {
			star = QimenStarTianRui
		}
		p.Set(targetPalace, star)
	}
	return p
}

// buildDoorPlate 构造八门盘。
func buildDoorPlate(yy tyme.YinYang, zhiFuPalace uint8, hour tyme.SixtyCycle) *Plate[QimenDoor] {
	p := NewPlate[QimenDoor]()
	xunStartBranchIdx := tenXunStartBranchIndex(hour)
	hourBranchIdx := hour.GetEarthBranch().GetIndex()
	branchSteps := (hourBranchIdx + 12 - xunStartBranchIdx) % 12
	zhiShiPalace := movePalaceBySteps(zhiFuPalace, branchSteps, yy)
	zhiFuIdx := int(LuoShuIndex[zhiFuPalace])
	zhiShiIdx := int(LuoShuIndex[zhiShiPalace])
	var steps int
	if yy == tyme.YANG {
		steps = (zhiShiIdx + 8 - zhiFuIdx) % 8
	} else {
		steps = (zhiFuIdx + 8 - zhiShiIdx) % 8
	}
	for i, originalPalace := range LuoShuOrder {
		var targetIdx int
		if yy == tyme.YANG {
			targetIdx = (i + steps) % 8
		} else {
			targetIdx = (i + 8 - steps) % 8
		}
		if d := QimenDoorFromPalace(originalPalace); d != nil {
			p.Set(LuoShuOrder[targetIdx], *d)
		}
	}
	return p
}

// buildGodPlate 构造九神盘。
func buildGodPlate(yy tyme.YinYang, zhiFuPalace uint8) *Plate[QimenGod] {
	p := NewPlate[QimenGod]()
	var order [8]uint8
	if yy == tyme.YANG {
		order = GodYangOrder
	} else {
		order = GodYinOrder
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
	for i, god := range GodsOrder {
		p.Set(order[(startIndex+i)%8], god)
	}
	return p
}

// buildHiddenPlate 构造暗干盘 (从 8 宫起)。
func buildHiddenPlate(yy tyme.YinYang) *Plate[tyme.HeavenStem] {
	p := NewPlate[tyme.HeavenStem]()
	palace := uint8(8)
	for _, stemIdx := range SanQiLiuYi {
		p.Set(palace, tyme.HeavenStem{}.FromIndex(int(stemIdx)))
		palace = stepPalace(palace, yy)
	}
	return p
}

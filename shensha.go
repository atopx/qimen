package qimen

import (
	"fmt"

	"github.com/6tail/tyme4go/tyme"
)

// ShenShaKind 神煞种类。
type ShenShaKind int

const (
	// ShenShaYiMa 驿马 — 主迁移、出行、变动。
	ShenShaYiMa ShenShaKind = iota
	// ShenShaTaoHua 桃花 — 主异性缘、魅力。又名"咸池"。
	ShenShaTaoHua
	// ShenShaHuaGai 华盖 — 主孤高、艺术、宗教、清高。
	ShenShaHuaGai
	// ShenShaTianYi 天乙贵人。
	ShenShaTianYi
	// ShenShaTianDe 天德贵人。
	ShenShaTianDe
	// ShenShaYueDe 月德贵人。
	ShenShaYueDe
	// ShenShaGuoYin 国印贵人。
	ShenShaGuoYin
	// ShenShaWenChang 文昌。
	ShenShaWenChang
	// ShenShaLuShen 禄神。
	ShenShaLuShen
	// ShenShaYangRen 羊刃 — 仅阳干。
	ShenShaYangRen
)

// Name 神煞中文名。
func (k ShenShaKind) Name() string {
	switch k {
	case ShenShaYiMa:
		return "驿马"
	case ShenShaTaoHua:
		return "桃花"
	case ShenShaHuaGai:
		return "华盖"
	case ShenShaTianYi:
		return "天乙贵人"
	case ShenShaTianDe:
		return "天德贵人"
	case ShenShaYueDe:
		return "月德贵人"
	case ShenShaGuoYin:
		return "国印贵人"
	case ShenShaWenChang:
		return "文昌"
	case ShenShaLuShen:
		return "禄神"
	case ShenShaYangRen:
		return "羊刃"
	}
	return ""
}

// Summary 神煞客观描述。
func (k ShenShaKind) Summary() string {
	switch k {
	case ShenShaYiMa:
		return "主迁移、出行、变动"
	case ShenShaTaoHua:
		return "主异性缘、魅力、风流韵事"
	case ShenShaHuaGai:
		return "主孤高、艺术、宗教、清高"
	case ShenShaTianYi:
		return "八字神煞之首,主逢凶化吉、贵人相助"
	case ShenShaTianDe:
		return "上天之德庇佑,主化灾解难"
	case ShenShaYueDe:
		return "月令之德庇佑,主品行端正"
	case ShenShaGuoYin:
		return "主掌权印信、仕途亨通"
	case ShenShaWenChang:
		return "主学业、文采、考试、功名"
	case ShenShaLuShen:
		return "代表俸禄、衣食、财源"
	case ShenShaYangRen:
		return "刚猛之气的极端表现,主血光、刑伤、果敢"
	}
	return ""
}

// Auspice 神煞吉凶定性。
//
//   - 天乙/天德/月德贵人 → 大吉
//   - 国印、文昌、禄神 → 吉
//   - 驿马、桃花、华盖 → 中和
//   - 羊刃 → 凶
func (k ShenShaKind) Auspice() Auspice {
	switch k {
	case ShenShaTianYi, ShenShaTianDe, ShenShaYueDe:
		return AuspiceGreatAuspicious
	case ShenShaGuoYin, ShenShaWenChang, ShenShaLuShen:
		return AuspiceAuspicious
	case ShenShaYiMa, ShenShaTaoHua, ShenShaHuaGai:
		return AuspiceNeutral
	case ShenShaYangRen:
		return AuspiceInauspicious
	}
	return AuspiceNeutral
}

// String 实现 fmt.Stringer。
func (k ShenShaKind) String() string { return k.Name() }

// ShenShaTarget 神煞对应的目标 (天干或地支)。Stem 和 Branch 二选一非 nil。
type ShenShaTarget struct {
	Stem   *tyme.HeavenStem
	Branch *tyme.EarthBranch
}

// String 实现 fmt.Stringer (输出天干名或地支名)。
func (t ShenShaTarget) String() string {
	if t.Stem != nil {
		return t.Stem.GetName()
	}
	if t.Branch != nil {
		return t.Branch.GetName()
	}
	return ""
}

// ShenSha 单个神煞实例。
type ShenSha struct {
	Kind       ShenShaKind
	Target     ShenShaTarget
	PalaceCell uint8
}

// Palace 神煞落入的宫位号 (1..=9)。
func (s ShenSha) Palace() uint8 { return s.PalaceCell }

// Name 委派到 [ShenShaKind.Name]。
func (s ShenSha) Name() string { return s.Kind.Name() }

// Summary 委派到 [ShenShaKind.Summary]。
func (s ShenSha) Summary() string { return s.Kind.Summary() }

// Auspice 委派到 [ShenShaKind.Auspice]。
func (s ShenSha) Auspice() Auspice { return s.Kind.Auspice() }

// String 实现 fmt.Stringer (形如 "驿马(寅→8宫)")。
func (s ShenSha) String() string {
	return fmt.Sprintf("%s(%s→%d宫)", s.Kind.Name(), s.Target.String(), s.PalaceCell)
}

// ===================== 查表常量 =====================

// 驿马: 日/年支 (0..=11) → 驿马地支索引。
var yiMaTable = [12]uint8{2, 11, 8, 5, 2, 11, 8, 5, 2, 11, 8, 5}

// 桃花: 日/年支 → 桃花地支索引。
var taoHuaTable = [12]uint8{9, 6, 3, 0, 9, 6, 3, 0, 9, 6, 3, 0}

// 华盖: 日/年支 → 华盖地支索引 (三合局之墓库)。
var huaGaiTable = [12]uint8{4, 1, 10, 7, 4, 1, 10, 7, 4, 1, 10, 7}

// 天乙贵人: 日干 (0..=9) → 两个贵人地支索引。
var tianYiTable = [10][2]uint8{
	{1, 7},  // 甲 → 丑、未
	{0, 8},  // 乙 → 子、申
	{11, 9}, // 丙 → 亥、酉
	{11, 9}, // 丁 → 亥、酉
	{1, 7},  // 戊 → 丑、未
	{0, 8},  // 己 → 子、申
	{6, 2},  // 庚 → 午、寅
	{6, 2},  // 辛 → 午、寅
	{5, 3},  // 壬 → 巳、卯
	{5, 3},  // 癸 → 巳、卯
}

// 天德贵人: 月支 (0..=11) → (kind, value)。
// kind=0 表示天干; kind=1 表示地支。
var tianDeTable = [12]struct {
	kind  uint8
	value uint8
}{
	{1, 5},  // 子月 → 巳
	{0, 6},  // 丑月 → 庚
	{0, 3},  // 寅月 → 丁
	{1, 8},  // 卯月 → 申
	{0, 8},  // 辰月 → 壬
	{0, 7},  // 巳月 → 辛
	{1, 11}, // 午月 → 亥
	{0, 0},  // 未月 → 甲
	{0, 9},  // 申月 → 癸
	{1, 2},  // 酉月 → 寅
	{0, 2},  // 戌月 → 丙
	{0, 1},  // 亥月 → 乙
}

// 月德贵人: 月支 → 天干索引 (三合局之阳干)。
var yueDeTable = [12]uint8{8, 6, 2, 0, 8, 6, 2, 0, 8, 6, 2, 0}

// 国印贵人: 年干 (0..=9) → 地支索引。
var guoYinTable = [10]uint8{10, 11, 1, 2, 1, 2, 4, 5, 7, 8}

// 文昌: 日/年干 → 地支索引 (天干食神之临官位)。
var wenChangTable = [10]uint8{5, 6, 8, 9, 8, 9, 11, 0, 2, 3}

// 禄神: 日干 → 地支索引 (天干临官位)。
var luShenTable = [10]uint8{2, 3, 5, 6, 5, 6, 8, 9, 11, 0}

// 羊刃: 阳日干 → 地支索引; 阴干为 -1。
var yangRenTable = [10]int{3, -1, 6, -1, 6, -1, 9, -1, 0, -1}

// ===================== 检测函数 =====================

// pushBranch 把地支神煞添加到 out (唯一落宫由 BranchToPalace 决定)。
func pushBranch(out *[]ShenSha, kind ShenShaKind, branchIdx int) {
	if branchIdx < 0 || branchIdx >= 12 {
		return
	}
	branch := tyme.EarthBranch{}.FromIndex(branchIdx)
	*out = append(*out, ShenSha{
		Kind:       kind,
		Target:     ShenShaTarget{Branch: &branch},
		PalaceCell: BranchToPalace[branchIdx],
	})
}

// pushStem 把天干神煞添加到 out (扫描地盘所有宫位, 该天干所在宫各生成一个实例)。
func pushStem(out *[]ShenSha, kind ShenShaKind, stemIdx int, earth *Plate[tyme.HeavenStem]) {
	if stemIdx < 0 || stemIdx >= 10 {
		return
	}
	stem := tyme.HeavenStem{}.FromIndex(stemIdx)
	for palace := uint8(1); palace <= 9; palace++ {
		if v := earth.Get(palace); v != nil && v.GetIndex() == stem.GetIndex() {
			s := stem
			*out = append(*out, ShenSha{
				Kind:       kind,
				Target:     ShenShaTarget{Stem: &s},
				PalaceCell: palace,
			})
		}
	}
}

// detectShenSha 检测全部 10 种神煞。每个落宫各生成一个 [ShenSha] 实例。
//
//   - 地支神煞: 由 BranchToPalace 直接映射 (单宫)
//   - 天干神煞: 扫描地盘三奇六仪所有宫位 (通常单宫)
//   - 天乙贵人有两支 → 生成两个实例
func detectShenSha(yearStem tyme.HeavenStem, monthBranch tyme.EarthBranch, dayStem tyme.HeavenStem, dayBranch tyme.EarthBranch, earth *Plate[tyme.HeavenStem]) []ShenSha {
	out := make([]ShenSha, 0, 12)

	dayBranchIdx := dayBranch.GetIndex()
	dayStemIdx := dayStem.GetIndex()
	monthBranchIdx := monthBranch.GetIndex()
	yearStemIdx := yearStem.GetIndex()

	if dayBranchIdx >= 0 && dayBranchIdx < 12 {
		pushBranch(&out, ShenShaYiMa, int(yiMaTable[dayBranchIdx]))
		pushBranch(&out, ShenShaTaoHua, int(taoHuaTable[dayBranchIdx]))
		pushBranch(&out, ShenShaHuaGai, int(huaGaiTable[dayBranchIdx]))
	}

	if dayStemIdx >= 0 && dayStemIdx < 10 {
		pair := tianYiTable[dayStemIdx]
		pushBranch(&out, ShenShaTianYi, int(pair[0]))
		pushBranch(&out, ShenShaTianYi, int(pair[1]))

		pushBranch(&out, ShenShaWenChang, int(wenChangTable[dayStemIdx]))
		pushBranch(&out, ShenShaLuShen, int(luShenTable[dayStemIdx]))
		if idx := yangRenTable[dayStemIdx]; idx >= 0 {
			pushBranch(&out, ShenShaYangRen, idx)
		}
	}

	if monthBranchIdx >= 0 && monthBranchIdx < 12 {
		entry := tianDeTable[monthBranchIdx]
		if entry.kind == 0 {
			pushStem(&out, ShenShaTianDe, int(entry.value), earth)
		} else {
			pushBranch(&out, ShenShaTianDe, int(entry.value))
		}
		pushStem(&out, ShenShaYueDe, int(yueDeTable[monthBranchIdx]), earth)
	}

	if yearStemIdx >= 0 && yearStemIdx < 10 {
		pushBranch(&out, ShenShaGuoYin, int(guoYinTable[yearStemIdx]))
	}

	return out
}

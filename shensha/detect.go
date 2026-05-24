package shensha

import (
	"iter"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/internal/tables"
)

// DetectInput is the chart view used by Detect.
//
// EarthStems is the 地盘 三奇六仪 layout (palace - 1 indexed) used to find
// stem-anchored 神煞 (天德/月德 when keyed by stem).
type DetectInput struct {
	YearStem    almanac.Stem
	MonthBranch almanac.Branch
	DayStem     almanac.Stem
	DayBranch   almanac.Branch

	EarthStems [9]almanac.Stem
}

// Detect yields every 神煞 instance landing in the chart.
func Detect(in DetectInput) iter.Seq[ShenSha] {
	return func(yield func(ShenSha) bool) {
		dayBranchIdx := in.DayBranch.Index()
		dayStemIdx := in.DayStem.Index()
		monthBranchIdx := in.MonthBranch.Index()
		yearStemIdx := in.YearStem.Index()

		// Branch-anchored shensha from day branch.
		if !pushBranch(yield, YiMa, int(yiMaTable[dayBranchIdx])) {
			return
		}
		if !pushBranch(yield, TaoHua, int(taoHuaTable[dayBranchIdx])) {
			return
		}
		if !pushBranch(yield, HuaGai, int(huaGaiTable[dayBranchIdx])) {
			return
		}

		// Stem-anchored shensha from day stem.
		pair := tianYiTable[dayStemIdx]
		if !pushBranch(yield, TianYi, int(pair[0])) {
			return
		}
		if !pushBranch(yield, TianYi, int(pair[1])) {
			return
		}
		if !pushBranch(yield, WenChang, int(wenChangTable[dayStemIdx])) {
			return
		}
		if !pushBranch(yield, LuShen, int(luShenTable[dayStemIdx])) {
			return
		}
		if idx := yangRenTable[dayStemIdx]; idx >= 0 {
			if !pushBranch(yield, YangRen, idx) {
				return
			}
		}

		// Month-anchored 天德/月德.
		entry := tianDeTable[monthBranchIdx]
		if entry.kind == 0 {
			if !pushStem(yield, TianDe, int(entry.value), &in.EarthStems) {
				return
			}
		} else {
			if !pushBranch(yield, TianDe, int(entry.value)) {
				return
			}
		}
		if !pushStem(yield, YueDe, int(yueDeTable[monthBranchIdx]), &in.EarthStems) {
			return
		}

		// Year-anchored 国印贵人.
		if !pushBranch(yield, GuoYin, int(guoYinTable[yearStemIdx])) {
			return
		}
	}
}

// pushBranch yields a single branch-anchored shensha.
// Precondition: branchIdx ∈ [0, 11] (always true for table-sourced values).
func pushBranch(yield func(ShenSha) bool, kind Kind, branchIdx int) bool {
	branch := almanac.BranchOf(branchIdx)
	return yield(ShenSha{
		Kind:   kind,
		Target: Target{Branch: branch},
		Palace: tables.BranchToPalace[branchIdx],
	})
}

// pushStem yields one shensha instance per palace where the target
// stem appears in the EarthStems layout.
// Precondition: stemIdx ∈ [0, 9] (always true for table-sourced values).
func pushStem(yield func(ShenSha) bool, kind Kind, stemIdx int, earth *[9]almanac.Stem) bool {
	stem := almanac.StemOf(stemIdx)
	for palace := uint8(1); palace <= 9; palace++ {
		if earth[palace-1].Index() == stem.Index() {
			if !yield(ShenSha{
				Kind:   kind,
				Target: Target{Stem: stem, IsStem: true},
				Palace: palace,
			}) {
				return false
			}
		}
	}
	return true
}

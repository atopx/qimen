// Package compute exposes the qimen layout-derivation primitives
// (yin/yang, yuan, ju, xun-shou, kong-wang).
//
// All operations are total over their closed-domain inputs (Term, Cycle,
// Yuan, etc.). Callers must validate Method / Style at the Chart entry
// point — these primitives assume MethodTime semantics.
package compute

import (
	"errors"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/internal/tables"
)

// ErrUnsupportedMethod is returned when a caller asks Chart to use a
// Method that is not yet implemented (MethodDay/Month/Year reserved).
//
// This is a Chart-level concern; compute primitives never validate it.
var ErrUnsupportedMethod = errors.New("qimen: unsupported method")

// YinYang returns the 阴/阳 遁 polarity for a solar instant.
//
//	冬至 ≤ t < 夏至 → 阳遁
//	夏至 ≤ t < 次冬至 → 阴遁
func YinYang(st almanac.SolarTime) almanac.YinYang {
	year := int(st.Year)
	winter := almanac.TermOf(year, 0).SolarTime()
	summer := almanac.TermOf(year, 12).SolarTime()
	nextWinter := almanac.TermOf(year+1, 0).SolarTime()
	if (!st.Before(winter) && st.Before(summer)) || !st.Before(nextWinter) {
		return almanac.Yang
	}
	return almanac.Yin
}

// Yuan returns the 三元 segment for a given day pillar.
// Index mod 15 lands in [0..4]→上元, [5..9]→中元, [10..14]→下元.
func Yuan(day almanac.Cycle) enum.Yuan {
	switch day.Index() % 15 {
	case 0, 1, 2, 3, 4:
		return enum.YuanUpper
	case 5, 6, 7, 8, 9:
		return enum.YuanMiddle
	}
	return enum.YuanLower
}

// Ju returns the local 局 number (1..9) from a solar term + 元.
//
// Both term.Index() ∈ [0, 23] and yuan ∈ [0, 2] are closed-domain;
// this is a total function on its inputs.
func Ju(term almanac.Term, yuan enum.Yuan) uint8 {
	return tables.TermJu[term.Index()][yuan]
}

// XunShou returns the 旬首 stem for the hour pillar (戊/己/庚/辛/壬/癸).
// hour.Ten() is always 0..5 by construction.
func XunShou(hour almanac.Cycle) almanac.Stem {
	return almanac.StemOf(int(tables.TenXunShou[hour.Ten().Index()]))
}

// KongWang returns the 旬空亡 branch pair for the hour pillar.
func KongWang(hour almanac.Cycle) [2]almanac.Branch {
	pair := tables.TenKongBranches[hour.Ten().Index()]
	return [2]almanac.Branch{
		almanac.BranchOf(int(pair[0])),
		almanac.BranchOf(int(pair[1])),
	}
}

// palaceBranchesCache holds the precomputed []almanac.Branch slices for
// each palace 1..9. Built once at package init from tables.PalaceBranches;
// callers MUST NOT mutate the returned slice.
var palaceBranchesCache [10][]almanac.Branch

func init() {
	for palace := uint8(1); palace <= 9; palace++ {
		idxs := tables.PalaceBranches[palace]
		if len(idxs) == 0 {
			continue
		}
		out := make([]almanac.Branch, len(idxs))
		for i, idx := range idxs {
			out[i] = almanac.BranchOf(int(idx))
		}
		palaceBranchesCache[palace] = out
	}
}

// BranchesForPalace returns the (shared, read-only) 地支 list for a palace.
// Precondition: palace ∈ [1, 9]. Palace 5 returns nil (center has no branches).
func BranchesForPalace(palace uint8) []almanac.Branch {
	return palaceBranchesCache[palace]
}

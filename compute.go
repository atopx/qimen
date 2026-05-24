package qimen

import (
	"github.com/6tail/tyme4go/tyme"
)

// computeYinYang 计算阴阳遁。
//
// 用冬至/夏至前后划分阴阳遁:
//   - 冬至 (含) 至夏至 (不含): 阳遁
//   - 夏至 (含) 至次年冬至 (不含): 阴遁
func computeYinYang(st tyme.SolarTime) tyme.YinYang {
	year := st.GetYear()
	winter := tyme.SolarTerm{}.FromIndex(year, 0).GetJulianDay().GetSolarTime()
	summer := tyme.SolarTerm{}.FromIndex(year, 12).GetJulianDay().GetSolarTime()
	nextWinter := tyme.SolarTerm{}.FromIndex(year+1, 0).GetJulianDay().GetSolarTime()
	// (!is_before(winter) && is_before(summer)) || !is_before(next_winter)
	if (!st.IsBefore(winter) && st.IsBefore(summer)) || !st.IsBefore(nextWinter) {
		return tyme.YANG
	}
	return tyme.YIN
}

// computeYuan 由日柱六十甲子索引计算三元 (上/中/下元)。
//
// 索引 mod 15 落入 [0,4]→上元, [5,9]→中元, [10,14]→下元。
func computeYuan(day tyme.SixtyCycle) QimenYuan {
	switch day.GetIndex() % 15 {
	case 0, 1, 2, 3, 4:
		return QimenYuanUpper
	case 5, 6, 7, 8, 9:
		return QimenYuanMiddle
	}
	return QimenYuanLower
}

// computeJu 由节气与三元计算局数 (1..=9)。
//
// 节气索引越界返回 ErrCodeUnsupportedTerm。
func computeJu(term tyme.SolarTerm, yuan QimenYuan) (uint8, error) {
	idx := term.GetIndex()
	if idx < 0 || idx >= len(TermJu) {
		return 0, newUnsupportedTerm(term.GetName())
	}
	row := TermJu[idx]
	switch yuan {
	case QimenYuanUpper:
		return row[0], nil
	case QimenYuanMiddle:
		return row[1], nil
	case QimenYuanLower:
		return row[2], nil
	}
	return 0, newUnsupportedTerm(term.GetName())
}

// computeXunShou 由时柱六十甲子取旬首天干 (戊/己/庚/辛/壬/癸 之一)。
func computeXunShou(hour tyme.SixtyCycle) tyme.HeavenStem {
	tenIdx := hour.GetTen().GetIndex()
	if tenIdx > 5 {
		tenIdx = 5
	}
	if tenIdx < 0 {
		tenIdx = 0
	}
	return tyme.HeavenStem{}.FromIndex(int(TenXunShou[tenIdx]))
}

// computeKongWang 由时柱六十甲子取旬空亡两支。
func computeKongWang(hour tyme.SixtyCycle) [2]tyme.EarthBranch {
	tenIdx := hour.GetTen().GetIndex()
	if tenIdx > 5 {
		tenIdx = 5
	}
	if tenIdx < 0 {
		tenIdx = 0
	}
	pair := TenKongBranches[tenIdx]
	return [2]tyme.EarthBranch{
		tyme.EarthBranch{}.FromIndex(int(pair[0])),
		tyme.EarthBranch{}.FromIndex(int(pair[1])),
	}
}

// branchesForPalace 取宫位所统辖的地支列表。
func branchesForPalace(palace uint8) []tyme.EarthBranch {
	idxs := palaceBranchIndices(palace)
	out := make([]tyme.EarthBranch, 0, len(idxs))
	for _, i := range idxs {
		out = append(out, tyme.EarthBranch{}.FromIndex(int(i)))
	}
	return out
}

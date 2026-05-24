package almanac

import (
	"fmt"
	"iter"
	"math"
	"sync"
)

// LunarYear is one Chinese lunar year, identified by the Gregorian year
// that contains its 春节 (1st day of 正月).
type LunarYear struct {
	Year      int   // Gregorian year containing 春节
	LeapMonth uint8 // 1..12 if year has a leap month, else 0
}

// LunarYearOf returns metadata for the lunar year corresponding to a
// given Gregorian year.
func LunarYearOf(year int) LunarYear {
	cache := lunarYearCache(year)
	return LunarYear{Year: year, LeapMonth: cache.leapMonth}
}

// Name returns the localized name, e.g. "农历甲辰年".
func (y LunarYear) Name() string {
	cycle := CycleOf(y.Year - 4)
	return fmt.Sprintf("农历%s年", cycle.Name())
}

// String implements fmt.Stringer.
func (y LunarYear) String() string { return y.Name() }

// Cycle returns the sexagenary cycle of the lunar year.
func (y LunarYear) Cycle() Cycle { return CycleOf(y.Year - 4) }

// MonthCount returns 12 (normal) or 13 (leap year).
func (y LunarYear) MonthCount() int {
	if y.LeapMonth == 0 {
		return 12
	}
	return 13
}

// Months iterates lunar months in calendar order. Pushes a leap month
// immediately after the month whose number matches y.LeapMonth.
func (y LunarYear) Months() iter.Seq[LunarMonth] {
	return func(yield func(LunarMonth) bool) {
		for m := uint8(1); m <= 12; m++ {
			if !yield(LunarMonth{Year: y, Month: m}) {
				return
			}
			if m == y.LeapMonth {
				if !yield(LunarMonth{Year: y, Month: m, Leap: true}) {
					return
				}
			}
		}
	}
}

// --- internal cache: leap-month + 12-month first-day Julian Days ---

type lunarYearInfo struct {
	leapMonth uint8
	// firstDays[i] = JD offset (from j2000, noon) of 1st day of month
	// at index i (0..11 or 0..12 with leap inserted after leapMonth-1).
	firstDays []int64
	// monthNumbers parallel to firstDays, holds month number (1..12).
	monthNumbers []uint8
	// leapMask[i]=true if month index i is the leap month.
	leapMask []bool
}

var (
	lunarYearCacheMu sync.RWMutex
	lunarYearMap     = map[int]*lunarYearInfo{}
)

func lunarYearCache(year int) *lunarYearInfo {
	lunarYearCacheMu.RLock()
	info, ok := lunarYearMap[year]
	lunarYearCacheMu.RUnlock()
	if ok {
		return info
	}
	lunarYearCacheMu.Lock()
	defer lunarYearCacheMu.Unlock()
	if info, ok = lunarYearMap[year]; ok {
		return info
	}
	info = computeLunarYearInfo(year)
	lunarYearMap[year] = info
	return info
}

// computeLunarYearInfo derives leap-month + month-start days for a year.
//
// Algorithm (GB/T 33661-2017):
//  1. Find 冬至 of year-1 (prevDZ) and 冬至 of year (currDZ).
//  2. Walk forward through 朔's starting before prevDZ until we have at least
//     16 朔's and cover ~60 days past currDZ — guarantees enough span for
//     up to 13 months of lunar year y plus a buffer.
//  3. The month containing prevDZ is "11月 of (y-1)"; the month containing
//     currDZ is "11月 of y". Months between them: 12 (normal) or 13 (leap).
//  4. If leap, the first month between the two 11月's that contains no 中气
//     (even-indexed solar term) is the leap month; it takes the previous
//     month's number.
//  5. 正月 of y is the 3rd 朔 after 11月 of (y-1), or 4th if leap fell before it.
func computeLunarYearInfo(year int) *lunarYearInfo {
	// TermOf(year, 0) returns the 冬至 in calendar year (year-1),
	// because the canonical 寿星 cycle starts at 冬至 which physically
	// occurs at the end of the prior calendar year. So:
	//   - 冬至 in calendar year (year-1) → TermOf(year, 0)
	//   - 冬至 in calendar year (year)   → TermOf(year+1, 0)
	prevDZ := math.Floor(TermOf(year, 0).JulianDay() + 0.5)
	currDZ := math.Floor(TermOf(year+1, 0).JulianDay() + 0.5)

	// Collect ~18 朔's starting before prevDZ.
	shuos := []float64{calcShuo(prevDZ - j2000 - 30)}
	for len(shuos) < 18 || shuos[len(shuos)-1]+j2000 < currDZ+90 {
		nextProbe := shuos[len(shuos)-1] + 30
		next := calcShuo(nextProbe)
		if next <= shuos[len(shuos)-1]+1 {
			// guard against numerical stagnation
			next = shuos[len(shuos)-1] + 29.5306
		}
		shuos = append(shuos, next)
		if len(shuos) > 30 {
			break
		}
	}

	// Index of 朔 that opens the lunar month containing prevDZ.
	month11Idx := -1
	for i := 0; i+1 < len(shuos); i++ {
		if shuos[i]+j2000 <= prevDZ && shuos[i+1]+j2000 > prevDZ {
			month11Idx = i
			break
		}
	}
	if month11Idx < 0 {
		return &lunarYearInfo{leapMonth: 0}
	}

	// Index of 朔 that opens the lunar month containing currDZ.
	month11NextIdx := -1
	for i := month11Idx + 1; i+1 < len(shuos); i++ {
		if shuos[i]+j2000 <= currDZ && shuos[i+1]+j2000 > currDZ {
			month11NextIdx = i
			break
		}
	}
	if month11NextIdx < 0 {
		return &lunarYearInfo{leapMonth: 0}
	}

	span := month11NextIdx - month11Idx
	hasLeap := span == 13

	// Find leap month: first month in (month11Idx, month11NextIdx] with no 中气.
	leapShuoIdx := -1
	if hasLeap {
		zhongqiJDs := collectZhongqi(prevDZ, currDZ+30)
		for i := month11Idx + 1; i <= month11NextIdx; i++ {
			start := shuos[i] + j2000
			end := shuos[i+1] + j2000
			hasZQ := false
			for _, z := range zhongqiJDs {
				if z >= start && z < end {
					hasZQ = true
					break
				}
			}
			if !hasZQ {
				leapShuoIdx = i
				break
			}
		}
	}

	// 正月 starts at month11Idx + 2 (or +3 if leap is between).
	zhengyueIdx := month11Idx + 2
	if hasLeap && leapShuoIdx > month11Idx && leapShuoIdx <= zhengyueIdx {
		zhengyueIdx = month11Idx + 3
	}

	monthCount := 12
	leapMonth := uint8(0)
	if hasLeap && leapShuoIdx >= zhengyueIdx && leapShuoIdx < zhengyueIdx+13 {
		var n uint8 = 0
		for i := zhengyueIdx; i < leapShuoIdx; i++ {
			n++
		}
		leapMonth = n
		monthCount = 13
	}

	// Bounds check: ensure shuos has zhengyueIdx + monthCount entries.
	if zhengyueIdx+monthCount > len(shuos) {
		// Truncate to what we have
		monthCount = len(shuos) - zhengyueIdx
		if monthCount < 12 {
			monthCount = 12
		}
		if zhengyueIdx+monthCount > len(shuos) {
			return &lunarYearInfo{leapMonth: 0}
		}
	}

	info := &lunarYearInfo{
		leapMonth:    leapMonth,
		firstDays:    make([]int64, monthCount),
		monthNumbers: make([]uint8, monthCount),
		leapMask:     make([]bool, monthCount),
	}

	if hasLeap && leapMonth > 0 {
		var monthNum uint8 = 0
		for out := 0; out < monthCount; out++ {
			i := zhengyueIdx + out
			info.firstDays[out] = int64(math.Floor(shuos[i] + j2000 + 0.5))
			if i == leapShuoIdx {
				info.monthNumbers[out] = monthNum
				info.leapMask[out] = true
			} else {
				monthNum++
				info.monthNumbers[out] = monthNum
				info.leapMask[out] = false
			}
		}
	} else {
		for i := 0; i < monthCount; i++ {
			info.firstDays[i] = int64(math.Floor(shuos[zhengyueIdx+i] + j2000 + 0.5))
			info.monthNumbers[i] = uint8(i + 1)
		}
	}

	return info
}

// collectZhongqi returns JD-floor of the 12 中气 (even-indexed terms)
// whose absolute JD falls in [start, end].
func collectZhongqi(start, end float64) []float64 {
	var out []float64
	// start/end are absolute JDs; convert to year estimate via J2000 anchor.
	startYear := int(math.Floor((start-j2000)/365.25 + 2000))
	endYear := int(math.Floor((end-j2000)/365.25 + 2000))
	// Sweep one extra year on each side to be safe (TermOf year convention
	// places 冬至 of calendar Y at TermOf(Y+1, 0)).
	for y := startYear - 1; y <= endYear+2; y++ {
		for i := 0; i < 24; i += 2 {
			jd := math.Floor(TermOf(y, i).JulianDay() + 0.5)
			if jd >= start && jd <= end {
				out = append(out, jd)
			}
		}
	}
	return out
}

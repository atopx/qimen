package almanac

import "math"

// Pillars holds the four sexagenary pillars (year/month/day/hour)
// derived from a solar instant using term-direct rules:
//
//   - Year:  立春 (term index 3) marks the year switch. (Y - 4) mod 60
//     gives 甲子 numbering from year 4 AD.
//   - Month: Each 节 (odd term) starts a new month. 五虎遁年起月 gives
//     the stem of 寅月 from the year stem.
//   - Day:   2000-01-07 noon UT = 甲子日 (JD 2451551). Hour 23 rolls
//     forward by 1 day.
//   - Hour:  五鼠遁日起时. Two-hour 时辰 from 23-01=子.
type Pillars struct {
	Year, Month, Day, Hour Cycle
}

// PillarsOf computes the four pillars for a solar instant.
func PillarsOf(s SolarTime) Pillars {
	yearCycle := yearPillar(s)
	dayCycle := dayPillar(s)
	monthCycle := monthPillar(s, yearCycle)
	hourCycle := hourPillarFromDay(s, dayCycle)
	return Pillars{
		Year:  yearCycle,
		Month: monthCycle,
		Day:   dayCycle,
		Hour:  hourCycle,
	}
}

// Pillars returns the cached pillars for s. Convenience wrapper.
func (s SolarTime) Pillars() Pillars { return PillarsOf(s) }

// Term returns the current solar term for s.
func (s SolarTime) Term() Term { return TermOfSolarTime(s) }

// yearPillar uses 立春 (term index 3) as the year boundary.
// Year is (Y - 4) mod 60 for instants ≥ 立春 of Y, else (Y - 5) mod 60.
func yearPillar(s SolarTime) Cycle {
	y := int(s.Year)
	lichun := TermOf(y, 3).SolarTime()
	if s.Before(lichun) {
		y--
	}
	return CycleOf(y - 4)
}

// monthPillar uses 节气-based month boundaries (节 = odd-index terms).
// Month sequence starts at 寅 (after 立春). 五虎遁: month stem of 寅月 =
// (year_stem * 2 + 2) mod 10, hence sixty-cycle index at seq 0 =
// ((year_stem mod 5) * 12 + 2) mod 60. Each subsequent month adds 1.
func monthPillar(s SolarTime, year Cycle) Cycle {
	// Determine month sequence: 0 = 寅月 (从立春起), 1 = 卯月 (从惊蛰起),
	// ..., 11 = 丑月 (从小寒起).
	seq := monthSequence(s)
	yearStem := int(year.Stem())
	idx0 := ((yearStem%5)*12 + 2) % 60
	return CycleOf(idx0 + seq)
}

// monthSequence returns the 0..11 月柱 offset since 寅月.
//
// Derived in O(1) from the current solar term: each pair of adjacent
// terms (one 节 + one 气) is exactly one month, with 寅月 starting at
// 立春 (term index 3). So seq = ((termIdx - 3 + 24) % 24) / 2.
func monthSequence(s SolarTime) int {
	t := TermOfSolarTime(s)
	return ((t.Index() - 3 + 24) % 24) / 2
}

// dayPillar uses JD 2451551 (2000-01-07 noon UT) as 甲子日 anchor.
// Hour 23 rolls to next day.
func dayPillar(s SolarTime) Cycle {
	noonJD := julianDayFromYmdhms(int(s.Year), int(s.Month), int(s.Day), 12, 0, 0)
	idx := int(math.Floor(noonJD+0.5)) - 2451551
	if s.Hour == 23 {
		idx++
	}
	return CycleOf(idx)
}

// hourPillarFromDay derives the hour pillar from the day pillar.
// At day stem ds, 子时 (hour branch 0) has stem (ds*2) mod 10, so the
// hour-pillar sixty-cycle index at 子时 is ((ds mod 5) * 12) mod 60.
// Each subsequent 时辰 adds 1 to the cycle index.
func hourPillarFromDay(s SolarTime, day Cycle) Cycle {
	branchIdx := hourBranchIndex(int(s.Hour))
	dayStem := int(day.Stem())
	idx0 := ((dayStem % 5) * 12) % 60
	return CycleOf(idx0 + branchIdx)
}

// hourBranchIndex maps 0..23 hours to 0..11 branches.
// 23, 0 → 0 (子); 1, 2 → 1 (丑); ...; 21, 22 → 11 (亥).
func hourBranchIndex(h int) int {
	return ((h + 1) / 2) % 12
}

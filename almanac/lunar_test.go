package almanac

import "testing"

// TestLunarLeapYears verifies leap-month detection for known leap years
// (verified against the 中科院 lunar calendar reference).
//
//   - 2020 闰四月
//   - 2023 闰二月
//   - 2025 闰六月
//   - 2028 闰五月
func TestLunarLeapYears(t *testing.T) {
	cases := []struct {
		year      int
		leapMonth uint8
	}{
		{2020, 4},
		{2023, 2},
		{2025, 6},
		{2028, 5},
	}
	for _, c := range cases {
		y := LunarYearOf(c.year)
		if y.LeapMonth != c.leapMonth {
			t.Errorf("year %d: leapMonth=%d, want %d", c.year, y.LeapMonth, c.leapMonth)
		}
	}
}

// TestLunarKnownDates anchors a few well-known lunar dates.
func TestLunarKnownDates(t *testing.T) {
	cases := []struct {
		desc        string
		y, m, d     int
		wantYear    int
		wantMonth   uint8
		wantLeap    bool
		wantDayName string
	}{
		{"春节 2025", 2025, 1, 29, 2025, 1, false, "初一"},
		{"中秋 2025", 2025, 10, 6, 2025, 8, false, "十五"},
		{"闰六月 2025 初一", 2025, 7, 25, 2025, 6, true, "初一"},
	}
	for _, c := range cases {
		st, _ := SolarTimeOf(c.y, c.m, c.d, 12, 0, 0)
		ld := st.LunarDay()
		if ld.Month.Year.Year != c.wantYear ||
			ld.Month.Month != c.wantMonth ||
			ld.Month.Leap != c.wantLeap ||
			ld.Name() != c.wantDayName {
			t.Errorf("%s: got (y=%d m=%d leap=%v d=%s), want (y=%d m=%d leap=%v d=%s)",
				c.desc,
				ld.Month.Year.Year, ld.Month.Month, ld.Month.Leap, ld.Name(),
				c.wantYear, c.wantMonth, c.wantLeap, c.wantDayName)
		}
	}
}

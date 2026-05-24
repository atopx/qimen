// Package qimen builds 奇门遁甲 (Qimen Dunjia) charts from a solar
// instant. The package is fully self-contained: no external dependencies.
//
// # Quick start
//
//	// Build a chart from a solar instant
//	chart := qimen.New()                       // current time, default options
//	chart, _ := qimen.FromTimestamp(1735632000) // Unix seconds (UTC+8)
//	t, _ := time.Parse("2006-01-02 15:04", "2026-01-14 18:45")
//	chart, _ := qimen.FromTime(t)              // standard library time.Time
//	st, _ := almanac.SolarTimeOf(2026, 1, 14, 18, 45, 0)
//	chart, _ := qimen.From(st)                 // first-class SolarTime
//
//	// Iterate the 9 palaces
//	for n, p := range chart.Palaces() {
//	    fmt.Printf("%d宫 %s — 地盘:%s 天盘:%s\n",
//	        n, p.Name, p.EarthStem.Name(), p.HeavenStem.Name())
//	}
//
//	// Stream patterns + shensha
//	for p := range chart.Patterns() {
//	    fmt.Printf("格局 %s [%s]\n", p.Name(), p.Auspice().Name())
//	}
//
// # Options
//
//   - [WithMethod] selects the 起局法门 (currently only [enum.MethodTime])
//   - [WithStyle] selects the 盘式 (currently only [enum.StyleRotate])
//
// # Error sentinels
//
// Errors returned from chart construction wrap one of [ErrUnsupportedMethod],
// [ErrUnsupportedStyle], [ErrUnsupportedTerm], or [ErrInvalidTime] —
// callers can match with errors.Is.
package qimen

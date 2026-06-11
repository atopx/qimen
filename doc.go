// Package qimen builds 奇门遁甲 (Qimen Dunjia) charts from a solar
// instant. The package is fully self-contained: no external dependencies.
//
// # Quick start
//
//	// Build a chart from a solar instant
//	chart := qimen.New()                        // current time, default options
//	chart, _ := qimen.FromTimestamp(1735632000) // Unix seconds (UTC+8 wall clock)
//	t, _ := time.Parse("2006-01-02 15:04", "2026-01-14 18:45")
//	chart, _ := qimen.FromTime(t)               // standard library time.Time
//	st, _ := almanac.SolarTimeOf(2026, 1, 14, 18, 45, 0)
//	chart, _ := qimen.From(st)                  // first-class SolarTime
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
// # Time zones
//
// Charts are cast in the local civil (wall clock) time of the event.
// [FromTime] uses the wall clock of the supplied time.Time as-is;
// [FromTimestamp] and [New] interpret the instant in UTC+8 (China
// Standard Time). For another zone, convert before calling:
// FromTime(time.Unix(unix, 0).In(loc)).
//
// # Options
//
//   - [WithMethod] selects the 起局法门 (currently only [enum.MethodTime])
//   - [WithStyle] selects the 盘式 (currently only [enum.StyleRotate])
//   - [WithJuRule] selects the 定局规则: 拆补 (default) or 置闰
//
// # Error sentinels
//
// Errors returned from chart construction wrap one of
// [ErrUnsupportedMethod], [ErrUnsupportedStyle], or [ErrInvalidTime] —
// callers can match with errors.Is.
//
// # Conventions
//
// The layout follows mainstream 时家转盘 practice; points where schools
// diverge are fixed as follows:
//
//   - 三元 is derived from the day pillar 符头 grid (index mod 15). The
//     局 is keyed to the astronomical term (拆补, default) or, with
//     [WithJuRule], to the solstice-aligned 置闰 schedule (超神 / 接气,
//     intercalating 芒种 / 大雪 when the leader's advance reaches nine
//     days counted inclusively); [Chart.JuTerm] exposes the term used.
//   - The day pillar switches at 23:00 (晚子时 counts as the next day).
//   - The center-palace stem stays in place on the heaven plate (palace
//     2 is not rendered with a second 寄宫 stem).
package qimen

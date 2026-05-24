package almanac

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

// ErrUnsupportedTerm signals an out-of-range solar term index or name.
var ErrUnsupportedTerm = errors.New("almanac: unsupported solar term")

// termNames lists the 24 solar terms in the canonical 寿星 order, starting
// from 冬至 (winter solstice). Index 0..23 maps directly to a Term.Index.
var termNames = [24]string{
	"冬至", "小寒", "大寒",
	"立春", "雨水", "惊蛰",
	"春分", "清明", "谷雨",
	"立夏", "小满", "芒种",
	"夏至", "小暑", "大暑",
	"立秋", "处暑", "白露",
	"秋分", "寒露", "霜降",
	"立冬", "小雪", "大雪",
}

// Term identifies one of the 24 solar terms occurring in a given Gregorian year.
// Term is an immutable value type — methods never mutate the receiver.
type Term struct {
	year             int     // Gregorian year that contains the term
	index            uint8   // 0..23 (冬至=0, 大雪=23)
	cursoryJulianDay float64 // calendar-rounded JD (noon precision)
	julianDay        float64 // second-precision JD (cached on construction)
}

type termKey struct {
	year  int
	index int
}

// termCache memoizes Term values across calls. Each Term entry already
// holds the expensive julianDay precomputation, so a hit is purely a map
// lookup with no further VSOP / qiAccurate work.
var termCache sync.Map

// TermOf returns the (year, index) solar term. index is wrapped modulo 24,
// with the year shifted forward when wrapping past 大雪.
func TermOf(year, index int) Term {
	size := 24
	y := (year*size + index) / size
	idx := ((index % size) + size) % size
	key := termKey{y, idx}
	if v, ok := termCache.Load(key); ok {
		return v.(Term)
	}
	cjd := initTermJD(y, idx)
	t := Term{
		year:             y,
		index:            uint8(idx),
		cursoryJulianDay: cjd,
		julianDay:        qiAccurate2(cjd) + j2000,
	}
	termCache.Store(key, t)
	return t
}

// TermOfName returns the term in the given year with the given Chinese name.
// Returns ErrUnsupportedTerm when name is unknown.
func TermOfName(year int, name string) (Term, error) {
	for i, n := range termNames {
		if n == name {
			return TermOf(year, i), nil
		}
	}
	return Term{}, fmt.Errorf("%w: %q", ErrUnsupportedTerm, name)
}

// initTermJD computes a coarse Julian Day for the term's start (≈noon)
// based on the calcQi table search.
func initTermJD(year, offset int) float64 {
	jd := math.Floor(float64(year-2000)*365.2422 + 180)
	// 355 ≈ JD offset of 2000-12-21 winter solstice (J2000 = 2451545)
	w := math.Floor((jd-355+183)/365.2422)*365.2422 + 355
	if calcQi(w) > jd {
		w -= 365.2422
	}
	return calcQi(w + 15.2184*float64(offset))
}

// Year returns the Gregorian year that contains this term.
func (t Term) Year() int { return t.year }

// Index returns the term's index in canonical order (0=冬至 .. 23=大雪).
func (t Term) Index() int { return int(t.index) }

// Name returns the Chinese name of the term.
func (t Term) Name() string { return termNames[t.index] }

// IsJie reports whether the term is a "节" (odd index — 小寒, 立春, 惊蛰, ...).
func (t Term) IsJie() bool { return t.index%2 == 1 }

// IsQi reports whether the term is a "气" (even index — 冬至, 大寒, 雨水, ...).
func (t Term) IsQi() bool { return t.index%2 == 0 }

// CursoryJulianDay returns the calendar-rounded Julian Day (noon-precision).
func (t Term) CursoryJulianDay() float64 { return t.cursoryJulianDay }

// JulianDay returns the second-precision Julian Day of the term's start
// (cached at construction).
func (t Term) JulianDay() float64 { return t.julianDay }

// SolarTime returns the solar instant when this term begins.
func (t Term) SolarTime() SolarTime {
	return solarFromJulianDay(t.julianDay)
}

// Next returns the term n steps ahead (negative = backward).
func (t Term) Next(n int) Term {
	idx := int(t.index) + n
	size := 24
	y := (t.year*size + idx) / size
	wrapped := ((idx % size) + size) % size
	return TermOf(y, wrapped)
}

// String implements fmt.Stringer.
func (t Term) String() string { return t.Name() }

// jdEpsilon is half a second in Julian-Day units. SolarTime ↔ JD
// round-trips lose sub-second precision because SolarTime is integer-
// second; this tolerance lets boundary instants (instants that were
// themselves derived from a Term's SolarTime) still match their term.
const jdEpsilon = 0.5 / 86400

// TermOfSolarTime returns the solar term whose start ≤ s, in the canonical
// order. This is the "current" term for the supplied instant.
//
// Strategy: TermOf(Y, 0) by convention sits in calendar year (Y-1), so the
// 冬至 of calendar year Y is TermOf(Y+1, 0). Probe that first; if it is
// already after s, reverse-walk the 24 terms keyed under calendar year Y.
func TermOfSolarTime(s SolarTime) Term {
	jd := s.JulianDay()
	year := int(s.Year)
	// 冬至 of current calendar year — last possible term ≤ s.
	if t := TermOf(year+1, 0); t.julianDay <= jd+jdEpsilon {
		return t
	}
	for idx := 23; idx >= 0; idx-- {
		t := TermOf(year, idx)
		if t.julianDay <= jd+jdEpsilon {
			return t
		}
	}
	// All terms after s — fall back to previous calendar year's tail.
	return TermOf(year, 0).Next(-1)
}

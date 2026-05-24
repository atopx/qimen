package qimen

import (
	"errors"
	"testing"
	"time"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
)

// solarTime is a test helper: fail on err.
func solarTime(t *testing.T, y, mo, d, h, mi, s int) almanac.SolarTime {
	t.Helper()
	st, err := almanac.SolarTimeOf(y, mo, d, h, mi, s)
	if err != nil {
		t.Fatalf("SolarTimeOf: %v", err)
	}
	return st
}

// TestEntryPointsAgree verifies the 4 chart constructors all produce
// the same chart for an equivalent instant.
func TestEntryPointsAgree(t *testing.T) {
	cst := time.FixedZone("CST", 8*3600)
	tt := time.Date(2026, 1, 14, 18, 45, 0, 0, cst)
	st, _ := almanac.SolarTimeFromTime(tt)

	c1, _ := From(st)
	c2, _ := FromTime(tt)
	c3, _ := FromTimestamp(tt.Unix())

	for _, pair := range [][2]*Chart{{c1, c2}, {c1, c3}} {
		a, b := pair[0], pair[1]
		if a.Year() != b.Year() ||
			a.Month() != b.Month() ||
			a.Day() != b.Day() ||
			a.Hour() != b.Hour() ||
			a.Ju() != b.Ju() {
			t.Errorf("entry points disagree: %+v vs %+v",
				a.SolarTime(), b.SolarTime())
		}
	}
}

// TestSelfPalaceJiaDay covers the 甲 → 值符原宫 fallback path.
func TestSelfPalaceJiaDay(t *testing.T) {
	c := MustFrom(solarTime(t, 2025, 5, 5, 5, 5, 0))
	if c.Day().Stem().Name() != "甲" {
		t.Skip("not a jia day")
	}
	if got, want := c.SelfPalace(), c.ZhiFu().OriginalPalace; got != want {
		t.Errorf("self palace on jia day: got %d, want %d", got, want)
	}
}

// TestUnsupportedMethod covers the error path for non-time methods.
func TestUnsupportedMethod(t *testing.T) {
	st := solarTime(t, 2026, 3, 2, 18, 30, 0)
	_, err := From(st, WithMethod(enum.MethodDay))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrUnsupportedMethod) {
		t.Errorf("err: got %v, want wrapping ErrUnsupportedMethod", err)
	}
}

// TestUnsupportedStyle covers the error path for non-rotate styles.
func TestUnsupportedStyle(t *testing.T) {
	st := solarTime(t, 2026, 3, 2, 18, 30, 0)
	_, err := From(st, WithStyle(enum.StyleFly))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrUnsupportedStyle) {
		t.Errorf("err: got %v, want wrapping ErrUnsupportedStyle", err)
	}
}

// TestPalaceBounds covers nil-returning out-of-range palace lookups.
func TestPalaceBounds(t *testing.T) {
	c := MustFrom(solarTime(t, 2026, 1, 14, 18, 45, 0))
	if c.Palace(0) != nil {
		t.Error("palace 0 should be nil")
	}
	if c.Palace(10) != nil {
		t.Error("palace 10 should be nil")
	}
}

// TestPalacesIter verifies the iter.Seq2 yields all 9 palaces in order.
func TestPalacesIter(t *testing.T) {
	c := MustFrom(solarTime(t, 2026, 1, 14, 18, 45, 0))
	count := 0
	for n, p := range c.Palaces() {
		count++
		if p.Number != n {
			t.Errorf("yield mismatch: n=%d, palace=%d", n, p.Number)
		}
	}
	if count != 9 {
		t.Errorf("expected 9 palaces, got %d", count)
	}
}

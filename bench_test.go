package qimen

import (
	"testing"

	"github.com/atopx/qimen/almanac"
)

// benchSamples 覆盖一年中分布在不同节气、阴阳遁、上中下三元的代表性时刻。
var benchSamples = []struct {
	y, m, d, h, mn, s int
}{
	{2025, 1, 1, 0, 30, 0},   // 冬至附近,阳遁
	{2025, 3, 20, 6, 30, 0},  // 春分附近
	{2025, 5, 5, 5, 5, 0},    // 立夏
	{2025, 6, 21, 12, 30, 0}, // 夏至前后
	{2025, 8, 8, 18, 30, 0},  // 立秋附近
	{2025, 9, 23, 22, 30, 0}, // 秋分附近
	{2025, 12, 22, 0, 30, 0}, // 冬至
}

func mustSolarTimes(b *testing.B) []almanac.SolarTime {
	b.Helper()
	out := make([]almanac.SolarTime, len(benchSamples))
	for i, s := range benchSamples {
		st, err := almanac.SolarTimeOf(s.y, s.m, s.d, s.h, s.mn, s.s)
		if err != nil {
			b.Fatalf("SolarTimeOf[%d]: %v", i, err)
		}
		out[i] = st
	}
	return out
}

// BenchmarkFrom builds a chart from a SolarTime; matches the historical
// FromSolarTime benchmark for direct comparison with the baseline.
func BenchmarkFrom(b *testing.B) {
	sts := mustSolarTimes(b)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = From(sts[i%len(sts)])
	}
}

// BenchmarkFrom_WithReads builds a chart and immediately consumes the
// pattern / shensha / grid aggregations.
func BenchmarkFrom_WithReads(b *testing.B) {
	sts := mustSolarTimes(b)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := From(sts[i%len(sts)])
		for range c.Patterns() {
		}
		for range c.ShenSha() {
		}
		_ = c.Grid()
	}
}

// BenchmarkPalaceAccess measures the cost of repeated palace lookups
// on an already-built chart.
func BenchmarkPalaceAccess(b *testing.B) {
	sts := mustSolarTimes(b)
	c := From(sts[0])
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = c.Palace(uint8(i%9) + 1)
	}
}

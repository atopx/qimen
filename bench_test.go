package qimen

import (
	"testing"

	"github.com/6tail/tyme4go/tyme"
)

// benchSamples 覆盖一年中分布在不同节气、阴阳遁、上中下三元的代表性时刻。
var benchSamples = []struct {
	y, m, d, h, mn, s int
}{
	{2025, 1, 1, 0, 30, 0},   // 冬至附近,阳遁
	{2025, 3, 20, 6, 30, 0},  // 春分附近
	{2025, 5, 5, 5, 5, 0},    // 立夏(沿用现有 example 用例)
	{2025, 6, 21, 12, 30, 0}, // 夏至前后
	{2025, 8, 8, 18, 30, 0},  // 立秋附近
	{2025, 9, 23, 22, 30, 0}, // 秋分附近
	{2025, 12, 22, 0, 30, 0}, // 冬至
}

func mustSolarTimes(b *testing.B) []tyme.SolarTime {
	b.Helper()
	out := make([]tyme.SolarTime, len(benchSamples))
	for i, s := range benchSamples {
		st, err := tyme.SolarTime{}.FromYmdHms(s.y, s.m, s.d, s.h, s.mn, s.s)
		if err != nil {
			b.Fatalf("FromYmdHms[%d]: %v", i, err)
		}
		out[i] = *st
	}
	return out
}

// BenchmarkFromSolarTime 衡量起盘本身的耗时。
func BenchmarkFromSolarTime(b *testing.B) {
	sts := mustSolarTimes(b)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = FromSolarTime(sts[i%len(sts)])
	}
}

// BenchmarkFromSolarTime_WithReads 起盘 + 取主要聚合 (Patterns/ShenSha/GridLayout),
// 更贴近实际调用路径。
func BenchmarkFromSolarTime_WithReads(b *testing.B) {
	sts := mustSolarTimes(b)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		q := FromSolarTime(sts[i%len(sts)])
		_ = q.Patterns()
		_ = q.ShenSha()
		_ = q.GridLayout()
	}
}

// BenchmarkPalaceAccess 衡量已起盘后频繁按宫位号查询的开销。
func BenchmarkPalaceAccess(b *testing.B) {
	sts := mustSolarTimes(b)
	q := FromSolarTime(sts[0])
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = q.Palace(uint8(i%9) + 1)
	}
}

package qimen

import (
	"testing"

	"github.com/6tail/tyme4go/tyme"
)

// TestSolarTermBoundaryYinYang 覆盖冬至/夏至阴阳遁的精确切换边界。
// examples/valid 按小时采样, 不会命中节气的精确秒级时刻, 因此该边界需独立测试。
func TestSolarTermBoundaryYinYang(t *testing.T) {
	// 冬至: 在前为阴 (上年阴遁尾), 边界即冬至时刻为阳
	winter := tyme.SolarTerm{}.FromIndex(2027, 0).GetJulianDay().GetSolarTime()
	if got := computeYinYang(winter); got != tyme.YANG {
		t.Errorf("winter boundary: got %v, want YANG", got)
	}
	prev := winter.Next(-1)
	if got := computeYinYang(prev); got != tyme.YIN {
		t.Errorf("winter-1s: got %v, want YIN", got)
	}

	summer, err := tyme.SolarTerm{}.FromName(2026, "夏至")
	if err != nil {
		t.Fatal(err)
	}
	summerTime := summer.GetJulianDay().GetSolarTime()
	if got := computeYinYang(summerTime); got != tyme.YIN {
		t.Errorf("summer boundary: got %v, want YIN", got)
	}
	if got := computeYinYang(summerTime.Next(-1)); got != tyme.YANG {
		t.Errorf("summer-1s: got %v, want YANG", got)
	}
}

package qimen

import (
	"testing"

	"github.com/6tail/tyme4go/tyme"
)

// solarTime 是一个测试辅助, panic on err。
func solarTime(t *testing.T, y, mo, d, h, mi, s int) tyme.SolarTime {
	t.Helper()
	st, err := tyme.SolarTime{}.FromYmdHms(y, mo, d, h, mi, s)
	if err != nil {
		t.Fatalf("invalid solar time: %v", err)
	}
	return *st
}

// TestSelfPalaceJiaDayUsesZhiFuOrig 覆盖 SelfPalace 在甲日的特殊分支
// (甲遁不可见, 退化为值符原宫)。该路径未被 examples/valid 校验覆盖。
func TestSelfPalaceJiaDayUsesZhiFuOrig(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2025, 5, 5, 5, 5, 0))
	if q.Day().GetHeavenStem().GetName() != "甲" {
		t.Skip("not a jia day")
	}
	if got, want := q.SelfPalace(), q.ZhiFu().OriginalPalace; got != want {
		t.Errorf("self palace on jia day: got %d, want %d", got, want)
	}
}

// TestUnsupportedMethod 覆盖错误路径; examples/valid 仅使用默认参数。
func TestUnsupportedMethod(t *testing.T) {
	st := solarTime(t, 2026, 3, 2, 18, 30, 0)
	_, err := FromSolarTimeWithOptions(st, QimenOptions{Method: QimenMethodDay, ChartType: QimenChartTypeSanYuan})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "unsupported qimen method: 日家" {
		t.Errorf("err: got %q", err.Error())
	}
}

// TestPalaceBounds 覆盖宫位号越界返回 nil 的边界条件。
func TestPalaceBounds(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 1, 14, 18, 45, 0))
	if q.Palace(0) != nil {
		t.Error("palace 0 should be nil")
	}
	if q.Palace(10) != nil {
		t.Error("palace 10 should be nil")
	}
}

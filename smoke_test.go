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

func TestSmoke(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 1, 14, 18, 45, 0))
	if q == nil {
		t.Fatal("FromSolarTime returned nil")
	}
}

func TestYangSnapshot(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 1, 14, 18, 45, 0))

	if got, want := q.Term().GetName(), "小寒"; got != want {
		t.Errorf("term: got %q, want %q", got, want)
	}
	if got, want := q.YinYang(), tyme.YANG; got != want {
		t.Errorf("yinYang: got %v, want %v", got, want)
	}
	if got, want := q.Yuan(), QimenYuanMiddle; got != want {
		t.Errorf("yuan: got %v, want %v", got, want)
	}
	if got, want := q.Ju(), uint8(8); got != want {
		t.Errorf("ju: got %d, want %d", got, want)
	}
	if got, want := q.XunShou().GetName(), "癸"; got != want {
		t.Errorf("xunShou: got %q, want %q", got, want)
	}
	if got, want := q.Year().GetName(), "乙巳"; got != want {
		t.Errorf("year: got %q, want %q", got, want)
	}
	if got, want := q.Month().GetName(), "己丑"; got != want {
		t.Errorf("month: got %q, want %q", got, want)
	}
	if got, want := q.Day().GetName(), "戊子"; got != want {
		t.Errorf("day: got %q, want %q", got, want)
	}
	if got, want := q.Hour().GetName(), "辛酉"; got != want {
		t.Errorf("hour: got %q, want %q", got, want)
	}
	if got, want := q.ZhiFu().Star, QimenStarTianFu; got != want {
		t.Errorf("zhiFu.star: got %v, want %v", got, want)
	}
	if got, want := q.ZhiFu().OriginalPalace, uint8(4); got != want {
		t.Errorf("zhiFu.original: got %d, want %d", got, want)
	}
	if got, want := q.ZhiFu().Palace, uint8(2); got != want {
		t.Errorf("zhiFu.palace: got %d, want %d", got, want)
	}
	if got, want := q.ZhiShi().Door, QimenDoorBlock; got != want {
		t.Errorf("zhiShi.door: got %v, want %v", got, want)
	}
	if got, want := q.ZhiShi().OriginalPalace, uint8(4); got != want {
		t.Errorf("zhiShi.original: got %d, want %d", got, want)
	}
	if got, want := q.ZhiShi().Palace, uint8(2); got != want {
		t.Errorf("zhiShi.palace: got %d, want %d", got, want)
	}

	// SanQiLiuYi: nameByPalace lookup
	nameByPalace := func(items []QimenHeavenStemPlacement, p uint8) string {
		for _, it := range items {
			if it.Palace == p {
				return it.HeavenStem.GetName()
			}
		}
		return ""
	}
	sq := q.SanQiLiuYi()
	for _, c := range []struct {
		palace uint8
		want   string
	}{
		{1, "庚"}, {2, "辛"}, {3, "壬"}, {4, "癸"}, {5, "丁"},
		{6, "丙"}, {7, "乙"}, {8, "戊"}, {9, "己"},
	} {
		if got := nameByPalace(sq, c.palace); got != c.want {
			t.Errorf("sanqi[%d]: got %q, want %q", c.palace, got, c.want)
		}
	}
	tian := q.TianPan()
	for _, c := range []struct {
		palace uint8
		want   string
	}{
		{1, "乙"}, {2, "癸"}, {3, "庚"}, {4, "戊"}, {5, "丁"},
		{6, "辛"}, {7, "己"}, {8, "丙"}, {9, "壬"},
	} {
		if got := nameByPalace(tian, c.palace); got != c.want {
			t.Errorf("tian[%d]: got %q, want %q", c.palace, got, c.want)
		}
	}
	hidden := q.HiddenHeavenStems()
	for _, c := range []struct {
		palace uint8
		want   string
	}{
		{1, "庚"}, {2, "辛"}, {3, "壬"}, {4, "癸"}, {5, "丁"},
		{6, "丙"}, {7, "乙"}, {8, "戊"}, {9, "己"},
	} {
		if got := nameByPalace(hidden, c.palace); got != c.want {
			t.Errorf("hidden[%d]: got %q, want %q", c.palace, got, c.want)
		}
	}

	// Kong wang: 子、丑
	kw := q.KongWang()
	if kw[0].GetName() != "子" || kw[1].GetName() != "丑" {
		t.Errorf("kongWang: got [%q,%q], want [子,丑]", kw[0].GetName(), kw[1].GetName())
	}

	// 空亡格局: 1宫子, 8宫丑
	kwPalaces := map[uint8]bool{}
	for _, p := range q.Patterns() {
		if p.Kind == PatternKongWang {
			kwPalaces[p.Palace] = true
		}
	}
	if !kwPalaces[1] {
		t.Errorf("expected KongWang in palace 1")
	}
	if !kwPalaces[8] {
		t.Errorf("expected KongWang in palace 8")
	}

	// 驿马 - 日支子→寅 落 8 宫
	var yiMa []ShenSha
	for _, s := range q.ShenSha() {
		if s.Kind == ShenShaYiMa {
			yiMa = append(yiMa, s)
		}
	}
	if len(yiMa) != 1 {
		t.Fatalf("expected 1 YiMa, got %d", len(yiMa))
	}
	if yiMa[0].PalaceCell != 8 {
		t.Errorf("YiMa palace: got %d, want 8", yiMa[0].PalaceCell)
	}

	// Palace 5 should have no star/door/god
	center := q.Palace(5)
	if center.Star != nil {
		t.Errorf("center star: got %v, want nil", *center.Star)
	}
	if center.Door != nil {
		t.Errorf("center door: got %v, want nil", *center.Door)
	}
	if center.God != nil {
		t.Errorf("center god: got %v, want nil", *center.God)
	}
	if center.HeavenHeavenStem.GetName() != "丁" {
		t.Errorf("center heaven: got %q, want 丁", center.HeavenHeavenStem.GetName())
	}
}

func TestYinSnapshot(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 10, 31, 12, 2, 0))

	if got, want := q.Term().GetName(), "霜降"; got != want {
		t.Errorf("term: got %q, want %q", got, want)
	}
	if got, want := q.YinYang(), tyme.YIN; got != want {
		t.Errorf("yinYang: got %v, want %v", got, want)
	}
	if got, want := q.Yuan(), QimenYuanLower; got != want {
		t.Errorf("yuan: got %v, want %v", got, want)
	}
	if got, want := q.Ju(), uint8(2); got != want {
		t.Errorf("ju: got %d, want %d", got, want)
	}
	if got, want := q.XunShou().GetName(), "癸"; got != want {
		t.Errorf("xunShou: got %q, want %q", got, want)
	}
}

func TestSelfPalaceJiaDayUsesZhiFuOrig(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2025, 5, 5, 5, 5, 0))
	if q.Day().GetHeavenStem().GetName() != "甲" {
		t.Skip("not a jia day")
	}
	if got, want := q.SelfPalace(), q.ZhiFu().OriginalPalace; got != want {
		t.Errorf("self palace on jia day: got %d, want %d", got, want)
	}
}

func TestHexagramsYinDun3(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2025, 9, 15, 10, 0, 0))
	want := map[uint8]string{
		1: "风水涣", 2: "水地比", 3: "地雷复", 4: "泽风大过",
		6: "雷天大壮", 7: "山泽损", 8: "火山旅", 9: "天火同人",
	}
	for n, w := range want {
		p := q.Palace(n)
		if p.Hexagram == nil {
			t.Errorf("palace %d: missing hexagram (want %q)", n, w)
			continue
		}
		if got := p.Hexagram.Name(); got != w {
			t.Errorf("palace %d: got %q, want %q", n, got, w)
		}
	}
	if q.Palace(5).Hexagram != nil {
		t.Error("center palace should have no hexagram")
	}
}

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

func TestPalaceBounds(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 1, 14, 18, 45, 0))
	if q.Palace(0) != nil {
		t.Error("palace 0 should be nil")
	}
	if q.Palace(10) != nil {
		t.Error("palace 10 should be nil")
	}
	for n := uint8(1); n <= 9; n++ {
		if q.Palace(n) == nil {
			t.Errorf("palace %d should be non-nil", n)
		}
	}
}

func TestConstantsBasic(t *testing.T) {
	if len(TermJu) != 24 {
		t.Errorf("TermJu length: got %d, want 24", len(TermJu))
	}
	for i, row := range TermJu {
		for j, v := range row {
			if v < 1 || v > 9 {
				t.Errorf("TermJu[%d][%d]=%d out of [1,9]", i, j, v)
			}
		}
	}
	if len(LuoShuOrder) != 8 {
		t.Errorf("LuoShuOrder length: got %d, want 8", len(LuoShuOrder))
	}
	for _, v := range LuoShuOrder {
		if v == 5 {
			t.Error("LuoShuOrder must not contain 5")
		}
	}
	if s := QimenStarFromPalace(1); s == nil || s.Name() != "天蓬" {
		t.Errorf("QimenStarFromPalace(1) failed: %v", s)
	}
	if QimenStarFromPalace(5) == nil {
		// center is fine (天禽)
	}
	if got := ElementWood.RelationTo(ElementFire); got != ElementRelationGenerates {
		t.Errorf("wood→fire: got %v, want Generates", got)
	}
	if AuspiceGreatAuspicious.Name() != "大吉" {
		t.Errorf("auspice name")
	}
}

func TestStepPalace(t *testing.T) {
	if stepPalace(9, tyme.YANG) != 1 {
		t.Error("step 9 yang should be 1")
	}
	if stepPalace(4, tyme.YANG) != 5 {
		t.Errorf("step 4 yang should be 5, got %d", stepPalace(4, tyme.YANG))
	}
	if stepPalace(1, tyme.YIN) != 9 {
		t.Errorf("step 1 yin should be 9, got %d", stepPalace(1, tyme.YIN))
	}
	if areOppositePalaces(1, 9) == false || areOppositePalaces(1, 8) == true {
		t.Error("oppositePalaces")
	}
}

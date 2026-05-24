package qimen

import (
	"testing"

	"github.com/6tail/tyme4go/tyme"
)

// TestHexagramStructure 覆盖卦象结构与 64 卦完备性。
// Symbol/Index 字段不在 examples/valid 比对范围内, 需独立测试。
func TestHexagramStructure(t *testing.T) {
	// 中宫无卦
	if TrigramFromPalace(5) != nil {
		t.Error("palace 5 should map to nil trigram")
	}
	// Symbol/Index 由静态表查询
	h := NewHexagram(TrigramDui, TrigramXun)
	if h.Symbol() != "䷛" {
		t.Errorf("symbol: got %q", h.Symbol())
	}
	if h.Index() != 27 {
		t.Errorf("index: got %d", h.Index())
	}
	// 8 卦 × 8 卦 = 64 卦, 索引必须唯一且元数据完整
	trigrams := []Trigram{TrigramQian, TrigramDui, TrigramLi, TrigramZhen, TrigramXun, TrigramKan, TrigramGen, TrigramKun}
	seen := map[uint8]bool{}
	for _, u := range trigrams {
		for _, l := range trigrams {
			hx := NewHexagram(u, l)
			seen[hx.Index()] = true
			if hx.Name() == "" || hx.Symbol() == "" {
				t.Errorf("missing data for %v %v", u, l)
			}
		}
	}
	if len(seen) != 64 {
		t.Errorf("got %d unique hexagrams, want 64", len(seen))
	}
}

// TestPatternNamesAndAuspice 枚举所有格局类型, 确保元数据完整。
// 罕见格局可能未在 examples/valid 的 2025 数据中出现, 需独立兜底。
func TestPatternNamesAndAuspice(t *testing.T) {
	cases := []struct {
		k    PatternKind
		name string
		a    Auspice
	}{
		{PatternFanYin, "反吟", AuspiceInauspicious},
		{PatternFuYin, "伏吟", AuspiceInauspicious},
		{PatternRuMu, "入墓", AuspiceInauspicious},
		{PatternKongWang, "落空亡", AuspiceInauspicious},
		{PatternMenPo, "门迫", AuspiceInauspicious},
		{PatternYiQiDeShi, "乙奇得使", AuspiceAuspicious},
		{PatternBingQiDeShi, "丙奇得使", AuspiceAuspicious},
		{PatternDingQiDeShi, "丁奇得使", AuspiceAuspicious},
		{PatternTianDun, "天遁", AuspiceAuspicious},
		{PatternDiDun, "地遁", AuspiceAuspicious},
		{PatternRenDun, "人遁", AuspiceAuspicious},
		{PatternShenDun, "神遁", AuspiceAuspicious},
		{PatternGuiDun, "鬼遁", AuspiceAuspicious},
		{PatternFengDun, "风遁", AuspiceAuspicious},
		{PatternYunDun, "云遁", AuspiceAuspicious},
		{PatternLongDun, "龙遁", AuspiceAuspicious},
		{PatternHuDun, "虎遁", AuspiceAuspicious},
		{PatternQingLongFanShou, "青龙返首", AuspiceGreatAuspicious},
		{PatternFeiNiaoDieXue, "飞鸟跌穴", AuspiceGreatAuspicious},
		{PatternDaGe, "大格", AuspiceGreatInauspicious},
		{PatternXiaoGe, "小格", AuspiceInauspicious},
		{PatternXingGe, "刑格", AuspiceInauspicious},
		{PatternBoGe, "悖格", AuspiceGreatInauspicious},
		{PatternTianWangSiZhang, "天网四张", AuspiceGreatInauspicious},
	}
	for _, c := range cases {
		p := Pattern{Kind: c.k}
		if p.Name() != c.name {
			t.Errorf("name: got %q, want %q", p.Name(), c.name)
		}
		if p.Summary() == "" {
			t.Errorf("%s missing summary", c.name)
		}
		if p.Auspice() != c.a {
			t.Errorf("%s auspice: got %v, want %v", c.name, p.Auspice(), c.a)
		}
	}
}

// TestShenShaKinds 枚举所有神煞类型, 确保元数据完整。
// 罕见神煞可能未在 examples/valid 的 2025 数据中出现, 需独立兜底。
func TestShenShaKinds(t *testing.T) {
	for _, c := range []struct {
		k    ShenShaKind
		name string
		a    Auspice
	}{
		{ShenShaYiMa, "驿马", AuspiceNeutral},
		{ShenShaTaoHua, "桃花", AuspiceNeutral},
		{ShenShaHuaGai, "华盖", AuspiceNeutral},
		{ShenShaTianYi, "天乙贵人", AuspiceGreatAuspicious},
		{ShenShaTianDe, "天德贵人", AuspiceGreatAuspicious},
		{ShenShaYueDe, "月德贵人", AuspiceGreatAuspicious},
		{ShenShaGuoYin, "国印贵人", AuspiceAuspicious},
		{ShenShaWenChang, "文昌", AuspiceAuspicious},
		{ShenShaLuShen, "禄神", AuspiceAuspicious},
		{ShenShaYangRen, "羊刃", AuspiceInauspicious},
	} {
		if c.k.Name() != c.name {
			t.Errorf("kind name: got %q, want %q", c.k.Name(), c.name)
		}
		if c.k.Auspice() != c.a {
			t.Errorf("%s auspice: got %v", c.name, c.k.Auspice())
		}
		if c.k.Summary() == "" {
			t.Errorf("%s missing summary", c.name)
		}
	}
}

// TestShenShaStringFormat 覆盖 String() 格式化输出, 该字段不在 examples/valid 比对内。
func TestShenShaStringFormat(t *testing.T) {
	branch := tyme.EarthBranch{}.FromIndex(2)
	s := ShenSha{Kind: ShenShaYiMa, Target: ShenShaTarget{Branch: &branch}, PalaceCell: 8}
	if got, want := s.String(), "驿马(寅→8宫)"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	stem := tyme.HeavenStem{}.FromIndex(2)
	s2 := ShenSha{Kind: ShenShaYueDe, Target: ShenShaTarget{Stem: &stem}, PalaceCell: 6}
	if got, want := s2.String(), "月德贵人(丙→6宫)"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

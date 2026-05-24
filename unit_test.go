package qimen

import (
	"testing"

	"github.com/6tail/tyme4go/tyme"
)

func TestAuspiceNames(t *testing.T) {
	cases := []struct {
		a    Auspice
		name string
	}{
		{AuspiceGreatAuspicious, "大吉"},
		{AuspiceAuspicious, "吉"},
		{AuspiceNeutral, "中和"},
		{AuspiceInauspicious, "凶"},
		{AuspiceGreatInauspicious, "大凶"},
	}
	for _, c := range cases {
		if got := c.a.Name(); got != c.name {
			t.Errorf("Name: got %q, want %q", got, c.name)
		}
		if got := c.a.String(); got != c.name {
			t.Errorf("String: got %q, want %q", got, c.name)
		}
	}
	if !AuspiceGreatAuspicious.IsAuspicious() || !AuspiceAuspicious.IsAuspicious() {
		t.Error("auspicious")
	}
	if !AuspiceGreatInauspicious.IsInauspicious() || !AuspiceInauspicious.IsInauspicious() {
		t.Error("inauspicious")
	}
	if !AuspiceGreatAuspicious.IsExtreme() || !AuspiceGreatInauspicious.IsExtreme() {
		t.Error("extreme")
	}
	if AuspiceNeutral.IsAuspicious() || AuspiceNeutral.IsInauspicious() {
		t.Error("neutral")
	}
}

func TestEnumNames(t *testing.T) {
	for _, c := range []struct {
		v    interface{ Name() string }
		want string
	}{
		{QimenMethodTime, "时家"},
		{QimenMethodDay, "日家"},
		{QimenMethodMonth, "月家"},
		{QimenMethodYear, "年家"},
		{QimenChartTypeSanYuan, "三元"},
		{QimenChartTypeSiZhu, "四柱"},
		{QimenYuanUpper, "上元"},
		{QimenYuanMiddle, "中元"},
		{QimenYuanLower, "下元"},
		{QimenStarTianPeng, "天蓬"},
		{QimenStarTianRui, "天芮"},
		{QimenStarTianChong, "天冲"},
		{QimenStarTianFu, "天辅"},
		{QimenStarTianQin, "天禽"},
		{QimenStarTianXin, "天心"},
		{QimenStarTianZhu, "天柱"},
		{QimenStarTianRen, "天任"},
		{QimenStarTianYing, "天英"},
		{QimenStarQinRui, "禽芮"},
		{QimenDoorRest, "休门"},
		{QimenDoorLife, "生门"},
		{QimenDoorHurt, "伤门"},
		{QimenDoorBlock, "杜门"},
		{QimenDoorView, "景门"},
		{QimenDoorDeath, "死门"},
		{QimenDoorFear, "惊门"},
		{QimenDoorOpen, "开门"},
		{QimenGodZhiFu, "值符"},
		{QimenGodTengShe, "腾蛇"},
		{QimenGodTaiYin, "太阴"},
		{QimenGodLiuHe, "六合"},
		{QimenGodBaiHu, "白虎"},
		{QimenGodXuanWu, "玄武"},
		{QimenGodJiuDi, "九地"},
		{QimenGodJiuTian, "九天"},
	} {
		if got := c.v.Name(); got != c.want {
			t.Errorf("Name: got %q, want %q", got, c.want)
		}
	}
}

func TestStarDoorHomePalace(t *testing.T) {
	cases := []struct {
		s    QimenStar
		want uint8
	}{
		{QimenStarTianPeng, 1},
		{QimenStarTianRui, 2},
		{QimenStarQinRui, 2},
		{QimenStarTianChong, 3},
		{QimenStarTianFu, 4},
		{QimenStarTianQin, 5},
		{QimenStarTianXin, 6},
		{QimenStarTianZhu, 7},
		{QimenStarTianRen, 8},
		{QimenStarTianYing, 9},
	}
	for _, c := range cases {
		if got := c.s.HomePalace(); got != c.want {
			t.Errorf("%v.HomePalace: got %d, want %d", c.s, got, c.want)
		}
	}
	for n := uint8(1); n <= 9; n++ {
		s := QimenStarFromPalace(n)
		if s == nil {
			t.Errorf("StarFromPalace(%d) is nil", n)
		}
	}
	if QimenStarFromPalace(0) != nil || QimenStarFromPalace(10) != nil {
		t.Error("StarFromPalace out-of-range should be nil")
	}
	if QimenDoorFromPalace(5) != nil {
		t.Error("DoorFromPalace(5) should be nil")
	}
	for _, p := range []uint8{1, 2, 3, 4, 6, 7, 8, 9} {
		if QimenDoorFromPalace(p) == nil {
			t.Errorf("DoorFromPalace(%d) should be non-nil", p)
		}
	}
}

func TestElements(t *testing.T) {
	for _, c := range []struct {
		idx  int
		want Element
	}{
		{0, ElementWood}, // 甲
		{1, ElementWood}, // 乙
		{2, ElementFire}, // 丙
		{4, ElementEarth},
		{6, ElementMetal},
		{9, ElementWater},
	} {
		if got := ElementFromHeavenStemIndex(c.idx); got != c.want {
			t.Errorf("stem %d: got %v, want %v", c.idx, got, c.want)
		}
	}
	for _, c := range []struct {
		idx  int
		want Element
	}{
		{0, ElementWater}, // 子
		{2, ElementWood},  // 寅
		{5, ElementFire},  // 巳
		{8, ElementMetal}, // 申
		{4, ElementEarth}, // 辰
	} {
		if got := ElementFromEarthBranchIndex(c.idx); got != c.want {
			t.Errorf("branch %d: got %v, want %v", c.idx, got, c.want)
		}
	}
	for _, c := range []struct {
		p    uint8
		want Element
	}{
		{1, ElementWater}, {2, ElementEarth}, {3, ElementWood}, {4, ElementWood},
		{5, ElementEarth}, {6, ElementMetal}, {7, ElementMetal}, {8, ElementEarth}, {9, ElementFire},
	} {
		if got := ElementFromPalace(c.p); got != c.want {
			t.Errorf("palace %d: got %v, want %v", c.p, got, c.want)
		}
	}
}

func TestElementRelations(t *testing.T) {
	cases := []struct {
		a, b Element
		want ElementRelation
	}{
		{ElementWood, ElementWood, ElementRelationSame},
		{ElementWood, ElementFire, ElementRelationGenerates},
		{ElementFire, ElementEarth, ElementRelationGenerates},
		{ElementEarth, ElementMetal, ElementRelationGenerates},
		{ElementMetal, ElementWater, ElementRelationGenerates},
		{ElementWater, ElementWood, ElementRelationGenerates},
		{ElementFire, ElementWood, ElementRelationGenerated},
		{ElementWood, ElementWater, ElementRelationGenerated},
		{ElementWood, ElementEarth, ElementRelationRestrains},
		{ElementEarth, ElementWater, ElementRelationRestrains},
		{ElementWater, ElementFire, ElementRelationRestrains},
		{ElementFire, ElementMetal, ElementRelationRestrains},
		{ElementMetal, ElementWood, ElementRelationRestrains},
		{ElementEarth, ElementWood, ElementRelationRestrained},
		{ElementWood, ElementMetal, ElementRelationRestrained},
	}
	for _, c := range cases {
		if got := c.a.RelationTo(c.b); got != c.want {
			t.Errorf("%v→%v: got %v, want %v", c.a, c.b, got, c.want)
		}
	}
	// Auspice mapping
	if ElementRelationGenerated.AuspiceAsSelf() != AuspiceAuspicious {
		t.Error("Generated should be auspicious")
	}
	if ElementRelationGenerates.AuspiceAsSelf() != AuspiceInauspicious {
		t.Error("Generates should be inauspicious")
	}
	if ElementRelationSame.AuspiceAsSelf() != AuspiceNeutral {
		t.Error("Same should be neutral")
	}
}

func TestStarDoorElements(t *testing.T) {
	if QimenStarTianPeng.Element() != ElementWater {
		t.Error("TianPeng element")
	}
	if QimenStarTianYing.Element() != ElementFire {
		t.Error("TianYing element")
	}
	if QimenDoorRest.Element() != ElementWater {
		t.Error("Rest element")
	}
	if QimenDoorView.Element() != ElementFire {
		t.Error("View element")
	}
	if QimenDoorOpen.Element() != ElementMetal {
		t.Error("Open element")
	}
}

func TestTrigramAndHexagram(t *testing.T) {
	if got := TrigramFromPalace(1); got == nil || *got != TrigramKan {
		t.Error("palace 1 → Kan")
	}
	if got := TrigramFromPalace(5); got != nil {
		t.Error("palace 5 → nil")
	}
	// 上兑下巽 = 泽风大过 (序号 27)
	h := NewHexagram(TrigramDui, TrigramXun)
	if h.Name() != "泽风大过" {
		t.Errorf("name: got %q", h.Name())
	}
	if h.Symbol() != "䷛" {
		t.Errorf("symbol: got %q", h.Symbol())
	}
	if h.Index() != 27 {
		t.Errorf("index: got %d", h.Index())
	}
	if h.Auspice() != AuspiceInauspicious {
		t.Errorf("auspice: got %v", h.Auspice())
	}
	// 乾乾 = 乾为天
	h2 := NewHexagram(TrigramQian, TrigramQian)
	if h2.Name() != "乾为天" || h2.Auspice() != AuspiceGreatAuspicious {
		t.Error("乾为天 mismatch")
	}
	// All 64 unique
	trigrams := []Trigram{TrigramQian, TrigramDui, TrigramLi, TrigramZhen, TrigramXun, TrigramKan, TrigramGen, TrigramKun}
	seen := map[uint8]bool{}
	for _, u := range trigrams {
		for _, l := range trigrams {
			h := NewHexagram(u, l)
			seen[h.Index()] = true
			if h.Name() == "" || h.Symbol() == "" {
				t.Errorf("missing data for %v %v", u, l)
			}
		}
	}
	if len(seen) != 64 {
		t.Errorf("got %d unique hexagrams", len(seen))
	}
}

func TestTerrainNames(t *testing.T) {
	cases := []struct {
		name string
		want Auspice
	}{
		{"长生", AuspiceAuspicious},
		{"沐浴", AuspiceNeutral},
		{"冠带", AuspiceAuspicious},
		{"临官", AuspiceAuspicious},
		{"帝旺", AuspiceGreatAuspicious},
		{"衰", AuspiceNeutral},
		{"病", AuspiceInauspicious},
		{"死", AuspiceGreatInauspicious},
		{"墓", AuspiceInauspicious},
		{"绝", AuspiceGreatInauspicious},
		{"胎", AuspiceNeutral},
		{"养", AuspiceAuspicious},
	}
	for _, c := range cases {
		raw, err := tyme.Terrain{}.FromName(c.name)
		if err != nil {
			t.Fatal(err)
		}
		w := NewTerrain(*raw)
		if w.Name() != c.name {
			t.Errorf("name: got %q, want %q", w.Name(), c.name)
		}
		if w.Auspice() != c.want {
			t.Errorf("%s auspice: got %v, want %v", c.name, w.Auspice(), c.want)
		}
		if w.Summary() == "" {
			t.Errorf("%s missing summary", c.name)
		}
	}
}

func TestPalaceRelationDescriptions(t *testing.T) {
	cases := []struct {
		subjectEl, palaceEl Element
		desc                string
		auspice             Auspice
	}{
		{ElementWood, ElementWood, "伤门与宫比和", AuspiceNeutral},
		{ElementWood, ElementFire, "伤门生宫", AuspiceInauspicious},
		{ElementWood, ElementWater, "宫生伤门", AuspiceAuspicious},
		{ElementWood, ElementEarth, "伤门克宫", AuspiceNeutral},
		{ElementWood, ElementMetal, "宫克伤门", AuspiceInauspicious},
	}
	for _, c := range cases {
		r := palaceRelationForSubject("伤门", c.subjectEl, c.palaceEl)
		if r.Description != c.desc {
			t.Errorf("desc: got %q, want %q", r.Description, c.desc)
		}
		if r.Auspice() != c.auspice {
			t.Errorf("auspice: got %v, want %v", r.Auspice(), c.auspice)
		}
	}
}

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

func TestShenShaStringFormat(t *testing.T) {
	branch := tyme.EarthBranch{}.FromIndex(2)
	s := ShenSha{Kind: ShenShaYiMa, Target: ShenShaTarget{Branch: &branch}, PalaceCell: 8}
	want := "驿马(寅→8宫)"
	if got := s.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	stem := tyme.HeavenStem{}.FromIndex(2)
	s2 := ShenSha{Kind: ShenShaYueDe, Target: ShenShaTarget{Stem: &stem}, PalaceCell: 6}
	if got := s2.String(); got != "月德贵人(丙→6宫)" {
		t.Errorf("got %q", got)
	}
}

func TestPlateSetGet(t *testing.T) {
	p := NewPlate[int]()
	if p.Get(5) != nil {
		t.Error("empty plate Get(5) should be nil")
	}
	if p.Get(0) != nil || p.Get(10) != nil {
		t.Error("out of range Get should be nil")
	}
	p.Set(5, 42)
	v := p.Get(5)
	if v == nil || *v != 42 {
		t.Errorf("Get(5) = %v, want 42", v)
	}
}

func TestPlatePanicsOnInvalidSet(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	p := NewPlate[int]()
	p.Set(10, 1)
}

func TestMovePalaceBySteps(t *testing.T) {
	// 阳遁顺行从 1, 4 步: 1→2→3→4→5, 末尾在 5 改寄 2
	if got := movePalaceBySteps(1, 4, tyme.YANG); got != 2 {
		t.Errorf("got %d, want 2 (5→2 kickout)", got)
	}
	// 阳遁顺行从 1, 5 步: 1→2→3→4→5→6 = 6 (末尾不在 5)
	if got := movePalaceBySteps(1, 5, tyme.YANG); got != 6 {
		t.Errorf("got %d, want 6", got)
	}
	// 阴遁逆行从 9, 1 步: 9→8
	if got := movePalaceBySteps(9, 1, tyme.YIN); got != 8 {
		t.Errorf("got %d, want 8", got)
	}
}

func TestIsStemInTomb(t *testing.T) {
	jia := tyme.HeavenStem{}.FromIndex(0) // 甲
	if !isStemInTomb(jia, 2) {
		t.Error("甲 should be tombed in 2")
	}
	if isStemInTomb(jia, 5) {
		t.Error("甲 should not be tombed in 5")
	}
}

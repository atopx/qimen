package qimen

import (
	"testing"

	"github.com/6tail/tyme4go/tyme"
)

func TestComputeYuan(t *testing.T) {
	mustCycle := func(name string) tyme.SixtyCycle {
		c, err := tyme.SixtyCycle{}.FromName(name)
		if err != nil {
			t.Fatal(err)
		}
		return *c
	}
	cases := []struct {
		name string
		want QimenYuan
	}{
		{"甲子", QimenYuanUpper},
		{"己巳", QimenYuanMiddle},
		{"甲戌", QimenYuanLower},
		{"戊子", QimenYuanMiddle},
	}
	for _, c := range cases {
		if got := computeYuan(mustCycle(c.name)); got != c.want {
			t.Errorf("computeYuan(%s): got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestComputeJuTable(t *testing.T) {
	xiaoHan, err := tyme.SolarTerm{}.FromName(2026, "小寒")
	if err != nil {
		t.Fatal(err)
	}
	check := func(label string, term tyme.SolarTerm, yuan QimenYuan, want uint8) {
		got, err := computeJu(term, yuan)
		if err != nil {
			t.Errorf("computeJu(%s,%v) err: %v", label, yuan, err)
			return
		}
		if got != want {
			t.Errorf("computeJu(%s,%v): got %d, want %d", label, yuan, got, want)
		}
	}
	check("小寒", *xiaoHan, QimenYuanUpper, 2)
	check("小寒", *xiaoHan, QimenYuanMiddle, 8)
	check("小寒", *xiaoHan, QimenYuanLower, 5)

	shuangJiang, err := tyme.SolarTerm{}.FromName(2026, "霜降")
	if err != nil {
		t.Fatal(err)
	}
	check("霜降", *shuangJiang, QimenYuanUpper, 5)
	check("霜降", *shuangJiang, QimenYuanMiddle, 8)
	check("霜降", *shuangJiang, QimenYuanLower, 2)
}

func TestComputeXunShou(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"甲子", "戊"},
		{"乙亥", "己"},
		{"戊子", "庚"},
		{"癸卯", "辛"},
		{"壬子", "壬"},
		{"辛酉", "癸"},
	}
	for _, c := range cases {
		cyc, err := tyme.SixtyCycle{}.FromName(c.name)
		if err != nil {
			t.Fatal(err)
		}
		got := computeXunShou(*cyc).GetName()
		if got != c.want {
			t.Errorf("xunShou(%s): got %q, want %q", c.name, got, c.want)
		}
	}
}

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

func TestGridLayout(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 3, 2, 18, 30, 0))
	grid := q.GridLayout()
	if grid[0][0].Number != 4 {
		t.Errorf("grid[0][0]: got %d, want 4", grid[0][0].Number)
	}
	if grid[1][1].Number != 5 {
		t.Errorf("grid[1][1]: got %d, want 5", grid[1][1].Number)
	}
	if grid[2][2].Number != 6 {
		t.Errorf("grid[2][2]: got %d, want 6", grid[2][2].Number)
	}
}

func TestHexagramsSnapshotPure8(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 5, 11, 15, 30, 0))
	want := map[uint8]string{
		1: "坎为水", 2: "坤为地", 3: "震为雷", 4: "巽为风",
		6: "乾为天", 7: "兑为泽", 8: "艮为山", 9: "离为火",
	}
	for n, w := range want {
		p := q.Palace(n)
		if p.Hexagram == nil {
			t.Errorf("palace %d: missing hexagram", n)
			continue
		}
		if got := p.Hexagram.Name(); got != w {
			t.Errorf("palace %d: got %q, want %q", n, got, w)
		}
	}
}

func TestHexagramsSnapshotYangDun3(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 3, 4, 9, 13, 0))
	want := map[uint8]string{
		1: "地水师", 2: "雷地豫", 3: "天雷无妄", 4: "水风井",
		6: "火天大有", 7: "风泽中孚", 8: "泽山咸", 9: "山火贲",
	}
	for n, w := range want {
		p := q.Palace(n)
		if p.Hexagram == nil {
			t.Errorf("palace %d: missing hexagram", n)
			continue
		}
		if got := p.Hexagram.Name(); got != w {
			t.Errorf("palace %d: got %q, want %q", n, got, w)
		}
	}
}

func TestShenShaIncludesMultipleKinds(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 1, 14, 18, 45, 0))
	kinds := map[ShenShaKind]bool{}
	for _, s := range q.ShenSha() {
		kinds[s.Kind] = true
	}
	for _, k := range []ShenShaKind{
		ShenShaYiMa, ShenShaTaoHua, ShenShaHuaGai,
		ShenShaTianYi, ShenShaTianDe, ShenShaYueDe,
		ShenShaGuoYin, ShenShaWenChang, ShenShaLuShen,
	} {
		if !kinds[k] {
			t.Errorf("missing kind %s", k.Name())
		}
	}
}

func TestTenStarsCenterIsNone(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 1, 14, 18, 45, 0))
	for n := uint8(1); n <= 9; n++ {
		p := q.Palace(n)
		if n == 5 {
			if p.TenStar != nil {
				t.Errorf("center TenStar should be nil")
			}
		} else if p.TenStar == nil {
			t.Errorf("palace %d TenStar should be non-nil", n)
		}
	}
}

func TestPalaceFullSnapshot(t *testing.T) {
	q := FromSolarTime(solarTime(t, 2026, 1, 14, 18, 45, 0))
	type want struct {
		num            uint8
		star           QimenStar
		door           QimenDoor
		god            QimenGod
		heaven, ground string
	}
	cases := []want{
		{1, QimenStarTianZhu, QimenDoorFear, QimenGodLiuHe, "乙", "庚"},
		{2, QimenStarTianFu, QimenDoorBlock, QimenGodZhiFu, "癸", "辛"},
		{3, QimenStarTianPeng, QimenDoorRest, QimenGodXuanWu, "庚", "壬"},
		{4, QimenStarTianRen, QimenDoorLife, QimenGodJiuDi, "戊", "癸"},
		{6, QimenStarQinRui, QimenDoorDeath, QimenGodTaiYin, "辛", "丙"},
		{7, QimenStarTianYing, QimenDoorView, QimenGodTengShe, "己", "乙"},
		{8, QimenStarTianXin, QimenDoorOpen, QimenGodBaiHu, "丙", "戊"},
		{9, QimenStarTianChong, QimenDoorHurt, QimenGodJiuTian, "壬", "己"},
	}
	for _, c := range cases {
		p := q.Palace(c.num)
		if p.Star == nil || *p.Star != c.star {
			t.Errorf("palace %d star: got %v, want %v", c.num, p.Star, c.star)
		}
		if p.Door == nil || *p.Door != c.door {
			t.Errorf("palace %d door: got %v, want %v", c.num, p.Door, c.door)
		}
		if p.God == nil || *p.God != c.god {
			t.Errorf("palace %d god: got %v, want %v", c.num, p.God, c.god)
		}
		if p.HeavenHeavenStem.GetName() != c.heaven {
			t.Errorf("palace %d heaven: got %q, want %q", c.num, p.HeavenHeavenStem.GetName(), c.heaven)
		}
		if p.EarthHeavenStem.GetName() != c.ground {
			t.Errorf("palace %d earth: got %q, want %q", c.num, p.EarthHeavenStem.GetName(), c.ground)
		}
	}
}

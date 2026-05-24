package qimen

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/6tail/tyme4go/tyme"
)

// dumpEnv 指向 Rust 端导出的 JSONL 文件路径。未设置或文件不存在时跳过测试。
const dumpEnv = "QIMEN_VERIFY_DUMP"

// dumpRecord 与 examples/dump.rs 输出 schema 对应。
type dumpRecord struct {
	T           dumpTime      `json:"t"`
	YearPillar  string        `json:"year_pillar"`
	MonthPillar string        `json:"month_pillar"`
	DayPillar   string        `json:"day_pillar"`
	HourPillar  string        `json:"hour_pillar"`
	Term        string        `json:"term"`
	YinYang     string        `json:"yin_yang"`
	Ju          uint8         `json:"ju"`
	Yuan        string        `json:"yuan"`
	XunShou     string        `json:"xun_shou"`
	ZhiFu       dumpDutyStar  `json:"zhi_fu"`
	ZhiShi      dumpDutyDoor  `json:"zhi_shi"`
	KongWang    [2]string     `json:"kong_wang"`
	Patterns    []dumpPattern `json:"patterns"`
	ShenSha     []dumpShenSha `json:"shen_sha"`
	Palaces     []dumpPalace  `json:"palaces"`
}

type dumpTime struct {
	Y  int   `json:"y"`
	M  uint8 `json:"m"`
	D  uint8 `json:"d"`
	H  uint8 `json:"H"`
	Mn uint8 `json:"M"`
	S  uint8 `json:"S"`
}

type dumpDutyStar struct {
	Star           string `json:"star"`
	OriginalPalace uint8  `json:"original_palace"`
	Palace         uint8  `json:"palace"`
}

type dumpDutyDoor struct {
	Door           string `json:"door"`
	OriginalPalace uint8  `json:"original_palace"`
	Palace         uint8  `json:"palace"`
}

type dumpPattern struct {
	Name      string `json:"name"`
	Palace    uint8  `json:"palace"`
	Auspice   string `json:"auspice"`
	DetailKey string `json:"detail_key"`
}

type dumpShenSha struct {
	Kind    string `json:"kind"`
	Target  string `json:"target"`
	Palace  uint8  `json:"palace"`
	Auspice string `json:"auspice"`
}

type dumpHexagram struct {
	Name    string `json:"name"`
	Auspice string `json:"auspice"`
}

type dumpTerrain struct {
	Name    string `json:"name"`
	Auspice string `json:"auspice"`
}

type dumpPalace struct {
	Number             uint8         `json:"number"`
	Name               string        `json:"name"`
	Direction          string        `json:"direction"`
	Element            string        `json:"element"`
	EarthBranches      []string      `json:"earth_branches"`
	EarthHeavenStem    string        `json:"earth_heaven_stem"`
	SanQiLiuYi         string        `json:"san_qi_liu_yi"`
	HeavenHeavenStem   string        `json:"heaven_heaven_stem"`
	HiddenHeavenStem   string        `json:"hidden_heaven_stem"`
	Star               *string       `json:"star"`
	Door               *string       `json:"door"`
	God                *string       `json:"god"`
	TenStar            *string       `json:"ten_star"`
	Terrain            *dumpTerrain  `json:"terrain"`
	Hexagram           *dumpHexagram `json:"hexagram"`
	DoorPalaceRelation *string       `json:"door_palace_relation"`
	StarPalaceRelation *string       `json:"star_palace_relation"`
	Patterns           []dumpPattern `json:"patterns"`
	ShenSha            []dumpShenSha `json:"shen_sha"`
}

// TestVerifyDumpAgainstRust 逐行读取 Rust dump 并用 Go 库重建后比对。
// 通过环境变量 QIMEN_VERIFY_DUMP=<path> 触发,默认 SKIP。
func TestVerifyDumpAgainstRust(t *testing.T) {
	path := os.Getenv(dumpEnv)
	if path == "" {
		t.Skipf("set %s=<path/to/jsonl> to enable cross-port verification", dumpEnv)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Skipf("dump file unavailable (%s): %v", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1<<20)
	line := 0
	for scanner.Scan() {
		line++
		raw := scanner.Bytes()
		if len(raw) == 0 {
			continue
		}
		var rec dumpRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			t.Fatalf("line %d: parse: %v", line, err)
		}
		label := fmt.Sprintf("%04d-%02d-%02d_%02d:%02d", rec.T.Y, rec.T.M, rec.T.D, rec.T.H, rec.T.Mn)
		t.Run(label, func(t *testing.T) {
			st, err := tyme.SolarTime{}.FromYmdHms(rec.T.Y, int(rec.T.M), int(rec.T.D), int(rec.T.H), int(rec.T.Mn), int(rec.T.S))
			if err != nil {
				t.Fatalf("FromYmdHms: %v", err)
			}
			q := FromSolarTime(*st)
			verifyOne(t, &rec, q)
		})
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if line == 0 {
		t.Fatalf("dump file is empty: %s", path)
	}
}

func verifyOne(t *testing.T, want *dumpRecord, q *Qimen) {
	t.Helper()
	cmpStr := func(field, got, exp string) {
		if got != exp {
			t.Errorf("%s: got %q, want %q", field, got, exp)
		}
	}
	cmpStr("year_pillar", q.Year().GetName(), want.YearPillar)
	cmpStr("month_pillar", q.Month().GetName(), want.MonthPillar)
	cmpStr("day_pillar", q.Day().GetName(), want.DayPillar)
	cmpStr("hour_pillar", q.Hour().GetName(), want.HourPillar)
	cmpStr("term", q.Term().GetName(), want.Term)
	cmpStr("yin_yang", q.YinYang().GetName(), want.YinYang)
	if q.Ju() != want.Ju {
		t.Errorf("ju: got %d, want %d", q.Ju(), want.Ju)
	}
	cmpStr("yuan", q.Yuan().Name(), want.Yuan)
	cmpStr("xun_shou", q.XunShou().GetName(), want.XunShou)

	zf := q.ZhiFu()
	cmpStr("zhi_fu.star", zf.Star.Name(), want.ZhiFu.Star)
	if zf.OriginalPalace != want.ZhiFu.OriginalPalace {
		t.Errorf("zhi_fu.original_palace: got %d, want %d", zf.OriginalPalace, want.ZhiFu.OriginalPalace)
	}
	if zf.Palace != want.ZhiFu.Palace {
		t.Errorf("zhi_fu.palace: got %d, want %d", zf.Palace, want.ZhiFu.Palace)
	}
	zs := q.ZhiShi()
	cmpStr("zhi_shi.door", zs.Door.Name(), want.ZhiShi.Door)
	if zs.OriginalPalace != want.ZhiShi.OriginalPalace {
		t.Errorf("zhi_shi.original_palace: got %d, want %d", zs.OriginalPalace, want.ZhiShi.OriginalPalace)
	}
	if zs.Palace != want.ZhiShi.Palace {
		t.Errorf("zhi_shi.palace: got %d, want %d", zs.Palace, want.ZhiShi.Palace)
	}
	kw := q.KongWang()
	cmpStr("kong_wang[0]", kw[0].GetName(), want.KongWang[0])
	cmpStr("kong_wang[1]", kw[1].GetName(), want.KongWang[1])

	verifyPatterns(t, "patterns", q.Patterns(), want.Patterns)
	verifyShenSha(t, "shen_sha", q.ShenSha(), want.ShenSha)

	wantByNum := make(map[uint8]*dumpPalace, 9)
	for i := range want.Palaces {
		wantByNum[want.Palaces[i].Number] = &want.Palaces[i]
	}
	for n := uint8(1); n <= 9; n++ {
		wp := wantByNum[n]
		if wp == nil {
			t.Errorf("missing palace %d in dump", n)
			continue
		}
		gp := q.Palace(n)
		if gp == nil {
			t.Errorf("palace %d: go returned nil", n)
			continue
		}
		verifyPalace(t, n, gp, wp)
	}
}

func verifyPalace(t *testing.T, n uint8, got *QimenPalace, want *dumpPalace) {
	t.Helper()
	prefix := fmt.Sprintf("palace%d", n)
	mk := func(field string) string { return prefix + "." + field }
	if got.Number != want.Number {
		t.Errorf("%s: got %d, want %d", mk("number"), got.Number, want.Number)
	}
	if got.PalaceName != want.Name {
		t.Errorf("%s: got %q, want %q", mk("name"), got.PalaceName, want.Name)
	}
	if got.Direction.GetName() != want.Direction {
		t.Errorf("%s: got %q, want %q", mk("direction"), got.Direction.GetName(), want.Direction)
	}
	if got.Element().Name() != want.Element {
		t.Errorf("%s: got %q, want %q", mk("element"), got.Element().Name(), want.Element)
	}
	branches := make([]string, 0, len(got.EarthBranches))
	for _, b := range got.EarthBranches {
		branches = append(branches, b.GetName())
	}
	if !reflect.DeepEqual(branches, want.EarthBranches) {
		t.Errorf("%s: got %v, want %v", mk("earth_branches"), branches, want.EarthBranches)
	}
	if got.EarthHeavenStem.GetName() != want.EarthHeavenStem {
		t.Errorf("%s: got %q, want %q", mk("earth_heaven_stem"), got.EarthHeavenStem.GetName(), want.EarthHeavenStem)
	}
	if got.SanQiLiuYi.GetName() != want.SanQiLiuYi {
		t.Errorf("%s: got %q, want %q", mk("san_qi_liu_yi"), got.SanQiLiuYi.GetName(), want.SanQiLiuYi)
	}
	if got.HeavenHeavenStem.GetName() != want.HeavenHeavenStem {
		t.Errorf("%s: got %q, want %q", mk("heaven_heaven_stem"), got.HeavenHeavenStem.GetName(), want.HeavenHeavenStem)
	}
	if got.HiddenHeavenStem.GetName() != want.HiddenHeavenStem {
		t.Errorf("%s: got %q, want %q", mk("hidden_heaven_stem"), got.HiddenHeavenStem.GetName(), want.HiddenHeavenStem)
	}
	verifyOptString(t, mk("star"), starName(got.Star), want.Star)
	verifyOptString(t, mk("door"), doorName(got.Door), want.Door)
	verifyOptString(t, mk("god"), godName(got.God), want.God)
	verifyOptString(t, mk("ten_star"), tenStarName(got.TenStar), want.TenStar)

	verifyOptTerrain(t, mk("terrain"), got.TerrainValue, want.Terrain)
	verifyOptHexagram(t, mk("hexagram"), got.Hexagram, want.Hexagram)

	verifyOptString(t, mk("door_palace_relation"), relationDesc(got.DoorPalaceRelation()), want.DoorPalaceRelation)
	verifyOptString(t, mk("star_palace_relation"), relationDesc(got.StarPalaceRelation()), want.StarPalaceRelation)

	verifyPatterns(t, mk("patterns"), got.Patterns, want.Patterns)
	verifyShenSha(t, mk("shen_sha"), got.ShenSha, want.ShenSha)
}

func verifyOptString(t *testing.T, field string, got *string, want *string) {
	t.Helper()
	switch {
	case got == nil && want == nil:
		return
	case got == nil:
		t.Errorf("%s: got nil, want %q", field, *want)
	case want == nil:
		t.Errorf("%s: got %q, want nil", field, *got)
	case *got != *want:
		t.Errorf("%s: got %q, want %q", field, *got, *want)
	}
}

func verifyOptTerrain(t *testing.T, field string, got *Terrain, want *dumpTerrain) {
	t.Helper()
	switch {
	case got == nil && want == nil:
		return
	case got == nil:
		t.Errorf("%s: got nil, want {name:%s auspice:%s}", field, want.Name, want.Auspice)
	case want == nil:
		t.Errorf("%s: got {name:%s auspice:%s}, want nil", field, got.Name(), got.Auspice().Name())
	default:
		if got.Name() != want.Name {
			t.Errorf("%s.name: got %q, want %q", field, got.Name(), want.Name)
		}
		if got.Auspice().Name() != want.Auspice {
			t.Errorf("%s.auspice: got %q, want %q", field, got.Auspice().Name(), want.Auspice)
		}
	}
}

func verifyOptHexagram(t *testing.T, field string, got *Hexagram, want *dumpHexagram) {
	t.Helper()
	switch {
	case got == nil && want == nil:
		return
	case got == nil:
		t.Errorf("%s: got nil, want {name:%s auspice:%s}", field, want.Name, want.Auspice)
	case want == nil:
		t.Errorf("%s: got {name:%s auspice:%s}, want nil", field, got.Name(), got.Auspice().Name())
	default:
		if got.Name() != want.Name {
			t.Errorf("%s.name: got %q, want %q", field, got.Name(), want.Name)
		}
		if got.Auspice().Name() != want.Auspice {
			t.Errorf("%s.auspice: got %q, want %q", field, got.Auspice().Name(), want.Auspice)
		}
	}
}

func verifyPatterns(t *testing.T, field string, got []Pattern, want []dumpPattern) {
	t.Helper()
	gotDumps := make([]dumpPattern, len(got))
	for i, p := range got {
		gotDumps[i] = dumpPattern{
			Name:      p.Name(),
			Palace:    p.Palace,
			Auspice:   p.Auspice().Name(),
			DetailKey: patternDetailKey(p),
		}
	}
	sort.Slice(gotDumps, func(i, j int) bool { return gotDumps[i].DetailKey < gotDumps[j].DetailKey })
	wantSorted := make([]dumpPattern, len(want))
	copy(wantSorted, want)
	sort.Slice(wantSorted, func(i, j int) bool { return wantSorted[i].DetailKey < wantSorted[j].DetailKey })
	if len(gotDumps) != len(wantSorted) {
		t.Errorf("%s length mismatch: got %d, want %d\n got:  %s\n want: %s",
			field, len(gotDumps), len(wantSorted), formatPatterns(gotDumps), formatPatterns(wantSorted))
		return
	}
	for i := range gotDumps {
		if gotDumps[i] != wantSorted[i] {
			t.Errorf("%s mismatch:\n got:  %s\n want: %s", field, formatPatterns(gotDumps), formatPatterns(wantSorted))
			return
		}
	}
}

func verifyShenSha(t *testing.T, field string, got []ShenSha, want []dumpShenSha) {
	t.Helper()
	gotDumps := make([]dumpShenSha, len(got))
	for i, s := range got {
		gotDumps[i] = dumpShenSha{
			Kind:    s.Kind.Name(),
			Target:  s.Target.String(),
			Palace:  s.Palace(),
			Auspice: s.Auspice().Name(),
		}
	}
	sortKey := func(s dumpShenSha) string {
		return fmt.Sprintf("%02d|%s|%s", s.Palace, s.Kind, s.Target)
	}
	sort.Slice(gotDumps, func(i, j int) bool { return sortKey(gotDumps[i]) < sortKey(gotDumps[j]) })
	wantSorted := make([]dumpShenSha, len(want))
	copy(wantSorted, want)
	sort.Slice(wantSorted, func(i, j int) bool { return sortKey(wantSorted[i]) < sortKey(wantSorted[j]) })
	if len(gotDumps) != len(wantSorted) {
		t.Errorf("%s length mismatch: got %d, want %d\n got:  %s\n want: %s",
			field, len(gotDumps), len(wantSorted), formatShenSha(gotDumps), formatShenSha(wantSorted))
		return
	}
	for i := range gotDumps {
		if gotDumps[i] != wantSorted[i] {
			t.Errorf("%s mismatch:\n got:  %s\n want: %s", field, formatShenSha(gotDumps), formatShenSha(wantSorted))
			return
		}
	}
}

func patternDetailKey(p Pattern) string {
	detail := ""
	switch p.Kind {
	case PatternFanYin:
		detail = strconv.Itoa(int(p.OriginalPalace))
	case PatternMenPo:
		if p.Door != nil {
			detail = p.Door.Name()
		}
	case PatternRuMu:
		if p.Stem != nil {
			detail = p.Stem.GetName()
		}
	case PatternKongWang:
		if p.Branch != nil {
			detail = p.Branch.GetName()
		}
	}
	return fmt.Sprintf("%d|%s|%s", p.Palace, p.Name(), detail)
}

func starName(s *QimenStar) *string {
	if s == nil {
		return nil
	}
	n := s.Name()
	return &n
}

func doorName(d *QimenDoor) *string {
	if d == nil {
		return nil
	}
	n := d.Name()
	return &n
}

func godName(g *QimenGod) *string {
	if g == nil {
		return nil
	}
	n := g.Name()
	return &n
}

func tenStarName(t *tyme.TenStar) *string {
	if t == nil {
		return nil
	}
	n := t.GetName()
	return &n
}

func relationDesc(r *PalaceRelation) *string {
	if r == nil {
		return nil
	}
	d := r.Description
	return &d
}

func formatPatterns(ps []dumpPattern) string {
	parts := make([]string, 0, len(ps))
	for _, p := range ps {
		parts = append(parts, fmt.Sprintf("{%d %s/%s/%s}", p.Palace, p.Name, p.Auspice, p.DetailKey))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func formatShenSha(ss []dumpShenSha) string {
	parts := make([]string, 0, len(ss))
	for _, s := range ss {
		parts = append(parts, fmt.Sprintf("{%d %s→%s/%s}", s.Palace, s.Kind, s.Target, s.Auspice))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

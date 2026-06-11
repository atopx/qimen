// Package main 校验 JSONL 中的奇门盘面数据。
//
// 运行: go run ./examples/valid -file <path/to/data.jsonl>
//
// 程序逐行读取 jsonl, 用本库依据 `t` 重新起局, 并将计算结果与 JSONL 中的
// 字段逐项比对, 报告不一致项及总体统计。
//
// 未传 -file 或路径不存在/非 .jsonl 直接 panic。
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/atopx/qimen"
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/palace"
	"github.com/atopx/qimen/pattern"
	"github.com/atopx/qimen/shensha"
)

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

func main() {
	filePath := flag.String("file", "", "path to a .jsonl validation file (required)")
	flag.Parse()

	if *filePath == "" {
		panic("-file is required (path to .jsonl validation data)")
	}
	if !strings.HasSuffix(*filePath, ".jsonl") {
		panic(fmt.Sprintf("-file must end with .jsonl, got %q", *filePath))
	}
	testData, err := os.ReadFile(*filePath)
	if err != nil {
		panic(fmt.Sprintf("read %q: %v", *filePath, err))
	}

	scanner := bufio.NewScanner(bytes.NewReader(testData))
	scanner.Buffer(make([]byte, 64*1024), 1<<20)

	totalRecords := 0
	failedRecords := 0
	totalMismatches := 0
	start := time.Now()

	for scanner.Scan() {
		raw := scanner.Bytes()
		if len(bytes.TrimSpace(raw)) == 0 {
			continue
		}
		totalRecords++

		var rec dumpRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			failedRecords++
			totalMismatches++
			_, _ = fmt.Fprintf(os.Stderr, "record %d: parse error: %v\n", totalRecords, err)
			if err != nil {
				return
			}
			continue
		}
		label := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
			rec.T.Y, rec.T.M, rec.T.D, rec.T.H, rec.T.Mn, rec.T.S)

		st, err := almanac.SolarTimeOf(
			rec.T.Y, int(rec.T.M), int(rec.T.D),
			int(rec.T.H), int(rec.T.Mn), int(rec.T.S),
		)
		if err != nil {
			failedRecords++
			totalMismatches++
			_, _ = fmt.Fprintf(os.Stderr, "[%s] SolarTimeOf error: %v\n", label, err)
			continue
		}

		c := qimen.From(st)
		mismatches := verifyOne(&rec, c)
		if len(mismatches) > 0 {
			failedRecords++
			totalMismatches += len(mismatches)
			_, _ = fmt.Fprintf(os.Stderr, "[%s] FAIL (%d mismatch):\n", label, len(mismatches))
			for _, m := range mismatches {
				_, _ = fmt.Fprintf(os.Stderr, "  - %s\n", m)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		os.Exit(2)
	}

	if totalRecords == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "no records found in 2025.jsonl")
		os.Exit(2)
	}

	elapsed := time.Since(start)
	passed := totalRecords - failedRecords
	avg := elapsed / time.Duration(totalRecords)
	fmt.Printf("校验完成: 总计 %d 条, 通过 %d 条, 失败 %d 条, 不一致字段 %d 项\n",
		totalRecords, passed, failedRecords, totalMismatches)
	fmt.Printf("耗时统计: 总耗时 %s, 平均耗时 %s/条\n", elapsed, avg)
	if failedRecords > 0 {
		os.Exit(1)
	}
}

func verifyOne(want *dumpRecord, c *qimen.Chart) []string {
	var errs []string
	cmpStr := func(field, got, exp string) {
		if got != exp {
			errs = append(errs, fmt.Sprintf("%s: got %q, want %q", field, got, exp))
		}
	}
	cmpStr("year_pillar", c.Year().Name(), want.YearPillar)
	cmpStr("month_pillar", c.Month().Name(), want.MonthPillar)
	cmpStr("day_pillar", c.Day().Name(), want.DayPillar)
	cmpStr("hour_pillar", c.Hour().Name(), want.HourPillar)
	cmpStr("term", c.Term().Name(), want.Term)
	cmpStr("yin_yang", c.YinYang().Name(), want.YinYang)
	if c.Ju() != want.Ju {
		errs = append(errs, fmt.Sprintf("ju: got %d, want %d", c.Ju(), want.Ju))
	}
	cmpStr("yuan", c.Yuan().Name(), want.Yuan)
	cmpStr("xun_shou", c.XunShou().Name(), want.XunShou)

	zf := c.ZhiFu()
	cmpStr("zhi_fu.star", zf.Star.Name(), want.ZhiFu.Star)
	if zf.OriginalPalace != want.ZhiFu.OriginalPalace {
		errs = append(errs, fmt.Sprintf("zhi_fu.original_palace: got %d, want %d",
			zf.OriginalPalace, want.ZhiFu.OriginalPalace))
	}
	if zf.Palace != want.ZhiFu.Palace {
		errs = append(errs, fmt.Sprintf("zhi_fu.palace: got %d, want %d",
			zf.Palace, want.ZhiFu.Palace))
	}
	zs := c.ZhiShi()
	cmpStr("zhi_shi.door", zs.Door.Name(), want.ZhiShi.Door)
	if zs.OriginalPalace != want.ZhiShi.OriginalPalace {
		errs = append(errs, fmt.Sprintf("zhi_shi.original_palace: got %d, want %d",
			zs.OriginalPalace, want.ZhiShi.OriginalPalace))
	}
	if zs.Palace != want.ZhiShi.Palace {
		errs = append(errs, fmt.Sprintf("zhi_shi.palace: got %d, want %d",
			zs.Palace, want.ZhiShi.Palace))
	}
	kw := c.KongWang()
	cmpStr("kong_wang[0]", kw[0].Name(), want.KongWang[0])
	cmpStr("kong_wang[1]", kw[1].Name(), want.KongWang[1])

	errs = append(errs, diffPatterns("patterns", collectPatterns(c), want.Patterns)...)
	errs = append(errs, diffShenSha("shen_sha", collectShenSha(c), want.ShenSha)...)

	wantByNum := make(map[uint8]*dumpPalace, 9)
	for i := range want.Palaces {
		wantByNum[want.Palaces[i].Number] = &want.Palaces[i]
	}
	for n := uint8(1); n <= 9; n++ {
		wp := wantByNum[n]
		if wp == nil {
			errs = append(errs, fmt.Sprintf("missing palace %d in dump", n))
			continue
		}
		gp := c.Palace(n)
		if gp == nil {
			errs = append(errs, fmt.Sprintf("palace %d: lib returned nil", n))
			continue
		}
		errs = append(errs, diffPalace(n, gp, wp)...)
	}
	return errs
}

func collectPatterns(c *qimen.Chart) []pattern.Pattern {
	var out []pattern.Pattern
	for p := range c.Patterns() {
		out = append(out, p)
	}
	return out
}

func collectShenSha(c *qimen.Chart) []shensha.ShenSha {
	var out []shensha.ShenSha
	for s := range c.ShenSha() {
		out = append(out, s)
	}
	return out
}

func diffPalace(n uint8, got *palace.Palace, want *dumpPalace) []string {
	var errs []string
	prefix := fmt.Sprintf("palace%d", n)
	mk := func(field string) string { return prefix + "." + field }

	if got.Number != want.Number {
		errs = append(errs, fmt.Sprintf("%s: got %d, want %d", mk("number"), got.Number, want.Number))
	}
	if got.Name != want.Name {
		errs = append(errs, fmt.Sprintf("%s: got %q, want %q", mk("name"), got.Name, want.Name))
	}
	if got.Direction.Name() != want.Direction {
		errs = append(errs, fmt.Sprintf("%s: got %q, want %q",
			mk("direction"), got.Direction.Name(), want.Direction))
	}
	if got.Element().Name() != want.Element {
		errs = append(errs, fmt.Sprintf("%s: got %q, want %q",
			mk("element"), got.Element().Name(), want.Element))
	}
	branches := make([]string, 0, len(got.Branches))
	for _, b := range got.Branches {
		branches = append(branches, b.Name())
	}
	if !stringsEqual(branches, want.EarthBranches) {
		errs = append(errs, fmt.Sprintf("%s: got %v, want %v",
			mk("earth_branches"), branches, want.EarthBranches))
	}
	if got.EarthStem.Name() != want.EarthHeavenStem {
		errs = append(errs, fmt.Sprintf("%s: got %q, want %q",
			mk("earth_heaven_stem"), got.EarthStem.Name(), want.EarthHeavenStem))
	}
	if got.EarthStem.Name() != want.SanQiLiuYi {
		errs = append(errs, fmt.Sprintf("%s: got %q, want %q",
			mk("san_qi_liu_yi"), got.EarthStem.Name(), want.SanQiLiuYi))
	}
	if got.HeavenStem.Name() != want.HeavenHeavenStem {
		errs = append(errs, fmt.Sprintf("%s: got %q, want %q",
			mk("heaven_heaven_stem"), got.HeavenStem.Name(), want.HeavenHeavenStem))
	}
	if got.HiddenStem.Name() != want.HiddenHeavenStem {
		errs = append(errs, fmt.Sprintf("%s: got %q, want %q",
			mk("hidden_heaven_stem"), got.HiddenStem.Name(), want.HiddenHeavenStem))
	}

	errs = append(errs, diffOptString(mk("star"), starName(got), want.Star)...)
	errs = append(errs, diffOptString(mk("door"), doorName(got), want.Door)...)
	errs = append(errs, diffOptString(mk("god"), godName(got), want.God)...)
	errs = append(errs, diffOptString(mk("ten_star"), tenStarName(got), want.TenStar)...)
	errs = append(errs, diffOptTerrain(mk("terrain"), got, want.Terrain)...)
	errs = append(errs, diffOptHexagram(mk("hexagram"), got, want.Hexagram)...)
	errs = append(errs, diffOptString(mk("door_palace_relation"),
		doorPalaceRelDesc(got), want.DoorPalaceRelation)...)
	errs = append(errs, diffOptString(mk("star_palace_relation"),
		starPalaceRelDesc(got), want.StarPalaceRelation)...)
	errs = append(errs, diffPatterns(mk("patterns"), got.Patterns, want.Patterns)...)
	errs = append(errs, diffShenSha(mk("shen_sha"), got.ShenSha, want.ShenSha)...)
	return errs
}

func diffOptString(field string, got, want *string) []string {
	switch {
	case got == nil && want == nil:
		return nil
	case got == nil:
		return []string{fmt.Sprintf("%s: got nil, want %q", field, *want)}
	case want == nil:
		return []string{fmt.Sprintf("%s: got %q, want nil", field, *got)}
	case *got != *want:
		return []string{fmt.Sprintf("%s: got %q, want %q", field, *got, *want)}
	}
	return nil
}

func diffOptTerrain(field string, got *palace.Palace, want *dumpTerrain) []string {
	hasGot := !got.IsCenter()
	switch {
	case !hasGot && want == nil:
		return nil
	case !hasGot:
		return []string{fmt.Sprintf("%s: got nil, want {name:%s auspice:%s}",
			field, want.Name, want.Auspice)}
	case want == nil:
		return []string{fmt.Sprintf("%s: got {name:%s auspice:%s}, want nil",
			field, got.Terrain.Name(), got.Terrain.Auspice().Name())}
	}
	var errs []string
	if got.Terrain.Name() != want.Name {
		errs = append(errs, fmt.Sprintf("%s.name: got %q, want %q",
			field, got.Terrain.Name(), want.Name))
	}
	if got.Terrain.Auspice().Name() != want.Auspice {
		errs = append(errs, fmt.Sprintf("%s.auspice: got %q, want %q",
			field, got.Terrain.Auspice().Name(), want.Auspice))
	}
	return errs
}

func diffOptHexagram(field string, got *palace.Palace, want *dumpHexagram) []string {
	hasGot := !got.IsCenter()
	switch {
	case !hasGot && want == nil:
		return nil
	case !hasGot:
		return []string{fmt.Sprintf("%s: got nil, want {name:%s auspice:%s}",
			field, want.Name, want.Auspice)}
	case want == nil:
		return []string{fmt.Sprintf("%s: got {name:%s auspice:%s}, want nil",
			field, got.Hexagram.Name(), got.Hexagram.Auspice().Name())}
	}
	var errs []string
	if got.Hexagram.Name() != want.Name {
		errs = append(errs, fmt.Sprintf("%s.name: got %q, want %q",
			field, got.Hexagram.Name(), want.Name))
	}
	if got.Hexagram.Auspice().Name() != want.Auspice {
		errs = append(errs, fmt.Sprintf("%s.auspice: got %q, want %q",
			field, got.Hexagram.Auspice().Name(), want.Auspice))
	}
	return errs
}

func diffPatterns(field string, got []pattern.Pattern, want []dumpPattern) []string {
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
		return []string{fmt.Sprintf("%s length mismatch: got %d, want %d\n      got:  %s\n      want: %s",
			field, len(gotDumps), len(wantSorted), formatPatterns(gotDumps), formatPatterns(wantSorted))}
	}
	for i := range gotDumps {
		if gotDumps[i] != wantSorted[i] {
			return []string{fmt.Sprintf("%s mismatch:\n      got:  %s\n      want: %s",
				field, formatPatterns(gotDumps), formatPatterns(wantSorted))}
		}
	}
	return nil
}

func diffShenSha(field string, got []shensha.ShenSha, want []dumpShenSha) []string {
	gotDumps := make([]dumpShenSha, len(got))
	for i, s := range got {
		gotDumps[i] = dumpShenSha{
			Kind:    s.Kind.Name(),
			Target:  s.Target.String(),
			Palace:  s.Palace,
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
		return []string{fmt.Sprintf("%s length mismatch: got %d, want %d\n      got:  %s\n      want: %s",
			field, len(gotDumps), len(wantSorted), formatShenSha(gotDumps), formatShenSha(wantSorted))}
	}
	for i := range gotDumps {
		if gotDumps[i] != wantSorted[i] {
			return []string{fmt.Sprintf("%s mismatch:\n      got:  %s\n      want: %s",
				field, formatShenSha(gotDumps), formatShenSha(wantSorted))}
		}
	}
	return nil
}

func patternDetailKey(p pattern.Pattern) string {
	detail := ""
	switch p.Kind {
	case pattern.FanYin:
		detail = strconv.Itoa(int(p.OriginalPalace))
	case pattern.MenPo:
		detail = p.Door.Name()
	case pattern.RuMu:
		detail = p.Stem.Name()
	case pattern.KongWang:
		detail = p.Branch.Name()
	default:
	}
	return fmt.Sprintf("%d|%s|%s", p.Palace, p.Name(), detail)
}

func starName(p *palace.Palace) *string {
	if p.IsCenter() {
		return nil
	}
	n := p.Star.Name()
	return &n
}

func doorName(p *palace.Palace) *string {
	if p.IsCenter() {
		return nil
	}
	n := p.Door.Name()
	return &n
}

func godName(p *palace.Palace) *string {
	if p.IsCenter() {
		return nil
	}
	n := p.God.Name()
	return &n
}

func tenStarName(p *palace.Palace) *string {
	if p.IsCenter() {
		return nil
	}
	n := p.TenStar.Name()
	return &n
}

func doorPalaceRelDesc(p *palace.Palace) *string {
	if p.IsCenter() {
		return nil
	}
	d := p.DoorPalaceRelation().Description
	return &d
}

func starPalaceRelDesc(p *palace.Palace) *string {
	if p.IsCenter() {
		return nil
	}
	d := p.StarPalaceRelation().Description
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

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

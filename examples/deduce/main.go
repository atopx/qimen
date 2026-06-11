// Package main demonstrates qimen chart construction and the full
// domain surface: 四柱/农历, 四家法门, 转盘/飞盘, 置闰/拆补, 三盘干,
// 星门神, 衍生属性 (十神/长生/演卦/生克), 格局与神煞.
//
// Usage:
//
//	go run ./examples/deduce                          # current time
//	go run ./examples/deduce -time 202601141845       # YYYYMMDDHHMM
//	go run ./examples/deduce -method year -style fly  # 年家飞盘
//	go run ./examples/deduce -rule chaibu             # 时家拆补
//
// Non-conforming flag values panic.
package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/atopx/qimen"
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/auspice"
	"github.com/atopx/qimen/element"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/palace"
	"github.com/atopx/qimen/pattern"
)

var hanNum = [10]string{"", "一", "二", "三", "四", "五", "六", "七", "八", "九"}

func main() {
	timeArg := flag.String("time", "", "instant as YYYYMMDDHHMM (12 digits); default = now")
	methodArg := flag.String("method", "time", "起局法门: time|day|month|year")
	styleArg := flag.String("style", "rotate", "盘式: rotate|fly")
	ruleArg := flag.String("rule", "zhirun", "时家/日家定局规则: zhirun|chaibu")
	flag.Parse()

	chart := qimen.From(parseInstant(*timeArg),
		qimen.WithMethod(parseFlag(*methodArg, methodNames)),
		qimen.WithStyle(parseFlag(*styleArg, styleNames)),
		qimen.WithJuRule(parseFlag(*ruleArg, ruleNames)),
	)

	printHeader(chart)
	printGrid(chart)
	printDetails(chart)
}

// ===================== flag parsing =====================

var methodNames = map[string]enum.Method{
	"time": enum.MethodTime, "day": enum.MethodDay,
	"month": enum.MethodMonth, "year": enum.MethodYear,
}

var styleNames = map[string]enum.Style{
	"rotate": enum.StyleRotate, "fly": enum.StyleFly,
}

var ruleNames = map[string]enum.JuRule{
	"zhirun": enum.JuRuleZhiRun, "chaibu": enum.JuRuleChaiBu,
}

func parseFlag[T any](v string, names map[string]T) T {
	if out, ok := names[v]; ok {
		return out
	}
	keys := make([]string, 0, len(names))
	for k := range names {
		keys = append(keys, k)
	}
	panic(fmt.Sprintf("unknown flag value %q, want one of %s", v, strings.Join(keys, "|")))
}

// parseInstant parses the -time flag value into a SolarTime.
// Empty string returns the current instant; any non-conforming value panics.
func parseInstant(s string) almanac.SolarTime {
	if s == "" {
		return almanac.Now()
	}
	if len(s) != 12 {
		panic(fmt.Sprintf("-time must be 12 digits YYYYMMDDHHMM, got %q (len=%d)", s, len(s)))
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			panic(fmt.Sprintf("-time must be all digits, got %q", s))
		}
	}
	atoi := func(part string) int { n, _ := strconv.Atoi(part); return n }
	st, err := almanac.SolarTimeOf(
		atoi(s[0:4]), atoi(s[4:6]), atoi(s[6:8]), atoi(s[8:10]), atoi(s[10:12]), 0)
	if err != nil {
		panic(fmt.Sprintf("-time %q is invalid: %v", s, err))
	}
	return st
}

// ===================== header =====================

func printHeader(c *qimen.Chart) {
	section("起局")
	row("盘法", fmt.Sprintf("%s%s%s", c.Method().Name(), c.Style().Name(), c.JuRule().Name()))
	row("公历", c.SolarTime().String())
	ld := c.LunarDay()
	row("农历", fmt.Sprintf("%s年%s%s", ld.Year().Cycle().Name(), ld.Month.Name(), ld.Name()))
	row("四柱", fmt.Sprintf("%s %s %s %s", c.Year().Name(), c.Month().Name(), c.Day().Name(), c.Hour().Name()))

	row("节气", fmt.Sprintf("%s%s元", termLine(c), c.Yuan().Name()))
	row("遁局", fmt.Sprintf("%s遁%s局", c.YinYang().Name(), hanNum[c.Ju()]))
	row("主柱", fmt.Sprintf("%s (%s旬, 遁干%s, 空亡%s%s)",
		c.Lead().Name(), c.Lead().Ten().Name(), c.XunShou().Name(),
		c.KongWang()[0].Name(), c.KongWang()[1].Name()))

	zf, zs := c.ZhiFu(), c.ZhiShi()
	row("值符", fmt.Sprintf("%s  %s → %s", zf.Star.Name(),
		palaceLabel(zf.OriginalPalace), palaceLabel(zf.Palace)))
	row("值使", fmt.Sprintf("%s门  %s → %s", zs.Door.Name(),
		palaceLabel(zs.OriginalPalace), palaceLabel(zs.Palace)))
	row("用神", fmt.Sprintf("日干%s落%s (自身)  值符落%s (对方)",
		c.Day().Stem().Name(), palaceLabel(c.SelfPalace()), palaceLabel(c.OpponentPalace())))
}

// termLine renders the astronomical term, annotating the working term
// when 置闰 has shifted it (超神 leads, 接气/闰 trails).
func termLine(c *qimen.Chart) string {
	term, ju := c.Term(), c.JuTerm()
	if term == ju {
		return term.Name()
	}
	seq := func(t almanac.Term) int { return t.Year()*24 + t.Index() }
	state := "接气"
	if seq(ju) > seq(term) {
		state = "超神"
	}
	return fmt.Sprintf("%s · 用局%s (%s)", term.Name(), ju.Name(), state)
}

func section(title string) {
	fmt.Printf("\n── %s %s\n", title, strings.Repeat("─", 40-displayWidth(title)))
}

func row(key, value string) { fmt.Printf("%s  %s\n", key, value) }

func palaceLabel(n uint8) string {
	return [10]string{"", "坎一", "坤二", "震三", "巽四", "中五", "乾六", "兑七", "艮八", "离九"}[n]
}

// ===================== 3×3 grid =====================

const cellWidth = 14

// printGrid renders the canonical 3×3 board (巽离坤 / 震中兑 / 艮坎乾).
// Each cell shows the layout facts only: 神 / 天盘干+星 / 地盘干+门 /
// 暗干 / 宫位 — derived attributes live in the detail section.
func printGrid(c *qimen.Chart) {
	section("盘面")
	border := func(l, m, r string) {
		bar := strings.Repeat("─", cellWidth)
		fmt.Println(l + bar + m + bar + m + bar + r)
	}
	border("┌", "┬", "┐")
	for rowIdx, gridRow := range c.Grid() {
		var lines [5][3]string
		for col, p := range gridRow {
			for i, line := range cellLines(p) {
				lines[i][col] = line
			}
		}
		for _, line := range lines {
			fmt.Printf("│%s│%s│%s│\n",
				padTo(" "+line[0], cellWidth),
				padTo(" "+line[1], cellWidth),
				padTo(" "+line[2], cellWidth))
		}
		if rowIdx < 2 {
			border("├", "┼", "┤")
		}
	}
	border("└", "┴", "┘")
}

// cellLines builds the five display lines of one palace cell. Empty
// slots (rotate-style center, the fly-style door/god gaps) stay blank.
func cellLines(p *palace.Palace) [5]string {
	var god, star, door string
	if p.GodSet {
		god = p.God.Name()
	} else if p.IsCenter() {
		god = "(寄坤二)"
	}
	if p.StarSet {
		star = " " + p.Star.Name()
	}
	if p.DoorSet {
		door = " " + p.Door.Name() + "门"
	}
	label := palaceLabel(p.Number) + " " + p.Direction.Name() + p.Element().Name()
	if p.IsCenter() {
		label = palaceLabel(p.Number) + " " + p.Element().Name()
	}
	return [5]string{
		god,
		p.HeavenStem.Name() + star,
		p.EarthStem.Name() + door,
		"暗 " + p.HiddenStem.Name(),
		label,
	}
}

// ===================== per-palace details =====================

// printDetails renders three aligned sections: a per-palace attribute
// table (十神 / 旺衰 / 演卦 / 星情 / 门情), then the patterns and the
// shensha grouped by palace. 中和 ratings are left untagged so that 吉
// / 凶 marks stand out.
func printDetails(c *qimen.Chart) {
	section("宫位")
	rows := [][]string{{"宫位", "十神", "旺衰", "演卦", "星情", "门情"}}
	for n, p := range c.Palaces() {
		if p.IsCenter() && !p.StarSet {
			continue // rotate-style center carries no derived attributes
		}
		rows = append(rows, []string{
			palaceLabel(n),
			tenStarText(p), terrainText(p), hexagramText(p),
			starRelationText(p), doorRelationText(p),
		})
	}
	printTable(rows)

	section("格局")
	empty := true
	for n, p := range c.Palaces() {
		if len(p.Patterns) == 0 {
			continue
		}
		empty = false
		parts := make([]string, 0, len(p.Patterns))
		for _, pat := range p.Patterns {
			parts = append(parts, patternText(pat))
		}
		row(palaceLabel(n), strings.Join(parts, "  "))
	}
	if empty {
		fmt.Println("无")
	}

	section("神煞")
	for n, p := range c.Palaces() {
		if len(p.ShenSha) == 0 {
			continue
		}
		parts := make([]string, 0, len(p.ShenSha))
		for _, ss := range p.ShenSha {
			parts = append(parts, tag(fmt.Sprintf("%s(%s)", ss.Name(), ss.Target.String()), ss.Auspice()))
		}
		row(palaceLabel(n), strings.Join(parts, "  "))
	}
}

// tag appends the auspice rating in fixed-width brackets, leaving the
// unremarkable 中和 untagged so 吉 / 凶 marks stand out.
func tag(name string, a auspice.Auspice) string {
	return name + "[" + a.Name() + "]"
}

func tenStarText(p *palace.Palace) string {
	if p.IsCenter() {
		return "-"
	}
	return p.TenStar.Name()
}

func terrainText(p *palace.Palace) string {
	if p.IsCenter() {
		return "-"
	}
	return tag(p.Terrain.Name(), p.Terrain.Auspice())
}

func hexagramText(p *palace.Palace) string {
	if !p.HexagramSet {
		return "-"
	}
	return tag(p.Hexagram.Name(), p.Hexagram.Auspice())
}

// starRelationText spells out the 星-宫 五行 relation with the star as
// the named subject, e.g. "天任克宫" / "宫生天柱·吉".
func starRelationText(p *palace.Palace) string {
	if !p.StarSet {
		return "-"
	}
	rel := p.StarPalaceRelation()
	return tag(relationText(p.Star.Name(), rel.Element), rel.Auspice)
}

// doorRelationText does the same for the 门-宫 relation, with the
// single-rune door name expanded to "X门" for readability.
func doorRelationText(p *palace.Palace) string {
	if !p.DoorSet {
		return "-"
	}
	rel := p.DoorPalaceRelation()
	return tag(relationText(p.Door.Name()+"门", rel.Element), rel.Auspice)
}

// relationText phrases a subject-vs-palace relation as a plain clause.
func relationText(subject string, rel element.Relation) string {
	switch rel {
	case element.Generates:
		return subject + "生宫"
	case element.Generated:
		return "宫生" + subject
	case element.Restrains:
		return subject + "克宫"
	case element.Restrained:
		return "宫克" + subject
	default: // element.Same
		return subject + "比和"
	}
}

// patternText renders one pattern with its kind-specific detail.
// Per pattern.Pattern docs, Kind discriminates which side field is set.
func patternText(p pattern.Pattern) string {
	detail := ""
	switch p.Kind {
	case pattern.FanYin:
		detail = palaceLabel(p.OriginalPalace) + "→" + palaceLabel(p.Palace)
	case pattern.MenPo:
		detail = p.Door.Name() + "门"
	case pattern.RuMu:
		detail = p.Stem.Name()
	case pattern.KongWang:
		detail = p.Branch.Name()
	}
	if detail != "" {
		detail = "(" + detail + ")"
	}
	return tag(p.Name()+detail, p.Auspice())
}

// printTable renders rows with CJK-aware column alignment.
func printTable(rows [][]string) {
	cols := len(rows[0])
	widths := make([]int, cols)
	for _, r := range rows {
		for i, cell := range r {
			widths[i] = max(widths[i], displayWidth(cell))
		}
	}
	for _, r := range rows {
		var b strings.Builder
		for i, cell := range r {
			if i > 0 {
				b.WriteString("  ")
			}
			b.WriteString(padTo(cell, widths[i]))
		}
		fmt.Println(strings.TrimRight(b.String(), " "))
	}
}

// ===================== CJK-aware padding =====================

// displayWidth counts terminal columns (CJK runes occupy two).
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		if r > 0x7F {
			w += 2
		} else {
			w++
		}
	}
	return w
}

// padTo right-pads s with spaces to the given display width.
func padTo(s string, width int) string {
	if pad := width - displayWidth(s); pad > 0 {
		return s + strings.Repeat(" ", pad)
	}
	return s
}

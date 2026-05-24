// Package main demonstrates qimen chart construction.
//
// Usage:
//
//	go run ./examples/deduce                 # current time
//	go run ./examples/deduce -time 202505050505  # YYYYMMDDHHMM (秒位补 0)
//
// Non-conforming -time values panic.
package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/atopx/qimen"
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/element"
	"github.com/atopx/qimen/palace"
	"github.com/atopx/qimen/pattern"
	"github.com/atopx/qimen/shensha"
)

func main() {
	timeArg := flag.String("time", "", "input time as YYYYMMDDHHMM (12 digits); default = current time")
	flag.Parse()

	c, err := qimen.From(parseInstant(*timeArg))
	if err != nil {
		panic(err)
	}
	printChart(c)
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
	y, _ := strconv.Atoi(s[0:4])
	mo, _ := strconv.Atoi(s[4:6])
	d, _ := strconv.Atoi(s[6:8])
	h, _ := strconv.Atoi(s[8:10])
	mi, _ := strconv.Atoi(s[10:12])
	st, err := almanac.SolarTimeOf(y, mo, d, h, mi, 0)
	if err != nil {
		panic(fmt.Sprintf("-time %q is invalid: %v", s, err))
	}
	return st
}

func printChart(c *qimen.Chart) {
	fmt.Printf("时间: %s\n", c.SolarTime().String())
	fmt.Printf("四柱: %s %s %s %s\n",
		c.Year().Name(), c.Month().Name(), c.Day().Name(), c.Hour().Name())
	fmt.Printf("节气: %s\n", c.Term().Name())
	fmt.Printf("遁局: %s遁%d局 %s\n", c.YinYang().Name(), c.Ju(), c.Yuan().Name())
	fmt.Printf("旬首: %s\n", c.XunShou().Name())

	zhiFu := c.ZhiFu()
	fmt.Printf("值符: %s 原宫%d 落宫%d\n", zhiFu.Star.Name(), zhiFu.OriginalPalace, zhiFu.Palace)
	zhiShi := c.ZhiShi()
	fmt.Printf("值使: %s 原宫%d 落宫%d\n", zhiShi.Door.Name(), zhiShi.OriginalPalace, zhiShi.Palace)

	kw := c.KongWang()
	fmt.Printf("空亡: %s%s\n", kw[0].Name(), kw[1].Name())

	patterns := collectPatterns(c)
	if len(patterns) == 0 {
		fmt.Println("格局: 无")
	} else {
		parts := make([]string, 0, len(patterns))
		for _, p := range patterns {
			parts = append(parts, patternLine(p))
		}
		fmt.Printf("格局: %s\n", strings.Join(parts, "、"))
	}

	ss := collectShenSha(c)
	if len(ss) == 0 {
		fmt.Println("神煞: 无")
	} else {
		parts := make([]string, 0, len(ss))
		for _, s := range ss {
			parts = append(parts, shenShaLine(s))
		}
		fmt.Printf("神煞: %s\n", strings.Join(parts, "、"))
	}

	fmt.Println()
	for _, row := range c.Grid() {
		for _, p := range row {
			if p.Number == 5 {
				continue
			}
			printPalace(p)
		}
		fmt.Println()
	}
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

// Per pattern.Pattern docs: Kind discriminates which side field is valid.
func patternLine(p pattern.Pattern) string {
	var detail string
	switch p.Kind {
	case pattern.FanYin:
		detail = fmt.Sprintf("%d→%d", p.OriginalPalace, p.Palace)
	case pattern.MenPo:
		detail = fmt.Sprintf("%d宫%s", p.Palace, p.Door.Name())
	case pattern.RuMu:
		detail = fmt.Sprintf("%d宫%s", p.Palace, p.Stem.Name())
	case pattern.KongWang:
		detail = fmt.Sprintf("%d宫%s", p.Palace, p.Branch.Name())
	default:
		detail = fmt.Sprintf("%d宫", p.Palace)
	}
	return fmt.Sprintf("%s[%s](%s)", p.Name(), p.Auspice().Name(), detail)
}

func shenShaLine(s shensha.ShenSha) string {
	return fmt.Sprintf("%s[%s]", s.String(), s.Auspice().Name())
}

func printPalace(p *palace.Palace) {
	starName, starEl, doorName, doorEl := "空", "-", "空", "-"
	godName, tenStarName, terrainText, hexaText := "空", "空", "空", "空"
	if !p.IsCenter() {
		starName = p.Star.Name()
		starEl = element.OfStar(int(p.Star)).Name()
		doorName = p.Door.Name()
		doorEl = element.OfDoor(int(p.Door)).Name()
		godName = p.God.Name()
		tenStarName = p.TenStar.Name()
		terrainText = fmt.Sprintf("%s[%s]", p.Terrain.Name(), p.Terrain.Auspice().Name())
		hexaText = fmt.Sprintf("%s [%s]", p.Hexagram.String(), p.Hexagram.Auspice().Name())
	}
	var branchText strings.Builder
	for _, b := range p.Branches {
		branchText.WriteString(b.Name())
	}

	fmt.Printf(
		"  %d宫%s(%s/%s) 地支:%s 地盘:%s 天盘:%s 暗干:%s 九星:%s(%s) 八门:%s(%s) 八神:%s 十神:%s 长生:%s\n",
		p.Number, p.Name, p.Direction.Name(), p.Element().Name(),
		branchText.String(),
		p.SanQiLiuYi.Name(), p.HeavenStem.Name(), p.HiddenStem.Name(),
		starName, starEl,
		doorName, doorEl,
		godName, tenStarName, terrainText,
	)
	fmt.Printf("       卦:%s\n", hexaText)
}

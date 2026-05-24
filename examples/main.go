// Package main 演示奇门遁甲起局并打印盘面信息。
//
// 运行: go run ./examples
package main

import (
	"fmt"
	"strings"

	"github.com/6tail/tyme4go/tyme"
	"github.com/atopx/qimen"
)

func optStr[T fmt.Stringer](v *T, placeholder string) string {
	if v == nil {
		return placeholder
	}
	return (*v).String()
}

func kongWangText(bs [2]tyme.EarthBranch) string {
	return bs[0].GetName() + bs[1].GetName()
}

func patternLine(p qimen.Pattern) string {
	var detail string
	switch p.Kind {
	case qimen.PatternFanYin:
		detail = fmt.Sprintf("%d→%d", p.OriginalPalace, p.Palace)
	case qimen.PatternMenPo:
		door := ""
		if p.Door != nil {
			door = p.Door.Name()
		}
		detail = fmt.Sprintf("%d宫%s", p.Palace, door)
	case qimen.PatternRuMu:
		stem := ""
		if p.Stem != nil {
			stem = p.Stem.GetName()
		}
		detail = fmt.Sprintf("%d宫%s", p.Palace, stem)
	case qimen.PatternKongWang:
		branch := ""
		if p.Branch != nil {
			branch = p.Branch.GetName()
		}
		detail = fmt.Sprintf("%d宫%s", p.Palace, branch)
	default:
		detail = fmt.Sprintf("%d宫", p.Palace)
	}
	return fmt.Sprintf("%s[%s](%s)", p.Name(), p.Auspice().Name(), detail)
}

func shenShaLine(s qimen.ShenSha) string {
	return fmt.Sprintf("%s[%s]", s.String(), s.Auspice().Name())
}

func printPalace(palace *qimen.QimenPalace) {
	starEl := "-"
	if palace.Star != nil {
		starEl = palace.Star.Element().Name()
	}
	doorEl := "-"
	if palace.Door != nil {
		doorEl = palace.Door.Element().Name()
	}
	terrainText := "空"
	if palace.TerrainValue != nil {
		terrainText = fmt.Sprintf("%s[%s]", palace.TerrainValue.Name(), palace.TerrainValue.Auspice().Name())
	}
	doorRel := "-"
	if r := palace.DoorPalaceRelation(); r != nil {
		doorRel = r.String()
	}
	starRel := "-"
	if r := palace.StarPalaceRelation(); r != nil {
		starRel = r.String()
	}
	hexaText := "空"
	if palace.Hexagram != nil {
		hexaText = fmt.Sprintf("%s [%s]", palace.Hexagram.String(), palace.Hexagram.Auspice().Name())
	}
	branchText := ""
	for _, b := range palace.EarthBranches {
		branchText += b.GetName()
	}
	starName := optStr(palace.Star, "空")
	doorName := optStr(palace.Door, "空")
	godName := optStr(palace.God, "空")
	tenStarName := "空"
	if palace.TenStar != nil {
		tenStarName = palace.TenStar.GetName()
	}
	fmt.Printf(
		"  %d宫%s(%s/%s) 地支:%s 地盘:%s 天盘:%s 暗干:%s 九星:%s(%s) 八门:%s(%s) 八神:%s 十神:%s 长生:%s\n",
		palace.Number, palace.PalaceName,
		palace.Direction.GetName(), palace.Element().Name(),
		branchText,
		palace.SanQiLiuYi.GetName(),
		palace.HeavenHeavenStem.GetName(),
		palace.HiddenHeavenStem.GetName(),
		starName, starEl,
		doorName, doorEl,
		godName, tenStarName, terrainText,
	)
	fmt.Printf("       门宫:%s 星宫:%s 卦:%s\n", doorRel, starRel, hexaText)
}

func main() {
	st, err := tyme.SolarTime{}.FromYmdHms(2025, 5, 5, 5, 5, 0)
	if err != nil {
		panic(err)
	}
	q := qimen.FromSolarTime(*st)

	fmt.Printf("时间: %s\n", st.String())
	fmt.Printf("四柱: %s %s %s %s\n", q.Year().GetName(), q.Month().GetName(), q.Day().GetName(), q.Hour().GetName())
	fmt.Printf("节气: %s\n", q.Term().GetName())
	fmt.Printf("遁局: %s遁%d局 %s\n", q.YinYang().GetName(), q.Ju(), q.Yuan().Name())
	fmt.Printf("旬首: %s\n", q.XunShou().GetName())

	zhiFu := q.ZhiFu()
	fmt.Printf("值符: %s 原宫%d 落宫%d\n", zhiFu.Star.Name(), zhiFu.OriginalPalace, zhiFu.Palace)
	zhiShi := q.ZhiShi()
	fmt.Printf("值使: %s 原宫%d 落宫%d\n", zhiShi.Door.Name(), zhiShi.OriginalPalace, zhiShi.Palace)

	kw := q.KongWang()
	fmt.Printf("空亡: %s\n", kongWangText(kw))

	patterns := q.Patterns()
	if len(patterns) == 0 {
		fmt.Println("格局: 无")
	} else {
		parts := make([]string, 0, len(patterns))
		for _, p := range patterns {
			parts = append(parts, patternLine(p))
		}
		fmt.Printf("格局: %s\n", strings.Join(parts, "、"))
	}

	shenSha := q.ShenSha()
	if len(shenSha) == 0 {
		fmt.Println("神煞: 无")
	} else {
		parts := make([]string, 0, len(shenSha))
		for _, s := range shenSha {
			parts = append(parts, shenShaLine(s))
		}
		fmt.Printf("神煞: %s\n", strings.Join(parts, "、"))
	}

	for _, row := range q.GridLayout() {
		for _, palace := range row {
			if palace.Number == 5 {
				continue
			}
			printPalace(palace)
		}
		fmt.Println()
	}
}

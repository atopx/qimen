// Package qimen 提供奇门遁甲 (Qimen Dunjia) 的 Go 实现, 基于 6tail/tyme4go 历法计算。
//
// # 起局
//
//	import (
//	    "github.com/atopx/qimen"
//	    "github.com/6tail/tyme4go/tyme"
//	)
//
//	st, _ := tyme.SolarTime{}.FromYmdHms(2026, 1, 14, 18, 45, 0)
//	q := qimen.FromSolarTime(*st)
//
//	fmt.Printf("节气: %s\n", q.Term().GetName())
//	fmt.Printf("局数: %s遁%d局\n", q.YinYang(), q.Ju())
//	fmt.Printf("旬首: %s\n", q.XunShou().GetName())
//	for _, p := range q.Patterns() {
//	    fmt.Printf("格局: %s [%s] - %s\n", p.Name(), p.Auspice(), p.Summary())
//	}
//	for _, s := range q.ShenSha() {
//	    fmt.Printf("神煞: %s [%s]\n", s.Name(), s.Auspice())
//	}
//
// # 统一接口
//
// 所有领域实体 (Pattern / Hexagram / ShenSha / ShenShaKind / Terrain)
// 都实现 [Auspicious] 接口, 提供 Name() / Summary() / Auspice() 三个方法。
//
// # 维度
//
//   - 盘面: 三奇六仪、天盘、暗干、九星、八门、九神
//   - 格局 [Pattern]: 二十四格局 (反吟/伏吟/入墓/落空亡/门迫/三奇得使/八遁/青龙返首/飞鸟跌穴/...)
//   - 神煞 [ShenSha]: 驿马、桃花、华盖、天乙贵人、天德贵人、月德贵人、国印贵人、文昌、禄神、羊刃
//   - 五行 [Element] / [ElementRelation]: 含九星/八门/天干/地支/宫位的五行映射
//   - 十神: 复用 tyme4go 的 TenStar, 以日柱天干为日主
//   - 长生十二宫 [Terrain]: 包装 tyme4go.Terrain 以接入 [Auspicious]
//   - 六十四卦 [Hexagram]: 门宫演卦, 门为上卦, 宫为下卦
package qimen

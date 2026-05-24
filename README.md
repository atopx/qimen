# qimen

`qimen` 是零依赖的 Go 奇门遁甲 (时家三元) 起局库。内含自包含的历法引擎 (节气、朔望、四柱、农历) 和完整的盘面演算 (三奇六仪/天盘/暗干/九星/八门/九神), 提供二十四格局识别、十种神煞查算、六十四卦演卦、长生十二宫与十神的统一领域接口。

## 安装

```bash
go get github.com/atopx/qimen
```

## 快速起局

```go
package main

import (
    "fmt"
    "time"

    "github.com/atopx/qimen"
    "github.com/atopx/qimen/almanac"
)

func main() {
    // 4 种构造入口任选其一
    chart := qimen.New()  // 当前时刻
    st, _ := almanac.SolarTimeOf(2026, 1, 14, 18, 45, 0) // 指定时刻
    chart, _  = qimen.From(st) // 内部 SolarTime
    chart, _  = qimen.FromTime(time.Now()) // 标准库 time.Time
    chart, _  = qimen.FromTimestamp(1768376700) // Unix 秒

    fmt.Printf("节气: %s\n", chart.Term().Name())
    fmt.Printf("局数: %s遁%d局 %s\n", chart.YinYang().Name(), chart.Ju(), chart.Yuan().Name())
    fmt.Printf("旬首: %s\n", chart.XunShou().Name())
    fmt.Printf("值符: %s 落 %d 宫\n", chart.ZhiFu().Star.Name(), chart.ZhiFu().Palace)
    fmt.Printf("值使: %s 落 %d 宫\n", chart.ZhiShi().Door.Name(), chart.ZhiShi().Palace)

    for p := range chart.Patterns() {
        fmt.Printf("格局: %s [%s] - %s\n", p.Name(), p.Auspice().Name(), p.Summary())
    }
    for s := range chart.ShenSha() {
        fmt.Printf("神煞: %s [%s] - %s\n", s.Name(), s.Auspice().Name(), s.Summary())
    }
}
```

完整示例: `examples/deduce/main.go`, 运行:

```bash
go run ./examples/deduce
```

## 起局选项

```go
chart, err := qimen.From(st,
    qimen.WithMethod(enum.MethodTime),    // 起局法门, 默认 MethodTime
    qimen.WithStyle(enum.StyleRotate),    // 盘式, 默认 StyleRotate
)
```

未实现的法门/盘式 (`MethodDay/Month/Year`, `StyleFly/SiZhu`) 返回 sentinel error, 可用 `errors.Is`:

```go
if errors.Is(err, qimen.ErrUnsupportedMethod) { /* ... */ }
if errors.Is(err, qimen.ErrUnsupportedStyle)  { /* ... */ }
```

## 顶层 `*Chart` API

| 方法                                                     | 返回                               | 说明                                |
| -------------------------------------------------------- | ---------------------------------- | ----------------------------------- |
| `SolarTime()`                                            | `almanac.SolarTime`                | 起局阳历时刻                        |
| `LunarDay()`                                             | `almanac.LunarDay`                 | 对应农历日                          |
| `Year/Month/Day/Hour()`                                  | `almanac.Cycle`                    | 四柱六十甲子                        |
| `Term()`                                                 | `almanac.Term`                     | 当前节气                            |
| `YinYang()`                                              | `almanac.YinYang`                  | 阴/阳遁                             |
| `Ju()`                                                   | `uint8`                            | 局数 (1..=9)                        |
| `Yuan()`                                                 | `enum.Yuan`                        | 三元                                |
| `XunShou()`                                              | `almanac.Stem`                     | 旬首天干                            |
| `ZhiFu()`                                                | `qimen.Duty`                       | 值符 (星 + 原宫 + 落宫)             |
| `ZhiShi()`                                               | `qimen.DutyDoor`                   | 值使 (门 + 原宫 + 落宫)             |
| `KongWang()`                                             | `[2]almanac.Branch`                | 旬空亡两支                          |
| `Palace(n)`                                              | `*palace.Palace`                   | 按宫号取宫 (越界 nil)               |
| `Palaces()`                                              | `iter.Seq2[uint8, *palace.Palace]` | 1..9 流式枚举                       |
| `Grid()`                                                 | `[3][3]*palace.Palace`             | 三行三列展示 (巽离坤/震中兑/艮坎乾) |
| `StemPalace(stem)` / `SelfPalace()` / `OpponentPalace()` | `uint8`                            | 用神查宫                            |
| `EarthStems/HeavenStems/HiddenStems()`                   | `iter.Seq2[uint8, almanac.Stem]`   | 三盘干流式枚举                      |
| `Patterns()`                                             | `iter.Seq[pattern.Pattern]`        | 全盘格局流                          |
| `ShenSha()`                                              | `iter.Seq[shensha.ShenSha]`        | 全盘神煞流                          |

## 单宫 `*palace.Palace`

- 基础: `Number`, `Name`, `Direction`, `Branches`
- 盘面: `EarthStem`, `SanQiLiuYi`, `HeavenStem`, `HiddenStem`
- 三盘 (中宫 5 多为未设): `(Star, StarSet)`, `(Door, DoorSet)`, `(God, GodSet)`
- 衍生: `(TenStar, TenStarSet)`, `(Terrain, TerrainSet)`, `(Hexagram, HexagramSet)`
- 聚合: `Patterns`, `ShenSha`
- 关系: `DoorPalaceRelation()`, `StarPalaceRelation()` 返回 `(Relation, bool)`

## 子包速览

| 包         | 用途                                                                    |
| ---------- | ----------------------------------------------------------------------- |
| `almanac`  | 历法引擎: SolarTime / Term / Cycle / Stem / Branch / Pillars / LunarDay |
| `auspice`  | `Auspice` 5 级吉凶 + `Auspicable` 统一接口                              |
| `enum`     | `Method`, `Style`, `Yuan`, `Star`, `Door`, `God` 业务枚举               |
| `element`  | 五行 `Element` + 生克关系 `Relation`                                    |
| `terrain`  | `Terrain` 长生十二宫 + auspice                                          |
| `hexagram` | `Trigram` 八卦 + `Hexagram` 六十四卦 + auspice                          |
| `plate`    | 泛型 `Plate[T]` + 六盘构造器 (`BuildEarth/Heaven/Star/Door/God/Hidden`) |
| `palace`   | `Palace` 单宫数据载体                                                   |
| `pattern`  | 二十四格局检测 (`Detect → iter.Seq[Pattern]`)                           |
| `shensha`  | 十种神煞检测 (`Detect → iter.Seq[ShenSha]`)                             |

## 覆盖能力

- **盘面**: 三奇六仪、天盘、暗干、九星、八门、九神
- **格局**: 二十四格局 (反吟/伏吟/入墓/落空亡/门迫/三奇得使/八遁/青龙返首/飞鸟跌穴/大格/小格/刑格/悖格/天网四张)
- **神煞**: 驿马、桃花、华盖、天乙贵人、天德贵人、月德贵人、国印贵人、文昌、禄神、羊刃
- **五行生克**: 九星 / 八门 / 天干 / 地支 / 宫位的五行映射与关系判定
- **六十四卦**: 门宫演卦, 门为上卦, 宫为下卦
- **长生十二宫**: `Terrain` 接入 `Auspicable`
- **历法**: 24 节气 (秒级)、朔望、四柱直推、农历日期 (含闰月)

## 测试

```bash
go test ./...
go test -race -cover ./...
```

## Benchmark

```bash
go test -bench=. -benchmem -run='^$' -benchtime=3s
```

涵盖 `From` 起盘耗时、`Patterns/ShenSha/Grid` 聚合调用,以及 `Palace(n)` 索引访问。基准时刻覆盖一年中分布在不同节气、阴阳遁、上中下三元的代表性样本 (见 [bench_test.go](bench_test.go))。


## License

MIT — 见 [LICENSE](LICENSE)。

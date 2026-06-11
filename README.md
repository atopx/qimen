# qimen

`qimen` 是零依赖的 Go 奇门遁甲起局库, 支持时/日/月/年四家法门与转盘/飞盘两种盘式。内含自包含的历法引擎 (节气、朔望、四柱、农历、紫白) 和完整的盘面演算 (三奇六仪/天盘/暗干/九星/八门/九神), 提供二十四格局识别、十种神煞查算、六十四卦演卦、长生十二宫与十神的统一领域接口。时家置闰盘面经 16 组权威排盘结果逐项校验。

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
    chart := qimen.New()  // 当前时刻 (UTC+8), 默认配置
    st, _ := almanac.SolarTimeOf(2026, 1, 14, 18, 45, 0) // 指定时刻
    chart  = qimen.From(st) // 内部 SolarTime, 全配置组合可构造, 无 error
    chart, _  = qimen.FromTime(time.Now()) // 标准库 time.Time, 取其挂钟时间
    chart, _  = qimen.FromTimestamp(1768376700) // Unix 秒, 按 UTC+8 解释

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
chart := qimen.From(st,
    qimen.WithMethod(enum.MethodTime),     // 起局法门: 时/日/月/年家, 默认时家
    qimen.WithStyle(enum.StyleRotate),     // 盘式: 转盘/飞盘, 默认转盘
    qimen.WithJuRule(enum.JuRuleChaiBu),   // 时家定局规则, 默认置闰
)
```

全部法门/盘式/定局规则组合均已实现, `From` 为全函数 (无 error)。
`FromTime`/`FromTimestamp` 仅在历法输入越界时返回 `qimen.ErrInvalidTime`
(与 `almanac.ErrInvalidTime` 同一哨兵, `errors.Is` 任一可匹配)。

### 法门与盘式

- **时家** (默认): 主柱为时柱, 局数由节气三元表给出;
- **日家**: 主柱为日柱, 局数与时家同源 (节气三元本就按日定局), 同样支持置闰/拆补;
- **月家**: 主柱为月柱, 恒阴遁; 寅月起局 8/5/2 (子午卯酉/辰戌丑未/寅申巳亥年),
  逐月逆行;
- **年家**: 主柱为年柱, 恒阴遁; 局按 60 年元固定 (上元一局/中元四局/下元七局,
  1864 上元甲子起), 立春换年;
- **转盘** (默认): 洛书环刚性转动, 中宫干与天禽寄坤二, 八神;
- **飞盘**: 按 1..9 宫序数飞, 中宫实际参与, 天禽落实宫, 九星布满九宫,
  九神 (含**太常**) 布满九宫, 八门留一宫空 (`Palace.StarSet/DoorSet/GodSet`
  标志存在性)。

日/月/年家与飞盘盘面已通过权威排盘结果逐项校验 (共 21 组金标准,
见 `qimen_golden_test.go`)。

### 时区语义

起局使用事件发生地的挂钟时间 (民用时):

- `FromTime(t)` 直接取 `t` 的挂钟时间, **不做时区换算**;
- `FromTimestamp(unix)` 与 `New()` 按 UTC+8 (中国标准时间) 解释;
- 其他时区先显式转换: `qimen.FromTime(time.Unix(unix, 0).In(loc))`。

### 排盘约定

流派分歧点的取舍如下 (时家转盘主流做法):

- 三元由日柱符头网格直推 (`index mod 15`), 两种定局规则共用;
- 时家定局规则默认**置闰** (`WithJuRule` 可切换**拆补**):
  符头与二至对齐, 超神/接气, 符头超前满九日 (含符头日的传统计数) 时在芒种/大雪置闰;
  `chart.JuTerm()` 返回实际用局节气, `chart.Term()` 恒为天文节气;
  置闰盘面已通过权威排盘结果逐项校验 (见 `qimen_golden_test.go`);
- 暗干随值使转动: 八门转几步, 地盘干随转几步 (门下藏干);
- 日柱 23 点换日 (晚子时归次日);
- 天盘中宫干原位保留, 坤二宫不另显寄宫双干。

## 顶层 `*Chart` API

| 方法                                                     | 返回                               | 说明                                |
| -------------------------------------------------------- | ---------------------------------- | ----------------------------------- |
| `SolarTime()`                                            | `almanac.SolarTime`                | 起局阳历时刻                        |
| `LunarDay()`                                             | `almanac.LunarDay`                 | 对应农历日                          |
| `Year/Month/Day/Hour()`                                  | `almanac.Cycle`                    | 四柱六十甲子                        |
| `Term()`                                                 | `almanac.Term`                     | 当前天文节气                        |
| `JuTerm()`                                               | `almanac.Term`                     | 用局节气 (拆补下同 `Term()`)        |
| `YinYang()`                                              | `almanac.YinYang`                  | 阴/阳遁                             |
| `Ju()`                                                   | `uint8`                            | 局数 (1..=9)                        |
| `Yuan()`                                                 | `enum.Yuan`                        | 三元                                |
| `XunShou()`                                              | `almanac.Stem`                     | 旬首天干                            |
| `ZhiFu()`                                                | `qimen.DutyStar`                   | 值符 (星 + 原宫 + 落宫)             |
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
- 盘面: `EarthStem` (三奇六仪), `HeavenStem`, `HiddenStem`
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
| `pattern`  | 二十四格局检测 (`AppendAll` / `Detect → iter.Seq[Pattern]`)             |
| `shensha`  | 十种神煞检测 (`AppendAll` / `Detect → iter.Seq[ShenSha]`)               |

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

# qimen

`qimen` 是基于 [`6tail/tyme4go`](https://github.com/6tail/tyme4go) 历法的 Go 奇门遁甲 (时家三元) 起局库。提供盘面构造、二十四格局识别、十种神煞查算、六十四卦 (门宫演卦)、长生十二宫、十神等统一的领域接口。

## 安装

```bash
go get github.com/atopx/qimen
```

依赖 `github.com/6tail/tyme4go` 提供历法计算; Go 1.26+ (使用泛型)。

## 快速起局

```go
package main

import (
    "fmt"

    "github.com/6tail/tyme4go/tyme"
    "github.com/atopx/qimen"
)

func main() {
    st, _ := tyme.SolarTime{}.FromYmdHms(2026, 1, 14, 18, 45, 0)
    q := qimen.FromSolarTime(*st)

    fmt.Printf("节气: %s\n", q.Term().GetName())
    fmt.Printf("局数: %s遁%d局 %s\n", q.YinYang(), q.Ju(), q.Yuan())
    fmt.Printf("旬首: %s\n", q.XunShou().GetName())
    fmt.Printf("值符: %s 落 %d 宫\n", q.ZhiFu().Star, q.ZhiFu().Palace)
    fmt.Printf("值使: %s 落 %d 宫\n", q.ZhiShi().Door, q.ZhiShi().Palace)

    for _, p := range q.Patterns() {
        fmt.Printf("格局: %s [%s] - %s\n", p.Name(), p.Auspice(), p.Summary())
    }
    for _, s := range q.ShenSha() {
        fmt.Printf("神煞: %s [%s] - %s\n", s.Name(), s.Auspice(), s.Summary())
    }
}
```

完整示例: 见 `examples/main.go`, 运行:

```bash
go run ./examples
```

## API 速查

### 起局

| 方法 | 说明 |
|---|---|
| `FromSolarTime(t tyme.SolarTime) *Qimen` | 默认参数 (时家三元) 起局; 默认参数不会失败 |
| `FromSolarTimeWithOptions(t, opts) (*Qimen, error)` | 自定义参数起局; 仅 `QimenMethodTime` + `QimenChartTypeSanYuan` 已实现 |

### `*Qimen` 顶层访问

| 方法 | 返回 | 说明 |
|---|---|---|
| `SolarTime()` | `tyme.SolarTime` | 起局阳历时刻 |
| `Year/Month/Day/Hour()` | `tyme.SixtyCycle` | 四柱 |
| `Term()` | `tyme.SolarTerm` | 节气 |
| `YinYang()` | `tyme.YinYang` | 阴阳遁 |
| `Ju()` | `uint8` | 局数 (1..=9) |
| `Yuan()` | `QimenYuan` | 三元 |
| `XunShou()` | `tyme.HeavenStem` | 旬首天干 |
| `ZhiFu()` | `QimenDutyStar` | 值符 (星 + 原宫 + 落宫) |
| `ZhiShi()` | `QimenDutyDoor` | 值使 (门 + 原宫 + 落宫) |
| `KongWang()` | `[2]tyme.EarthBranch` | 旬空亡两支 |
| `Palace(n uint8)` | `*QimenPalace` | 取指定宫位; 越界返回 nil |
| `Palaces()` | `[9]*QimenPalace` | 全部 9 宫 |
| `GridLayout()` | `[3][3]*QimenPalace` | 三行三列展示 (巽离坤/震中兑/艮坎乾) |
| `SanQiLiuYi() / TianPan() / HiddenHeavenStems()` | `[]QimenHeavenStemPlacement` | 三奇六仪 / 天盘 / 暗干 |
| `StemPalace(stem) / SelfPalace() / OpponentPalace()` | `uint8` | 用神查宫 |
| `Patterns()` | `[]Pattern` | 全盘格局 |
| `ShenSha()` | `[]ShenSha` | 全盘神煞 |

### 单宫 `*QimenPalace`

- 基础字段: `Number`, `PalaceName`, `Direction`, `EarthBranches`
- 盘面字段: `EarthHeavenStem`, `SanQiLiuYi`, `HeavenHeavenStem`, `HiddenHeavenStem`
- 三盘: `Star *QimenStar`, `Door *QimenDoor`, `God *QimenGod` (中宫 5 均为 nil)
- 衍生属性: `TenStar *tyme.TenStar`, `TerrainValue *Terrain`, `Hexagram *Hexagram`
- 聚合: `Patterns []Pattern`, `ShenSha []ShenSha`
- 关系派生: `DoorPalaceRelation()`, `StarPalaceRelation()` 返回 `*PalaceRelation`

### 统一 `Auspicious` 接口

`Pattern`, `Hexagram`, `Terrain`, `ShenSha`, `ShenShaKind` 均实现:

```go
type Auspicious interface {
    Name() string
    Summary() string
    Auspice() Auspice
}
```

`Auspice` 5 级: `AuspiceGreatAuspicious`, `AuspiceAuspicious`, `AuspiceNeutral`, `AuspiceInauspicious`, `AuspiceGreatInauspicious`。

## 覆盖能力

- **盘面**: 三奇六仪、天盘、暗干、九星、八门、九神
- **格局**: 二十四格局 (反吟/伏吟/入墓/落空亡/门迫/三奇得使/八遁/青龙返首/飞鸟跌穴/大格/小格/刑格/悖格/天网四张)
- **神煞**: 驿马、桃花、华盖、天乙贵人、天德贵人、月德贵人、国印贵人、文昌、禄神、羊刃
- **五行生克**: 含九星/八门/天干/地支/宫位的五行映射与关系判定
- **六十四卦**: 门宫演卦, 门为上卦, 宫为下卦
- **长生十二宫**: 包装 `tyme.Terrain` 接入 `Auspicious`

## 测试

```bash
go test ./...
go test -race -cover ./...
```

## Benchmark

```bash
go test -bench=. -benchmem -run=^$ -benchtime=3s
```

涵盖 `FromSolarTime` 起盘耗时、`Patterns/ShenSha/GridLayout` 聚合调用,以及 `Palace(n)` 索引访问。基准时刻覆盖一年中分布在不同节气、阴阳遁、上中下三元的代表性样本(见 [bench_test.go](bench_test.go))。

## 跨端对照测试 (Rust ↔ Go)

仓库内并存 Rust crate 实现(`src/*.rs`,依赖 `tyme4rs`)与 Go 端口(根目录 `*.go`,依赖 `tyme4go`)。`examples/dump.rs` + `verify_dump_test.go` 用于验证两端同一时刻输出完全一致。

```bash
# 1. 用 Rust 端导出指定年份每天 12 时辰的 JSONL (中点采样,4380 行/年)
mkdir -p testdata
cargo run --release --example dump -- 2025 > testdata/qimen_2025.jsonl

# 2. 用 Go 端逐行重建并比对所有字段 (env var gated, 默认 SKIP)
QIMEN_VERIFY_DUMP=testdata/qimen_2025.jsonl \
    go test -run TestVerifyDumpAgainstRust -v -timeout 5m
```

比对覆盖:四柱、节气、阴阳遁/局数/三元/旬首、值符值使、空亡;每宫的 16 个字段(基础 + 三盘 + 十神 + 长生 + 卦象 + 门宫/星宫关系);全盘聚合 patterns / shen_sha (按 detail_key 与 (palace, kind, target) 排序后比对)。`testdata/` 已加入 `.gitignore`,不入仓。

## License

MIT — 见 [LICENSE](LICENSE)。

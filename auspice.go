package qimen

// Auspice 吉凶等级 (5 级)。这是格局/神煞自身在传统文献中的固定吉凶定性,
// 而非由计算得出的得分。调用方仍可结合具体盘面、用神生克自行加减判断;
// 本库不做任何动态评分。
type Auspice int

const (
	// AuspiceGreatAuspicious 大吉 — 极佳格局。
	AuspiceGreatAuspicious Auspice = iota
	// AuspiceAuspicious 吉 — 寻常吉格。
	AuspiceAuspicious
	// AuspiceNeutral 中和 — 平稳无显著吉凶。
	AuspiceNeutral
	// AuspiceInauspicious 凶 — 普通凶格。
	AuspiceInauspicious
	// AuspiceGreatInauspicious 大凶 — 极差格局。
	AuspiceGreatInauspicious
)

// Name 中文名称 (大吉/吉/中和/凶/大凶)。
func (a Auspice) Name() string {
	switch a {
	case AuspiceGreatAuspicious:
		return "大吉"
	case AuspiceAuspicious:
		return "吉"
	case AuspiceNeutral:
		return "中和"
	case AuspiceInauspicious:
		return "凶"
	case AuspiceGreatInauspicious:
		return "大凶"
	default:
		return ""
	}
}

// String 实现 fmt.Stringer。
func (a Auspice) String() string { return a.Name() }

// IsAuspicious 是否为吉象 (大吉或吉)。
func (a Auspice) IsAuspicious() bool {
	return a == AuspiceGreatAuspicious || a == AuspiceAuspicious
}

// IsInauspicious 是否为凶象 (大凶或凶)。
func (a Auspice) IsInauspicious() bool {
	return a == AuspiceGreatInauspicious || a == AuspiceInauspicious
}

// IsExtreme 是否为极端 (大吉或大凶)。
func (a Auspice) IsExtreme() bool {
	return a == AuspiceGreatAuspicious || a == AuspiceGreatInauspicious
}

// Auspicious 是具有"名称 / 客观描述 / 吉凶等级"三个属性的领域实体的统一接口。
//
// 库内所有"格局/神煞/卦/长生十二宫"等领域实体都实现该接口。
type Auspicious interface {
	// Name 中文名称。
	Name() string
	// Summary 客观一句话描述。
	Summary() string
	// Auspice 传统文献既定的吉凶等级 (5 级)。
	Auspice() Auspice
}

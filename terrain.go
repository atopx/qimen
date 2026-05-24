package qimen

import "github.com/6tail/tyme4go/tyme"

// Terrain 长生十二宫的本包包装, 实现 [Auspicious] 接口。
//
// 12 宫: 长生、沐浴、冠带、临官、帝旺、衰、病、死、墓、绝、胎、养。
type Terrain struct {
	Inner tyme.Terrain
}

// NewTerrain 包装 tyme.Terrain。
func NewTerrain(t tyme.Terrain) Terrain { return Terrain{Inner: t} }

// Index 在 12 宫中的索引 (0..11)。
func (t Terrain) Index() int { return t.Inner.GetIndex() }

// Name 中文名 (长生/沐浴/...)。
func (t Terrain) Name() string {
	i := t.Index()
	if i < 0 || i >= len(terrainNames) {
		return ""
	}
	return terrainNames[i]
}

// Summary 客观描述。
func (t Terrain) Summary() string {
	i := t.Index()
	if i < 0 || i >= len(terrainSummary) {
		return ""
	}
	return terrainSummary[i]
}

// Auspice 传统易学既定的吉凶定性。
//
//   - 长生 / 冠带 / 临官 / 养 → 吉
//   - 帝旺 → 大吉
//   - 沐浴 / 衰 / 胎 → 中和
//   - 病 / 墓 → 凶
//   - 死 / 绝 → 大凶
func (t Terrain) Auspice() Auspice {
	i := t.Index()
	if i < 0 || i >= len(terrainAuspice) {
		return AuspiceNeutral
	}
	return terrainAuspice[i]
}

// String 实现 fmt.Stringer。
func (t Terrain) String() string { return t.Name() }

var terrainNames = [12]string{
	"长生", "沐浴", "冠带", "临官", "帝旺", "衰", "病", "死", "墓", "绝", "胎", "养",
}

var terrainAuspice = [12]Auspice{
	AuspiceAuspicious,        // 0  长生
	AuspiceNeutral,           // 1  沐浴
	AuspiceAuspicious,        // 2  冠带
	AuspiceAuspicious,        // 3  临官
	AuspiceGreatAuspicious,   // 4  帝旺
	AuspiceNeutral,           // 5  衰
	AuspiceInauspicious,      // 6  病
	AuspiceGreatInauspicious, // 7  死
	AuspiceInauspicious,      // 8  墓
	AuspiceGreatInauspicious, // 9  绝
	AuspiceNeutral,           // 10 胎
	AuspiceAuspicious,        // 11 养
}

var terrainSummary = [12]string{
	"长生十二宫之首。万物初生之态,如人之出生,充满生机与希望。",
	"长生十二宫第二位。又名\"桃花\",如人初生后沐浴净身,象征不稳定与桃花。",
	"长生十二宫第三位。如人成年加冠,开始走向社会,代表成长、自信。",
	"长生十二宫第四位。又名\"禄\",如人入仕做官,代表壮年、权力。",
	"长生十二宫第五位。五行力量巅峰,如帝王之旺,极盛之态。盛极必衰。",
	"长生十二宫第六位。帝旺之后力量开始衰退,如人步入中年体力下降。",
	"长生十二宫第七位。力量进一步衰弱,如人生病,需休养调理。",
	"长生十二宫第八位。万物气绝,如人之死亡,代表终结与静止。",
	"长生十二宫第九位。又名\"库\",如人入土归葬,代表收藏、积蓄、隐藏。",
	"长生十二宫第十位。万物气息断绝,旧形全消,但蕴含新生的转机。",
	"长生十二宫第十一位。新生命在母体中孕育,代表新计划、新希望的萌芽。",
	"长生十二宫第十二位。胎儿在母体中成长发育,即将出世,代表蓄养、培育。",
}

package qimen

import (
	"fmt"

	"github.com/6tail/tyme4go/tyme"
)

// QimenOptions 起局参数。
type QimenOptions struct {
	// Method 起局方法 (时家/日家/...)。
	Method QimenMethod
	// ChartType 盘式 (三元/四柱)。
	ChartType QimenChartType
}

// DefaultOptions 默认起局参数 (时家三元)。
func DefaultOptions() QimenOptions {
	return QimenOptions{Method: QimenMethodTime, ChartType: QimenChartTypeSanYuan}
}

// QimenDutyStar 值符: 落在某宫的九星 (含原宫与落宫)。
type QimenDutyStar struct {
	Star           QimenStar
	OriginalPalace uint8
	Palace         uint8
}

// QimenDutyDoor 值使: 落在某宫的八门 (含原宫与落宫)。
type QimenDutyDoor struct {
	Door           QimenDoor
	OriginalPalace uint8
	Palace         uint8
}

// QimenHeavenStemPlacement 某天干在某宫的位置。
type QimenHeavenStemPlacement struct {
	Palace     uint8
	HeavenStem tyme.HeavenStem
}

// QimenPalace 单宫数据 — 既包含盘面客观信息, 也聚合该宫所有衍生属性
// (十神/长生/卦象/格局/神煞)。
type QimenPalace struct {
	Number           uint8
	PalaceName       string
	Direction        tyme.Direction
	EarthBranches    []tyme.EarthBranch
	EarthHeavenStem  tyme.HeavenStem
	SanQiLiuYi       tyme.HeavenStem
	HeavenHeavenStem tyme.HeavenStem
	HiddenHeavenStem tyme.HeavenStem
	Star             *QimenStar
	Door             *QimenDoor
	God              *QimenGod
	TenStar          *tyme.TenStar
	TerrainValue     *Terrain
	Hexagram         *Hexagram
	Patterns         []Pattern
	ShenSha          []ShenSha
}

// Element 宫位五行属性。
func (p *QimenPalace) Element() Element { return ElementFromPalace(p.Number) }

// DoorPalaceRelation 八门与所在宫位的五行生克关系。中宫无门返回 nil。
//
// 吉凶以**门**为用神视角 — 受生为吉、受克为凶、耗泄为凶、其余中和。
func (p *QimenPalace) DoorPalaceRelation() *PalaceRelation {
	if p.Door == nil {
		return nil
	}
	r := palaceRelationForSubject(p.Door.Name(), p.Door.Element(), p.Element())
	return &r
}

// StarPalaceRelation 九星与所在宫位的五行生克关系。中宫无星返回 nil。
func (p *QimenPalace) StarPalaceRelation() *PalaceRelation {
	if p.Star == nil {
		return nil
	}
	r := palaceRelationForSubject(p.Star.Name(), p.Star.Element(), p.Element())
	return &r
}

// PalaceRelation 主体 (门 / 星) 与所在宫位的五行生克关系结果。
type PalaceRelation struct {
	Description     string
	ElementRelation ElementRelation
	AuspiceLevel    Auspice
}

// Auspice 以主体为用神视角的吉凶等级。
func (r PalaceRelation) Auspice() Auspice { return r.AuspiceLevel }

// String 实现 fmt.Stringer。
func (r PalaceRelation) String() string {
	return fmt.Sprintf("%s[%s]", r.Description, r.AuspiceLevel.Name())
}

// palaceRelationForSubject 构造主体—宫位关系。
func palaceRelationForSubject(subjectName string, subjectEl, palaceEl Element) PalaceRelation {
	rel := subjectEl.RelationTo(palaceEl)
	var desc string
	switch rel {
	case ElementRelationSame:
		desc = subjectName + "与宫比和"
	case ElementRelationGenerated:
		desc = "宫生" + subjectName
	case ElementRelationGenerates:
		desc = subjectName + "生宫"
	case ElementRelationRestrained:
		desc = "宫克" + subjectName
	case ElementRelationRestrains:
		desc = subjectName + "克宫"
	}
	return PalaceRelation{
		Description:     desc,
		ElementRelation: rel,
		AuspiceLevel:    rel.AuspiceAsSelf(),
	}
}

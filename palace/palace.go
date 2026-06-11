// Package palace defines the per-cell Palace value that aggregates
// every derived attribute (stems, star, door, god, ten-star, terrain,
// hexagram, patterns, shensha) for one of the 9 grid positions.
package palace

import (
	"fmt"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/auspice"
	"github.com/atopx/qimen/element"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/hexagram"
	"github.com/atopx/qimen/pattern"
	"github.com/atopx/qimen/shensha"
	"github.com/atopx/qimen/terrain"
)

// Palace holds the data for one of the 9 grid cells of a qimen chart.
//
// Field validity follows a single invariant: center palace (Number == 5)
// has NO Star / Door / God / TenStar / Terrain / Hexagram — these slots
// hold zero values. Every non-center palace has all six populated.
//
// Use [Palace.IsCenter] to gate access; the per-field `Set` flags that
// previous versions exposed are now redundant and have been removed.
type Palace struct {
	Number     uint8
	Name       string
	Direction  almanac.Direction
	Branches   []almanac.Branch
	EarthStem  almanac.Stem // 地盘干 (三奇六仪)
	HeavenStem almanac.Stem
	HiddenStem almanac.Stem

	// Populated iff Number != 5. Reading any of these for the center
	// palace yields zero values.
	Star     enum.Star
	Door     enum.Door
	God      enum.God
	TenStar  almanac.TenStar
	Terrain  terrain.Terrain
	Hexagram hexagram.Hexagram

	Patterns []pattern.Pattern
	ShenSha  []shensha.ShenSha
}

// IsCenter reports whether this is the center palace (Number == 5).
// Derived attributes (Star / Door / God / TenStar / Terrain / Hexagram)
// are zero for the center; callers should gate by !IsCenter().
func (p *Palace) IsCenter() bool { return p.Number == 5 }

// Element returns the 五行 of this palace.
func (p *Palace) Element() element.Element { return element.FromPalace(p.Number) }

// DoorPalaceRelation computes how the 八门 of this palace relates to
// its 宫位 五行. Returns the zero Relation for the center palace.
//
// Callers concerned with center-palace handling should check IsCenter()
// first; the zero value is a defined sentinel (empty Description).
func (p *Palace) DoorPalaceRelation() Relation {
	if p.IsCenter() {
		return Relation{}
	}
	return relationFromSubject(p.Door.Name(), element.OfDoor(p.Door), p.Element())
}

// StarPalaceRelation computes how the 九星 of this palace relates to
// its 宫位 五行. Returns the zero Relation for the center palace.
func (p *Palace) StarPalaceRelation() Relation {
	if p.IsCenter() {
		return Relation{}
	}
	return relationFromSubject(p.Star.Name(), element.OfStar(p.Star), p.Element())
}

// Relation is the result of a subject (door/star) vs. palace 五行 comparison.
type Relation struct {
	Description string
	Element     element.Relation
	Auspice     auspice.Auspice
}

// String formats as "<desc>[<auspice>]".
func (r Relation) String() string {
	return fmt.Sprintf("%s[%s]", r.Description, r.Auspice.Name())
}

// relationDescPrefix / Suffix wrap the subject name with the relation phrase.
// Indexed by element.Relation.
var (
	relationDescPrefix = [5]string{
		element.Same:       "",
		element.Generates:  "",
		element.Generated:  "宫生",
		element.Restrains:  "",
		element.Restrained: "宫克",
	}
	relationDescSuffix = [5]string{
		element.Same:       "与宫比和",
		element.Generates:  "生宫",
		element.Generated:  "",
		element.Restrains:  "克宫",
		element.Restrained: "",
	}
)

func relationFromSubject(subjectName string, subjEl, palaceEl element.Element) Relation {
	rel := subjEl.RelationTo(palaceEl)
	return Relation{
		Description: relationDescPrefix[rel] + subjectName + relationDescSuffix[rel],
		Element:     rel,
		Auspice:     rel.AuspiceAsSelf(),
	}
}

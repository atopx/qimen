// Package plate implements the 9-palace plate types and layout builders
// for qimen charts.
//
// Each plate stores up to 9 typed values (one per palace). Six concrete
// plates exist in any chart:
//
//   - 三奇六仪 (StemPlate, ground layer)
//   - 天盘 (StemPlate, heaven layer)
//   - 暗干 (StemPlate, hidden layer)
//   - 九星 (StarPlate)
//   - 八门 (DoorPlate)
//   - 九神 (GodPlate)
//
// The builders take an `enum.Style` parameter so future implementations
// of StyleFly can dispatch internally without changing call sites.
package plate

import (
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
)

// Plate is a generic 9-cell plate keyed by palace number (1..9).
// Unset cells return the zero value of T with ok=false.
type Plate[T any] struct {
	cells [9]T
	set   [9]bool
}

// New returns an empty plate.
func New[T any]() *Plate[T] { return &Plate[T]{} }

// Set stores v in palace n. Precondition: n ∈ [1, 9].
//
// No bounds check — qimen only calls Set with values from internal layout
// tables (LuoShuOrder, palace indices); an out-of-range value indicates
// a bug in the calling builder that should fail loudly via the Go runtime.
func (p *Plate[T]) Set(palace uint8, v T) {
	p.cells[palace-1] = v
	p.set[palace-1] = true
}

// Get returns (value, true) for a set cell, (zero, false) otherwise.
func (p *Plate[T]) Get(palace uint8) (T, bool) {
	var zero T
	if palace < 1 || palace > 9 {
		return zero, false
	}
	if !p.set[palace-1] {
		return zero, false
	}
	return p.cells[palace-1], true
}

// MustGet returns the value at palace, falling back to a zero value if unset.
//
// Useful for hot paths where the cell is known to be set.
func (p *Plate[T]) MustGet(palace uint8) T {
	var zero T
	if palace < 1 || palace > 9 || !p.set[palace-1] {
		return zero
	}
	return p.cells[palace-1]
}

// Type aliases for the four concrete plate flavors.
type (
	// StemPlate holds heavenly stems (used for 三奇六仪 / 天盘 / 暗干).
	StemPlate = Plate[almanac.Stem]
	// StarPlate holds nine-star values.
	StarPlate = Plate[enum.Star]
	// DoorPlate holds eight-door values (center palace is unset).
	DoorPlate = Plate[enum.Door]
	// GodPlate holds nine-god values (center palace is unset).
	GodPlate = Plate[enum.God]
)

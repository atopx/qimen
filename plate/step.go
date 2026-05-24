package plate

import "github.com/atopx/qimen/almanac"

// StepPalace advances one step in 阴/阳 遁 order, skipping the center.
func StepPalace(palace uint8, yy almanac.YinYang) uint8 {
	if yy == almanac.Yang {
		if palace == 9 {
			return 1
		}
		return palace + 1
	}
	if palace == 1 {
		return 9
	}
	return palace - 1
}

// MoveBy steps n times in 阴/阳 遁 order. If the destination is center
// (palace 5), it is reassigned to 2 (坤 fallback per qimen tradition).
func MoveBy(palace uint8, steps int, yy almanac.YinYang) uint8 {
	target := palace
	for i := 0; i < steps; i++ {
		target = StepPalace(target, yy)
	}
	if target == 5 {
		return 2
	}
	return target
}

// AreOpposite reports whether two palaces are diametrically opposite.
//
//	1↔9, 2↔8, 3↔7, 4↔6.
func AreOpposite(a, b uint8) bool {
	switch {
	case a == 1 && b == 9, a == 9 && b == 1:
		return true
	case a == 2 && b == 8, a == 8 && b == 2:
		return true
	case a == 3 && b == 7, a == 7 && b == 3:
		return true
	case a == 4 && b == 6, a == 6 && b == 4:
		return true
	}
	return false
}

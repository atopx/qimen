package almanac

import "fmt"

// LunarHour identifies a 时辰 within a LunarDay.
type LunarHour struct {
	Day  LunarDay
	Hour uint8 // 0..23 (solar hour)
}

// Name returns the Chinese 时辰 name, e.g. "子时".
func (h LunarHour) Name() string {
	b := hourBranchIndex(int(h.Hour))
	return BranchOf(b).Name() + "时"
}

// String returns "<day><hour>".
func (h LunarHour) String() string {
	return fmt.Sprintf("%s%s", h.Day, h.Name())
}

// Branch returns the earth branch of the time.
func (h LunarHour) Branch() Branch {
	return BranchOf(hourBranchIndex(int(h.Hour)))
}

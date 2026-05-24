// Package enum holds the qimen-specific domain enumerations
// (Method / Style / Yuan / Star / Door / God).
//
// All values are exported with type-qualifying prefixes (`MethodTime`,
// `StarTianPeng`, ...) so they read naturally as `enum.StarTianPeng`
// without the package name doubling up.
package enum

// Method 奇门起局法门 (time-based, day-based, ...).
type Method uint8

const (
	// MethodTime 时家 (default, currently the only fully implemented method).
	MethodTime Method = iota
	// MethodDay 日家 (reserved for future implementation).
	MethodDay
	// MethodMonth 月家 (reserved for future implementation).
	MethodMonth
	// MethodYear 年家 (reserved for future implementation).
	MethodYear
)

var methodNames = [4]string{"时家", "日家", "月家", "年家"}

// Name returns the Chinese label.
func (m Method) Name() string { return methodNames[m] }

// String implements fmt.Stringer.
func (m Method) String() string { return m.Name() }

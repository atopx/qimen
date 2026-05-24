package enum

// Yuan 三元 (上/中/下元), determined by the day pillar.
type Yuan uint8

const (
	// YuanUpper 上元
	YuanUpper Yuan = iota
	// YuanMiddle 中元
	YuanMiddle
	// YuanLower 下元
	YuanLower
)

var yuanNames = [3]string{"上元", "中元", "下元"}

// Name returns the Chinese label.
func (y Yuan) Name() string { return yuanNames[y] }

// String implements fmt.Stringer.
func (y Yuan) String() string { return y.Name() }

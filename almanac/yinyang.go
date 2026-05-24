package almanac

// YinYang denotes the yin/yang polarity of a stem, year, or other entity.
type YinYang uint8

const (
	// Yang 阳
	Yang YinYang = iota
	// Yin 阴
	Yin
)

// Name returns the Chinese character "阳" or "阴".
func (y YinYang) Name() string {
	if y == Yin {
		return "阴"
	}
	return "阳"
}

// String implements fmt.Stringer.
func (y YinYang) String() string { return y.Name() }

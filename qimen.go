package qimen

import (
	"fmt"
	"iter"
	"time"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/hexagram"
	"github.com/atopx/qimen/internal/compute"
	"github.com/atopx/qimen/internal/tables"
	"github.com/atopx/qimen/palace"
	"github.com/atopx/qimen/pattern"
	"github.com/atopx/qimen/plate"
	"github.com/atopx/qimen/shensha"
	"github.com/atopx/qimen/terrain"
)

// Chart is the result of one qimen 起局. Once built, all fields are
// read-only and safe to share across goroutines.
type Chart struct {
	cfg chartCfg

	// Pillars + solar context
	solar    almanac.SolarTime
	pillars  almanac.Pillars
	term     almanac.Term
	yinYang  almanac.YinYang
	ju       uint8
	yuan     enum.Yuan
	xunShou  almanac.Stem
	zhiFu    Duty
	zhiShi   DutyDoor
	kongWang [2]almanac.Branch
	lunarDay almanac.LunarDay // cached at build for O(1) LunarDay() access

	// Palaces (1..9; index = palace number - 1).
	// Stored by value to avoid 9 heap allocations per chart.
	palaces [9]palace.Palace
}

// config holds optional construction parameters set via Option.
// loc applies only at the entry points that accept a time / Unix instant;
// it is consumed and discarded before reaching build().
type config struct {
	method enum.Method
	style  enum.Style
	loc    *time.Location
}

// chartCfg is the slimmed config persisted inside Chart — only the
// fields needed by Method()/Style() accessors.
type chartCfg struct {
	method enum.Method
	style  enum.Style
}

// Option configures a Chart constructor.
type Option func(*config)

// WithMethod selects the 起局 method (currently only [enum.MethodTime]).
func WithMethod(m enum.Method) Option {
	return func(c *config) { c.method = m }
}

// WithStyle selects the chart style (currently only [enum.StyleRotate]).
func WithStyle(s enum.Style) Option {
	return func(c *config) { c.style = s }
}

// Duty is the 值符 entry: the 九星 currently acting as 值符, its
// originating palace, and the palace it has rotated into.
type Duty struct {
	Star           enum.Star
	OriginalPalace uint8
	Palace         uint8
}

// DutyDoor is the 值使 entry: the 八门 currently acting as 值使,
// its originating palace, and where it has rotated.
type DutyDoor struct {
	Door           enum.Door
	OriginalPalace uint8
	Palace         uint8
}

// defaultConfig returns the default construction options:
// 时家 / 转盘 / Asia/Shanghai (UTC+8).
func defaultConfig() config {
	return config{
		method: enum.MethodTime,
		style:  enum.StyleRotate,
		loc:    time.FixedZone("CST", 8*3600),
	}
}

// New builds a chart for the current solar instant in UTC+8.
//
// Construction with default options is infallible; this avoids forcing
// callers to handle an error in the most common code path.
func New(opts ...Option) *Chart {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	st := almanac.Now()
	chart, err := build(st, cfg)
	if err != nil {
		panic("qimen.New: default options must succeed: " + err.Error())
	}
	return chart
}

// From builds a chart from a [almanac.SolarTime].
func From(t almanac.SolarTime, opts ...Option) (*Chart, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return build(t, cfg)
}

// MustFrom is like From but panics on error. Useful for tests and
// scripts that supply hard-coded inputs.
func MustFrom(t almanac.SolarTime, opts ...Option) *Chart {
	chart, err := From(t, opts...)
	if err != nil {
		panic("qimen.MustFrom: " + err.Error())
	}
	return chart
}

// FromTime builds a chart from a standard library [time.Time]. The
// for this entry point.
func FromTime(t time.Time, opts ...Option) (*Chart, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	st, err := almanac.SolarTimeFromTime(t)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidTime, err)
	}
	return build(st, cfg)
}

// FromTimestamp builds a chart from a Unix-seconds timestamp.
// The default location is UTC+8.
func FromTimestamp(unix int64, opts ...Option) (*Chart, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	t := time.Unix(unix, 0).In(cfg.loc)
	st, err := almanac.SolarTimeFromTime(t)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidTime, err)
	}
	return build(st, cfg)
}

// validateConfig checks Method + Style at the single Chart entry point.
// Returns sentinel errors that wrap to errors.Is matchers.
func validateConfig(cfg config) error {
	if cfg.method != enum.MethodTime {
		return fmt.Errorf("%w: %s", ErrUnsupportedMethod, cfg.method.Name())
	}
	if cfg.style != enum.StyleRotate {
		return fmt.Errorf("%w: %s", ErrUnsupportedStyle, cfg.style.Name())
	}
	return nil
}

// build is the shared chart-construction routine.
func build(t almanac.SolarTime, cfg config) (*Chart, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	pillars := t.Pillars()
	term := t.Term()
	yinYang := compute.YinYang(t)
	yuan := compute.Yuan(pillars.Day)
	ju := compute.Ju(term, yuan)

	// Build the 6 plates. Style was validated at entry; primitives are
	// total functions assuming StyleRotate semantics.
	earth := plate.BuildEarth(yinYang, ju)

	xunShou := compute.XunShou(pillars.Hour)
	hourStem := pillars.Hour.Stem()

	zhiFuOriginalPalace := uint8(2)
	if p, ok := plate.FindStem(earth, xunShou, true); ok {
		zhiFuOriginalPalace = p
	}
	zhiFuPalace := plate.FindHourStem(earth, hourStem, zhiFuOriginalPalace)

	heaven := plate.BuildHeaven(earth, yinYang, zhiFuOriginalPalace, zhiFuPalace)
	stars := plate.BuildStar(zhiFuOriginalPalace, zhiFuPalace)
	doors := plate.BuildDoor(yinYang, zhiFuOriginalPalace, pillars.Hour)
	gods := plate.BuildGod(yinYang, zhiFuPalace)
	hidden := plate.BuildHidden(yinYang)

	// Resolve 值符 (star landing on zhiFuPalace).
	zhiFuStar := enum.StarQinRui
	if s, ok := stars.Get(zhiFuPalace); ok {
		zhiFuStar = s
	}
	// Resolve 值使 door (originating door from zhiFuOriginalPalace).
	// zhiFuOriginalPalace is sourced from earth-plate stem search which
	// only returns non-center palaces; DoorOfPalace is safe here.
	zhiShiDoor := enum.DoorDeath
	if zhiFuOriginalPalace != 5 {
		zhiShiDoor = enum.DoorOfPalace(zhiFuOriginalPalace)
	}
	zhiShiPalace := zhiFuPalace
	if p, ok := plate.FindDoor(doors, zhiShiDoor); ok {
		zhiShiPalace = p
	}

	kongWang := compute.KongWang(pillars.Hour)

	// Assemble 9 palaces in place (no heap allocations).
	var palaces [9]palace.Palace
	for i := 0; i < 9; i++ {
		n := uint8(i + 1)
		earthStem, _ := earth.Get(n)
		heavenStem := earthStem
		if v, ok := heaven.Get(n); ok {
			heavenStem = v
		}
		hiddenStem := earthStem
		if v, ok := hidden.Get(n); ok {
			hiddenStem = v
		}
		p := &palaces[i]
		p.Number = n
		p.Name = tables.PalaceNames[n]
		p.Direction = almanac.DirectionOfPalace(n)
		p.Branches = compute.BranchesForPalace(n)
		p.EarthStem = earthStem
		p.HeavenStem = heavenStem
		p.HiddenStem = hiddenStem
		p.SanQiLiuYi = earthStem
		// Star / Door / God only populated for non-center palaces
		// (BuildStar/Door/God iterate LuoShuOrder which excludes palace 5).
		if n != 5 {
			p.Star, _ = stars.Get(n)
			p.Door, _ = doors.Get(n)
			p.God, _ = gods.Get(n)
		}
	}

	// Derived attributes (十神 / 长生 / 六十四卦) — only for non-center palaces.
	dayStem := pillars.Day.Stem()
	for i := range palaces {
		p := &palaces[i]
		if p.IsCenter() {
			continue
		}
		p.TenStar = dayStem.TenStarOf(p.EarthStem)
		p.Terrain = terrain.Of(dayStem.TerrainOf(p.Branches[0]))
		// Hexagram: upper = door's home palace trigram, lower = this palace's trigram.
		lower := hexagram.TrigramOfPalace(p.Number)
		upper := hexagram.TrigramOfPalace(p.Door.HomePalace())
		p.Hexagram = hexagram.Of(upper, lower)
	}

	// Pattern detection (build input view, then dispatch results to palaces).
	pIn := pattern.DetectInput{
		ZhiFuOriginalPalace: zhiFuOriginalPalace,
		ZhiFuPalace:         zhiFuPalace,
		KongWang:            kongWang,
	}
	for i := 0; i < 9; i++ {
		p := &palaces[i]
		pIn.EarthStems[i] = p.EarthStem
		pIn.HeavenStems[i] = p.HeavenStem
		if !p.IsCenter() {
			pIn.Doors[i] = p.Door
			pIn.DoorsSet[i] = true
			pIn.Gods[i] = p.God
			pIn.GodsSet[i] = true
		}
		pIn.Branches[i] = p.Branches
	}
	for pat := range pattern.Detect(pIn) {
		n := pat.Palace
		if n >= 1 && n <= 9 {
			palaces[n-1].Patterns = append(palaces[n-1].Patterns, pat)
		}
	}

	// 神煞 detection.
	ssIn := shensha.DetectInput{
		YearStem:    pillars.Year.Stem(),
		MonthBranch: pillars.Month.Branch(),
		DayStem:     dayStem,
		DayBranch:   pillars.Day.Branch(),
	}
	for i := 0; i < 9; i++ {
		ssIn.EarthStems[i] = palaces[i].EarthStem
	}
	for ss := range shensha.Detect(ssIn) {
		n := ss.Palace
		if n >= 1 && n <= 9 {
			palaces[n-1].ShenSha = append(palaces[n-1].ShenSha, ss)
		}
	}

	return &Chart{
		cfg:      chartCfg{method: cfg.method, style: cfg.style},
		solar:    t,
		pillars:  pillars,
		term:     term,
		yinYang:  yinYang,
		ju:       ju,
		yuan:     yuan,
		xunShou:  xunShou,
		zhiFu:    Duty{Star: zhiFuStar, OriginalPalace: zhiFuOriginalPalace, Palace: zhiFuPalace},
		zhiShi:   DutyDoor{Door: zhiShiDoor, OriginalPalace: zhiFuOriginalPalace, Palace: zhiShiPalace},
		kongWang: kongWang,
		lunarDay: t.LunarDay(),
		palaces:  palaces,
	}, nil
}

// ===================== context accessors =====================

// SolarTime returns the solar instant used to build the chart.
func (c *Chart) SolarTime() almanac.SolarTime { return c.solar }

// LunarDay returns the lunar date of the chart's solar instant
// (cached at build for O(1) access).
func (c *Chart) LunarDay() almanac.LunarDay { return c.lunarDay }

// Method returns the chart's 起局法门 (currently always MethodTime).
func (c *Chart) Method() enum.Method { return c.cfg.method }

// Style returns the chart's 盘式 (currently always StyleRotate).
func (c *Chart) Style() enum.Style { return c.cfg.style }

// Year returns the 年柱 sixty cycle.
func (c *Chart) Year() almanac.Cycle { return c.pillars.Year }

// Month returns the 月柱 sixty cycle.
func (c *Chart) Month() almanac.Cycle { return c.pillars.Month }

// Day returns the 日柱 sixty cycle.
func (c *Chart) Day() almanac.Cycle { return c.pillars.Day }

// Hour returns the 时柱 sixty cycle.
func (c *Chart) Hour() almanac.Cycle { return c.pillars.Hour }

// Term returns the current solar term.
func (c *Chart) Term() almanac.Term { return c.term }

// YinYang returns the 阴/阳 遁.
func (c *Chart) YinYang() almanac.YinYang { return c.yinYang }

// Ju returns the local 局 number (1..9).
func (c *Chart) Ju() uint8 { return c.ju }

// Yuan returns the 三元 segment.
func (c *Chart) Yuan() enum.Yuan { return c.yuan }

// XunShou returns the 旬首 stem.
func (c *Chart) XunShou() almanac.Stem { return c.xunShou }

// ZhiFu returns the 值符 duty record.
func (c *Chart) ZhiFu() Duty { return c.zhiFu }

// ZhiShi returns the 值使 duty record.
func (c *Chart) ZhiShi() DutyDoor { return c.zhiShi }

// KongWang returns the 旬空亡 branch pair.
func (c *Chart) KongWang() [2]almanac.Branch { return c.kongWang }

// ===================== palace access =====================

// Palace returns the palace at number n (1..9), or nil when out of range.
// Returned pointer aliases the chart's internal storage — do not mutate.
func (c *Chart) Palace(n uint8) *palace.Palace {
	if n < 1 || n > 9 {
		return nil
	}
	return &c.palaces[n-1]
}

// Palaces streams the 9 palaces in canonical (1..9) order.
//
// Yields (n, palace) pairs where n is the palace number (1-indexed).
// The yielded pointer aliases the chart's internal storage.
func (c *Chart) Palaces() iter.Seq2[uint8, *palace.Palace] {
	return func(yield func(uint8, *palace.Palace) bool) {
		for i := 0; i < 9; i++ {
			if !yield(uint8(i+1), &c.palaces[i]) {
				return
			}
		}
	}
}

// Grid returns the 3×3 display grid [巽离坤; 震中兑; 艮坎乾].
func (c *Chart) Grid() [3][3]*palace.Palace {
	var out [3][3]*palace.Palace
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			out[row][col] = &c.palaces[tables.Grid[row][col]-1]
		}
	}
	return out
}

// ===================== user-stem queries =====================

// StemPalace returns the palace where a heavenly stem appears in the 地盘.
//
//   - 甲 (idx 0) has no visible position; falls back to the 值符原宫.
//   - Other stems are looked up; center palace stems are mapped to 2 (坤).
func (c *Chart) StemPalace(stem almanac.Stem) uint8 {
	if stem.Index() == 0 {
		return c.zhiFu.OriginalPalace
	}
	for i := range c.palaces {
		p := &c.palaces[i]
		if p.EarthStem.Index() == stem.Index() {
			if p.Number == 5 {
				return 2
			}
			return p.Number
		}
	}
	return c.zhiFu.OriginalPalace
}

// SelfPalace returns the palace of the user (subject), keyed by 日干.
func (c *Chart) SelfPalace() uint8 { return c.StemPalace(c.pillars.Day.Stem()) }

// OpponentPalace returns the palace of the counterpart / topic, which
// classical qimen places at 值符落宫.
func (c *Chart) OpponentPalace() uint8 { return c.zhiFu.Palace }

// ===================== aggregated streams =====================

// Patterns streams every 格局 across the 9 palaces in palace order.
func (c *Chart) Patterns() iter.Seq[pattern.Pattern] {
	return func(yield func(pattern.Pattern) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			for _, pat := range p.Patterns {
				if !yield(pat) {
					return
				}
			}
		}
	}
}

// ShenSha streams every 神煞 across the 9 palaces in palace order.
func (c *Chart) ShenSha() iter.Seq[shensha.ShenSha] {
	return func(yield func(shensha.ShenSha) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			for _, ss := range p.ShenSha {
				if !yield(ss) {
					return
				}
			}
		}
	}
}

// EarthStems streams the 地盘 stems with palace numbers.
func (c *Chart) EarthStems() iter.Seq2[uint8, almanac.Stem] {
	return func(yield func(uint8, almanac.Stem) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			if !yield(p.Number, p.EarthStem) {
				return
			}
		}
	}
}

// HeavenStems streams the 天盘 stems with palace numbers.
func (c *Chart) HeavenStems() iter.Seq2[uint8, almanac.Stem] {
	return func(yield func(uint8, almanac.Stem) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			if !yield(p.Number, p.HeavenStem) {
				return
			}
		}
	}
}

// HiddenStems streams the 暗干 stems with palace numbers.
func (c *Chart) HiddenStems() iter.Seq2[uint8, almanac.Stem] {
	return func(yield func(uint8, almanac.Stem) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			if !yield(p.Number, p.HiddenStem) {
				return
			}
		}
	}
}

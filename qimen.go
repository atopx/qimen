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
	method enum.Method
	style  enum.Style

	// Pillars + solar context
	solar    almanac.SolarTime
	pillars  almanac.Pillars
	term     almanac.Term
	yinYang  almanac.YinYang
	ju       uint8
	yuan     enum.Yuan
	xunShou  almanac.Stem
	zhiFu    DutyStar
	zhiShi   DutyDoor
	kongWang [2]almanac.Branch
	lunarDay almanac.LunarDay // cached at build for O(1) LunarDay() access

	// Palaces (1..9; index = palace number - 1).
	// Stored by value to avoid 9 heap allocations per chart.
	palaces [9]palace.Palace
}

// config holds optional construction parameters set via Option.
type config struct {
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

// DutyStar is the 值符 entry: the 九星 currently acting as 值符, its
// originating palace, and the palace it has rotated into.
type DutyStar struct {
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

// defaultConfig returns the default construction options: 时家 / 转盘.
func defaultConfig() config {
	return config{method: enum.MethodTime, style: enum.StyleRotate}
}

// New builds a chart for the current instant (UTC+8 wall clock) with
// default options (时家 / 转盘). Construction with defaults cannot fail,
// so the most common code path needs no error handling. Use From /
// FromTime for non-default options or instants.
func New() *Chart {
	return build(almanac.Now(), defaultConfig())
}

// From builds a chart from a [almanac.SolarTime].
func From(t almanac.SolarTime, opts ...Option) (*Chart, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	return build(t, cfg), nil
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

// FromTime builds a chart from a standard library [time.Time], using the
// time's own wall clock — its Location is respected as-is, never
// converted. Classical qimen casts charts in the local civil time of the
// event; convert explicitly (t.In(loc)) before calling when the value's
// zone is not the intended one.
func FromTime(t time.Time, opts ...Option) (*Chart, error) {
	st, err := almanac.SolarTimeFromTime(t)
	if err != nil {
		return nil, err // already wraps ErrInvalidTime
	}
	return From(st, opts...)
}

// FromTimestamp builds a chart from a Unix-seconds timestamp, interpreted
// as a UTC+8 (China Standard Time) wall clock. For another zone use
// FromTime(time.Unix(unix, 0).In(loc)).
func FromTimestamp(unix int64, opts ...Option) (*Chart, error) {
	st, err := almanac.SolarTimeFromUnix(unix)
	if err != nil {
		return nil, err // already wraps ErrInvalidTime
	}
	return From(st, opts...)
}

// validateConfig checks Method + Style at the Chart entry points.
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

// build is the shared chart-construction routine. cfg must already be
// validated; build itself is total over its closed-domain inputs.
func build(t almanac.SolarTime, cfg config) *Chart {
	// Resolve the solar term once; pillars, 阴阳遁 and 局 all derive
	// from it without further term lookups.
	term := t.Term()
	pillars := almanac.PillarsAt(t, term)
	yinYang := term.YinYang()
	yuan := compute.Yuan(pillars.Day)
	ju := compute.Ju(term, yuan)

	// Build the 6 plates. Style was validated at entry; primitives are
	// total functions assuming StyleRotate semantics.
	earth := plate.BuildEarth(yinYang, ju)

	xunShou := compute.XunShou(pillars.Hour)
	hourStem := pillars.Hour.Stem()

	// 值符原宫: the earth plate always contains the 旬首 stem (one of
	// 戊..癸) outside the center; the fallback is defensive only.
	zhiFuOriginalPalace := uint8(2)
	if p, ok := plate.FindStem(&earth, xunShou, true); ok {
		zhiFuOriginalPalace = p
	}
	zhiFuPalace := plate.FindHourStem(&earth, hourStem, zhiFuOriginalPalace)

	heaven := plate.BuildHeaven(&earth, yinYang, zhiFuOriginalPalace, zhiFuPalace)
	stars := plate.BuildStar(zhiFuOriginalPalace, zhiFuPalace)
	doors, zhiShiPalace := plate.BuildDoor(yinYang, zhiFuOriginalPalace, pillars.Hour)
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

	kongWang := pillars.Hour.EmptyBranches()

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

	// Pattern detection (build input view, then bucket results per palace).
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
	pats := pattern.AppendAll(make([]pattern.Pattern, 0, 12), &pIn)
	for i, ps := range bucketByPalace(pats, patternPalace) {
		palaces[i].Patterns = ps
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
	// The 10 神煞 kinds yield at most 12 instances (each 地盘 stem is
	// unique, so stem-anchored kinds land in at most one palace each).
	sss := shensha.AppendAll(make([]shensha.ShenSha, 0, 12), &ssIn)
	for i, ss := range bucketByPalace(sss, shenshaPalace) {
		palaces[i].ShenSha = ss
	}

	return &Chart{
		method:   cfg.method,
		style:    cfg.style,
		solar:    t,
		pillars:  pillars,
		term:     term,
		yinYang:  yinYang,
		ju:       ju,
		yuan:     yuan,
		xunShou:  xunShou,
		zhiFu:    DutyStar{Star: zhiFuStar, OriginalPalace: zhiFuOriginalPalace, Palace: zhiFuPalace},
		zhiShi:   DutyDoor{Door: zhiShiDoor, OriginalPalace: zhiFuOriginalPalace, Palace: zhiShiPalace},
		kongWang: kongWang,
		lunarDay: t.LunarDay(),
		palaces:  palaces,
	}
}

func patternPalace(p pattern.Pattern) uint8 { return p.Palace }
func shenshaPalace(s shensha.ShenSha) uint8 { return s.Palace }

// bucketByPalace regroups detection results (each tagged with a palace
// number 1..9) into 9 per-palace sub-slices of the SAME backing array:
// items are stably sorted by palace in place (they arrive nearly sorted,
// so insertion sort is effectively linear), then sliced per palace with
// the capacity clamped so appending to one palace's slice cannot clobber
// a neighbour's. Zero additional allocations.
func bucketByPalace[T any](items []T, palaceOf func(T) uint8) [9][]T {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && palaceOf(items[j]) < palaceOf(items[j-1]); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
	var out [9][]T
	for start := 0; start < len(items); {
		p := palaceOf(items[start])
		end := start + 1
		for end < len(items) && palaceOf(items[end]) == p {
			end++
		}
		out[p-1] = items[start:end:end]
		start = end
	}
	return out
}

// ===================== context accessors =====================

// SolarTime returns the solar instant used to build the chart.
func (c *Chart) SolarTime() almanac.SolarTime { return c.solar }

// LunarDay returns the lunar date of the chart's solar instant
// (cached at build for O(1) access).
func (c *Chart) LunarDay() almanac.LunarDay { return c.lunarDay }

// Method returns the chart's 起局法门 (currently always MethodTime).
func (c *Chart) Method() enum.Method { return c.method }

// Style returns the chart's 盘式 (currently always StyleRotate).
func (c *Chart) Style() enum.Style { return c.style }

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
func (c *Chart) ZhiFu() DutyStar { return c.zhiFu }

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
	if stem == almanac.Jia {
		return c.zhiFu.OriginalPalace
	}
	for i := range c.palaces {
		p := &c.palaces[i]
		if p.EarthStem == stem {
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

package qimen

import (
	"github.com/6tail/tyme4go/tyme"
)

// Qimen 奇门遁甲盘。一次起局后所有盘面信息的不可变快照。
type Qimen struct {
	solarTime tyme.SolarTime
	options   QimenOptions
	year      tyme.SixtyCycle
	month     tyme.SixtyCycle
	day       tyme.SixtyCycle
	hour      tyme.SixtyCycle
	term      tyme.SolarTerm
	yinYang   tyme.YinYang
	ju        uint8
	yuan      QimenYuan
	xunShou   tyme.HeavenStem
	zhiFu     QimenDutyStar
	zhiShi    QimenDutyDoor
	kongWang  [2]tyme.EarthBranch
	palaces   [9]*QimenPalace
}

// FromSolarTime 由阳历时间使用默认参数 (时家三元) 起局。默认参数下不会失败。
func FromSolarTime(t tyme.SolarTime) *Qimen {
	q, err := FromSolarTimeWithOptions(t, DefaultOptions())
	if err != nil {
		panic("default options must succeed: " + err.Error())
	}
	return q
}

// FromSolarTimeWithOptions 由阳历时间与自定义参数起局。
//
// 错误情形:
//   - ErrCodeUnsupportedMethod — 当前仅支持 [QimenMethodTime]
//   - ErrCodeUnsupportedChartType — 当前仅支持 [QimenChartTypeSanYuan]
//   - ErrCodeUnsupportedTerm — 节气索引越界
func FromSolarTimeWithOptions(t tyme.SolarTime, opts QimenOptions) (*Qimen, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	sch := t.GetSixtyCycleHour()
	year := sch.GetYear()
	month := sch.GetMonth()
	day := sch.GetDay()
	hour := sch.GetSixtyCycle()
	term := t.GetTerm()
	yinYang := computeYinYang(t)
	yuan := computeYuan(day)
	ju, err := computeJu(term, yuan)
	if err != nil {
		return nil, err
	}

	earth := buildEarthPlate(yinYang, ju)
	xunShou := computeXunShou(hour)
	hourStem := hour.GetHeavenStem()

	zhiFuOriginalPalace := uint8(2)
	if p := findStemPalace(earth, xunShou, true); p != nil {
		zhiFuOriginalPalace = *p
	}
	zhiFuPalace := findHourStemPalace(earth, hourStem, zhiFuOriginalPalace)

	heaven := buildHeavenPlate(earth, yinYang, zhiFuOriginalPalace, zhiFuPalace)
	stars := buildStarPlate(zhiFuOriginalPalace, zhiFuPalace)
	doors := buildDoorPlate(yinYang, zhiFuOriginalPalace, hour)
	gods := buildGodPlate(yinYang, zhiFuPalace)
	hidden := buildHiddenPlate(yinYang)

	zhiFuStar := QimenStarQinRui
	if s := stars.Get(zhiFuPalace); s != nil {
		zhiFuStar = *s
	}
	zhiShiDoor := QimenDoorDeath
	if d := QimenDoorFromPalace(zhiFuOriginalPalace); d != nil {
		zhiShiDoor = *d
	}
	zhiShiPalace := zhiFuPalace
	if p := findDoorPalace(doors, zhiShiDoor); p != nil {
		zhiShiPalace = *p
	}

	kongWang := computeKongWang(hour)

	// 构造盘面骨架
	var palaces [9]*QimenPalace
	for i := 0; i < 9; i++ {
		n := uint8(i + 1)
		var earthStem tyme.HeavenStem
		if v := earth.Get(n); v != nil {
			earthStem = *v
		} else {
			earthStem = tyme.HeavenStem{}.FromIndex(0)
		}
		heavenStem := earthStem
		if v := heaven.Get(n); v != nil {
			heavenStem = *v
		}
		hiddenStem := earthStem
		if v := hidden.Get(n); v != nil {
			hiddenStem = *v
		}
		palaces[i] = &QimenPalace{
			Number:           n,
			PalaceName:       PalaceNames[n],
			Direction:        tyme.Direction{}.FromIndex(int(n) - 1),
			EarthBranches:    branchesForPalace(n),
			EarthHeavenStem:  earthStem,
			SanQiLiuYi:       earthStem,
			HeavenHeavenStem: heavenStem,
			HiddenHeavenStem: hiddenStem,
			Star:             stars.Get(n),
			Door:             doors.Get(n),
			God:              gods.Get(n),
		}
	}

	// 衍生属性 (十神 / 长生 / 64 卦)
	dayStem := day.GetHeavenStem()
	for _, p := range palaces {
		if p.Number != 5 {
			ts := dayStem.GetTenStar(p.EarthHeavenStem)
			p.TenStar = &ts
		}
		if len(p.EarthBranches) > 0 {
			tt := NewTerrain(dayStem.GetTerrain(p.EarthBranches[0]))
			p.TerrainValue = &tt
		}
		// 卦: 上 = 该宫门的本宫卦, 下 = 该宫宫位卦; 中宫或无门时为 nil
		lower := TrigramFromPalace(p.Number)
		if lower != nil && p.Door != nil {
			upper := TrigramFromPalace(p.Door.HomePalace())
			if upper != nil {
				h := NewHexagram(*upper, *lower)
				p.Hexagram = &h
			}
		}
	}

	// 格局检测 + 分发
	patternList := detectPatterns(zhiFuOriginalPalace, zhiFuPalace, palaces, kongWang)
	for _, pat := range patternList {
		n := pat.Palace
		if n >= 1 && n <= 9 {
			palaces[n-1].Patterns = append(palaces[n-1].Patterns, pat)
		}
	}

	// 神煞检测 + 分发
	shenShaList := detectShenSha(
		year.GetHeavenStem(),
		month.GetEarthBranch(),
		dayStem,
		day.GetEarthBranch(),
		earth,
	)
	for _, ss := range shenShaList {
		n := ss.PalaceCell
		if n >= 1 && n <= 9 {
			palaces[n-1].ShenSha = append(palaces[n-1].ShenSha, ss)
		}
	}

	return &Qimen{
		solarTime: t,
		options:   opts,
		year:      year,
		month:     month,
		day:       day,
		hour:      hour,
		term:      term,
		yinYang:   yinYang,
		ju:        ju,
		yuan:      yuan,
		xunShou:   xunShou,
		zhiFu:     QimenDutyStar{Star: zhiFuStar, OriginalPalace: zhiFuOriginalPalace, Palace: zhiFuPalace},
		zhiShi:    QimenDutyDoor{Door: zhiShiDoor, OriginalPalace: zhiFuOriginalPalace, Palace: zhiShiPalace},
		kongWang:  kongWang,
		palaces:   palaces,
	}, nil
}

func validateOptions(opts QimenOptions) error {
	if opts.Method != QimenMethodTime {
		return newUnsupportedMethod(opts.Method)
	}
	if opts.ChartType != QimenChartTypeSanYuan {
		return newUnsupportedChartType(opts.ChartType)
	}
	return nil
}

// ===================== 公共访问器 =====================

// SolarTime 起局所用阳历时间。
func (q *Qimen) SolarTime() tyme.SolarTime { return q.solarTime }

// Options 起局参数。
func (q *Qimen) Options() QimenOptions { return q.options }

// Year 年柱。
func (q *Qimen) Year() tyme.SixtyCycle { return q.year }

// Month 月柱。
func (q *Qimen) Month() tyme.SixtyCycle { return q.month }

// Day 日柱。
func (q *Qimen) Day() tyme.SixtyCycle { return q.day }

// Hour 时柱。
func (q *Qimen) Hour() tyme.SixtyCycle { return q.hour }

// Term 节气。
func (q *Qimen) Term() tyme.SolarTerm { return q.term }

// YinYang 阴阳遁。
func (q *Qimen) YinYang() tyme.YinYang { return q.yinYang }

// Ju 局数 (1..=9)。
func (q *Qimen) Ju() uint8 { return q.ju }

// Yuan 三元 (上/中/下元)。
func (q *Qimen) Yuan() QimenYuan { return q.yuan }

// XunShou 旬首 (戊/己/庚/辛/壬/癸 之一)。
func (q *Qimen) XunShou() tyme.HeavenStem { return q.xunShou }

// ZhiFu 值符。
func (q *Qimen) ZhiFu() QimenDutyStar { return q.zhiFu }

// ZhiShi 值使。
func (q *Qimen) ZhiShi() QimenDutyDoor { return q.zhiShi }

// KongWang 旬空亡两支地支。
func (q *Qimen) KongWang() [2]tyme.EarthBranch { return q.kongWang }

// Palaces 九宫数据数组 (索引 0..=8 对应宫位 1..=9)。
func (q *Qimen) Palaces() [9]*QimenPalace { return q.palaces }

// Palace 通过宫位号 O(1) 取宫。越界返回 nil。
func (q *Qimen) Palace(number uint8) *QimenPalace {
	if number < 1 || number > 9 {
		return nil
	}
	return q.palaces[number-1]
}

// GridLayout 三行三列九宫展示 ([巽离坤; 震中兑; 艮坎乾]), 返回引用网格。
func (q *Qimen) GridLayout() [3][3]*QimenPalace {
	var out [3][3]*QimenPalace
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			out[row][col] = q.palaces[Grid[row][col]-1]
		}
	}
	return out
}

// StemPalace 任意天干在地盘所临之宫 (奇门"用神"基础查询)。
//
//   - 甲 (idx 0): 甲遁不可见, 以值符原宫 (旬首所临之宫) 代之
//   - 其他天干: 在 9 个 palace 的地盘干中线性查找
//   - 中宫干: 寄 2 宫 (坤宫)
func (q *Qimen) StemPalace(stem tyme.HeavenStem) uint8 {
	idx := stem.GetIndex()
	if idx == 0 {
		return q.zhiFu.OriginalPalace
	}
	for _, p := range q.palaces {
		if p.EarthHeavenStem.GetIndex() == idx {
			if p.Number == 5 {
				return 2
			}
			return p.Number
		}
	}
	return q.zhiFu.OriginalPalace
}

// SelfPalace 日干用神 (自身位/主位) 所临之宫。
//
// 占测自己时, 以日柱天干为用神, 该天干在地盘所临之宫即"我"。
func (q *Qimen) SelfPalace() uint8 { return q.StemPalace(q.day.GetHeavenStem()) }

// OpponentPalace 时干用神 (彼位/客位/事位) 所临之宫。
//
// 占测对方/事物时, 以时柱天干为用神, 与值符落宫同宫。
func (q *Qimen) OpponentPalace() uint8 { return q.zhiFu.Palace }

// SanQiLiuYi 三奇六仪 (按宫位顺序枚举)。
func (q *Qimen) SanQiLiuYi() []QimenHeavenStemPlacement {
	out := make([]QimenHeavenStemPlacement, 0, 9)
	for _, p := range q.palaces {
		out = append(out, QimenHeavenStemPlacement{Palace: p.Number, HeavenStem: p.SanQiLiuYi})
	}
	return out
}

// TianPan 天盘干 (按宫位顺序枚举)。
func (q *Qimen) TianPan() []QimenHeavenStemPlacement {
	out := make([]QimenHeavenStemPlacement, 0, 9)
	for _, p := range q.palaces {
		out = append(out, QimenHeavenStemPlacement{Palace: p.Number, HeavenStem: p.HeavenHeavenStem})
	}
	return out
}

// HiddenHeavenStems 暗干 (按宫位顺序枚举)。
func (q *Qimen) HiddenHeavenStems() []QimenHeavenStemPlacement {
	out := make([]QimenHeavenStemPlacement, 0, 9)
	for _, p := range q.palaces {
		out = append(out, QimenHeavenStemPlacement{Palace: p.Number, HeavenStem: p.HiddenHeavenStem})
	}
	return out
}

// Patterns 全盘所有格局 (聚合 9 宫的 Patterns)。
func (q *Qimen) Patterns() []Pattern {
	var out []Pattern
	for _, p := range q.palaces {
		out = append(out, p.Patterns...)
	}
	return out
}

// ShenSha 全盘所有神煞 (聚合 9 宫的 ShenSha)。
func (q *Qimen) ShenSha() []ShenSha {
	var out []ShenSha
	for _, p := range q.palaces {
		out = append(out, p.ShenSha...)
	}
	return out
}

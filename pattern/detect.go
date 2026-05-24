package pattern

import (
	"iter"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/element"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/internal/stemconst"
	"github.com/atopx/qimen/plate"
)

// DetectInput is the unpacked view of a chart used by Detect. Using
// flat arrays here avoids a circular dependency between pattern and
// palace packages.
//
// Arrays are indexed by (palace - 1); center (palace 5) is included
// but the detection rules skip it.
type DetectInput struct {
	ZhiFuOriginalPalace uint8
	ZhiFuPalace         uint8

	EarthStems  [9]almanac.Stem
	HeavenStems [9]almanac.Stem

	Doors    [9]enum.Door
	DoorsSet [9]bool

	Gods    [9]enum.God
	GodsSet [9]bool

	Branches [9][]almanac.Branch

	KongWang [2]almanac.Branch
}

// Detect yields every 格局 found in the chart view. Use slices.Collect
// to materialize if needed.
func Detect(in DetectInput) iter.Seq[Pattern] {
	return func(yield func(Pattern) bool) {
		// Global patterns
		if in.ZhiFuOriginalPalace == in.ZhiFuPalace {
			if !yield(Pattern{Kind: FuYin, Palace: in.ZhiFuPalace}) {
				return
			}
		} else if plate.AreOpposite(in.ZhiFuOriginalPalace, in.ZhiFuPalace) {
			if !yield(Pattern{
				Kind:           FanYin,
				Palace:         in.ZhiFuPalace,
				OriginalPalace: in.ZhiFuOriginalPalace,
			}) {
				return
			}
		}

		// Per-palace patterns
		for i := 0; i < 9; i++ {
			n := uint8(i + 1)
			if n == 5 {
				continue
			}
			if !detectPalace(n, &in, yield) {
				return
			}
		}
	}
}

func detectPalace(n uint8, in *DetectInput, yield func(Pattern) bool) bool {
	idx := n - 1
	heaven := in.HeavenStems[idx]
	earth := in.EarthStems[idx]
	h := heaven.Index()
	e := earth.Index()
	door := in.Doors[idx]
	doorOK := in.DoorsSet[idx]
	god := in.Gods[idx]
	godOK := in.GodsSet[idx]

	// 入墓
	if plate.IsStemInTomb(heaven, n) {
		if !yield(Pattern{Kind: RuMu, Palace: n, Stem: heaven}) {
			return false
		}
	}

	// 落空亡 — both kongwang branches may match this palace independently
	for _, kb := range in.KongWang {
		for _, pb := range in.Branches[idx] {
			if pb.Index() == kb.Index() {
				if !yield(Pattern{Kind: KongWang, Palace: n, Branch: kb}) {
					return false
				}
				break
			}
		}
	}

	// 门迫: door element 受克 by palace element (palace.element clobbers door.element)
	if doorOK && element.OfDoor(int(door)).RelationTo(element.FromPalace(n)) == element.Restrained {
		if !yield(Pattern{Kind: MenPo, Palace: n, Door: door}) {
			return false
		}
	}

	// 三奇得使
	if doorOK {
		if h == stemconst.Yi && door == enum.DoorOpen {
			if !yield(Pattern{Kind: YiQiDeShi, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Bing && door == enum.DoorRest {
			if !yield(Pattern{Kind: BingQiDeShi, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Ding && door == enum.DoorLife {
			if !yield(Pattern{Kind: DingQiDeShi, Palace: n}) {
				return false
			}
		}
	}

	// 八遁
	if doorOK {
		if h == stemconst.Bing && e == stemconst.Ding && door == enum.DoorLife {
			if !yield(Pattern{Kind: TianDun, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Yi && door == enum.DoorOpen && godOK && god == enum.GodJiuDi {
			if !yield(Pattern{Kind: DiDun, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Ding && door == enum.DoorRest && godOK && god == enum.GodTaiYin {
			if !yield(Pattern{Kind: RenDun, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Bing && door == enum.DoorLife && godOK && god == enum.GodJiuTian {
			if !yield(Pattern{Kind: ShenDun, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Ding && door == enum.DoorBlock && godOK && god == enum.GodJiuDi {
			if !yield(Pattern{Kind: GuiDun, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Yi && (door == enum.DoorOpen || door == enum.DoorBlock) && n == 4 {
			if !yield(Pattern{Kind: FengDun, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Yi && door == enum.DoorOpen && n == 6 {
			if !yield(Pattern{Kind: YunDun, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Yi && door == enum.DoorRest && n == 1 {
			if !yield(Pattern{Kind: LongDun, Palace: n}) {
				return false
			}
		}
		if h == stemconst.Yi && door == enum.DoorOpen && n == 7 {
			if !yield(Pattern{Kind: HuDun, Palace: n}) {
				return false
			}
		}
	}

	// 青龙返首 / 飞鸟跌穴
	if h == stemconst.Wu && e == stemconst.Bing {
		if !yield(Pattern{Kind: QingLongFanShou, Palace: n}) {
			return false
		}
	}
	if h == stemconst.Bing && e == stemconst.Wu {
		if !yield(Pattern{Kind: FeiNiaoDieXue, Palace: n}) {
			return false
		}
	}

	// 凶格
	if h == stemconst.Geng && e == stemconst.Gui {
		if !yield(Pattern{Kind: DaGe, Palace: n}) {
			return false
		}
	}
	if h == stemconst.Geng && e == stemconst.Ren {
		if !yield(Pattern{Kind: XiaoGe, Palace: n}) {
			return false
		}
	}
	if h == stemconst.Geng && e == stemconst.Ji {
		if !yield(Pattern{Kind: XingGe, Palace: n}) {
			return false
		}
	}

	// 悖格
	if (h == stemconst.Bing && e == stemconst.Geng) || (h == stemconst.Geng && e == stemconst.Bing) {
		if !yield(Pattern{Kind: BoGe, Palace: n}) {
			return false
		}
	}

	// 天网四张
	if h == stemconst.Gui && e == stemconst.Gui {
		if !yield(Pattern{Kind: TianWangSiZhang, Palace: n}) {
			return false
		}
	}

	return true
}

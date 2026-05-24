package qimen

import (
	"github.com/6tail/tyme4go/tyme"
)

// 天干索引 (按 tyme4go 顺序: 0甲 1乙 2丙 3丁 4戊 5己 6庚 7辛 8壬 9癸)
const (
	stemYi   = 1
	stemBing = 2
	stemDing = 3
	stemWu   = 4
	stemJi   = 5
	stemGeng = 6
	stemRen  = 8
	stemGui  = 9
)

// PatternKind 奇门格局种类。
type PatternKind int

const (
	// PatternFanYin 反吟格 — 值符落入与原宫对冲的宫位。
	PatternFanYin PatternKind = iota
	// PatternFuYin 伏吟格 — 值符落入原宫, 与原宫位重合。
	PatternFuYin
	// PatternRuMu 入墓格 — 天盘天干恰好落在自己的墓库宫。
	PatternRuMu
	// PatternKongWang 空亡 — 旬空亡支所在的宫位。
	PatternKongWang
	// PatternMenPo 门迫 — 门被所在宫位五行所克。
	PatternMenPo
	// PatternYiQiDeShi 乙奇得使 — 天盘乙奇加临地盘开门。
	PatternYiQiDeShi
	// PatternBingQiDeShi 丙奇得使 — 天盘丙奇加临地盘休门。
	PatternBingQiDeShi
	// PatternDingQiDeShi 丁奇得使 — 天盘丁奇加临地盘生门。
	PatternDingQiDeShi
	// PatternTianDun 天遁 — 天盘丙、地盘丁、加临生门。
	PatternTianDun
	// PatternDiDun 地遁 — 天盘乙奇、加临开门、临神盘九地。
	PatternDiDun
	// PatternRenDun 人遁 — 天盘丁奇、加临休门、临神盘太阴。
	PatternRenDun
	// PatternShenDun 神遁 — 天盘丙奇、加临生门、临神盘九天。
	PatternShenDun
	// PatternGuiDun 鬼遁 — 天盘丁奇、加临杜门、临神盘九地。
	PatternGuiDun
	// PatternFengDun 风遁 — 天盘乙奇、加临开门(或杜门)、落巽 4 宫。
	PatternFengDun
	// PatternYunDun 云遁 — 天盘乙奇、加临开门、落乾 6 宫。
	PatternYunDun
	// PatternLongDun 龙遁 — 天盘乙奇、加临休门、落坎 1 宫。
	PatternLongDun
	// PatternHuDun 虎遁 — 天盘乙奇、加临开门、落兑 7 宫。
	PatternHuDun
	// PatternQingLongFanShou 青龙返首 — 天盘戊加地盘丙。
	PatternQingLongFanShou
	// PatternFeiNiaoDieXue 飞鸟跌穴 — 天盘丙加地盘戊。
	PatternFeiNiaoDieXue
	// PatternDaGe 大格 — 天盘庚加地盘癸 (大凶)。
	PatternDaGe
	// PatternXiaoGe 小格 — 天盘庚加地盘壬。
	PatternXiaoGe
	// PatternXingGe 刑格 — 天盘庚加地盘己。
	PatternXingGe
	// PatternBoGe 悖格 — 天盘丙地盘庚 或 天盘庚地盘丙 (大凶)。
	PatternBoGe
	// PatternTianWangSiZhang 天网四张 — 天盘癸加地盘癸 (大凶)。
	PatternTianWangSiZhang
)

// Pattern 奇门格局实例。Kind 决定语义, 可选字段按 Kind 填充。
type Pattern struct {
	Kind           PatternKind
	Palace         uint8
	OriginalPalace uint8             // 仅 FanYin 使用
	Stem           *tyme.HeavenStem  // 仅 RuMu 使用
	Door           *QimenDoor        // 仅 MenPo 使用
	Branch         *tyme.EarthBranch // 仅 KongWang 使用
}

// Name 格局中文名。
func (p Pattern) Name() string {
	switch p.Kind {
	case PatternFanYin:
		return "反吟"
	case PatternFuYin:
		return "伏吟"
	case PatternRuMu:
		return "入墓"
	case PatternKongWang:
		return "落空亡"
	case PatternMenPo:
		return "门迫"
	case PatternYiQiDeShi:
		return "乙奇得使"
	case PatternBingQiDeShi:
		return "丙奇得使"
	case PatternDingQiDeShi:
		return "丁奇得使"
	case PatternTianDun:
		return "天遁"
	case PatternDiDun:
		return "地遁"
	case PatternRenDun:
		return "人遁"
	case PatternShenDun:
		return "神遁"
	case PatternGuiDun:
		return "鬼遁"
	case PatternFengDun:
		return "风遁"
	case PatternYunDun:
		return "云遁"
	case PatternLongDun:
		return "龙遁"
	case PatternHuDun:
		return "虎遁"
	case PatternQingLongFanShou:
		return "青龙返首"
	case PatternFeiNiaoDieXue:
		return "飞鸟跌穴"
	case PatternDaGe:
		return "大格"
	case PatternXiaoGe:
		return "小格"
	case PatternXingGe:
		return "刑格"
	case PatternBoGe:
		return "悖格"
	case PatternTianWangSiZhang:
		return "天网四张"
	}
	return ""
}

// Summary 格局的客观描述。
func (p Pattern) Summary() string {
	switch p.Kind {
	case PatternFanYin:
		return "奇门凶格。星门反吟,反复无常,事情易变。"
	case PatternFuYin:
		return "奇门凶格。星门伏吟,事情停滞,宜守不宜进。"
	case PatternRuMu:
		return "奇门凶格。日干或时干入墓,艰难阻塞,难以发展。"
	case PatternKongWang:
		return "奇门凶格。落入空亡宫,心愿落空,难以实现。"
	case PatternMenPo:
		return "奇门凶格。门克宫位,门被迫害,谋事不成,阻碍重重。"
	case PatternYiQiDeShi:
		return "奇门吉格。乙奇临开门,利于谋划,贵人相助,诸事吉利。"
	case PatternBingQiDeShi:
		return "奇门吉格。丙奇临休门,光明正大,官司必胜,声名可得。"
	case PatternDingQiDeShi:
		return "奇门吉格。丁奇临生门,才思敏捷,求财必得,生意兴隆。"
	case PatternTianDun:
		return "奇门吉格。丙丁同临生门,天助之,万事亨通,大吉大利。"
	case PatternDiDun:
		return "奇门吉格。乙奇临开门加九地,隐匿藏形,避凶趋吉。"
	case PatternRenDun:
		return "奇门吉格。丁奇临休门加太阴,人和之象,贵人暗助。"
	case PatternShenDun:
		return "奇门吉格。丙奇临生门加九天,神助之象,心想事成。"
	case PatternGuiDun:
		return "奇门吉格。丁奇临杜门加九地,神秘莫测,暗中成事。"
	case PatternFengDun:
		return "奇门吉格。乙奇临杜门在巽宫,运筹帷幄,避开祸端。"
	case PatternYunDun:
		return "奇门吉格。乙奇临开门在乾宫,腾云驾雾,步步高升。"
	case PatternLongDun:
		return "奇门吉格。乙奇临休门在坎宫,龙入大海,鸿图大展。"
	case PatternHuDun:
		return "奇门吉格。乙奇临开门在兑宫,猛虎添翼,势不可挡。"
	case PatternQingLongFanShou:
		return "奇门大吉格。天盘戊临地盘丙,大吉大利,名利双收。"
	case PatternFeiNiaoDieXue:
		return "奇门大吉格。天盘丙临地盘戊,诸事顺遂,不求自得。"
	case PatternDaGe:
		return "奇门凶格。庚临癸上,谋事难成,处处受制,大凶。"
	case PatternXiaoGe:
		return "奇门凶格。庚临壬上,小有阻碍,谋事迟缓。"
	case PatternXingGe:
		return "奇门凶格。庚临己上,官司牢狱,纷争不断。"
	case PatternBoGe:
		return "奇门凶格。庚金克制三奇,悖逆阻碍,诸事不顺,主凶。"
	case PatternTianWangSiZhang:
		return "奇门凶格。癸水入火域,身陷天网,行动招祸,主凶。"
	}
	return ""
}

// Auspice 传统文献既定的吉凶等级。
func (p Pattern) Auspice() Auspice {
	switch p.Kind {
	case PatternQingLongFanShou, PatternFeiNiaoDieXue:
		return AuspiceGreatAuspicious
	case PatternYiQiDeShi, PatternBingQiDeShi, PatternDingQiDeShi,
		PatternTianDun, PatternDiDun, PatternRenDun, PatternShenDun,
		PatternGuiDun, PatternFengDun, PatternYunDun, PatternLongDun, PatternHuDun:
		return AuspiceAuspicious
	case PatternDaGe, PatternBoGe, PatternTianWangSiZhang:
		return AuspiceGreatInauspicious
	case PatternFanYin, PatternFuYin, PatternRuMu, PatternKongWang,
		PatternMenPo, PatternXiaoGe, PatternXingGe:
		return AuspiceInauspicious
	}
	return AuspiceNeutral
}

// String 实现 fmt.Stringer。
func (p Pattern) String() string { return p.Name() }

// detectPatterns 检测全部当前已实现的格局。
//
// 先扫全局格 (反吟/伏吟), 再逐宫扫描 (跳过中宫) 22 个本宫格。
func detectPatterns(zhiFuOriginalPalace, zhiFuPalace uint8, palaces [9]*QimenPalace, kongWang [2]tyme.EarthBranch) []Pattern {
	var out []Pattern

	if zhiFuOriginalPalace == zhiFuPalace {
		out = append(out, Pattern{Kind: PatternFuYin, Palace: zhiFuPalace})
	} else if areOppositePalaces(zhiFuOriginalPalace, zhiFuPalace) {
		out = append(out, Pattern{Kind: PatternFanYin, OriginalPalace: zhiFuOriginalPalace, Palace: zhiFuPalace})
	}

	for _, palace := range palaces {
		if palace == nil || palace.Number == 5 {
			continue
		}
		detectPalacePatterns(palace, kongWang, &out)
	}

	return out
}

func detectPalacePatterns(palace *QimenPalace, kongWang [2]tyme.EarthBranch, out *[]Pattern) {
	n := palace.Number
	heaven := palace.HeavenHeavenStem
	earth := palace.EarthHeavenStem
	h := heaven.GetIndex()
	e := earth.GetIndex()
	door := palace.Door
	god := palace.God

	// 入墓
	if isStemInTomb(heaven, n) {
		stem := heaven
		*out = append(*out, Pattern{Kind: PatternRuMu, Palace: n, Stem: &stem})
	}

	// 落空亡 — 两支若同落本宫则各 push 一次
	for _, branch := range kongWang {
		for _, pb := range palace.EarthBranches {
			if pb.GetIndex() == branch.GetIndex() {
				b := branch
				*out = append(*out, Pattern{Kind: PatternKongWang, Palace: n, Branch: &b})
				break
			}
		}
	}

	// 门迫
	if door != nil && door.Element().RelationTo(ElementFromPalace(n)) == ElementRelationRestrained {
		d := *door
		*out = append(*out, Pattern{Kind: PatternMenPo, Palace: n, Door: &d})
	}

	// 三奇得使
	if door != nil {
		if h == stemYi && *door == QimenDoorOpen {
			*out = append(*out, Pattern{Kind: PatternYiQiDeShi, Palace: n})
		}
		if h == stemBing && *door == QimenDoorRest {
			*out = append(*out, Pattern{Kind: PatternBingQiDeShi, Palace: n})
		}
		if h == stemDing && *door == QimenDoorLife {
			*out = append(*out, Pattern{Kind: PatternDingQiDeShi, Palace: n})
		}
	}

	// 八遁
	if door != nil {
		if h == stemBing && e == stemDing && *door == QimenDoorLife {
			*out = append(*out, Pattern{Kind: PatternTianDun, Palace: n})
		}
		if h == stemYi && *door == QimenDoorOpen && god != nil && *god == QimenGodJiuDi {
			*out = append(*out, Pattern{Kind: PatternDiDun, Palace: n})
		}
		if h == stemDing && *door == QimenDoorRest && god != nil && *god == QimenGodTaiYin {
			*out = append(*out, Pattern{Kind: PatternRenDun, Palace: n})
		}
		if h == stemBing && *door == QimenDoorLife && god != nil && *god == QimenGodJiuTian {
			*out = append(*out, Pattern{Kind: PatternShenDun, Palace: n})
		}
		if h == stemDing && *door == QimenDoorBlock && god != nil && *god == QimenGodJiuDi {
			*out = append(*out, Pattern{Kind: PatternGuiDun, Palace: n})
		}
		if h == stemYi && (*door == QimenDoorOpen || *door == QimenDoorBlock) && n == 4 {
			*out = append(*out, Pattern{Kind: PatternFengDun, Palace: n})
		}
		if h == stemYi && *door == QimenDoorOpen && n == 6 {
			*out = append(*out, Pattern{Kind: PatternYunDun, Palace: n})
		}
		if h == stemYi && *door == QimenDoorRest && n == 1 {
			*out = append(*out, Pattern{Kind: PatternLongDun, Palace: n})
		}
		if h == stemYi && *door == QimenDoorOpen && n == 7 {
			*out = append(*out, Pattern{Kind: PatternHuDun, Palace: n})
		}
	}

	// 大吉格
	if h == stemWu && e == stemBing {
		*out = append(*out, Pattern{Kind: PatternQingLongFanShou, Palace: n})
	}
	if h == stemBing && e == stemWu {
		*out = append(*out, Pattern{Kind: PatternFeiNiaoDieXue, Palace: n})
	}

	// 凶格
	if h == stemGeng && e == stemGui {
		*out = append(*out, Pattern{Kind: PatternDaGe, Palace: n})
	}
	if h == stemGeng && e == stemRen {
		*out = append(*out, Pattern{Kind: PatternXiaoGe, Palace: n})
	}
	if h == stemGeng && e == stemJi {
		*out = append(*out, Pattern{Kind: PatternXingGe, Palace: n})
	}

	// 悖格 (丙↔庚)
	if (h == stemBing && e == stemGeng) || (h == stemGeng && e == stemBing) {
		*out = append(*out, Pattern{Kind: PatternBoGe, Palace: n})
	}

	// 天网四张
	if h == stemGui && e == stemGui {
		*out = append(*out, Pattern{Kind: PatternTianWangSiZhang, Palace: n})
	}
}

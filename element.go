package qimen

// Element 五行。
type Element int

const (
	// ElementWood 木。
	ElementWood Element = iota
	// ElementFire 火。
	ElementFire
	// ElementEarth 土。
	ElementEarth
	// ElementMetal 金。
	ElementMetal
	// ElementWater 水。
	ElementWater
)

// Name 中文名称 (木/火/土/金/水)。
func (e Element) Name() string {
	switch e {
	case ElementWood:
		return "木"
	case ElementFire:
		return "火"
	case ElementEarth:
		return "土"
	case ElementMetal:
		return "金"
	case ElementWater:
		return "水"
	}
	return ""
}

// String 实现 fmt.Stringer。
func (e Element) String() string { return e.Name() }

// ElementFromHeavenStemIndex 由天干索引 (0..=9: 甲乙丙丁戊己庚辛壬癸) 取五行。
//
// 甲乙=木, 丙丁=火, 戊己=土, 庚辛=金, 壬癸=水。
func ElementFromHeavenStemIndex(index int) Element {
	switch index {
	case 0, 1:
		return ElementWood
	case 2, 3:
		return ElementFire
	case 4, 5:
		return ElementEarth
	case 6, 7:
		return ElementMetal
	default:
		return ElementWater
	}
}

// ElementFromEarthBranchIndex 由地支索引 (0..=11: 子丑寅卯辰巳午未申酉戌亥) 取五行。
func ElementFromEarthBranchIndex(index int) Element {
	switch index {
	case 0, 11:
		return ElementWater
	case 1, 4, 7, 10:
		return ElementEarth
	case 2, 3:
		return ElementWood
	case 5, 6:
		return ElementFire
	default:
		return ElementMetal
	}
}

// ElementFromPalace 由九宫宫位号 (1..=9) 取所属五行。
//
// 1坎=水 2坤=土 3震=木 4巽=木 5中=土 6乾=金 7兑=金 8艮=土 9离=火。
// 越界时返回 [ElementEarth] 作为安全默认。
func ElementFromPalace(palace uint8) Element {
	switch palace {
	case 1:
		return ElementWater
	case 3, 4:
		return ElementWood
	case 6, 7:
		return ElementMetal
	case 9:
		return ElementFire
	default:
		return ElementEarth
	}
}

// ElementRelation 两个五行之间的关系 (以视角主体为 self)。
type ElementRelation int

const (
	// ElementRelationSame 比和 (双方五行相同)。
	ElementRelationSame ElementRelation = iota
	// ElementRelationGenerates self 生 other (耗泄 self)。
	ElementRelationGenerates
	// ElementRelationGenerated other 生 self (生扶 self)。
	ElementRelationGenerated
	// ElementRelationRestrains self 克 other (攻击 other)。
	ElementRelationRestrains
	// ElementRelationRestrained other 克 self (压制 self)。
	ElementRelationRestrained
)

// Name 简短中文描述 (比和/生出/受生/克出/受克)。
func (r ElementRelation) Name() string {
	switch r {
	case ElementRelationSame:
		return "比和"
	case ElementRelationGenerates:
		return "生出"
	case ElementRelationGenerated:
		return "受生"
	case ElementRelationRestrains:
		return "克出"
	case ElementRelationRestrained:
		return "受克"
	}
	return ""
}

// String 实现 fmt.Stringer。
func (r ElementRelation) String() string { return r.Name() }

// AuspiceAsSelf 以 self 为用神视角的吉凶等级。
//
//   - 比和 → 中和
//   - 受生 (other 生 self) → 吉, 主体得力
//   - 生出 (self 生 other, 耗泄) → 凶, 主体失力
//   - 克出 (self 克 other) → 中和
//   - 受克 (other 克 self) → 凶
func (r ElementRelation) AuspiceAsSelf() Auspice {
	switch r {
	case ElementRelationSame, ElementRelationRestrains:
		return AuspiceNeutral
	case ElementRelationGenerated:
		return AuspiceAuspicious
	case ElementRelationGenerates, ElementRelationRestrained:
		return AuspiceInauspicious
	}
	return AuspiceNeutral
}

// RelationTo 计算 e 相对于 other 的五行关系。
func (e Element) RelationTo(other Element) ElementRelation {
	if e == other {
		return ElementRelationSame
	}
	// 五行相生: 木→火→土→金→水→木
	switch {
	case e == ElementWood && other == ElementFire,
		e == ElementFire && other == ElementEarth,
		e == ElementEarth && other == ElementMetal,
		e == ElementMetal && other == ElementWater,
		e == ElementWater && other == ElementWood:
		return ElementRelationGenerates
	case e == ElementFire && other == ElementWood,
		e == ElementEarth && other == ElementFire,
		e == ElementMetal && other == ElementEarth,
		e == ElementWater && other == ElementMetal,
		e == ElementWood && other == ElementWater:
		return ElementRelationGenerated
	}
	// 五行相克: 木→土→水→火→金→木
	switch {
	case e == ElementWood && other == ElementEarth,
		e == ElementEarth && other == ElementWater,
		e == ElementWater && other == ElementFire,
		e == ElementFire && other == ElementMetal,
		e == ElementMetal && other == ElementWood:
		return ElementRelationRestrains
	}
	return ElementRelationRestrained
}

// Element 九星五行属性。
//
// 天蓬=水 天芮=土 天冲=木 天辅=木 天禽=土 天心=金 天柱=金 天任=土 天英=火 禽芮=土。
func (s QimenStar) Element() Element {
	switch s {
	case QimenStarTianPeng:
		return ElementWater
	case QimenStarTianChong, QimenStarTianFu:
		return ElementWood
	case QimenStarTianYing:
		return ElementFire
	case QimenStarTianXin, QimenStarTianZhu:
		return ElementMetal
	case QimenStarTianRui, QimenStarTianQin, QimenStarTianRen, QimenStarQinRui:
		return ElementEarth
	}
	return ElementEarth
}

// Element 八门五行属性。
//
// 休=水 生=土 伤=木 杜=木 景=火 死=土 惊=金 开=金。
func (d QimenDoor) Element() Element {
	switch d {
	case QimenDoorRest:
		return ElementWater
	case QimenDoorLife, QimenDoorDeath:
		return ElementEarth
	case QimenDoorHurt, QimenDoorBlock:
		return ElementWood
	case QimenDoorView:
		return ElementFire
	case QimenDoorFear, QimenDoorOpen:
		return ElementMetal
	}
	return ElementEarth
}

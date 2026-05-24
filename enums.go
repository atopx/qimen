package qimen

// QimenMethod 奇门起局方法。当前实现仅支持 [QimenMethodTime] (时家奇门)。
type QimenMethod int

const (
	// QimenMethodTime 时家。
	QimenMethodTime QimenMethod = iota
	// QimenMethodDay 日家 (尚未实现)。
	QimenMethodDay
	// QimenMethodMonth 月家 (尚未实现)。
	QimenMethodMonth
	// QimenMethodYear 年家 (尚未实现)。
	QimenMethodYear
)

// Name 中文名称 (时家/日家/月家/年家)。
func (m QimenMethod) Name() string {
	switch m {
	case QimenMethodTime:
		return "时家"
	case QimenMethodDay:
		return "日家"
	case QimenMethodMonth:
		return "月家"
	case QimenMethodYear:
		return "年家"
	}
	return ""
}

// String 实现 fmt.Stringer。
func (m QimenMethod) String() string { return m.Name() }

// QimenChartType 奇门盘式。当前实现仅支持 [QimenChartTypeSanYuan] (三元盘)。
type QimenChartType int

const (
	// QimenChartTypeSanYuan 三元。
	QimenChartTypeSanYuan QimenChartType = iota
	// QimenChartTypeSiZhu 四柱 (尚未实现)。
	QimenChartTypeSiZhu
)

// Name 中文名称 (三元/四柱)。
func (c QimenChartType) Name() string {
	switch c {
	case QimenChartTypeSanYuan:
		return "三元"
	case QimenChartTypeSiZhu:
		return "四柱"
	}
	return ""
}

// String 实现 fmt.Stringer。
func (c QimenChartType) String() string { return c.Name() }

// QimenYuan 三元 (上元/中元/下元)。
type QimenYuan int

const (
	// QimenYuanUpper 上元。
	QimenYuanUpper QimenYuan = iota
	// QimenYuanMiddle 中元。
	QimenYuanMiddle
	// QimenYuanLower 下元。
	QimenYuanLower
)

// Name 中文名称 (上元/中元/下元)。
func (y QimenYuan) Name() string {
	switch y {
	case QimenYuanUpper:
		return "上元"
	case QimenYuanMiddle:
		return "中元"
	case QimenYuanLower:
		return "下元"
	}
	return ""
}

// String 实现 fmt.Stringer。
func (y QimenYuan) String() string { return y.Name() }

// QimenStar 九星 (天蓬/天芮/天冲/天辅/天禽/天心/天柱/天任/天英)。
//
// 额外的 [QimenStarQinRui] (禽芮) 是当天禽与天芮同落 2 宫时使用的合并标记。
type QimenStar int

const (
	// QimenStarTianPeng 天蓬。
	QimenStarTianPeng QimenStar = iota
	// QimenStarTianRui 天芮。
	QimenStarTianRui
	// QimenStarTianChong 天冲。
	QimenStarTianChong
	// QimenStarTianFu 天辅。
	QimenStarTianFu
	// QimenStarTianQin 天禽。
	QimenStarTianQin
	// QimenStarTianXin 天心。
	QimenStarTianXin
	// QimenStarTianZhu 天柱。
	QimenStarTianZhu
	// QimenStarTianRen 天任。
	QimenStarTianRen
	// QimenStarTianYing 天英。
	QimenStarTianYing
	// QimenStarQinRui 禽芮 (天禽与天芮合并落 2 宫的标记)。
	QimenStarQinRui
)

// Name 九星中文名。
func (s QimenStar) Name() string {
	switch s {
	case QimenStarTianPeng:
		return "天蓬"
	case QimenStarTianRui:
		return "天芮"
	case QimenStarTianChong:
		return "天冲"
	case QimenStarTianFu:
		return "天辅"
	case QimenStarTianQin:
		return "天禽"
	case QimenStarTianXin:
		return "天心"
	case QimenStarTianZhu:
		return "天柱"
	case QimenStarTianRen:
		return "天任"
	case QimenStarTianYing:
		return "天英"
	case QimenStarQinRui:
		return "禽芮"
	}
	return ""
}

// String 实现 fmt.Stringer。
func (s QimenStar) String() string { return s.Name() }

// HomePalace 九星本位宫 (天蓬→1, 天芮→2, ..., 禽芮合并寄 2)。
func (s QimenStar) HomePalace() uint8 {
	switch s {
	case QimenStarTianPeng:
		return 1
	case QimenStarTianRui, QimenStarQinRui:
		return 2
	case QimenStarTianChong:
		return 3
	case QimenStarTianFu:
		return 4
	case QimenStarTianQin:
		return 5
	case QimenStarTianXin:
		return 6
	case QimenStarTianZhu:
		return 7
	case QimenStarTianRen:
		return 8
	case QimenStarTianYing:
		return 9
	}
	return 0
}

// QimenStarFromPalace 由九宫宫位号取该宫本位星 (越界返回 nil)。
func QimenStarFromPalace(palace uint8) *QimenStar {
	var s QimenStar
	switch palace {
	case 1:
		s = QimenStarTianPeng
	case 2:
		s = QimenStarTianRui
	case 3:
		s = QimenStarTianChong
	case 4:
		s = QimenStarTianFu
	case 5:
		s = QimenStarTianQin
	case 6:
		s = QimenStarTianXin
	case 7:
		s = QimenStarTianZhu
	case 8:
		s = QimenStarTianRen
	case 9:
		s = QimenStarTianYing
	default:
		return nil
	}
	return &s
}

// QimenDoor 八门 (休/生/伤/杜/景/死/惊/开)。
type QimenDoor int

const (
	// QimenDoorRest 休门。
	QimenDoorRest QimenDoor = iota
	// QimenDoorLife 生门。
	QimenDoorLife
	// QimenDoorHurt 伤门。
	QimenDoorHurt
	// QimenDoorBlock 杜门。
	QimenDoorBlock
	// QimenDoorView 景门。
	QimenDoorView
	// QimenDoorDeath 死门。
	QimenDoorDeath
	// QimenDoorFear 惊门。
	QimenDoorFear
	// QimenDoorOpen 开门。
	QimenDoorOpen
)

// Name 八门中文名。
func (d QimenDoor) Name() string {
	switch d {
	case QimenDoorRest:
		return "休门"
	case QimenDoorLife:
		return "生门"
	case QimenDoorHurt:
		return "伤门"
	case QimenDoorBlock:
		return "杜门"
	case QimenDoorView:
		return "景门"
	case QimenDoorDeath:
		return "死门"
	case QimenDoorFear:
		return "惊门"
	case QimenDoorOpen:
		return "开门"
	}
	return ""
}

// String 实现 fmt.Stringer。
func (d QimenDoor) String() string { return d.Name() }

// HomePalace 八门本位宫 (休→1, 生→8, 伤→3, 杜→4, 景→9, 死→2, 惊→7, 开→6)。
func (d QimenDoor) HomePalace() uint8 {
	switch d {
	case QimenDoorRest:
		return 1
	case QimenDoorDeath:
		return 2
	case QimenDoorHurt:
		return 3
	case QimenDoorBlock:
		return 4
	case QimenDoorView:
		return 9
	case QimenDoorOpen:
		return 6
	case QimenDoorFear:
		return 7
	case QimenDoorLife:
		return 8
	}
	return 0
}

// QimenDoorFromPalace 由九宫宫位号取该宫本位门 (中宫 5 无门返回 nil)。
func QimenDoorFromPalace(palace uint8) *QimenDoor {
	var d QimenDoor
	switch palace {
	case 1:
		d = QimenDoorRest
	case 2:
		d = QimenDoorDeath
	case 3:
		d = QimenDoorHurt
	case 4:
		d = QimenDoorBlock
	case 6:
		d = QimenDoorOpen
	case 7:
		d = QimenDoorFear
	case 8:
		d = QimenDoorLife
	case 9:
		d = QimenDoorView
	default:
		return nil
	}
	return &d
}

// QimenGod 九神 (值符/腾蛇/太阴/六合/白虎/玄武/九地/九天)。
type QimenGod int

const (
	// QimenGodZhiFu 值符。
	QimenGodZhiFu QimenGod = iota
	// QimenGodTengShe 腾蛇。
	QimenGodTengShe
	// QimenGodTaiYin 太阴。
	QimenGodTaiYin
	// QimenGodLiuHe 六合。
	QimenGodLiuHe
	// QimenGodBaiHu 白虎。
	QimenGodBaiHu
	// QimenGodXuanWu 玄武。
	QimenGodXuanWu
	// QimenGodJiuDi 九地。
	QimenGodJiuDi
	// QimenGodJiuTian 九天。
	QimenGodJiuTian
)

// Name 九神中文名。
func (g QimenGod) Name() string {
	switch g {
	case QimenGodZhiFu:
		return "值符"
	case QimenGodTengShe:
		return "腾蛇"
	case QimenGodTaiYin:
		return "太阴"
	case QimenGodLiuHe:
		return "六合"
	case QimenGodBaiHu:
		return "白虎"
	case QimenGodXuanWu:
		return "玄武"
	case QimenGodJiuDi:
		return "九地"
	case QimenGodJiuTian:
		return "九天"
	}
	return ""
}

// String 实现 fmt.Stringer。
func (g QimenGod) String() string { return g.Name() }

package qimen

import "fmt"

// Trigram 八卦。索引顺序: 0乾, 1兑, 2离, 3震, 4巽, 5坎, 6艮, 7坤。
type Trigram int

const (
	// TrigramQian 乾 ☰ 天。
	TrigramQian Trigram = iota
	// TrigramDui 兑 ☱ 泽。
	TrigramDui
	// TrigramLi 离 ☲ 火。
	TrigramLi
	// TrigramZhen 震 ☳ 雷。
	TrigramZhen
	// TrigramXun 巽 ☴ 风。
	TrigramXun
	// TrigramKan 坎 ☵ 水。
	TrigramKan
	// TrigramGen 艮 ☶ 山。
	TrigramGen
	// TrigramKun 坤 ☷ 地。
	TrigramKun
)

// TrigramFromPalace 由后天九宫宫位号取八卦 (中宫 5 返回 nil)。
//
// 1坎/2坤/3震/4巽/6乾/7兑/8艮/9离。
func TrigramFromPalace(palace uint8) *Trigram {
	var t Trigram
	switch palace {
	case 1:
		t = TrigramKan
	case 2:
		t = TrigramKun
	case 3:
		t = TrigramZhen
	case 4:
		t = TrigramXun
	case 6:
		t = TrigramQian
	case 7:
		t = TrigramDui
	case 8:
		t = TrigramGen
	case 9:
		t = TrigramLi
	default:
		return nil
	}
	return &t
}

// Name 中文卦名 (乾/兑/...)。
func (t Trigram) Name() string {
	switch t {
	case TrigramQian:
		return "乾"
	case TrigramDui:
		return "兑"
	case TrigramLi:
		return "离"
	case TrigramZhen:
		return "震"
	case TrigramXun:
		return "巽"
	case TrigramKan:
		return "坎"
	case TrigramGen:
		return "艮"
	case TrigramKun:
		return "坤"
	}
	return ""
}

// Symbol 八卦符号 (☰☱☲☳☴☵☶☷)。
func (t Trigram) Symbol() string {
	switch t {
	case TrigramQian:
		return "☰"
	case TrigramDui:
		return "☱"
	case TrigramLi:
		return "☲"
	case TrigramZhen:
		return "☳"
	case TrigramXun:
		return "☴"
	case TrigramKan:
		return "☵"
	case TrigramGen:
		return "☶"
	case TrigramKun:
		return "☷"
	}
	return ""
}

// ElementName 八卦自然元素名 (天/泽/火/雷/风/水/山/地)。
func (t Trigram) ElementName() string {
	switch t {
	case TrigramQian:
		return "天"
	case TrigramDui:
		return "泽"
	case TrigramLi:
		return "火"
	case TrigramZhen:
		return "雷"
	case TrigramXun:
		return "风"
	case TrigramKan:
		return "水"
	case TrigramGen:
		return "山"
	case TrigramKun:
		return "地"
	}
	return ""
}

// String 实现 fmt.Stringer。
func (t Trigram) String() string { return t.Name() }

// Hexagram 六十四卦, 由上下两卦组成。
type Hexagram struct {
	upper Trigram
	lower Trigram
	index uint8
}

// NewHexagram 由上下两卦构造一个六十四卦。
func NewHexagram(upper, lower Trigram) Hexagram {
	idx := hexagramIndex[upper][lower]
	return Hexagram{upper: upper, lower: lower, index: idx}
}

// Upper 上卦。
func (h Hexagram) Upper() Trigram { return h.upper }

// Lower 下卦。
func (h Hexagram) Lower() Trigram { return h.lower }

// Index 卦序 (周易序卦传 0..64)。
func (h Hexagram) Index() uint8 { return h.index }

// Symbol Unicode 卦象符号 (䷀..䷿)。
func (h Hexagram) Symbol() string { return hexagramData[h.index].symbol }

// Name 中文卦名。
func (h Hexagram) Name() string { return hexagramData[h.index].name }

// Summary 客观描述。
func (h Hexagram) Summary() string { return hexagramData[h.index].summary }

// Auspice 传统易学既定的吉凶定性。
func (h Hexagram) Auspice() Auspice { return hexagramAuspice[h.index] }

// String 实现 fmt.Stringer (符号 + 卦名)。
func (h Hexagram) String() string { return fmt.Sprintf("%s %s", h.Symbol(), h.Name()) }

// hexagramIndex 是 [upper][lower] → 周易序卦传索引的反查表。
var hexagramIndex = [8][8]uint8{
	{0, 9, 12, 24, 43, 5, 32, 11},    // 乾
	{42, 57, 48, 16, 27, 46, 30, 44}, // 兑
	{13, 37, 29, 20, 49, 63, 55, 34}, // 离
	{33, 53, 54, 50, 31, 39, 61, 15}, // 震
	{8, 60, 36, 41, 56, 58, 52, 19},  // 巽
	{4, 59, 62, 2, 47, 28, 38, 7},    // 坎
	{25, 40, 21, 26, 17, 3, 51, 22},  // 艮
	{10, 18, 35, 23, 45, 6, 14, 1},   // 坤
}

type hexagramEntry struct {
	name    string
	symbol  string
	summary string
}

// hexagramData 卦序对应的 (中文名, Unicode 符号, 客观描述)。
var hexagramData = [64]hexagramEntry{
	{"乾为天", "䷀", "六十四卦之首。元亨利贞,刚健中正,自强不息。"},
	{"坤为地", "䷁", "六十四卦第二卦。厚德载物,柔顺利贞,包容万物。"},
	{"水雷屯", "䷂", "六十四卦之一。元亨利贞,初始艰难,宜守不宜进。"},
	{"山水蒙", "䷃", "六十四卦之一。亨通成功,童蒙求教,循循善诱。"},
	{"水天需", "䷄", "六十四卦之一。有孚,光亨,贞吉。待时而动,饮食宴乐。"},
	{"天水讼", "䷅", "六十四卦之一。有孚窒惕,中吉终凶。息事宁人,防范争端。"},
	{"地水师", "䷆", "六十四卦之一。师贞,丈人吉,无咎。纪律严明,统帅有道。"},
	{"水地比", "䷇", "六十四卦之一。吉。原筮,元永贞,无咎。团结亲比,辅佐共赢。"},
	{"风天小畜", "䷈", "六十四卦之一。亨,密云不雨,自我西郊。小有积蓄,力量尚微。"},
	{"天泽履", "䷉", "六十四卦之一。履虎尾,不咥人,亨。如履薄冰,循礼而行。"},
	{"地天泰", "䷊", "六十四卦之一。小往大来,吉亨,天地交泰。"},
	{"天地否", "䷋", "六十四卦之一。否之匪人,不利君子贞,闭塞不通。"},
	{"天火同人", "䷌", "六十四卦之一。同人于野,亨。志同道合,团结协作。"},
	{"火天大有", "䷍", "六十四卦之一。元亨。火在天上,普照万物。繁荣昌盛,大获成功。"},
	{"地山谦", "䷎", "六十四卦之一。亨,君子有终。卑以自牧,谦受益。"},
	{"雷地豫", "䷏", "六十四卦之一。利建侯行师。和乐喜悦,顺时而动。"},
	{"泽雷随", "䷐", "六十四卦之一。元亨利贞,无咎。随遇而安,顺应潮流。"},
	{"山风蛊", "䷑", "六十四卦之一。元亨,利涉大川。振疲起顿,改革弊端。"},
	{"地泽临", "䷒", "六十四卦之一。元亨利贞,至于八月有凶。督导推进,把握时机。"},
	{"风地观", "䷓", "六十四卦之一。盥而不荐,有孚颙若。观察感悟,风行地上。"},
	{"火雷噬嗑", "䷔", "六十四卦之一。亨。利用狱。障碍阻隔,果断铲除。"},
	{"山火贲", "䷕", "六十四卦之一。亨,小利有攸往。装饰与文饰,文明与礼仪。"},
	{"山地剥", "䷖", "六十四卦之一。不利有攸往。阴长阳消,剥落殆尽。"},
	{"地雷复", "䷗", "六十四卦之一。亨,出入无疾,朋来无咎。生机重现,剥极必复。"},
	{"天雷无妄", "䷘", "六十四卦之一。元亨利贞,匪正有眚。顺应自然,不妄求。"},
	{"山天大畜", "䷙", "六十四卦之一。利贞,不家食吉,利涉大川。积蓄硕大,大有作为。"},
	{"山雷颐", "䷚", "六十四卦之一。贞吉。观颐,自求口实。颐养天年,言行有节。"},
	{"泽风大过", "䷛", "六十四卦之一。栋桡。利有攸往,亨。重压之下,果断行动。"},
	{"坎为水", "䷜", "六十四卦之一。习坎,有孚维心,亨。险难重重,需内怀诚信。"},
	{"离为火", "䷝", "六十四卦之一。利贞,亨。畜牝牛吉。明亮美丽,依附守正。"},
	{"泽山咸", "䷞", "六十四卦之一。亨,利贞。取女吉。感应交融,心心相印。"},
	{"雷风恒", "䷟", "六十四卦之一。亨,无咎,利贞。持之以恒,恒久而不变。"},
	{"天山遁", "䷠", "六十四卦之一。遁亨,小利贞。退避守正,以退为进。"},
	{"雷天大壮", "䷡", "六十四卦之一。利贞。盛大强壮,需坚守正道防偏激。"},
	{"火地晋", "䷢", "六十四卦之一。康侯用锡马蕃庶,昼日三接。光明普照,晋升腾达。"},
	{"地火明夷", "䷣", "六十四卦之一。利艰贞。韬光养晦,在黑暗中保存光明。"},
	{"风火家人", "䷤", "六十四卦之一。利女贞。各尽其道,家和万事兴。"},
	{"火泽睽", "䷥", "六十四卦之一。小事吉。二女同居,其志不同。异中求同,化解对立。"},
	{"水山蹇", "䷦", "六十四卦之一。利西南,不利东北。利见大人,贞吉。进退维谷,宜反求诸己。"},
	{"雷水解", "䷧", "六十四卦之一。利西南。无所往,其来复吉。解脱困境,化解矛盾。"},
	{"山泽损", "䷨", "六十四卦之一。有孚,元吉,无咎,可贞。损下益上,诚信为本。"},
	{"风雷益", "䷩", "六十四卦之一。利有攸往,利涉大川。损上益下,获利匪浅。"},
	{"泽天夬", "䷪", "六十四卦之一。扬于王庭,孚号有厉。果断抉择,清除障碍。"},
	{"天风姤", "䷫", "六十四卦之一。天下有风,阴长阳消,防范微隐。"},
	{"泽地萃", "䷬", "六十四卦之一。亨。王假有庙。利见大人,亨。精英聚会,资源整合。"},
	{"地风升", "䷭", "六十四卦之一。元亨,用见大人,勿恤。步步高升,积微成大。"},
	{"泽水困", "䷮", "六十四卦之一。亨,贞,大人吉,无咎。穷困潦倒,言有不信。"},
	{"水风井", "䷯", "六十四卦之一。改邑不改井,无丧无得。源远流长,取之不尽。"},
	{"泽火革", "䷰", "六十四卦之一。己日乃孚。元亨利贞,悔亡。变革创新,面貌一新。"},
	{"火风鼎", "䷱", "六十四卦之一。元吉,亨。革故鼎新,稳重权柄。"},
	{"震为雷", "䷲", "六十四卦之一。震亨,震惊百里,不丧匕鬯。警示振作,恐惧修省。"},
	{"艮为山", "䷳", "六十四卦之一。艮其背,不获其身。动静不失其时,止其所当止。"},
	{"风山渐", "䷴", "六十四卦之一。女归吉,利贞。循序渐进,草木渐茂。"},
	{"雷泽归妹", "䷵", "六十四卦之一。征凶,无攸利。不当其位,需防偏离正道。"},
	{"雷火丰", "䷶", "六十四卦之一。亨,王假之。日中则昃,盛极之时当忧。"},
	{"火山旅", "䷷", "六十四卦之一。小亨,旅贞吉。身在异乡,寻求安身。"},
	{"巽为风", "䷸", "六十四卦之一。小亨,利攸往。随风潜入,谦逊顺服。"},
	{"兑为泽", "䷹", "六十四卦之一。亨,利贞。悦随其后,和悦交流。"},
	{"风水涣", "䷺", "六十四卦之一。亨,王假有庙。涣散离析,需重新整合。"},
	{"水泽节", "䷻", "六十四卦之一。亨,苦节不可贞。节制有度,不可过头。"},
	{"风泽中孚", "䷼", "六十四卦之一。豚鱼吉。中心诚信,感化万物。"},
	{"雷山小过", "䷽", "六十四卦之一。亨,利贞。可小事,不可大事。过犹不及,宜下不宜上。"},
	{"水火既济", "䷾", "六十四卦之一。亨,小利贞,初吉终乱。"},
	{"火水未济", "䷿", "六十四卦末卦。亨,但须谨慎,事尚未完成。"},
}

// hexagramAuspice 卦序对应的传统易学吉凶定性 (5 级)。
var hexagramAuspice = [64]Auspice{
	AuspiceGreatAuspicious,   // 0  乾为天
	AuspiceAuspicious,        // 1  坤为地
	AuspiceInauspicious,      // 2  水雷屯
	AuspiceNeutral,           // 3  山水蒙
	AuspiceAuspicious,        // 4  水天需
	AuspiceInauspicious,      // 5  天水讼
	AuspiceNeutral,           // 6  地水师
	AuspiceAuspicious,        // 7  水地比
	AuspiceNeutral,           // 8  风天小畜
	AuspiceAuspicious,        // 9  天泽履
	AuspiceGreatAuspicious,   // 10 地天泰
	AuspiceGreatInauspicious, // 11 天地否
	AuspiceAuspicious,        // 12 天火同人
	AuspiceGreatAuspicious,   // 13 火天大有
	AuspiceAuspicious,        // 14 地山谦
	AuspiceAuspicious,        // 15 雷地豫
	AuspiceAuspicious,        // 16 泽雷随
	AuspiceNeutral,           // 17 山风蛊
	AuspiceAuspicious,        // 18 地泽临
	AuspiceNeutral,           // 19 风地观
	AuspiceNeutral,           // 20 火雷噬嗑
	AuspiceNeutral,           // 21 山火贲
	AuspiceGreatInauspicious, // 22 山地剥
	AuspiceAuspicious,        // 23 地雷复
	AuspiceNeutral,           // 24 天雷无妄
	AuspiceAuspicious,        // 25 山天大畜
	AuspiceAuspicious,        // 26 山雷颐
	AuspiceInauspicious,      // 27 泽风大过
	AuspiceGreatInauspicious, // 28 坎为水
	AuspiceAuspicious,        // 29 离为火
	AuspiceAuspicious,        // 30 泽山咸
	AuspiceNeutral,           // 31 雷风恒
	AuspiceNeutral,           // 32 天山遁
	AuspiceAuspicious,        // 33 雷天大壮
	AuspiceAuspicious,        // 34 火地晋
	AuspiceInauspicious,      // 35 地火明夷
	AuspiceAuspicious,        // 36 风火家人
	AuspiceNeutral,           // 37 火泽睽
	AuspiceGreatInauspicious, // 38 水山蹇
	AuspiceAuspicious,        // 39 雷水解
	AuspiceNeutral,           // 40 山泽损
	AuspiceAuspicious,        // 41 风雷益
	AuspiceNeutral,           // 42 泽天夬
	AuspiceNeutral,           // 43 天风姤
	AuspiceAuspicious,        // 44 泽地萃
	AuspiceAuspicious,        // 45 地风升
	AuspiceGreatInauspicious, // 46 泽水困
	AuspiceNeutral,           // 47 水风井
	AuspiceNeutral,           // 48 泽火革
	AuspiceAuspicious,        // 49 火风鼎
	AuspiceNeutral,           // 50 震为雷
	AuspiceNeutral,           // 51 艮为山
	AuspiceAuspicious,        // 52 风山渐
	AuspiceInauspicious,      // 53 雷泽归妹
	AuspiceNeutral,           // 54 雷火丰
	AuspiceNeutral,           // 55 火山旅
	AuspiceNeutral,           // 56 巽为风
	AuspiceAuspicious,        // 57 兑为泽
	AuspiceNeutral,           // 58 风水涣
	AuspiceNeutral,           // 59 水泽节
	AuspiceAuspicious,        // 60 风泽中孚
	AuspiceNeutral,           // 61 雷山小过
	AuspiceNeutral,           // 62 水火既济
	AuspiceNeutral,           // 63 火水未济
}

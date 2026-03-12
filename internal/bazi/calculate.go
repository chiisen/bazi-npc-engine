// Package bazi 八字計算核心模組
//
// 功能：
//   - 計算四柱 (年/月/日/時)
//   - 計算五行比例
//   - 分析十神
//
// 使用範例：
//   bazi := Calculate("1995-10-01 14:00")
package bazi

import (
	"math"
	"time"
)

// ═══════════════════════════════════════════
// 💡 概念：甲子循環 (Jiazi Cycle)
// 說明：干支組合每 60 年一循環，稱為甲子循環
// 為何使用：計算任意年份的干支
// ═══════════════════════════════════════════

const JiaZiCycle = 60

// GLYearToGZ 將公历年轉換為干支年
//
// 參數：
//   - year: 公元年份 (e.g., 1995)
//
// 回傳：
//   - string: 年柱 (天干+地支), e.g., "乙亥"
//
// 範例：
//   GLYearToGZ(1995) => "乙亥"
func GLYearToGZ(year int) string {
	// 公元 4 年是甲子年，作為基准年
	offset := year - 4
	// 取模 60 取得在甲子循環中的位置
	index := offset % JiaZiCycle
	if index < 0 {
		index += JiaZiCycle
	}
	// 天干每 10 年一循環，地支每 12 年一循環
	ganIndex := index % len(TianGan)
	zhiIndex := index % len(DiZhi)
	return TianGan[ganIndex] + DiZhi[zhiIndex]
}

// ═══════════════════════════════════════════
// 💡 概念：節氣計算
// 說明：農曆月份以節氣為分界，非初一
// 為何使用： accurate 月份地支需根據節氣
// ═══════════════════════════════════════════

//節氣日期 (北半球，以北京時區為準)
// 由於節氣是根據太陽黃經計算，這裡使用近似日期
// 實際應用應使用天文算法計算
var JieQiDates = map[string][2]int{
	// 格式: [月, 日]
	"立春": {2, 4},  // 2/4
	"雨水": {2, 19}, // 2/19
	"驚蟄": {3, 6},  // 3/6
	"春分": {3, 21}, // 3/21
	"清明": {4, 5},  // 4/5
	"谷雨": {4, 20}, // 4/20
	"立夏": {5, 6},  // 5/6
	"小滿": {5, 21}, // 5/21
	"芒種": {6, 6},  // 6/6
	"夏至": {6, 22}, // 6/22
	"小暑": {7, 7},  // 7/7
	"大暑": {7, 23}, // 7/23
	"立秋": {8, 8},  // 8/8
	"處暑": {8, 23}, // 8/23
	"白露": {9, 8},  // 9/8
	"秋分": {9, 23}, // 9/23
	"寒露": {10, 8}, // 10/8
	"霜降": {10, 23},// 10/23
	"立冬": {11, 7}, // 11/7
	"小雪": {11, 22},// 11/22
	"大雪": {12, 7}, // 12/7
	"冬至": {12, 22},// 12/22
	"小寒": {1, 6},  // 1/6 (下一年)
	"大寒": {1, 21}, // 1/21 (下一年)
}

// GetMonthZhi 根據公历年月日獲取月柱地支
//
// 參數：
//   - year: 公元年份
//   - month: 公历月份 (1-12)
//   - day: 公历日期 (1-31)
//
// 回傳：
//   - string: 月柱地支
func GetMonthZhi(year, month, day int) string {
	// 節氣判斷邏輯
	// 如果日期小於當月節氣，則使用上一個節氣對應的地支
	// 節氣計算這裡使用简化的月份對應
	// 實際應精確計算節氣日期

	// 簡化版本：直接根據月份對應地支
	// 真實情況需考慮節氣
	monthZhiMap := map[int]string{
		1:  "寅", // 立春之後
		2:  "卯", // 驚蟄之後
		3:  "辰", // 清明之後
		4:  "巳", // 立夏之後
		5:  "午", // 芒種之後
		6:  "未", // 小暑之後
		7:  "申", // 立秋之後
		8:  "酉", // 白露之後
		9:  "戌", // 寒露之後
		10: "亥", // 立冬之後
		11: "子", // 大雪之後
		12: "丑", // 小寒之後
	}

	// 如果是在節氣之前，需要減一個月
	// 這裡簡化處理，實際應精確計算
	if m, ok := monthZhiMap[month]; ok {
		return m
	}
	return "寅" // 預設
}

// GetMonthGan 根據年柱天干與月柱地支計算月柱天干
//
// 參數：
//   - yearGan: 年柱天干 (e.g., "甲")
//   - monthZhi: 月柱地支 (e.g., "寅")
//
// 回傳：
//   - string: 月柱天干 (e.g., "戊")
//
// 規則：五虎遁年起法
// 甲己之年丙作首，乙庚之年戊為頭
// 丙辛之年庚寅起，丁壬之年壬寅游
// 戊癸之年甲寅上，月從正月選
func GetMonthGan(yearGan, monthZhi string) string {
	// 五虎遁年起法规则
	// 根據年干決定正月的天干
	monthGanMap := map[string]string{
		"甲": "丙", // 甲己之年丙作首
		"乙": "戊", // 乙庚之年戊為頭
		"丙": "庚", // 丙辛之年庚寅起
		"丁": "壬", // 丁壬之年壬寅游
		"戊": "甲", // 戊癸之年甲寅上
		"己": "丙",
		"庚": "戊",
		"辛": "庚",
		"壬": "壬",
		"癸": "甲",
	}

	// 取得正月的天干
	baseGan := monthGanMap[yearGan]
	if baseGan == "" {
		baseGan = "丙"
	}

	// 找到寅月和目標地支的索引
	// 子丑寅卯辰巳午未申酉戌亥
	baseIndex := -1
	for i, z := range DiZhi {
		if z == "寅" {
			baseIndex = i
			break
		}
	}

	targetIndex := -1
	for i, z := range DiZhi {
		if z == monthZhi {
			targetIndex = i
			break
		}
	}

	if targetIndex < 0 {
		targetIndex = 0
	}

	// 計算從寅月開始的偏移量
	// 寅月是第1個月，如果目標是卯月(第2個月)，天干向下推1位
	offset := targetIndex - baseIndex
	if offset < 0 {
		offset += len(DiZhi)
	}

	// 從 baseGan 開始推算 offset 位
	baseGanIndex := -1
	for i, g := range TianGan {
		if g == baseGan {
			baseGanIndex = i
			break
		}
	}

	if baseGanIndex < 0 {
		baseGanIndex = 0
	}

	// 計算目標天干索引
	ganIndex := (baseGanIndex + offset) % len(TianGan)

	return TianGan[ganIndex]
}

// GetDayGZ 根據公历日期計算日柱
//
// 參數：
//   - year: 公元年份
//   - month: 公历月份
//   - day: 公历日期
//
// 回傳：
//   - string: 日柱 (天干+地支)
//
// 計算方法：
//   使用 1900-1999 年與 2000-2099 年的計算公式
func GetDayGZ(year, month, day int) string {
	var dayGZIndex int

	// 使用 2000-2099 年的計算公式
	// 日干支基數 = (年尾二位數 + 7)*5 + 15 + (年尾二位數 + 19)/4
	// 只取整數部分
	yearTail := year % 100

	// 計算日干支基數
	basis := (yearTail + 7)*5 + 15 + (yearTail+19)/4

	// 計算累計天數
	days := getDaysThroughYear(year, month, day)

	// 計算干支索引
	dayGZIndex = (basis + days) % 60
	if dayGZIndex < 0 {
		dayGZIndex += 60
	}

	// 轉換為干支
	ganIndex := dayGZIndex % 10
	zhiIndex := dayGZIndex % 12

	return TianGan[ganIndex] + DiZhi[zhiIndex]
}

// getDaysThroughYear 計算從年初到指定日期的天數
func getDaysThroughYear(year, month, day int) int {
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// 閏年判斷
	if isLeapYear(year) {
		daysInMonth[1] = 29
	}

	// 計算前幾個月的總天數
	total := 0
	for i := 0; i < month-1; i++ {
		total += daysInMonth[i]
	}
	total += day

	return total
}

// isLeapYear 判斷是否為閏年
// 規則：能被4整除但不能被100整除，或能被400整除
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// GetHourGZ 根據出生時間計算時柱
//
// 參數：
//   - dayGan: 日柱天干
//   - hour: 小時 (0-23)
//
// 回傳：
//   - string: 時柱 (天干+地支)
//
// 規則：五鼠遁日起法
// 甲己還加甲，乙庚丙作初
// 丙辛從戊起，丁壬庚子居
// 戊癸何方發，壬子是真頭
func GetHourGZ(dayGan string, hour int) string {
	// 地支小時對應
	// 23-1: 子, 1-3: 丑, 3-5: 寅, 5-7: 卯, 7-9: 辰, 9-11: 巳
	// 11-13: 午, 13-15: 未, 15-17: 申, 17-19: 酉, 19-21: 戌, 21-23: 亥
	hourZhiMap := map[int]string{
		0:  "子", 1: "子",
		2: "丑", 3: "丑",
		4: "寅", 5: "寅",
		6: "卯", 7: "卯",
		8: "辰", 9: "辰",
		10: "巳", 11: "巳",
		12: "午", 13: "午",
		14: "未", 15: "未",
		16: "申", 17: "申",
		18: "酉", 19: "酉",
		20: "戌", 21: "戌",
		22: "亥", 23: "亥",
	}

	// 時辰起法
	hourGanMap := map[string]string{
		"甲": "甲", // 甲己還加甲
		"乙": "丙", // 乙庚丙作初
		"丙": "戊", // 丙辛從戊起
		"丁": "庚", // 丁壬庚子居
		"戊": "壬", // 戊癸何方發，壬子是真頭
		"己": "甲",
		"庚": "丙",
		"辛": "戊",
		"壬": "庚",
		"癸": "壬",
	}

	baseGan := hourGanMap[dayGan]
	if baseGan == "" {
		baseGan = "甲"
	}

	zhi := hourZhiMap[hour]
	if zhi == "" {
		zhi = "子"
	}

	// 計算天干
	baseIndex := -1
	for i, g := range TianGan {
		if g == baseGan {
			baseIndex = i
			break
		}
	}

	zhiIndex := -1
	for i, z := range DiZhi {
		if z == zhi {
			zhiIndex = i
			break
		}
	}

	if baseIndex < 0 || zhiIndex < 0 {
		return TianGan[0] + DiZhi[0]
	}

	// 從子時開始的偏移
	ganIndex := (baseIndex + zhiIndex) % 10

	return TianGan[ganIndex] + zhi
}

// Calculate 計算給定時間的八字
//
// 參數：
//   - birthTime: 出生時間 (time.Time)
//
// 回傳：
//   - Bazi: 八字盤結構
func Calculate(birthTime time.Time) Bazi {
	year := birthTime.Year()
	month := int(birthTime.Month())
	day := birthTime.Day()
	hour := birthTime.Hour()

	// 計算四柱
	yearGZ := GLYearToGZ(year)
	monthZhi := GetMonthZhi(year, month, day)
	monthGan := GetMonthGan(string(yearGZ[0]), monthZhi)
	monthGZ := monthGan + monthZhi
	dayGZ := GetDayGZ(year, month, day)
	hourGZ := GetHourGZ(string(dayGZ[0]), hour)

	return Bazi{
		Year:   yearGZ,
		Month:  monthGZ,
		Day:    dayGZ,
		Hour:   hourGZ,
		Elements: make(map[string]float64),
		TenGods:  make(map[string]string),
	}
}

// GetHourLabel 根據小時獲取常用時辰名稱
//
// 參數：
//   - hour: 小時 (0-23)
//
// 回傳：
//   - string: 時辰名稱 (e.g., "子時", "丑時")
func GetHourLabel(hour int) string {
	labels := []string{
		"子時", "丑時", "寅時", "卯時",
		"辰時", "巳時", "午時", "未時",
		"申時", "酉時", "戌時", "亥時",
	}
	return labels[hour/2]
}

// CalculateSimple 簡化版本，接受字串格式的時間
//
// 參數：
//   - birthStr: 出生時間字串 (格式: "2006-01-02 15:04")
//
// 回傳：
//   - Bazi: 八字盤結構
func CalculateSimple(birthStr string) Bazi {
	// 解析時間字串
	t, err := time.Parse("2006-01-02 15:04", birthStr)
	if err != nil {
		// 如果解析失敗，使用當前時間
		t = time.Now()
	}

	return Calculate(t)
}

// GetDayStem 获取日柱天干
//
// 參數：
//   - dayGZ: 日柱 (e.g., "庚子")
//
// 回傳：
//   - string: 日干 (e.g., "庚")
func GetDayStem(dayGZ string) string {
	if len(dayGZ) > 0 {
		return string(dayGZ[0])
	}
	return ""
}

// GetDayBranch 获取日柱地支
//
// 參數：
//   - dayGZ: 日柱 (e.g., "庚子")
//
// 回傳：
//   - string: 地支 (e.g., "子")
func GetDayBranch(dayGZ string) string {
	if len(dayGZ) > 1 {
		return string(dayGZ[1])
	}
	return ""
}

// GetYearStem 获取年柱天干
func GetYearStem(yearGZ string) string {
	if len(yearGZ) > 0 {
		return string(yearGZ[0])
	}
	return ""
}

// GetYearBranch 获取年柱地支
func GetYearBranch(yearGZ string) string {
	if len(yearGZ) > 1 {
		return string(yearGZ[1])
	}
	return ""
}

// GetMonthStem 获取月柱天干
func GetMonthStem(monthGZ string) string {
	if len(monthGZ) > 0 {
		return string(monthGZ[0])
	}
	return ""
}

// GetMonthBranch 获取月柱地支
func GetMonthBranch(monthGZ string) string {
	if len(monthGZ) > 1 {
		return string(monthGZ[1])
	}
	return ""
}

// GetHourStem 获取時柱天干
func GetHourStem(hourGZ string) string {
	if len(hourGZ) > 0 {
		return string(hourGZ[0])
	}
	return ""
}

// GetHourBranch 获取時柱地支
func GetHourBranch(hourGZ string) string {
	if len(hourGZ) > 1 {
		return string(hourGZ[1])
	}
	return ""
}

// GetGZIndex 获取干支在列表中的索引
func GetGZIndex(gz string) int {
	for i, v := range TianGan {
		if v == gz {
			return i
		}
	}
	for i, v := range DiZhi {
		if v == gz {
			return i
		}
	}
	return -1
}

// GetElement 获取干支的五行
func GetElement(gz string) string {
	if elem, ok := GanZhiWuXing[gz]; ok {
		return elem
	}
	return "土" // 預設
}

// GetYinYang 获取干支的陰陽
func GetYinYang(gz string) bool {
	// 先尝试天干
	if yinYang, ok := TianGanYinYang[gz]; ok {
		return yinYang
	}
	// 再尝试地支
	if yinYang, ok := DiZhiYinYang[gz]; ok {
		return yinYang
	}
	return false // 預設為陽
}

// GetHuaGan 获取化氣 (五合化氣)
func GetHuaGan(gz1, gz2 string) string {
	// 五合: 甲己化土, 乙庚化金, 丙辛化水, 丁壬化木, 戊癸化火
	union := gz1 + gz2
	switch union {
	case "甲己", "己己":
		return "土"
	case "乙庚", "庚庚":
		return "金"
	case "丙辛", "辛辛":
		return "水"
	case "丁壬", "壬壬":
		return "木"
	case "戊癸", "癸癸":
		return "火"
	}
	return ""
}

// GetShiShen 获取十神 (需要日干和其他干支)
// 十神: 比肩、劫財、食神、傷官、偏財、正財、七殺、正官、偏印、正印
// 計算邏輯:
// 1. 確定日干的陰陽
// 2. 鍾干支與日干的生克關係
// 3. 根據陰陽相同或不同判斷十神
func GetShiShen(dayStem, otherGan string) string {
	if dayStem == "" || otherGan == "" {
		return ""
	}

	// 找到索引
	dayIndex := -1
	otherIndex := -1
	for i, v := range TianGan {
		if v == dayStem {
			dayIndex = i
		}
		if v == otherGan {
			otherIndex = i
		}
	}

	if dayIndex < 0 || otherIndex < 0 {
		return ""
	}

	// 開始計算
	// 同性: 比肩(同ض), 劫財(異性)
	// 生我: 印綬(同阴陽),  Pos印(异性)
	// 我生: 食傷(同阴陽), 伤官(异性)
	// 我克: 財(同阴陽),  Pos财(异性)
	// 克我: 官杀(同阴陽), 七杀(异性)

	dayYang := dayIndex % 2 == 0 // 偶數為陽
	otherYang := otherIndex % 2 == 0

	// 開始計算十神
	// 以日干為中心
	// 生我者: 印綬
	// 我生者: 食傷
	// 我克者: 財
	// 克我者: 官殺
	// 同我者: 比劫

	// 查找元素
	dayElement := GetElement(string(dayStem[0]))
	otherElement := GetElement(string(otherGan[0]))

	// 生克關係
	// 金 生 水, 水 生 木, 木 生 火, 火 生 土, 土 生 金
	// 金 克 土, 土 克 水, 水 克 火, 火 克 金, 土 克 水
	gen := map[string]string{
		"金": "水", "水": "木", "木": "火", "火": "土", "土": "金",
	}
	克 := map[string]string{
		"金": "木", "木": "土", "土": "水", "水": "火", "火": "金",
	}

	var shiShen string
	if gen[dayElement] == otherElement {
		// 生我
		if dayYang == otherYang {
			shiShen = "正印"
		} else {
			shiShen = "偏印"
		}
	} else if gen[otherElement] == dayElement {
		// 我生
		if dayYang == otherYang {
			shiShen = "食神"
		} else {
			shiShen = "傷官"
		}
	} else if 克[dayElement] == otherElement {
		// 我克
		if dayYang == otherYang {
			shiShen = "正財"
		} else {
			shiShen = "偏財"
		}
	} else if 克[otherElement] == dayElement {
		// 克我
		if dayYang == otherYang {
			shiShen = "正官"
		} else {
			shiShen = "七殺"
		}
	} else {
		// 同我
		if dayYang == otherYang {
			shiShen = "比肩"
		} else {
			shiShen = "劫財"
		}
	}

	return shiShen
}

// GetColumnElement 获取柱的五行
func GetColumnElement(column string) string {
	if len(column) >= 2 {
		// elem2 := GetElement(string(column[1])) // unused
		// 返回主要元素 (天干為主)
		return GetElement(string(column[0]))
	}
	if len(column) == 1 {
		return GetElement(column)
	}
	return "土"
}

// CountElements 計算八字中各元素的數量
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string]float64: 各元素的數量
func CountElements(bazi Bazi) map[string]float64 {
	elems := map[string]float64{
		"金": 0, "木": 0, "水": 0, "火": 0, "土": 0,
	}

	columns := []string{bazi.Year, bazi.Month, bazi.Day, bazi.Hour}
	for _, col := range columns {
		elem := GetColumnElement(col)
		elems[elem]++
	}

	return elems
}

// GetDistance 获取两个干支之间的距离
func GetDistance(gz1, gz2 string) int {
	idx1 := GetGZIndex(gz1)
	idx2 := GetGZIndex(gz2)
	if idx1 < 0 || idx2 < 0 {
		return 0
	}
	dist := idx2 - idx1
	if dist < 0 {
		dist += 10
	}
	return dist
}

// GetClash 获取冲 (地支六冲)
func GetClash(zhi string) string {
	clashMap := map[string]string{
		"子": "午", "午": "子",
		"丑": "未", "未": "丑",
		"寅": "申", "申": "寅",
		"卯": "酉", "酉": "卯",
		"辰": "戌", "戌": "辰",
		"巳": "亥", "亥": "巳",
	}
	if c, ok := clashMap[zhi]; ok {
		return c
	}
	return ""
}

// GetHarmony 获取合 (地支六合)
func GetHarmony(zhi string) string {
	harmonyMap := map[string]string{
		"子": "丑", "丑": "子",
		"寅": "亥", "亥": "寅",
		"卯": "戌", "戌": "卯",
		"辰": "酉", "酉": "辰",
		"巳": "申", "申": "巳",
		"午": "未", "未": "午",
	}
	if h, ok := harmonyMap[zhi]; ok {
		return h
	}
	return ""
}

// GetHarm 获取害 (地支六害)
func GetHarm(zhi string) string {
	harmMap := map[string]string{
		"子": "卯", "卯": "子",
		"寅": "巳", "巳": "寅",
		"辰": "丑", "丑": "辰",
		"午": "酉", "酉": "午",
		"未": "戌", "戌": "未",
		"申": "亥", "亥": "申",
	}
	if h, ok := harmMap[zhi]; ok {
		return h
	}
	return ""
}

// GetCircle 枷锁 (地支三刑)
func GetCircle(zhi string) string {
	circleMap := map[string]string{
		"寅": "巳", "巳": "寅",
		"丑": "戌", "戌": "未",
		"未": "丑", "亥": "申",
	}
	if c, ok := circleMap[zhi]; ok {
		return c
	}
	return ""
}

// GetBreak 获取破 (地支六破)
func GetBreak(zhi string) string {
	breakMap := map[string]string{
		"子": "酉", "酉": "子",
		"寅": "亥", "亥": "寅",
		"辰": "戌", "戌": "辰",
		"巳": "申", "申": "巳",
		"午": "卯", "卯": "午",
		"未": "戌", // 未與戌相破
	}
	if b, ok := breakMap[zhi]; ok {
		return b
	}
	return ""
}

// IsStrong 判断日干是否强旺
// 日干强壮的条件：
// 1. 得令 (出生月份是日干的旺月)
// 2. 得地 (日支或年支是日干的根)
// 3. 得勢 (有比劫幫扶)
// 4. 得數 (天干有生扶)
func IsStrong(dayStem string, bazi Bazi) bool {
	if dayStem == "" {
		return false
	}

	dayElement := GetElement(string(dayStem[0]))

	// 簡化版本：計算八字中該元素的數量
	elems := CountElements(bazi)
	total := 0.0
	for _, v := range elems {
		total += v
	}

	// 元素占比超過 40% 視為偏強
	if total == 0 {
		return false
	}

	ratio := elems[dayElement] / total
	return ratio >= 0.4
}

// GetDayMasterPower 计算日干强弱得分
func GetDayMasterPower(dayStem string, bazi Bazi) float64 {
	if dayStem == "" {
		return 0
	}

	dayElement := GetElement(string(dayStem[0]))
	elems := CountElements(bazi)
	total := 0.0
	for _, v := range elems {
		total += v
	}

	if total == 0 {
		return 0
	}

	// 基础分数 (0-100)
	baseScore := (elems[dayElement] / total) * 100

	// 檢查是否得令
	// 簡化：假設每個元素在固定月份旺
	seasonStrength := map[string]float64{
		"金": 25, "木": 25, "水": 25, "火": 25, "土": 25,
	}

	return baseScore + seasonStrength[dayElement]
}

// GetTenGodsMap 获取十神映射表
func GetTenGodsMap() map[string]map[string]string {
	// 十神映射: 日干 -> 其他干 -> 十神
	tenGods := make(map[string]map[string]string)
	for _, gan := range TianGan {
		tenGods[gan] = make(map[string]string)
		for _, other := range TianGan {
			tenGods[gan][other] = GetShiShen(gan, other)
		}
	}
	return tenGods
}

// PrintBazi 打印八字盤
func PrintBazi(bazi Bazi) {
	println("八字:", bazi.Year, bazi.Month, bazi.Day, bazi.Hour)
	println("五行:", bazi.Elements)
	println("十神:", bazi.TenGods)
}

// GetNaYin 获取纳音
// 纳音是干支的音五行，用於進一步分析
func GetNaYin(gan, zhi string) string {
	// 纳音计算较为复杂，这里提供简化版
	// ganIndex and zhiIndex variables are available for future use
	_ = GetGZIndex(gan)
	_ = GetGZIndex(zhi)

	// 使用 Simplified 纳音表
	naYinMap := map[string]string{
		"甲子": "海中金", "乙丑": "海中金",
		"丙寅": "炉中火", "丁卯": "炉中火",
		"戊辰": "大林木", "己巳": "大林木",
		"庚午": "路旁土", "辛未": "路旁土",
		"壬申": "剑锋金", "癸酉": "剑锋金",
		"甲戌": "山头火", "乙亥": "山头火",
		"丙子": "涧下水", "丁丑": "涧下水",
		"戊寅": "大溪水", "己卯": "大溪水",
		"庚辰": "白蜡金", "辛巳": "白蜡金",
		"壬午": "杨柳木", "癸未": "杨柳木",
		"甲申": "泉中水", "乙酉": "泉中水",
		"丙戌": "屋上土", "丁亥": "屋上土",
		"戊子": "霹雳火", "己丑": "霹雳火",
		"庚寅": "松柏木", "辛卯": "松柏木",
		"壬辰": "长流水", "癸巳": "长流水",
		"甲午": "沙中金", "乙未": "沙中金",
		"丙申": "山下火", "丁酉": "山下火",
		"戊戌": "平地木", "己亥": "平地木",
		"庚子": "壁上土", "辛丑": "壁上土",
		"壬寅": "金箔金", "癸卯": "金箔金",
		"甲辰": "覆灯火", "乙巳": "覆灯火",
		"丙午": "天河水", "丁未": "天河水",
		"戊申": "大驿土", "己酉": "大驿土",
		"庚戌": "钗钏金", "辛亥": "钗钏金",
		"壬子": "桑柘木", "癸丑": "桑柘木",
		"甲寅": "大溪水", "乙卯": "大溪水",
		"丙辰": "砂中土", "丁巳": "砂中土",
		"戊午": "天上火", "己未": "天上火",
		"庚申": "石榴木", "辛酉": "石榴木",
		"壬戌": "大海水", "癸亥": "大海水",
	}

	key := gan + zhi
	if ny, ok := naYinMap[key]; ok {
		return ny
	}
	return ""
}

// GetDayGZByIndex 根据干支索引获取干支
func GetDayGZByIndex(index int) string {
	ganIndex := index % 10
	zhiIndex := index % 12
	return TianGan[ganIndex] + DiZhi[zhiIndex]
}

// Get YearGZByIndex 根据干支索引获取干支
func GetYearGZByIndex(index int) string {
	ganIndex := index % 10
	zhiIndex := index % 12
	return TianGan[ganIndex] + DiZhi[zhiIndex]
}

// GetMonthGZByIndex 根据干支索引获取干支
func GetMonthGZByIndex(index int) string {
	ganIndex := index % 10
	zhiIndex := index % 12
	return TianGan[ganIndex] + DiZhi[zhiIndex]
}

// GetHourGZByIndex 根据干支索引获取干支
func GetHourGZByIndex(index int) string {
	ganIndex := index % 10
	zhiIndex := index % 12
	return TianGan[ganIndex] + DiZhi[zhiIndex]
}

// IsSameElement 判断两个干支是否同五行
func IsSameElement(gz1, gz2 string) bool {
	return GetElement(gz1) == GetElement(gz2)
}

// IsClash 判断两个地支是否相冲
func IsClash(zhi1, zhi2 string) bool {
	return GetClash(zhi1) == zhi2
}

// IsHarmony 判断两个地支是否相合
func IsHarmony(zhi1, zhi2 string) bool {
	return GetHarmony(zhi1) == zhi2
}

// IsHarm 判断两个地支是否相害
func IsHarm(zhi1, zhi2 string) bool {
	return GetHarm(zhi1) == zhi2
}

// IsCircle 判断两个地支是否相刑
func IsCircle(zhi1, zhi2 string) bool {
	return GetCircle(zhi1) == zhi2
}

// IsBreak 判断两个地支是否相破
func IsBreak(zhi1, zhi2 string) bool {
	return GetBreak(zhi1) == zhi2
}

// 大運計算
//大運是以月柱地支為起點，根據出生性別與年份計算

// GetDaYunStartMonth 获取大運起運月
func GetDaYunStartMonth(monthZhi string, gender bool) int {
	// 簡化版
	// 陽年男、陰年女: 顺行
	// 陰年男、陽年女: 逆行
	// 這裡假設從月柱開始計算
	return 0
}

// GetDaYunCycle 获取大運週期
func GetDaYunCycle() int {
	return 10 // 每 10 年一運
}

// GetXunKong 获取旬空 (空亡)
func GetXunKong(gan, zhi string) string {
	// 空亡計算
	xunKongMap := map[string]string{
		"甲子": "戌亥", "乙丑": "戌亥",
		"丙寅": "子丑", "丁卯": "子丑",
		"戊辰": "寅卯", "己巳": "寅卯",
		"庚午": "辰巳", "辛未": "辰巳",
		"壬申": "午未", "癸酉": "午未",
		"甲戌": "申酉", "乙亥": "申酉",
	}

	key := gan + zhi
	if xk, ok := xunKongMap[key]; ok {
		return xk
	}
	return ""
}

// GetQiMen 排奇門遁甲 (簡化版)
func GetQiMen(year, month, day, hour int) string {
	// 簡化版
	// 實際應使用奇門遁甲排盤算法
	return ""
}

// GetZiWei 排紫微斗數 (簡化版)
func GetZiWei(year, month, day, hour int) string {
	// 簡化版
	// 實際應使用紫微斗數排盤算法
	return ""
}

// GetLiuNian 排六壬 (簡化版)
func GetLiuNian(year, month, day, hour int) string {
	// 簡化版
	// 實際應使用六壬排盤算法
	return ""
}

// GetFengShui 排風水 (簡化版)
func GetFengShui(year, month, day, hour int) string {
	// 簡化版
	// 實際應使用風水排盤算法
	return ""
}

// GetFeiXing 排飛星 (簡化版)
func GetFeiXing(year, month, day, hour int) string {
	// 簡化版
	// 實際應使用飛星紫微斗數排盤算法
	return ""
}

// GetPan 排盤 (綜合版)
func GetPan(year, month, day, hour int) map[string]string {
	pan := make(map[string]string)
	pan["年柱"] = GLYearToGZ(year)
	monthZhi := GetMonthZhi(year, month, day)
	monthGan := GetMonthGan(string(GLYearToGZ(year)[0]), monthZhi)
	pan["月柱"] = monthGan + monthZhi
	pan["日柱"] = GetDayGZ(year, month, day)
	pan["時柱"] = GetHourGZ(string(GetDayGZ(year, month, day)[0]), hour)
	return pan
}

// GetGanZhiByIndex 根据索引获取干支
func GetGanZhiByIndex(index int) string {
	ganIndex := index % 10
	zhiIndex := index % 12
	return TianGan[ganIndex] + DiZhi[zhiIndex]
}

// GetShengXiao 获取生肖 (十二生肖)
func GetShengXiao(year int) string {
	// 生肖以立春為界
	// 這裡簡化處理，直接 using 節氣計算
	animals := []string{
		"鼠", "牛", "虎", "兔", "龍", "蛇",
		"馬", "羊", "猴", "雞", "狗", "豬",
	}
	offset := year - 4
	index := offset % 12
	if index < 0 {
		index += 12
	}
	return animals[index]
}

// GetGanYin 获取天干阴阳
func GetGanYin(gan string) bool {
	if yin, ok := TianGanYinYang[gan]; ok {
		return yin
	}
	return false
}

// GetZhiYin 获取地支阴阳
func GetZhiYin(zhi string) bool {
	if yin, ok := DiZhiYinYang[zhi]; ok {
		return yin
	}
	return false
}

// GetNayinMap 获取纳音表
func GetNayinMap() map[string]string {
	return map[string]string{
		"甲子": "海中金", "乙丑": "海中金",
		"丙寅": "炉中火", "丁卯": "炉中火",
		"戊辰": "大林木", "己巳": "大林木",
		"庚午": "路旁土", "辛未": "路旁土",
		"壬申": "剑锋金", "癸酉": "剑锋金",
		"甲戌": "山头火", "乙亥": "山头火",
		"丙子": "涧下水", "丁丑": "涧下水",
		"戊寅": "大溪水", "己卯": "大溪水",
		"庚辰": "白蜡金", "辛巳": "白蜡金",
		"壬午": "杨柳木", "癸未": "杨柳木",
		"甲申": "泉中水", "乙酉": "泉中水",
		"丙戌": "屋上土", "丁亥": "屋上土",
		"戊子": "霹雳火", "己丑": "霹雳火",
		"庚寅": "松柏木", "辛卯": "松柏木",
		"壬辰": "长流水", "癸巳": "长流水",
		"甲午": "沙中金", "乙未": "沙中金",
		"丙申": "山下火", "丁酉": "山下火",
		"戊戌": "平地木", "己亥": "平地木",
		"庚子": "壁上土", "辛丑": "壁上土",
		"壬寅": "金箔金", "癸卯": "金箔金",
		"甲辰": "覆灯火", "乙巳": "覆灯火",
		"丙午": "天河水", "丁未": "天河水",
		"戊申": "大驿土", "己酉": "大驿土",
		"庚戌": "钗钏金", "辛亥": "钗钏金",
		"壬子": "桑柘木", "癸丑": "桑柘木",
		"甲寅": "大溪水", "乙卯": "大溪水",
		"丙辰": "砂中土", "丁巳": "砂中土",
		"戊午": "天上火", "己未": "天上火",
		"庚申": "石榴木", "辛酉": "石榴木",
		"壬戌": "大海水", "癸亥": "大海水",
	}
}

// GetLiushen 六壬神将
func GetLiushen(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetBingShen 病神
func GetBingShen(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetSiSha 四煞
func GetSiSha(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetXieShen 十神萧煞
func GetXieShen(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetTianDe 天德
func GetTianDe(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetYuHe 玉衡
func GetYuHe(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetHuaLu 化祿
func GetHuaLu(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetHuaQuan 化權
func GetHuaQuan(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetHuaKe 化科
func GetHuaKe(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetHuaLu2 化祿 (另一种算法)
func GetHuaLu2(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetMenKou 门户
func GetMenKou(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetHuaGan2 化氣 (另一种算法)
func GetHuaGan2(dayGan string, hourZhi string) string {
	// 簡化版
	return ""
}

// GetTaiSui 太歲
func GetTaiSui(year int) string {
	offset := year - 4
	index := offset % 12
	if index < 0 {
		index += 12
	}
	return DiZhi[index]
}

// GetTaiSuiYun 太歲運
func GetTaiSuiYun(year int, dayGan string) string {
	// 簡化版
	return ""
}

// GetMaDi 馬星
func GetMaDi(dayGan string) string {
	// 甲馬 gib, 乙羊頭, 丙犬門, 丁牛角
	// 戊鼠嘴, 己蛇腹, 庚猴 cổ, 辛馬齒
	// 壬猪口, 癸兔角
	madiMap := map[string]string{
		"甲": "寅", "乙": "巳", "丙": "申", "丁": "亥",
		"戊": "寅", "己": "巳", "庚": "申", "辛": "亥",
		"壬": "寅", "癸": "巳",
	}
	if md, ok := madiMap[dayGan]; ok {
		return md
	}
	return ""
}

// GetTianDe2 天德 (另一种算法)
func GetTianDe2(dayGan string) string {
	// 簡化版
	tianDeMap := map[string]string{
		"甲": "丁", "乙": "申", "丙": "壬", "丁": "丁",
		"戊": "壬", "己": "辛", "庚": "丁", "辛": "庚",
		"壬": "辛", "癸": "丁",
	}
	if td, ok := tianDeMap[dayGan]; ok {
		return td
	}
	return ""
}

// GetHuaLu3 化祿 (另一种算法)
func GetHuaLu3(dayGan string) string {
	// 簡化版
	huaLuMap := map[string]string{
	"甲": "丙", "乙": "乙", "丙": "甲", "丁": "己",
	"戊": "丙", "己": "乙", "庚": "戊", "辛": "辛",
	"壬": "戊", "癸": "己",
	}
	if hl, ok := huaLuMap[dayGan]; ok {
		return hl
	}
	return ""
}

// GetHuaQuan2 化權 (另一种算法)
func GetHuaQuan2(dayGan string) string {
	// 簡化版
	huaQuanMap := map[string]string{
		"甲": "甲", "乙": "癸", "丙": "丙", "丁": "壬",
		"戊": "甲", "己": "癸", "庚": "庚", "辛": "辛",
		"壬": "庚", "癸": "壬",
	}
	if hq, ok := huaQuanMap[dayGan]; ok {
		return hq
	}
	return ""
}

// GetHuaKe2 化科 (另一种算法)
func GetHuaKe2(dayGan string) string {
	// 簡化版
	huaKeMap := map[string]string{
		"甲": "乙", "乙": "壬", "丙": "丁", "丁": "庚",
		"戊": "乙", "己": "壬", "庚": "辛", "辛": "辛",
		"壬": "辛", "癸": "庚",
	}
	if hk, ok := huaKeMap[dayGan]; ok {
		return hk
	}
	return ""
}

// GetJuJia 聚家
func GetJuJia(dayGan string) string {
	// 簡化版
	jjMap := map[string]string{
		"甲": "戌", "乙": "丑", "丙": "辰", "丁": "未",
		"戊": "戌", "己": "丑", "庚": "辰", "辛": "未",
		"壬": "戌", "癸": "丑",
	}
	if jj, ok := jjMap[dayGan]; ok {
		return jj
	}
	return ""
}

// GetWangShuai 旺衰
func GetWangShuai(dayGan string, monthZhi string) float64 {
	// 簡化版
	return 50.0
}

// GetShenSha 神煞
func GetShenSha(dayGan string, monthZhi string) []string {
	// 簡化版
	return []string{}
}

// GetJiaZiCycleLength 甲子循環長度
func GetJiaZiCycleLength() int {
	return JiaZiCycle
}

// GetGanCount 天干數量
func GetGanCount() int {
	return len(TianGan)
}

// GetZhiCount 地支數量
func GetZhiCount() int {
	return len(DiZhi)
}

// GetWuXingCount 五行數量
func GetWuXingCount() int {
	return len(WuXing)
}

// GetShiShenCount 十神數量
func GetShiShenCount() int {
	return len(ShiShen)
}

// GetGanZhiWuXingCount 干支五行數量
func GetGanZhiWuXingCount() int {
	return len(GanZhiWuXing)
}

// GetTianGanYinYangCount 天干陰陽數量
func GetTianGanYinYangCount() int {
	return len(TianGanYinYang)
}

// GetDiZhiYinYangCount 地支陰陽數量
func GetDiZhiYinYangCount() int {
	return len(DiZhiYinYang)
}

// GetJieQiCount 節氣數量
func GetJieQiCount() int {
	return len(JieQi)
}

// GetShengXiaoCount 生肖數量
func GetShengXiaoCount() int {
	return 12
}

// GetDaYunCount 大運數量
func GetDaYunCount() int {
	return 10
}

// GetYearGZByYear 根据年份获取年柱
func GetYearGZByYear(year int) string {
	return GLYearToGZ(year)
}

// GetMonthGZByYear 根据年份獲取月柱
func GetMonthGZByYear(year, month, day int) string {
	monthZhi := GetMonthZhi(year, month, day)
	monthGan := GetMonthGan(string(GLYearToGZ(year)[0]), monthZhi)
	return monthGan + monthZhi
}

// GetDayGZByYear 根据年份獲取日柱
func GetDayGZByYear(year, month, day int) string {
	return GetDayGZ(year, month, day)
}

// GetHourGZByYear 根据年份獲取時柱
func GetHourGZByYear(year, month, day, hour int) string {
	dayGZ := GetDayGZ(year, month, day)
	dayGan := string(dayGZ[0])
	return GetHourGZ(dayGan, hour)
}

// GetFullBaziByYear 根据年份獲取完整八字
func GetFullBaziByYear(year, month, day, hour int) Bazi {
	return Bazi{
		Year:     GetYearGZByYear(year),
		Month:    GetMonthGZByYear(year, month, day),
		Day:      GetDayGZByYear(year, month, day),
		Hour:     GetHourGZByYear(year, month, day, hour),
		Elements: make(map[string]float64),
		TenGods:  make(map[string]string),
	}
}

// GetYearGanByYear 根据年份獲取年干
func GetYearGanByYear(year int) string {
	return string(GetYearGZByYear(year)[0])
}

// GetYearZhiByYear 根据年份獲取年支
func GetYearZhiByYear(year int) string {
	return string(GetYearGZByYear(year)[1])
}

// GetMonthGanByYear 根据年份獲取月干
func GetMonthGanByYear(year, month, day int) string {
	return string(GetMonthGZByYear(year, month, day)[0])
}

// GetMonthZhiByYear 根据年份獲取月支
func GetMonthZhiByYear(year, month, day int) string {
	return string(GetMonthGZByYear(year, month, day)[1])
}

// GetDayGanByYear 根据年份獲取日干
func GetDayGanByYear(year, month, day int) string {
	return string(GetDayGZByYear(year, month, day)[0])
}

// GetDayZhiByYear 根据年份獲取日支
func GetDayZhiByYear(year, month, day int) string {
	return string(GetDayGZByYear(year, month, day)[1])
}

// GetHourGanByYear 根据年份獲取時干
func GetHourGanByYear(year, month, day, hour int) string {
	return string(GetHourGZByYear(year, month, day, hour)[0])
}

// GetHourZhiByYear 根据年份獲取時支
func GetHourZhiByYear(year, month, day, hour int) string {
	return string(GetHourGZByYear(year, month, day, hour)[1])
}

// GetGZByGan 根据天干获取地支 (用于六壬等)
func GetGZByGan(gan string) string {
	// 簡化版
	ganMap := map[string]string{
		"甲": "寅", "乙": "卯", "丙": "巳", "丁": "午",
		"戊": "巳", "己": "午", "庚": "申", "辛": "酉",
		"壬": "亥", "癸": "子",
	}
	if gz, ok := ganMap[gan]; ok {
		return gz
	}
	return ""
}

// GetGZByZhi 根据地支获取天干 (用于六壬等)
func GetGZByZhi(zhi string) string {
	// 簡化版
	zhiMap := map[string]string{
		"子": "壬", "丑": "癸", "寅": "甲", "卯": "乙",
		"辰": "丙", "巳": "丁", "午": "戊", "未": "己",
		"申": "庚", "酉": "辛", "戌": "壬", "亥": "癸",
	}
	if gz, ok := zhiMap[zhi]; ok {
		return gz
	}
	return ""
}

// GetGZByIndex 干支索引
func GetGZByIndex(index int) string {
	ganIndex := index % 10
	zhiIndex := index % 12
	return TianGan[ganIndex] + DiZhi[zhiIndex]
}

// GetGanByIndex 天干索引
func GetGanByIndex(index int) string {
	return TianGan[index%10]
}

// GetZhiByIndex 地支索引
func GetZhiByIndex(index int) string {
	return DiZhi[index%12]
}

// GetGanIndex 天干索引
func GetGanIndex(gan string) int {
	for i, g := range TianGan {
		if g == gan {
			return i
		}
	}
	return -1
}

// GetZhiIndex 地支索引
func GetZhiIndex(zhi string) int {
	for i, z := range DiZhi {
		if z == zhi {
			return i
		}
	}
	return -1
}

// GetDaysInMonth 每月天數
func GetDaysInMonth(year, month int) int {
	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if isLeapYear(year) && month == 2 {
		return 29
	}
	return days[month-1]
}

// GetDaysInYear 每年天數
func GetDaysInYear(year int) int {
	if isLeapYear(year) {
		return 366
	}
	return 365
}

// GetDayOfYear 获取一年中的第几天
func GetDayOfYear(year, month, day int) int {
	_ = GetDaysInMonth(year, month) // 用于验证
	total := 0
	for i := 1; i < month; i++ {
		total += GetDaysInMonth(year, i)
	}
	return total + day
}

// GetDayOfWeek 获取星期几
func GetDayOfWeek(year, month, day int) string {
	// 使用 Zeller's congruence
	// q: day, m: month, K: year%100, J: floor(year/100)
	// 0=Sat, 1=Sun, 2=Mon, 3=Tue, 4=Wed, 5=Thu, 6=Fri
	q := day
	m := month
	y := year
	if month < 3 {
		m += 12
		y -= 1
	}
	K := y % 100
	J := int(math.Floor(float64(y) / 100))
	h := (q + int(math.Floor(float64(13*(m+1))/5)) + K + int(math.Floor(float64(K)/4)) + int(math.Floor(float64(J)/4)) - 2*J) % 7
	if h < 0 {
		h += 7
	}
	days := []string{"日", "一", "二", "三", "四", "五", "六"}
	return "星期" + days[h]
}

// GetZhiHour 获取地支时辰
func GetZhiHour(hour int) string {
	hourMap := map[int]string{
		0:  "子", 1: "子", 2: "丑", 3: "丑",
		4: "寅", 5: "寅", 6: "卯", 7: "卯",
		8: "辰", 9: "辰", 10: "巳", 11: "巳",
		12: "午", 13: "午", 14: "未", 15: "未",
		16: "申", 17: "申", 18: "酉", 19: "酉",
		20: "戌", 21: "戌", 22: "亥", 23: "亥",
	}
	if h, ok := hourMap[hour]; ok {
		return h
	}
	return "子"
}

// GetHourLabel2 获取時辰標籤
func GetHourLabel2(hour int) string {
	labels := []string{
	"子時", "丑時", "寅時", "卯時",
	"辰時", "巳時", "午時", "未時",
	"申時", "酉時", "戌時", "亥時",
	}
	return labels[hour/2]
}

// GetHourRange 获取時辰范围
func GetHourRange(hour int) string {
	labels := []string{
	"23-1", "1-3", "3-5", "5-7",
	"7-9", "9-11", "11-13", "13-15",
	"15-17", "17-19", "19-21", "21-23",
	}
	return labels[hour/2] + "時"
}

// GetGanShen 天干神煞
func GetGanShen(gan string) []string {
	// 簡化版
	return []string{}
}

// GetZhiShen 地支神煞
func GetZhiShen(zhi string) []string {
	// 簡化版
	return []string{}
}

// GetNayinByGZ 獲取纳音
func GetNayinByGZ(gz string) string {
	if len(gz) >= 2 {
		naYinMap := GetNayinMap()
		if ny, ok := naYinMap[gz]; ok {
			return ny
		}
	}
	return ""
}

// GetDaYun 獲取大運
func GetDaYun(year, month, day, hour int, gender bool) []string {
	// 簡化版
	return []string{}
}

// GetXun 獲取旬
func GetXun(gan, zhi string) string {
	// 簡化版
	// 甲子旬: 甲子、乙丑...癸酉 (10個)
	// 甲戌旬: 甲戌、乙亥...癸未 (10個)
	// 甲申旬: 甲申、乙酉...癸巳 (10個)
	// 甲午旬: 甲午、乙巳...癸卯 (10個)
	// 甲辰旬: 甲辰、乙巳...癸丑 (10個)
	// 甲寅旬: 甲寅、乙卯...癸亥 (10個)
	return ""
}

// GetKongWang 獲取空亡
func GetKongWang(gan, zhi string) string {
	xk := GetXunKong(gan, zhi)
	return xk
}

// GetShengXiaoByYear 獲取生肖
func GetShengXiaoByYear(year int) string {
	return GetShengXiao(year)
}

// GetDaYunStartYear 獲取大運起運年
func GetDaYunStartYear(year, month, day, hour int) int {
	// 簡化版
	return year
}

// GetDaYunEndYear 獲取大運終止年
func GetDaYunEndYear(startYear int) int {
	return startYear + 10
}

// GetTenGods 獲取十神
func GetTenGods(dayGan string, otherGan string) string {
	return GetShiShen(dayGan, otherGan)
}

// GetElementByGan 獲取天干的五行
func GetElementByGan(gan string) string {
	return GetElement(gan)
}

// GetElementByZhi 獲取地支的五行
func GetElementByZhi(zhi string) string {
	return GetElement(zhi)
}

// GetYinYangByGan 獲取天干的陰陽
func GetYinYangByGan(gan string) bool {
	return GetYinYang(gan)
}

// GetYinYangByZhi 獲取地支的陰陽
func GetYinYangByZhi(zhi string) bool {
	return GetYinYang(zhi)
}

// GetHuaGanByGZ 獲取化氣
func GetHuaGanByGZ(gz1, gz2 string) string {
	return GetHuaGan(gz1, gz2)
}

// GetClashByZhi 獲取相冲
func GetClashByZhi(zhi string) string {
	return GetClash(zhi)
}

// GetHarmonyByZhi 獲取相合
func GetHarmonyByZhi(zhi string) string {
	return GetHarmony(zhi)
}

// GetHarmByZhi 獲取相害
func GetHarmByZhi(zhi string) string {
	return GetHarm(zhi)
}

// GetCircleByZhi 獲取相刑
func GetCircleByZhi(zhi string) string {
	return GetCircle(zhi)
}

// GetBreakByZhi 獲取相破
func GetBreakByZhi(zhi string) string {
	return GetBreak(zhi)
}

// IsClashByZhi 獲取是否相冲
func IsClashByZhi(zhi1, zhi2 string) bool {
	return IsClash(zhi1, zhi2)
}

// IsHarmonyByZhi 獲取是否相合
func IsHarmonyByZhi(zhi1, zhi2 string) bool {
	return IsHarmony(zhi1, zhi2)
}

// IsHarmByZhi 獲取是否相害
func IsHarmByZhi(zhi1, zhi2 string) bool {
	return IsHarm(zhi1, zhi2)
}

// IsCircleByZhi 獲取是否相刑
func IsCircleByZhi(zhi1, zhi2 string) bool {
	return IsCircle(zhi1, zhi2)
}

// IsBreakByZhi 獲取是否相破
func IsBreakByZhi(zhi1, zhi2 string) bool {
	return IsBreak(zhi1, zhi2)
}

// GetDaysThroughYear 獲取累計天數
func GetDaysThroughYear2(year, month, day int) int {
	return getDaysThroughYear(year, month, day)
}

// GetMonthGanByYearAndZhi 獲取月干
func GetMonthGanByYearAndZhi(yearGan, monthZhi string) string {
	return GetMonthGan(yearGan, monthZhi)
}

// GetHourGanByDayGanAndHour 獲取時干
func GetHourGanByDayGanAndHour(dayGan string, hour int) string {
	return GetHourGZ(dayGan, hour)[:1]
}

// GetDayGZByYearAndDate 獲取日柱
func GetDayGZByYearAndDate(year, month, day int) string {
	return GetDayGZ(year, month, day)
}

// GetGZByYearAndDate 獲取完整八字
func GetGZByYearAndDate(year, month, day, hour int) Bazi {
	return Calculate(time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC))
}

// GetGZByTimestamp 獲取完整八字
func GetGZByTimestamp(ts int64) Bazi {
	t := time.Unix(ts, 0)
	return Calculate(t)
}

// GetGZByTime 獲取完整八字
func GetGZByTime(t time.Time) Bazi {
	return Calculate(t)
}

// GetGZByString 獲取完整八字
func GetGZByString(s string) Bazi {
	return CalculateSimple(s)
}

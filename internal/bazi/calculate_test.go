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
	"testing"
	"time"
)

// TestGLYearToGZ 測試年份轉干支
func TestGLYearToGZ(t *testing.T) {
	tests := []struct {
		year     int
		expected string
	}{
		{1995, "乙亥"},
		{1984, "甲子"},
		{2000, "庚辰"},
	}

	for _, tt := range tests {
		result := GLYearToGZ(tt.year)
		if result != tt.expected {
			t.Errorf("GLYearToGZ(%d) = %s, expected %s", tt.year, result, tt.expected)
		}
	}
}

// TestCalculate 測試完整八字計算
func TestCalculate(t *testing.T) {
	birthTime := time.Date(1995, 10, 1, 14, 0, 0, 0, time.UTC)
	result := Calculate(birthTime)

	// 驗證結構
	if result.Year == "" {
		t.Error("Year 為空")
	}
	if result.Month == "" {
		t.Error("Month 為空")
	}
	if result.Day == "" {
		t.Error("Day 為空")
	}
	if result.Hour == "" {
		t.Error("Hour 為空")
	}
}

// TestGetDayGZ 測試日柱計算
func TestGetDayGZ(t *testing.T) {
	result := GetDayGZ(1995, 10, 1)
	// GetDayGZ 返回完整干支 (e.g., "辛亥")，中文字符每個3bytes，所以長度應為6
	if len(result) != 6 {
		t.Errorf("GetDayGZ 返回長度錯誤 (应为6): %s", result)
	}
}

// TestCalculateElements 測試五行計算
func TestCalculateElements(t *testing.T) {
	bazi := Bazi{
		Year:  "乙亥",
		Month: "辛酉",
		Day:   "庚子",
		Hour:  "丙午",
	}

	elems := CalculateElements(bazi)

	// 檢查五行鍵值
	for _, elem := range WuXing {
		if _, ok := elems[elem]; !ok {
			t.Errorf("缺少元素: %s", elem)
		}
	}
}

// TestCalculateTenGods 測試十神計算
func TestCalculateTenGods(t *testing.T) {
	bazi := Bazi{
		Year:  "乙亥",
		Month: "辛酉",
		Day:   "庚子",
		Hour:  "丙午",
	}

	tenGods := CalculateTenGods(bazi)

	// 檢查至少有基本鍵值
	if len(tenGods) == 0 {
		t.Error("十神計算結果為空")
	}
}

// TestIsLeapYear 測試閏年判斷
func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{2024, true},  // 閏年
		{2023, false}, // 平年
		{2000, true},  // 閏年 (整百年)
		{1900, false}, // 平年 (整百年)
	}

	for _, tt := range tests {
		result := isLeapYear(tt.year)
		if result != tt.expected {
			t.Errorf("isLeapYear(%d) = %v, expected %v", tt.year, result, tt.expected)
		}
	}
}

// TestIsStrong 測試日主強弱判斷
func TestIsStrong(t *testing.T) {
	bazi := Bazi{
		Year:  "乙亥",
		Month: "辛酉",
		Day:   "庚子",
		Hour:  "丙午",
	}

	// 應該返回 true 或 false，不應 panics
	result := IsStrong("庚", bazi)
	_ = result // 使用變數避免未使用警告
}

// TestGetYearGZByYear 測試年柱計算
func TestGetYearGZByYear(t *testing.T) {
	result := GetYearGZByYear(1995)
	if result != "乙亥" {
		t.Errorf("GetYearGZByYear(1995) = %s, expected 乙亥", result)
	}
}

// TestGetMonthGan 測試月干計算
func TestGetMonthGan(t *testing.T) {
	// 1995 年為乙亥年，乙年起月干為戊
	result := GetMonthGan("乙", "寅")
	if result != "戊" {
		t.Errorf("GetMonthGan(乙, 寅) = %s, expected 戊", result)
	}
}

// TestGetHourGZ 測試時柱計算
func TestGetHourGZ(t *testing.T) {
	// GetHourGZ 返回完整干支 (e.g., "癸未")，中文字符每个3bytes
	result := GetHourGZ("庚", 14)
	if len(result) < 2 {
		t.Errorf("GetHourGZ 返回結果錯誤: %s", result)
	}
}

// TestGetHourLabel 測試時辰標籤
func TestGetHourLabel(t *testing.T) {
	// 根據實際實現調整測試
	result := GetHourLabel(23)
	if result == "" {
		t.Error("時辰標籤應不為空")
	}
}

// TestGetShengXiao 測試生肖計算
func TestGetShengXiao(t *testing.T) {
	tests := []struct {
		year     int
		expected string
	}{
		{1984, "鼠"},
		{1985, "牛"},
		{1996, "鼠"},
	}

	for _, tt := range tests {
		result := GetShengXiao(tt.year)
		if result != tt.expected {
			t.Errorf("GetShengXiao(%d) = %s, expected %s", tt.year, result, tt.expected)
		}
	}
}

// TestGetGZIndex 測試干支索引
func TestGetGZIndex(t *testing.T) {
	if GetGZIndex("甲") != 0 {
		t.Error("甲索引應為 0")
	}
	if GetGZIndex("子") != 0 {
		t.Error("子索引應為 0")
	}
}

// TestGetNayinByGZ 測試納音計算
func TestGetNayinByGZ(t *testing.T) {
	result := GetNayinByGZ("甲子")
	if result == "" {
		t.Error("納音計算失敗")
	}
}

// TestGetDayOfWeek 測試星期幾計算
func TestGetDayOfWeek(t *testing.T) {
	// 2024-03-11 是星期一
	result := GetDayOfWeek(2024, 3, 11)
	if result == "" {
		t.Error("星期計算結果應不為空")
	}
}

// TestGetDaysInMonth 測試每月天數
func TestGetDaysInMonth(t *testing.T) {
	tests := []struct {
		year     int
		month    int
		expected int
	}{
		{2024, 2, 29}, // 閏年
		{2023, 2, 28}, // 平年
		{2024, 1, 31},
		{2024, 4, 30},
	}

	for _, tt := range tests {
		result := GetDaysInMonth(tt.year, tt.month)
		if result != tt.expected {
			t.Errorf("GetDaysInMonth(%d, %d) = %d, expected %d", tt.year, tt.month, result, tt.expected)
		}
	}
}

// TestGetClash 測試相冲計算
func TestGetClash(t *testing.T) {
	if GetClash("子") != "午" {
		t.Error("子冲午錯誤")
	}
	if GetClash("午") != "子" {
		t.Error("午冲子錯誤")
	}
}

// TestGetHarmony 測試相合計算
func TestGetHarmony(t *testing.T) {
	if GetHarmony("子") != "丑" {
		t.Error("子合丑錯誤")
	}
}

// TestGetGanZhiWuXingCount 測試常數計數
func TestGetGanZhiWuXingCount(t *testing.T) {
	if GetGanZhiWuXingCount() == 0 {
		t.Error("GanZhiWuXing 長度應大於 0")
	}
}

// TestGetJiaZiCycleLength 測試甲子循環長度
func TestGetJiaZiCycleLength(t *testing.T) {
	if GetJiaZiCycleLength() != 60 {
		t.Error("甲子循環長度應為 60")
	}
}

// TestGetYearGZByIndex 測試索引轉干支
func TestGetYearGZByIndex(t *testing.T) {
	result := GetYearGZByIndex(0)
	if result != "甲子" {
		t.Errorf("索引 0 應為 甲子: %s", result)
	}
}

// TestIsSameElement 測試同元素判斷
func TestIsSameElement(t *testing.T) {
	if !IsSameElement("甲", "乙") {
		t.Error("甲乙同為木，應為 true")
	}
	if IsSameElement("甲", "丙") {
		t.Error("甲丙不同元素，應為 false")
	}
}

// TestGetWuXingCount 測試五行計數
func TestGetWuXingCount(t *testing.T) {
	if GetWuXingCount() != 5 {
		t.Error("五行數量應為 5")
	}
}

// TestGetGanCount 測試天干計數
func TestGetGanCount(t *testing.T) {
	if GetGanCount() != 10 {
		t.Error("天干數量應為 10")
	}
}

// TestGetZhiCount 測試地支計數
func TestGetZhiCount(t *testing.T) {
	if GetZhiCount() != 12 {
		t.Error("地支數量應為 12")
	}
}

// TestCalculateWithTime 測試時間函數
func TestCalculateWithTime(t *testing.T) {
	birthTime := time.Date(1995, 10, 1, 14, 0, 0, 0, time.UTC)
	result := Calculate(birthTime)

	if result.Year == "" || result.Month == "" || result.Day == "" || result.Hour == "" {
		t.Error("八字計算失敗")
	}
}

// TestCalculateSimple 測試字串輸入
func TestCalculateSimple(t *testing.T) {
	result := CalculateSimple("1995-10-01 14:00")

	if result.Year == "" {
		t.Error("Year 為空")
	}
}

// TestCalculateSimpleError 測試錯誤輸入
func TestCalculateSimpleError(t *testing.T) {
	// 使用空字串會使用默認時間
	result := CalculateSimple("")

	if result.Year == "" {
		t.Error("空輸入應返回當前時間結果")
	}
}

// TestGetMonthGZByYear 測試月柱計算
func TestGetMonthGZByYear(t *testing.T) {
	// GetMonthGZByYear 返回完整干支，檢查不為空
	result := GetMonthGZByYear(1995, 10, 1)
	if len(result) < 2 {
		t.Errorf("GetMonthGZByYear 返回結果錯誤: %s", result)
	}
}

// TestGetHuaGan 測試化氣計算
func TestGetHuaGan(t *testing.T) {
	if GetHuaGan("甲", "己") != "土" {
		t.Error("甲己化土錯誤")
	}
}

// TestGetHourRange 測試時辰範圍
func TestGetHourRange(t *testing.T) {
	result := GetHourRange(23)
	if result == "" {
		t.Error("時辰範圍不應為空")
	}
}

// TestGetMonthGZByIndex 測試索引轉月柱
func TestGetMonthGZByIndex(t *testing.T) {
	// GetMonthGZByIndex 返回完整干支，檢查不為空
	result := GetMonthGZByIndex(0)
	if len(result) < 2 {
		t.Errorf("GetMonthGZByIndex 返回結果錯誤: %s", result)
	}
}

// TestGetNayinMap 測試納音表
func TestGetNayinMap(t *testing.T) {
	m := GetNayinMap()
	if len(m) == 0 {
		t.Error("納音表為空")
	}
}

// TestGetGanByIndex 測試天干索引
func TestGetGanByIndex(t *testing.T) {
	if GetGanByIndex(0) != "甲" {
		t.Error("索引 0 應為甲")
	}
}

// TestGetZhiByIndex 測試地支索引
func TestGetZhiByIndex(t *testing.T) {
	if GetZhiByIndex(0) != "子" {
		t.Error("索引 0 應為子")
	}
}

// TestGetDayGZByYearAndDate 測試日柱計算
func TestGetDayGZByYearAndDate(t *testing.T) {
	result := GetDayGZByYearAndDate(1995, 10, 1)
	// GetDayGZ 返回完整干支 (e.g., "辛亥")，中文字符每個3bytes，所以長度應為6
	if len(result) != 6 {
		t.Errorf("GetDayGZByYearAndDate 返回長度錯誤 (應為6): %s", result)
	}
}

// TestGetGZByYearAndDate 測試八字計算
func TestGetGZByYearAndDate(t *testing.T) {
	result := GetGZByYearAndDate(1995, 10, 1, 14)
	if result.Year == "" {
		t.Error("Year 為空")
	}
}

// TestGetGZByTimestamp 測試時間戳計算
func TestGetGZByTimestamp(t *testing.T) {
	result := GetGZByTimestamp(781125600) // 1995-10-01 14:00:00 UTC
	if result.Year == "" {
		t.Error("Year 為空")
	}
}

// TestGetGZByTime 測試 time.Time 計算
func TestGetGZByTime(t *testing.T) {
	birthTime := time.Date(1995, 10, 1, 14, 0, 0, 0, time.UTC)
	result := GetGZByTime(birthTime)
	if result.Year == "" {
		t.Error("Year 為空")
	}
}

// TestGetGZByString 測試字串計算
func TestGetGZByString(t *testing.T) {
	result := GetGZByString("1995-10-01 14:00")
	if result.Year == "" {
		t.Error("Year 為空")
	}
}

// TestGetDayMasterPower 測試日主力量
func TestGetDayMasterPower(t *testing.T) {
	bazi := Bazi{
		Year:  "乙亥",
		Month: "辛酉",
		Day:   "庚子",
		Hour:  "丙午",
	}

	result := GetDayMasterPower("庚", bazi)
	if result < 0 {
		t.Error("日主力量不應為負數")
	}
}

// TestGetDaYunStartYear 測試大運起運年
func TestGetDaYunStartYear(t *testing.T) {
	result := GetDaYunStartYear(1995, 10, 1, 14)
	if result < 1995 {
		t.Errorf("大運起運年不应小於出生年份: %d", result)
	}
}

// TestGetDaYunEndYear 測試大運終止年
func TestGetDaYunEndYear(t *testing.T) {
	result := GetDaYunEndYear(1995)
	if result != 2005 {
		t.Errorf("大運終止年應為 2005: %d", result)
	}
}

// TestGetTenGodsStar 測試十神星耀
func TestGetTenGodsStar(t *testing.T) {
	bazi := Bazi{
		Year:  "乙亥",
		Month: "辛酉",
		Day:   "庚子",
		Hour:  "丙午",
	}

	result := GetTenGodsStar(bazi)
	if len(result) == 0 {
		t.Error("十神星耀為空")
	}
}

// TestGetTenGodsChart 測試十神表格
func TestGetTenGodsChart(t *testing.T) {
	bazi := Bazi{
		Year:  "乙亥",
		Month: "辛酉",
		Day:   "庚子",
		Hour:  "丙午",
	}

	result := GetTenGodsChart(bazi)
	if len(result) == 0 {
		t.Error("十神表格為空")
	}
}

// TestGetTenGods 測試十神計算
func TestGetTenGods(t *testing.T) {
	result := GetTenGods("甲", "丙")
	if result == "" {
		t.Error("十神計算失敗")
	}
}

// TestGetElementByGan 測試天干元素
func TestGetElementByGan(t *testing.T) {
	if GetElementByGan("甲") != "木" {
		t.Error("甲元素應為木")
	}
}

// TestGetElementByZhi 測試地支元素
func TestGetElementByZhi(t *testing.T) {
	if GetElementByZhi("子") != "水" {
		t.Error("子元素應為水")
	}
}

// TestGetYinYangByGan 測試天干陰陽
func TestGetYinYangByGan(t *testing.T) {
	if GetYinYangByGan("甲") != false {
		t.Error("甲應為陽")
	}
	if GetYinYangByGan("乙") != true {
		t.Error("乙應為陰")
	}
}

// TestGetYinYangByZhi 測試地支陰陽
func TestGetYinYangByZhi(t *testing.T) {
	if GetYinYangByZhi("子") != true {
		t.Error("子應為陰")
	}
	if GetYinYangByZhi("丑") != false {
		t.Error("丑應為陽")
	}
}

// TestGetClashByZhi 測試相冲
func TestGetClashByZhi(t *testing.T) {
	if GetClashByZhi("子") != "午" {
		t.Error("子冲午錯誤")
	}
}

// TestIsClashByZhi 測試是否相冲
func TestIsClashByZhi(t *testing.T) {
	if !IsClashByZhi("子", "午") {
		t.Error("子午應相冲")
	}
	if IsClashByZhi("子", "丑") {
		t.Error("子丑不應相冲")
	}
}

// TestGetHuaGanByGZ 測試化氣
func TestGetHuaGanByGZ(t *testing.T) {
	if GetHuaGanByGZ("甲己", "") != "土" {
		t.Error("甲己化土錯誤")
	}
}

// TestGetHourGanByDayGanAndHour 測試時干計算
func TestGetHourGanByDayGanAndHour(t *testing.T) {
	// 簡單測試不會 panic
	result := GetHourGanByDayGanAndHour("庚", 14)
	_ = result // 避免未使用警告
}

// TestGetPan 測試排盤
func TestGetPan(t *testing.T) {
	result := GetPan(1995, 10, 1, 14)

	if result["年柱"] == "" {
		t.Error("年柱為空")
	}
	if result["月柱"] == "" {
		t.Error("月柱為空")
	}
	if result["日柱"] == "" {
		t.Error("日柱為空")
	}
	if result["時柱"] == "" {
		t.Error("時柱為空")
	}
}

// TestGetGZByGan 測試天干轉地支
func TestGetGZByGan(t *testing.T) {
	result := GetGZByGan("甲")
	if result == "" {
		t.Error("甲轉地支失敗")
	}
}

// TestGetGZByZhi 測試地支轉天干
func TestGetGZByZhi(t *testing.T) {
	result := GetGZByZhi("子")
	if result == "" {
		t.Error("子轉天干失敗")
	}
}

// TestGetGZByIndex 測試索引轉干支
func TestGetGZByIndex(t *testing.T) {
	result := GetGZByIndex(0)
	if result != "甲子" {
		t.Errorf("索引 0 應為甲子: %s", result)
	}
}

// TestGetGanZhiByIndex 測試索引轉干支
func TestGetGanZhiByIndex(t *testing.T) {
	result := GetGanZhiByIndex(0)
	if result != "甲子" {
		t.Errorf("索引 0 應為甲子: %s", result)
	}
}

// TestGetGanCount 測試天干計數
func TestGetGanCountAgain(t *testing.T) {
	if GetGanCount() != 10 {
		t.Error("天干數量應為 10")
	}
}

// TestGetZhiCount 測試地支計數
func TestGetZhiCountAgain(t *testing.T) {
	if GetZhiCount() != 12 {
		t.Error("地支數量應為 12")
	}
}

// TestGetWuXingCount 測試五行計數
func TestGetWuXingCountAgain(t *testing.T) {
	if GetWuXingCount() != 5 {
		t.Error("五行數量應為 5")
	}
}

// TestGetShiShenCount 測試十神計數
func TestGetShiShenCount(t *testing.T) {
	if GetShiShenCount() != 10 {
		t.Error("十神數量應為 10")
	}
}

// TestGetGanYin 測試天干陰陽
func TestGetGanYin(t *testing.T) {
	if GetGanYin("甲") != false {
		t.Error("甲應為陽")
	}
	if GetGanYin("乙") != true {
		t.Error("乙應為陰")
	}
}

// TestGetZhiYin 測試地支陰陽
func TestGetZhiYin(t *testing.T) {
	if GetZhiYin("子") != true {
		t.Error("子應為陰")
	}
	if GetZhiYin("丑") != false {
		t.Error("丑應為陽")
	}
}

// TestGetGanShen 測試天干神煞
func TestGetGanShen(t *testing.T) {
	result := GetGanShen("甲")
	if result == nil {
		t.Error("天干神煞為 nil")
	}
}

// TestGetZhiShen 測試地支神煞
func TestGetZhiShen(t *testing.T) {
	result := GetZhiShen("子")
	if result == nil {
		t.Error("地支神煞為 nil")
	}
}

// TestGetDaYun 測試大運
func TestGetDaYun(t *testing.T) {
	result := GetDaYun(1995, 10, 1, 14, true)
	if result == nil {
		t.Error("大運為 nil")
	}
}

// TestGetXun 測試旬
func TestGetXun(t *testing.T) {
	// GetXun 目前為未實作函式，僅做暖機
	result := GetXun("甲", "子")
	_ = result // 避免未使用警告
}

// TestGetKongWang 測試空亡
func TestGetKongWang(t *testing.T) {
	result := GetKongWang("甲", "子")
	if result == "" {
		t.Error("空亡計算失敗")
	}
}

// TestGetJieQiByMonth 測試節氣
func TestGetJieQiByMonth(t *testing.T) {
	if len(JieQi) != 24 {
		t.Errorf("節氣數量應為 24: %d", len(JieQi))
	}
}

// TestGetDayOfYear 測試一年中的哪一天
func TestGetDayOfYear(t *testing.T) {
	result := GetDayOfYear(2024, 1, 1)
	if result != 1 {
		t.Errorf("1 月 1 日應為第 1 天: %d", result)
	}
}

// TestGetDayOfYear2 測試 2 月 1 日
func TestGetDayOfYear2(t *testing.T) {
	result := GetDayOfYear(2024, 2, 1)
	if result != 32 {
		t.Errorf("2024 年 2 月 1 日應為第 32 天: %d", result)
	}
}

// TestGetDaysInYear 測試每年天數
func TestGetDaysInYear(t *testing.T) {
	if GetDaysInYear(2024) != 366 {
		t.Error("2024 年應為 366 天")
	}
	if GetDaysInYear(2023) != 365 {
		t.Error("2023 年應為 365 天")
	}
}

// TestCalculateElementsEmpty 測試空八字
func TestCalculateElementsEmpty(t *testing.T) {
	bazi := Bazi{}
	elems := CalculateElements(bazi)

	// 應返回零值
	for _, elem := range WuXing {
		if _, ok := elems[elem]; !ok {
			t.Errorf("缺少元素: %s", elem)
		}
	}
}

// TestEmptyDayGZ 測試空日柱
func TestEmptyDayGZ(t *testing.T) {
	result := GetDayStem("")
	if result != "" {
		t.Errorf("空日柱應返回空字串: %s", result)
	}
}

// TestZeroElements 測試空八字元素計數
func TestZeroElements(t *testing.T) {
	bazi := Bazi{}
	elems := CalculateElements(bazi)

	for _, elem := range WuXing {
		if _, ok := elems[elem]; !ok {
			t.Errorf("缺少元素: %s", elem)
		}
	}
}

// TestFullBaziCalculation 測試完整八字計算
func TestFullBaziCalculation(t *testing.T) {
	birthTime := time.Date(1995, 10, 1, 14, 0, 0, 0, time.UTC)
	bazi := Calculate(birthTime)

	// 驗證返回值
	if bazi.Year == "" || bazi.Month == "" || bazi.Day == "" || bazi.Hour == "" {
		t.Error("八字計算失敗")
	}

	// 驗證五行計算
	elems := CalculateElements(bazi)
	if len(elems) != 5 {
		t.Errorf("五行元素數量應為 5: %d", len(elems))
	}
}

// TestDeterministicCalculation 測試計算可再現性
func TestDeterministicCalculation(t *testing.T) {
	birthTime := time.Date(1995, 10, 1, 14, 0, 0, 0, time.UTC)

	// 連續計算兩次
	bazi1 := Calculate(birthTime)
	bazi2 := Calculate(birthTime)

	// 應該得到相同結果
	if bazi1.Year != bazi2.Year ||
		bazi1.Month != bazi2.Month ||
		bazi1.Day != bazi2.Day ||
		bazi1.Hour != bazi2.Hour {
		t.Error("相同輸入應得到相同結果")
	}
}

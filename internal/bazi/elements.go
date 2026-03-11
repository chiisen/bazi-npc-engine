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

import "math"

// ═══════════════════════════════════════════
// 💡 概念：五行Strength計算
// 說明：根據八字中的干支計算五行的強弱比例
// 為何使用：分析八字中金木水火土的平衡狀態
// 注意事項：五行失衡可能表示性格或命運上的特點
// ═══════════════════════════════════════════

// CalculateElements 計算八字中各元素的強弱比例
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string]float64: 各元素的強弱比例 (總和為 1.0)
//
// 計算邏輯：
//   1. 統計八字中每個元素的出現次數
//   2. 考慮地支藏干的影響
//   3. 計算比例並回傳
func CalculateElements(bazi Bazi) map[string]float64 {
	elems := map[string]float64{
		"金": 0, "木": 0, "水": 0, "火": 0, "土": 0,
	}

	// 統計八字中各元素的出現次數
	columns := []string{bazi.Year, bazi.Month, bazi.Day, bazi.Hour}
	for _, col := range columns {
		elem := GetColumnElement(col)
		elems[elem]++
	}

	// 考慮地支藏干的影響
	// 地支中隱藏的天干也會對應五行
	zhiList := []string{
		GetYearBranch(bazi.Year),
		GetMonthBranch(bazi.Month),
		GetDayBranch(bazi.Day),
		GetHourBranch(bazi.Hour),
	}

	for _, zhi := range zhiList {
		// 簡化：每個地支藏干加 0.3 的權重
		// 實際應根據不同地支的藏干詳細計算
		if elem := GetElement(zhi); elem != "" {
			elems[elem] += 0.3
		}
	}

	// 計算總和
	total := 0.0
	for _, v := range elems {
		total += v
	}

	// 轉換為比例
	if total > 0 {
		for k, v := range elems {
			elems[k] = v / total
		}
	}

	return elems
}

// CalculateElementScore 計算五行得分
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string]float64: 各元素的得分
//
// 得分計算：
//   - 基礎分：每個天干/地支各 1 分
//   - 藏干加成：地支藏干額外加分
//   - 生助加成：被生的元素額外加分
func CalculateElementScore(bazi Bazi) map[string]float64 {
	scores := map[string]float64{
		"金": 0, "木": 0, "水": 0, "火": 0, "土": 0,
	}

	// 計算天干地支的基礎分數
	columns := []string{bazi.Year, bazi.Month, bazi.Day, bazi.Hour}
	for _, col := range columns {
		elem := GetColumnElement(col)
		scores[elem]++

		// 地支加 0.5 分 (含藏干)
		if len(col) > 1 {
			elem2 := GetElement(string(col[1]))
			if elem2 != "" {
				scores[elem2] += 0.5
			}
		}
	}

	// 考慮生克關係
	// 被生的元素會間接獲得力量
	dayStem := GetDayStem(bazi.Day)
	dayElement := GetElement(dayStem)

	// 生我的元素 (印綬) 為我提供support
	gen := map[string]string{
		"金": "水", "水": "木", "木": "火", "火": "土", "土": "金",
	}
	for elem, source := range gen {
		if dayElement == source {
			scores[elem] += 0.5 // 印綬加成
		}
	}

	// 我生的元素 (食傷) 我的力量會減弱
	for elem, target := range gen {
		if dayElement == elem {
			scores[target] += 0.3 // 食傷發散
		}
	}

	return scores
}

// IsElementBalanced 判斷五行是否平衡
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - bool: 是否平衡
//   - []string: 失衡的元素列表
//
// 平衡標準：
//   - 每個元素的比例約在 0.15-0.35 之間
//   - 沒有元素過於強烈或過於弱勢
func IsElementBalanced(bazi Bazi) (bool, []string) {
	elems := CalculateElements(bazi)

	var imbalanced []string
	for elem, ratio := range elems {
		// 簡化標準：單個元素不能超過 0.4 或低於 0.1
		if ratio > 0.4 {
			imbalanced = append(imbalanced, elem+"過旺")
		} else if ratio < 0.1 {
			imbalanced = append(imbalanced, elem+"過弱")
		}
	}

	return len(imbalanced) == 0, imbalanced
}

// GetMissingElements 獲取缺失的元素
//
// 參數：
//   - bazi: 八字盤
//   - threshold: 閾值 (小於此值視為缺失)
//
// 回傳：
//   - []string: 缺失的元素列表
func GetMissingElements(bazi Bazi, threshold float64) []string {
	elems := CalculateElements(bazi)

	var missing []string
	for elem, ratio := range elems {
		if ratio < threshold {
			missing = append(missing, elem)
		}
	}

	return missing
}

// GetStrongElements 獲取強勢的元素
//
// 參數：
//   - bazi: 八字盤
//   - threshold: 閾值 (大於此值視為強勢)
//
// 回傳：
//   - []string: 強勢的元素列表
func GetStrongElements(bazi Bazi, threshold float64) []string {
	elems := CalculateElements(bazi)

	var strong []string
	for elem, ratio := range elems {
		if ratio >= threshold {
			strong = append(strong, elem)
		}
	}

	return strong
}

// GetWeakElements 獲取弱勢的元素
//
// 參數：
//   - bazi: 八字盤
//   - threshold: 閾值 (小於此值視為弱勢)
//
// 回傳：
//   - []string: 弱勢的元素列表
func GetWeakElements(bazi Bazi, threshold float64) []string {
	elems := CalculateElements(bazi)

	var weak []string
	for elem, ratio := range elems {
		if ratio < threshold {
			weak = append(weak, elem)
		}
	}

	return weak
}

// ElementPower 五行StrengthPOWER (0-100)
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string]float64: 各元素的POWER (0-100)
func ElementPower(bazi Bazi) map[string]float64 {
	scores := CalculateElementScore(bazi)

	// 計算最大值用於歸一化
	maxScore := 0.0
	for _, v := range scores {
		if v > maxScore {
			maxScore = v
		}
	}

	// 歸一化到 0-100
	power := make(map[string]float64)
	for elem, score := range scores {
		if maxScore > 0 {
			power[elem] = math.Round(score/maxScore*100) / 100 * 100
		} else {
			power[elem] = 0
		}
	}

	return power
}

// ElementCompatibility 元素相容性分析
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string]string: 相容性分析結果
//
// 分析內容：
//   - 元素是否過旺/過弱
//   - 元素之間的生克關係
//   - 是否有缺失元素
func ElementCompatibility(bazi Bazi) map[string]string {
	result := make(map[string]string)

	elems := CalculateElements(bazi)
	scores := CalculateElementScore(bazi)

	// 分析各元素狀態
	for elem, ratio := range elems {
		if ratio > 0.35 {
			result[elem] = "過旺"
		} else if ratio < 0.12 {
			result[elem] = "過弱"
		} else {
			result[elem] = "適中"
		}
	}

	// 分析生克關係
	gen := map[string]string{
		"金": "水", "水": "木", "木": "火", "火": "土", "土": "金",
	}
	克 := map[string]string{
		"金": "木", "木": "土", "土": "水", "水": "火", "火": "金",
	}

	dayStem := GetDayStem(bazi.Day)
	dayElement := GetElement(dayStem)

	// 我生 (食傷)
	for elem, target := range gen {
		if dayElement == elem {
			if scores[target] > 0 {
				result["食傷"] = "有表現力"
			}
		}
	}

	// 我克 (財)
	for elem, target := range 克 {
		if dayElement == elem {
			if scores[target] > 0 {
				result["財富"] = "有能力获取"
			}
		}
	}

	return result
}

// ElementHealth 元素健康度
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string]float64: 各元素的健康度 (0-100)
func ElementHealth(bazi Bazi) map[string]float64 {
	power := ElementPower(bazi)
	health := make(map[string]float64)

	// 健康度 = POWER 的平方根 * 10
	// 這樣可以壓縮高低差異
	for elem, p := range power {
		health[elem] = math.Sqrt(p) * 10
	}

	return health
}

// ElementRecommendation 元素補救建議
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string][]string: 各元素的建議
//
// 建議內容：
//   - 補救方式 (顏色、數字、方向等)
//   - 忌諱事項
func ElementRecommendation(bazi Bazi) map[string][]string {
	recommendations := make(map[string][]string)

	// 獲取缺失元素
	missing := GetMissingElements(bazi, 0.1)

	// 獲取過旺元素
	strong := GetStrongElements(bazi, 0.35)

	// 火建議
	if contains(missing, "火") {
		recommendations["火"] = []string{
			"補救：多使用紅色、紫色",
			"方向：南方",
			"數字：2, 7",
			"季節：夏季",
		}
	}
	if contains(strong, "火") {
		recommendations["火"] = append(recommendations["火"],
			"忌諱：避免過度熱情、急躁",
		)
	}

	// 土建議
	if contains(missing, "土") {
		recommendations["土"] = []string{
			"補救：多使用黃色、棕色",
			"方向：中潯",
			"數字：5, 6",
			"季節：四季月",
		}
	}
	if contains(strong, "土") {
		recommendations["土"] = append(recommendations["土"],
			"忌諱：避免過度固執、優柔寡斷",
		)
	}

	// 金建議
	if contains(missing, "金") {
		recommendations["金"] = []string{
			"補救：多使用白色、金色",
			"方向：西方",
			"數字：4, 9",
			"季節：秋季",
		}
	}
	if contains(strong, "金") {
		recommendations["金"] = append(recommendations["金"],
			"忌諱：避免過度剛硬、冷漠",
		)
	}

	// 水建議
	if contains(missing, "水") {
		recommendations["水"] = []string{
			"補救：多使用黑色、藍色",
			"方向：北方",
			"數字：1, 6",
			"季節：冬季",
		}
	}
	if contains(strong, "水") {
		recommendations["水"] = append(recommendations["水"],
			"忌諱：避免過度Fluid、情緒化",
		)
	}

	// 木建議
	if contains(missing, "木") {
		recommendations["木"] = []string{
			"補救：多使用綠色、青色",
			"方向：東方",
			"數字：3, 8",
			"季節：春季",
		}
	}
	if contains(strong, "木") {
		recommendations["木"] = append(recommendations["木"],
			"忌諱：避免過度生長、競爭",
		)
	}

	return recommendations
}

// helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

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

// ═══════════════════════════════════════════
// 💡 概念：十神分析 (Shi Shen)
// 說明：根據日主 (出生日的天干) 與其他干支的生克關係判斷
// 為何使用：描述人生各方面的人格特質與命運走向
// 注意事項：十神代表與日主的關係，用於分析性格與人生領域
// ═══════════════════════════════════════════

// ShiShenMap 十神對應關係
// 生我者：印綬 (正印、偏印)
// 我生者：食傷 (食神、傷官)
// 我克者：財 (正財、偏財)
// 克我者：官殺 (正官、七殺)
// 同我者：比劫 (比肩、劫財)
var ShiShenMap = map[string]map[string]string{
	// 日干 -> 其他干 -> 十神
	"甲": {
		"甲": "比肩", "乙": "劫財",
		"丙": "食神", "丁": "傷官",
		"戊": "偏財", "己": "正財",
		"庚": "七殺", "辛": "正官",
		"壬": "偏印", "癸": "正印",
	},
	"乙": {
		"乙": "比肩", "甲": "劫財",
		"丁": "食神", "丙": "傷官",
		"己": "偏財", "戊": "正財",
		"辛": "七殺", "庚": "正官",
		"癸": "偏印", "壬": "正印",
	},
	"丙": {
		"丙": "比肩", "丁": "劫財",
		"戊": "食神", "己": "傷官",
		"庚": "偏財", "辛": "正財",
		"壬": "七殺", "癸": "正官",
		"甲": "偏印", "乙": "正印",
	},
	"丁": {
		"丁": "比肩", "丙": "劫財",
		"己": "食神", "戊": "傷官",
		"辛": "偏財", "庚": "正財",
		"癸": "七殺", "壬": "正官",
		"乙": "偏印", "甲": "正印",
	},
	"戊": {
		"戊": "比肩", "己": "劫財",
		"庚": "食神", "辛": "傷官",
		"壬": "偏財", "癸": "正財",
		"甲": "七殺", "乙": "正官",
		"丙": "偏印", "丁": "正印",
	},
	"己": {
		"己": "比肩", "戊": "劫財",
		"辛": "食神", "庚": "傷官",
		"癸": "偏財", "壬": "正財",
		"乙": "七殺", "甲": "正官",
		"丁": "偏印", "丙": "正印",
	},
	"庚": {
		"庚": "比肩", "辛": "劫財",
		"壬": "食神", "癸": "傷官",
		"甲": "偏財", "乙": "正財",
		"丙": "七殺", "丁": "正官",
		"戊": "偏印", "己": "正印",
	},
	"辛": {
		"辛": "比肩", "庚": "劫財",
		"癸": "食神", "壬": "傷官",
		"乙": "偏財", "甲": "正財",
		"丁": "七殺", "丙": "正官",
		"己": "偏印", "戊": "正印",
	},
	"壬": {
		"壬": "比肩", "癸": "劫財",
		"甲": "食神", "乙": "傷官",
		"丙": "偏財", "丁": "正財",
		"戊": "七殺", "己": "正官",
		"庚": "偏印", "辛": "正印",
	},
	"癸": {
		"癸": "比肩", "壬": "劫財",
		"乙": "食神", "甲": "傷官",
		"丁": "偏財", "丙": "正財",
		"己": "七殺", "戊": "正官",
		"辛": "偏印", "庚": "正印",
	},
}

// CalculateTenGods 計算八字中各柱的十神
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string]string: 各柱的十神分析
//
// 範例：
//   map[string]string{
//     "年柱": "正官",
//     "月柱": "偏印",
//     "日柱": "日主",
//     "時柱": "偏財",
//   }
func CalculateTenGods(bazi Bazi) map[string]string {
	result := make(map[string]string)

	// 日柱為日主
	dayStem := GetDayStem(bazi.Day)
	result["日柱"] = "日主"

	// 計算其他柱的十神
	tenGodsMap := ShiShenMap[dayStem]
	if tenGodsMap == nil {
		tenGodsMap = ShiShenMap["甲"] // 預設
	}

	// 年柱
	yearGan := GetYearStem(bazi.Year)
	if yearGan != "" {
		if shen, ok := tenGodsMap[yearGan]; ok {
			result["年柱"] = shen
		}
	}

	// 月柱
	monthGan := GetMonthStem(bazi.Month)
	if monthGan != "" {
		if shen, ok := tenGodsMap[monthGan]; ok {
			result["月柱"] = shen
		}
	}

	// 時柱
	hourGan := GetHourStem(bazi.Hour)
	if hourGan != "" {
		if shen, ok := tenGodsMap[hourGan]; ok {
			result["時柱"] = shen
		}
	}

	return result
}

// GetTenGodsDescription 獲取十神的描述
//
// 參數：
//   - shen: 十神名稱
//
// 回傳：
//   - string: 十神描述
//
// 十神描述：
//   - 比肩：同性同源，代表平等、競爭
//   - 劫財：异性同源，代表衝動、冒險
//   - 食神：我生する，代表溫和、表達
//   - 傷官：我生する，代表才華、叛逆
//   - 偏財：我克する，代表直覺、投機
//   - 正財：我克する，代表穩定、節儉
//   - 七殺：克我する，代表壓力、進取
//   - 正官：克我する，代表責任、守矩
//   - 偏印：生我する，代表偏門、特長
//   - 正印：生我する，代表正統、學習
func GetTenGodsDescription(shen string) string {
	description := map[string]string{
		"比肩": "同性同源，代表平等、競爭、朋友",
		"劫財": "异性同源，代表衝動、冒險、兄弟",
		"食神": "我生對象，代表溫和、表達、美食",
		"傷官": "我生對象，代表才華、叛逆、創作",
		"偏財": "我克對象，代表直覺、投機、意外財",
		"正財": "我克對象，代表穩定、節儉、現實財",
		"七殺": "克我對象，代表壓力、進取、權威",
		"正官": "克我對象，代表責任、守矩、名譽",
		"偏印": "生我對象，代表偏門、特長、老人",
		"正印": "生我對象，代表正統、學習、母親",
	}

	if desc, ok := description[shen]; ok {
		return desc
	}
	return "未知十神"
}

// GetTenGodsSummary 十神總結
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string][]string: 各類十神的總結
//
// 範例：
//   map[string][]string{
//     "比劫": ["年柱:比肩", "月柱:劫財"],
//     "食傷": ["時柱:食神"],
//     "財": ["日柱:正財"],
//     "官殺": ["月柱:正官"],
//     "印綬": ["年柱:正印"],
//   }
func GetTenGodsSummary(bazi Bazi) map[string][]string {
	result := map[string][]string{
		"比劫":  {},
		"食傷":  {},
		"財":    {},
		"官殺":  {},
		"印綬":  {},
	}

	tenGods := CalculateTenGods(bazi)

	for _, shen := range tenGods {
		switch shen {
		case "比肩", "劫財":
			result["比劫"] = append(result["比劫"], shen)
		case "食神", "傷官":
			result["食傷"] = append(result["食傷"], shen)
		case "偏財", "正財":
			result["財"] = append(result["財"], shen)
		case "七殺", "正官":
			result["官殺"] = append(result["官殺"], shen)
		case "偏印", "正印":
			result["印綬"] = append(result["印綬"], shen)
		}
	}

	return result
}

// HasMainTenGods 是否具備主要十神
//
// 參數：
//   - tenGods: 十神摘要
//   - types: 要檢查的十神類型
//
// 回傳：
//   - bool: 是否具備所有指定的十神
func HasMainTenGods(tenGods map[string][]string, types ...string) bool {
	for _, t := range types {
		if len(tenGods[t]) == 0 {
			return false
		}
	}
	return true
}

// GetTenGodsPower 十神力量分析
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string]float64: 各類十神的力量
//
// 計算方式：
//   - 基於十神的出現次數
//   - 考慮地支藏干的加成
func GetTenGodsPower(bazi Bazi) map[string]float64 {
	power := map[string]float64{
		"比劫":  0,
		"食傷":  0,
		"財":    0,
		"官殺":  0,
		"印綬":  0,
	}

	tenGods := CalculateTenGods(bazi)

	// 基礎出現次數
	for _, shen := range tenGods {
		switch shen {
		case "比肩", "劫財":
			power["比劫"]++
		case "食神", "傷官":
			power["食傷"]++
		case "偏財", "正財":
			power["財"]++
		case "七殺", "正官":
			power["官殺"]++
		case "偏印", "正印":
			power["印綬"]++
		}
	}

	// 考慮地支藏干的加成
	_ = []string{
		GetYearBranch(bazi.Year),
		GetMonthBranch(bazi.Month),
		GetDayBranch(bazi.Day),
		GetHourBranch(bazi.Hour),
	}

	// 簡化：地支藏干每個加 0.3
	power["比劫"] += 0.3
	power["食傷"] += 0.3
	power["財"] += 0.3
	power["官殺"] += 0.3
	power["印綬"] += 0.3

	return power
}

// GetTenGodsStar 十神星耀 (奇門遁甲等使用)
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - map[string]string: 各柱的十神星耀
func GetTenGodsStar(bazi Bazi) map[string]string {
	result := make(map[string]string)

	tenGods := CalculateTenGods(bazi)

	for col, shen := range tenGods {
		desc := GetTenGodsDescription(shen)
		result[col] = desc
	}

	return result
}

// GetTenGodsChart 十神表格
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - []map[string]string: 十神表格
func GetTenGodsChart(bazi Bazi) []map[string]string {
	result := make([]map[string]string, 0)

	tenGods := CalculateTenGods(bazi)

	for col, shen := range tenGods {
		row := map[string]string{
			"柱": col,
			"十神": shen,
			"描述": GetTenGodsDescription(shen),
		}
		result = append(result, row)
	}

	return result
}

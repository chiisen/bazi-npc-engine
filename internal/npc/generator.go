// Package npc NPC 設定生成器模組
//
// 功能：
//   - 生成完整 NPC 設定
//   - 姓名、職業、背景故事
//   - 重要人生事件
//
// 使用範例：
//   npc := GenerateNPC(personality)
package npc

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/chiisen/bazi-npc-engine/internal/personality"
)

// ═══════════════════════════════════════════
// 💡 概念：NPC 生成邏輯
// 說明：根據人格特徵生成完整的 NPC 設定
// 為何使用：提供 RPG 遊戲中可使用的 NPC 資料
// 注意事項：使用種子確保可再現性
// ═══════════════════════════════════════════

// NPCProfile NPC 設定結構
type NPCProfile struct {
	Name        string   // 角色姓名
	Age         int      // 年齡
	Occupation  string   // 職業
	Personality []string // 總結描述
	Background  string   // 背景故事
	LifeEvents  []string // 重要事件
}

// 姓名库
var FirstNames = []string{
	"文", "武", "志", "強", "明", "華", "民", "建", "國",
	"偉", "芳", "娜", "秀英", "俊", "杰", "瀾", "浩", "宇", "軒",
	"晨", "辰", "子", "之間", "天", "地", "玄", "黃", "宇宙", "洪荒",
}

var LastNames = []string{
	"張", "王", "李", "趙", "劉", "陳", "楊", "黃", "周", "吳",
	"徐", "孫", "胡", "朱", "高", "林", "何", "郭", "馬", "羅",
	"梁", "宋", "謝", "韓", "唐", "馮", "于", "董", "蕭", "程",
}

// 職業列表
var Occupations = []string{
	"商人", "農夫", "工匠", "書生", "武師", "醫師", "道士", "和尚",
	"店小二", "船夫", "鐵匠", "織工", "書吏", "捕快", "賭徒", "俠客",
	"漁夫", "獵人", "畫師", "琴師", "書畫家", "商贩", "藥商",
	"酒保", "教書先生", "算命師", "風水師", "裁縫", "木匠", "車夫",
}

// 背景故事模板
var Backgrounds = []string{
	"出生於/%s家庭，自幼{%s}",
	"成長於{%s}的環境，{%s}",
	"原本是{%s}，因{%s}而改變人生",
	"来自{%s}，{%s}",
}

// 事件模板
var Events = []string{
	"十歲那年，{%s}",
	"二十歲時，{%s}",
	"一次意外中，{%s}",
	"遇到了{%s}，影響了我一生",
	"在{%s}的經歷，讓我{%s}",
}

// Generate 生成 NPC
//
// 參數：
//   - pers: 人格特徵
//   - seed: 隨機種子 (0 表示使用當前時間)
//
// 回傳：
//   - *NPCProfile: NPC 設定
func Generate(pers *personality.Personality, seed int) *NPCProfile {
	// 使用種子確保可再現性
	if seed == 0 {
		seed = int(time.Now().Unix())
	}
	rand.Seed(int64(seed))

	// 生成姓名
	name := generateName()

	// 計算年齡 (根据日柱簡化計算)
	age := generateAge()

	// 選擇職業
	occupation := generateOccupation()

	// 生成背景故事
	background := generateBackground(name)

	// 生成重要事件
	lifeEvents := generateLifeEvents()

	// 總結描述
	personalityDesc := []string{}
	if len(pers.Traits) > 0 {
		personalityDesc = append(personalityDesc, fmt.Sprintf("性格%s", pers.Traits[0]))
	}
	if len(pers.Strengths) > 0 {
		personalityDesc = append(personalityDesc, fmt.Sprintf("優點：%s", pers.Strengths[0]))
	}
	if len(pers.Weaknesses) > 0 {
		personalityDesc = append(personalityDesc, fmt.Sprintf("缺點：%s", pers.Weaknesses[0]))
	}

	return &NPCProfile{
		Name:        name,
		Age:         age,
		Occupation:  occupation,
		Personality: personalityDesc,
		Background:  background,
		LifeEvents:  lifeEvents,
	}
}

// generateName 生成姓名
func generateName() string {
	firstName := FirstNames[rand.Intn(len(FirstNames))]
	lastName := LastNames[rand.Intn(len(LastNames))]
	return lastName + firstName
}

// generateAge 生成年齡
func generateAge() int {
	// 簡化：根據當前年份與某個年份做差
	currentYear := time.Now().Year()
	year := currentYear - rand.Intn(40) - 20 // 20-60 歲
	return currentYear - year
}

// generateOccupation 生成職業
func generateOccupation() string {
	return Occupations[rand.Intn(len(Occupations))]
}

// generateBackground 生成背景故事
func generateBackground(name string) string {
	parts := []string{
		"家境普通", "家境富裕", "家境貧寒", "學習優異", "天資聰穎", "體弱多病",
		"性格活潑", "性格沉穩", "熱情好客", "沉默寡言", "勤勞刻苦", "聰明伶俐",
	}
	parts2 := []string{
		"立志成為一名優秀的%s", "尋找人生的意義", "追求變化的生活", "尋找失散的家人",
		"報仇雪恨", "尋求長生之道", "為國雙民", "開創一番事業",
	}

	p1 := parts[rand.Intn(len(parts))]
	p2 := parts2[rand.Intn(len(parts2))]

	return fmt.Sprintf(Backgrounds[rand.Intn(len(Backgrounds))], p1, p2)
}

// generateLifeEvents 生成重要事件
func generateLifeEvents() []string {
	templates := []string{
		"遇見了人生的關鍵導師",
		"經歷了一場重大變故",
		"遇到了生命中的摯愛",
		"完成了一件重要事業",
		"學會了珍貴的技能",
		"看透了人生的真相",
		"加入了重要的組織",
		"經歷了崢嶸歲月",
	}

	events := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		events = append(events, templates[rand.Intn(len(templates))])
	}

	return events
}

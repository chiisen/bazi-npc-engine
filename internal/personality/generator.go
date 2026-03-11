// Package personality 人格生成器模組
//
// 功能：
//   - 將八字轉換為人格特徵
//   - 分析優點與缺點
//   - 生成行為模式
//
// 使用範例：
//   pers := Generate(bazi)
package personality

import (
	"fmt"

	"github.com/chiisen/bazi-npc-engine/internal/bazi"
)

// ═══════════════════════════════════════════
// 💡 概念：人格生成邏輯
// 說明：根據八字的五行與十神分析生成人格特徵
// 為何使用：提供可解釋的 NPC 人格設計
// 注意事項：這是一個簡化的模型，實際應用可加入更多規則
// ═══════════════════════════════════════════

// Personality 人格特徵結構
type Personality struct {
	Traits     []string
	Strengths  []string
	Weaknesses []string
	Behavior   []string
}

// Generate 生成人格特徵
//
// 參數：
//   - bazi: 八字盤
//
// 回傳：
//   - Personality: 人格特徵
func Generate(baziData bazi.Bazi) *Personality {
	personality := &Personality{
		Traits:     []string{},
		Strengths:  []string{},
		Weaknesses: []string{},
		Behavior:   []string{},
	}

	// 分析日主 (日柱天干)
	dayStem := bazi.GetDayStem(baziData.Day)
	_ = bazi.GetElement(dayStem) // unused

	// 分析五行平衡
	elems := bazi.CalculateElements(baziData)

	// 根據五行平衡生成特徵
	if elems["木"] > 0.3 {
		personality.Traits = append(personality.Traits, "富有遠見")
		personality.Strengths = append(personality.Strengths, "創造力強")
		personality.Behavior = append(personality.Behavior, "喜歡探索新事物")
	}

	if elems["火"] > 0.3 {
		personality.Traits = append(personality.Traits, "充滿熱情")
		personality.Strengths = append(personality.Strengths, "領導能力強")
		personality.Behavior = append(personality.Behavior, "善於溝通表達")
	}

	if elems["土"] > 0.3 {
		personality.Traits = append(personality.Traits, "踏實穩重")
		personality.Strengths = append(personality.Strengths, "責任感強")
		personality.Behavior = append(personality.Behavior, "長於規劃與執行")
	}

	if elems["金"] > 0.3 {
		personality.Traits = append(personality.Traits, "果断堅毅")
		personality.Strengths = append(personality.Strengths, "邏輯思維清晰")
		personality.Behavior = append(personality.Behavior, "講求效率與規則")
	}

	if elems["水"] > 0.3 {
		personality.Traits = append(personality.Traits, "機智靈活")
		personality.Strengths = append(personality.Strengths, "適應力強")
		personality.Behavior = append(personality.Behavior, "善於變通與應對")
	}

	// 分析十神
	tenGods := bazi.CalculateTenGods(baziData)

	// 正官多：重視規則與名譽
	if hasTenGod(tenGods, "正官") {
		personality.Strengths = append(personality.Strengths, "守規矩")
		personality.Behavior = append(personality.Behavior, "重視社會地位")
	}

	// 七殺多：具有進取心
	if hasTenGod(tenGods, "七殺") {
		personality.Traits = append(personality.Traits, "野心勃勃")
		personality.Strengths = append(personality.Strengths, "敢於挑戰")
	}

	// 正財多：注重現實利益
	if hasTenGod(tenGods, "正財") {
		personality.Traits = append(personality.Traits, "務實節儉")
		personality.Strengths = append(personality.Strengths, "善於理財")
	}

	// 偏財多：善於投機
	if hasTenGod(tenGods, "偏財") {
		personality.Traits = append(personality.Traits, "會號稱")
		personality.Strengths = append(personality.Strengths, "把握機會")
	}

	// 印綬多：重視學習
	if hasTenGod(tenGods, "正印") || hasTenGod(tenGods, "偏印") {
		personality.Strengths = append(personality.Strengths, "好學 not")
		personality.Behavior = append(personality.Behavior, "喜歡思考")
	}

	// 分析日主強弱
	if bazi.IsStrong(dayStem, baziData) {
		personality.Strengths = append(personality.Strengths, "自信滿滿")
		personality.Weaknesses = append(personality.Weaknesses, "可能過於自我")
	} else {
		personality.Strengths = append(personality.Strengths, "謙虛謹慎")
		personality.Weaknesses = append(personality.Weaknesses, "可能缺乏自信")
	}

	// 確保至少有一些特徵
	if len(personality.Traits) == 0 {
		personality.Traits = append(personality.Traits, "獨特")
	}
	if len(personality.Strengths) == 0 {
		personality.Strengths = append(personality.Strengths, "誠實")
	}
	if len(personality.Weaknesses) == 0 {
		personality.Weaknesses = append(personality.Weaknesses, "急躁")
	}
	if len(personality.Behavior) == 0 {
		personality.Behavior = append(personality.Behavior, "謹慎")
	}

	return personality
}

// hasTenGod 檢查是否具有某個十神
func hasTenGod(tenGods map[string]string, shen string) bool {
	for _, s := range tenGods {
		if s == shen {
			return true
		}
	}
	return false
}

// ToDescription 將人格轉換為描述字串
func (p *Personality) ToDescription() string {
	desc := "這是一個"

	// 組合特質
	if len(p.Traits) > 0 {
		desc += fmt.Sprintf("、%s的", p.Traits[0])
	}
	desc += "的人。"

	return desc
}

// GetTraitsString 獲取特質字串
func (p *Personality) GetTraitsString() string {
	if len(p.Traits) == 0 {
		return "未知"
	}
	return fmt.Sprintf("%s", p.Traits)
}

// GetStrengthsString 獲取優點字串
func (p *Personality) GetStrengthsString() string {
	if len(p.Strengths) == 0 {
		return "未知"
	}
	return fmt.Sprintf("%s", p.Strengths)
}

// GetWeaknessesString 獲取缺點字串
func (p *Personality) GetWeaknessesString() string {
	if len(p.Weaknesses) == 0 {
		return "未知"
	}
	return fmt.Sprintf("%s", p.Weaknesses)
}

// Package llm LLM Persona 生成器模組
//
// 功能：
//   - 生成 LLM 可使用的 Prompt
//   - 支援多種 LLM 提供者
//
// 使用範例：
//   prompt := BuildPrompt(npc)
package llm

import (
	"fmt"
	"strings"

	"github.com/chiisen/bazi-npc-engine/internal/npc"
)

// ═══════════════════════════════════════════
// 💡 概念：LLM Prompt 生成
// 說明：將 NPC 設定轉換為 LLM 可使用的 Prompt
// 為何使用：讓 LLM 能夠精準扮演指定角色
// 注意事項：Prompt 應包含完整的角色背景與行為指引
// ═══════════════════════════════════════════

// BuildPrompt 建構 LLM Prompt
//
// 參數：
//   - npc: NPC 設定
//
// 回傳：
//   - string: Prompt 文本
func BuildPrompt(n *npc.NPCProfile) string {
	var sb strings.Builder

	sb.WriteString("You are an NPC with the following characteristics:\n\n")

	// 基本資訊
	sb.WriteString("Name: ")
	sb.WriteString(n.Name)
	sb.WriteString("\n")

	sb.WriteString("Age: ")
	sb.WriteString(fmt.Sprintf("%d", n.Age))
	sb.WriteString("\n")

	sb.WriteString("Occupation: ")
	sb.WriteString(n.Occupation)
	sb.WriteString("\n\n")

	// 人格特質
	sb.WriteString("Personality:\n")
	for _, trait := range n.Personality {
		sb.WriteString("- ")
		sb.WriteString(trait)
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// 背景故事
	sb.WriteString("Background:\n")
	sb.WriteString(n.Background)
	sb.WriteString("\n\n")

	// 重要事件
	sb.WriteString("Life Events:\n")
	for i, event := range n.LifeEvents {
		sb.WriteString("- ")
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, event))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// 行為指引
	sb.WriteString("Behavior Guidelines:\n")
	sb.WriteString("- Respond naturally and consistently with your personality\n")
	sb.WriteString("- Refer to your background when appropriate\n")
	sb.WriteString("- Maintain historical and cultural authenticity\n")
	sb.WriteString("- Keep responses concise and focused\n\n")

	sb.WriteString("Please respond as this character.")

	return sb.String()
}

// BuildSystemPrompt 建構系統 Prompt
// 用於設定模型的角色扮演模式
//
// 參數：
//   - npc: NPC 設定
//
// 回傳：
//   - string: 系統 Prompt 文本
func BuildSystemPrompt(n *npc.NPCProfile) string {
	return fmt.Sprintf(`
You are role-playing as %s, a %d-year-old %s.
%s
Your background: %s

Respond in character as %s. Be authentic to their personality and experiences.
Maintain consistency with their background and life events. Keep responses
concise but expressive.

Key traits: %s
`,
		n.Name,
		n.Age,
		n.Occupation,
		strings.Join(n.Personality, ", "),
		n.Background,
		n.Name,
		strings.Join(n.Personality, ", "),
	)
}

// BuildScenePrompt 建構情境 Prompt
// 用於特定互動場景
//
// 參數：
//   - npc: NPC 設定
//   - scene: 情境描述
//   - user: 用戶描述
//
// 回傳：
//   - string: 情境 Prompt 文本
func BuildScenePrompt(n *npc.NPCProfile, scene, user string) string {
	return fmt.Sprintf(`
You are %s, a %d-year-old %s.

Background: %s
Personality: %s

Current Scene: %s
User: %s

Respond naturally as %s would in this situation.
`,
		n.Name,
		n.Age,
		n.Occupation,
		n.Background,
		strings.Join(n.Personality, ", "),
		scene,
		user,
		n.Name,
	)
}

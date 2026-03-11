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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chiisen/bazi-npc-engine/internal/npc"
)

// ═══════════════════════════════════════════
// 💡 概念：策略模式 (Strategy Pattern)
// 說明：定義介面，讓不同 LLM 提供者可相互替換
// 為何使用：未來可更換 LLM 提供者而不影響使用者
// 注意事項：所有 LLM Client 必須實作 LLMClient 介面
// ═══════════════════════════════════════════

// LLMClient LLM 用戶端介面
type LLMClient interface {
	// Generate 生成回應
	Generate(prompt string) (string, error)
}

// ═══════════════════════════════════════════
// 💡 概念：HTTP Client
// 說明：使用 Go 標準庫 net/http 進行 API 呼叫
// 為何使用：無外部依賴，穩定可靠
// 注意事項：設定合理的 timeout 避免長時間等待
// ═══════════════════════════════════════════

const (
	defaultTimeout = 30 * time.Second
)

// HTTPClient HTTP API 用戶端
type HTTPClient struct {
	BaseURL   string
	APIKey    string
	Model     string
	Timeout   time.Duration
	httpClient *http.Client
}

// APIRequest API 請求結構
type APIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message 訊息結構
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// APIResponse API 回應結構
type APIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 選擇結構
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用量結構
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewHTTPClient 建立新的 HTTP 用戶端
//
// 參數：
//   - baseURL: API 基礎 URL
//   - apiKey: API 金鑰
//   - model: 模型名稱
//
// 回傳：
//   - *HTTPClient: HTTP 用戶端
func NewHTTPClient(baseURL, apiKey, model string) *HTTPClient {
	return &HTTPClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Model:      model,
		Timeout:    defaultTimeout,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
}

// Generate 生成回應
//
// 參數：
//   - prompt: Prompt 文本
//
// 回傳：
//   - string: 生成的回應
//   - error: 錯誤訊息
func (c *HTTPClient) Generate(prompt string) (string, error) {
	// 構建請求
	messages := []Message{
		{Role: "system", Content: "You are a helpful NPC assistant."},
		{Role: "user", Content: prompt},
	}

	reqBody := APIRequest{
		Model:       c.Model,
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化請求失敗: %v", err)
	}

	// 建立 HTTP 請求
	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("建立請求失敗: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	// 發送請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("API 呼叫失敗: %v", err)
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("讀取回應失敗: %v", err)
	}

	// 檢查狀態碼
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 回應錯誤 (%d): %s", resp.StatusCode, string(body))
	}

	// 解析回應
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("解析回應失敗: %v", err)
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("API 回應無選擇")
	}

	return apiResp.Choices[0].Message.Content, nil
}

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

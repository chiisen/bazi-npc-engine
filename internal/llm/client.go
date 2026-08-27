// Package llm LLM Persona 生成器模組
//
// 功能：
//   - 生成 LLM 可使用的 Prompt
//   - 透過 LLMClient 介面支援多種 LLM 提供者
//
// 使用範例：
//   client, _ := NewClient(ClientConfig{Provider: "openai", APIKey: "sk-..."})
//   resp, _ := client.Generate(prompt)
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

// ═══════════════════════════════════════════
// 概念：策略模式 (Strategy Pattern)
// 說明：定義介面，讓不同 LLM 提供者可替換
// 為何使用：未來加入非 OpenAI 相容 provider 不影響呼叫端
// ═══════════════════════════════════════════

// LLMClient LLM 用戶端介面
type LLMClient interface {
	Generate(prompt string) (string, error)
}

// ═══════════════════════════════════════════
// 概念：OpenAI 相容 Client
// 說明：處理所有 OpenAI Chat Completions 相容格式的 API
// 適用：OpenAI、OpenCode Go、MiniMax、其他相容服務
// ═══════════════════════════════════════════

// APIRequest OpenAI 相容 API 請求結構
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

// APIResponse OpenAI 相容 API 回應結構
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

// OpenAIClient OpenAI 相容 API 用戶端
type OpenAIClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

// NewOpenAIClient 建立 OpenAI 相容用戶端
//
// 參數：
//   - baseURL: API 基礎 URL（含版本路徑，例如 https://api.openai.com/v1）
//   - apiKey: API 金鑰
//   - model: 模型名稱
func NewOpenAIClient(baseURL, apiKey, model string) *OpenAIClient {
	return &OpenAIClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Model:      model,
		HTTPClient: &http.Client{Timeout: defaultTimeout},
	}
}

// Generate 呼叫 LLM API 生成回應
//
// 參數：
//   - prompt: Prompt 文本
//
// 回傳：
//   - string: 模型回傳的第一個 choice 內容
//   - error: 呼叫失敗時的錯誤（含位置與技術細節）
func (c *OpenAIClient) Generate(prompt string) (string, error) {
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
		return "", fmt.Errorf("序列化請求失敗（位置: client.go Generate）: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("建立請求失敗（位置: client.go Generate）: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("API 呼叫失敗（位置: client.go Generate）: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("讀取回應失敗（位置: client.go Generate）: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		snippet := string(body)
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}
		return "", fmt.Errorf("API 回應錯誤 %d（位置: client.go Generate）: %s", resp.StatusCode, snippet)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("解析回應失敗（位置: client.go Generate）: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("API 回應無選擇（位置: client.go Generate）")
	}

	return apiResp.Choices[0].Message.Content, nil
}
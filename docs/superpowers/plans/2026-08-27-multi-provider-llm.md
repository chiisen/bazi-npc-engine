# 多 LLM Provider 支援 實作計畫

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 讓 `npcgen` CLI 支援 OpenAI、OpenCode Go、MiniMax 三個 LLM Provider，使用者可透過 `--provider` 旗標切換。

**Architecture:** 保留 `LLMClient` 介面作為未來擴充點。將現有 `HTTPClient` 重新命名為 `OpenAIClient`，所有 OpenAI 相容格式由單一實作處理。新增 `Provider` 預設表（`providers.go`）與 Factory 函式（`factory.go`）依 provider 名稱切換 BaseURL/API Key/Model。

**Tech Stack:** Go 1.23、Go standard `net/http`、`httptest` 測試。

**Spec:** `docs/README.md`（2026-08-27-multi-provider-llm-design.md）

---

## 檔案結構總覽

| 動作 | 檔案 | 職責 |
|------|------|------|
| 新增 | `internal/llm/providers.go` | `Provider` 結構 + `DefaultProviders` 預設表 |
| 新增 | `internal/llm/providers_test.go` | 預設表完整性測試 |
| 修改 | `internal/llm/generator.go` | 移除 HTTPClient，僅保留 Prompt 建構函式 |
| 新增 | `internal/llm/client.go` | `LLMClient` 介面 + `OpenAIClient` 實作（含 APIRequest/Message/APIResponse/Choice/Usage 結構） |
| 新增 | `internal/llm/client_test.go` | `OpenAIClient.Generate` 錯誤路徑測試 |
| 新增 | `internal/llm/factory.go` | `ClientConfig` + `NewClient` Factory 函式 |
| 新增 | `internal/llm/factory_test.go` | `NewClient` 各路徑測試 |
| 修改 | `cmd/npcgen/main.go` | 新增 `--provider`/`--model`/`--api-key` 旗標 + 整合 factory |
| 修改 | `CHANGELOG.md` | 記錄 v0.2.0 新功能 |

---

## Task 1: 建立 Provider 預設表

**Files:**
- Create: `internal/llm/providers.go`
- Create: `internal/llm/providers_test.go`

- [ ] **Step 1: 寫失敗測試**

`internal/llm/providers_test.go`：

```go
package llm

import (
	"strings"
	"testing"
)

func TestDefaultProviders_完整性(t *testing.T) {
	required := []string{"openai", "opencode-go", "MiniMax"}
	for _, name := range required {
		p, ok := DefaultProviders[name]
		if !ok {
			t.Errorf("缺少 provider: %s", name)
			continue
		}
		if p.Name == "" {
			t.Errorf("%s: Name 為空", name)
		}
		if p.BaseURL == "" {
			t.Errorf("%s: BaseURL 為空", name)
		}
		if p.DefaultModel == "" {
			t.Errorf("%s: DefaultModel 為空", name)
		}
		if p.EnvVar == "" {
			t.Errorf("%s: EnvVar 為空", name)
		}
	}
}

func TestDefaultProviders_BaseURL_必須HTTPS(t *testing.T) {
	for name, p := range DefaultProviders {
		if !strings.HasPrefix(p.BaseURL, "https://") {
			t.Errorf("%s: BaseURL 必須以 https:// 開頭，實際: %s", name, p.BaseURL)
		}
	}
}
```

- [ ] **Step 2: 執行測試確認失敗**

Run: `go test ./internal/llm/ -run TestDefaultProviders -v`
Expected: FAIL（`providers.go` 不存在，編譯失敗）

- [ ] **Step 3: 實作 Provider 結構與預設表**

`internal/llm/providers.go`：

```go
// Package llm LLM Provider 預設設定
package llm

// Provider LLM 提供者預設設定
type Provider struct {
	Name         string
	BaseURL      string
	DefaultModel string
	EnvVar       string
}

// DefaultProviders 內建支援的 LLM Provider 預設表
var DefaultProviders = map[string]Provider{
	"openai": {
		Name:         "openai",
		BaseURL:      "https://api.openai.com/v1",
		DefaultModel: "gpt-4o-mini",
		EnvVar:       "OPENAI_API_KEY",
	},
	"opencode-go": {
		Name:         "opencode-go",
		BaseURL:      "https://opencode.ai/zen/go/v1",
		DefaultModel: "claude-sonnet-4-5",
		EnvVar:       "OPENCODE_API_KEY",
	},
	"MiniMax": {
		Name:         "MiniMax",
		BaseURL:      "https://api.minimax.io/v1",
		DefaultModel: "MiniMax-M2",
		EnvVar:       "MiniMax_API_KEY",
	},
}
```

- [ ] **Step 4: 執行測試確認通過**

Run: `go test ./internal/llm/ -run TestDefaultProviders -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/llm/providers.go internal/llm/providers_test.go
git commit -m "✨ feat(llm): 新增 Provider 預設表（OpenAI/OpenCode Go/MiniMax）"
```

---

## Task 2: 重構 HTTPClient → OpenAIClient 並提取至 client.go

**Files:**
- Create: `internal/llm/client.go`
- Modify: `internal/llm/generator.go`（移除 HTTPClient 相關程式碼，保留 Prompt 建構函式）

> 注意：本任務為純重構（不改行為），執行後跑既有測試確認沒破壞。

- [ ] **Step 1: 確認既有測試通過**

Run: `go test ./... -v 2>&1 | tail -20`
Expected: PASS（目前無 llm 測試，至少其他模組測試都通過）

- [ ] **Step 2: 建立新 client.go 檔案**

`internal/llm/client.go`：

```go
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
```

- [ ] **Step 3: 從 generator.go 移除重複的 HTTPClient 程式碼**

修改 `internal/llm/generator.go`，**只保留** Prompt 建構函式（`BuildPrompt`、`BuildSystemPrompt`、`BuildScenePrompt`）。刪除以下內容：

- `import` 中的 `"bytes"`、`"encoding/json"`、`"io"`、`"net/http"`、`"time"`
- `defaultTimeout` 常數
- `HTTPClient` 結構
- `APIRequest`、`Message`、`APIResponse`、`Choice`、`Usage` 結構
- `NewHTTPClient` 函式
- `HTTPClient.Generate` 方法

`generator.go` 最終內容：

```go
// Package llm LLM Persona 生成器模組
//
// 功能：
//   - 生成 LLM 可使用的 Prompt
//   - 透過 LLMClient 介面支援多種 LLM 提供者
//
// 使用範例：
//   prompt := llm.BuildPrompt(npcProfile)
package llm

import (
	"fmt"
	"strings"

	"github.com/chiisen/bazi-npc-engine/internal/npc"
)

// ═══════════════════════════════════════════
// 概念：LLM Prompt 生成
// 說明：將 NPC 設定轉換為 LLM 可使用的 Prompt
// 為何使用：讓 LLM 能夠精準扮演指定角色
// ═══════════════════════════════════════════

// BuildPrompt 建構 LLM Prompt
//
// 參數：
//   - n: NPC 設定
//
// 回傳：
//   - string: Prompt 文本
func BuildPrompt(n *npc.NPCProfile) string {
	var sb strings.Builder

	sb.WriteString("You are an NPC with the following characteristics:\n\n")

	sb.WriteString("Name: ")
	sb.WriteString(n.Name)
	sb.WriteString("\n")

	sb.WriteString("Age: ")
	sb.WriteString(fmt.Sprintf("%d", n.Age))
	sb.WriteString("\n")

	sb.WriteString("Occupation: ")
	sb.WriteString(n.Occupation)
	sb.WriteString("\n\n")

	sb.WriteString("Personality:\n")
	for _, trait := range n.Personality {
		sb.WriteString("- ")
		sb.WriteString(trait)
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	sb.WriteString("Background:\n")
	sb.WriteString(n.Background)
	sb.WriteString("\n\n")

	sb.WriteString("Life Events:\n")
	for i, event := range n.LifeEvents {
		sb.WriteString("- ")
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, event))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	sb.WriteString("Behavior Guidelines:\n")
	sb.WriteString("- Respond naturally and consistently with your personality\n")
	sb.WriteString("- Refer to your background when appropriate\n")
	sb.WriteString("- Maintain historical and cultural authenticity\n")
	sb.WriteString("- Keep responses concise and focused\n\n")

	sb.WriteString("Please respond as this character.")

	return sb.String()
}

// BuildSystemPrompt 建構系統 Prompt
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
```

- [ ] **Step 4: 編譯確認**

Run: `go build ./...`
Expected: 無錯誤

- [ ] **Step 5: 執行既有測試確認未破壞**

Run: `go test ./... 2>&1 | tail -10`
Expected: PASS（所有既有測試仍通過）

- [ ] **Step 6: Commit**

```bash
git add internal/llm/client.go internal/llm/generator.go
git commit -m "♻️ refactor(llm): 拆分 HTTPClient 至 client.go 並改名為 OpenAIClient"
```

---

## Task 3: 為 OpenAIClient.Generate 新增錯誤路徑測試

**Files:**
- Create: `internal/llm/client_test.go`

- [ ] **Step 1: 寫失敗測試（5 個案例）**

`internal/llm/client_test.go`：

```go
package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 驗證 helper：解析請求 body 為 APIRequest
func decodeRequest(t *testing.T, r *http.Request) APIRequest {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("讀取請求 body 失敗: %v", err)
	}
	var req APIRequest
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("解析請求 body 失敗: %v", err)
	}
	return req
}

func TestOpenAIClient_Generate_正常回應(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"choices":[{"index":0,"message":{"role":"assistant","content":"hi"}}]}`))
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key", "test-model")
	out, err := client.Generate("hello")
	if err != nil {
		t.Fatalf("預期無錯誤，實際: %v", err)
	}
	if out != "hi" {
		t.Errorf("預期 'hi'，實際: %s", out)
	}
}

func TestOpenAIClient_Generate_非200狀態(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"invalid api key"}}`))
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "bad-key", "model")
	_, err := client.Generate("hi")
	if err == nil {
		t.Fatal("預期錯誤，實際 nil")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("錯誤訊息應含 401，實際: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "invalid api key") {
		t.Errorf("錯誤訊息應含 body 片段，實際: %s", err.Error())
	}
}

func TestOpenAIClient_Generate_無Choices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[]}`))
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "k", "m")
	_, err := client.Generate("hi")
	if err == nil {
		t.Fatal("預期錯誤，實際 nil")
	}
	if !strings.Contains(err.Error(), "API 回應無選擇") {
		t.Errorf("錯誤訊息應含 'API 回應無選擇'，實際: %s", err.Error())
	}
}

func TestOpenAIClient_Generate_JSON解析失敗(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{not valid json`))
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "k", "m")
	_, err := client.Generate("hi")
	if err == nil {
		t.Fatal("預期錯誤，實際 nil")
	}
	if !strings.Contains(err.Error(), "解析回應失敗") {
		t.Errorf("錯誤訊息應含 '解析回應失敗'，實際: %s", err.Error())
	}
}

func TestOpenAIClient_Generate_請求格式正確(t *testing.T) {
	var capturedReq APIRequest
	var capturedAuth string
	var capturedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = decodeRequest(t, r)
		capturedAuth = r.Header.Get("Authorization")
		capturedPath = r.URL.Path
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "secret-key", "gpt-test")
	_, err := client.Generate("test prompt")
	if err != nil {
		t.Fatalf("預期無錯誤，實際: %v", err)
	}

	if capturedPath != "/chat/completions" {
		t.Errorf("預期路徑 /chat/completions，實際: %s", capturedPath)
	}
	if capturedAuth != "Bearer secret-key" {
		t.Errorf("預期 Authorization 'Bearer secret-key'，實際: %s", capturedAuth)
	}
	if capturedReq.Model != "gpt-test" {
		t.Errorf("預期 model 'gpt-test'，實際: %s", capturedReq.Model)
	}
	if len(capturedReq.Messages) != 2 {
		t.Errorf("預期 2 個 messages，實際: %d", len(capturedReq.Messages))
	}
	if capturedReq.Messages[0].Role != "system" {
		t.Errorf("預期第一則 role=system，實際: %s", capturedReq.Messages[0].Role)
	}
	if capturedReq.Messages[1].Role != "user" || capturedReq.Messages[1].Content != "test prompt" {
		t.Errorf("預期第二則 user='test prompt'，實際: %+v", capturedReq.Messages[1])
	}
}
```

- [ ] **Step 2: 執行測試確認全部通過**

Run: `go test ./internal/llm/ -run TestOpenAIClient -v`
Expected: 5 個測試全部 PASS（Task 2 已實作完整 Generate 邏輯）

> 若失敗：依錯誤訊息修正 `client.go` 後重跑。

- [ ] **Step 3: Commit**

```bash
git add internal/llm/client_test.go
git commit -m "✅ test(llm): 為 OpenAIClient.Generate 新增 5 個錯誤路徑測試"
```

---

## Task 4: 實作 Factory 函式 NewClient

**Files:**
- Create: `internal/llm/factory.go`
- Create: `internal/llm/factory_test.go`

- [ ] **Step 1: 寫失敗測試（6 個案例）**

`internal/llm/factory_test.go`：

```go
package llm

import (
	"strings"
	"testing"
	"time"
)

func TestNewClient_已知Provider_帶APIKey(t *testing.T) {
	cfg := ClientConfig{
		Provider: "openai",
		APIKey:   "explicit-key",
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("預期無錯誤，實際: %v", err)
	}
	if client == nil {
		t.Fatal("預期非 nil client")
	}

	// 驗證回傳的是 *OpenAIClient 且設定正確
	oc, ok := client.(*OpenAIClient)
	if !ok {
		t.Fatalf("預期 *OpenAIClient，實際: %T", client)
	}
	if oc.APIKey != "explicit-key" {
		t.Errorf("APIKey 應為 'explicit-key'，實際: %s", oc.APIKey)
	}
	if oc.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("BaseURL 不正確: %s", oc.BaseURL)
	}
}

func TestNewClient_從環境變數讀取APIKey(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "env-key")

	cfg := ClientConfig{Provider: "openai"}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("預期無錯誤，實際: %v", err)
	}
	oc := client.(*OpenAIClient)
	if oc.APIKey != "env-key" {
		t.Errorf("APIKey 應從環境變數讀取 'env-key'，實際: %s", oc.APIKey)
	}
}

func TestNewClient_缺少APIKey_錯誤訊息含EnvVar名稱(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")

	cfg := ClientConfig{Provider: "openai"}
	_, err := NewClient(cfg)
	if err == nil {
		t.Fatal("預期錯誤，實際 nil")
	}
	if !strings.Contains(err.Error(), "OPENAI_API_KEY") {
		t.Errorf("錯誤訊息應含環境變數名稱 'OPENAI_API_KEY'，實際: %s", err.Error())
	}
}

func TestNewClient_未知Provider_錯誤訊息含可用清單(t *testing.T) {
	cfg := ClientConfig{Provider: "unknown", APIKey: "k"}
	_, err := NewClient(cfg)
	if err == nil {
		t.Fatal("預期錯誤，實際 nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "openai") {
		t.Errorf("錯誤訊息應含 'openai'，實際: %s", msg)
	}
	if !strings.Contains(msg, "MiniMax") {
		t.Errorf("錯誤訊息應含 'MiniMax'，實際: %s", msg)
	}
}

func TestNewClient_指定Model_覆寫預設(t *testing.T) {
	cfg := ClientConfig{
		Provider: "openai",
		APIKey:   "k",
		Model:    "gpt-4-turbo",
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("預期無錯誤，實際: %v", err)
	}
	oc := client.(*OpenAIClient)
	if oc.Model != "gpt-4-turbo" {
		t.Errorf("Model 應為 'gpt-4-turbo'，實際: %s", oc.Model)
	}
}

func TestNewClient_未指定Model_使用預設(t *testing.T) {
	cfg := ClientConfig{Provider: "openai", APIKey: "k"}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("預期無錯誤，實際: %v", err)
	}
	oc := client.(*OpenAIClient)
	if oc.Model != "gpt-4o-mini" {
		t.Errorf("Model 應為預設 'gpt-4o-mini'，實際: %s", oc.Model)
	}
}

func TestNewClient_自訂Timeout(t *testing.T) {
	cfg := ClientConfig{
		Provider: "openai",
		APIKey:   "k",
		Timeout:  10 * time.Second,
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("預期無錯誤，實際: %v", err)
	}
	oc := client.(*OpenAIClient)
	if oc.HTTPClient == nil || oc.HTTPClient.Timeout != 10*time.Second {
		t.Errorf("Timeout 應為 10s，實際: %v", oc.HTTPClient.Timeout)
	}
}
```

- [ ] **Step 2: 執行測試確認失敗**

Run: `go test ./internal/llm/ -run TestNewClient -v`
Expected: FAIL（`factory.go` 不存在）

- [ ] **Step 3: 實作 factory.go**

`internal/llm/factory.go`：

```go
package llm

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// ClientConfig Factory 輸入設定
type ClientConfig struct {
	Provider string        // provider 名稱（key of DefaultProviders）
	APIKey   string        // API 金鑰（空則讀環境變數）
	Model    string        // 模型名稱（空則用 provider 預設）
	Timeout  time.Duration // HTTP timeout（0 則用 defaultTimeout）
}

// NewClient 根據設定建立 LLMClient
//
// 解析順序：
//   1. 查 DefaultProviders[cfg.Provider]，找不到回傳錯誤
//   2. cfg.APIKey 空 → 讀 os.Getenv(provider.EnvVar)
//   3. cfg.APIKey 仍空 → 回傳錯誤（訊息含 EnvVar 名稱）
//   4. cfg.Model 空 → 用 provider.DefaultModel
//   5. cfg.Timeout 0 → 用 defaultTimeout
func NewClient(cfg ClientConfig) (LLMClient, error) {
	provider, ok := DefaultProviders[cfg.Provider]
	if !ok {
		available := make([]string, 0, len(DefaultProviders))
		for k := range DefaultProviders {
			available = append(available, k)
		}
		return nil, fmt.Errorf(
			"不支援的 provider %q（位置: factory.go NewClient），可用: %s",
			cfg.Provider,
			strings.Join(available, ", "),
		)
	}

	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv(provider.EnvVar)
	}
	if apiKey == "" {
		return nil, fmt.Errorf(
			"缺少 API key（位置: factory.go NewClient），請設定 --api-key 或環境變數 %s",
			provider.EnvVar,
		)
	}

	model := cfg.Model
	if model == "" {
		model = provider.DefaultModel
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return &OpenAIClient{
		BaseURL:    provider.BaseURL,
		APIKey:     apiKey,
		Model:      model,
		HTTPClient: &http.Client{Timeout: timeout},
	}, nil
}
```

- [ ] **Step 4: 執行測試確認全部通過**

Run: `go test ./internal/llm/ -run TestNewClient -v`
Expected: 7 個測試全部 PASS

- [ ] **Step 5: 執行全部 llm 模組測試**

Run: `go test ./internal/llm/ -v`
Expected: 全部 PASS

- [ ] **Step 6: Commit**

```bash
git add internal/llm/factory.go internal/llm/factory_test.go
git commit -m "✨ feat(llm): 新增 Factory 函式 NewClient 支援多 Provider 切換"
```

---

## Task 5: 更新 CLI 支援 --provider/--model/--api-key

**Files:**
- Modify: `cmd/npcgen/main.go`

- [ ] **Step 1: 確認 CLI 既有編譯與測試通過**

Run: `go build ./... && go test ./cmd/... -v 2>&1 | tail -10`
Expected: 無錯誤（cmd 目前無測試）

- [ ] **Step 2: 修改 CLIOptions 結構加入新欄位**

`cmd/npcgen/main.go` line 19-28 改為：

```go
// CLI 參數結構
type CLIOptions struct {
	Birth     string
	Format    string
	Output    string
	Seed      int
	Verbose   bool
	Help      bool
	Version  bool
	Provider  string
	Model     string
	APIKey    string
}
```

- [ ] **Step 3: 修改 ParseOptions 加入新旗標**

`cmd/npcgen/main.go` line 31-50 改為：

```go
// ParseOptions 解析命令列選項
func ParseOptions() *CLIOptions {
	birth := flag.String("birth", "", "出生時間 (格式: YYYY-MM-DD HH:MM)")
	format := flag.String("format", "text", "輸出格式 (json/text)")
	output := flag.String("output", "", "輸出檔案路徑")
	seed := flag.Int("seed", 0, "隨機種子 (用於再現相同結果)")
	verbose := flag.Bool("verbose", false, "顯示詳細資訊")
	showHelp := flag.Bool("help", false, "顯示帮助訊息")
	showVersion := flag.Bool("version", false, "顯示版本訊息")
	provider := flag.String("provider", "openai", "LLM Provider (openai/opencode-go/MiniMax)")
	model := flag.String("model", "", "LLM 模型名稱（覆寫 provider 預設）")
	apiKey := flag.String("api-key", "", "LLM API 金鑰（測試用，建議用環境變數）")
	flag.Parse()

	return &CLIOptions{
		Birth:   *birth,
		Format:  *format,
		Output:  *output,
		Seed:    *seed,
		Verbose: *verbose,
		Help:    *showHelp,
		Version: *showVersion,
		Provider: *provider,
		Model:    *model,
		APIKey:   *apiKey,
	}
}
```

- [ ] **Step 4: 修改 PrintHelp 加入新旗標說明**

`cmd/npcgen/main.go` line 53-73 改為：

```go
// PrintHelp 顯示幫助訊息
func PrintHelp() {
	fmt.Println(`Bazi NPC Generator - 使用八字生成 RPG NPC

用法：
  npcgen [選項]

選項：
  --birth string       出生時間 (格式: YYYY-MM-DD HH:MM) [必填]
  --format string      輸出格式 (json/text, 預設: text)
  --output string      輸出檔案路徑 (預設: 標準輸出)
  --seed int           隨機種子 (用於再現相同結果)
  --verbose            顯示詳細八字資訊
  --provider string    LLM Provider (openai/opencode-go/MiniMax, 預設: openai)
  --model string       LLM 模型名稱（覆寫 provider 預設）
  --api-key string     LLM API 金鑰（測試用，建議用環境變數）
  --help               顯示幫助訊息
  --version            顯示版本訊息

範例：
  npcgen --birth "1995-10-01 14:00"
  npcgen --birth "1995-10-01 14:00" --format json
  npcgen --birth "1995-10-01 14:00" --output npc.json
  npcgen --birth "1995-10-01 14:00" --verbose
  npcgen --birth "1995-10-01 14:00" --provider MiniMax`)
}
```

- [ ] **Step 5: 在 Run() 加上 API key 明文警告**

在 `cmd/npcgen/main.go` 的 `Run()` 函式（line 211 附近），在處理 Help/Version 後、檢查 Birth 之前，加入：

```go
	// 警告：明文傳入 API key
	if opts.APIKey != "" {
		log.Println("WARN: 不建議在 CLI 傳入金鑰，建議使用環境變數")
	}
```

完整 `Run()` 函式（從 line 211 開始）：

```go
// Run CLI 主流程
func Run() int {
	opts := ParseOptions()

	// 處理 special 選項
	if opts.Help {
		PrintHelp()
		return 0
	}

	if opts.Version {
		PrintVersion()
		return 0
	}

	// 警告：明文傳入 API key
	if opts.APIKey != "" {
		log.Println("WARN: 不建議在 CLI 傳入金鑰，建議使用環境變數")
	}

	// 檢查出生時間是否提供
	if opts.Birth == "" {
		fmt.Println("錯誤：請提供出生時間")
		fmt.Println("使用 --help 查看使用說明")
		return 1
	}

	// 檢查格式是否有效
	if opts.Format != "text" && opts.Format != "json" {
		fmt.Printf("錯誤：無效的格式 '%s'，請使用 'text' 或 'json'\n", opts.Format)
		return 1
	}

	// 生成 NPC
	err := GenerateNPC(opts)
	if err != nil {
		log.Printf("錯誤：生成 NPC 失敗: %v", err)
		return 1
	}

	return 0
}
```

> 注意：本次實作**不**串接 LLM 呼叫，僅新增 CLI 旗標解析與幫助文字。LLM 串接留待後續 PR（避免本次範圍過大）。

- [ ] **Step 6: 編譯確認**

Run: `go build ./...`
Expected: 無錯誤

- [ ] **Step 7: 執行 --help 確認旗標顯示正確**

Run: `go run ./cmd/npcgen/ --help`
Expected: 輸出含 `--provider`, `--model`, `--api-key` 說明

- [ ] **Step 8: 執行既有測試確認未破壞**

Run: `go test ./... 2>&1 | tail -10`
Expected: PASS

- [ ] **Step 9: Commit**

```bash
git add cmd/npcgen/main.go
git commit -m "✨ feat(cli): 新增 --provider/--model/--api-key 旗標"
```

---

## Task 6: 執行全部測試 + 更新 CHANGELOG

**Files:**
- Modify: `CHANGELOG.md`

- [ ] **Step 1: 執行全部單元測試**

Run: `go test ./... -v 2>&1 | tail -30`
Expected: 全部 PASS（含既有 69 個 + 新增 llm 測試 14 個）

- [ ] **Step 2: 執行 go vet**

Run: `go vet ./...`
Expected: 無警告

- [ ] **Step 3: 檢查專案 Makefile 是否有 lint 指令，若有則執行**

Run: `make lint 2>/dev/null || echo "no lint target"`
Expected: 全部通過或無此 target

- [ ] **Step 4: 更新 CHANGELOG.md**

在 CHANGELOG.md 開開加入新版本記錄（依照專案既有格式）：

```markdown
## [Unreleased]

### 新增

- 多 LLM Provider 支援（OpenAI / OpenCode Go / MiniMax）
  - `--provider` 旗標切換 provider
  - `--model` 覆寫預設模型
  - `--api-key` 臨時指定金鑰（建議用環境變數）
  - 工廠函式 `llm.NewClient()` 統一處理設定解析
- `internal/llm/providers.go` 新增 Provider 預設表
- `internal/llm/client.go` 新增 OpenAIClient 實作（從 generator.go 拆分）
- 14 個新單元測試（provider 預設表 2 + client 錯誤路徑 5 + factory 路徑 7）
```

- [ ] **Step 5: Commit CHANGELOG**

```bash
git add CHANGELOG.md
git commit -m "📝 docs: 更新 CHANGELOG 記錄 v0.2.0 多 Provider 功能"
```

- [ ] **Step 6: 最終 git status 確認乾淨**

Run: `git status`
Expected: nothing to commit, working tree clean

---

## 驗收標準

- [ ] 所有任務完成且 commit
- [ ] `go test ./...` 全綠
- [ ] `go vet ./...` 無警告
- [ ] `go run ./cmd/npcgen/ --help` 顯示三個新旗標
- [ ] CHANGELOG.md 已更新
- [ ] `git log --oneline -10` 顯示 6 個新 commit（每個 task 一個）
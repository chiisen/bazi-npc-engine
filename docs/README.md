# 多 LLM Provider 支援設計

**日期**：2026-08-27
**狀態**：待審核
**範圍**：`internal/llm/` 與 `cmd/npcgen/main.go`

## 目標

讓 `npcgen` CLI 支援三個 LLM Provider，使用者可透過旗標切換：

- **OpenAI**（既有）
- **OpenCode Go**（新）
- **MiniMax**（新）

## 背景

目前 `internal/llm/generator.go` 只有單一 `HTTPClient` 實作，採用 OpenAI Chat Completions 格式。經調查三個 Provider 皆為 OpenAI 相容格式（Bearer 認證、`/chat/completions` 端點、相同 JSON 結構），因此可用單一 client 實作 + Provider 預設表滿足需求。

來源：
- OpenCode Go：[官方文件](https://opencode.ai/docs/go/)、[GitHub Issue #668](https://github.com/craft-ai-agents/craft-agents-oss/issues/668)
- MiniMax：[Chat Completions API](https://platform.minimax.io/docs/api-reference/text-chat-openai)、[OpenAI SDK 相容性](https://platform.minimax.io/docs/api-reference/text-openai-api)

## 架構

```
┌─────────────────────────────────────────────┐
│ CLI (cmd/npcgen/main.go)                    │
│   └─ --provider openai|opencode-go|MiniMax │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│ Factory (internal/llm/factory.go)           │
│   └─ NewClient(cfg ClientConfig) -> ... │
└──────────────────┬──────────────────────────┘
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
┌──────────────┐    ┌──────────────────────┐
│ providers.go │    │ client.go (介面+實作)│
│ 預設表       │    │  LLMClient (介面)    │
│  - openai    │    │   └─ OpenAIClient   │
│  - opencode │    │      (所有 OpenAI   │
│  - MiniMax    │    │       相容格式)     │
└──────────────┘    └──────────────────────┘
```

**核心決策**：
- 保留 `LLMClient` 介面（未來擴充點）
- 將現有 `HTTPClient` 重新命名為 `OpenAIClient`（職責更精準）
- 新增 `Provider` 註冊表儲存三組 BaseURL 與預設模型
- Factory 函式根據 provider 名稱查找預設值，加上使用者自訂參數後建立 client

## 元件設計

### `internal/llm/client.go`

```go
// LLMClient 介面（保留）
type LLMClient interface {
    Generate(prompt string) (string, error)
}

// OpenAIClient（從 HTTPClient 改名）
type OpenAIClient struct {
    BaseURL    string
    APIKey     string
    Model      string
    HTTPClient *http.Client
}

// 既有 APIRequest / Message / APIResponse / Choice / Usage 結構不變
func NewOpenAIClient(baseURL, apiKey, model string) *OpenAIClient
func (c *OpenAIClient) Generate(prompt string) (string, error)
```

要點：
- `HTTPClient` 改為可注入欄位（方便測試）
- `NewOpenAIClient` 不再讀取設定，僅設定結構
- 所有 JSON 結構（APIRequest/Message/APIResponse/Choice/Usage）維持不變
- Prompt 建構函式（BuildPrompt/BuildSystemPrompt/BuildScenePrompt）維持不變

### `internal/llm/providers.go`（新增）

```go
type Provider struct {
    Name         string
    BaseURL      string
    DefaultModel string
    EnvVar       string
}

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

### `internal/llm/factory.go`（新增）

```go
type ClientConfig struct {
    Provider string
    APIKey   string
    Model    string
    Timeout  time.Duration
}

func NewClient(cfg ClientConfig) (LLMClient, error) {
    // 1. 查 DefaultProviders[cfg.Provider]
    // 2. 沒找到 → 不支援的 provider 錯誤
    // 3. cfg.APIKey 空 → 讀 os.Getenv(provider.EnvVar)
    // 4. cfg.APIKey 仍空 → 缺少 API key 錯誤
    // 5. cfg.Model 空 → 用 provider.DefaultModel
    // 6. cfg.Timeout 0 → 用 defaultTimeout
    // 7. 回傳 *OpenAIClient
}
```

### CLI 變更（`cmd/npcgen/main.go`）

新增旗標：
- `--provider <name>`：選擇 provider，預設 `"openai"`
- `--model <name>`：覆寫預設模型
- `--api-key <key>`：直接給定金鑰（測試/一次性用）

既有 `--birth`、`--format`、`--output`、`--verbose`、`--seed` 維持不變。

## 資料流

### 完整呼叫鏈

```
1. 使用者下指令 $ npcgen --birth "1995-10-01 14:00" --provider MiniMax --model MiniMax-M2

2. CLI 解析 (cmd/npcgen/main.go)
   ├─ 解析 --birth → time.Time
   ├─ 解析 --provider → "MiniMax"
   ├─ 解析 --model → "MiniMax-M2"
   ├─ 解析 --api-key（若有）→ string
   └─ 建立 ClientConfig{Provider, APIKey, Model, Timeout}

3. Factory (factory.go)
   ├─ 查 DefaultProviders["MiniMax"] → {BaseURL, DefaultModel, EnvVar}
   ├─ APIKey 空 → 讀 os.Getenv("MiniMax_API_KEY")
   ├─ Model 空 → 用 "MiniMax-M2"
   └─ 回傳 *OpenAIClient

4. NPC 生成 (內部既有流程，不變)
   ├─ bazi.Calculate(birth) → FourPillars
   ├─ personality.Generate(fourPillars) → []Trait
   └─ npc.Generate(...) → *NPCProfile

5. Prompt 建構 (既有，不變)
   └─ llm.BuildPrompt(npc) → string

6. OpenAIClient.Generate(prompt) (client.go)
   ├─ 組裝 APIRequest{Model, Messages, MaxTokens, Temperature}
   ├─ json.Marshal → 請求 body
   ├─ http.NewRequest("POST", BaseURL+"/chat/completions", body)
   ├─ 設定 Header：Content-Type、Authorization: Bearer
   ├─ httpClient.Do(req)
   ├─ 檢查 status code（非 200 → 錯誤含 body 片段）
   ├─ json.Unmarshal → APIResponse
   └─ 回傳 choices[0].message.content

7. CLI 輸出 (既有，不變)
   └─ 依 --format 印出 text 或 json
```

### 設定解析優先序

對於 `APIKey` 與 `Model`：

```
CLI flag > 環境變數 > Provider 預設
```

- `--api-key xxx` → 直接用 xxx（最高優先，適合 CI/測試）
- 否則讀 `os.Getenv(Provider.EnvVar)`
- 環境變數也空 → 錯誤「請設定 --api-key 或環境變數 <EnvVar>」

`Model`：
- `--model xxx` → 用 xxx
- 否則用 `Provider.DefaultModel`

## 錯誤處理

依 CLAUDE.md「錯誤訊息三要素」原則（位置/上下文/技術細節）。

| 階段 | 錯誤情境 | 訊息格式 |
|------|---------|---------|
| Factory | 不支援的 provider | `不支援的 provider "xxx"（位置: factory.go NewClient），可用: openai, opencode-go, MiniMax` |
| Factory | 缺少 API key | `缺少 API key（位置: factory.go NewClient），請設定 --api-key 或環境變數 <EnvVar>` |
| Client | JSON 序列化失敗 | `序列化請求失敗（位置: client.go Generate）: %v` |
| Client | HTTP 請求建立失敗 | `建立請求失敗（位置: client.go Generate）: %v` |
| Client | 網路錯誤 | `API 呼叫失敗（位置: client.go Generate）: %v` |
| Client | API 回應非 200 | `API 回應錯誤 401（位置: client.go Generate）: <body 前 200 字>` |
| Client | 回應解析失敗 | `解析回應失敗（位置: client.go Generate）: %v` |
| Client | 回應無 choices | `API 回應無選擇（位置: client.go Generate）` |

### 設計原則

1. 錯誤一律 wrap：`fmt.Errorf("...: %w", err)` 保留堆疊
2. 位置明確：每個錯誤訊息含「位置: 檔名 + 函式」
3. 回應 body 截斷：非 200 時只回傳前 200 字
4. Fail-Fast：參數驗證在 Factory 階段完成
5. 環境變數缺失早警告：CLI 啟動時若缺 API key 直接退出

### 警告 vs 錯誤

- **Warn**：`--api-key` 明文傳入時印 `WARN: 不建議在 CLI 傳入金鑰，建議使用環境變數`
- **Error**：所有上述錯誤情境一律 exit code 1

## 測試策略

### 單元測試

| 檔案 | 測試目標 | 案例數 |
|------|---------|-------|
| `client_test.go` | `OpenAIClient.Generate` | 5 |
| `providers_test.go` | `DefaultProviders` 完整性 | 2 |
| `factory_test.go` | `NewClient` 各路徑 | 6 |

### `client_test.go` 案例（使用 `httptest.NewServer`）

1. 正常回應：mock 200 + 合法 JSON → 驗證回傳 content
2. 非 200 狀態：mock 401 + 錯誤 body → 驗證錯誤訊息含狀態碼與 body
3. 回應無 choices：mock 200 但 `choices: []` → 驗證「API 回應無選擇」錯誤
4. JSON 解析失敗：mock 200 + 壞 JSON → 驗證錯誤訊息
5. HTTP 請求驗證：mock server 驗證收到的 path、header、body 結構正確

### `providers_test.go` 案例

1. 預設表完整性：三個 provider 都存在且欄位非空
2. BaseURL 格式：每個 BaseURL 以 `https://` 開頭

### `factory_test.go` 案例

1. 已知 provider + 給 API key → 回傳正確設定的 client
2. 已知 provider + API key 空 + 環境變數有值 → 從環境讀取
3. 已知 provider + API key 空 + 環境變數空 → 錯誤訊息含 EnvVar 名稱
4. 未知 provider → 錯誤訊息含可用 provider 列表
5. 指定 model → 覆寫預設
6. 未指定 model → 用 DefaultModel

用 `t.Setenv("OPENAI_API_KEY", "test")` 設定環境變數測試。

### 整合測試（不寫實作、僅占位）

實際呼叫需要真實 API key，放到 `tests/integration/llm_integration_test.go` 並加 build tag：

```go
//go:build integration
```

執行：`go test -tags=integration ./tests/integration/...`（需設定對應環境變數）

### 覆蓋率目標

- Factory `NewClient`：所有 branch 100%
- Client `Generate`：所有錯誤路徑 100%
- 整體 `internal/llm/`：≥ 90%

### 不寫的測試

- 不 mock HTTP（用 `httptest` 模擬更貼近真實）
- 不測 Prompt 建構（既有測試已覆蓋）
- 不寫 benchmark

## 檔案異動清單

| 動作 | 檔案 |
|------|------|
| 修改 | `internal/llm/generator.go`（拆分為 client.go + factory.go） |
| 新增 | `internal/llm/client.go`（介面 + OpenAIClient） |
| 新增 | `internal/llm/factory.go`（NewClient） |
| 新增 | `internal/llm/providers.go`（預設表） |
| 新增 | `internal/llm/client_test.go` |
| 新增 | `internal/llm/providers_test.go` |
| 新增 | `internal/llm/factory_test.go` |
| 修改 | `cmd/npcgen/main.go`（新增 --provider/--model/--api-key） |

## 風險與緩解

| 風險 | 影響 | 緩解 |
|------|------|------|
| OpenCode Go 端點變更 | 連線失敗 | 預設表集中管理，未來變更只需改 `providers.go` |
| 三個 Provider 看似相同但實際有微差 | 某 Provider 特定行為失敗 | 保留 `LLMClient` 介面，未來可加 provider-specific client |
| API key 明文傳入 | 安全風險 | CLI 印出 WARN 警告，但不阻擋（測試/CI 仍可用） |
| 預設模型被 Provider 下架 | 預設呼叫失敗 | 使用者可用 `--model` 覆寫；文件標註模型清單需查 Provider 官網 |

## 不做的事（Out of Scope）

- 串流回應（streaming）：目前為一次性回傳，未來再考慮
- Function calling / Tool use：與 NPC 對話無關
- 重試邏輯（retry with backoff）：不在此次範圍
- 其他 Provider（Anthropic、Gemini、Ollama）：保留介面，未來獨立 PR
- 設定檔（YAML config）：使用者可用環境變數，CLI 不引入新機制

## 後續

通過審核後呼叫 `superpowers:writing-plans` skill 撰寫實作計畫。
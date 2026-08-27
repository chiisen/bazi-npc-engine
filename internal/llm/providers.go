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
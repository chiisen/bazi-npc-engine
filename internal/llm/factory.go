// Package llm - Factory 函式
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
//  1. 查 DefaultProviders[cfg.Provider]，找不到回傳錯誤
//  2. cfg.APIKey 空 → 讀 os.Getenv(provider.EnvVar)
//  3. cfg.APIKey 仍空 → 回傳錯誤（訊息含 EnvVar 名稱）
//  4. cfg.Model 空 → 用 provider.DefaultModel
//  5. cfg.Timeout 0 → 用 defaultTimeout
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

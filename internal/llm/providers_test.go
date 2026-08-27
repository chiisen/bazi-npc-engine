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
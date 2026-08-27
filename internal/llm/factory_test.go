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

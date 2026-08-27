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

# 變更記錄 (CHANGELOG)

> [!IMPORTANT]
> 本專案遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/) 規範。

## [Unreleased] - 未發佈版本

### 新增
- 多 LLM Provider 支援（OpenAI / OpenCode Go / MiniMax）
  - `--provider` 旗標切換 provider
  - `--model` 覆寫預設模型
  - `--api-key` 臨時指定金鑰（建議用環境變數）
  - 工廠函式 `llm.NewClient()` 統一處理設定解析
  - **CLI 已串接 LLM 呼叫**（無 API key 時跳過；失敗時警告但不中斷）
- `internal/llm/providers.go` 新增 Provider 預設表
- `internal/llm/client.go` 新增 OpenAIClient 實作（從 generator.go 拆分）
- 14 個新單元測試（provider 預設表 2 + client 錯誤路徑 5 + factory 路徑 7）
- 初始專案結構建立
- 八字計算引擎設計
- 人格生成器設計

### •Cancelled• 取消
- 待定

### 已變更
- 待定

---

## [0.1.0] - 2026-03-11

### 首次發佈

#### 新增功能
- 八字計算核心 (四柱、五行、十神)
- 人格生成引擎
- CLI 工具基本框架
- JSON 與文字雙輸出格式
- 支援種子再現功能

#### 技術ファイル
- `go.mod` 初始化
- Clean Architecture 目錄結構
- Makefile 建立基本指令
- `docs/` 目錄文件初始化

#### 文件
- `docs/PRD.md` - 產品需求文件
- `docs/ARCHITECTURE.md` - 架構設計文檔
- `docs/API.md` - API 規格文件
- `docs/DIRECTORY.md` - 目錄結構說明
- `docs/DECISIONS.md` - 技術決策記錄
- `docs/CHANGELOG.md` - 本文件

#### 作者
- [😸SAM]

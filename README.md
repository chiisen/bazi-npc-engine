# Bazi NPC Engine

> 使用中國傳統八字（四柱命理）生成 RPG NPC 的人工智能工具

[![Go Version](https://img.shields.io/badge/go-1.23-green)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Makefile](https://img.shields.io/badge/Makefile-available-blue)](Makefile)

## 📖 簡介

Bazi NPC Engine 是一個基於中國傳統命理學「八字」（Four Pillars of Destiny）的 NPC 生成引擎。透過分析出生時間的干支組合，計算五行比例與十神關係，自動生成具有獨特人格特徵與背景故事的 RPG 角色。

### 專案狀態

> [!NOTE]
> **當前版本**: v0.1.0 (初始釋出版)
>
> - ✅ 八字計算核心已完成
> - ✅ 人格生成引擎已實作
> - ✅ CLI 工具可使用
> - ✅ LLM 整合介面已定義
> - 🚧 內建 LLM API 連線需額外配置

### 核心特色

- 🌟 **純八字計算**：不依賴外部 API，完全本地計算
- 🧠 **人格生成引擎**：基於五行平衡與十神分析
- 🤖 **LLM 整合**：提供標準 Prompt 格式，輕鬆銜接大語言模型
- 🎮 **RPG 專用**：包含姓名、年齡、職業、背景故事、重要事件
- 🔮 **中國傳統文化**：完整實作干支、五行、十神系統
- 🛠️ **開發友好**：提供 Makefile 簡化開發流程

## 🏗️ 架構

```
bazi-npc-engine/
├── cmd/
│   └── npcgen/           # CLI 工具入口
│       └── main.go
├── internal/
│   ├── bazi/            # 八字計算核心
│   │   ├── calculate.go # 四柱計算（年/月/日/時柱）
│   │   ├── elements.go  # 五行比例與平衡分析
│   │   ├── shishen.go   # 十神關係分析
│   │   ├── constants.go # 干支常數定義
│   │   └── calculate_test.go # 單元測試
│   ├── personality/     # 人格生成模組
│   │   └── generator.go
│   ├── npc/             # NPC 設定生成模組
│   │   └── generator.go
│   └── llm/             # LLM 整合模組
│       └── generator.go
├── pkg/
│   └── types/           # 共用型別定義
│       └── types.go
├── tests/
│   ├── unit/            # 單元測試
│   └── integration/     # 整合測試
├── docs/
│   ├── PRD.md           # 產品需求文件
│   ├── ARCHITECTURE.md  # 架構設計文檔
│   ├── API.md           # API 規格文件
│   ├── DIRECTORY.md     # 目錄結構說明
│   └── DECISIONS.md     # 技術決策記錄
├── Makefile             # 開發指令集
└── CHANGELOG.md         # 變更記錄
```

### 模組說明

| 模組 | 功能 | 輸入 | 輸出 |
|------|------|------|------|
| `bazi` | 八字計算 | 出生時間 | 四柱干支、五行比例、十神 |
| `personality` | 人格生成 | 八字數據 | 人格特徵、優缺點 |
| `npc` | NPC 生成 | 人格數據 | 完整 NPC 設定 |
| `llm` | LLM 整合 | NPC 設定 | Prompt 文本 |

## 💻 安裝與使用

### 前置需求

- **Go 1.23** 或更高版本
- **Make** (建議使用，簡化開發流程)
- 網路連線（僅 LLM 整合時需要）

### 編譯

```bash
# 使用 Makefile (建議)
make build

# 或直接使用 go
go build -v ./...

# 編譯 CLI 工具
go build -o npcgen ./cmd/npcgen
```

### 執行測試

```bash
# 使用 Makefile (建議)
make test

# 執行所有測試
go test -v ./...

# 執行特定套件測試
go test ./internal/bazi/... -v

# 產生測試覆蓋率報告
go test ./... -cover

# 執行單元測試
make test-unit

# 執行整合測試
make test-integration
```

### 格式化程式碼

```bash
# 使用 Makefile (建議)
make fmt
make vet
make lint  # 結合 fmt 和 vet

# 或直接使用 go
go fmt ./...
go vet ./...

# 運行 npcgen
make run --birth "1995-10-01 14:00"

# 清理快取
make clean
```

## 🎮 CLI 使用指南

### 基本用法

```bash
# 生成 NPC（預設文字輸出）
npcgen --birth "1995-10-01 14:00"

# 輸出為 JSON 格式
npcgen --birth "1995-10-01 14:00" --format json

# 輸出到檔案
npcgen --birth "1995-10-01 14:00" --output npc.json

# 顯示詳細八字資訊
npcgen --birth "1995-10-01 14:00" --verbose

# 使用 Makefile 運行
make run --birth "1995-10-01 14:00" --verbose
```

### 命令列選項

| 參數 | 格式 | 預設值 | 說明 |
|------|------|--------|------|
| `--birth` | `YYYY-MM-DD HH:MM` | 必填 | 出生時間（24小時制） |
| `--format` | `text` / `json` | `text` | 輸出格式 |
| `--output` | 檔案路徑 | 標準輸出 | 輸出檔案路徑 |
| `--seed` | 整數 | 當前時間 | 隨機種子（用於再現） |
| `--verbose` | - | `false` | 顯示詳細八字資訊 |
| `--help` | - | - | 顯示幫助訊息 |
| `--version` | - | - | 顯示版本資訊 |

### 範例輸出

```bash
$ npcgen --birth "1995-10-01 14:00" --verbose

=== NPC 生成結果 ===

姓名: 張志強
年齡: 29
職業: 商人

人格特質:
  - 踏實穩重
  - 勤勞刻苦

背景故事:
  出生於家境普通家庭，自幼立志成為一名優秀的商人

重要事件:
  1. 遇見了人生的關鍵導師
  2. 經歷了一場重大變故
  3. 學會了珍貴的技能

=== 八字資訊 ===
年柱: 乙亥
月柱: 辛酉
日柱: 庚子
時柱: 丙午

五行比例:
  木: 16.67%
  火: 16.67%
  土: 8.33%
  金: 33.33%
  水: 25.00%

十神分析:
  年柱: 正財
  月柱: 七殺
  日柱: 日主
  時柱: 偏官
```

### 🐶 GEMINI / Claude Agent 使用提示

本專案遵循 **GEMINI.md** 規範。AI Agent 在貢獻代碼時應：

1. 每個函式上方必須附上繁體中文註解
2. 核心概念使用 `// ═══════════════════════════════════════════` 分隔線標示
3. 風險操作（如 `unwrap()`、`Unsafe`）需主動警示
4. 錯誤訊息需包含「位置、上下文、技術細節」三要素

## 📚 八字基礎知識

### 四柱（Four Pillars）

| 柱位 | 定義 | 代表意義 |
|------|------|----------|
| **年柱** | 出生年的干支 | 家庭背景、祖業 |
| **月柱** | 出生月的干支 | 事業、生涯發展 |
| **日柱** | 出生日的干支 | 命主自身、核心性格 |
| **時柱** | 出生時間的干支 | 子女、晚年、潛意識 |

### 天干地支（Heavenly Stems & Earthly Branches）

- **天干（10個）**：甲、乙、丙、丁、戊、己、庚、辛、壬、癸
- **地支（12個）**：子、丑、寅、卯、辰、巳、午、未、申、酉、戌、亥
- **干支組合**：60甲子循環（甲子→癸亥）

### 五行（Five Elements）

| 元素 | 性質 | 代表特質 |
|------|------|----------|
| **木**（Wood） | 生長、發展 | 富有遠見、創造力強 |
| **火**（Fire） | 熱情、亮點 | 充滿熱情、領導能力強 |
| **土**（Earth） | 穩定、承載 | 踏實穩重、責任感強 |
| **金**（Metal） | 堅毅、收斂 | 果斷堅毅、邏輯思維清晰 |
| **水**（Water） | 流動、智慧 | 機智靈活、適應力強 |

### 十神（Ten Gods）

| 十神 | 吉凶 | 代表意義 |
|------|------|----------|
| **正官** | 吉 | 約束、規範、社會地位 |
| **七殺** | 平 | 進取、挑戰、權威 |
| **正財** | 吉 | 現實利益、勤儉 |
| **偏財** | 吉 | 意外之財、投機 |
| **正印** | 吉 | 學習、支持、仁慈 |
| **偏印** | 平 | 偏門學問、奇特 |
| **食神** | 吉 | 享受、表達、溫和 |
| **傷官** | 平 | 才華、表達、反叛 |
| **比肩** | 平 | 同類、幫助、獨立 |
| **劫財** | 平 | 競爭、拼搏、冒險 |

## 🎯 人格生成邏輯

### 五行過旺分析

| 五行 | 過旺特質 | 建議方位 |
|------|----------|----------|
| 木過旺 | 富有遠見、創造力強 | 東方、東南 |
| 火過旺 | 充滿熱情、領導能力強 | 南方 |
| 土過旺 | 踏實穩重、責任感強 | 中央、東南 |
| 金過旺 | 果斷堅毅、邏輯思維清晰 | 西方、西北 |
| 水過旺 | 機智靈活、適應力強 | 北方、東北 |

### 十神對人格影響

- **正官多**：重視規則與名譽，適合公職、管理
- **七殺多**：具有進取心，適合創業、武職
- **正財多**：注重現實利益，適合商業、理財
- **偏財多**：善於把握機會，適合投機、經商
- **印綬多**：重視學習，適合教育、研究

## 🤖 LLM 整合

### API 介面

```go
package llm

// 生成 NPC Profile
npc := npc.Generate(personality, seed)

// 生成 LLM Prompt
prompt := llm.BuildPrompt(npc)

// 生成 System Prompt
system := llm.BuildSystemPrompt(npc)

// 生成 Scene Prompt（特定情境）
scenePrompt := llm.BuildScenePrompt(npc, "場景描述", "用戶描述")
```

### 使用範例

```go
import (
    "github.com/chiisen/bazi-npc-engine/internal/bazi"
    "github.com/chiisen/bazi-npc-engine/internal/personality"
    "github.com/chiisen/bazi-npc-engine/internal/npc"
    "github.com/chiisen/bazi-npc-engine/internal/llm"
)

func main() {
    // 計算八字
    birthTime, _ := time.Parse("2006-01-02 15:04", "1995-10-01 14:00")
    baziData := bazi.Calculate(birthTime)

    // 分析五行與十神
    baziData.Elements = bazi.CalculateElements(baziData)
    baziData.TenGods = bazi.CalculateTenGods(baziData)

    // 生成人格
    pers := personality.Generate(baziData)

    // 生成 NPC
    npcProfile := npc.Generate(pers, 0)

    // 生成 LLM Prompt
    prompt := llm.BuildPrompt(npcProfile)

    // 發送至 LLM API
    client := llm.NewHTTPClient("https://api.openai.com/v1/chat/completions",
        "your-api-key", "gpt-4")
    response, err := client.Generate(prompt)
}
```

### LLM 內建提示詞

**System Prompt**：
```
You are role-playing as [姓名], a [年齡]-year-old [職業].
[人格特質]
Your background: [背景故事]
Respond in character as [姓名]. Be authentic to their personality and experiences.
```

**Scene Prompt**：
```
You are [姓名], a [年齡]-year-old [職業].
Current Scene: [場景描述]
User: [用戶描述]
Respond naturally as [姓名] would in this situation.
```

## 🛠️ 開發指南

### 新增功能

1. 修改 `internal/bazi/calculate.go` 添加八字計算邏輯
2. 修改 `internal/personality/generator.go` 調整人格生成
3. 修改 `internal/npc/generator.go` 增加 NPC 設定
4. 編寫測試於 `internal/bazi/calculate_test.go`

### 程式碼規範

- 遵循 Go 標準格式（`go fmt`）
- 函式必須有繁體中文註解
- 新增功能需包含單元測試

### 常用指令

```bash
# 編譯
go build ./...

# 測試
go test ./...

# 格式化
go fmt ./...

# 靜態分析
go vet ./...

# 產生文件
godoc -http=:6060
```

## 📜 變更記錄

### v0.1.0

- ✅ 八字計算（年/月/日/時柱）
- ✅ 五行比例與平衡分析
- ✅ 十神關係分析
- ✅ 人格生成引擎
- ✅ NPC 設定生成
- ✅ CLI 工具
- ✅ LLM 整合支援

## 📄 授權

MIT License

## 👤 作者

[😸SAM](https://github.com/chiisen)

## 🙏 致謝

- 八字計算參考傳統命理學原理
- 簡化模型以適合 RPG 應用

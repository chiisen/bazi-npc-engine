# Bazi NPC Engine 目錄結構

```
bazi-npc-engine/
├── cmd/                    # 應用程式入口點
│   └── npcgen/            # CLI 主程式
│       └── main.go        # CLI 啟動檔案
│       └── cmd.go         # CLI 命令定義
│       └── options.go     # CLI 參數定義
├── internal/              # 專業邏輯 (不對外公開)
│   ├── bazi/              # 八字計算核心
│   │   ├── bazi.go        # 八字結構與計算
│   │   ├── saizhu.go      # 四柱計算
│   │   ├── wuxing.go      # 五行計算
│   │   └── shishen.go     # 十神計算
│   ├── personality/       # 人格生成引擎
│   │   ├── generator.go   # 人格生成器
│   │   └── traits.go      # 特質映射表
│   ├── npc/               # NPC 設定生成
│   │   ├── generator.go   # NPC 生成器
│   │   └── names.go       # 姓名生成
│   └── llm/               # LLM Prompt 生成
│       ├── builder.go     # Prompt 建構器
│       └── api.go         # LLM API 串接
├── pkg/                   # 公開ライブラリ (可被其他專案引用)
│   └── types/             # 共用資料型別
│       └── bazi.go        # 八定義
│       └── personality.go # 人格定義
│       └── npc.go         # NPC 定義
├── tests/                 # 測試檔案
│   ├── unit/              # 單元測試
│   └── integration/       # 整合測試
├── docs/                  # 文件
│   ├── PRD.md             # 產品需求文件
│   ├── ARCHITECTURE.md    # 架構設計文檔
│   ├── API.md             # API 規格文件
│   └── DECISIONS.md       # 技術決策記錄
├── .cursorrules           # Cursor AI 助手指引
├── .gitignore
├── go.mod                 # Go module 定義
├── go.sum                 # Go 依賴版本鎖定
├── Makefile               # 常用指令集
└── README.md              # 專案說明
```

## 目錄說明

### cmd/
放置應用程式的 main 函式，每個子目錄代表一個可執行檔。

### internal/
放置專案的核心邏輯，這些 package 不會被其他專案 import。

### pkg/
放置可被其他專案使用的公開 API，適合未來模組化發布。

### tests/
存放所有測試檔案，分為 unit 和 integration。

### docs/
存放所有技術文件。

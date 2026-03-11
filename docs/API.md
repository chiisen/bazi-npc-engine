# API 規格文件 (API Specification)

## CLI 使用說明

### 基本用法

```bash
npcgen --birth "1995-10-01 14:00"
```

### 完整參數

```bash
npcgen [flags]

Flags:
  --birth string     出生時間 (格式: YYYY-MM-DD HH:MM)
  --format string    輸出格式 (json/text, default: "text")
  --output string    輸出檔案路徑 (預設: 標準輸出)
  --seed int         隨機種子 (用於再現相同結果)
  --verbose          顯示詳細資訊
  --help             顯示幫助訊息
```

### 範例

#### 1. 基本生成 (文字輸出)

```bash
$ npcgen --birth "1995-10-01 14:00"

=== NPC 生成結果 ===

Name: 林濤 (Lin Tao)
Age: 30
Occupation: 商人

人格特質:
- 深思熟慮
- 冷靜沉著
- 善於投機

背景故事:
出生於 商人家庭，童年經歷戰亂导致家道中落。

重要事件:
- 10 歲：戰亂中與家人失散
- 20 歲：獨立經商
- 30 歲：建立工商網絡
```

#### 2. JSON 輸出

```bash
$ npcgen --birth "1995-10-01 14:00" --format json

{
  "npc": {
    "name": "Lin Tao",
    "age": 30,
    "occupation": "Merchant",
    "personality": ["strategic", "calm", "opportunistic"],
    "background": "Merchant family ruined by war.",
    "life_events": [
      "Lost family at age 10 during war",
      "Started independent business at 20",
      "Built commercial network at 30"
    ]
  },
  "bazi": {
    "year": "乙亥",
    "month": "辛酉",
    "day": "庚子",
    "hour": "壬辰",
    "elements": {
      "wood": 2.5,
      "fire": 0,
      "earth": 3.0,
      "metal": 3.5,
      "water": 2.0
    }
  }
}
```

#### 3. 儲存結果到檔案

```bash
$ npcgen --birth "1995-10-01 14:00" --output npc.json
```

#### 4. 再現相同結果 (使用 seed)

```bash
$ npcgen --birth "1995-10-01 14:00" --seed 12345
```

## 資料格式 (Data Formats)

### Bazi Format (八字格式)

| 欄位 | 型別 | 說明 |
| :--- | :--- | :--- |
| year | string | 年柱 (天干+地支) |
| month | string | 月柱 (天干+地支) |
| day | string | 日柱 (天干+地支) |
| hour | string | 時柱 (天干+地支) |
| elements | map | 五行比例 (金木水火土) |
| ten_gods | map | 十神分析結果 |

### Personality Format (人格格式)

| 欄位 | 型別 | 說明 |
| :--- | :--- |
| traits | []string | 人格特質 |
| strengths | []string | 優點 |
| weaknesses | []string | 缺點 |
| behavior | []string | 行為模式 |

### NPC Profile Format (NPC 設定格式)

| 欄位 | 型別 | 說明 |
| :--- | :--- | :--- |
| name | string | 角色姓名 |
| age | int | 年齡 |
| occupation | string | 職業 |
| personality | []string | 總結描述 |
| background | string | 背景故事 |
| life_events | []string | 重要事件 |

## HTTP API (未來擴展)

### 生成 NPC

```http
POST /api/v1/npc/generate
Content-Type: application/json

{
  "birth_datetime": "1995-10-01T14:00:00+08:00",
  "format": "json"
}

Response 200 OK
{
  "npc": { ... },
  "bazi": { ... }
}
```

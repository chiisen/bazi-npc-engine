# PRD --- Bazi NPC Generator

Version: 0.1 Author: Sam Project Draft Date: 2026-03-11

------------------------------------------------------------------------

# 1. 專案概述 (Project Overview)

**Bazi NPC Generator** 是一個使用「八字（Bazi）」作為角色種子 (seed) 的
RPG NPC 生成系統。

系統將：

出生時間 → 八字 → 人格模型 → NPC Persona → LLM 對話 / 行為

目標是建立 **具有一致性人格與人生邏輯的 AI NPC**。

------------------------------------------------------------------------

# 2. 問題定義 (Problem Statement)

現有 NPC 生成器通常有以下問題：

-   人格隨機，缺乏邏輯一致性
-   背景故事與性格不匹配
-   AI 對話缺乏長期人格穩定

例如：

勇敢但極端逃避責任\
極度忠誠但同時背叛組織

缺乏一個 **deterministic persona engine**。

------------------------------------------------------------------------

# 3. 解決方案 (Solution)

使用 **八字系統作為人格生成算法**。

八字本質：

seed = 出生時間

人格 = f(五行, 十神, 格局)

人生 = g(大運, 流年)

因此可建立：

Deterministic Persona System

------------------------------------------------------------------------

# 4. 專案目標 (Goals)

主要目標

1.  建立 NPC 人格生成引擎
2.  生成可解釋的人格模型
3.  提供 LLM Persona Prompt
4.  支援 RPG NPC 生成功能

次要目標

-   支援遊戲 NPC
-   支援 AI Agent Persona
-   支援故事生成

------------------------------------------------------------------------

# 5. 非目標 (Non Goals)

以下不在 MVP 範圍：

-   完整 RPG 遊戲
-   完整命理分析系統
-   商業命理服務

------------------------------------------------------------------------

# 6. 系統架構 (System Architecture)

High Level Architecture

Birth Datetime ↓ Bazi Engine ↓ Personality Generator ↓ NPC Profile ↓ LLM
Persona Prompt ↓ NPC Dialogue / Behavior

------------------------------------------------------------------------

# 7. 模組設計 (Modules)

## 7.1 Bazi Engine

功能：

-   計算四柱
-   五行比例
-   十神分析

輸入

出生時間

輸出

{ "year": "","month": "","day": "","hour": "","elements": {},
"ten_gods": {} }

------------------------------------------------------------------------

## 7.2 Personality Generator

將八字轉換為人格特徵

輸出

{ "traits": \[\], "strengths": \[\], "weakness": \[\],
"behavior_pattern": \[\] }

------------------------------------------------------------------------

## 7.3 NPC Profile Generator

輸出完整 NPC 設定

{ "name": "","occupation": "","personality": \[\], "background":
"","life_events": \[\] }

------------------------------------------------------------------------

## 7.4 LLM Persona Builder

生成 Prompt

Example

You are an NPC with the following personality:

Traits: - ambitious - strategic - suspicious

Background: Merchant family ruined by war.

------------------------------------------------------------------------

# 8. API 設計 (API Design)

POST /npc/generate

Input

{ "birth_datetime": "1990-03-12T10:00:00" }

Output

{ "npc": { "name": "","occupation": "","personality": \[\],
"background": "" } }

------------------------------------------------------------------------

# 9. 技術架構 (Tech Stack)

Backend

Python / Go

AI

LLM API

Data

JSON / SQLite

------------------------------------------------------------------------

# 10. MVP 功能 (Minimum Viable Product)

MVP 包含：

1.  八字計算
2.  人格生成
3.  NPC JSON 輸出
4.  CLI 生成工具

Example CLI

npcgen "1995-10-01 14:00"

------------------------------------------------------------------------

# 11. CLI 設計

Command

npcgen

Example

npcgen --birth "1995-10-01 14:00"

Output

NPC: Name: Lin Tao Personality: Strategic, Calm, Opportunistic
Occupation: Merchant

------------------------------------------------------------------------

# 12. 未來擴展 (Future Work)

-   NPC 人生模擬
-   流年事件
-   NPC Agent
-   Game Engine integration

------------------------------------------------------------------------

# 13. 開源計畫 (Open Source Plan)

License

MIT

Repository Structure

/bazi-engine /personality-engine /npc-generator /cli

------------------------------------------------------------------------

# 14. 成功指標 (Success Metrics)

1.  NPC 人格一致性
2.  NPC 故事合理性
3.  LLM Persona 表現

------------------------------------------------------------------------

# 15. 範例 NPC

Birth Time

1992-07-18 11:00

Generated NPC

Name: Chen Rui

Traits

-   strategic
-   ambitious
-   secretive

Background

Former imperial scholar turned underground information broker.

// Package types 共用資料型別定義
//
// 此 package 包含所有公開的資料結構，可供其他專案引用。
package types

// Bazi 八字盤結構 (與 internal/bazi 包中的 Bazi 類似，
// 這裡保留以供未來外部 use 时使用)
type Bazi struct {
	Year   string
	Month  string
	Day    string
	Hour   string
}

// Personality 人格特徵結構
type Personality struct {
	Traits     []string
	Strengths  []string
	Weaknesses []string
	Behavior   []string
}

// NPCProfile NPC 設定結構
type NPCProfile struct {
	Name        string
	Age         int
	Occupation  string
	Personality []string
	Background  string
	LifeEvents  []string
}

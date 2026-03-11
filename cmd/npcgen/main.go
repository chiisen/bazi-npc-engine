// Package main CLI 入口點
//
// 功能：處理命令列參數並啟動 NPC 生成流程
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/chiisen/bazi-npc-engine/internal/bazi"
	"github.com/chiisen/bazi-npc-engine/internal/npc"
	"github.com/chiisen/bazi-npc-engine/internal/personality"
)

// CLI 參數結構
type CLIOptions struct {
	Birth     string
	Format    string
	Output    string
	Seed      int
	Verbose   bool
}

// ParseOptions 解析命令列選項
func ParseOptions() *CLIOptions {
	birth := flag.String("birth", "", "出生時間 (格式: YYYY-MM-DD HH:MM)")
	format := flag.String("format", "text", "輸出格式 (json/text)")
	output := flag.String("output", "", "輸出檔案路徑")
	seed := flag.Int("seed", 0, "隨機種子 (用於再現相同結果)")
	verbose := flag.Bool("verbose", false, "顯示詳細資訊")
	flag.Parse()

	return &CLIOptions{
		Birth:  *birth,
		Format: *format,
		Output: *output,
		Seed:   *seed,
		Verbose: *verbose,
	}
}

// FormatBirth 輸出出生時間
func FormatBirth(opts *CLIOptions) string {
	if opts.Birth == "" {
		return "未提供"
	}
	return opts.Birth
}

// PrintResult 輸出結果
func PrintResult(npcProfile *npc.NPCProfile, bazi *bazi.Bazi, opts *CLIOptions) error {
	if opts.Format == "json" {
		data, err := json.MarshalIndent(npcProfile, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	} else {
		// Text 格式輸出
		fmt.Println("=== NPC 生成結果 ===")
		fmt.Println()
		fmt.Printf("姓名: %s\n", npcProfile.Name)
		fmt.Printf("年齡: %d\n", npcProfile.Age)
		fmt.Printf("職業: %s\n", npcProfile.Occupation)
		fmt.Println()
		fmt.Println("人格特質:")
		for _, trait := range npcProfile.Personality {
			fmt.Printf("  - %s\n", trait)
		}
		fmt.Println()
		fmt.Println("背景故事:")
		fmt.Printf("  %s\n", npcProfile.Background)
		fmt.Println()
		fmt.Println("重要事件:")
		for i, event := range npcProfile.LifeEvents {
			fmt.Printf("  %d. %s\n", i+1, event)
		}
		fmt.Println()

		// 如果開啟 verbose，顯示八字資訊
		if opts.Verbose {
			fmt.Println("=== 八字資訊 ===")
			fmt.Printf("年柱: %s\n", bazi.Year)
			fmt.Printf("月柱: %s\n", bazi.Month)
			fmt.Printf("日柱: %s\n", bazi.Day)
			fmt.Printf("時柱: %s\n", bazi.Hour)
			fmt.Println()
			fmt.Println("五行比例:")
			for elem, ratio := range bazi.Elements {
				fmt.Printf("  %s: %.2f\n", elem, ratio*100)
			}
			fmt.Println()
			fmt.Println("十神分析:")
			for col, shen := range bazi.TenGods {
				fmt.Printf("  %s: %s\n", col, shen)
			}
		}
	}
	return nil
}

// GenerateNPC 生成 NPC
func GenerateNPC(opts *CLIOptions) error {
	// 解析出生時間
	birthTime, err := time.Parse("2006-01-02 15:04", opts.Birth)
	if err != nil {
		return fmt.Errorf("出生時間格式錯誤，請使用 YYYY-MM-DD HH:MM 格式: %v", err)
	}

	// 計算八字
	baziData := bazi.Calculate(birthTime)

	// 計算五行比例
	baziData.Elements = bazi.CalculateElements(baziData)

	// 計算十神
	baziData.TenGods = bazi.CalculateTenGods(baziData)

	// 生成人格
	pers := personality.Generate(baziData)

	// 生成 NPC
	npcProfile := npc.Generate(pers, opts.Seed)

	// 輸出結果
	return PrintResult(npcProfile, &baziData, opts)
}

// 定義 CLI 命令
func Run() int {
	fmt.Println("Bazi NPC Generator v0.1.0")
	fmt.Println("使用八字生成 RPG NPC")
	fmt.Println()

	opts := ParseOptions()

	// 檢查出生時間是否提供
	if opts.Birth == "" {
		fmt.Println("錯誤：請提供出生時間")
		fmt.Println("用法: npcgen --birth \"1995-10-01 14:00\"")
		return 1
	}

	// 生成 NPC
	err := GenerateNPC(opts)
	if err != nil {
		log.Fatalf("生成 NPC 失敗: %v", err)
		return 1
	}

	return 0
}

func main() {
	code := Run()
	if code != 0 {
		fmt.Println()
		fmt.Println("提示：使用 --help 查看更多資訊")
	}
}

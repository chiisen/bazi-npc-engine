// Package main CLI 入口點
//
// 功能：處理命令列參數並啟動 NPC 生成流程
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
	Help      bool
	Version   bool
}

// ParseOptions 解析命令列選項
func ParseOptions() *CLIOptions {
	birth := flag.String("birth", "", "出生時間 (格式: YYYY-MM-DD HH:MM)")
	format := flag.String("format", "text", "輸出格式 (json/text)")
	output := flag.String("output", "", "輸出檔案路徑")
	seed := flag.Int("seed", 0, "隨機種子 (用於再現相同結果)")
	verbose := flag.Bool("verbose", false, "顯示詳細資訊")
	showHelp := flag.Bool("help", false, "顯示帮助訊息")
	showVersion := flag.Bool("version", false, "顯示版本訊息")
	flag.Parse()

	return &CLIOptions{
		Birth:   *birth,
		Format:  *format,
		Output:  *output,
		Seed:    *seed,
		Verbose: *verbose,
		Help:    *showHelp,
		Version: *showVersion,
	}
}

// PrintHelp 顯示幫助訊息
func PrintHelp() {
	fmt.Println(`Bazi NPC Generator - 使用八字生成 RPG NPC

用法：
  npcgen [選項]

選項：
  --birth string    出生時間 (格式: YYYY-MM-DD HH:MM) [必填]
  --format string   輸出格式 (json/text, 預設: text)
  --output string   輸出檔案路徑 (預設: 標準輸出)
  --seed int        隨機種子 (用於再現相同結果)
  --verbose         顯示詳細八字資訊
  --help            顯示幫助訊息
  --version         顯示版本訊息

範例：
  npcgen --birth "1995-10-01 14:00"
  npcgen --birth "1995-10-01 14:00" --format json
  npcgen --birth "1995-10-01 14:00" --output npc.json
  npcgen --birth "1995-10-01 14:00" --verbose
`)
}

// PrintVersion 顯示版本訊息
func PrintVersion() {
	fmt.Println("Bazi NPC Generator v0.1.0")
}

// PrintResult 輸出結果
func PrintResult(npcProfile *npc.NPCProfile, baziData *bazi.Bazi, opts *CLIOptions) error {
	// 如果指定了輸出檔案
	if opts.Output != "" {
		file, err := os.Create(opts.Output)
		if err != nil {
			return fmt.Errorf("無法建立輸出檔案: %v", err)
		}
		defer file.Close()

		if opts.Format == "json" {
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			err = encoder.Encode(npcProfile)
		} else {
			// 將 text 輸出寫入檔案
			var content string
			content = fmt.Sprintf("=== NPC 生成結果 ===\n\n")
			content += fmt.Sprintf("姓名: %s\n", npcProfile.Name)
			content += fmt.Sprintf("年齡: %d\n", npcProfile.Age)
			content += fmt.Sprintf("職業: %s\n", npcProfile.Occupation)
			content += "\n人格特質:\n"
			for _, trait := range npcProfile.Personality {
				content += fmt.Sprintf("  - %s\n", trait)
			}
			content += fmt.Sprintf("\n背景故事:\n  %s\n", npcProfile.Background)
			content += "\n重要事件:\n"
			for i, event := range npcProfile.LifeEvents {
				content += fmt.Sprintf("  %d. %s\n", i+1, event)
			}

			if opts.Verbose {
				content += "\n=== 八字資訊 ===\n"
				content += fmt.Sprintf("年柱: %s\n", baziData.Year)
				content += fmt.Sprintf("月柱: %s\n", baziData.Month)
				content += fmt.Sprintf("日柱: %s\n", baziData.Day)
				content += fmt.Sprintf("時柱: %s\n", baziData.Hour)
				content += "\n五行比例:\n"
				for elem, ratio := range baziData.Elements {
					content += fmt.Sprintf("  %s: %.2f\n", elem, ratio*100)
				}
				content += "\n十神分析:\n"
				for col, shen := range baziData.TenGods {
					content += fmt.Sprintf("  %s: %s\n", col, shen)
				}
			}

			_, err = file.WriteString(content)
		}
		return err
	}

	// 標準輸出
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
			fmt.Printf("年柱: %s\n", baziData.Year)
			fmt.Printf("月柱: %s\n", baziData.Month)
			fmt.Printf("日柱: %s\n", baziData.Day)
			fmt.Printf("時柱: %s\n", baziData.Hour)
			fmt.Println()
			fmt.Println("五行比例:")
			for elem, ratio := range baziData.Elements {
				fmt.Printf("  %s: %.2f\n", elem, ratio*100)
			}
			fmt.Println()
			fmt.Println("十神分析:")
			for col, shen := range baziData.TenGods {
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

// Run CLI 主流程
func Run() int {
	opts := ParseOptions()

	// 處理 special 選項
	if opts.Help {
		PrintHelp()
		return 0
	}

	if opts.Version {
		PrintVersion()
		return 0
	}

	// 檢查出生時間是否提供
	if opts.Birth == "" {
		fmt.Println("錯誤：請提供出生時間")
		fmt.Println("使用 --help 查看使用說明")
		return 1
	}

	// 檢查格式是否有效
	if opts.Format != "text" && opts.Format != "json" {
		fmt.Printf("錯誤：無效的格式 '%s'，請使用 'text' 或 'json'\n", opts.Format)
		return 1
	}

	// 生成 NPC
	err := GenerateNPC(opts)
	if err != nil {
		log.Printf("錯誤：生成 NPC 失敗: %v", err)
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
	os.Exit(code)
}

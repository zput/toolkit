package main

import (
	"fmt"
	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite/auto/extractor"
	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite/auto/generator"

	"os"
	"path/filepath"
)

// 步骤1: 递归追踪函数调用链，提取所有函数
// 步骤2: 逐个分析函数，提取模型

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法: go run main.go <工作目录> <函数路径> [模型路径...] [输出目录]")
		fmt.Println("示例: go run main.go .        internal/service/web/game.go:115 model/db,xxx/model/db")
		fmt.Println("      go run main.go .       internal/service/web/game.go:115 model/db,xxx/model/db ./output")
		os.Exit(1)
	}

	workDir := os.Args[1]
	funcPath := os.Args[2]
	modelPaths := []string{}
	if len(os.Args) > 3 {
		// 解析模型路径（逗号分隔）
		for _, path := range splitByComma(os.Args[3]) {
			if path != "" {
				modelPaths = append(modelPaths, path)
			}
		}
	}

	fmt.Printf("工作目录: %s\n", workDir)
	fmt.Printf("函数路径: %s\n", funcPath)
	fmt.Printf("模型路径: %v\n", modelPaths)

	// 步骤1: 递归追踪函数调用链，提取所有函数
	fmt.Println("\n=== 步骤1: 递归追踪函数调用链 ===")
	allFuncs, err := extractor.TrackFunctionCalls(workDir, funcPath)
	if err != nil {
		fmt.Printf("错误: 追踪函数调用失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("找到 %d 个函数:\n", len(allFuncs))
	for i, fn := range allFuncs {
		fmt.Printf("  %d. %s\n", i+1, fn.Key())
	}

	// 步骤2: 逐个分析函数，提取模型
	fmt.Println("\n=== 步骤2: 提取模型 ===")
	models, err := extractor.ExtractModels(workDir, allFuncs, modelPaths)
	if err != nil {
		fmt.Printf("错误: 提取模型失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("找到 %d 个模型:\n", len(models))
	for i, model := range models {
		fmt.Printf("  %d. %s (%s)\n", i+1, model.FullName, model.PackagePath)
	}

	// 步骤3: 生成 Mock 代码
	if len(models) > 0 {
		fmt.Println("\n=== 步骤3: 生成 Mock 代码 ===")

		// 确定输出目录（默认使用工作目录下的 stub_test/unittest/mock_modules）
		outputDir := filepath.Join(workDir, "stub_test", "unittest", "mock_modules")

		// 如果命令行参数提供了输出目录，使用命令行参数
		if len(os.Args) > 4 {
			outputDir = os.Args[4]
		}

		fmt.Printf("输出目录: %s\n", outputDir)

		err := generator.GenerateMockCode(workDir, funcPath, models, outputDir)
		if err != nil {
			fmt.Printf("错误: 生成 Mock 代码失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ Mock 代码生成成功")
	} else {
		fmt.Println("\n警告: 未找到任何模型，跳过代码生成")
	}
}

// splitByComma 按逗号分割字符串
func splitByComma(s string) []string {
	var result []string
	start := 0
	for i, char := range s {
		if char == ',' {
			if i > start {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

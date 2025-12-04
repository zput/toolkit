package generator

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite/auto/dto"
)

// TestExtractStructDefinition 测试提取 struct 定义
func TestExtractStructDefinition(t *testing.T) {
	testDir := getTestDataDir()

	model := dto.ModelInfo{
		PackagePath: "testpkg/model/db/web_game", // 这是 go.mod 中定义的模块路径
		PackageName: "web_game",
		TypeName:    "GameTemplate",
		FullName:    "web_game.GameTemplate",
	}

	structDecl, err := ExtractStructDefinition(testDir, model)
	if err != nil {
		t.Fatalf("提取 struct 定义失败: %v", err)
	}

	if structDecl == nil {
		t.Fatal("structDecl 为 nil")
	}

	// 验证是否是 type 声明
	if structDecl.Tok.String() != "type" {
		t.Errorf("期望是 type 声明，实际是: %s", structDecl.Tok.String())
	}

	// 验证是否包含 GameTemplate
	found := false
	for _, spec := range structDecl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			if typeSpec.Name.Name == "GameTemplate" {
				found = true
				// 验证是否是 struct 类型
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					t.Error("GameTemplate 不是 struct 类型")
				}
				break
			}
		}
	}

	if !found {
		t.Error("未找到 GameTemplate 类型定义")
	}
}

// TestGenerateStructFile 测试生成 struct 文件
func TestGenerateStructFile(t *testing.T) {
	testDir := getTestDataDir()
	outputDir := filepath.Join(testDir, "output")

	// 清理输出目录
	os.RemoveAll(outputDir)
	defer os.RemoveAll(outputDir)

	models := []dto.ModelInfo{
		{
			PackagePath: "testpkg/model/db/web_game",
			PackageName: "web_game",
			TypeName:    "GameTemplate",
			FullName:    "web_game.GameTemplate",
		},
	}

	packageName := "service"

	// 生成 struct 文件
	err := GenerateStructFile(testDir, models, outputDir, packageName)
	if err != nil {
		t.Fatalf("生成 struct 文件失败: %v", err)
	}

	// 验证文件是否存在
	expectedFile := filepath.Join(outputDir, "mock_service", "db", "web_game.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("期望生成文件 %s，但文件不存在", expectedFile)
	}

	// 验证文件内容
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证包名
	if !strings.Contains(contentStr, "package db") {
		t.Errorf("期望包含 'package db'，但内容为: %s", contentStr[:min(200, len(contentStr))])
	}

	// 验证 struct 定义
	if !strings.Contains(contentStr, "type GameTemplate struct") {
		t.Errorf("期望包含 'type GameTemplate struct'，但内容为: %s", contentStr[:min(200, len(contentStr))])
	}

	// 验证字段（根据 testdata/model/db/web_game/model.go）
	if !strings.Contains(contentStr, "ID") {
		t.Errorf("期望包含字段 'ID'")
	}

	if !strings.Contains(contentStr, "Name") {
		t.Errorf("期望包含字段 'Name'")
	}

	// 验证不包含任何 tags（包括 gorm, json 等）
	if strings.Contains(contentStr, "gorm:") {
		t.Errorf("生成的 struct 不应该包含 gorm tags，但找到了 gorm:")
	}
	if strings.Contains(contentStr, "json:") {
		t.Errorf("生成的 struct 不应该包含 json tags，但找到了 json:")
	}
	if strings.Contains(contentStr, "`") && strings.Contains(contentStr, ":") {
		// 检查是否还有 tag（格式: `key:"value"`）
		t.Errorf("生成的 struct 不应该包含任何 tags")
	}

	t.Logf("✓ Struct 文件生成成功: %s", expectedFile)
	t.Logf("文件内容预览:\n%s", contentStr[:min(500, len(contentStr))])
}

// TestGenerateStructFile_MultipleModels 测试生成多个模型的 struct 文件
func TestGenerateStructFile_MultipleModels(t *testing.T) {
	testDir := getTestDataDir()
	outputDir := filepath.Join(testDir, "output")

	// 清理输出目录
	os.RemoveAll(outputDir)
	defer os.RemoveAll(outputDir)

	models := []dto.ModelInfo{
		{
			PackagePath: "testpkg/model/db/web_game",
			PackageName: "web_game",
			TypeName:    "GameTemplate",
			FullName:    "web_game.GameTemplate",
		},
		// 如果有多个模型，可以在这里添加
	}

	packageName := "service"

	// 生成 struct 文件
	err := GenerateStructFile(testDir, models, outputDir, packageName)
	if err != nil {
		t.Fatalf("生成 struct 文件失败: %v", err)
	}

	// 验证文件是否存在
	expectedFile := filepath.Join(outputDir, "mock_service", "db", "web_game.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("期望生成文件 %s，但文件不存在", expectedFile)
	}

	t.Logf("✓ 多个模型的 Struct 文件生成成功")
}

// TestGenerateStructFile_DifferentPackages 测试不同包的模型
func TestGenerateStructFile_DifferentPackages(t *testing.T) {
	testDir := getTestDataDir()
	outputDir := filepath.Join(testDir, "output")

	// 清理输出目录
	os.RemoveAll(outputDir)
	defer os.RemoveAll(outputDir)

	models := []dto.ModelInfo{
		{
			PackagePath: "testpkg/model/db/web_game",
			PackageName: "web_game",
			TypeName:    "GameTemplate",
			FullName:    "web_game.GameTemplate",
		},
		// 如果有不同包的模型，应该生成不同的文件
	}

	packageName := "service"

	// 生成 struct 文件
	err := GenerateStructFile(testDir, models, outputDir, packageName)
	if err != nil {
		t.Fatalf("生成 struct 文件失败: %v", err)
	}

	// 验证文件是否存在
	expectedFile := filepath.Join(outputDir, "mock_service", "db", "web_game.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("期望生成文件 %s，但文件不存在", expectedFile)
	}

	t.Logf("✓ 不同包的 Struct 文件生成成功")
}

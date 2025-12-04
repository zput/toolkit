package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite/auto/dto"
)

// getTestDataDir 获取测试数据目录
func getTestDataDir() string {
	return filepath.Join(".", "testdata")
}

// TestGenerateMockCode 测试生成 Mock Go 代码
func TestGenerateMockCode(t *testing.T) {
	testDir := getTestDataDir()
	outputDir := filepath.Join(testDir, "output")

	// 清理输出目录
	os.RemoveAll(outputDir)
	defer os.RemoveAll(outputDir)

	// 使用 testdata 中的实际数据
	funcPath := filepath.Join(testDir, "service", "service.go") + ":9" // TestVarDecl
	models := []dto.ModelInfo{
		{
			PackagePath: "testpkg/model/db/web_game",
			PackageName: "web_game",
			TypeName:    "GameTemplate",
			FullName:    "web_game.GameTemplate",
		},
	}

	// 生成代码
	err := GenerateMockCode(testDir, funcPath, models, outputDir)
	if err != nil {
		t.Fatalf("生成 Mock 代码失败: %v", err)
	}

	// 验证文件是否存在（输出到 testdata 内部）
	expectedFile := filepath.Join(outputDir, "mock_modules/mock_service/service_service_go_9_mock.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("期望生成文件 %s，但文件不存在", expectedFile)
	}

	// 验证文件内容
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("读取生成的文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证包名（使用 testdata 中的实际包名）
	if !strings.Contains(contentStr, "package mock_service") {
		t.Errorf("期望包含 'package mock_service'，但内容为: %s", contentStr[:min(200, len(contentStr))])
	}

	// 验证类型名（使用 testdata 中的实际行号）
	if !strings.Contains(contentStr, "ServiceLine9Mock") {
		t.Errorf("期望包含 'ServiceLine9Mock'，但内容为: %s", contentStr[:min(200, len(contentStr))])
	}

	// 验证模型导入（使用 testdata 中的实际包路径）
	if !strings.Contains(contentStr, "testpkg/model/db/web_game") {
		t.Errorf("期望包含模型导入路径 'testpkg/model/db/web_game'")
	}

	// 验证模型列表（使用 testdata 中的实际模型）
	if !strings.Contains(contentStr, "new(web_game.GameTemplate)") {
		t.Errorf("期望包含 'new(web_game.GameTemplate)'")
	}

	// 验证接口实现
	if !strings.Contains(contentStr, "testfixtures.CallerDbInterface") {
		t.Errorf("期望包含 'testfixtures.CallerDbInterface'")
	}

	t.Logf("✓ Mock Go 文件生成成功: %s", expectedFile)
	t.Logf("文件内容预览:\n%s", contentStr[:min(500, len(contentStr))])
}

// TestGenerateYAMLFixtures 测试生成 YAML fixture 文件
func TestGenerateYAMLFixtures(t *testing.T) {
	testDir := getTestDataDir()
	outputDir := filepath.Join(testDir, "output")

	// 清理输出目录
	os.RemoveAll(outputDir)
	defer os.RemoveAll(outputDir)

	// 使用 testdata 中的实际数据
	packageName := "service" // 注意：这里应该是包名，不是 mock_ 前缀的包名
	models := []dto.ModelInfo{
		{
			PackagePath: "testpkg/model/db/web_game",
			PackageName: "web_game",
			TypeName:    "GameTemplate",
			FullName:    "web_game.GameTemplate",
		},
	}

	// 生成 YAML fixture 文件
	err := GenerateYAMLFixtures(outputDir, packageName, models)
	if err != nil {
		t.Fatalf("生成 YAML fixture 文件失败: %v", err)
	}

	// 验证文件是否存在（输出到 testdata 内部）
	expectedFiles := []string{
		filepath.Join(outputDir, "mock_modules", "mock_service", "testdata", "fixtures", "sign_template.yml"),
	}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("期望生成文件 %s，但文件不存在", expectedFile)
		} else {
			t.Logf("✓ YAML fixture 文件生成成功: %s", expectedFile)

			// 验证文件内容
			content, err := os.ReadFile(expectedFile)
			if err != nil {
				t.Fatalf("读取生成的文件失败: %v", err)
			}

			contentStr := string(content)

			// 验证是有效的 YAML 格式（至少不是空文件）
			if len(contentStr) == 0 {
				t.Errorf("YAML 文件内容为空")
			}
		}
	}
}

// TestModelNameToFileName 测试模型名转文件名
func TestModelNameToFileName(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     string
	}{
		{
			name:     "SimpleCamelCase",
			typeName: "GameTemplate",
			want:     "sign_template",
		},
		{
			name:     "MultipleWords",
			typeName: "HelperCryptoKey",
			want:     "helper_crypto_key",
		},
		{
			name:     "SingleWord",
			typeName: "User",
			want:     "user",
		},
		{
			name:     "AllCaps",
			typeName: "API",
			want:     "api",
		},
		{
			name:     "MixedCase",
			typeName: "SignGameType",
			want:     "sign_Game_type",
		},
		{
			name:     "WithNumbers",
			typeName: "Model2DB",
			want:     "model2_db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ModelNameToFileName(tt.typeName)
			if got != tt.want {
				t.Errorf("ModelNameToFileName(%q) = %q, want %q", tt.typeName, got, tt.want)
			}
		})
	}
}

// TestExtractPackageName 测试提取包名
func TestExtractPackageName(t *testing.T) {
	testDir := getTestDataDir()

	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			name:     "TestdataService",
			filePath: filepath.Join(testDir, "service", "service.go"),
			want:     "service",
		},
		{
			name:     "TestdataBusiness",
			filePath: filepath.Join(testDir, "business", "business.go"),
			want:     "business",
		},
		{
			name:     "TestdataLogic",
			filePath: filepath.Join(testDir, "logic", "template", "logic.go"),
			want:     "template",
		},
		{
			name:     "RootPath",
			filePath: "main.go",
			want:     ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractPackageName(tt.filePath)
			if got != tt.want {
				t.Errorf("ExtractPackageName(%q) = %q, want %q", tt.filePath, got, tt.want)
			}
		})
	}
}

// TestGenerateTypeName 测试生成类型名
func TestGenerateTypeName(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
		lineNumber  string
		want        string
	}{
		{
			name:        "NormalCase",
			packageName: "sign",
			lineNumber:  "115",
			want:        "SignLine115Mock",
		},
		{
			name:        "WithUnderscore",
			packageName: "sign_api",
			lineNumber:  "120",
			want:        "SignApiLine120Mock",
		},
		{
			name:        "SingleChar",
			packageName: "a",
			lineNumber:  "1",
			want:        "ALine1Mock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateTypeName(tt.packageName, tt.lineNumber)
			if got != tt.want {
				t.Errorf("GenerateTypeName(%q, %q) = %q, want %q", tt.packageName, tt.lineNumber, got, tt.want)
			}
		})
	}
}

// TestConvertFilePathToFileName 测试文件路径转文件名
func TestConvertFilePathToFileName(t *testing.T) {
	testDir := getTestDataDir()

	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			name:     "TestdataService",
			filePath: filepath.Join(testDir, "service", "service.go"),
			want:     "service_service_go",
		},
		{
			name:     "TestdataBusiness",
			filePath: filepath.Join(testDir, "business", "business.go"),
			want:     "business_business_go",
		},
		{
			name:     "TestdataLogic",
			filePath: filepath.Join(testDir, "logic", "template", "logic.go"),
			want:     "logic_template_logic_go",
		},
		{
			name:     "RootFile",
			filePath: "main.go",
			want:     "main_go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertFilePathToFileName(tt.filePath)
			if got != tt.want {
				t.Errorf("ConvertFilePathToFileName(%q) = %q, want %q", tt.filePath, got, tt.want)
			}
		})
	}
}

// TestGenerateCompleteMock 测试完整生成流程（Mock Go + YAML）
func TestGenerateCompleteMock(t *testing.T) {
	testDir := getTestDataDir()
	outputDir := filepath.Join(testDir, "output")

	// 清理输出目录
	os.RemoveAll(outputDir)
	defer os.RemoveAll(outputDir)

	// 使用 testdata 中的实际数据
	funcPath := filepath.Join(testDir, "business", "business.go") + ":9" // CreateGamePreCheck
	models := []dto.ModelInfo{
		{
			PackagePath: "testpkg/model/db/web_game",
			PackageName: "web_game",
			TypeName:    "GameTemplate",
			FullName:    "web_game.GameTemplate",
		},
	}

	// 生成完整 Mock 代码（包括 Go 文件和 YAML 文件）
	err := GenerateMockCode(testDir, funcPath, models, outputDir)
	if err != nil {
		t.Fatalf("生成 Mock 代码失败: %v", err)
	}

	// 验证 Mock Go 文件（输出到 testdata 内部）
	mockGoFile := filepath.Join(outputDir, "mock_modules/mock_business/business_business_go_9_mock.go")
	if _, err := os.Stat(mockGoFile); os.IsNotExist(err) {
		t.Errorf("期望生成 Mock Go 文件 %s，但文件不存在", mockGoFile)
	}

	// 验证 YAML fixture 文件（输出到 testdata 内部）
	yamlFiles := []string{
		filepath.Join(outputDir, "mock_modules/mock_business/testdata/fixtures/sign_template.yml"),
	}

	for _, yamlFile := range yamlFiles {
		if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
			t.Errorf("期望生成 YAML fixture 文件 %s，但文件不存在", yamlFile)
		} else {
			t.Logf("✓ YAML fixture 文件生成成功: %s", yamlFile)
		}
	}

	t.Logf("✓ 完整 Mock 代码生成成功")
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

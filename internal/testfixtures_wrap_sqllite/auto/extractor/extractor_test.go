package extractor

import (
	"path/filepath"
	"testing"
)

// TestExtractModels_FromVarDecl 测试从变量声明中提取模型
func TestExtractModels_FromVarDecl(t *testing.T) {
	testDir := getTestDataDir()
	funcPath := filepath.Join(testDir, "service", "service.go") + ":9" // TestVarDecl (第9行)
	t.Logf("source path: %v", funcPath)

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}
	t.Logf("all functions: %v", allFuncs)

	modelPaths := []string{"model/db", "testpkg/model/db"}
	models, err := ExtractModels(testDir, allFuncs, modelPaths)
	if err != nil {
		t.Fatalf("提取模型失败: %v", err)
	}

	// 验证：应该找到 GameTemplate 模型
	found := false
	for _, model := range models {
		if model.TypeName == "GameTemplate" && model.PackagePath == "testpkg/model/db/web_game" {
			found = true
			t.Logf("✓ 找到模型: %s", model.FullName)
			break
		}
	}

	if !found {
		t.Errorf("期望找到 GameTemplate 模型，但没有找到。实际找到: %v", models)
	}
}

// TestExtractModels_FromAssign 测试从赋值语句中提取模型
func TestExtractModels_FromAssign(t *testing.T) {
	testDir := getTestDataDir()
	funcPath := filepath.Join(testDir, "service", "service.go") + ":15" // TestAssignFromCall (第15行)

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}

	modelPaths := []string{"model/db", "testpkg/model/db"}
	models, err := ExtractModels(testDir, allFuncs, modelPaths)
	if err != nil {
		t.Fatalf("提取模型失败: %v", err)
	}

	// 验证：应该找到 GameTemplate 模型（从 getModel 的返回类型中提取）
	found := false
	for _, model := range models {
		if model.TypeName == "GameTemplate" && model.PackagePath == "testpkg/model/db/web_game" {
			found = true
			t.Logf("✓ 找到模型: %s", model.FullName)
			break
		}
	}

	if !found {
		t.Errorf("期望找到 GameTemplate 模型，但没有找到。实际找到: %v", models)
	}
}

// TestExtractModels_FromCompositeLit 测试从复合字面量中提取模型
func TestExtractModels_FromCompositeLit(t *testing.T) {
	testDir := getTestDataDir()
	funcPath := filepath.Join(testDir, "service", "service.go") + ":25" // TestCompositeLit (第25行)

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}

	modelPaths := []string{"model/db", "testpkg/model/db"}
	models, err := ExtractModels(testDir, allFuncs, modelPaths)
	if err != nil {
		t.Fatalf("提取模型失败: %v", err)
	}

	// 验证：应该找到 GameTemplate 模型
	found := false
	for _, model := range models {
		if model.TypeName == "GameTemplate" && model.PackagePath == "testpkg/model/db/web_game" {
			found = true
			t.Logf("✓ 找到模型: %s", model.FullName)
			break
		}
	}

	if !found {
		t.Errorf("期望找到 GameTemplate 模型，但没有找到。实际找到: %v", models)
	}
}

// TestExtractModels_FromTypeAssert 测试从类型断言中提取模型
func TestExtractModels_FromTypeAssert(t *testing.T) {
	testDir := getTestDataDir()
	funcPath := filepath.Join(testDir, "service", "service.go") + ":31" // TestTypeAssert (第31行)

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}

	modelPaths := []string{"model/db", "testpkg/model/db"}
	models, err := ExtractModels(testDir, allFuncs, modelPaths)
	if err != nil {
		t.Fatalf("提取模型失败: %v", err)
	}

	// 验证：应该找到 GameTemplate 模型
	found := false
	for _, model := range models {
		if model.TypeName == "GameTemplate" && model.PackagePath == "testpkg/model/db/web_game" {
			found = true
			t.Logf("✓ 找到模型: %s", model.FullName)
			break
		}
	}

	if !found {
		t.Errorf("期望找到 GameTemplate 模型，但没有找到。实际找到: %v", models)
	}
}

// TestExtractModels_FromReturnType 测试从函数返回类型中提取模型
func TestExtractModels_FromReturnType(t *testing.T) {
	testDir := getTestDataDir()
	funcPath := filepath.Join(testDir, "service", "service.go") + ":35" // TestReturnType (第35行)

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}

	modelPaths := []string{"model/db", "testpkg/model/db"}
	models, err := ExtractModels(testDir, allFuncs, modelPaths)
	if err != nil {
		t.Fatalf("提取模型失败: %v", err)
	}

	// 验证：应该找到 GameTemplate 模型
	found := false
	for _, model := range models {
		if model.TypeName == "GameTemplate" && model.PackagePath == "testpkg/model/db/web_game" {
			found = true
			t.Logf("✓ 找到模型: %s", model.FullName)
			break
		}
	}

	if !found {
		t.Errorf("期望找到 GameTemplate 模型，但没有找到。实际找到: %v", models)
	}
}

// TestExtractModels_FromSliceType 测试从切片类型中提取模型
func TestExtractModels_FromSliceType(t *testing.T) {
	testDir := getTestDataDir()
	funcPath := filepath.Join(testDir, "service", "service.go") + ":39" // TestSliceType (第39行)

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}

	modelPaths := []string{"model/db", "testpkg/model/db"}
	models, err := ExtractModels(testDir, allFuncs, modelPaths)
	if err != nil {
		t.Fatalf("提取模型失败: %v", err)
	}

	// 验证：应该找到 GameTemplate 模型
	found := false
	for _, model := range models {
		if model.TypeName == "GameTemplate" && model.PackagePath == "testpkg/model/db/web_game" {
			found = true
			t.Logf("✓ 找到模型: %s", model.FullName)
			break
		}
	}

	if !found {
		t.Errorf("期望找到 GameTemplate 模型，但没有找到。实际找到: %v", models)
	}
}

// TestExtractModels_MultipleModels 测试多个模型类型（去重）
func TestExtractModels_MultipleModels(t *testing.T) {
	testDir := getTestDataDir()
	funcPath := filepath.Join(testDir, "service", "service.go") + ":43" // TestMultipleModels (第43行)

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}

	modelPaths := []string{"model/db", "testpkg/model/db"}
	models, err := ExtractModels(testDir, allFuncs, modelPaths)
	if err != nil {
		t.Fatalf("提取模型失败: %v", err)
	}

	// 验证：应该只找到1个 GameTemplate 模型（去重后）
	count := 0
	for _, model := range models {
		if model.TypeName == "GameTemplate" && model.PackagePath == "testpkg/model/db/web_game" {
			count++
			t.Logf("✓ 找到模型: %s", model.FullName)
		}
	}

	if count == 0 {
		t.Errorf("期望找到 GameTemplate 模型，但没有找到。实际找到: %v", models)
	}
	if count > 1 {
		t.Errorf("期望只找到1个 GameTemplate 模型（去重后），但找到 %d 个", count)
	}
}

// TestExtractModels_RealWorldCase 测试真实场景：从 business 层追踪到 logic 层，提取所有模型
func TestExtractModels_RealWorldCase(t *testing.T) {
	testDir := getTestDataDir()
	funcPath := filepath.Join(testDir, "business", "business.go") + ":10" // CreateGamePreCheck

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}

	// 验证：应该找到至少1个函数（CreateGamePreCheck）
	if len(allFuncs) < 1 {
		t.Errorf("期望找到至少1个函数，但只找到 %d 个", len(allFuncs))
	}

	// 打印找到的函数，用于调试
	for _, fn := range allFuncs {
		t.Logf("找到函数: %s", fn.Key())
	}

	modelPaths := []string{"model/db", "testpkg/model/db"}
	models, err := ExtractModels(testDir, allFuncs, modelPaths)
	if err != nil {
		t.Fatalf("提取模型失败: %v", err)
	}

	// 验证：应该找到 GameTemplate 模型（从 GetTmplDetail 的返回类型中提取）
	found := false
	for _, model := range models {
		if model.TypeName == "GameTemplate" && model.PackagePath == "testpkg/model/db/web_game" {
			found = true
			t.Logf("✓ 找到模型: %s", model.FullName)
			break
		}
	}

	if !found {
		t.Errorf("期望找到 GameTemplate 模型，但没有找到。实际找到: %v", models)
	}

	t.Logf("找到 %d 个模型: %v", len(models), models)
}

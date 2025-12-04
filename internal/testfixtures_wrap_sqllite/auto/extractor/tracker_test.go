package extractor

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/packages"
)

// getTestDataDir 获取测试数据目录
func getTestDataDir() string {
	return filepath.Join(".", "testdata")
}

// TestTrackFunctionCalls_NestedSelector 测试嵌套选择器表达式：st.GameTemplateLogic.GetTmplDetail
func TestTrackFunctionCalls_NestedSelector(t *testing.T) {
	testDir := getTestDataDir()

	// 测试：从 business.go:10 (CreateGamePreCheck) 开始追踪
	funcPath := filepath.Join(testDir, "business", "business.go") + ":10"

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}

	// 验证：应该找到以下函数
	expectedFuncs := []string{
		"testpkg/business.GameTemplate2Bus.CreateGamePreCheck",
		"testpkg/logic/template.GameTemplateLogic.GetTmplDetail",
		"testpkg/logic/template.Repo.Find",
	}

	funcMap := make(map[string]bool)
	for _, fn := range allFuncs {
		key := fn.Key()
		funcMap[key] = true
		t.Logf("找到函数: %s", key)
	}

	for _, expected := range expectedFuncs {
		if !funcMap[expected] {
			t.Errorf("期望找到函数 %s，但没有找到。实际找到: %v", expected, funcMap)
		}
	}

	// 验证：应该找到至少3个函数
	if len(allFuncs) < 3 {
		t.Errorf("期望找到至少3个函数，但只找到 %d 个", len(allFuncs))
	}
}

// TestTrackFunctionCalls_SimpleCall 测试简单函数调用
func TestTrackFunctionCalls_SimpleCall(t *testing.T) {
	testDir := getTestDataDir()

	// 测试：从 logic.go:10 (GetTmplDetail) 开始追踪
	funcPath := filepath.Join(testDir, "logic", "template", "logic.go") + ":10"

	allFuncs, err := TrackFunctionCalls(testDir, funcPath)
	if err != nil {
		t.Fatalf("追踪函数调用失败: %v", err)
	}

	// 验证：应该找到以下函数
	expectedFuncs := []string{
		"testpkg/logic/template.GameTemplateLogic.GetTmplDetail",
		"testpkg/logic/template.Repo.Find",
	}

	funcMap := make(map[string]bool)
	for _, fn := range allFuncs {
		key := fn.Key()
		funcMap[key] = true
		t.Logf("找到函数: %s", key)
	}

	for _, expected := range expectedFuncs {
		if !funcMap[expected] {
			t.Errorf("期望找到函数 %s，但没有找到。实际找到: %v", expected, funcMap)
		}
	}
}

// TestExtractPackageAndType 测试 extractPackageAndType
func TestExtractPackageAndType(t *testing.T) {
	testCases := []struct {
		input    string
		pkgPath  string
		typeName string
	}{
		{"*testpkg/logic/template.GameTemplateLogic", "testpkg/logic/template", "GameTemplateLogic"},
		{"testpkg/logic/template.GameTemplateLogic", "testpkg/logic/template", "GameTemplateLogic"},
		{"*testpkg/business.GameTemplate2Bus", "testpkg/business", "GameTemplate2Bus"},
		{"[]*testpkg/model/db/web_game.GameTemplate", "testpkg/model/db/web_game", "GameTemplate"},
		{"[]testpkg/model/db/web_game.GameTemplate", "testpkg/model/db/web_game", "GameTemplate"},
		{"[]*github.com/web_game.GameTemplate", "github.com/web_game", "GameTemplate"},
	}

	for _, tc := range testCases {
		pkgPath, typeName := extractPackageAndType(tc.input)
		if pkgPath != tc.pkgPath || typeName != tc.typeName {
			t.Errorf("输入: %s, 期望: (%s, %s), 实际: (%s, %s)", tc.input, tc.pkgPath, tc.typeName, pkgPath, typeName)
		}
	}
}

// TestFindMethodInPackage 测试 findMethodInPackage
func TestFindMethodInPackage(t *testing.T) {
	testDir := getTestDataDir()

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports,
		Dir:  testDir,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		t.Fatalf("加载包失败: %v", err)
	}

	// 建立包映射
	pkgMap := make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		pkgMap[pkg.PkgPath] = pkg
	}

	// 测试查找方法
	pkgPath := "testpkg/logic/template"
	typeName := "GameTemplateLogic"
	methodName := "GetTmplDetail"

	foundFunc, foundPkg := findMethodInPackage(pkgPath, typeName, methodName, pkgMap)
	if foundFunc == nil || foundPkg == nil {
		t.Errorf("未找到方法 %s.%s.%s", pkgPath, typeName, methodName)
	} else {
		t.Logf("✓ 找到方法: %s.%s.%s", foundPkg.PkgPath, typeName, methodName)
	}
}

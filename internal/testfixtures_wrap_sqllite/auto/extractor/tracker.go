package extractor

import (
	"fmt"
	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite/auto/dto"

	"go/ast"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// TrackFunctionCalls 递归追踪函数调用链，返回所有函数
func TrackFunctionCalls(workDir, funcPath string) ([]dto.FunctionInfo, error) {
	// 解析函数路径
	target, err := parseFunctionPath(funcPath)
	if err != nil {
		return nil, fmt.Errorf("解析函数路径失败: %w", err)
	}

	// 加载包
	pkgs, err := loadPackages(workDir)
	if err != nil {
		return nil, fmt.Errorf("加载包失败: %w", err)
	}

	// 建立包路径到包的映射
	pkgMap := make(map[string]*packages.Package)
	var modulePath string
	for _, pkg := range pkgs {
		pkgMap[pkg.PkgPath] = pkg
		// 从第一个包中提取模块路径（假设所有包都在同一个模块中）
		if modulePath == "" && pkg.Module != nil {
			modulePath = pkg.Module.Path
		}
	}

	// 找到目标函数
	targetFunc, targetPkg, err := findTargetFunction(pkgs, target)
	if err != nil {
		return nil, fmt.Errorf("查找目标函数失败: %w", err)
	}

	// 递归追踪
	visited := make(map[string]bool)
	allFuncs := []dto.FunctionInfo{}

	var track func(*ast.FuncDecl, *packages.Package)
	track = func(fn *ast.FuncDecl, pkg *packages.Package) {
		funcInfo := functionToInfo(fn, pkg)
		key := funcInfo.Key()

		if visited[key] {
			return
		}
		visited[key] = true
		allFuncs = append(allFuncs, funcInfo)

		// 收集变量类型映射
		varMap := collectVariableTypes(fn, pkg)

		// 查找所有函数调用
		ast.Inspect(fn, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			calledFunc, calledPkg := findCalledFunction(call, pkg, varMap, fn, pkgMap, modulePath)
			if calledFunc != nil && calledPkg != nil {
				// 只追踪项目内的函数
				if isProjectPackage(calledPkg.PkgPath, modulePath) {
					track(calledFunc, calledPkg)
				}
			}

			return true
		})
	}

	track(targetFunc, targetPkg)

	return allFuncs, nil
}

// functionToInfo 将函数转换为函数信息
func functionToInfo(fn *ast.FuncDecl, pkg *packages.Package) dto.FunctionInfo {
	info := dto.FunctionInfo{
		PackagePath: pkg.PkgPath,
	}

	// 获取文件路径和行号
	if len(pkg.GoFiles) > 0 {
		info.File = pkg.GoFiles[0]
		info.Line = int(pkg.Fset.Position(fn.Pos()).Line)
	}

	// 判断是方法还是函数
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		// 是方法
		recvType := fn.Recv.List[0].Type
		info.TypeName = extractTypeName(recvType)
		info.MethodName = fn.Name.Name
	} else {
		// 是包级函数
		info.FuncName = fn.Name.Name
	}

	return info
}

// parseFunctionPath 解析函数路径
func parseFunctionPath(funcPath string) (target, error) {
	// 支持格式: filepath:line
	if strings.Contains(funcPath, ":") && (strings.Contains(funcPath, "/") || strings.Contains(funcPath, "\\")) {
		parts := strings.SplitN(funcPath, ":", 2)
		if len(parts) == 2 {
			return target{
				Type:     "file",
				FilePath: parts[0],
				Line:     parts[1],
			}, nil
		}
	}

	// 支持格式: package.Function 或 Type.Method
	parts := strings.Split(funcPath, ".")
	if len(parts) == 2 {
		return target{
			Type:     "function",
			Package:  parts[0],
			Function: parts[1],
		}, nil
	}
	if len(parts) == 3 {
		return target{
			Type:     "method",
			TypeName: parts[1],
			Method:   parts[2],
		}, nil
	}

	return target{}, fmt.Errorf("不支持的函数路径格式: %s", funcPath)
}

// target 目标函数信息
type target struct {
	Type     string // "file", "function", "method"
	FilePath string
	Line     string
	Package  string
	Function string
	TypeName string
	Method   string
}

// loadPackages 加载包
func loadPackages(workDir string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports,
		Dir:  workDir,
	}

	return packages.Load(cfg, "./...")
}

// findTargetFunction 查找目标函数
func findTargetFunction(pkgs []*packages.Package, t target) (*ast.FuncDecl, *packages.Package, error) {
	switch t.Type {
	case "file":
		return findFunctionByFileLine(pkgs, t.FilePath, t.Line)
	case "function":
		return findFunctionByName(pkgs, t.Package, t.Function)
	case "method":
		return findMethodByName(pkgs, t.TypeName, t.Method)
	default:
		return nil, nil, fmt.Errorf("不支持的目标类型: %s", t.Type)
	}
}

// findFunctionByFileLine 通过文件路径和行号查找函数
func findFunctionByFileLine(pkgs []*packages.Package, filePath, lineStr string) (*ast.FuncDecl, *packages.Package, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, nil, err
	}

	var line int
	fmt.Sscanf(lineStr, "%d", &line)

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			pkgFilePath := pkg.Fset.File(file.Pos()).Name()
			pkgAbsPath, err := filepath.Abs(pkgFilePath)
			if err != nil {
				continue
			}
			if pkgAbsPath == absPath {
				var targetFunc *ast.FuncDecl
				ast.Inspect(file, func(n ast.Node) bool {
					if fn, ok := n.(*ast.FuncDecl); ok {
						fnLine := pkg.Fset.Position(fn.Pos()).Line
						if fnLine <= line && line <= pkg.Fset.Position(fn.End()).Line {
							targetFunc = fn
							return false
						}
					}
					return true
				})
				if targetFunc != nil {
					return targetFunc, pkg, nil
				}
			}
		}
	}

	return nil, nil, fmt.Errorf("未找到函数: %s:%s", filePath, lineStr)
}

// findFunctionByName 通过名称查找包级函数
func findFunctionByName(pkgs []*packages.Package, pkgPath, funcName string) (*ast.FuncDecl, *packages.Package, error) {
	for _, pkg := range pkgs {
		if pkg.PkgPath != pkgPath {
			continue
		}
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok {
					if fn.Name.Name == funcName && fn.Recv == nil {
						return false
					}
				}
				return true
			})
		}
	}
	return nil, nil, fmt.Errorf("未找到函数: %s.%s", pkgPath, funcName)
}

// findMethodByName 通过名称查找方法
func findMethodByName(pkgs []*packages.Package, typeName, methodName string) (*ast.FuncDecl, *packages.Package, error) {
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok {
					if fn.Name.Name == methodName && fn.Recv != nil {
						recvTypeName := extractTypeName(fn.Recv.List[0].Type)
						if recvTypeName == typeName {
							return false
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil, fmt.Errorf("未找到方法: %s.%s", typeName, methodName)
}

// extractTypeName 提取类型名
func extractTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return extractTypeName(t.X)
	default:
		return ""
	}
}

// isProjectPackage 判断是否是项目内的包
func isProjectPackage(pkgPath, modulePath string) bool {
	// 如果提供了模块路径，检查包路径是否以模块路径开头
	if modulePath != "" {
		if strings.HasPrefix(pkgPath, modulePath) {
			return true
		}
		// 如果包路径是模块路径的子路径（如 testpkg/logic/template 对于模块 testpkg）
		if strings.HasPrefix(pkgPath, modulePath+"/") {
			return true
		}
	}

	// 排除标准库（不包含点号或斜杠的包路径，如 "fmt", "os" 等）
	if !strings.Contains(pkgPath, ".") && !strings.Contains(pkgPath, "/") {
		return false
	}

	// 排除常见的第三方库路径
	excludePrefixes := []string{
		"github.com",
		"golang.org",
		"google.golang.org",
		"gopkg.in",
	}
	for _, prefix := range excludePrefixes {
		if strings.HasPrefix(pkgPath, prefix) {
			return false
		}
	}

	// 如果包路径包含斜杠，可能是项目内的包（如 testpkg/logic/template）
	if strings.Contains(pkgPath, "/") {
		return true
	}

	return false
}

// collectVariableTypes 收集变量类型映射
func collectVariableTypes(fn *ast.FuncDecl, pkg *packages.Package) map[string]string {
	varMap := make(map[string]string)
	if pkg == nil || pkg.TypesInfo == nil {
		return varMap
	}

	ast.Inspect(fn, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for i, lhs := range assign.Lhs {
			if i >= len(assign.Rhs) {
				continue
			}

			ident, ok := lhs.(*ast.Ident)
			if !ok {
				continue
			}

			// 从右侧表达式获取类型
			if call, ok := assign.Rhs[i].(*ast.CallExpr); ok {
				if callType := pkg.TypesInfo.TypeOf(call); callType != nil {
					// 处理多返回值（Tuple）
					if tuple, ok := callType.(*types.Tuple); ok && tuple.Len() > i {
						varMap[ident.Name] = tuple.At(i).Type().String()
					} else if i == 0 {
						// 单返回值
						varMap[ident.Name] = callType.String()
					}
				}
			}
		}

		return true
	})

	return varMap
}

// findCalledFunction 查找被调用的函数
func findCalledFunction(call *ast.CallExpr, pkg *packages.Package, varMap map[string]string, fn *ast.FuncDecl, pkgMap map[string]*packages.Package, modulePath string) (*ast.FuncDecl, *packages.Package) {
	if pkg == nil || pkg.TypesInfo == nil {
		return nil, nil
	}

	var sel *ast.SelectorExpr
	var ident *ast.Ident

	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		sel = fun
		ident = fun.Sel
	case *ast.Ident:
		ident = fun
	default:
		return nil, nil
	}

	if ident == nil {
		return nil, nil
	}

	// 处理选择器表达式（方法调用或包函数调用）
	if sel != nil {
		// 检查是否是包函数调用（pkg.Function）
		if x, ok := sel.X.(*ast.Ident); ok {
			if obj := pkg.TypesInfo.ObjectOf(x); obj != nil {
				if pkgName, ok := obj.(*types.PkgName); ok {
					importedPkg := pkgName.Imported()
					if importedPkg == nil {
						return nil, nil
					}
					pkgPath := importedPkg.Path()
					if !isProjectPackage(pkgPath, modulePath) {
						return nil, nil
					}
					// 查找该包中的函数
					return findFunctionInPackage(pkgPath, ident.Name, pkgMap)
				}
			}
		}

		// 检查是否是方法调用（obj.Method 或 st.Field.Method）
		// 处理嵌套选择器表达式：st.GameTemplateLogic.GetTmplDetail

		// 方法1: 直接使用 TypeOf 获取 sel.X 的类型（对于嵌套选择器，应该返回字段类型）
		var exprType string
		if pkg.TypesInfo != nil {
			if selType := pkg.TypesInfo.TypeOf(sel.X); selType != nil {
				exprType = selType.String()
			}
		}

		// 方法2: 如果 TypeOf 失败，使用 resolveExprTypeForMethodCall
		if exprType == "" {
			exprType = resolveExprTypeForMethodCall(sel.X, pkg, varMap, fn, pkgMap)
		}

		if exprType != "" {
			// 从类型中提取包路径和类型名
			pkgPath, typeName := extractPackageAndType(exprType)
			if pkgPath != "" && typeName != "" {
				if isProjectPackage(pkgPath, modulePath) {
					foundFunc, foundPkg := findMethodInPackage(pkgPath, typeName, ident.Name, pkgMap)
					if foundFunc != nil && foundPkg != nil {
						return foundFunc, foundPkg
					}
				}
			}
		}
	} else {
		// 直接函数调用
		if obj := pkg.TypesInfo.ObjectOf(ident); obj != nil {
			if _, ok := obj.(*types.Func); ok {
				// 在当前包中查找函数定义
				return findFunctionInPackage(pkg.PkgPath, ident.Name, pkgMap)
			}
		}
	}

	return nil, nil
}

// resolveExprTypeForMethodCall 解析方法调用表达式的类型（支持嵌套选择器）
func resolveExprTypeForMethodCall(expr ast.Expr, pkg *packages.Package, varMap map[string]string, fn *ast.FuncDecl, pkgMap map[string]*packages.Package) string {
	// 方法1: 直接使用 TypeOf（最可靠，支持嵌套选择器）
	if pkg.TypesInfo != nil {
		if exprType := pkg.TypesInfo.TypeOf(expr); exprType != nil {
			return exprType.String()
		}
	}

	// 方法2: 处理嵌套选择器表达式（如 st.GameTemplateLogic）
	if nestedSel, ok := expr.(*ast.SelectorExpr); ok {
		// 先尝试直接使用 TypeOf 获取嵌套选择器的类型（字段类型）
		if pkg.TypesInfo != nil {
			if fieldType := pkg.TypesInfo.TypeOf(nestedSel); fieldType != nil {
				return fieldType.String()
			}
		}

		// 如果 TypeOf 失败，递归解析嵌套选择器的基础部分
		baseType := resolveExprTypeForMethodCall(nestedSel.X, pkg, varMap, fn, pkgMap)
		if baseType != "" {
			// 从基础类型中查找字段类型
			fieldType := getFieldTypeFromTypeString(baseType, nestedSel.Sel.Name, pkg, pkgMap)
			if fieldType != "" {
				return fieldType
			}
		}
	}

	// 方法3: 处理标识符（变量或接收者）
	if ident, ok := expr.(*ast.Ident); ok {
		// 从 varMap 查找
		if typeStr, ok := varMap[ident.Name]; ok {
			return typeStr
		}

		// 从类型信息获取
		if obj := pkg.TypesInfo.ObjectOf(ident); obj != nil {
			if varObj, ok := obj.(*types.Var); ok {
				return varObj.Type().String()
			}
		}

		// 可能是方法接收者
		if fn != nil && fn.Recv != nil && len(fn.Recv.List) > 0 {
			recv := fn.Recv.List[0]
			if len(recv.Names) > 0 && recv.Names[0].Name == ident.Name {
				if recvType := pkg.TypesInfo.TypeOf(recv.Type); recvType != nil {
					return recvType.String()
				}
			}
		}
	}

	return ""
}

// getFieldTypeFromTypeString 从类型字符串中查找字段类型
func getFieldTypeFromTypeString(typeStr, fieldName string, pkg *packages.Package, pkgMap map[string]*packages.Package) string {
	if pkg == nil || pkg.TypesInfo == nil {
		return ""
	}

	// 从类型字符串中提取包路径和类型名
	pkgPath, typeName := extractPackageAndType(typeStr)
	if pkgPath == "" || typeName == "" {
		return ""
	}

	// 在当前包中查找类型
	var actualType types.Type
	for _, obj := range pkg.TypesInfo.Defs {
		if typeNameObj, ok := obj.(*types.TypeName); ok {
			if typeNameObj.Name() == typeName {
				actualType = typeNameObj.Type()
				break
			}
		}
	}

	// 如果找不到，尝试从导入的包中查找
	if actualType == nil {
		for _, file := range pkg.Syntax {
			for _, imp := range file.Imports {
				importPath := strings.Trim(imp.Path.Value, "\"")
				if importPath == pkgPath {
					// 查找导入的包名
					var importedPkg *types.Package
					if imp.Name != nil {
						// 有别名
						if obj := pkg.TypesInfo.ObjectOf(imp.Name); obj != nil {
							if pkgName, ok := obj.(*types.PkgName); ok {
								importedPkg = pkgName.Imported()
							}
						}
					} else {
						// 无别名，从 Uses 中查找
						for _, obj := range pkg.TypesInfo.Uses {
							if pkgName, ok := obj.(*types.PkgName); ok {
								if pkgName.Imported().Path() == pkgPath {
									importedPkg = pkgName.Imported()
									break
								}
							}
						}
					}

					if importedPkg != nil {
						obj := importedPkg.Scope().Lookup(typeName)
						if obj != nil {
							if typeNameObj, ok := obj.(*types.TypeName); ok {
								actualType = typeNameObj.Type()
								break
							}
						}
					}
				}
			}
			if actualType != nil {
				break
			}
		}
	}

	if actualType == nil {
		return ""
	}

	// 处理指针类型
	if ptrType, ok := actualType.(*types.Pointer); ok {
		actualType = ptrType.Elem()
	}

	// 处理命名类型
	namedType, ok := actualType.(*types.Named)
	if !ok {
		return ""
	}

	// 获取结构体类型
	structType, ok := namedType.Underlying().(*types.Struct)
	if !ok {
		return ""
	}

	// 查找字段
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if field.Name() == fieldName {
			return field.Type().String()
		}
	}

	return ""
}

// extractPackageAndType 从类型字符串中提取包路径和类型名
func extractPackageAndType(typeStr string) (string, string) {
	// 处理切片类型 []pkg.Type 或 []*pkg.Type
	if strings.HasPrefix(typeStr, "[]") {
		typeStr = typeStr[2:]
	}

	// 处理指针类型 *pkg.Type
	if strings.HasPrefix(typeStr, "*") {
		typeStr = typeStr[1:]
	}

	// 查找最后一个点，分割包路径和类型名
	lastDot := strings.LastIndex(typeStr, ".")
	if lastDot == -1 {
		return "", ""
	}

	pkgPath := typeStr[:lastDot]
	typeName := typeStr[lastDot+1:]

	return pkgPath, typeName
}

// findFunctionInPackage 在包中查找函数
func findFunctionInPackage(pkgPath, funcName string, pkgMap map[string]*packages.Package) (*ast.FuncDecl, *packages.Package) {
	pkg, ok := pkgMap[pkgPath]
	if !ok {
		return nil, nil
	}

	var foundFunc *ast.FuncDecl
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			if foundFunc != nil {
				return false
			}
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Name.Name == funcName && fn.Recv == nil {
					foundFunc = fn
					return false
				}
			}
			return true
		})
		if foundFunc != nil {
			break
		}
	}

	if foundFunc != nil {
		return foundFunc, pkg
	}
	return nil, nil
}

// findMethodInPackage 在包中查找方法
func findMethodInPackage(pkgPath, typeName, methodName string, pkgMap map[string]*packages.Package) (*ast.FuncDecl, *packages.Package) {
	pkg, ok := pkgMap[pkgPath]
	if !ok {
		return nil, nil
	}

	var foundFunc *ast.FuncDecl
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			if foundFunc != nil {
				return false
			}
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Name.Name == methodName && fn.Recv != nil && len(fn.Recv.List) > 0 {
					recvTypeName := extractTypeName(fn.Recv.List[0].Type)
					// 支持类型名匹配（去掉指针符号后比较）
					if recvTypeName == typeName ||
						(strings.HasPrefix(typeName, "*") && recvTypeName == typeName[1:]) ||
						(strings.HasPrefix(recvTypeName, "*") && recvTypeName[1:] == typeName) {
						foundFunc = fn
						return false
					}
				}
			}
			return true
		})
		if foundFunc != nil {
			break
		}
	}

	if foundFunc != nil {
		return foundFunc, pkg
	}
	return nil, nil
}

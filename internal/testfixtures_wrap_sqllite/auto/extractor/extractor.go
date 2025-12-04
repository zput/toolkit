package extractor

import (
	"fmt"
	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite/auto/dto"

	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ExtractModels 从函数列表中提取模型
func ExtractModels(workDir string, funcs []dto.FunctionInfo, modelPaths []string) ([]dto.ModelInfo, error) {
	if len(funcs) == 0 {
		return []dto.ModelInfo{}, nil
	}

	// 加载包（复用已加载的包）
	pkgs, err := loadPackages(workDir)
	if err != nil {
		return nil, fmt.Errorf("加载包失败: %w", err)
	}

	// 建立包路径到包的映射
	pkgMap := make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		pkgMap[pkg.PkgPath] = pkg
	}

	// 收集所有模型
	allModels := make(map[string]*dto.ModelInfo) // key: PackagePath.TypeName

	for _, fnInfo := range funcs {
		// 找到函数定义
		pkg, ok := pkgMap[fnInfo.PackagePath]
		if !ok {
			continue
		}

		// 找到函数 AST
		var fn *ast.FuncDecl
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				if f, ok := n.(*ast.FuncDecl); ok {
					if fnInfo.TypeName != "" {
						// 是方法
						if f.Name.Name == fnInfo.MethodName {
							if f.Recv != nil && len(f.Recv.List) > 0 {
								recvTypeName := extractTypeName(f.Recv.List[0].Type)
								if recvTypeName == fnInfo.TypeName {
									fn = f
									return false
								}
							}
						}
					} else {
						// 是包级函数
						if f.Name.Name == fnInfo.FuncName {
							fn = f
							return false
						}
					}
				}
				return true
			})
			if fn != nil {
				break
			}
		}

		if fn == nil {
			continue
		}

		// 从函数中提取模型
		models := extractModelsFromFunc(fn, pkg, modelPaths)
		for _, model := range models {
			key := model.PackagePath + "." + model.TypeName
			allModels[key] = model
		}
	}

	// 转换为切片
	result := make([]dto.ModelInfo, 0, len(allModels))
	for _, model := range allModels {
		result = append(result, *model)
	}

	return result, nil
}

// extractModelsFromFunc 从单个函数中提取模型
func extractModelsFromFunc(fn *ast.FuncDecl, pkg *packages.Package, modelPaths []string) []*dto.ModelInfo {
	var models []*dto.ModelInfo
	seen := make(map[string]bool)

	if pkg == nil || pkg.TypesInfo == nil {
		return models
	}

	// 收集变量类型映射
	varMap := collectVariableTypesForExtraction(fn, pkg)

	// 遍历函数体，查找所有类型引用
	ast.Inspect(fn, func(n ast.Node) bool {
		// 1. 检查变量声明中的类型
		if genDecl, ok := n.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					if valueSpec.Type != nil {
						model := extractModelFromTypeExpr(valueSpec.Type, pkg, modelPaths)
						if model != nil {
							key := model.PackagePath + "." + model.TypeName
							if !seen[key] {
								models = append(models, model)
								seen[key] = true
							}
						}
					}
				}
			}
		}

		// 2. 检查赋值语句中的类型
		if assign, ok := n.(*ast.AssignStmt); ok {
			// 检查右侧函数调用的返回类型
			for i, rhs := range assign.Rhs {
				if call, ok := rhs.(*ast.CallExpr); ok {
					if pkg.TypesInfo != nil {
						if callType := pkg.TypesInfo.TypeOf(call); callType != nil {
							if tuple, ok := callType.(*types.Tuple); ok && tuple.Len() > i {
								// 多返回值，取对应位置的类型
								model := extractModelFromType(tuple.At(i).Type(), pkg, modelPaths)
								if model != nil {
									key := model.PackagePath + "." + model.TypeName
									if !seen[key] {
										models = append(models, model)
										seen[key] = true
									}
								}
							} else if i == 0 {
								// 单返回值
								model := extractModelFromType(callType, pkg, modelPaths)
								if model != nil {
									key := model.PackagePath + "." + model.TypeName
									if !seen[key] {
										models = append(models, model)
										seen[key] = true
									}
								}
							}
						}
					}
				}
				// 检查复合字面量
				if compLit, ok := rhs.(*ast.CompositeLit); ok {
					model := extractModelFromTypeExpr(compLit.Type, pkg, modelPaths)
					if model != nil {
						key := model.PackagePath + "." + model.TypeName
						if !seen[key] {
							models = append(models, model)
							seen[key] = true
						}
					}
				}
			}
			// 检查左侧变量的类型（从 varMap 获取）
			for _, lhs := range assign.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if typeStr, ok := varMap[ident.Name]; ok {
						model := extractModelFromTypeString(typeStr, pkg, modelPaths)
						if model != nil {
							key := model.PackagePath + "." + model.TypeName
							if !seen[key] {
								models = append(models, model)
								seen[key] = true
							}
						}
					}
				}
			}
		}

		// 3. 检查类型断言
		if typeAssert, ok := n.(*ast.TypeAssertExpr); ok {
			model := extractModelFromTypeExpr(typeAssert.Type, pkg, modelPaths)
			if model != nil {
				key := model.PackagePath + "." + model.TypeName
				if !seen[key] {
					models = append(models, model)
					seen[key] = true
				}
			}
		}

		// 4. 检查返回语句中的类型
		if retStmt, ok := n.(*ast.ReturnStmt); ok {
			if pkg.TypesInfo != nil {
				for _, result := range retStmt.Results {
					if resultType := pkg.TypesInfo.TypeOf(result); resultType != nil {
						model := extractModelFromType(resultType, pkg, modelPaths)
						if model != nil {
							key := model.PackagePath + "." + model.TypeName
							if !seen[key] {
								models = append(models, model)
								seen[key] = true
							}
						}
					}
				}
			}
		}

		return true
	})

	// 5. 检查函数返回类型
	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			if field.Type != nil {
				model := extractModelFromTypeExpr(field.Type, pkg, modelPaths)
				if model != nil {
					key := model.PackagePath + "." + model.TypeName
					if !seen[key] {
						models = append(models, model)
						seen[key] = true
					}
				}
			}
		}
	}

	return models
}

// collectVariableTypesForExtraction 收集变量类型映射（用于模型提取）
func collectVariableTypesForExtraction(fn *ast.FuncDecl, pkg *packages.Package) map[string]string {
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
					if tuple, ok := callType.(*types.Tuple); ok && tuple.Len() > i {
						// 多返回值，取对应位置的类型
						varMap[ident.Name] = tuple.At(i).Type().String()
					} else if i == 0 {
						// 单返回值
						varMap[ident.Name] = callType.String()
					}
				}
			} else if pkg.TypesInfo != nil {
				if rhsType := pkg.TypesInfo.TypeOf(assign.Rhs[i]); rhsType != nil {
					varMap[ident.Name] = rhsType.String()
				}
			}
		}

		return true
	})

	return varMap
}

// extractModelFromTypeExpr 从 AST 类型表达式中提取模型
func extractModelFromTypeExpr(typeExpr ast.Expr, pkg *packages.Package, modelPaths []string) *dto.ModelInfo {
	if typeExpr == nil {
		return nil
	}

	// 处理标识符类型（如 Type）
	if ident, ok := typeExpr.(*ast.Ident); ok {
		if pkg != nil && pkg.TypesInfo != nil {
			if obj := pkg.TypesInfo.ObjectOf(ident); obj != nil {
				if typeObj, ok := obj.(*types.TypeName); ok {
					if namedType, ok := typeObj.Type().(*types.Named); ok {
						return extractModelFromType(namedType, pkg, modelPaths)
					}
				}
			}
		}
		return nil
	}

	// 处理选择器表达式（如 pkg.Type）
	if sel, ok := typeExpr.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			if pkg != nil && pkg.TypesInfo != nil {
				if obj := pkg.TypesInfo.ObjectOf(ident); obj != nil {
					if pkgName, ok := obj.(*types.PkgName); ok {
						importedPkg := pkgName.Imported()
						if importedPkg != nil {
							pkgPath := importedPkg.Path()
							if isModelPath(pkgPath, modelPaths) {
								typeName := sel.Sel.Name
								pkgName := getPackageName(pkgPath)
								return &dto.ModelInfo{
									PackagePath: pkgPath,
									PackageName: pkgName,
									TypeName:    typeName,
									FullName:    pkgName + "." + typeName,
								}
							}
						}
					}
				}
			}
		}
		return nil
	}

	// 处理指针类型（如 *Type）
	if star, ok := typeExpr.(*ast.StarExpr); ok {
		return extractModelFromTypeExpr(star.X, pkg, modelPaths)
	}

	// 处理数组/切片类型（如 []Type）
	if array, ok := typeExpr.(*ast.ArrayType); ok {
		return extractModelFromTypeExpr(array.Elt, pkg, modelPaths)
	}

	return nil
}

// extractModelFromType 从 types.Type 中提取模型
func extractModelFromType(typ types.Type, pkg *packages.Package, modelPaths []string) *dto.ModelInfo {
	if typ == nil {
		return nil
	}

	// 处理指针类型
	if ptrType, ok := typ.(*types.Pointer); ok {
		return extractModelFromType(ptrType.Elem(), pkg, modelPaths)
	}

	// 处理切片类型
	if sliceType, ok := typ.(*types.Slice); ok {
		return extractModelFromType(sliceType.Elem(), pkg, modelPaths)
	}

	// 处理数组类型
	if arrayType, ok := typ.(*types.Array); ok {
		return extractModelFromType(arrayType.Elem(), pkg, modelPaths)
	}

	// 处理命名类型
	if namedType, ok := typ.(*types.Named); ok {
		obj := namedType.Obj()
		if obj == nil {
			return nil
		}

		pkg := obj.Pkg()
		if pkg == nil {
			return nil
		}

		pkgPath := pkg.Path()
		if !isModelPath(pkgPath, modelPaths) {
			return nil
		}

		pkgName := getPackageName(pkgPath)
		return &dto.ModelInfo{
			PackagePath: pkgPath,
			PackageName: pkgName,
			TypeName:    obj.Name(),
			FullName:    pkgName + "." + obj.Name(),
		}
	}

	return nil
}

// extractModelFromTypeString 从类型字符串中提取模型
func extractModelFromTypeString(typeStr string, pkg *packages.Package, modelPaths []string) *dto.ModelInfo {
	// 提取包路径和类型名
	pkgPath, typeName := extractPackageAndType(typeStr)
	if pkgPath == "" || typeName == "" {
		return nil
	}

	if !isModelPath(pkgPath, modelPaths) {
		return nil
	}

	pkgName := getPackageName(pkgPath)
	return &dto.ModelInfo{
		PackagePath: pkgPath,
		PackageName: pkgName,
		TypeName:    typeName,
		FullName:    pkgName + "." + typeName,
	}
}

// isModelPath 检查包路径是否匹配模型路径模式
func isModelPath(pkgPath string, modelPaths []string) bool {
	for _, pattern := range modelPaths {
		if strings.Contains(pkgPath, pattern) {
			return true
		}
	}
	return false
}

// getPackageName 从包路径中提取包名
func getPackageName(pkgPath string) string {
	parts := strings.Split(pkgPath, "/")
	return parts[len(parts)-1]
}

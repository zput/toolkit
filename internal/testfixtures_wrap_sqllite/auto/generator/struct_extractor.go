package generator

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite/auto/dto"
	"golang.org/x/tools/go/packages"
)

// ExtractStructDefinition 从模型中提取 struct 定义
func ExtractStructDefinition(workDir string, model dto.ModelInfo) (*ast.GenDecl, error) {
	// 加载包
	pkgs, err := loadPackagesForStruct(workDir)
	if err != nil {
		return nil, fmt.Errorf("加载包失败: %w", err)
	}

	// 找到模型所在的包
	var targetPkg *packages.Package
	for _, pkg := range pkgs {
		if pkg.PkgPath == model.PackagePath {
			targetPkg = pkg
			break
		}
	}

	if targetPkg == nil {
		// 调试：列出所有加载的包路径
		var loadedPaths []string
		for _, pkg := range pkgs {
			loadedPaths = append(loadedPaths, pkg.PkgPath)
		}
		return nil, fmt.Errorf("未找到包: %s (已加载的包: %v)", model.PackagePath, loadedPaths)
	}

	// 在包中查找 struct 定义
	var structDecl *ast.GenDecl
	for _, file := range targetPkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			if genDecl, ok := n.(*ast.GenDecl); ok {
				if genDecl.Tok == token.TYPE {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if typeSpec.Name.Name == model.TypeName {
								// 检查是否是 struct 类型
								if _, ok := typeSpec.Type.(*ast.StructType); ok {
									structDecl = genDecl
									return false
								}
							}
						}
					}
				}
			}
			return true
		})
		if structDecl != nil {
			break
		}
	}

	if structDecl == nil {
		return nil, fmt.Errorf("未找到 struct 定义: %s.%s", model.PackagePath, model.TypeName)
	}

	return structDecl, nil
}

// GenerateStructFile 生成 struct 定义文件
func GenerateStructFile(workDir string, models []dto.ModelInfo, outputDir, packageName string) error {
	// 创建 db 目录（用于存放 struct 定义）
	dbDir := filepath.Join(outputDir, "mock_"+packageName, "db")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("创建 db 目录失败: %w", err)
	}

	// 按包分组模型
	modelsByPackage := make(map[string][]dto.ModelInfo)
	for _, model := range models {
		modelsByPackage[model.PackagePath] = append(modelsByPackage[model.PackagePath], model)
	}

	// 为每个包生成一个文件
	for pkgPath, pkgModels := range modelsByPackage {
		// 生成文件名（使用包名）
		fileName := getPackageName(pkgPath) + ".go"
		filePath := filepath.Join(dbDir, fileName)

		// 提取第一个模型的包名（同一包内的模型使用相同的包名）
		pkgName := pkgModels[0].PackageName

		// 生成文件内容
		content, err := generateStructFileContent(workDir, pkgModels, pkgName)
		if err != nil {
			return fmt.Errorf("生成 struct 文件内容失败: %w", err)
		}

		// 写入文件
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return fmt.Errorf("写入 struct 文件失败 %s: %w", filePath, err)
		}
	}

	return nil
}

// generateStructFileContent 生成 struct 文件内容
func generateStructFileContent(workDir string, models []dto.ModelInfo, pkgName string) ([]byte, error) {
	// 加载包以获取类型信息
	pkgs, err := loadPackagesForStruct(workDir)
	if err != nil {
		return nil, fmt.Errorf("加载包失败: %w", err)
	}

	// 建立包路径到包的映射
	pkgMap := make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		pkgMap[pkg.PkgPath] = pkg
	}

	fset := token.NewFileSet()

	// 创建文件 AST
	file := &ast.File{
		Name:  ast.NewIdent("db"),
		Decls: []ast.Decl{},
	}

	// 收集所有需要的导入
	imports := make(map[string]bool)

	// 为每个模型提取 struct 定义
	for _, model := range models {
		structDecl, err := ExtractStructDefinition(workDir, model)
		if err != nil {
			return nil, fmt.Errorf("提取 struct 定义失败 %s: %w", model.FullName, err)
		}

		// 获取模型所在包的类型信息
		modelPkg, ok := pkgMap[model.PackagePath]
		if !ok {
			return nil, fmt.Errorf("未找到包: %s", model.PackagePath)
		}

		// 复制 type spec
		for _, spec := range structDecl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				// 处理 struct 类型
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					// 创建新的 struct 类型，移除所有 tags
					newStructType := removeAllTagsFromStruct(structType)

					// 创建新的 type spec
					newTypeSpec := &ast.TypeSpec{
						Name: typeSpec.Name,
						Type: newStructType,
					}

					// 收集导入（使用类型信息，从字段类型中收集）
					collectImportsFromStructWithTypes(newStructType, modelPkg, imports)

					// 添加到文件
					file.Decls = append(file.Decls, &ast.GenDecl{
						Tok:   token.TYPE,
						Specs: []ast.Spec{newTypeSpec},
					})
				} else {
					// 非 struct 类型，直接复制
					newTypeSpec := &ast.TypeSpec{
						Name: typeSpec.Name,
						Type: typeSpec.Type,
					}
					file.Decls = append(file.Decls, &ast.GenDecl{
						Tok:   token.TYPE,
						Specs: []ast.Spec{newTypeSpec},
					})
				}
			}
		}
	}

	// 添加导入语句
	if len(imports) > 0 {
		importSpecs := []ast.Spec{}
		for imp := range imports {
			importSpecs = append(importSpecs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", imp),
				},
			})
		}
		file.Decls = append([]ast.Decl{&ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: importSpecs,
		}}, file.Decls...)
	}

	// 格式化代码
	var buf strings.Builder
	if err := format.Node(&buf, fset, file); err != nil {
		return nil, fmt.Errorf("格式化代码失败: %w", err)
	}

	return []byte(buf.String()), nil
}

// collectImportsFromStructWithTypes 从 struct 中收集导入（使用类型信息）
func collectImportsFromStructWithTypes(structType *ast.StructType, pkg *packages.Package, imports map[string]bool) {
	if structType.Fields == nil || pkg == nil || pkg.TypesInfo == nil {
		return
	}

	for _, field := range structType.Fields.List {
		if field.Type != nil {
			collectImportsFromTypeWithTypes(field.Type, pkg, imports)
		}
	}
}

// collectImportsFromTypeWithTypes 从类型中收集导入（使用类型信息）
func collectImportsFromTypeWithTypes(expr ast.Expr, pkg *packages.Package, imports map[string]bool) {
	if pkg == nil || pkg.TypesInfo == nil {
		return
	}

	switch t := expr.(type) {
	case *ast.SelectorExpr:
		// 选择器表达式，如 soft_delete.DeletedAt
		if ident, ok := t.X.(*ast.Ident); ok {
			// 方法1: 从 ObjectOf 获取
			if obj := pkg.TypesInfo.ObjectOf(ident); obj != nil {
				if pkgName, ok := obj.(*types.PkgName); ok {
					importedPkg := pkgName.Imported()
					if importedPkg != nil {
						imports[importedPkg.Path()] = true
						return
					}
				}
			}

			// 方法2: 从 Uses 中查找
			for _, useObj := range pkg.TypesInfo.Uses {
				if usePkgName, ok := useObj.(*types.PkgName); ok {
					if usePkgName.Name() == ident.Name {
						importedPkg := usePkgName.Imported()
						if importedPkg != nil {
							imports[importedPkg.Path()] = true
							return
						}
					}
				}
			}

			// 方法3: 从类型信息中获取（通过 TypeOf）
			if exprType := pkg.TypesInfo.TypeOf(expr); exprType != nil {
				if namedType, ok := exprType.(*types.Named); ok {
					if namedType.Obj() != nil && namedType.Obj().Pkg() != nil {
						// 检查是否是外部包的类型
						importedPkg := namedType.Obj().Pkg()
						if importedPkg.Path() != pkg.PkgPath {
							imports[importedPkg.Path()] = true
							return
						}
					}
				}
			}
		}
	case *ast.ArrayType:
		collectImportsFromTypeWithTypes(t.Elt, pkg, imports)
	case *ast.StarExpr:
		collectImportsFromTypeWithTypes(t.X, pkg, imports)
	case *ast.MapType:
		collectImportsFromTypeWithTypes(t.Key, pkg, imports)
		collectImportsFromTypeWithTypes(t.Value, pkg, imports)
	case *ast.ChanType:
		collectImportsFromTypeWithTypes(t.Value, pkg, imports)
	}
}

// loadPackagesForStruct 加载包（用于提取 struct）
func loadPackagesForStruct(workDir string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports,
		Dir:  workDir,
	}

	// 使用 ./... 加载所有包
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, err
	}

	// 检查是否有错误
	var errs []error
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			errs = append(errs, err)
		}
	})

	if len(errs) > 0 {
		// 如果有错误，记录但不立即返回（可能只是警告）
		// 实际使用时可以根据需要处理
	}

	return pkgs, nil
}

// getPackageName 从包路径中提取包名
func getPackageName(pkgPath string) string {
	parts := strings.Split(pkgPath, "/")
	return parts[len(parts)-1]
}

// removeAllTagsFromStruct 从 struct 中移除所有 tags
func removeAllTagsFromStruct(structType *ast.StructType) *ast.StructType {
	if structType.Fields == nil {
		return structType
	}

	// 创建新的 struct 类型
	newStructType := &ast.StructType{
		Fields: &ast.FieldList{
			Opening: structType.Fields.Opening,
			Closing: structType.Fields.Closing,
			List:    make([]*ast.Field, 0, len(structType.Fields.List)),
		},
	}

	// 复制字段，移除所有 tags
	for _, field := range structType.Fields.List {
		newField := &ast.Field{
			Names: field.Names,
			Type:  field.Type,
			Tag:   nil, // 移除所有 tags
		}
		newStructType.Fields.List = append(newStructType.Fields.List, newField)
	}

	return newStructType
}

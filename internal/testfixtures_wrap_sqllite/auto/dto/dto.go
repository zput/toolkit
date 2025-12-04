package dto

import "fmt"

// FunctionInfo 函数信息
type FunctionInfo struct {
	PackagePath string // 包路径
	TypeName    string // 类型名（如果是方法）
	MethodName  string // 方法名（如果是方法）
	FuncName    string // 函数名（如果是包级函数）
	File        string // 文件路径
	Line        int    // 行号
}

// ModelInfo 模型信息
type ModelInfo struct {
	PackagePath string // 包路径
	PackageName string // 包名
	TypeName    string // 类型名
	FullName    string // 完整名称（包名.类型名）
}

func (f FunctionInfo) Key() string {
	if f.TypeName != "" {
		return fmt.Sprintf("%s.%s.%s", f.PackagePath, f.TypeName, f.MethodName)
	}
	return fmt.Sprintf("%s.%s", f.PackagePath, f.FuncName)
}

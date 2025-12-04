package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/zput/toolkit/internal/testfixtures_wrap_sqllite/auto/dto"
)

// GenerateMockCode 生成 Mock Go 代码和 YAML fixture 文件
func GenerateMockCode(workDir, funcPath string, models []dto.ModelInfo, outputDir string) error {
	// 解析函数路径
	filePath, lineNumber, err := parseFunctionPath(funcPath)
	if err != nil {
		return fmt.Errorf("解析函数路径失败: %w", err)
	}

	// 提取包名
	packageName := ExtractPackageName(filePath)
	if packageName == "." {
		// 如果无法提取包名，使用文件名（去掉扩展名）
		baseName := filepath.Base(filePath)
		packageName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
	}

	// 生成类型名
	typeName := GenerateTypeName(packageName, lineNumber)

	// 生成 Mock Go 代码
	mockGoCode, err := generateMockGoCode(filePath, lineNumber, packageName, typeName, models)
	if err != nil {
		return fmt.Errorf("生成 Mock Go 代码失败: %w", err)
	}

	// 生成文件名
	fileName := ConvertFilePathToFileName(filePath)
	mockFileName := fmt.Sprintf("%s_%s_mock.go", fileName, lineNumber)

	// 创建输出目录（outputDir 已经包含了 mock_modules）
	packageOutputDir := filepath.Join(outputDir, "mock_"+packageName)
	if err := os.MkdirAll(packageOutputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 写入 Mock Go 文件
	mockGoPath := filepath.Join(packageOutputDir, mockFileName)
	if err := os.WriteFile(mockGoPath, []byte(mockGoCode), 0644); err != nil {
		return fmt.Errorf("写入 Mock Go 文件失败: %w", err)
	}

	// 生成 YAML fixture 文件
	if err := GenerateYAMLFixtures(outputDir, packageName, models); err != nil {
		return fmt.Errorf("生成 YAML fixture 文件失败: %w", err)
	}

	// 生成 struct 定义文件
	if err := GenerateStructFile(workDir, models, outputDir, packageName); err != nil {
		return fmt.Errorf("生成 struct 文件失败: %w", err)
	}

	return nil
}

// GenerateYAMLFixtures 生成 YAML fixture 文件
func GenerateYAMLFixtures(outputDir, packageName string, models []dto.ModelInfo) error {
	// 创建 fixtures 目录（outputDir 已经包含了 mock_modules）
	fixturesDir := filepath.Join(outputDir, "mock_"+packageName, "testdata", "fixtures")
	if err := os.MkdirAll(fixturesDir, 0755); err != nil {
		return fmt.Errorf("创建 fixtures 目录失败: %w", err)
	}

	// 为每个模型生成 YAML 文件
	for _, model := range models {
		yamlFileName := ModelNameToFileName(model.TypeName) + ".yml"
		yamlPath := filepath.Join(fixturesDir, yamlFileName)

		// 生成 YAML 模板内容
		yamlContent, err := generateYAMLTemplate(model)
		if err != nil {
			return fmt.Errorf("生成 YAML 模板失败: %w", err)
		}

		// 写入文件
		if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
			return fmt.Errorf("写入 YAML 文件失败 %s: %w", yamlPath, err)
		}
	}

	return nil
}

// ModelNameToFileName 将模型名转换为文件名（驼峰转下划线）
func ModelNameToFileName(typeName string) string {
	if len(typeName) == 0 {
		return ""
	}

	var result strings.Builder
	runes := []rune(typeName)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// 如果是第一个字符，直接转小写
		if i == 0 {
			result.WriteRune(unicode.ToLower(r))
			continue
		}

		prevR := runes[i-1]
		isUpper := unicode.IsUpper(r)
		prevIsUpper := unicode.IsUpper(prevR)
		prevIsLower := unicode.IsLower(prevR)
		prevIsDigit := unicode.IsDigit(prevR)
		isDigit := unicode.IsDigit(r)

		if isUpper {
			// 大写字母
			// 如果前一个是小写字母或数字，需要加下划线
			if prevIsLower || prevIsDigit {
				result.WriteByte('_')
			} else if prevIsUpper && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				// 前一个是大写，当前是大写，但下一个是小写（如 "APIKey" -> "api_key"）
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else if isDigit {
			// 数字：数字和字母之间不需要下划线
			result.WriteRune(r)
		} else {
			// 小写字母直接写入
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ExtractPackageName 从文件路径提取包名
func ExtractPackageName(filePath string) string {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "" {
		return "."
	}
	return filepath.Base(dir)
}

// GenerateTypeName 生成类型名
func GenerateTypeName(packageName, lineNumber string) string {
	// 将包名首字母大写，下划线后的字母也大写
	parts := strings.Split(packageName, "_")
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return fmt.Sprintf("%sLine%sMock", result.String(), lineNumber)
}

// ConvertFilePathToFileName 将文件路径转换为文件名（用于生成 Mock 文件名）
func ConvertFilePathToFileName(filePath string) string {
	// 标准化路径
	filePath = filepath.Clean(filePath)

	// 移除 testdata/ 前缀（如果存在）
	filePath = strings.TrimPrefix(filePath, "testdata/")
	filePath = strings.TrimPrefix(filePath, "testdata"+string(filepath.Separator))

	// 如果路径是相对路径，去掉开头的 ./
	if strings.HasPrefix(filePath, "."+string(filepath.Separator)) {
		filePath = strings.TrimPrefix(filePath, "."+string(filepath.Separator))
	}

	// 替换路径分隔符为下划线
	name := strings.ReplaceAll(filePath, string(filepath.Separator), "_")

	// 将 .go 替换为 _go（保留扩展名信息）
	name = strings.ReplaceAll(name, ".go", "_go")

	// 如果还有其他点号，替换为下划线
	name = strings.ReplaceAll(name, ".", "_")

	return name
}

// parseFunctionPath 解析函数路径，返回文件路径和行号
func parseFunctionPath(funcPath string) (string, string, error) {
	parts := strings.SplitN(funcPath, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("无效的函数路径格式，期望 'filepath:line': %s", funcPath)
	}
	return parts[0], parts[1], nil
}

// cleanPackagePath 清理包路径，移除 []* 等前缀
func cleanPackagePath(pkgPath string) string {
	// 处理切片类型 []pkg.Type 或 []*pkg.Type
	if strings.HasPrefix(pkgPath, "[]") {
		pkgPath = strings.TrimPrefix(pkgPath, "[]")
	}
	// 处理指针类型 *pkg.Type
	if strings.HasPrefix(pkgPath, "*") {
		pkgPath = strings.TrimPrefix(pkgPath, "*")
	}
	return pkgPath
}

// mockCodeTemplate Mock Go 代码模板
const mockCodeTemplate = `package mock_{{.PackageName}}

import (
	"fmt"
	"github.com/zput/toolkit/pkg/testfixtures"
	"github.com/zput/toolkit/pkg/utils"
	"path/filepath"
{{range .Imports}}
	"{{.}}"
{{end}}
)

// {{.TypeName}} 自动生成的Mock实现
// 用于测试函数: {{.FilePath}}:{{.LineNumber}}
var _ testfixtures.CallerDbInterface = new({{.TypeName}})

// New{{.TypeName}} 创建新的Mock实例
func New{{.TypeName}}(dbName string) *{{.TypeName}} {
	return &{{.TypeName}}{
		dBName: dbName,
	}
}

type {{.TypeName}} struct {
	dBName string
}

// DefineDbLocationByWillGen 定义数据库文件位置
func (o *{{.TypeName}}) DefineDbLocationByWillGen() string {
	currentPath := utils.Getwd()
	dbPath := filepath.Join(currentPath, "mock_modules/mock_{{.PackageName}}/testdata")
	return filepath.Join(dbPath, fmt.Sprintf("test_%s.db", o.dBName))
}

// DefineDbTableModels 定义需要mock的数据库表模型
func (o *{{.TypeName}}) DefineDbTableModels() (tables []interface{}) {
	return []interface{}{
{{range .Models}}
		new({{.PackageName}}.{{.TypeName}}),
{{end}}
	}
}

// DefineMappingLocationByWillAutoInit2Db 定义YAML fixture数据的位置
func (o *{{.TypeName}}) DefineMappingLocationByWillAutoInit2Db() string {
	return filepath.Join(filepath.Dir(o.DefineDbLocationByWillGen()), "fixtures")
}
`

// yamlTemplate YAML fixture 模板
const yamlTemplate = `# {{.TypeName}} 表数据
# 字段说明：
#   请根据实际模型字段填写数据
#   示例：
#     - id: 1
#       field1: value1
#       field2: value2
#
# 数据记录（请在此处添加实际数据）
# - id: 1
#   # 其他字段...
`

// templateData Mock 代码模板数据
type templateData struct {
	PackageName string
	TypeName    string
	FilePath    string
	LineNumber  string
	Imports     []string
	Models      []modelData
}

// modelData 模型数据
type modelData struct {
	PackageName string
	TypeName    string
}

// generateMockGoCode 生成 Mock Go 代码
func generateMockGoCode(filePath, lineNumber, packageName, typeName string, models []dto.ModelInfo) (string, error) {
	// 准备模板数据
	data := templateData{
		PackageName: packageName,
		TypeName:    typeName,
		FilePath:    filePath,
		LineNumber:  lineNumber,
		Imports:     []string{},
		Models:      []modelData{},
	}

	// 收集模型导入（去重）
	importMap := make(map[string]bool)
	for _, model := range models {
		// 清理包路径，确保不包含 []* 前缀
		cleanPkgPath := cleanPackagePath(model.PackagePath)
		if !importMap[cleanPkgPath] {
			importMap[cleanPkgPath] = true
			data.Imports = append(data.Imports, cleanPkgPath)
		}
		data.Models = append(data.Models, modelData{
			PackageName: model.PackageName,
			TypeName:    model.TypeName,
		})
	}

	// 解析模板
	tmpl, err := template.New("mockCode").Parse(mockCodeTemplate)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	// 执行模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	return buf.String(), nil
}

// yamlTemplateData YAML 模板数据
type yamlTemplateData struct {
	TypeName string
}

// generateYAMLTemplate 生成 YAML 模板内容
func generateYAMLTemplate(model dto.ModelInfo) (string, error) {
	// 准备模板数据
	data := yamlTemplateData{
		TypeName: model.TypeName,
	}

	// 解析模板
	tmpl, err := template.New("yaml").Parse(yamlTemplate)
	if err != nil {
		return "", fmt.Errorf("解析 YAML 模板失败: %w", err)
	}

	// 执行模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行 YAML 模板失败: %w", err)
	}

	return buf.String(), nil
}

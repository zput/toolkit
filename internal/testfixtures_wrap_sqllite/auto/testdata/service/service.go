package service

import "testpkg/model/db/web_game"

type Service struct {
}

// TestCase1: 变量声明中的模型类型
func (s *Service) TestVarDecl() {
	var model *web_game.GameTemplate
	_ = model
}

// TestCase2: 赋值语句中的模型类型（从函数返回类型提取）
func (s *Service) TestAssignFromCall() {
	model, _ := s.getModel()
	_ = model
}

func (s *Service) getModel() (*web_game.GameTemplate, error) {
	return nil, nil
}

// TestCase3: 复合字面量中的模型类型
func (s *Service) TestCompositeLit() {
	model := &web_game.GameTemplate{ID: 1}
	_ = model
}

// TestCase4: 类型断言中的模型类型
func (s *Service) TestTypeAssert(x interface{}) {
	if model, ok := x.(*web_game.GameTemplate); ok {
		_ = model
	}
}

// TestCase5: 函数返回类型中的模型类型
func (s *Service) TestReturnType() *web_game.GameTemplate {
	return nil
}

// TestCase6: 切片类型
func (s *Service) TestSliceType() []*web_game.GameTemplate {
	return nil
}

// TestCase7: 多个模型类型
func (s *Service) TestMultipleModels() {
	var model1 *web_game.GameTemplate
	var model2 *web_game.GameTemplate
	_ = model1
	_ = model2
}

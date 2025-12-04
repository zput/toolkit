package business

import "testpkg/logic/template"

type GameTemplate2Bus struct {
	GameTemplateLogic *template.GameTemplateLogic
}

func (st *GameTemplate2Bus) CreateGamePreCheck() {
	model, _ := st.GameTemplateLogic.GetTmplDetail(1)
	_ = model
}

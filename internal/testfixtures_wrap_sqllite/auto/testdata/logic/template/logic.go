package template

import "testpkg/model/db/web_game"

type GameTemplateLogic struct {
	repo *Repo
}

func (s *GameTemplateLogic) GetTmplDetail(id int) (*web_game.GameTemplate, error) {
	return s.repo.Find(id)
}

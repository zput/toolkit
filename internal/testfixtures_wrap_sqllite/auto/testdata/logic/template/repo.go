package template

import "testpkg/model/db/web_game"

type Repo struct{}

func (r *Repo) Find(id int) (*web_game.GameTemplate, error) {
	return &web_game.GameTemplate{ID: id}, nil
}

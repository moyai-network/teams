package lang

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
)

type Lang struct {
	Name string `cmd:"name"`
}

func (la Lang) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())

	if err != nil {
		return
	}

	for l, t := range translations {
		if t.Properties.Name == la.Name {
			u.Language = data.Language{
				Tag: l,
			}
			data.SaveUser(u)
			return
		}
	}
}
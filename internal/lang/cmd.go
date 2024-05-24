package lang

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
)

type Lang struct {
	Name languages `cmd:"language"`
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
		if t.Properties.Name == string(la.Name) {
			u.Language = data.Language{
				Tag: l,
			}
			data.SaveUser(u)
			return
		}
	}
}

type (
	languages string
)

func (c languages) Type() string {
	return "Language"
}

func (c languages) Options(_ cmd.Source) []string {
	return []string{"English", "French", "Spanish"}
}

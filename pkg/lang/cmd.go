package lang

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"golang.org/x/text/language"
)

type Lang struct {
	Name languages `cmd:"language"`
}

func (la Lang) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, err := data2.LoadUserFromName(p.Name())

	if err != nil {
		return
	}

	for l, t := range translations {
		if t.Properties.Name == string(la.Name) {
			*u.Language = data2.Language{
				Tag: l,
			}
			data2.SaveUser(u)
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
	var langs []string
	for k, _ := range translations {
		langs = append(langs, langString(k))
	}
	return langs
}

func langString(l language.Tag) string {
	switch l {
	case language.English:
		return "English"
	case language.French:
		return "Français"
	case language.Spanish:
		return "Español"
	}
	panic("should never happen: unknown language tag: " + l.String())
}

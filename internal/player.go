package internal

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/roles"
	"github.com/moyai-network/teams/internal/ports/model"
	"github.com/moyai-network/teams/pkg/lang"
	"golang.org/x/text/language"
)

func Broadcastf(tx *world.Tx, key string, a ...interface{}) {
	for p := range Players(tx) {
		Messagef(p, key, a...)
	}
}

func Alertf(tx *world.Tx, s cmd.Source, key string, args ...any) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	for t := range Players(tx) {
		if u, _ := data2.LoadUserFromName(t.Name()); roles.Staff(u.Roles.Highest()) {
			t.Message(lang.Translatef(*u.Language, "staff.alert", p.Name(), fmt.Sprintf(lang.Translate(*u.Language, key), args...)))
		}
	}
}

func Messagef(src cmd.Source, key string, a ...interface{}) {
	out := &cmd.Output{}
	defer src.SendCommandOutput(out)
	l := model.Language{Tag: language.English}

	p, ok := src.(*player.Player)
	if ok {
		u, err := data2.LoadUserFromName(p.Name())
		if err != nil {
			out.Error("An error occurred while loading your user data.")
			return
		}
		l = *u.Language
	}
	out.Print(lang.Translatef(l, key, a...))
}

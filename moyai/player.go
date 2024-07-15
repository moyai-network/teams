package moyai

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/roles"
	"golang.org/x/text/language"
)

func Broadcastf(key string, a ...interface{}) {
	for _, p := range Players() {
		Messagef(p, key, a...)
	}
}

func Alertf(s cmd.Source, key string, args ...any) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	for _, t := range Players() {
		if u, _ := data.LoadUserFromName(t.Name()); roles.Staff(u.Roles.Highest()) {
			t.Message(lang.Translatef(*u.Language, "staff.alert", p.Name(), fmt.Sprintf(lang.Translate(*u.Language, key), args...)))
		}
	}
}

func Messagef(src cmd.Source, key string, a ...interface{}) {
	out := &cmd.Output{}
	defer src.SendCommandOutput(out)
	l := data.Language{Tag: language.English}

	p, ok := src.(*player.Player)
	if ok {
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			out.Error("An error occurred while loading your user data.")
			return
		}
		l = *u.Language
	}
	out.Print(lang.Translatef(l, key, a...))
}

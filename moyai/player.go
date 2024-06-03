package moyai

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
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
		if u, _ := data.LoadUserFromName(t.Name()); role.Staff(u.Roles.Highest()) {
			t.Message(lang.Translatef(*u.Language, "staff.alert", p.Name(), fmt.Sprintf(lang.Translate(*u.Language, key), args...)))
		}
	}
}

func Messagef(p *player.Player, key string, a ...interface{}) {
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		p.Message("An error occurred while loading your user data.")
		return
	}
	p.Message(lang.Translatef(*u.Language, key, a...))
}

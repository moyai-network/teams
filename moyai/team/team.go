package team

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
)

func Broadcast(t data.Team, key string, args ...interface{}) {
	for _, m := range t.Members {
		if h, ok := user.Lookup(m.XUID); ok {
			h.Player().Message(lang.Translatef(h.Player().Locale(), key, args))
		}
	}
}

func FocusingPlayers(t data.Team) (pl []*player.Player) {
	switch t.Focus.Kind {
	case 0:
		if h, ok := user.Lookup(t.Focus.Value); ok {
			pl = append(pl, h.Player())
			return
		}
	case 1:
		tm, err := data.LoadTeam(t.Focus.Value)
		if err != nil {
			t.Focus = data.Focus{}
			_ = data.SaveTeam(t)
			return
		}
		for _, m := range tm.Members {
			if h, ok := user.Lookup(m.XUID); ok {
				pl = append(pl, h.Player())
			}
		}
	default:
		panic("should never happen")
	}
	return
}

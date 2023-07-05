package user

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
)

func BroadcastTeam(t data.Team, key string, args ...interface{}) {
	for _, m := range t.Members {
		if h, ok := Lookup(m.XUID); ok {
			h.Player().Message(lang.Translatef(h.Player().Locale(), key, args...))
		}
	}
}

func FocusingPlayers(t data.Team) (pl []*player.Player) {
	switch t.Focus.Type() {
	case data.FocusTypePlayer():
		if h, ok := Lookup(t.Focus.Value()); ok {
			pl = append(pl, h.p)
			return
		}
	case data.FocusTypeTeam():
		tm, ok := data.LoadTeam(t.Focus.Value())
		if !ok {
			t.Focus = data.Focus{}
			data.SaveTeam(t)
			return
		}
		for _, m := range tm.Members {
			if h, ok := Lookup(m.XUID); ok {
				pl = append(pl, h.p)
			}
		}
	case data.FocusTypeNone():
		return
	default:
		panic("should never happen")
	}
	return
}

func TeamOnlineCount(t data.Team) int {
	var count int
	for _, m := range t.Members {
		if _, ok := Lookup(m.XUID); ok {
			count++
		}
	}
	return count
}

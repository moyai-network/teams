package team

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/data"
	"github.com/moyai-network/teams/moyai/user"
)

func OnlineMembers(tm data.Team) (players []*player.Player) {
	for _, m := range tm.Members {
		if p, ok := user.OnlineFromName(m.Name); ok {
			players = append(players, p)
		}
	}
	return
}

func Broadcastf(tm data.Team, key string, args ...interface{}) {
	for _, p := range OnlineMembers(tm) {
		user.Messagef(p, key, args...)
	}
}

func FocusedOnlinePlayers(t data.Team) (pl []*player.Player) {
	switch t.Focus.Type() {
	case data.FocusTypePlayer():
		if h, ok := user.OnlineFromName(t.Focus.Value()); ok {
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
			if h, ok := Lookup(m.Name); ok {
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

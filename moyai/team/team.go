package team

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
)

func OnlineMembers(tm data.Team) (players []*player.Player) {
	for _, m := range tm.Members {
		if p, ok := user.Lookup(m.Name); ok {
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
		if p, ok := user.Lookup(t.Focus.Value()); ok {
			pl = append(pl, p)
			return
		}
	case data.FocusTypeTeam():
		pl = append(pl, OnlineMembers(t)...)
	case data.FocusTypeNone():
		return
	default:
		panic("should never happen")
	}
	return
}

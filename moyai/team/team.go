package team

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
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
		moyai.Messagef(p, key, args...)
	}
}

package team

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/user"
	"github.com/moyai-network/teams/internal/model"
)

func OnlineMembers(tx *world.Tx, tm model.Team) (players []*player.Player) {
	for _, m := range tm.Members {
		if p, ok := user.Lookup(tx, m.Name); ok {
			players = append(players, p)
		}
	}
	return
}

func Broadcastf(tx *world.Tx, tm model.Team, key string, args ...interface{}) {
	for _, p := range OnlineMembers(tx, tm) {
		internal.Messagef(p, key, args...)
	}
}

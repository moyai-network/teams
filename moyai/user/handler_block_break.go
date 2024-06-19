package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/moyai-network/teams/moyai/data"
)

func (h *Handler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	teams, _ := data.LoadAllTeams()
	if posWithinProtectedArea(h.p, pos, teams) {
		ctx.Cancel()
	}
}

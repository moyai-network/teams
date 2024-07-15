package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/moyai-network/teams/moyai/data"
)

func (h *Handler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil || u.StaffMode || u.Frozen {
		ctx.Cancel()
		return
	}

	teams, _ := data.LoadAllTeams()
	if posWithinProtectedArea(h.p, pos, teams) {
		ctx.Cancel()
	}
}

package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/core"
)

func (h *Handler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	p := ctx.Val()
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok || u.StaffMode || u.Frozen {
		ctx.Cancel()
		return
	}

	teams := core.TeamRepository.FindAll()
	if posWithinProtectedArea(p, pos, teams) {
		ctx.Cancel()
	}
}

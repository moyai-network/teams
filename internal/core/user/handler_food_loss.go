package user

import (
	"github.com/df-mc/dragonfly/server/player"
)

func (h *Handler) HandleFoodLoss(ctx *player.Context, _ int, _ *int) {
	ctx.Cancel()
}

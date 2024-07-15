package user

import "github.com/df-mc/dragonfly/server/event"

func (h *Handler) HandleFoodLoss(ctx *event.Context, _ int, _ *int) {
	ctx.Cancel()
}

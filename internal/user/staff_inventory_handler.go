package user

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
)

type StaffInventoryHandler struct{}

func (StaffInventoryHandler) HandleTake(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}
func (StaffInventoryHandler) HandlePlace(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}
func (StaffInventoryHandler) HandleDrop(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

package user

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
)

type StaffInventoryHandler struct{}

func (StaffInventoryHandler) HandleTake(ctx *inventory.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}
func (StaffInventoryHandler) HandlePlace(ctx *inventory.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}
func (StaffInventoryHandler) HandleDrop(ctx *inventory.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

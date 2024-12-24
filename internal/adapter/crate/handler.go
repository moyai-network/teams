package crate

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
)

type Handler struct {
	inventory.NopHandler
}

// HandleTake ...
func (h Handler) HandleTake(ctx *inventory.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

// HandlePlace ...
func (h Handler) HandlePlace(ctx *inventory.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

// HandleDrop ...
func (h Handler) HandleDrop(ctx *inventory.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

package crate

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
)

type Handler struct {
	inventory.NopHandler
}

// HandleTake ...
func (h Handler) HandleTake(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

// HandlePlace ...
func (h Handler) HandlePlace(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

// HandleDrop ...
func (h Handler) HandleDrop(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

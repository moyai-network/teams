package user

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
)

func (h *Handler) HandleItemPickup(ctx *event.Context, i *item.Stack) {
	u, err := data.LoadUserFromName(h.p.Name())
	if err == nil && u.Teams.PVP.Active() {
		ctx.Cancel()
		return
	}

	for _, sp := range append(it.SpecialItems(), it.PartnerItems()...) {
		if _, ok := i.Value(sp.Key()); ok {
			*i = it.NewSpecialItem(sp, i.Count())
		}
	}
}

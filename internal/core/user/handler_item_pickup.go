package user

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/area"
	it "github.com/moyai-network/teams/internal/core/item"
)

func (h *Handler) HandleItemPickup(ctx *player.Context, i *item.Stack) {
	p := ctx.Val()

	u, ok := core.UserRepository.FindByName(p.Name())
	if ok && (u.Teams.PVP.Active() && !area.Spawn(p.Tx().World()).Vec3WithinOrEqualFloorXZ(p.Position())) {
		ctx.Cancel()
		return
	}

	for _, sp := range append(it.SpecialItems(), it.PartnerItems()...) {
		if _, ok := i.Value(sp.Key()); ok {
			*i = it.NewSpecialItem(sp, i.Count())
		}
	}
}

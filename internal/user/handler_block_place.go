package user

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/data"
	"github.com/moyai-network/teams/internal/eotw"
	it "github.com/moyai-network/teams/internal/item"
)

func (h *Handler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil || u.StaffMode || u.Frozen {
		ctx.Cancel()
		return
	}

	if h.coolDownBonedEffect.Active() {
		internal.Messagef(h.p, "bone.interact")
		ctx.Cancel()
		return
	}

	w := h.p.World()
	teams, _ := data.LoadAllTeams()

	if posWithinProtectedArea(h.p, pos, teams) {
		ctx.Cancel()
	}

	switch bl := b.(type) {
	case block.TNT, it.TripwireHook:
		ctx.Cancel()
	case block.Chest:
		for _, dir := range []cube.Direction{bl.Facing.RotateLeft(), bl.Facing.RotateRight()} {
			sidePos := pos.Side(dir.Face())
			for _, t := range teams {
				if !t.Member(h.p.Name()) {
					c := w.Block(sidePos)
					_, eotw := eotw.Running()
					if _, ok := c.(block.Chest); ok && !eotw && t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(sidePos.Vec3()) {
						ctx.Cancel()
					}
				}
			}
		}
	case block.EnderChest:
		held, _ := h.p.HeldItems()
		if _, ok := held.Value("partner_package"); !ok {
			break
		}
		if typ, ok := it.SpecialItem(held); ok {
			if _, ok := typ.(it.PartnerPackageType); ok {
				ctx.Cancel()
				return
			}
		}
	}
}

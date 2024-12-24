package user

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/eotw"
	"github.com/moyai-network/teams/internal/core/item"
)

func (h *Handler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	p := ctx.Val()
	u, err := data2.LoadUserFromName(p.Name())
	if err != nil || u.StaffMode || u.Frozen {
		ctx.Cancel()
		return
	}

	if h.coolDownBonedEffect.Active() {
		internal.Messagef(p, "bone.interact")
		ctx.Cancel()
		return
	}

	tx := p.Tx()
	teams, _ := data2.LoadAllTeams()

	if posWithinProtectedArea(p, pos, teams) {
		ctx.Cancel()
	}

	switch bl := b.(type) {
	case block.TNT, item.TripwireHook:
		ctx.Cancel()
	case block.Chest:
		for _, dir := range []cube.Direction{bl.Facing.RotateLeft(), bl.Facing.RotateRight()} {
			sidePos := pos.Side(dir.Face())
			for _, t := range teams {
				if !t.Member(p.Name()) {
					c := tx.Block(sidePos)
					_, eotw := eotw.Running()
					if _, ok := c.(block.Chest); ok && !eotw && t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(sidePos.Vec3()) {
						ctx.Cancel()
					}
				}
			}
		}
	case block.EnderChest:
		held, _ := p.HeldItems()
		if _, ok := held.Value("partner_package"); !ok {
			break
		}
		if typ, ok := item.SpecialItem(held); ok {
			if _, ok := typ.(item.PartnerPackageType); ok {
				ctx.Cancel()
				return
			}
		}
	}
}

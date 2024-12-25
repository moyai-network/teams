package user

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/adapter/crate"
	"github.com/moyai-network/teams/internal/core"
	it "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/pkg/lang"
)

func (h *Handler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	p := ctx.Val()
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}

	b := p.Tx().Block(pos)
	if _, ok := b.(block.ItemFrame); ok {
		teams := core.TeamRepository.FindAll()
		if posWithinProtectedArea(p, pos, teams) {
			ctx.Cancel()
			return
		}
	}

	held, _ := p.HeldItems()
	typ, ok := it.PartnerItem(held)
	if ok {
		if cd := h.coolDownGlobalAbilities; cd.Active() {
			p.SendJukeboxPopup(lang.Translatef(*u.Language, "popup.cooldown.partner_item", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		if spi := h.coolDownSpecificAbilities; spi.Active(typ) {
			p.SendJukeboxPopup(lang.Translatef(*u.Language, "popup.cooldown.partner_item.item", typ.Name(), spi.Remaining(typ).Seconds()))
			ctx.Cancel()
		} else {
			p.SendJukeboxPopup(lang.Translatef(*u.Language, "popup.ready.partner_item.item", typ.Name()))
			ctx.Cancel()
		}
	}

	for _, c := range crate.All() {
		if _, ok := b.(block.Chest); ok && pos.Vec3Middle() == c.Position() {
			p.OpenBlockContainer(pos, p.Tx())
			ctx.Cancel()
		}
	}
}

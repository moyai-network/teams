package user

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/crate"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
)

func (h *Handler) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	p := h.p
	u, err := data.LoadUserFromXUID(h.p.XUID())
	if err != nil {
		return
	}

	w := p.World()
	b := w.Block(pos)

	if _, ok := b.(block.ItemFrame); ok {
		teams, _ := data.LoadAllTeams()
		if posWithinProtectedArea(h.p, pos, teams) {
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
			p.OpenBlockContainer(pos)
			ctx.Cancel()
		}
	}
}

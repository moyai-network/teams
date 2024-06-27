package user

import (
	"time"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
)

func (h *Handler) HandleItemConsume(ctx *event.Context, i item.Stack) {
	switch i.Item().(type) {
	case item.GoldenApple:
		cd := h.coolDownGoldenApple
		if cd.Active() {
			moyai.Messagef(h.p, "gapple.cooldown")
			ctx.Cancel()
			return
		}
		cd.Set(time.Second * 30)
	case item.EnchantedApple:
		u, err := data.LoadUserFromName(h.p.Name())
		if err != nil {
			ctx.Cancel()
			return
		}
		if u.Teams.GodApple.Active() {
			moyai.Messagef(h.p, "godapple.cooldown")
			ctx.Cancel()
			return
		}
		moyai.Messagef(h.p, "godapple.active")
		u.Teams.GodApple.Set(time.Hour * 4)
		data.SaveUser(u)
	}

	if _, ok := it.SpecialItem(i); ok {
		ctx.Cancel()
	}

	if _, ok := it.PartnerItem(i); ok {
		ctx.Cancel()
	}
}

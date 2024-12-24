package user

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/core/data"
	it "github.com/moyai-network/teams/internal/core/item"
	"time"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/moyai-network/teams/internal"
)

func (h *Handler) HandleItemConsume(ctx *player.Context, i item.Stack) {
	p := ctx.Val()

	switch i.Item().(type) {
	case item.GoldenApple:
		cd := h.coolDownGoldenApple
		if cd.Active() {
			internal.Messagef(p, "gapple.cooldown")
			ctx.Cancel()
			return
		}
		cd.Set(time.Second * 30)
	case item.EnchantedApple:
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			ctx.Cancel()
			return
		}
		if u.Teams.GodApple.Active() {
			internal.Messagef(p, "godapple.cooldown")
			ctx.Cancel()
			return
		}
		internal.Messagef(p, "godapple.active")
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

package user

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

func (h *Handler) HandleItemDrop(ctx *player.Context, s item.Stack) {
	/*p := ctx.Val()
	w := p.Tx().World()

	u, err := data.LoadUserFromName(p.Name())
	if err != nil || u.StaffMode {
		ctx.Cancel()
		return
	}
	if h.lastArea.Load() != area.Spawn(w) {
		return
	}
	for ent := range p.Tx().Entities() {
		if p, ok := ent.(*player.Player); ok {
		//	p.HideEntity(e)
		}
	}*/
}

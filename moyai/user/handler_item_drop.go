package user

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai/area"
)

func (h *Handler) HandleItemDrop(ctx *event.Context, e world.Entity) {
	w := h.p.World()
	if h.lastArea.Load() != area.Spawn(w) {
		return
	}
	for _, ent := range w.Entities() {
		if p, ok := ent.(*player.Player); ok {
			p.HideEntity(e)
		}
	}
}

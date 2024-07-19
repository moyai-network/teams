package minecraft

import (
    "github.com/df-mc/dragonfly/server/event"
    "github.com/df-mc/dragonfly/server/player"
    "github.com/df-mc/dragonfly/server/world"
    "time"
)

type temporaryHandler struct {
    player.NopHandler
}

func (temporaryHandler) HandleChat(ctx *event.Context, message *string) {
    ctx.Cancel()
}

func (temporaryHandler) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src world.DamageSource) {
    ctx.Cancel()
}

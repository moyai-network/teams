package minecraft

import (
    cube "github.com/df-mc/dragonfly/server/block/cube"
    "github.com/df-mc/dragonfly/server/cmd"
    "github.com/df-mc/dragonfly/server/event"
    "github.com/df-mc/dragonfly/server/item"
    "github.com/df-mc/dragonfly/server/player"
    "github.com/df-mc/dragonfly/server/player/skin"
    "github.com/df-mc/dragonfly/server/world"
    "github.com/go-gl/mathgl/mgl64"
    "time"
)

type temporaryHandler struct {
    player.NopHandler
}

func (temporaryHandler) HandleItemDrop(ctx *event.Context, _ world.Entity) { ctx.Cancel() }
func (temporaryHandler) HandleMove(ctx *event.Context, _ mgl64.Vec3, _ float64, _ float64) {
    ctx.Cancel()
}
func (temporaryHandler) HandleTeleport(ctx *event.Context, _ mgl64.Vec3) { ctx.Cancel() }
func (temporaryHandler) HandleToggleSprint(ctx *event.Context, _ bool)   { ctx.Cancel() }
func (temporaryHandler) HandleToggleSneak(ctx *event.Context, _ bool)    { ctx.Cancel() }
func (temporaryHandler) HandleCommandExecution(ctx *event.Context, _ cmd.Command, _ []string) {
    ctx.Cancel()
}
func (temporaryHandler) HandleChat(ctx *event.Context, _ *string)          { ctx.Cancel() }
func (temporaryHandler) HandleSkinChange(ctx *event.Context, _ *skin.Skin) { ctx.Cancel() }
func (temporaryHandler) HandleStartBreak(ctx *event.Context, _ cube.Pos)   { ctx.Cancel() }
func (temporaryHandler) HandleBlockBreak(ctx *event.Context, _ cube.Pos, _ *[]item.Stack, _ *int) {
    ctx.Cancel()
}
func (temporaryHandler) HandleBlockPlace(ctx *event.Context, _ cube.Pos, _ world.Block) { ctx.Cancel() }
func (temporaryHandler) HandleBlockPick(ctx *event.Context, _ cube.Pos, _ world.Block)  { ctx.Cancel() }
func (temporaryHandler) HandleSignEdit(ctx *event.Context, _ cube.Pos, _ bool, _ string, _ string) {
    ctx.Cancel()
}
func (temporaryHandler) HandleLecternPageTurn(ctx *event.Context, _ cube.Pos, _ int, _ *int) {
    ctx.Cancel()
}
func (temporaryHandler) HandleItemPickup(ctx *event.Context, _ *item.Stack) { ctx.Cancel() }
func (temporaryHandler) HandleItemUse(ctx *event.Context)                   { ctx.Cancel() }
func (temporaryHandler) HandleItemUseOnBlock(ctx *event.Context, _ cube.Pos, _ cube.Face, _ mgl64.Vec3) {
    ctx.Cancel()
}
func (temporaryHandler) HandleItemUseOnEntity(ctx *event.Context, _ world.Entity) { ctx.Cancel() }
func (temporaryHandler) HandleItemConsume(ctx *event.Context, _ item.Stack)       { ctx.Cancel() }
func (temporaryHandler) HandleItemDamage(ctx *event.Context, _ item.Stack, _ int) { ctx.Cancel() }
func (temporaryHandler) HandleAttackEntity(ctx *event.Context, _ world.Entity, _ *float64, _ *float64, _ *bool) {
    ctx.Cancel()
}
func (temporaryHandler) HandleExperienceGain(ctx *event.Context, _ *int) { ctx.Cancel() }
func (temporaryHandler) HandleHurt(ctx *event.Context, _ *float64, _ *time.Duration, _ world.DamageSource) {
    ctx.Cancel()
}
func (temporaryHandler) HandleHeal(ctx *event.Context, _ *float64, _ world.HealingSource) {
    ctx.Cancel()
}
func (temporaryHandler) HandleFoodLoss(ctx *event.Context, _ int, _ *int) { ctx.Cancel() }

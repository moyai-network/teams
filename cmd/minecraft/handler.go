//nolint:unused
package minecraft

import (
	"time"

	cube "github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type temporaryHandler struct {
	player.NopHandler
}

func (temporaryHandler) HandleItemDrop(ctx *player.Context, _ world.Entity) { ctx.Cancel() }
func (temporaryHandler) HandleMove(ctx *player.Context, _ mgl64.Vec3, _ float64, _ float64) {
	ctx.Cancel()
}
func (temporaryHandler) HandleTeleport(ctx *player.Context, _ mgl64.Vec3) { ctx.Cancel() }
func (temporaryHandler) HandleToggleSprint(ctx *player.Context, _ bool)   { ctx.Cancel() }
func (temporaryHandler) HandleToggleSneak(ctx *player.Context, _ bool)    { ctx.Cancel() }
func (temporaryHandler) HandleCommandExecution(ctx *player.Context, _ cmd.Command, _ []string) {
	ctx.Cancel()
}
func (temporaryHandler) HandleChat(ctx *player.Context, _ *string)          { ctx.Cancel() }
func (temporaryHandler) HandleSkinChange(ctx *player.Context, _ *skin.Skin) { ctx.Cancel() }
func (temporaryHandler) HandleStartBreak(ctx *player.Context, _ cube.Pos)   { ctx.Cancel() }
func (temporaryHandler) HandleBlockBreak(ctx *player.Context, _ cube.Pos, _ *[]item.Stack, _ *int) {
	ctx.Cancel()
}
func (temporaryHandler) HandleBlockPlace(ctx *player.Context, _ cube.Pos, _ world.Block) {
	ctx.Cancel()
}
func (temporaryHandler) HandleBlockPick(ctx *player.Context, _ cube.Pos, _ world.Block) { ctx.Cancel() }
func (temporaryHandler) HandleSignEdit(ctx *player.Context, _ cube.Pos, _ bool, _ string, _ string) {
	ctx.Cancel()
}
func (temporaryHandler) HandleLecternPageTurn(ctx *player.Context, _ cube.Pos, _ int, _ *int) {
	ctx.Cancel()
}
func (temporaryHandler) HandleItemPickup(ctx *player.Context, _ *item.Stack) { ctx.Cancel() }
func (temporaryHandler) HandleItemUse(ctx *player.Context)                   { ctx.Cancel() }
func (temporaryHandler) HandleItemUseOnBlock(ctx *player.Context, _ cube.Pos, _ cube.Face, _ mgl64.Vec3) {
	ctx.Cancel()
}
func (temporaryHandler) HandleItemUseOnEntity(ctx *player.Context, _ world.Entity) { ctx.Cancel() }
func (temporaryHandler) HandleItemConsume(ctx *player.Context, _ item.Stack)       { ctx.Cancel() }
func (temporaryHandler) HandleItemDamage(ctx *player.Context, _ item.Stack, _ int) { ctx.Cancel() }
func (temporaryHandler) HandleAttackEntity(ctx *player.Context, _ world.Entity, _ *float64, _ *float64, _ *bool) {
	ctx.Cancel()
}
func (temporaryHandler) HandleExperienceGain(ctx *player.Context, _ *int) { ctx.Cancel() }
func (temporaryHandler) HandleHurt(ctx *player.Context, _ *float64, _ *time.Duration, _ world.DamageSource) {
	ctx.Cancel()
}
func (temporaryHandler) HandleHeal(ctx *player.Context, _ *float64, _ world.HealingSource) {
	ctx.Cancel()
}
func (temporaryHandler) HandleFoodLoss(ctx *player.Context, _ int, _ *int) { ctx.Cancel() }

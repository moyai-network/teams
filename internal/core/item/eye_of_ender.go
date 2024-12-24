package item

import (
	b "github.com/moyai-network/teams/internal/core/block"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

func init() {
	world.RegisterItem(EyeOfEnder{})
	creative.RegisterItem(item.NewStack(EyeOfEnder{}, 1))
}

// EyeOfEnder is an item that can be used to activate an end portal
type EyeOfEnder struct{}

func (f EyeOfEnder) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, u item.User, ctx *item.UseContext) bool {
	if b, ok := tx.Block(pos).(b.PortalFrame); ok {
		b.Filled = true
		tx.SetBlock(pos, b, nil)
		tx.ScheduleBlockUpdate(pos, b, time.Second/4)
		return true
	}
	return false
}

func (EyeOfEnder) EncodeItem() (name string, meta int16) {
	return "minecraft:ender_eye", 0
}

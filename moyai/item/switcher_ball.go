package item

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	ent "github.com/moyai-network/hcf/moyai/entity"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type SwitcherBallType struct{}

func (SwitcherBallType) Name() string {
	return text.Colourf("<aqua>Switcher Ball</aqua>")
}

func (SwitcherBallType) Item() world.Item {
	return SwitcherBall{}
}

func (SwitcherBallType) Lore() []string {
	return []string{text.Colourf("<grey>Throw at a player to switch positions with them.</grey>")}
}

func (SwitcherBallType) Key() string {
	return "switcher_ball"
}

// SwitcherBall is a throwable combat item obtained through shovelling snow.
type SwitcherBall struct{}

// MaxCount ...
func (s SwitcherBall) MaxCount() int {
	return 16
}

// Use ...
func (s SwitcherBall) Use(w *world.World, user item.User, ctx *item.UseContext) bool {
	e := ent.NewSwitcherBall(entity.EyePosition(user), user.Rotation().Vec3().Mul(1.5), user)
	w.AddEntity(e)

	w.PlaySound(user.Position(), sound.ItemThrow{})
	ctx.SubtractFromCount(1)
	return true
}

// EncodeItem ...
func (s SwitcherBall) EncodeItem() (name string, meta int16) {
	return "minecraft:snowball", 0
}

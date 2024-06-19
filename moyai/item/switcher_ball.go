package item

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
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
	e := NewSwitcherBall(entity.EyePosition(user), user.Rotation().Vec3().Mul(1.5), user)
	w.AddEntity(e)

	w.PlaySound(user.Position(), sound.ItemThrow{})
	ctx.SubtractFromCount(1)
	return true
}

// EncodeItem ...
func (s SwitcherBall) EncodeItem() (name string, meta int16) {
	return "minecraft:snowball", 0
}

// NewSwitcherBall creates a switcher ball entity at a position with an owner entity.
func NewSwitcherBall(pos mgl64.Vec3, vel mgl64.Vec3, owner world.Entity) *entity.Ent {
	e := entity.Config{Behaviour: switcherBallConf.New(owner)}.New(SwitcherBallType{}, pos)
	e.SetVelocity(vel)
	return e
}

var switcherBallConf = entity.ProjectileBehaviourConfig{
	Gravity:       0.03,
	Drag:          0.01,
	Particle:      particle.SnowballPoof{},
	ParticleCount: 6,
}

// SwitcherBallEntType is a world.EntityType implementation for snowballs.
type SwitcherBallEntType struct{}

func (SwitcherBallType) EncodeEntity() string { return "minecraft:snowball" }

func (SwitcherBallType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}

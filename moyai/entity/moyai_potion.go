package entity

import (
	"image/color"
	"time"

	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose/nbtconv"
)

const (
	maxDebuffHit    = 1.0993 // Original Value: 1.0393
	maxDebuffMiss   = 0.9593 // Original Value: 0.9093
	maxNoDebuffHit  = 1.0925 // Original Value: 1.0325
	maxNoDebuffMiss = 0.9525 // Original Value: 0.9025
)

// NewMoyaiPotion creates a splash potion. SplashPotion is an item that grants
// effects when thrown.
func NewMoyaiPotion(pos mgl64.Vec3, vel mgl64.Vec3, owner world.Entity, t potion.Potion) world.Entity {
	colour, _ := effect.ResultingColour(t.Effects())

	d := atomic.NewBool(false)
	if t == potion.StrongPoison() {
		d.Toggle()
	}

	s := &SplashPotionType{d: *d, owner: owner}

	conf := moyaiPotionConf
	conf.Potion = t
	conf.Particle = particle.Splash{Colour: colour}
	conf.Hit = s.potionSplash(t, false)

	e := entity.Config{Behaviour: conf.New(owner)}.New(&SplashPotionType{d: *d}, pos)
	e.SetVelocity(vel)
	return e
}

var moyaiPotionConf = entity.ProjectileBehaviourConfig{
	Gravity: 0.06,
	Drag:    0.0025,
	Damage:  -1,
	Sound:   sound.GlassBreak{},
}

// SplashPotionType is a world.EntityType implementation for SplashPotion.
type SplashPotionType struct {
	d     atomic.Bool
	owner world.Entity
}

func (SplashPotionType) EncodeEntity() string { return "minecraft:splash_potion" }
func (SplashPotionType) Glint() bool          { return true }
func (SplashPotionType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}

func (SplashPotionType) DecodeNBT(m map[string]any) world.Entity {
	pot := NewMoyaiPotion(nbtconv.Vec3(m, "Pos"), mgl64.Vec3{}, nil, potion.From(nbtconv.Int32(m, "PotionId"))).(*entity.Ent)
	pot.SetVelocity(nbtconv.Vec3(m, "Motion"))
	return pot
}

func (SplashPotionType) EncodeNBT(e world.Entity) map[string]any {
	pot := e.(*entity.Ent)
	return map[string]any{
		"Pos":      nbtconv.Vec3ToFloat32Slice(pot.Position()),
		"Motion":   nbtconv.Vec3ToFloat32Slice(pot.Velocity()),
		"PotionId": int32(pot.Behaviour().(*entity.ProjectileBehaviour).Potion().Uint8()),
	}
}

// debuff returns true if the potion is a debuff potion.
func (s SplashPotionType) debuff() bool {
	return s.d.Load()
}

// expansion returns the expansion that should be used for the bounding box.
func (s SplashPotionType) expansion() mgl64.Vec3 {
	if s.debuff() {
		return mgl64.Vec3{2.5, 3.5, 2.5}
	}
	return mgl64.Vec3{2, 3, 2}
}

// potionSplash returns a function that creates a potion splash with a specific
// duration multiplier and potion type.
func (s SplashPotionType) potionSplash(pot potion.Potion, linger bool) func(e *entity.Ent, res trace.Result) {
	return func(e *entity.Ent, res trace.Result) {
		w, pos := e.World(), e.Position()

		effects := pot.Effects()
		box := e.Type().BBox(e).Translate(pos)

		_ = color.RGBA{R: 0x38, G: 0x5d, B: 0xc6, A: 0xff}
		if len(effects) > 0 {
			_, _ = effect.ResultingColour(effects)

			debuff := s.debuff()
			expansion := s.expansion()
			ignore := func(e world.Entity) bool {
				_, living := e.(entity.Living)
				return !living
			}

			affected := make(map[entity.Living]float64)
			if entityResult, ok := res.(*trace.EntityResult); ok {
				if splashed, ok := entityResult.Entity().(entity.Living); ok {
					if debuff {
						affected[splashed] = maxDebuffHit
					} else {
						affected[splashed] = maxNoDebuffHit
					}
				}
			}

			for _, e := range w.EntitiesWithin(box.GrowVec3(expansion.Mul(2)), ignore) {
				pos := e.Position()
				if e.Type().BBox(e).Translate(pos).IntersectsWith(box.GrowVec3(expansion)) {
					splashed := e.(entity.Living)
					if debuff {
						affected[splashed] = maxDebuffMiss
					} else {
						affected[splashed] = maxNoDebuffMiss
					}
				}
			}

			for splashed, potency := range affected {
				for _, eff := range effects {
					if p, ok := eff.Type().(effect.PotentType); ok {
						splashed.AddEffect(effect.NewInstant(p.WithPotency(potency), eff.Level()))
						continue
					}

					dur := time.Duration(float64(eff.Duration()) * 0.75 * potency)
					if dur < time.Second {
						continue
					}
					splashed.AddEffect(effect.New(eff.Type().(effect.LastingType), eff.Level(), dur))
				}
			}
		} else if pot == potion.Water() {
			if blockResult, ok := res.(*trace.BlockResult); ok {
				pos := blockResult.BlockPosition().Side(blockResult.Face())
				if _, ok := w.Block(pos).(block.Fire); ok {
					w.SetBlock(pos, nil, nil)
				}

				for _, f := range cube.HorizontalFaces() {
					h := pos.Side(f)
					if _, ok := w.Block(h).(block.Fire); ok {
						w.SetBlock(h, nil, nil)
					}
				}
			}
		}
		if linger {
			w.AddEntity(entity.NewAreaEffectCloud(pos, pot))
		}
	}
}

package check

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// KillAuraB checks if the player hits too many entities in an instance
type KillAuraB struct {
	basic
	attacking map[uint64]any
}

// NewKillAuraB creates a new KillAuraB check.
func NewKillAuraB() *KillAuraB {
	return &KillAuraB{
		attacking: make(map[uint64]any),
	}
}

func (*KillAuraB) Name() (string, string) {
	return "KillAura", "B"
}

func (*KillAuraB) Description() string {
	return "This checks if the player hits too many entities in an instance"
}

// MaxViolations ...
func (*KillAuraB) MaxViolations() float64 {
	return 15
}

// Process ...
func (k *KillAuraB) Process(p Processor, pk packet.Packet) bool {
	switch pk := pk.(type) {
	case *packet.InventoryTransaction:
		if data, ok := pk.TransactionData.(*protocol.UseItemOnEntityTransactionData); ok && data.ActionType == protocol.UseItemOnEntityActionAttack {
			if _, ok := k.attacking[data.TargetEntityRuntimeID]; !ok {
				k.attacking[data.TargetEntityRuntimeID] = true
			}
		}

	case *packet.PlayerAuthInput:
		if len(k.attacking) > 1 {
			var lastBBox cube.BBox
			var collides bool
			for id := range k.attacking {
				ent, ok := p.SearchEntity(id)
				if !ok {
					continue
				}
				loc := ent.LastPosition()
				bbox := cube.Box(float64(loc.X()-0.8), float64(loc.Y()), float64(loc.Z()-0.8), float64(loc.X()+0.8), float64(loc.Y()+1.8), float64(loc.Z()+0.8)).Extend(mgl64.Vec3{0.3, 0.3, 0.3})
				if lastBBox.Length() != 0 {
					collides = bbox.IntersectsWith(lastBBox)
					if collides {
						break
					}
				}
				lastBBox = bbox
			}
			if !collides {
				p.Flag(k, k.violationAfterTicks(p.ClientFrame(), 300), map[string]any{
					"entities": len(k.attacking),
				})
			}
		}
		p.Debug(k, map[string]any{
			"entities": len(k.attacking),
		})
		k.attacking = make(map[uint64]any)
	}

	return false
}

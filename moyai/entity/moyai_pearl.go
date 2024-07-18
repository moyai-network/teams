package entity

import (
	"fmt"
	"math"
	"strings"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"

	//"github.com/paroxity/portal/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

var directions = map[world.Entity]cube.Direction{}

var yaws = map[world.Entity]float64{}
var pitches = map[world.Entity]float64{}

// NewMoyaiPearl creates an EnderPearl entity. EnderPearl is a smooth, greenish-
// blue item used to teleport.
func NewMoyaiPearl(pos mgl64.Vec3, vel mgl64.Vec3, owner world.Entity) world.Entity {
	e := entity.Config{Behaviour: moyaiPearlConf.New(owner)}.New(entity.EnderPearlType{}, pos)
	e.SetVelocity(vel.Mul(1.35))

	directions[owner] = owner.Rotation().Direction()
	yaws[owner] = math.Round(owner.Rotation().Yaw())
	pitches[owner] = math.Round(owner.Rotation().Pitch())
	return e
}

var moyaiPearlConf = entity.ProjectileBehaviourConfig{
	Gravity:               0.03,
	Drag:                  0.01,
	KnockBackHeightAddend: 0.4 - 0.45,
	KnockBackForceAddend:  0.42 - 0.3608,
	Particle:              particle.EndermanTeleport{},
	Sound:                 sound.Teleport{},
	Hit:                   teleport,
}

// teleporter represents a living entity that can teleport.
type teleporter interface {
	// Teleport teleports the entity to the position given.
	Teleport(pos mgl64.Vec3)
	entity.Living
}

// teleport teleports the owner of an Ent to a trace.Result's position.
func teleport(e *entity.Ent, target trace.Result) {
	// if tlp, ok := e.Behaviour().(*entity.ProjectileBehaviour).Owner().(teleporter); ok {
	// 	if p, ok := tlp.(*player.Player); ok {
	// 		h, ok := p.Handler().(*user.Handler)
	// 		if !ok {
	// 			return
	// 		}
	// 		if h.Combat().Active() && area.Spawn(tlp.World()).Vec3WithinOrEqualXZ(target.Position()) {
	// 			h.Pearl().Reset()
	// 			return
	// 		}

	// 		u, _ := data.LoadUserFromName(p.Name())
	// 		teams, _ := data.LoadAllTeams()
	// 		if u.Teams.PVP.Active() {
	// 			for _, t := range teams {
	// 				a := t.Claim
	// 				if a.Vec3WithinOrEqualXZ(target.Position()) {
	// 					h.Pearl().Reset()
	// 					return
	// 				}
	// 			}
	// 		}
	// 	}

	// 	b := e.World().Block(cube.PosFromVec3(target.Position()))
	// 	p := e.Behaviour().(*entity.ProjectileBehaviour).Owner().(teleporter).(*player.Player)
	// 	rot := p.Rotation()
	// 	if f, ok := b.(block.WoodFenceGate); ok || f.Open {
	// 		session_writePacket(player_session(p), &packet.MovePlayer{
	// 			EntityRuntimeID: 1,
	// 			Position:        mgl32.Vec3{float32(e.Position()[0]), float32(e.Position()[1] + 1.621), float32(e.Position()[2])},
	// 			Pitch:           float32(rot[1]),
	// 			Yaw:             float32(rot[0]),
	// 			HeadYaw:         float32(rot[0]),
	// 			Mode:            packet.MoveModeNormal,
	// 		})
	// 	}

	// 	if slabPos, ok := validSlabPosition(e, target, directions[p]); ok {
	// 		session_writePacket(player_session(p), &packet.MovePlayer{
	// 			EntityRuntimeID: 1,
	// 			Position:        mgl32.Vec3{float32(slabPos[0]), float32(slabPos[1] + 1.621), float32(slabPos[2])},
	// 			Pitch:           float32(rot[1]),
	// 			Yaw:             float32(rot[0]),
	// 			HeadYaw:         float32(rot[0]),
	// 			Mode:            packet.MoveModeNormal,
	// 		})
	// 	}

	// 	if taliPos, ok := validTaliPearl(e, target, directions[p]); ok {
	// 		session_writePacket(player_session(p), &packet.MovePlayer{
	// 			EntityRuntimeID: 1,
	// 			Position:        mgl32.Vec3{float32(taliPos[0]), float32(taliPos[1] + 1.621), float32(taliPos[2])},
	// 			Pitch:           float32(rot[1]),
	// 			Yaw:             float32(rot[0]),
	// 			HeadYaw:         float32(rot[0]),
	// 			Mode:            packet.MoveModeNormal,
	// 		})
	// 	}

	// 	onGround := p.OnGround()
	// 	for _, v := range p.World().Viewers(p.Position()) {
	// 		v.ViewEntityMovement(p, e.Position(), rot, onGround)
	// 	}

	// 	e.World().PlaySound(tlp.Position(), sound.Teleport{})

	// 	data, err := data.LoadUserFromName(p.Name())
	// 	if err == nil && data.Teams.Settings.Visual.PearlAnimation {
	// 		session_writePacket(player_session(p), &packet.MovePlayer{
	// 			EntityRuntimeID: 1,
	// 			Position:        mgl32.Vec3{float32(target.Position()[0]), float32(target.Position()[1] + 1.621), float32(target.Position()[2])},
	// 			Pitch:           float32(rot[1]),
	// 			Yaw:             float32(rot[0]),
	// 			HeadYaw:         float32(rot[0]),
	// 			Mode:            packet.MoveModeNormal,
	// 		})
	// 	} else {
	// 		tlp.Teleport(target.Position())
	// 	}

	// 	//tlp.Hurt(5, entity.FallDamageSource{})
	// }
	p, ok := e.Behaviour().(*entity.ProjectileBehaviour).Owner().(*player.Player)
	usr, ok2 := p.Handler().(*user.Handler)

	if !ok || !ok2 {
		e.World().RemoveEntity(e)
		_ = e.Close()
		return
	}

	if usr.Combat().Active() && area.Spawn(p.World()).Vec3WithinOrEqualXZ(target.Position()) {
		usr.Pearl().Reset()
		return
	}

	u, _ := data.LoadUserFromName(p.Name())
	if u.Teams.PVP.Active() {
		teams, _ := data.LoadAllTeams()
		for _, t := range teams {
			a := t.Claim
			if a.Vec3WithinOrEqualXZ(target.Position()) {
				usr.Pearl().Reset()
				return
			}
		}
	}

	if pos, ok := validPosition(e, target, directions[p]); ok {
		p.Teleport(pos)
		p.PlaySound(sound.Teleport{})
		p.Hurt(5, entity.FallDamageSource{})
	} else {
		usr.Pearl().Reset()
		p.SendPopup(text.Colourf("<red>Pearl Refunded</red>"))
		if !p.GameMode().CreativeInventory() {
			_, _ = p.Inventory().AddItem(item.NewStack(item.EnderPearl{}, 1))
		}
	}
}

func validSlabPosition(e *entity.Ent, target trace.Result, direction cube.Direction) (mgl64.Vec3, bool) {
	pos := cube.Pos{int(target.Position().X()), int(target.Position().Y()), int(target.Position().Z())}

	newPos := cube.Pos{}

	switch direction.String() {
	case "west":
		newPos = pos.Sub(cube.Pos{1, 0, 0})
	case "east":
		newPos = pos.Add(cube.Pos{1, 0, 0})
	case "north":
		newPos = pos.Sub(cube.Pos{0, 0, 1})
	case "south":
		newPos = pos.Add(cube.Pos{0, 0, 1})
	}

	if _, ok := e.World().Block(newPos).(block.Air); ok {
		if _, ok := e.World().Block(newPos.Add(cube.Pos{0, 1, 0})).(block.Air); ok {
			return newPos.Vec3(), true
		}
	}

	return mgl64.Vec3{}, false
}

func validAirPosition(e *entity.Ent, target trace.Result, direction cube.Direction) (mgl64.Vec3, bool) {
	pos := cube.Pos{int(target.Position().X()), int(target.Position().Y()), int(target.Position().Z())}
	if _, ok := e.World().Block(pos.Add(cube.Pos{0, 1, 0})).(block.Air); !ok {
		newPos := cube.Pos{}

		switch direction.String() {
		case "west":
			newPos = pos.Sub(cube.Pos{1, 0, 0})
		case "east":
			newPos = pos.Add(cube.Pos{1, 0, 0})
		case "north":
			newPos = pos.Sub(cube.Pos{0, 0, 1})
		case "south":
			newPos = pos.Add(cube.Pos{0, 0, 1})
		}

		if _, ok := e.World().Block(newPos).(block.Air); ok {
			if _, ok := e.World().Block(newPos.Add(cube.Pos{0, 1, 0})).(block.Air); ok {
				return newPos.Vec3(), true
			}
		}
		return mgl64.Vec3{}, false
	}

	return mgl64.Vec3{}, false
}

func validFencePosition(e *entity.Ent, target trace.Result, direction cube.Direction) (mgl64.Vec3, bool) {
	pos := target.Position()
	b, ok := e.World().Block(cube.PosFromVec3(pos)).(block.WoodFenceGate)
	if ok && b.Open {
		newPos := mgl64.Vec3{}

		switch direction {
		case cube.West:
			newPos = pos.Add(mgl64.Vec3{-1, 0, 0})
		case cube.East:
			newPos = pos.Add(mgl64.Vec3{1, 0, 0})
		case cube.North:
			newPos = pos.Add(mgl64.Vec3{0, 0, -1})
		case cube.South:
			newPos = pos.Add(mgl64.Vec3{0, 0, 1})
		}

		if _, ok := e.World().Block(cube.PosFromVec3(newPos)).(block.Air); ok {
			if _, ok := e.World().Block(cube.PosFromVec3(newPos.Add(mgl64.Vec3{0, 1, 0}))).(block.Air); ok {
				return newPos, true
			}
		}
	}
	return pos, true
}

func validPosition(e *entity.Ent, target trace.Result, direction cube.Direction) (mgl64.Vec3, bool) {
	pos := cube.Pos{int(target.Position().X()), int(target.Position().Y()), int(target.Position().Z())}

	b := e.World().Block(pos)
	name, properties := b.EncodeBlock()

	if fencePos, fenceOk := validFencePosition(e, target, direction); fenceOk {
		fmt.Println("Fence Pearl")
		return fencePos, true
	}

	if underPos, underOk := validUnderPearl(e, target, direction); underOk {
		// // fmt.Println("Under Pearl")
		return underPos, true
	}

	if strings.Contains(name, "slab") {
		if strings.Contains(name, "double") {
			return mgl64.Vec3{}, false
		}

		if slabPos, slabOk := validSlabPosition(e, target, direction); slabOk {
			// // fmt.Println("Slab Pearl")
			return slabPos, true
		}
	}

	if name == "minecraft:air" || strings.Contains(name, "stairs") {
		if strings.Contains(name, "stairs") {
			if len(properties) != 0 {
				if properties["upside_down_bit"] == false && properties["weirdo_direction"] == int32(1) {
					return mgl64.Vec3{}, false
				}
			}
		}

		if airPos, airOk := validAirPosition(e, target, direction); airOk {
			// // fmt.Println("Air Pearl")
			return airPos, true
		}
	}

	if taliPos, ok := validTaliPearl(e, target, direction); ok {
		return taliPos, true
	}

	if validBlock(b) && validBlock(e.World().Block(pos.Add(cube.Pos{0, 1, 0}))) {
		// fmt.Println("Valid Block Pearl")
		return pos.Vec3(), true
	}

	return pos.Vec3(), false
}

func validUnderPearl(e *entity.Ent, target trace.Result, direction cube.Direction) (mgl64.Vec3, bool) {
	pos := cube.Pos{int(target.Position().X()), int(target.Position().Y()), int(target.Position().Z())}

	pitch := pitches[e]

	if pitch > -45 && pitch < 45 {
		if _, ok := e.World().Block(pos.Add(cube.Pos{0, 1, 0})).(block.Air); ok {
			return pos.Vec3(), true
		}
	}

	return mgl64.Vec3{}, false
}

func validTaliPearl(e *entity.Ent, target trace.Result, direction cube.Direction) (mgl64.Vec3, bool) {
	p, _ := e.Behaviour().(*entity.ProjectileBehaviour).Owner().(*player.Player)
	pos := cube.Pos{int(target.Position().X()), int(target.Position().Y()), int(target.Position().Z())}

	newPos := cube.Pos{}
	yaw := yaws[p]

	switch direction.String() {
	case "west":
		if yaw >= 120 && yaw <= 140 {
			newPos = pos.Add(cube.Pos{-1, 0, -1})
		}

		if yaw <= 60 && yaw >= 40 {
			newPos = pos.Add(cube.Pos{-1, 0, 1})
		}

		// fmt.Println("west")
	case "east":
		if yaw >= -70 && yaw <= -50 {
			newPos = pos.Add(cube.Pos{1, 0, 1})
			// fmt.Println("option 1")
		}

		if yaw >= -130 && yaw <= -110 {
			newPos = pos.Add(cube.Pos{1, 0, -1})
			// fmt.Println("option 2")
		}

		// fmt.Println("east")
	case "north":
		if yaw <= -160 && yaw >= -180 {
			newPos = pos.Add(cube.Pos{-1, 0, 1})
		}

		if yaw <= 160 && yaw >= 140 {
			newPos = pos.Add(cube.Pos{-1, 0, -1})
		}
		// fmt.Println("north")
	case "south":
		if yaw >= 20 && yaw <= 40 {
			newPos = pos.Add(cube.Pos{-1, 0, 1})
			// fmt.Println("option 1")
		}

		if yaw <= -20 && yaw >= -40 {
			newPos = pos.Add(cube.Pos{1, 0, 1})
			// fmt.Println("option 2")
		}

		// fmt.Println("south")
	}

	if _, ok := e.World().Block(newPos).(block.Air); ok {
		if _, ok := e.World().Block(newPos.Add(cube.Pos{0, 1, 0})).(block.Air); ok {
			return newPos.Vec3(), true
		}
	}

	return mgl64.Vec3{}, false
}

func validBlock(block2 world.Block) bool {
	//TODO: this needs a lot of work :pain:

	var blocks = []world.Block{block.DoubleFlower{}, block.DoubleTallGrass{}, block.ShortGrass{}, block.Flower{}, block.DoubleFlower{}, block.Air{}, block.Grass{}, block.Sand{}, block.Leaves{}, block.GlassPane{}, block.StainedGlassPane{}}

	name, _ := block2.EncodeBlock()
	for _, bk := range blocks {
		targetName, _ := bk.EncodeBlock()

		if strings.HasSuffix(name, "glass_pane") {
			return true
		}

		if targetName == name {
			return true
		}
	}

	return true
}

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session

// noinspection ALL
//
//go:linkname session_writePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
func session_writePacket(*session.Session, packet.Packet)

// noinspection ALL
//
//go:linkname session_ViewEntityTeleport github.com/df-mc/dragonfly/server/session.(*Session).ViewEntityTeleport
func session_ViewEntityTeleport(*session.Session, world.Entity, mgl64.Vec3)

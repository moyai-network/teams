package moyai

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/moyai/area"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math/rand"
	"time"
)

var lastDropPos cube.Pos

func tickAirDrop(w *world.World) {
	for {
		<-time.After(time.Minute * 10)
		pos := findAirDropPosition(w)
		Broadcastf("airdrop.incoming", pos.X(), pos.Z())
		for _, p := range Players() {
			p.PlaySound(sound.BarrelOpen{})
			p.PlaySound(sound.FireworkHugeBlast{})
			p.PlaySound(sound.FireworkLaunch{})
			p.PlaySound(sound.Note{})
		}
		destroyAirDrop(w, lastDropPos)
		dropAirDrop(w, pos)
	}
}

func dropAirDrop(w *world.World, pos cube.Pos) {
	bl := generateAirDrop(w)
	w.SetBlock(pos, bl, nil)
	for _, p := range Players() {
		p.PlaySound(sound.BarrelClose{})
		p.PlaySound(sound.FireworkBlast{})
		p.PlaySound(sound.FireworkTwinkle{})
		p.PlaySound(sound.Note{})
	}

	for {
		<-time.After(time.Second)

		oldPos := pos
		pos = pos.Add(cube.Pos{0, -1, 0})
		if _, ok := w.Block(pos).(block.Air); !ok {
			fillAirDrop(bl.Inventory())

			ch, ok := w.Block(oldPos).(block.Chest)
			if !ok {
				break
			}
			h, ok := ch.Inventory().Handler().(*airDropInventoryHandler)
			if !ok {
				break
			}
			h.pos = oldPos

			for _, p := range w.Viewers(pos.Vec3()) {
				p.ViewSound(pos.Vec3Centre(), sound.Fall{})
			}

			lastDropPos = oldPos
			return
		}

		w.SetBlock(oldPos, block.Air{}, nil)
		w.SetBlock(pos, bl, nil)
	}
}

func generateAirDrop(w *world.World) block.Chest {
	bl := block.NewChest()
	bl = bl.WithName(text.Colourf("<red>Air Drop</red>")).(block.Chest)
	bl.Inventory().Handle(&airDropInventoryHandler{
		inv: bl.Inventory(),
		w:   w,
	})
	return bl
}

func fillAirDrop(inventory *inventory.Inventory) {
	for i := 0; i < 27; i++ {
		if rand.Intn(2) == 0 {
			_ = inventory.SetItem(i, randomItem())
		}
	}
}

func randomItem() item.Stack {
	items := it.SpecialItems()
	return item.NewStack(items[rand.Intn(len(items))].Item(), rand.Intn(5))
}

func findAirDropPosition(w *world.World) cube.Pos {
	spawn := area.Spawn(w)
	warzone := area.WarZone(w)

	for {
		x := warzone.Min().X() + (warzone.Max().X()-warzone.Min().X())*rand.Float64()
		z := warzone.Min().Y() + (warzone.Max().Y()-warzone.Min().Y())*rand.Float64()
		pos := cube.Pos{int(x), 255, int(z)}
		if !spawn.Vec3WithinOrEqualFloorXZ(pos.Vec3Centre()) {
			return pos
		}
	}
}

type airDropInventoryHandler struct {
	inventory.NopHandler
	inv *inventory.Inventory
	pos cube.Pos
	w   *world.World
}

func (h airDropInventoryHandler) HandleTake(ctx *event.Context, _ int, st item.Stack) {
	stacks := h.inv.Items()
	if len(stacks) == 1 && stacks[0].Equal(st) {
		time.AfterFunc(time.Second, func() {
			destroyAirDrop(h.w, h.pos)
		})
		return
	}
}
func (airDropInventoryHandler) HandlePlace(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

func (h airDropInventoryHandler) HandleDrop(ctx *event.Context, _ int, st item.Stack) {
	stacks := h.inv.Items()
	if len(stacks) == 1 && stacks[0].Equal(st) {
		time.AfterFunc(time.Second, func() {
			destroyAirDrop(h.w, h.pos)
		})
		return
	}
}

func destroyAirDrop(w *world.World, pos cube.Pos) {
	if _, ok := w.Block(pos).(block.Air); ok {
		return
	}
	for _, v := range w.Viewers(pos.Vec3()) {
		v.ViewSound(pos.Vec3Centre(), sound.Explosion{})
	}
	w.SetBlock(pos, block.Air{}, nil)
}

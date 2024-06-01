package moyai

import (
	"fmt"
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

var (
	lastDropPos          cube.Pos
	parachuteWoolOffsets = []cube.Pos{
		{2, 6, 0},
		{3, 6, 0},
		{3, 6, 1},
		{3, 6, -1},
		{2, 6, -2},
		{1, 6, -3},
		{0, 6, -3},
		{-1, 6, -3},
		{-2, 6, -3},
		{-2, 6, -2},
		{-3, 6, -2},
		{-3, 6, -1},
		{-3, 6, 0},
		{-3, 6, 1},
		{-3, 6, 2},
		{-2, 6, 2},
		{-2, 6, 3},
		{-1, 6, 3},
		{0, 6, 3},
		{1, 6, 3},
		{2, 6, 2},

		{2, 7, 1},
		{2, 7, 0},
		{2, 7, -1},
		{1, 7, -2},
		{0, 7, -2},
		{-1, 7, -2},
		{-2, 7, -1},
		{-2, 7, 0},
		{-2, 7, 1},
		{-1, 7, 2},
		{0, 7, 2},
		{1, 7, 2},

		{0, 8, 0},
		{0, 8, -1},
		{0, 8, 1},
		{1, 8, -1},
		{1, 8, 0},
		{1, 8, 1},
		{-1, 8, -1},
		{-1, 8, 0},
		{-1, 8, 1},
	}

	parachuteRootOffsets = []cube.Pos{
		{0, 1, 0},
		{0, 2, 0},
		{1, 2, 0},
		{0, 2, -1},
		{0, 2, 1},

		{-1, 2, -1},
		{-1, 3, -1},
		{-2, 3, -1},
		{-2, 4, -1},
		{-2, 4, -2},
		{-2, 5, -2},

		{-1, 2, 1},
		{-1, 3, 1},
		{-1, 4, 1},
		{-1, 4, 2},
		{-2, 4, 2},
		{-2, 5, 2},

		{1, 2, 0},
		{1, 3, 0},
		{1, 4, 0},
		{2, 4, 0},
		{2, 5, 0},
	}
)

func tickAirDrop(w *world.World) {
	for {
		<-time.After(time.Second * 20)
		pos := findAirDropPosition(w)
		destroyAirDrop(w, lastDropPos)
		dropAirDrop(w, pos)
		return
	}
}

func dropAirDrop(w *world.World, pos cube.Pos) {
	Broadcastf("airdrop.incoming", pos.X(), pos.Z())

	bl := generateAirDrop(w)
	w.SetBlock(pos, bl, nil)
	for _, p := range Players() {
		p.Teleport(pos.Vec3())

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

		removeParachute(w, oldPos)
		w.SetBlock(oldPos, block.Air{}, nil)

		placeParachute(w, pos)
		w.SetBlock(pos, bl, nil)
	}
}

func removeParachute(w *world.World, pos cube.Pos) {
	for _, root := range parachuteRoots(pos) {
		if _, ok := w.Block(root).(block.WoodFence); !ok {
			continue
		}
		w.SetBlock(root, block.Air{}, nil)
	}
	for _, wool := range parachuteWool(pos) {
		if _, ok := w.Block(wool).(block.Wool); !ok {
			continue
		}
		w.SetBlock(wool, block.Air{}, nil)
	}
}

func placeParachute(w *world.World, pos cube.Pos) {
	for _, root := range parachuteRoots(pos) {
		if _, ok := w.Block(root).(block.Air); !ok {
			continue
		}
		w.SetBlock(root, block.WoodFence{Wood: block.OakWood()}, nil)
	}
	for _, wool := range parachuteWool(pos) {
		if _, ok := w.Block(wool).(block.Air); !ok {
			continue
		}
		w.SetBlock(wool, block.Wool{Colour: item.ColourRed()}, nil)
	}
}

func parachuteRoots(pos cube.Pos) []cube.Pos {
	roots := make([]cube.Pos, len(parachuteRootOffsets))
	for i, off := range parachuteRootOffsets {
		roots[i] = pos.Add(off)
	}
	return roots
}

func parachuteWool(pos cube.Pos) []cube.Pos {
	wool := make([]cube.Pos, len(parachuteWoolOffsets))
	for i, off := range parachuteWoolOffsets {
		if pos.Y() == 3 {
			fmt.Println(pos)
		}
		wool[i] = pos.Add(off)
	}
	return wool
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
	return it.NewSpecialItem(items[rand.Intn(len(items))], rand.Intn(5))
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

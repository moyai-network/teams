package moyai

import (
	_ "embed"
	"math/rand"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/moyai/area"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	err := gophig.GetConfComplex("assets/air_drop/parachute_offset_fence_1.json", gophig.JSONMarshaler{}, &parachuteFenceOffsets)
	checkErr(err)
	err = gophig.GetConfComplex("assets/air_drop/parachute_offset_wool_1.json", gophig.JSONMarshaler{}, &parachuteWoolOffsets)
	checkErr(err)
	err = gophig.GetConfComplex("assets/air_drop/parachute_offset_fence_2.json", gophig.JSONMarshaler{}, &parachuteFenceOffsets2)
	checkErr(err)
	err = gophig.GetConfComplex("assets/air_drop/parachute_offset_wool_2.json", gophig.JSONMarshaler{}, &parachuteWoolOffsets2)
	checkErr(err)
	err = gophig.GetConfComplex("assets/air_drop/parachute_offset_fence_3.json", gophig.JSONMarshaler{}, &parachuteFenceOffsets3)
	checkErr(err)
	err = gophig.GetConfComplex("assets/air_drop/parachute_offset_wool_3.json", gophig.JSONMarshaler{}, &parachuteWoolOffsets3)
	checkErr(err)
	err = gophig.GetConfComplex("assets/air_drop/parachute_offset_fence_4.json", gophig.JSONMarshaler{}, &parachuteFenceOffsets4)
	checkErr(err)
	err = gophig.GetConfComplex("assets/air_drop/parachute_offset_wool_4.json", gophig.JSONMarshaler{}, &parachuteWoolOffsets4)
	checkErr(err)
	err = gophig.GetConfComplex("assets/air_drop/parachute_offset_fence_5.json", gophig.JSONMarshaler{}, &parachuteWoolOffsets5)
	checkErr(err)
	err = gophig.GetConfComplex("assets/air_drop/parachute_offset_wool_5.json", gophig.JSONMarshaler{}, &parachuteWoolOffsets5)
	checkErr(err)
}

var (
	lastDropPos cube.Pos

	parachuteFenceOffsets  []cube.Pos
	parachuteWoolOffsets   []cube.Pos
	parachuteFenceOffsets2 []cube.Pos
	parachuteFenceOffsets3 []cube.Pos
	parachuteFenceOffsets4 []cube.Pos
	parachuteFenceOffsets5 []cube.Pos
	parachuteWoolOffsets2  []cube.Pos
	parachuteWoolOffsets3  []cube.Pos
	parachuteWoolOffsets4  []cube.Pos
	parachuteWoolOffsets5  []cube.Pos
)

func tickAirDrop(w *world.World) {
	for {
		<-time.After(time.Hour)
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
			removeParachute(w, oldPos)
			placeParachuteBlock(parachuteFenceOffsets2, w, pos, block.WoodFence{Wood: block.OakWood()})
			placeParachuteBlock(parachuteWoolOffsets2, w, pos, block.Wool{Colour: item.ColourRed()})
			<-time.After(time.Second / 4)
			removeParachuteBlock(parachuteFenceOffsets2, w, pos)
			removeParachuteBlock(parachuteWoolOffsets2, w, pos)
			placeParachuteBlock(parachuteFenceOffsets3, w, pos, block.WoodFence{Wood: block.OakWood()})
			placeParachuteBlock(parachuteWoolOffsets3, w, pos, block.Wool{Colour: item.ColourRed()})
			<-time.After(time.Second / 4)
			removeParachuteBlock(parachuteFenceOffsets3, w, pos)
			removeParachuteBlock(parachuteWoolOffsets3, w, pos)
			placeParachuteBlock(parachuteFenceOffsets4, w, pos, block.WoodFence{Wood: block.OakWood()})
			placeParachuteBlock(parachuteWoolOffsets4, w, pos, block.Wool{Colour: item.ColourRed()})
			<-time.After(time.Second / 4)
			removeParachuteBlock(parachuteFenceOffsets4, w, pos)
			removeParachuteBlock(parachuteWoolOffsets4, w, pos)
			placeParachuteBlock(parachuteFenceOffsets5, w, pos, block.WoodFence{Wood: block.OakWood()})
			placeParachuteBlock(parachuteWoolOffsets5, w, pos, block.Wool{Colour: item.ColourRed()})
			<-time.After(time.Second / 4)
			removeParachuteBlock(parachuteFenceOffsets5, w, pos)
			removeParachuteBlock(parachuteWoolOffsets5, w, pos)
			fillAirDrop(bl.Inventory(w, pos))

			ch, ok := w.Block(oldPos).(block.Chest)
			if !ok {
				break
			}
			h, ok := ch.Inventory(w, pos).Handler().(*airDropInventoryHandler)
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

func placeParachute(w *world.World, pos cube.Pos) {
	placeParachuteBlock(parachuteFenceOffsets, w, pos, block.WoodFence{Wood: block.OakWood()})
	placeParachuteBlock(parachuteWoolOffsets, w, pos, block.Wool{Colour: item.ColourRed()})
}

func removeParachute(w *world.World, pos cube.Pos) {
	removeParachuteBlock(parachuteFenceOffsets, w, pos)
	removeParachuteBlock(parachuteWoolOffsets, w, pos)
}

func placeParachuteBlock(offsets []cube.Pos, w *world.World, pos cube.Pos, bl world.Block) {
	for _, off := range offsets {
		newPos := pos.Add(off)
		if _, ok := w.Block(newPos).(block.Air); !ok {
			continue
		}
		w.SetBlock(newPos, bl, nil)
	}
}

func removeParachuteBlock(offsets []cube.Pos, w *world.World, pos cube.Pos) {
	for _, off := range offsets {
		newPos := pos.Add(off)
		_, fence := w.Block(newPos).(block.WoodFence)
		if _, wool := w.Block(newPos).(block.Wool); !fence && !wool {
			continue
		}
		w.SetBlock(newPos, block.Air{}, nil)
	}
}

func generateAirDrop(w *world.World) block.Chest {
	bl := block.NewChest()
	bl = bl.WithName(text.Colourf("<red>Air Drop</red>")).(block.Chest)
	inv := bl.Inventory(w, cube.Pos{})
	inv.Handle(&airDropInventoryHandler{
		inv: inv,
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
	items := it.PartnerItems()
	return it.NewSpecialItem(items[rand.Intn(len(items))], rand.Intn(5))
}

func findAirDropPosition(w *world.World) cube.Pos {
	spawn := area.Spawn(w)
	warzone := area.WarZone(w)

	for {
		x := warzone.Min().X() + (warzone.Max().X()-warzone.Min().X())*rand.Float64()
		z := warzone.Min().Y() + (warzone.Max().Y()-warzone.Min().Y())*rand.Float64()
		pos := cube.Pos{int(x), 255, int(z)}
		roads := area.Roads(w)
		valid := true
		for _, r := range roads {
			if r.Vec3WithinOrEqualFloorXZ(pos.Vec3Centre()) {
				valid = false
			}
		}
		if !spawn.Vec3WithinOrEqualFloorXZ(pos.Vec3Centre()) && valid {
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

func (h airDropInventoryHandler) HandleTake(_ *event.Context, _ int, st item.Stack) {
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

func (h airDropInventoryHandler) HandleDrop(_ *event.Context, _ int, st item.Stack) {
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

package internal

import (
	_ "embed"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/core/area"
	it "github.com/moyai-network/teams/internal/core/item"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math/rand"
	"time"
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
		w.Exec(func(tx *world.Tx) {
			destroyAirDrop(tx, lastDropPos)
			dropAirDrop(tx, pos)
		})
		return
	}
}

func dropAirDrop(tx *world.Tx, pos cube.Pos) {
	Broadcastf(tx, "airdrop.incoming", pos.X(), pos.Z())

	bl := generateAirDrop(tx)
	tx.SetBlock(pos, bl, nil)
	for p := range Players(tx) {
		p.PlaySound(sound.BarrelClose{})
		p.PlaySound(sound.FireworkBlast{})
		p.PlaySound(sound.FireworkTwinkle{})
		p.PlaySound(sound.Note{})
	}

	for {
		<-time.After(time.Second)

		oldPos := pos
		pos = pos.Add(cube.Pos{0, -1, 0})
		if _, ok := tx.Block(pos).(block.Air); !ok {
			removeParachute(tx, oldPos)
			placeParachuteBlock(parachuteFenceOffsets2, tx, pos, block.WoodFence{Wood: block.OakWood()})
			placeParachuteBlock(parachuteWoolOffsets2, tx, pos, block.Wool{Colour: item.ColourRed()})
			<-time.After(time.Second / 4)
			removeParachuteBlock(parachuteFenceOffsets2, tx, pos)
			removeParachuteBlock(parachuteWoolOffsets2, tx, pos)
			placeParachuteBlock(parachuteFenceOffsets3, tx, pos, block.WoodFence{Wood: block.OakWood()})
			placeParachuteBlock(parachuteWoolOffsets3, tx, pos, block.Wool{Colour: item.ColourRed()})
			<-time.After(time.Second / 4)
			removeParachuteBlock(parachuteFenceOffsets3, tx, pos)
			removeParachuteBlock(parachuteWoolOffsets3, tx, pos)
			placeParachuteBlock(parachuteFenceOffsets4, tx, pos, block.WoodFence{Wood: block.OakWood()})
			placeParachuteBlock(parachuteWoolOffsets4, tx, pos, block.Wool{Colour: item.ColourRed()})
			<-time.After(time.Second / 4)
			removeParachuteBlock(parachuteFenceOffsets4, tx, pos)
			removeParachuteBlock(parachuteWoolOffsets4, tx, pos)
			placeParachuteBlock(parachuteFenceOffsets5, tx, pos, block.WoodFence{Wood: block.OakWood()})
			placeParachuteBlock(parachuteWoolOffsets5, tx, pos, block.Wool{Colour: item.ColourRed()})
			<-time.After(time.Second / 4)
			removeParachuteBlock(parachuteFenceOffsets5, tx, pos)
			removeParachuteBlock(parachuteWoolOffsets5, tx, pos)
			fillAirDrop(bl.Inventory(tx, pos))

			ch, ok := tx.Block(oldPos).(block.Chest)
			if !ok {
				break
			}
			h, ok := ch.Inventory(tx, pos).Handler().(*airDropInventoryHandler)
			if !ok {
				break
			}
			h.pos = oldPos

			for _, p := range tx.Viewers(pos.Vec3()) {
				p.ViewSound(pos.Vec3Centre(), sound.Fall{})
			}

			lastDropPos = oldPos
			return
		}
		removeParachute(tx, oldPos)
		tx.SetBlock(oldPos, block.Air{}, nil)

		placeParachute(tx, pos)
		tx.SetBlock(pos, bl, nil)
	}
}

func placeParachute(tx *world.Tx, pos cube.Pos) {
	placeParachuteBlock(parachuteFenceOffsets, tx, pos, block.WoodFence{Wood: block.OakWood()})
	placeParachuteBlock(parachuteWoolOffsets, tx, pos, block.Wool{Colour: item.ColourRed()})
}

func removeParachute(tx *world.Tx, pos cube.Pos) {
	removeParachuteBlock(parachuteFenceOffsets, tx, pos)
	removeParachuteBlock(parachuteWoolOffsets, tx, pos)
}

func placeParachuteBlock(offsets []cube.Pos, tx *world.Tx, pos cube.Pos, bl world.Block) {
	for _, off := range offsets {
		newPos := pos.Add(off)
		if _, ok := tx.Block(newPos).(block.Air); !ok {
			continue
		}
		tx.SetBlock(newPos, bl, nil)
	}
}

func removeParachuteBlock(offsets []cube.Pos, tx *world.Tx, pos cube.Pos) {
	for _, off := range offsets {
		newPos := pos.Add(off)
		_, fence := tx.Block(newPos).(block.WoodFence)
		if _, wool := tx.Block(newPos).(block.Wool); !fence && !wool {
			continue
		}
		tx.SetBlock(newPos, block.Air{}, nil)
	}
}

func generateAirDrop(tx *world.Tx) block.Chest {
	bl := block.NewChest()
	bl = bl.WithName(text.Colourf("<red>Air Drop</red>")).(block.Chest)
	inv := bl.Inventory(tx, cube.Pos{})
	inv.Handle(&airDropInventoryHandler{
		inv: inv,
		w:   tx.World(),
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
		x := warzone.Min.X() + (warzone.Max.X()-warzone.Min.X())*rand.Float64()
		z := warzone.Min.Y() + (warzone.Max.Y()-warzone.Min.Y())*rand.Float64()
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

func (h airDropInventoryHandler) HandleTake(_ *inventory.Context, _ int, st item.Stack) {
	stacks := h.inv.Items()
	if len(stacks) == 1 && stacks[0].Equal(st) {
		time.AfterFunc(time.Second, func() {
			h.w.Exec(func(tx *world.Tx) {
				destroyAirDrop(tx, h.pos)
			})
		})
		return
	}
}
func (airDropInventoryHandler) HandlePlace(ctx *inventory.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

func (h airDropInventoryHandler) HandleDrop(_ *inventory.Context, _ int, st item.Stack) {
	stacks := h.inv.Items()
	if len(stacks) == 1 && stacks[0].Equal(st) {
		time.AfterFunc(time.Second, func() {
			h.w.Exec(func(tx *world.Tx) {
				destroyAirDrop(tx, h.pos)
			})
		})
		return
	}
}

func destroyAirDrop(tx *world.Tx, pos cube.Pos) {
	if _, ok := tx.Block(pos).(block.Air); ok {
		return
	}
	for _, v := range tx.Viewers(pos.Vec3()) {
		v.ViewSound(pos.Vec3Centre(), sound.Explosion{})
	}
	tx.SetBlock(pos, block.Air{}, nil)
}

package minecraft

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unsafe"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/bedrock-gophers/spawner/spawner"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/npc"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/crate"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	ent "github.com/moyai-network/teams/moyai/entity"
	"github.com/moyai-network/teams/moyai/menu"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var (
	shopSigns   = []shopSign{}
	cowSpawners = []cube.Pos{
		{-17, 63, -14},
		{-15, 63, -18},
		{-13, 63, -22},
		{-11, 63, -26},
	}
	endermanSpawners = []cube.Pos{
		{28, 63, -33},
		{27, 63, -36},
		{26, 63, -39},
		{25, 63, -42},
	}
)

func configureWorlds() {
	for _, w := range moyai.Worlds() {
		w.SetDifficulty(world.DifficultyHard)
		w.StopWeatherCycle()
		w.SetDefaultGameMode(world.GameModeSurvival)
		w.SetTime(6000)
		w.StopTime()
		w.SetTickRange(0)
		w.StopThundering()
		w.StopRaining()
		w.SetSpawn(cube.Pos{0, 80, 0})

		l := world.NewLoader(8, w, world.NopViewer{})
		l.Move(w.Spawn().Vec3Middle())
		l.Load(math.MaxInt)

		inv.PlaceFakeContainer(w, cube.Pos{0, 127, 0})
	}
}

func tickClearLag() {
	t := time.NewTicker(time.Minute / 2)
	defer t.Stop()

	for range t.C {
		for _, w := range moyai.Worlds() {
			clearAgedEntities(w)
		}
	}
}

func clearAgedEntities(w *world.World) {
	for _, e := range w.Entities() {
		if et, ok := e.(*entity.Ent); ok && et.Type() == (entity.ItemType{}) {
			age := fetchPrivateField[time.Duration](et, "age")
			if age > (time.Minute*5)/2 {
				w.RemoveEntity(e)
			}
		}
	}
}

// fetchPrivateField fetches a private field of a session.
func fetchPrivateField[T any](v any, name string) T {
	reflectedValue := reflect.ValueOf(v).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)
	privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

	return privateFieldValue.Interface().(T)
}

func placeSlapper() {
	w := moyai.Overworld()
	_ = npc.Create(npc.Settings{
		Name:       text.Colourf("<green>Click to use kits</green>"),
		Skin:       skin.Skin{},
		Scale:      1,
		Yaw:        215,
		MainHand:   item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, 1)),
		Helmet:     item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(ench.Protection{}, 1)),
		Chestplate: item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(ench.Protection{}, 1)),
		Leggings:   item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(ench.Protection{}, 1)),
		Boots:      item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(ench.Protection{}, 1)),

		Position: mgl64.Vec3{-7, 65, 38.5},
	}, w, func(p *player.Player) {
		if men, ok := menu.NewKitsMenu(p); ok {
			inv.SendMenu(p, men)
		}
	})
}

func placeSpawners() {
	w := moyai.Overworld()
	for _, pos := range cowSpawners {
		sp := spawner.New(ent.NewCow, pos.Vec3Centre(), w, time.Second*30, 25, true)
		w.SetBlock(pos, sp, nil)
	}
	for _, pos := range endermanSpawners {
		sp := spawner.New(ent.NewEnderman, pos.Vec3Centre(), w, time.Second*5, 25, true)
		w.SetBlock(pos, sp, nil)
	}
}

// shopSign is a sign that can be placed in the world to create a shop. It can be used to buy or sell items.
type shopSign struct {
	buy      bool
	it       world.Item
	quantity int
	price    int
	pos      cube.Pos
}

func placeShopSigns() {
	w := moyai.Overworld()
	for _, s := range shopSigns {
		var txt string
		state := text.Colourf("<green>[Buy]</green>")
		if !s.buy {
			state = text.Colourf("<red>[Sell]</red>")
		}

		name, _ := s.it.EncodeItem()
		txt = fmt.Sprintf("%s\n%s\n%d\n$%d", state, formatItemName(name), s.quantity, s.price)

		b := block.Sign{Front: block.SignText{
			Text: txt,
		}}
		w.SetBlock(s.pos, b, nil)
	}
}

func formatItemName(s string) string {
	split := strings.Split(s, "_")
	for i, str := range split {
		upperCasesPrefix := unicode.ToUpper(rune(str[0]))
		split[i] = string(upperCasesPrefix) + str[1:]
	}
	return strings.Join(split, " ")
}

// placeCrates places all crates in the world.
func placeCrates() {
	w := moyai.Overworld()
	for _, c := range crate.All() {
		b := block.NewChest()
		b.Facing = c.Facing().Direction()
		b.CustomName = text.Colourf("%s <grey>Crate</grey>", c.Name())

		*b.Inventory() = *inventory.New(27, nil)

		var items [27]item.Stack
		for i, r := range c.Rewards() {
			if r.Stack().Empty() {
				continue // Ignore this, ill fix it later
			}
			st := ench.AddEnchantmentLore(r.Stack())
			st = st.WithLore(append(st.Lore(), text.Colourf("<yellow>Chance: %d%%</yellow>", r.Chance()))...)
			items[i] = st
		}
		for i, s := range items {
			if s.Empty() {
				items[i] = item.NewStack(block.StainedGlass{Colour: item.ColourRed()}, 1)
			}
		}

		for s, i := range items {
			_ = b.Inventory().SetItem(s, i)
		}

		b.Inventory().Handle(crate.Handler{})

		w.SetBlock(cube.PosFromVec3(c.Position()), b, nil)
		t := entity.NewText(text.Colourf("%s <grey>Crate</grey>\n<yellow>Left click to open crate</yellow>\n<grey>Right click to see rewards</grey>", c.Name()), c.Position().Add(mgl64.Vec3{0, 2, 0}))
		w.AddEntity(t)
	}
}

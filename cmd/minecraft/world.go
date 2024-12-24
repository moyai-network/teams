package minecraft

import (
	"fmt"
	crate2 "github.com/moyai-network/teams/internal/adapter/crate"
	"github.com/moyai-network/teams/internal/core/enchantment"
	menu2 "github.com/moyai-network/teams/internal/core/menu"
	"math"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unsafe"

	"github.com/bedrock-gophers/inv/inv"
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
	"github.com/moyai-network/teams/internal"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var (
	shopSigns = []shopSign{
		{buy: true, it: block.Emerald{}, quantity: 16, price: 2500, pos: cube.Pos{28, 67, 8}, direction: cube.North},
		{buy: true, it: block.Diamond{}, quantity: 16, price: 2000, pos: cube.Pos{27, 67, 8}, direction: cube.North},
		{buy: true, it: block.Gold{}, quantity: 16, price: 1000, pos: cube.Pos{26, 67, 8}, direction: cube.North},
		{buy: true, it: block.Iron{}, quantity: 16, price: 500, pos: cube.Pos{25, 67, 8}, direction: cube.North},
		{buy: true, it: block.Lapis{}, quantity: 16, price: 500, pos: cube.Pos{24, 67, 8}, direction: cube.North},

		{it: block.Emerald{}, quantity: 16, price: 2000, pos: cube.Pos{28, 68, 8}, direction: cube.North},
		{it: block.Diamond{}, quantity: 16, price: 1500, pos: cube.Pos{27, 68, 8}, direction: cube.North},
		{it: block.Gold{}, quantity: 16, price: 500, pos: cube.Pos{26, 68, 8}, direction: cube.North},
		{it: block.Iron{}, quantity: 16, price: 250, pos: cube.Pos{25, 68, 8}, direction: cube.North},
		{it: block.Lapis{}, quantity: 16, price: 250, pos: cube.Pos{24, 68, 8}, direction: cube.North},

		//{buy: true, it: block.RedstoneWire{}, quantity: 16, price: 100, pos: cube.Pos{28, 69, 8}, direction: cube.North},
		//{buy: true, it: block.Lever{}, quantity: 8, price: 100, pos: cube.Pos{27, 69, 8}, direction: cube.North},
		//{buy: true, it: block.Button{}, quantity: 8, price: 100, pos: cube.Pos{26, 69, 8}, direction: cube.North},
		//{buy: true, it: block.RedstoneBlock{}, quantity: 8, price: 100, pos: cube.Pos{25, 69, 8}, direction: cube.North},
		{buy: true, it: block.Hopper{}, quantity: 8, price: 3000, pos: cube.Pos{24, 69, 8}, direction: cube.North},
	}
	cowSpawners = []cube.Pos{
		{-17, 63, -14},
		{-15, 63, -18},
		{-13, 63, -22},
		{-11, 63, -26},
		{-17, 63, -14},
		{-18, 63, -17},
		{-19, 63, -20},
		{-14, 63, -15},
		{-16, 63, -21},
		{-17, 63, -24},
		{-14, 63, -25},
		{-12, 63, -19},
		{-11, 63, -16},
		{-8, 63, -17},
		{-9, 63, -20},
		{-10, 63, -23},
	}
	endermanSpawners = []cube.Pos{
		{27, 28, 121},
		{27, 28, 131},
		{21, 28, 131},
		{23, 28, 134},
		{18, 28, 138},
		{12, 28, 138},
		{14, 28, 126},
		{16, 28, 132},

		{-13, 25, 114},
		{-21, 25, 115},
		{-24, 25, 120},
		{-33, 27, 110},
		{-36, 27, 107},
		{-40, 27, 110},
		{-38, 28, 103},

		{-60, 28, 62},
		{-56, 28, 59},
		{-51, 28, 62},

		{23, 27, 40},
		{19, 27, 44},
		{26, 27, 45},
		{30, 27, 39},
		{27, 27, 52},
	}
	blazeSpawners = []cube.Pos{
		{-513, 83, -188},
		{-508, 83, -198},
		{-520, 83, -188},
		{-519, 83, -199},

		{-437, 87, -187},
		{-434, 87, -177},
		{-427, 87, -175},
	}
)

func configureWorlds() {
	for _, w := range internal.Worlds() {
		w.Handle(&worldHandler{w: w})
		w.SetDifficulty(world.DifficultyHard)
		w.StopWeatherCycle()
		w.SetDefaultGameMode(world.GameModeSurvival)
		w.SetTime(6000)
		w.StopTime()
		w.SetTickRange(1)
		w.StopThundering()
		w.StopRaining()
		w.SetSpawn(cube.Pos{0, 80, 0})

		l := world.NewLoader(8, w, world.NopViewer{})
		w.Exec(func(tx *world.Tx) {
			l.Move(tx, w.Spawn().Vec3Middle())
			l.Load(tx, math.MaxInt)
		})
	}
}

func tickClearLag() {
	t := time.NewTicker(time.Minute / 2)
	defer t.Stop()

	for range t.C {
		for _, w := range internal.Worlds() {
			clearAgedEntities(w)
		}
	}
}

func clearAgedEntities(w *world.World) {
	/*for _, e := range w.Entities() {
		if et, ok := e.(*entity.Ent); ok && et.Type() == (entity.ItemType{}) {
			age := fetchPrivateField[time.Duration](et, "age")
			if age > (time.Minute*5)/2 {
				w.RemoveEntity(e)
			}
		}
	}*/
}

// fetchPrivateField fetches a private field of a session.
func fetchPrivateField[T any](v any, name string) T {
	reflectedValue := reflect.ValueOf(v).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)
	privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

	return privateFieldValue.Interface().(T)
}

func placeSlapper() {
	internal.Overworld().Exec(func(w *world.Tx) {
		_ = npc.Create(npc.Settings{
			Name:       text.Colourf("<green>Click to use kits</green>"),
			Skin:       skin.Skin{},
			Scale:      1,
			Yaw:        120,
			MainHand:   item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Sharpness{}, 1)),
			Helmet:     item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection{}, 1)),
			Chestplate: item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection{}, 1)),
			Leggings:   item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection{}, 1)),
			Boots:      item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection{}, 1)),

			Position: mgl64.Vec3{7, 67, 47.5},
		}, w, func(p *player.Player) {
			if men, ok := menu2.NewKitsMenu(p); ok {
				inv.SendMenu(p, men)
			}
		})
		_ = npc.Create(npc.Settings{
			Name:     text.Colourf("<gold>Block Shop</gold>"),
			Skin:     skin.Skin{},
			Scale:    1,
			Yaw:      -120,
			MainHand: item.NewStack(block.Diamond{}, 1),

			Position: mgl64.Vec3{-6, 67, 45.5},
		}, w, func(p *player.Player) {
			inv.SendMenu(p, menu2.NewBlocksMenu(p))
		})
	})
}

// shopSign is a sign that can be placed in the world to create a shop. It can be used to buy or sell items.
type shopSign struct {
	buy       bool
	it        world.Item
	quantity  int
	price     int
	pos       cube.Pos
	direction cube.Direction
}

func placeShopSigns() {
	w := internal.Overworld()
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
		b.Attach = block.WallAttachment(s.direction)
		w.Exec(func(tx *world.Tx) {
			tx.SetBlock(s.pos, b, nil)
		})
	}
}

func formatItemName(s string) string {
	split := strings.Split(s, ":")

	var formattedParts []string

	for i, str := range split {
		if i == 0 {
			continue
		}
		words := strings.Split(str, "_")
		for j, word := range words {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[j] = string(runes)
		}
		formattedParts = append(formattedParts, strings.Join(words, " "))
	}

	return strings.Join(formattedParts, " ")
}

// placeCrates places all crates in the world.
func placeCrates() {
	for _, c := range crate2.All() {
		internal.Overworld().Exec(func(w *world.Tx) {
			b := block.NewChest()
			b.Facing = c.Facing().Direction()
			b.CustomName = text.Colourf("%s <grey>Crate</grey>", c.Name())

			pos := cube.PosFromVec3(c.Position())
			*b.Inventory(w, pos) = *inventory.New(27, nil)

			var items [27]item.Stack
			for i, r := range c.Rewards() {
				if r.Stack().Empty() {
					continue // Ignore this, ill fix it later
				}
				st := enchantment.AddEnchantmentLore(r.Stack())
				st = st.WithLore(append(st.Lore(), text.Colourf("<yellow>Chance: %d%%</yellow>", r.Chance()))...)
				items[i] = st
			}
			for i, s := range items {
				if s.Empty() {
					items[i] = item.NewStack(block.StainedGlass{Colour: item.ColourRed()}, 1)
				}
			}

			for s, i := range items {
				_ = b.Inventory(w, pos).SetItem(s, i)
			}

			b.Inventory(w, pos).Handle(crate2.Handler{})

			w.SetBlock(cube.PosFromVec3(c.Position()), b, nil)
			t := entity.NewText(text.Colourf("%s <grey>Crate</grey>\n<yellow>Right click to open crate</yellow>\n<grey>Left click to see rewards</grey>", c.Name()), c.Position().Add(mgl64.Vec3{0, 2, 0}))
			w.AddEntity(t)
		})
	}
}

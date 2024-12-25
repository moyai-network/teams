package minecraft

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"
	crate2 "github.com/moyai-network/teams/internal/adapter/crate"
	"github.com/moyai-network/teams/internal/core/enchantment"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
	"strings"
	"unicode"
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

		{buy: true, it: block.Hopper{}, quantity: 8, price: 3000, pos: cube.Pos{24, 69, 8}, direction: cube.North},
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

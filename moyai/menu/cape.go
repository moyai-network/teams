package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/moyai/cape"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/roles"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Cape struct{}

func NewCape(p *player.Player) inv.Menu {
	m := inv.NewMenu(Cape{}, "Cape", inv.ContainerChest{})
	stacks := glassFilledStack(54)

	u, _ := data.LoadUserFromName(p.Name())

	for i, c := range cape.All() {
		if u.Teams.Settings.Advanced.Cape == c.Name() {
			stacks[i] = item.NewStack(block.Banner{Colour: item.ColourBlue()}, 1).
				WithCustomName(text.Colourf("<blue>%s</blue>", c.Name())).
				WithLore(text.Colourf("<green>Selected</green>")).
				WithEnchantments(item.NewEnchantment(glint{}, 1))
			continue
		}
		col := item.ColourGreen()
		name := text.Colourf("<green>%s</green>", c.Name())
		if !roles.Premium(u.Roles.Highest()) && c.Premium() {
			col = item.ColourRed()
			name = text.Colourf("<red>%s</red>", c.Name())
		}
		stacks[i] = item.NewStack(block.Banner{Colour: col}, 1).WithCustomName(name)
	}

	return m.WithStacks(stacks...)
}

func (Cape) Submit(p *player.Player, it item.Stack) {
	b, ok := it.Item().(block.Banner)
	if !ok {
		return
	}

	if b.Colour == item.ColourRed() {
		p.Message(text.Colourf("<red>You need a premium rank to use this cape.</red>"))
		inv.CloseContainer(p)
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil || b.Colour == item.ColourBlue() {
		return
	}

	c, ok := cape.ByName(colour.StripMinecraftColour(it.CustomName()))
	if !ok {
		return
	}
	u.Teams.Settings.Advanced.Cape = c.Name()

	sk := p.Skin()
	sk.Cape = c.Cape()
	p.SetSkin(sk)

	data.SaveUser(u)
	p.PlaySound(sound.Experience{})

	inv.UpdateMenu(p, NewCape(p))
}

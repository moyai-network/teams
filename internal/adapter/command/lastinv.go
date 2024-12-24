package command

import (
	"fmt"
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type LastInv struct {
	adminAllower
	Target []cmd.Target `cmd:"target"`
}

func (i LastInv) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tg, ok := i.Target[0].(*player.Player)
	if !ok {
		return
	}

	t, err := data.LoadUserFromName(tg.Name())
	if err != nil {
		internal.Messagef(p, "command.target.unknown", tg.Name())
		return
	}
	d := t.Teams.DeathInventory
	iv := inventory.New(54, nil)
	for i, d := range d.Items {
		var x int
		if i <= 35 {
			x = -27
		}
		if i <= 26 {
			x = -9
		}
		if i <= 17 {
			x = 9
		}
		if i <= 8 {
			x = 27
		}
		iv.SetItem(i+x, d)
	}
	iv.SetItem(45, d.Helmet)
	iv.SetItem(46, d.Chestplate)
	iv.SetItem(47, d.Leggings)
	iv.SetItem(48, d.Boots)

	rev := item.NewStack(item.EnchantedBook{}, 1).
		WithCustomName(text.Colourf("<red>Last Inventory Info</red>")).
		WithLore(text.Colourf("<red>Dead</red>: %t", t.Teams.DeathBan.Active()), text.Colourf("<green>Click to revive user with current inventory</green>"))

	iv.SetItem(53, rev)

	iv.Handle(&handler{})
	menu := inv.NewCustomMenu(fmt.Sprintf("Last Inventory of %s", t.DisplayName), inv.ContainerChest{DoubleChest: true}, iv, nil)

	inv.SendMenu(p, menu)
}

type handler struct {
	inventory.NopHandler
}

func (*handler) HandleTake(ctx *inventory.Context, a int, it item.Stack) {
	ctx.Cancel()
}

func (*handler) HandlePlace(ctx *inventory.Context, a int, it item.Stack) {
	ctx.Cancel()
}

func (*handler) HandleDrop(ctx *inventory.Context, a int, it item.Stack) {
	ctx.Cancel()
}

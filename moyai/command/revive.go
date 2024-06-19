package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
)

type Revive struct {
	Target []cmd.Target `cmd:"target"`
}

func (r Revive) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	target, ok := r.Target[0].(*player.Player)
	if !ok {
		return
	}
	tg, err := data.LoadUserFromName(target.Name())
	if err != nil {
		return
	}
	tg.Teams.Stats.Deaths--
	inv := tg.Teams.DeathInventory
	addDataInventory(p, *inv)

	data.SaveUser(tg)
}

func addDataInventory(p *player.Player, inv data.Inventory) {
	for _, i := range inv.Items {
		it.AddOrDrop(p, i)
	}
	it.AddArmorOrDrop(p, inv.Boots)
	it.AddArmorOrDrop(p, inv.Leggings)
	it.AddArmorOrDrop(p, inv.Chestplate)
	it.AddArmorOrDrop(p, inv.Helmet)
	it.AddOrDrop(p, inv.OffHand)
}

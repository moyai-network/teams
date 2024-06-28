package command

import (
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
)

type Revive struct {
	adminAllower
	Target []cmd.Target `cmd:"target"`
}

func (r Revive) Run(src cmd.Source, _ *cmd.Output) {
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

	if tg.Teams.DeathBan.After(time.Now()) {
		tg.Teams.DeathBan = time.Time{}
		moyai.Overworld().AddEntity(target)
		target.Teleport(mgl64.Vec3{0, 80, 0})
	}

	addDataInventory(target, *inv)

	data.SaveUser(tg)}

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

package command

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/core"
	it "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/model"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type Revive struct {
	adminAllower
	Target []cmd.Target `cmd:"target"`
}

func (r Revive) Run(src cmd.Source, _ *cmd.Output, tx *world.Tx) {
	target, ok := r.Target[0].(*player.Player)
	if !ok {
		return
	}
	tg, ok := core.UserRepository.FindByName(target.Name())
	if !ok {
		return
	}
	tg.Teams.Stats.Deaths--
	inv := tg.Teams.DeathInventory

	if tg.Teams.DeathBan.Active() {
		internal.Overworld().Exec(func(tx *world.Tx) {
			tx.AddEntity(target.H())
		})
		target.Teleport(mgl64.Vec3{0, 80, 0})

		target.Inventory().Clear()
		target.Armour().Clear()

		tg.Teams.DeathBan.Reset()
		tg.Teams.DeathBanned = false

		tg.Teams.PVP.Set(time.Hour + (time.Millisecond * 500))
		if !tg.Teams.PVP.Paused() {
			tg.Teams.PVP.TogglePause()
		}
	}

	addDataInventory(target, *inv)
	core.UserRepository.Save(tg)
}

func addDataInventory(p *player.Player, inv model.Inventory) {
	for _, i := range inv.Items {
		it.AddOrDrop(p, i)
	}
	it.AddArmorOrDrop(p, inv.Boots)
	it.AddArmorOrDrop(p, inv.Leggings)
	it.AddArmorOrDrop(p, inv.Chestplate)
	it.AddArmorOrDrop(p, inv.Helmet)
	it.AddOrDrop(p, inv.OffHand)
}

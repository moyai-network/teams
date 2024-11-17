package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/data"
	"github.com/moyai-network/teams/internal/kit"
	"github.com/moyai-network/teams/internal/user"
)

type StaffMode struct {
	trialAllower
	Sub cmd.SubCommand `cmd:"mode"`
}

func (StaffMode) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	mode := p.GameMode()
	u.StaffMode = !u.StaffMode

	if u.StaffMode {
		_ = internal.PlayerProvider().SavePlayer(p)
		if !u.Vanished {
			p.SetGameMode(vanishGameMode{lastMode: mode})
			internal.Messagef(p, "command.vanish.enabled")
			u.Vanished = true
		}

		p.Inventory().Clear()
		p.Armour().Clear()
		kit.Apply(kit.Staff{}, p)
		p.Inventory().Handle(user.StaffInventoryHandler{})
	} else {
		if u.Vanished {
			u.Vanished = false

			vanishMode, ok := mode.(vanishGameMode)
			if !ok {
				p.SetGameMode(world.GameModeCreative)
			} else {
				p.SetGameMode(vanishMode.lastMode)
			}
			internal.Messagef(p, "command.vanish.disabled")
		}

		p.Inventory().Handle(inventory.NopHandler{})
		dat, err := internal.LoadPlayerData(p.UUID())
		if err != nil {
			return
		}
		p.Teleport(dat.Position)
		p.SetGameMode(dat.GameMode)

		newInv := dat.Inventory
		p.Inventory().Clear()
		p.Armour().Clear()
		p.Armour().Set(newInv.Helmet, newInv.Chestplate, newInv.Leggings, newInv.Boots)
		for slot, it := range newInv.Items {
			if it.Empty() {
				continue
			}
			_ = p.Inventory().SetItem(slot, it)
		}
	}

	data.SaveUser(u)
	user.UpdateVanishState(p, u)
}

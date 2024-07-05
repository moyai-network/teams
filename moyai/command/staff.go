package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/kit"
	"github.com/moyai-network/teams/moyai/user"
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

	u, err := data.LoadUserFromXUID(p.XUID())
	if err != nil {
		return
	}
	mode := p.GameMode()

	if u.Vanished {
		//moyai.Alertf(s, "staff.alert.vanish.off")
		vanishMode, ok := mode.(vanishGameMode)
		if !ok {
			p.SetGameMode(world.GameModeCreative)
		} else {
			p.SetGameMode(vanishMode.lastMode)
		}
		moyai.Messagef(p, "command.vanish.disabled")
	} else {
		//moyai.Alertf(s, "staff.alert.vanish.on")
		p.SetGameMode(vanishGameMode{lastMode: mode})
		moyai.Messagef(p, "command.vanish.enabled")
		u.Vanished = true
	}
	u.StaffMode = !u.StaffMode

	if u.StaffMode {
		*u.PlayerData.Inventory = data.InventoryData(p)
		u.PlayerData.Position = p.Position()
		u.PlayerData.GameMode, _ = world.GameModeID(p.GameMode())

		p.Inventory().Clear()
		p.Armour().Clear()
		kit.Apply(kit.Staff{}, p)
		p.Inventory().Handle(user.StaffInventoryHandler{})
	} else {
		p.Inventory().Handle(inventory.NopHandler{})
		u.PlayerData.Inventory.Apply(p)
		p.Teleport(u.PlayerData.Position)
		mode, _ = world.GameModeByID(u.PlayerData.GameMode)
		p.SetGameMode(mode)
	}

	data.SaveUser(u)
	user.UpdateVanishState(p, u)
}

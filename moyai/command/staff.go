package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/user"
)

type StaffMode struct {
	adminAllower
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
		//user.Alertf(s, "staff.alert.vanish.off")
		vanishMode, ok := mode.(vanishGameMode)
		if !ok {
			return
		}
		p.SetGameMode(vanishMode.lastMode)
		user.Messagef(p, "command.vanish.disabled")
	} else {
		//user.Alertf(s, "staff.alert.vanish.on")
		p.SetGameMode(vanishGameMode{lastMode: mode})
		user.Messagef(p, "command.vanish.enabled")
	}

	u.StaffMode = !u.StaffMode
	u.Vanished = !u.Vanished
	data.SaveUser(u)
	user.UpdateVanishState(p, u)
}

func (StaffMode) Allow(s cmd.Source) bool {
	return allow(s, false, role.Trial{})
}

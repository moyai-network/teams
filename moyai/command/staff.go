package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
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
			return
		}
		p.SetGameMode(vanishMode.lastMode)
		moyai.Messagef(p, "command.vanish.disabled")
	} else {
		//moyai.Alertf(s, "staff.alert.vanish.on")
		p.SetGameMode(vanishGameMode{lastMode: mode})
		moyai.Messagef(p, "command.vanish.enabled")
	}

	u.StaffMode = !u.StaffMode
	u.Vanished = !u.Vanished
	data.SaveUser(u)
	user.UpdateVanishState(p, u)
}

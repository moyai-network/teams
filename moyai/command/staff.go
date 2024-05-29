package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
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

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}
	mode := p.GameMode()

	if h.Vanished() {
		vanishMode, ok := mode.(vanishGameMode)
		if !ok {
			return
		}
		p.SetGameMode(vanishMode.lastMode)
		user.Messagef(p, "command.vanish.disabled")
	} else {
		p.SetGameMode(vanishGameMode{lastMode: mode})
		user.Messagef(p, "command.vanish.enabled")
	}

	h.ToggleVanish()
}

func (StaffMode) Allow(s cmd.Source) bool {
	return allow(s, false, role.Trial{})
}

package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/user"
)

// PvpEnable is a command to enable PVP.
type PvpEnable struct {
	Sub cmd.SubCommand `cmd:"enable"`
}

// Run ...
func (c PvpEnable) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}

	if u.Teams.PVP.Active() {
		user.UpdateState(p)
		u.Teams.PVP.Reset()
		internal.Messagef(p, "command.pvp.enable")
	} else {
		internal.Messagef(p, "command.pvp.enabled-already")
	}
}

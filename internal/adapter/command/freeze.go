package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/user"
)

// Freeze is a command used to freeze a player.
type Freeze struct {
	trialAllower
	Targets []cmd.Target `cmd:"target"`
}

// Run ...
func (f Freeze) Run(src cmd.Source, _ *cmd.Output, tx *world.Tx) {
	if len(f.Targets) > 1 {
		internal.Messagef(src, "command.targets.exceed")
		return
	}
	target, ok := f.Targets[0].(*player.Player)
	if !ok {
		internal.Messagef(src, "command.target.unknown")
		return
	}
	if src == target {
		internal.Messagef(src, "command.usage.self")
		return
	}
	t, ok := user.Lookup(tx, target.Name())
	if !ok {
		internal.Messagef(src, "command.target.unknown")
		return
	}
	u, ok := core.UserRepository.FindByName(t.Name())
	if !ok {
		return
	}
	if u.Frozen {
		internal.Alertf(tx, src, "staff.alert.unfreeze", target.Name())
		internal.Messagef(src, "command.freeze.unfreeze", target.Name())
		internal.Messagef(t, "command.freeze.unfrozen")
		t.SetMobile()
	} else {
		internal.Alertf(tx, src, "staff.alert.freeze", target.Name())
		internal.Messagef(src, "command.freeze.freeze", target.Name())
		internal.Messagef(t, "command.freeze.frozen")
		t.SetImmobile()
	}
	u.Frozen = !u.Frozen
	core.UserRepository.Save(u)
}

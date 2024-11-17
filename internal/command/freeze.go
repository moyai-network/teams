package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/data"
	"github.com/moyai-network/teams/internal/user"
)

// Freeze is a command used to freeze a player.
type Freeze struct {
	trialAllower
	Targets []cmd.Target `cmd:"target"`
}

// Run ...
func (f Freeze) Run(src cmd.Source, _ *cmd.Output) {
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
	t, ok := user.Lookup(target.Name())
	if !ok {
		internal.Messagef(src, "command.target.unknown")
		return
	}
	u, err := data.LoadUserFromName(t.Name())
	if err != nil {
		return
	}
	if u.Frozen {
		internal.Alertf(src, "staff.alert.unfreeze", target.Name())
		internal.Messagef(src, "command.freeze.unfreeze", target.Name())
		internal.Messagef(t, "command.freeze.unfrozen")
		t.SetMobile()
	} else {
		internal.Alertf(src, "staff.alert.freeze", target.Name())
		internal.Messagef(src, "command.freeze.freeze", target.Name())
		internal.Messagef(t, "command.freeze.frozen")
		t.SetImmobile()
	}
	u.Frozen = !u.Frozen
	data.SaveUser(u)
}

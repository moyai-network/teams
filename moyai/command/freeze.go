package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/user"
)

// Freeze is a command used to freeze a player.
type Freeze struct {
	Targets []cmd.Target `cmd:"target"`
}

// Run ...
func (f Freeze) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	if len(f.Targets) > 1 {
		user.Messagef(p, "command.targets.exceed")
		return
	}
	target, ok := f.Targets[0].(*player.Player)
	if !ok {
		user.Messagef(p, "command.target.unknown")
		return
	}
	if s == target {
		user.Messagef(p, "command.usage.self")
		return
	}
	t, ok := user.Lookup(target.Name())
	if !ok {
		user.Messagef(p, "command.target.unknown")
		return
	}
	u, err := data.LoadUserFromName(t.Name())
	if err != nil {
		return
	}
	if u.Frozen {
		//user.Alertf(s, "staff.alert.unfreeze", target.Name())
		//o.Print(lang.Translatef(l, "command.freeze.unfreeze", target.Name()))
		//t.Player().Message(lang.Translatef(t.Player().Locale(), "command.freeze.unfrozen"))
		t.SetMobile()
	} else {
		//user.Alertf(s, "staff.alert.freeze", target.Name())
		//o.Print(lang.Translatef(l, "command.freeze.freeze", target.Name()))
		//t.Player().Message(lang.Translatef(t.Player().Locale(), "command.freeze.frozen"))
		t.Immobile()
	}
	u.Frozen = !u.Frozen
	data.SaveUser(u)
}

// Allow ...
func (f Freeze) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

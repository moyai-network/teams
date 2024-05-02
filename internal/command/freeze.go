package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/data"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/user"
)

// Freeze is a command used to freeze a player.
type Freeze struct {
	Targets []cmd.Target `cmd:"target"`
}

// Run ...
func (f Freeze) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	if len(f.Targets) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	target, ok := f.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if s == target {
		o.Error(lang.Translatef(l, "command.usage.self"))
		return
	}
	t, ok := user.Lookup(target.Name())
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	u, err := data.LoadUserOrCreate(t.Player().Name())
	if err != nil {
		return
	}
	if u.Frozen {
		user.Alert(s, "staff.alert.unfreeze", target.Name())
		o.Print(lang.Translatef(l, "command.freeze.unfreeze", target.Name()))
		t.Player().Message(lang.Translatef(t.Player().Locale(), "command.freeze.unfrozen"))
		t.Player().SetMobile()
	} else {
		user.Alert(s, "staff.alert.freeze", target.Name())
		o.Print(lang.Translatef(l, "command.freeze.freeze", target.Name()))
		t.Player().Message(lang.Translatef(t.Player().Locale(), "command.freeze.frozen"))
		t.Player().Immobile()
	}
	u.Frozen = !u.Frozen
	data.SaveUser(u)
}

// Allow ...
func (f Freeze) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

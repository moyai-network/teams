package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
)

// Freeze is a command used to freeze a player.
type Freeze struct {
	trialAllower
	Targets []cmd.Target `cmd:"target"`
}

// Run ...
func (f Freeze) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	if len(f.Targets) > 1 {
		moyai.Messagef(p, "command.targets.exceed")
		return
	}
	target, ok := f.Targets[0].(*player.Player)
	if !ok {
		moyai.Messagef(p, "command.target.unknown")
		return
	}
	if s == target {
		moyai.Messagef(p, "command.usage.self")
		return
	}
	t, ok := user.Lookup(target.Name())
	if !ok {
		moyai.Messagef(p, "command.target.unknown")
		return
	}
	u, err := data.LoadUserFromName(t.Name())
	if err != nil {
		return
	}
	if u.Frozen {
		moyai.Alertf(s, "staff.alert.unfreeze", target.Name())
		o.Print(lang.Translatef(*u.Language, "command.freeze.unfreeze", target.Name()))
		t.Message(lang.Translatef(*u.Language, "command.freeze.unfrozen"))
		t.SetMobile()
	} else {
		moyai.Alertf(s, "staff.alert.freeze", target.Name())
		o.Print(lang.Translatef(*u.Language, "command.freeze.freeze", target.Name()))
		t.Message(lang.Translatef(*u.Language, "command.freeze.frozen"))
		t.SetImmobile()
	}
	u.Frozen = !u.Frozen
	data.SaveUser(u)
}

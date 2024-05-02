package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/data"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
)

// Kick is a command that disconnects another player from the server.
type Kick struct {
	Targets []cmd.Target `cmd:"target"`
}

// Run ...
func (k Kick) Run(s cmd.Source, o *cmd.Output) {
	l, single := locale(s), true
	if len(k.Targets) > 1 {
		if p, ok := s.(*player.Player); ok {
			if u, err := data.LoadUserOrCreate(p.Name()); err == nil && !u.Roles.Contains(role.Operator{}) {
				o.Error(lang.Translatef(l, "command.targets.exceed"))
				return
			}
		}
		single = false
	}

	var kicked int
	for _, p := range k.Targets {
		if p, ok := p.(*player.Player); ok {
			u, err := data.LoadUserOrCreate(p.Name())
			if err != nil || u.Roles.Contains(role.Operator{}) {
				o.Print(lang.Translatef(l, "command.kick.fail"))
				continue
			}
			p.Disconnect(lang.Translatef(p.Locale(), "command.kick.reason"))
			if single {
				//webhook.SendPunishment(s.Name(), p.Name(), "", "Kick")
				o.Print(lang.Translatef(l, "command.kick.success", p.Name()))
				return
			}
			kicked++
		} else if single {
			o.Print(lang.Translatef(l, "command.target.unknown"))
			return
		}
	}
	if !single {
		return
	}
	o.Print(lang.Translatef(l, "command.kick.multiple", kicked))
}

// Allow ...
func (Kick) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	data2 "github.com/moyai-network/teams/internal/core/data"
	rls "github.com/moyai-network/teams/internal/core/roles"
	"github.com/moyai-network/teams/pkg/lang"
	"golang.org/x/text/language"
)

// Kick is a command that disconnects another player from the server.
type Kick struct {
	trialAllower
	Targets []cmd.Target `cmd:"target"`
}

// Run ...
func (k Kick) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l, single := locale(s), true
	if len(k.Targets) > 1 {
		if p, ok := s.(*player.Player); ok {
			if u, err := data2.LoadUserFromName(p.Name()); err == nil && !u.Roles.Contains(rls.Operator()) {
				o.Error(lang.Translatef(l, "command.targets.exceed"))
				return
			}
		}
		single = false
	}

	var kicked int
	for _, p := range k.Targets {
		if p, ok := p.(*player.Player); ok {
			u, err := data2.LoadUserFromName(p.Name())
			if err != nil || u.Roles.Contains(rls.Operator()) {
				o.Print(lang.Translatef(l, "command.kick.fail"))
				continue
			}
			p.Disconnect(lang.Translatef(data2.Language{Tag: language.English}, "command.kick.reason"))
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

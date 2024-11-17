package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/data"
	"github.com/moyai-network/teams/pkg/lang"
)

// Ping represents the Ping command.
type Ping struct {
	Targets cmd.Optional[[]cmd.Target] `cmd:"target"`
}

// Run ...
func (p Ping) Run(s cmd.Source, o *cmd.Output) {
	targets := p.Targets.LoadOr(nil)
	if len(targets) > 1 {
		o.Error(lang.Translatef(data.Language{}, "command.targets.exceed"))
		return
	}

	if len(targets) == 1 {
		target, ok := targets[0].(*player.Player)
		if !ok {
			o.Error(lang.Translatef(data.Language{}, "command.target.unknown", ""))
			return
		}
		internal.Messagef(s, "command.ping.output", target.Name(), (target.Latency() * 2).Milliseconds())
		return
	}

	if p, ok := s.(*player.Player); ok {
		internal.Messagef(p, "command.ping.output", p.Name(), (p.Latency() * 2).Milliseconds())
	}
}

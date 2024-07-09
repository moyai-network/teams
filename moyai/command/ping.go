package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
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
			o.Error(lang.Translatef(data.Language{}, "command.target.unknown"))
			return
		}
		moyai.Messagef(target, "command.ping.output", target.Name(), (target.Latency() * 2).Milliseconds())
		return
	}

	if p, ok := s.(*player.Player); ok {
		moyai.Messagef(p, "command.ping.output", p.Name(), (p.Latency() * 2).Milliseconds())
	}
}

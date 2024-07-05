package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
)

// Ping represents the Ping command.
type Ping struct {
	Target cmd.Optional[[]cmd.Target] `cmd:"target"`
}

// Run ...
func (p Ping) Run(src cmd.Source, out *cmd.Output) {
	var t cmd.Target
	pl, _ := src.(*player.Player)
	if ta, ok := p.Target.Load(); ok { 
		t = ta[0]
	} else {
		t = pl
	}

	if p, ok := t.(*player.Player); ok {
		moyai.Messagef(p, "command.ping.output", p.Name(), (p.Latency() * 2).Milliseconds())
	}
}

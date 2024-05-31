package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
)

// Ping represents the Ping command.
type Ping struct {
	Target cmd.Optional[[]cmd.Target]
}

// Run ...
func (p Ping) Run(src cmd.Source, out *cmd.Output) {
	var t []cmd.Target
	pl, ok := src.(*player.Player)
	t = append(t, pl)
	if !ok {
		t = p.Target.LoadOr(t)
	}
	if pl, ok := t[0].(*player.Player); ok {
		moyai.Messagef(pl, "command.ping.output", pl.Name(), (pl.Latency() * 2).Milliseconds())
	}
}

package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal"
)

// Nether implements the /end command.
type Nether struct{ operatorAllower }

// Run ...
func (Nether) Run(src cmd.Source, o *cmd.Output) {
	p, _ := src.(*player.Player)
	internal.Nether().AddEntity(p)
}

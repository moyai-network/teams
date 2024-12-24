package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

// Nether implements the /end command.
type Nether struct{ operatorAllower }

// Run ...
func (Nether) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	/*p, _ := src.(*player.Player)
	internal.Nether().AddEntity(p)*/
	panic("todo")
}

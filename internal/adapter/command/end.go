package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

// End implements the /end command.
type End struct{ operatorAllower }

// Run ...
func (e End) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	/*p, _ := src.(*player.Player)
	internal.End().AddEntity(p)*/
	panic("todo")
}

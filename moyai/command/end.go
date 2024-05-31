package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
)

// End implements the /end command.
type End struct{ operatorAllower }

// Run ...
func (e End) Run(src cmd.Source, o *cmd.Output) {
	p, _ := src.(*player.Player)
	moyai.End().AddEntity(p)
}

package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/role"
)

// End implements the /end command.
type End struct{}

// Run ...
func (e End) Run(src cmd.Source, o *cmd.Output) {
	p, _ := src.(*player.Player)
	moyai.End().AddEntity(p)
}

// Allow ...
func (End) Allow(s cmd.Source) bool {
	return allow(s, false, role.Admin{})
}
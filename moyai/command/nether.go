package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/role"
)

// Nether implements the /end command.
type Nether struct{}

// Run ...
func (Nether) Run(src cmd.Source, o *cmd.Output) {
	p, _ := src.(*player.Player)
	moyai.Nether().AddEntity(p)
}

// Allow ...
func (Nether) Allow(s cmd.Source) bool {
	return allow(s, false, role.Admin{})
}

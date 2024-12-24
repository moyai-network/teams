package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/team"
)

// TL is a command that allows players to see the coordinates of their team members.
type TL struct{}

// Run ...
func (TL) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		internal.Messagef(p, "user.team-less")
		return
	}

	for _, m := range team.OnlineMembers(tx, tm) {
		internal.Messagef(m, "command.tl", p.Name(), int(p.Position().X()), int(p.Position().Y()), int(p.Position().Z()))
	}
}

package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/team"
	"github.com/moyai-network/teams/moyai/user"
)

// TL is a command that allows players to see the coordinates of their team members.
type TL struct{}

// Run ...
func (TL) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}

	for _, m := range team.OnlineMembers(tm) {
		user.Messagef(m, "command.tl", p.Name(), int(p.Position().X()), int(p.Position().Y()), int(p.Position().Z()))
	}
}

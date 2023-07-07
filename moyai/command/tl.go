package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// TL is a command that allows players to see the coordinates of their team members.
type TL struct{}

// Run ...
func (TL) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name(), p.Handler().(*user.Handler).XUID())
	if err != nil {
		o.Error(lang.Translate(locale(s), "user.data.load.error"))
		return
	}
	tm, ok := u.Team()
	if !ok {
		p.Message(text.Colourf("<red>%s</red>", "You are not in a team."))
		return
	}
	for _, t := range tm.Members {
		if uTarget, ok := user.Lookup(t.Name); ok {
			uTarget.Player().Message(text.Colourf("<green>%s</green><grey>:</grey> <yellow>%d<grey>,</grey> %d<grey>,</grey> %d</yellow>", p.Name(), int(p.Position().X()), int(p.Position().Y()), int(p.Position().Z())))
		}
	}
}

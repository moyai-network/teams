package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

// Reply is a command that allows a player to reply to their most recent private message.
type Reply struct {
	Message cmd.Varargs `cmd:"message"`
}

// Run ...
func (r Reply) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	h, ok := user.Lookup(p.Name())
	if !ok {
		// The user somehow left in the middle of this, so just stop in our tracks.
		return
	}
	u, err := data.LoadUser(p.Name(), h.XUID())
	if err != nil {
		return
	}
	/*if !u.Settings().Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "user.whisper.disabled"))
		return
	}*/
	msg := strings.TrimSpace(string(r.Message))
	if len(msg) <= 0 {
		o.Error(lang.Translatef(l, "message.empty"))
		return
	}

	target, ok := user.Lookup(u.LastMessageFrom)
	if !ok {
		o.Error(lang.Translatef(l, "command.reply.none"))
		return
	}
	t, err := data.LoadUser(u.LastMessageFrom, target.XUID())
	if err != nil {
		return
	}
	/*if !t.Settings().Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "target.whisper.disabled"))
		return
	}*/

	uColour, uMsg := text.Colourf("<white>%s</white>", u.DisplayName), text.Colourf("<white>%s</white>", msg)
	tColour, tMsg := text.Colourf("<white>%s</white>", t.DisplayName), text.Colourf("<white>%s</white>", msg)
	if _, ok := u.Roles.Highest().(role.Default); !ok {
		uMsg = t.Roles.Highest().Colour(msg)
		uColour = u.Roles.Highest().Colour(u.DisplayName)
	}
	if _, ok := t.Roles.Highest().(role.Default); !ok {
		tMsg = u.Roles.Highest().Colour(msg)
		tColour = t.Roles.Highest().Colour(t.DisplayName)
	}

	t.LastMessageFrom = u.Name
	//target.SendCustomSound("random.orb", 1, 1, false)
	h.Message("command.whisper.to", tColour, tMsg)
	target.Message("command.whisper.from", uColour, uMsg)
}

// Allow ...
func (Reply) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

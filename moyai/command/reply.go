package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
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
	u, err := data.LoadUserFromName(p.Name())
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
	t, err := data.LoadUserFromName(u.LastMessageFrom)
	if err != nil {
		return
	}
	/*if !t.Settings().Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "target.whisper.disabled"))
		return
	}*/

	uMsg := t.Roles.Highest().Coloured(msg)
	uColour := u.Roles.Highest().Coloured(u.DisplayName)
	tMsg := u.Roles.Highest().Coloured(msg)
	tColour := t.Roles.Highest().Coloured(t.DisplayName)

	t.LastMessageFrom = u.Name
	data.SaveUser(t)

	target.PlaySound(sound.Experience{})
	h.Message("command.whisper.to", tColour, tMsg)
	target.Message("command.whisper.from", uColour, uMsg)
}

// Allow ...
func (Reply) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

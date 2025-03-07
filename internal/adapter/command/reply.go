package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/user"
	"github.com/moyai-network/teams/pkg/lang"
	"strings"
)

// Reply is a command that allows a player to reply to their most recent private message.
type Reply struct {
	Message cmd.Varargs `cmd:"message"`
}

// Run ...
func (r Reply) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
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

	target, ok := user.Lookup(tx, u.LastMessageFrom)
	if !ok {
		o.Error(lang.Translatef(l, "command.reply.none"))
		return
	}
	t, ok := core.UserRepository.FindByName(u.LastMessageFrom)
	if !ok {
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
	core.UserRepository.Save(t)

	target.PlaySound(sound.Experience{})
	internal.Messagef(p, "command.whisper.to", tColour, tMsg)
	internal.Messagef(target, "command.whisper.from", uColour, uMsg)
}

// Allow ...
func (Reply) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

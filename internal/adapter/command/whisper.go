package command

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/pkg/lang"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
)

// Whisper is a command that allows a player to send a private message to another player.
type Whisper struct {
	Target  []cmd.Target `cmd:"target"`
	Message cmd.Varargs  `cmd:"message"`
}

// Run ...
func (w Whisper) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
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
	msg := strings.TrimSpace(string(w.Message))
	if len(msg) <= 0 {
		internal.Messagef(p, "message.empty")
		return
	}
	if len(w.Target) > 1 {
		internal.Messagef(p, "command.targets.exceed")
		return
	}

	tP, ok := w.Target[0].(*player.Player)
	if !ok {
		internal.Messagef(p, "command.target.unknown")
		return
	}
	t, ok := core.UserRepository.FindByName(tP.Name())
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	/*if !t.Settings().Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "target.whisper.disabled"))
		return
	}*/

	uMsg := t.Roles.Highest().Coloured(msg)
	uTag := u.Roles.Highest().Coloured(u.DisplayName)
	tMsg := u.Roles.Highest().Coloured(msg)
	tTag := t.Roles.Highest().Coloured(t.DisplayName)

	t.LastMessageFrom = u.Name
	core.UserRepository.Save(t)

	tP.PlaySound(sound.Experience{})
	internal.Messagef(p, "command.whisper.to", tTag, tMsg)
	internal.Messagef(tP, "command.whisper.from", uTag, uMsg)
}

// Allow ...
func (Whisper) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

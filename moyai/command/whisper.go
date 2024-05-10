package command

import (
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Whisper is a command that allows a player to send a private message to another player.
type Whisper struct {
	Target  []cmd.Target `cmd:"target"`
	Message cmd.Varargs  `cmd:"message"`
}

// Run ...
func (w Whisper) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
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
	msg := strings.TrimSpace(string(w.Message))
	if len(msg) <= 0 {
		user.Messagef(p, "message.empty")
		return
	}
	if len(w.Target) > 1 {
		user.Messagef(p, "command.targets.exceed")
		return
	}

	tP, ok := w.Target[0].(*player.Player)
	if !ok {
		user.Messagef(p, "command.target.unknown")
		return
	}
	t, err := data.LoadUserFromName(tP.Name())
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	/*if !t.Settings().Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "target.whisper.disabled"))
		return
	}*/

	uTag, uMsg := text.Colourf("<white>%s</white>", u.DisplayName), text.Colourf("<white>%s</white>", msg)
	tTag, tMsg := text.Colourf("<white>%s</white>", t.DisplayName), text.Colourf("<white>%s</white>", msg)
	if _, ok := u.Roles.Highest().(role.Default); !ok {
		uMsg = t.Roles.Highest().Color(msg)
		uTag = u.Roles.Highest().Color(u.DisplayName)
	}
	if _, ok := t.Roles.Highest().(role.Default); !ok {
		tMsg = u.Roles.Highest().Color(msg)
		tTag = t.Roles.Highest().Color(t.DisplayName)
	}

	uH, ok := user.Lookup(u.Name)
	if !ok {
		return
	}

	tH, ok := user.Lookup(t.Name)
	if !ok {
		return
	}

	t.LastMessageFrom = u.Name
	data.SaveUser(t)

	tP.PlaySound(sound.Experience{})
	uH.Message("command.whisper.to", tTag, tMsg)
	tH.Message("command.whisper.from", uTag, uMsg)
}

// Allow ...
func (Whisper) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

package command

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/data"
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
	u, err := data.LoadUser(s.(*player.Player).Name(), "")
	if err != nil {
		return
	}
	/*if !u.Settings().Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "user.whisper.disabled"))
		return
	}*/
	msg := strings.TrimSpace(string(w.Message))
	if len(msg) <= 0 {
		o.Error(lang.Translatef(l, "message.empty"))
		return
	}
	if len(w.Target) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}

	tP, ok := w.Target[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	t, err := data.LoadUser(tP.Name(), "")
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	/*if !t.Settings().Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "target.whisper.disabled"))
		return
	}*/

	uTag, uMsg := text.Colourf("<white>%s</white>", u.Name), text.Colourf("<white>%s</white>", msg)
	tTag, tMsg := text.Colourf("<white>%s</white>", t.Name), text.Colourf("<white>%s</white>", msg)
	if _, ok := u.Roles.Highest().(role.Default); !ok {
		uMsg = t.Roles.Highest().Colour(msg)
		uTag = u.Roles.Highest().Colour(u.Name)
	}
	if _, ok := t.Roles.Highest().(role.Default); !ok {
		tMsg = u.Roles.Highest().Colour(msg)
		tTag = t.Roles.Highest().Colour(t.Name)
	}

	uH, ok := user.Lookup(u.Name)
	if !ok {
		return
	}

	tH, ok := user.Lookup(t.Name)
	if !ok {
		return
	}

	tH.SetLastMessageFrom(uH.Player())
	tH.Player().PlaySound(sound.Experience{})
	uH.Player().Message("command.whisper.to", tTag, tMsg)
	tH.Player().Message("command.whisper.from", uTag, uMsg)
}

// Allow ...
func (Whisper) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

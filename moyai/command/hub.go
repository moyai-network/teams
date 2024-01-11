package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Hub is a command that teleports the player to the hub.
type Hub struct {
}

// Run ...
func (h Hub) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	u, ok := user.Lookup(p.Name())
	if !ok {
		return
	}

	if u.Combat().Active() {
		o.Error(lang.Translate(p.Locale(), "command.error.combat-tagged"))
	}

	o.Print(text.Colourf("<green>Travelling to <black>The</black> <gold>Hub</gold>...</green>"))
	_ = p.Transfer("127.0.0.1:19132")
}

// Allow ...
func (Hub) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

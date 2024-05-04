package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Hub is a command that teleports the player to the hub.
type Hub struct {
}

// Run ...
func (Hub) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	if h.Combat().Active() {
		user.Messagef(p, "command.error.combat-tagged")
	}

	o.Print(text.Colourf("<green>Travelling to <black>The</black> <gold>Hub</gold>...</green>"))

	// jk this doesn't work
	p.RemoveScoreboard()
	_ = p.Transfer("127.0.0.1:19133")
}

// Allow ...
func (Hub) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

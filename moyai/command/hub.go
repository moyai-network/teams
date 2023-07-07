package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/paroxity/portal/socket"
	proxypacket "github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Hub is a command that teleports the player to the hub.
type Hub struct {
	s *socket.Client
}

func NewHub(s *socket.Client) *Hub {
	return &Hub{s: s}
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
	h.s.WritePacket(&proxypacket.TransferRequest{
		PlayerUUID: p.UUID(),
		Server:     "syn.lobby",
	})
}

// Allow ...
func (Hub) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

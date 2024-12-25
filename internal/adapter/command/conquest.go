package command

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/conquest"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// ConquestStart is a command that starts a KOTH.
type ConquestStart struct {
	adminAllower
	Sub cmd.SubCommand `cmd:"start"`
}

// ConquestStop is a command that stops a KOTH.
type ConquestStop struct {
	adminAllower
	Sub cmd.SubCommand `cmd:"stop"`
}

// Run ...
func (k ConquestStart) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	p, ok := s.(*player.Player)
	if ok {
		if u, ok := core.UserRepository.FindByName(p.Name()); ok {
			r := u.Roles.Highest()
			name = r.Coloured(p.Name())
		}
	}
	if conquest.Running() {
		internal.Messagef(p, "command.koth.running")
		return
	}
	conquest.Start()

	areas := conquest.All()
	for u := range internal.Players(p.Tx()) {
		for _, a := range areas {
			if a.Area().Vec3WithinOrEqualXZ(u.Position()) {
				a.StartCapturing(u)
			}
		}
	}

	internal.Broadcastf(p.Tx(), "koth.start", name, "Conquest", -100.0, -500.0)
	st := `
 §e█████████§r
 §e█████████§r
 §e█§6█§e█§6█§e█§6█§e█§6█§e█§r
 §e█§6███████§e█§r
 §e█§6█§b█§6█§b█§6█§b█§6█§e█§r §eConquest§r
 §e█§6███████§e█§r §6can be contested now!§r
 §e█████████§r
 §e█████████§r
 §e█████████§r

`

	p.Message(st)
}

// Run ...
func (ConquestStop) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	p, ok := s.(*player.Player)
	if ok {
		if u, ok := core.UserRepository.FindByName(p.Name()); ok {
			r := u.Roles.Highest()
			name = r.Coloured(p.Name())
		}
	}
	if !conquest.Running() {
		internal.Messagef(p, "command.koth.not.running")
	} else {
		conquest.Stop()
		internal.Broadcastf(p.Tx(), "koth.stop", name, "Conquest")
	}
}

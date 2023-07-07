package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
)

// Fly is a command that allows the player to fly in spawn.
type Fly struct{}

// Run ...
func (Fly) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	if f, ok := p.GameMode().(flyGameMode); ok {
		o.Print(lang.Translatef(p.Locale(), "command.fly.disabled"))
		p.SetGameMode(f.GameMode)
		return
	}
	o.Print(lang.Translatef(p.Locale(), "command.fly.enabled"))
	p.SetGameMode(flyGameMode{GameMode: p.GameMode()})
}

// Allow ...
func (Fly) Allow(s cmd.Source) bool {
	return allow(s, false, role.Admin{})
}

// flyGameMode is a game mode that allows the player to fly.
type flyGameMode struct {
	world.GameMode
}

func (flyGameMode) AllowsFlying() bool { return true }

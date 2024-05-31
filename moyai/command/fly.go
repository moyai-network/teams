package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/role"
)

// Fly is a command that allows the player to fly in spawn.
type Fly struct{}

// Run ...
func (Fly) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	if f, ok := p.GameMode().(flyGameMode); ok {
		moyai.Messagef(p, "command.fly.disabled")
		p.SetGameMode(f.GameMode)
		return
	}
	moyai.Messagef(p, "command.fly.enabled")
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

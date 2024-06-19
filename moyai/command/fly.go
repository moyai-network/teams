package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai"
)

// Fly is a command that allows the player to fly in spawn.
type Fly struct{ adminAllower }

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

// flyGameMode is a game mode that allows the player to fly.
type flyGameMode struct {
	world.GameMode
}

func (flyGameMode) AllowsFlying() bool { return true }

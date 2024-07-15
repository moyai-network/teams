package command

import (
	"strings"

	"github.com/moyai-network/teams/moyai"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// GameMode is a command for a player to change their own game mode or another player's.
type GameMode struct {
	managerAllower
	GameMode gameMode                   `cmd:"gamemode"`
	Targets  cmd.Optional[[]cmd.Target] `cmd:"target"`
}

// Run ...
func (g GameMode) Run(src cmd.Source, _ *cmd.Output) {
	var name string
	var mode world.GameMode
	switch strings.ToLower(string(g.GameMode)) {
	case "survival", "0", "s":
		name, mode = "survival", world.GameModeSurvival
	case "creative", "1", "c":
		name, mode = "creative", world.GameModeCreative
	case "adventure", "2", "a":
		name, mode = "adventure", world.GameModeAdventure
	case "spectator", "3", "sp":
		name, mode = "spectator", world.GameModeSpectator
	}

	targets := g.Targets.LoadOr(nil)
	if len(targets) > 1 {
		moyai.Messagef(src, "command.targets.exceed")
		return
	}
	if len(targets) == 1 {
		target, ok := targets[0].(*player.Player)
		if !ok {
			moyai.Messagef(src, "command.target.unknown")
			return
		}

		target.SetGameMode(mode)
		moyai.Alertf(src, "staff.alert.gamemode.change.other", target.Name(), name)
		moyai.Messagef(src, "command.gamemode.update.other", target.Name(), name)
		return
	}
	if p, ok := src.(*player.Player); ok {
		p.SetGameMode(mode)
		moyai.Alertf(src, "staff.alert.gamemode.change", name)
		moyai.Messagef(src, "command.gamemode.update.self", name)
		return
	}
	moyai.Messagef(src, "command.gamemode.console")
}

type gameMode string

// Type ...
func (gameMode) Type() string {
	return "GameMode"
}

// Options ...
func (gameMode) Options(cmd.Source) []string {
	return []string{
		"survival", "0", "s",
		"creative", "1", "c",
		"adventure", "2", "a",
		"spectator", "3", "sp",
	}
}

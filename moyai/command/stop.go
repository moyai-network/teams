package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"syscall"
)

type Stop struct{ operatorAllower }

func (s Stop) Run(src cmd.Source, out *cmd.Output) {
	if _, ok := src.(*player.Player); ok {
		out.Error("This command can only be run from the console.")
		return
	}

	moyai.Close()
	syscall.Exit(0)
}

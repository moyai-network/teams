package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/menu"
)

type Prizes struct{}

func (Prizes) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	menu.SendPrizesMenu(p)
}

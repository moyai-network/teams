package command

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/menu"
)

type Settings struct{}

func (Settings) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	inv.SendMenu(p, menu.NewSettings())
}

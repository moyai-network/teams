package command

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/menu"
)

// Kit is a command that allows players to select a kit.
type Kit struct{}

// Run ...
func (Kit) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	inv.SendMenu(p, menu.NewKitsMenu())
	//p.SendForm(form.NewKitForm(p))
}

// Allow ...
func (Kit) Allow(src cmd.Source) bool {
	return allow(src, false)
}

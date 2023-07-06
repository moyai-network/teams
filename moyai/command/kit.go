package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/form"
	"github.com/moyai-network/teams/moyai/user"
)

// Kit is a command that allows players to select a kit.
type Kit struct{}

// Run ...
func (Kit) Run(src cmd.Source, out *cmd.Output) {
	p := src.(*player.Player)
	u, ok := user.Lookup(p.Name())
	if !ok {
		return
	}
	p.SendForm(form.NewKitForm(u.Player()))
}

// Allow ...
func (Kit) Allow(src cmd.Source) bool {
	return allow(src, false)
}

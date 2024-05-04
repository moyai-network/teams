package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

// Kit is a command that allows players to select a kit.
type Kit struct{}

// Run ...
func (Kit) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	_ = p
	panic("inv menu")
	//p.SendForm(form.NewKitForm(p))
}

// Allow ...
func (Kit) Allow(src cmd.Source) bool {
	return allow(src, false)
}

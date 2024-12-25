package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/user"
)

// Nick is a command that allows the player to change their nickname.
type Nick struct {
	adminAllower
	Name string `cmd:"name"`
}

// NickReset is a command that allows the player to reset their nickname.
type NickReset struct {
	adminAllower
	Sub cmd.SubCommand `cmd:"reset"`
}

// Run ...
func (n Nick) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	u.DisplayName = n.Name
	core.UserRepository.Save(u)
	user.UpdateState(p)
	internal.Messagef(p, "nick.changed", n.Name)
}

// Run ...
func (n NickReset) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	u.DisplayName = p.Name()
	core.UserRepository.Save(u)
	internal.Messagef(p, "nick.reset")
}

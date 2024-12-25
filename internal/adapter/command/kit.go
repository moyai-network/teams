package command

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/menu"
)

// Kit is a command that allows players to select a kit.
type Kit struct{}

// KitReset is a command to reset a players kit.
type KitReset struct {
	adminAllower
	Sub    cmd.SubCommand             `cmd:"reset"`
	Target cmd.Optional[[]cmd.Target] `cmd:"target"`
}

// Run ...
func (Kit) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	if men, ok := menu.NewKitsMenu(p); ok {
		inv.SendMenu(p, men)
	}
}

// Run ...
func (k KitReset) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	t, ok := k.Target.Load()
	if !ok {
		u, ok := core.UserRepository.FindByName(p.Name())
		if !ok {
			return
		}
		for _, kt := range u.Teams.Kits {
			kt.Reset()
		}
		core.UserRepository.Save(u)
		return
	}
	tg, ok := t[0].(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(tg.Name())
	if !ok {
		return
	}
	for _, kt := range u.Teams.Kits {
		kt.Reset()
	}
	core.UserRepository.Save(u)
}

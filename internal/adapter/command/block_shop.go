package command

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/area"
	"github.com/moyai-network/teams/internal/core/menu"
)

// BlockShop is a command that allows players to use blockshop.
type BlockShop struct{}

// Run ...
func (BlockShop) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	if !area.Spawn(tx.World()).Vec3WithinOrEqualXZ(p.Position()) {
		internal.Messagef(p, "in.spawn")
		return
	}

	inv.SendMenu(p, menu.NewBlocksMenu(p))

}

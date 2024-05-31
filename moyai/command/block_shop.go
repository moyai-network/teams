package command

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/menu"
)

// BlockShop is a command that allows players to use blockshop.
type BlockShop struct{}

// Run ...
func (BlockShop) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	if !area.Spawn(p.World()).Vec3WithinOrEqualXZ(p.Position()) {
		moyai.Messagef(p, "in.spawn")
		return
	}

	inv.SendMenu(p, menu.NewBlocksMenu(p))

}

// Allow ...
func (BlockShop) Allow(src cmd.Source) bool {
	return allow(src, false)
}

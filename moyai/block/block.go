package block

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/world"
)

func init() {
	for _, t := range block.OreTypes() {
		it := RedstoneOre{Type: t}
		creative.RegisterItem(item.NewStack(it, 1))
		world.RegisterBlock(it)
		world.RegisterItem(it)
	}
}

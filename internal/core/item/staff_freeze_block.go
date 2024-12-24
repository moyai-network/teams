package item

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type StaffFreezeBlockType struct{}

func (StaffFreezeBlockType) Name() string {
	return text.Colourf("<yellow>Freeze Block</yellow>")
}

func (StaffFreezeBlockType) Item() world.Item {
	return block.PackedIce{}
}

func (StaffFreezeBlockType) Lore() []string {
	return []string{text.Colourf("<yellow>Hit a player with this block to freeze or unfreeze them.</yellow>")}
}

func (StaffFreezeBlockType) Key() string {
	return "staff_freeze"
}

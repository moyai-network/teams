package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type TimeWarpType struct{}

func (TimeWarpType) Name() string {
	return text.Colourf("<blue>Time Warp</blue>")
}

func (TimeWarpType) Item() world.Item {
	return item.Feather{}
}

func (TimeWarpType) Lore() []string {
	return []string{text.Colourf("<grey>Teleport back to where you last threw your pearl.</grey>")}
}

func (TimeWarpType) Key() string {
	return "time_warp"
}

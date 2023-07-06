package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type NinjaStarType struct{}

func (NinjaStarType) Name() string {
	return text.Colourf("<blue>Ninja Star</blue>")
}

func (NinjaStarType) Item() world.Item {
	return item.NetherStar{}
}

func (NinjaStarType) Lore() []string {
	return []string{text.Colourf("<grey>Right click to teleport to the last player that it you.</grey>")}
}

func (NinjaStarType) Key() string {
	return "ninja_star"
}

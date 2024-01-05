package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type PearlDisablerType struct{}

func (PearlDisablerType) Name() string {
	return text.Colourf("<purple>xNatsuri's Pearl Disabler</purple>")
}

func (PearlDisablerType) Item() world.Item {
	return item.BlazeRod{}
}

func (PearlDisablerType) Lore() []string {
	return []string{text.Colourf("<grey>Hit an enemy to disable their next pearl</grey>")}
}

func (PearlDisablerType) Key() string {
	return "pearl_disabler"
}

package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type EffectDisablerType struct{}

func (EffectDisablerType) Name() string {
	return text.Colourf("<gold>Effect Disabler</gold>")
}

func (EffectDisablerType) Item() world.Item {
	return item.Slimeball{}
}

func (EffectDisablerType) Lore() []string {
	return []string{text.Colourf("<grey>Hit an opponent to disable their effects for 10 seconds</grey>")}
}

func (EffectDisablerType) Key() string {
	return "effect_disabler"
}

package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type ComboAbilityType struct{}

func (ComboAbilityType) Name() string {
	return text.Colourf("<green>Combo Ability</green>")
}

func (ComboAbilityType) Item() world.Item {
	return item.Pufferfish{}
}

func (ComboAbilityType) Lore() []string {
	return []string{text.Colourf("<grey>Begin a 10 second period where every hit on an opponent gives 1 second of Strength II</grey>")}
}

func (ComboAbilityType) Key() string {
	return "combo_ability"
}

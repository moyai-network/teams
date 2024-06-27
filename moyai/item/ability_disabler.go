package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type AbilityDisablerType struct{}

func (AbilityDisablerType) Name() string {
	return text.Colourf("<purple>Effect Disabler</purple>")
}

func (AbilityDisablerType) Item() world.Item {
	return item.Coal{}
}

func (AbilityDisablerType) Lore() []string {
	return []string{text.Colourf("<grey>Right click to disable all partner item usage by opponents for 10 seconds for opponents in a 15 block radius.</grey>")}
}

func (AbilityDisablerType) Key() string {
	return "ability_disabler"
}

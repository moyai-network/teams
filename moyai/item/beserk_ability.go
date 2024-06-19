package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type BeserkAbilityType struct{}

func (BeserkAbilityType) Name() string {
	return text.Colourf("<red>Beserk Ability</red>")
}

func (BeserkAbilityType) Item() world.Item {
	return item.Dye{Colour: item.ColourRed()}
}

func (BeserkAbilityType) Lore() []string {
	return []string{text.Colourf("<grey>Right click to receive 12 seconds of Strength II, Resistance III, and Regen III</grey>")}
}

func (BeserkAbilityType) Key() string {
	return "beserk_ability"
}

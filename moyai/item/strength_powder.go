package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type StrengthPowderType struct{}

func (StrengthPowderType) Name() string {
	return text.Colourf("<red>Strength Powder</red>")
}

func (StrengthPowderType) Item() world.Item {
	return item.BlazePowder{}
}

func (StrengthPowderType) Lore() []string {
	return []string{text.Colourf("<grey>Right click to receive Strength II for 7 seconds.</grey>")}
}

func (StrengthPowderType) Key() string {
	return "strength_powder"
}

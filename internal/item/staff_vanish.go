package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type StaffVanishType struct{}

func (StaffVanishType) Name() string {
	return text.Colourf("<yellow>Vanish</yellow>")
}

func (StaffVanishType) Item() world.Item {
	return item.Dye{Colour: item.ColourGrey()}
}

func (StaffVanishType) Lore() []string {
	return []string{text.Colourf("<yellow>Vanish yourself.</yellow>")}
}

func (StaffVanishType) Key() string {
	return "staff_vanish"
}

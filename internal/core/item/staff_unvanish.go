package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type StaffUnVanishType struct{}

func (StaffUnVanishType) Name() string {
	return text.Colourf("<yellow>Unvanish</yellow>")
}

func (StaffUnVanishType) Item() world.Item {
	return item.Dye{Colour: item.ColourGreen()}
}

func (StaffUnVanishType) Lore() []string {
	return []string{text.Colourf("<yellow>Unvanish yourself.</yellow>")}
}

func (StaffUnVanishType) Key() string {
	return "staff_unvanish"
}

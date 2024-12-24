package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type TankIngotType struct{}

func (TankIngotType) Name() string {
	return text.Colourf("<aqua>Tank Ingot</aqua>")
}

func (TankIngotType) Item() world.Item {
	return item.IronIngot{}
}

func (TankIngotType) Lore() []string {
	return []string{text.Colourf("<grey>Right click to receive Resistance III for 7 seconds.</grey>")}
}

func (TankIngotType) Key() string {
	return "tank_ingot"
}

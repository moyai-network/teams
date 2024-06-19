package item

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type StaffRandomTeleportType struct{}

func (StaffRandomTeleportType) Name() string {
	return text.Colourf("<yellow>Random TP</yellow>")
}

func (StaffRandomTeleportType) Item() world.Item {
	return EyeOfEnder{}
}

func (StaffRandomTeleportType) Lore() []string {
	return []string{text.Colourf("<yellow>Right click to teleport to a random player.</yellow>")}
}

func (StaffRandomTeleportType) Key() string {
	return "staff_random_teleport"
}

package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type StaffTeleportStickType struct{}

func (StaffTeleportStickType) Name() string {
	return text.Colourf("<yellow>TP Stick</yellow>")
}

func (StaffTeleportStickType) Item() world.Item {
	return item.BlazeRod{}
}

func (StaffTeleportStickType) Lore() []string {
	return []string{text.Colourf("<yellow>Teleport to a specific player.</yellow>")}
}

func (StaffTeleportStickType) Key() string {
	return "staff_teleport_stick"
}

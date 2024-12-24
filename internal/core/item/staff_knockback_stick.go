package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type StaffKnockBackStickType struct{}

func (StaffKnockBackStickType) Name() string {
	return text.Colourf("<yellow>Knockback Stick</yellow>")
}

func (StaffKnockBackStickType) Item() world.Item {
	return item.Stick{}
}

func (StaffKnockBackStickType) Lore() []string {
	return []string{text.Colourf("<yellow>Hit a player with this stick to test their knockback.</yellow>")}
}

func (StaffKnockBackStickType) Key() string {
	return "staff_knockback_stick"
}

package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type FullInvisibilityType struct{}

func (FullInvisibilityType) Name() string {
	return text.Colourf("<black>Full Invisibility</black>")
}

func (FullInvisibilityType) Item() world.Item {
	return item.InkSac{}
}

func (FullInvisibilityType) Lore() []string {
	return []string{text.Colourf("<grey>Become completely invisible to enemies for 60 seconds or until they hit you.</grey>")}
}

func (FullInvisibilityType) Key() string {
	return "full_invisibility"
}

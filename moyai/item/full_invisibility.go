package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type FullInvisibilityType struct{}

func (FullInvisibilityType) Name() string {
	return text.Colourf("<black>xCqzzz's Full Invisibility</black>")
}

func (FullInvisibilityType) Item() world.Item {
	return item.InkSac{}
}

func (FullInvisibilityType) Lore() []string {
	return []string{text.Colourf("<grey>Become completely invisible to enemies.</grey>")}
}

func (FullInvisibilityType) Key() string {
	return "full_invisibility"
}

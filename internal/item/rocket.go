package item

import (
	"time"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type RocketType struct{}

func (RocketType) Name() string {
	return text.Colourf("<red>Rocket</red>")
}

func (RocketType) Item() world.Item {
	return item.Firework{
		Duration: time.Second,
	}
}

func (RocketType) Lore() []string {
	return []string{text.Colourf("<grey>Right click to boost up in the sky!</grey>")}
}

func (RocketType) Key() string {
	return "rocket"
}

package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type RageBrickType struct{}

func (RageBrickType) Name() string {
	return text.Colourf("<red>Rage Brick</red>")
}

func (RageBrickType) Item() world.Item {
	return item.NetherBrick{}
}

func (RageBrickType) Lore() []string {
	return []string{text.Colourf("<grey>Right click to receive 1 second Strength II for every player in a 15 block radius.</grey>")}
}

func (RageBrickType) Key() string {
	return "rage_brick"
}

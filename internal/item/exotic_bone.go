package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type ExoticBoneType struct{}

func (ExoticBoneType) Name() string {
	return text.Colourf("<red>RestartFU's Exotic Bone</red>")
}

func (ExoticBoneType) Item() world.Item {
	return item.Bone{}
}

func (ExoticBoneType) Lore() []string {
	return []string{text.Colourf("<grey>Hit a player to prevent them from\nbreaking, placing or interacting with blocks for 15 seconds</grey>")}
}

func (ExoticBoneType) Key() string {
	return "exotic_bone"
}

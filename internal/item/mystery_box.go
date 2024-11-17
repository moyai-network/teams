package item

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func NewMysteryBox(n int) item.Stack {
	return item.NewStack(block.Chest{}, n).WithValue("MYSTERY_BOX", true).WithCustomName(text.Colourf("§r<redstone>Mystery Box</redstone>"))
}

package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func NewCrowBar() item.Stack {
	return item.NewStack(item.Hoe{Tier: item.ToolTierDiamond}, 1).WithValue("CROWBAR", true).WithCustomName(text.Colourf("§r<yellow>Crowbar</yellow>"))
}

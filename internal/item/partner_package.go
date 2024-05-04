package item

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func NewPartnerPackage(n int) item.Stack {
	return item.NewStack(block.EnderChest{}, n).WithValue("PARTNER_PACKAGE", true).WithCustomName(text.Colourf("§r<amethyst>Partner Package</amethyst>"))
}

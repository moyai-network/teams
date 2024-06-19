package item

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type PartnerPackageType struct{}

func (PartnerPackageType) Name() string {
	return text.Colourf("<amethyst>Partner Package</amethyst>")
}

func (PartnerPackageType) Item() world.Item {
	return block.EnderChest{}
}

func (PartnerPackageType) Lore() []string {
	return []string{text.Colourf("<grey>Right click to open a partner package.</grey>")}
}

func (PartnerPackageType) Key() string {
	return "partner_package"
}

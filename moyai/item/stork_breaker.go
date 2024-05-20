package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type StormBreakerType struct{}

func (StormBreakerType) Name() string {
	return text.Colourf("<yellow>Storm Breaker</yellow>")
}

func (StormBreakerType) Item() world.Item {
	return item.Axe{Tier: item.ToolTierGold}
}

func (StormBreakerType) Lore() []string {
	return []string{text.Colourf("<grey>TODO</grey>")}
}

func (StormBreakerType) Key() string {
	return "storm_breaker"
}

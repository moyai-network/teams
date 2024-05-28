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
	return []string{text.Colourf("<grey>Hit a player to summon a lightning bolt on them and replace their helmet by a leather one for 5 seconds.</grey>")}
}

func (StormBreakerType) Key() string {
	return "storm_breaker"
}

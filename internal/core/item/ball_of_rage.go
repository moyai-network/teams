package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type BallOfRageType struct{}

func (BallOfRageType) Name() string {
	return text.Colourf("<red>Ball of Rage</red>")
}

func (BallOfRageType) Item() world.Item {
	return item.Egg{}
}

func (BallOfRageType) Lore() []string {
	return []string{text.Colourf("<grey>Throw this at opponent to create a cloud of Strength II and Resistance III to players in your faction and Weakness II and Wither II to opponents for 7 seconds.</grey>")}
}

func (BallOfRageType) Key() string {
	return "ball_of_rage"
}

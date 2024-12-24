package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type VampireAbilityType struct{}

func (VampireAbilityType) Name() string {
	return text.Colourf("<red>Vampire Ability</red>")
}

func (VampireAbilityType) Item() world.Item {
	return item.RabbitFoot{}
}

func (VampireAbilityType) Lore() []string {
	return []string{text.Colourf("<grey>Begin a 10 second period where you heal 50%% of the damage dealt to an opponent.</grey>")}
}

func (VampireAbilityType) Key() string {
	return "vampire_ability"
}

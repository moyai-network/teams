package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type FocusModeType struct{}

func (FocusModeType) Name() string {
	return text.Colourf("<gold>soraalex's Focus Mode</gold>")
}

func (FocusModeType) Item() world.Item {
	return item.GoldNugget{}
}

func (FocusModeType) Lore() []string {
	return []string{text.Colourf("<grey>Right click to deal 25%% more damage to non-class and non-archer-tagged opponents for 10 seconds.</grey>")}
}

func (FocusModeType) Key() string {
	return "focus_mode"
}

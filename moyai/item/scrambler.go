package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type ScramblerType struct{}

func (ScramblerType) Name() string {
	return text.Colourf("<gold>Scrambler</gold>")
}

func (ScramblerType) Item() world.Item {
	return item.Stick{}
}

func (ScramblerType) Lore() []string {
	return []string{text.Colourf("<grey>Hit a player 3 times to scramble their hotbar</grey>")}
}

func (ScramblerType) Key() string {
	return "scrambler"
}

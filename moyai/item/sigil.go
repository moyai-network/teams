package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type SigilType struct{}

func (SigilType) Name() string {
	return text.Colourf("<aqua>Allah's Sigil</aqua>")
}

func (SigilType) Item() world.Item {
	return item.Clock{}
}

func (SigilType) Lore() []string {
	return []string{text.Colourf("<grey>Call upon the All-Powerful Allah.</grey>")}
}

func (SigilType) Key() string {
	return "sigil"
}

package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type CloseCallType struct{}

func (CloseCallType) Name() string {
	return text.Colourf("<gold>AUTOClICK's Close Call</gold>")
}

func (CloseCallType) Item() world.Item {
	return item.Cookie{}
}

func (CloseCallType) Lore() []string {
	return []string{text.Colourf("<grey>Use while under 3 hearts to receive Regeneration IV and Strength II for 5 seconds.</grey>")}
}

func (CloseCallType) Key() string {
	return "close_call"
}

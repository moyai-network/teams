package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type KeyType struct {
	key int
}

func KeyTypeKOTH() KeyType {
	return KeyType{key: 0}
}

func KeyTypePharaoh() KeyType {
	return KeyType{key: 1}
}

func AllKeyTypes() []KeyType {
	return []KeyType{KeyTypeKOTH(), KeyTypePharaoh()}
}

func NewKey(keyType KeyType, n int) item.Stack {
	var colour item.Colour
	var value string
	var customName string

	switch keyType.key {
	case 0:
		colour = item.ColourRed()
		value = "crate-key_KOTH"
		customName = text.Colourf("<red>KOTH Crate Key</red>")
	case 1:
		colour = item.ColourCyan()
		value = "crate-key_Pharaoh"
		customName = text.Colourf("<black>Pharaoh Crate Key</black>")
	default:
		panic("should never happen")
	}
	return item.NewStack(item.Dye{Colour: colour}, n).WithValue(value, true).WithCustomName(customName)
}

func init() {
	for _, t := range AllKeyTypes() {
		creative.RegisterItem(NewKey(t, 1))
	}
}

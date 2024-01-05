package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type KeyType struct {
	key int
}

const (
	KeyTypeKOTH     = 0
	KeyTypeRevenant = 1
	KeyTypePartner  = 2
	KeyTypeNekros   = 3
	KeyTypeNova     = 4
)

func AllKeyTypes() []KeyType {
	return []KeyType{
		{key: KeyTypeKOTH},
		{key: KeyTypeRevenant},
		{key: KeyTypePartner},
		{key: KeyTypeNekros},
		{key: KeyTypeNova},
	}
}

func NewKey(keyType int, n int) item.Stack {
	var colour item.Colour
	var value string
	var customName string

	switch keyType {
	case KeyTypeKOTH:
		colour = item.ColourRed()
		value = "crate-key_KOTH"
		customName = text.Colourf("<red>KOTH Crate Key</red>")
	case KeyTypeRevenant:
		colour = item.ColourCyan()
		value = "crate-key_Revenant"
		customName = text.Colourf("<redstone>Revenant Crate Key</redstone>")
	case KeyTypePartner:
		colour = item.ColourGreen()
		value = "crate-key_Partner"
		customName = text.Colourf("<green>Partner Crate Key</green>")
	case KeyTypeNekros:
		colour = item.ColourPurple()
		value = "crate-key_Nekros"
		customName = text.Colourf("<purple>Nekros Crate Key</purple>")
	case KeyTypeNova:
		colour = item.ColourYellow()
		value = "crate-key_Nova"
		customName = text.Colourf("<yellow>Nova Crate Key</yellow>")
	default:
		panic("should never happen")
	}

	return item.NewStack(item.Dye{Colour: colour}, n).WithValue(value, true).WithCustomName(customName)
}

func init() {
	for _, t := range AllKeyTypes() {
		creative.RegisterItem(NewKey(t.key, 1))
	}
}

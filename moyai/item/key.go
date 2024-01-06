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
	KeyTypeKOTH    = 0
	KeyTypePharaoh = 1
	KeyTypePartner = 2
	KeyTypeMenes   = 3
	KeyTypeRamses  = 4
)

func AllKeyTypes() []KeyType {
	return []KeyType{
		{key: KeyTypeKOTH},
		{key: KeyTypePharaoh},
		{key: KeyTypePartner},
		{key: KeyTypeMenes},
		{key: KeyTypeRamses},
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
	case KeyTypePharaoh:
		colour = item.ColourBlack()
		value = "crate-key_Pharaoh"
		customName = text.Colourf("<black>Pharaoh Crate Key</black>")
	case KeyTypePartner:
		colour = item.ColourPurple()
		value = "crate-key_Partner"
		customName = text.Colourf("<green>Partner Crate Key</green>")
	case KeyTypeMenes:
		colour = item.ColourGreen()
		value = "crate-key_Menes"
		customName = text.Colourf("<emerald>Menes Crate Key</emerald>")
	case KeyTypeRamses:
		colour = item.ColourLightBlue()
		value = "crate-key_Ramses"
		customName = text.Colourf("<diamond>Ramses Crate Key</diamond>")
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

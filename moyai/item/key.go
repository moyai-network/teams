package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	ench "github.com/moyai-network/teams/moyai/enchantment"
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
	KeyTypeConquest = 5
)

func AllKeyTypes() []KeyType {
	return []KeyType{
		{key: KeyTypeKOTH},
		{key: KeyTypePharaoh},
		{key: KeyTypePartner},
		{key: KeyTypeMenes},
		{key: KeyTypeRamses},
		{key: KeyTypeConquest},
	}
}

func NewKey(keyType int, n int) item.Stack {
	var value string
	var customName string

	switch keyType {
	case KeyTypeKOTH:
		value = "crate-key_KOTH"
		customName = text.Colourf("<red>KOTH Crate Key</red>")
	case KeyTypePharaoh:
		value = "crate-key_Pharaoh"
		customName = text.Colourf("<black>Pharaoh Crate Key</black>")
	case KeyTypePartner:
		value = "crate-key_Partner"
		customName = text.Colourf("<green>Partner Crate Key</green>")
	case KeyTypeMenes:
		value = "crate-key_Menes"
		customName = text.Colourf("<emerald>Menes Crate Key</emerald>")
	case KeyTypeRamses:
		value = "crate-key_Ramses"
		customName = text.Colourf("<diamond>Ramses Crate Key</diamond>")
	case KeyTypeConquest:
		value = "crate-key_Conquest"
		customName = text.Colourf("<blue>Conquest Crate Key</blue>")
	default:
		panic("should never happen")
	}

	prot := item.NewEnchantment(ench.Protection{}, 1)
	return item.NewStack(TripwireHook{}, n).WithValue(value, true).WithCustomName(customName).WithEnchantments(prot)
}

func init() {
	for _, t := range AllKeyTypes() {
		creative.RegisterItem(NewKey(t.key, 1))
	}
}

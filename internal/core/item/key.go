package item

import (
    "github.com/df-mc/dragonfly/server/item"
    "github.com/df-mc/dragonfly/server/item/creative"
    ench "github.com/moyai-network/teams/internal/core/enchantment"
    "github.com/sandertv/gophertunnel/minecraft/text"
)

type KeyType struct {
    key int
}

var (
    KeyTypeKOTH     = KeyType{key: 0}
    KeyTypePharaoh  = KeyType{key: 1}
    KeyTypePartner  = KeyType{key: 2}
    KeyTypeMenes    = KeyType{key: 3}
    KeyTypeRamses   = KeyType{key: 4}
    KeyTypeConquest = KeyType{key: 5}
    KeyTypeSeasonal = KeyType{key: 6}
)

func AllKeyTypes() []KeyType {
    return []KeyType{
        KeyTypeKOTH,
        KeyTypePharaoh,
        KeyTypePartner,
        KeyTypeMenes,
        KeyTypeRamses,
        KeyTypeConquest,
        KeyTypeSeasonal,
    }
}

func NewKey(keyType KeyType, n int) item.Stack {
    var value string
    var customName string

    switch keyType {
    case KeyTypeKOTH:
        value = "crate-key_KOTH"
        customName = text.Colourf("<red>KOTH Crate Key</red>")
    case KeyTypePharaoh:
        value = "crate-key_Pharaoh"
        customName = text.Colourf("<dark-red>Pharaoh Crate Key</dark-red>")
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
    case KeyTypeSeasonal:
        value = "crate-key_Seasonal"
        customName = text.Colourf("<gold>Seaonal Crate Key</gold>")
    default:
        panic("should never happen")
    }

    prot := item.NewEnchantment(ench.Protection{}, 1)
    return item.NewStack(TripwireHook{}, n).WithValue(value, true).WithCustomName(customName).WithEnchantments(prot)
}

func init() {
    for _, t := range AllKeyTypes() {
        creative.RegisterItem(NewKey(t, 1))
    }
}

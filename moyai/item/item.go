package item

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/world"
)

type NopType struct{}

func (NopType) Name() string {
	return ""
}

func (NopType) Item() world.Item {
	return block.Air{}
}

func (NopType) Lore() []string {
	return nil
}

func (NopType) Key() string {
	return ""
}

type SpecialItemType interface {
	Name() string
	Item() world.Item
	Lore() []string
	Key() string
}

var (
	specialItems []SpecialItemType
)

func SpecialItems() []SpecialItemType {
	return specialItems
}

func NewSpecialItem(typ SpecialItemType, n int) item.Stack {
	return item.NewStack(typ.Item(), n).WithValue(typ.Key(), true).WithCustomName(typ.Name()).WithLore(typ.Lore()...)
}

func SpecialItem(i item.Stack) (SpecialItemType, bool) {
	for k, _ := range i.Values() {
		for _, v := range specialItems {
			if k == v.Key() && i.Item() == v.Item() {
				return v, true
			}
		}
	}
	return NopType{}, false
}

func RegisterSpecialItem(typ SpecialItemType) {
	creative.RegisterItem(NewSpecialItem(typ, 1))
	specialItems = append(specialItems, typ)
}

func init() {
	RegisterSpecialItem(SwitcherBallType{})
	//RegisterSpecialItem(NinjaStarType{})
	RegisterSpecialItem(ScramblerType{})
	RegisterSpecialItem(ExoticBoneType{})
	RegisterSpecialItem(PearlDisablerType{})
	RegisterSpecialItem(FullInvisibilityType{})
	RegisterSpecialItem(SigilType{})
	creative.RegisterItem(NewPartnerPackage(1))
}

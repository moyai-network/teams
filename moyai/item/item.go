package item

import (
	"math/rand"
	"reflect"
	"strings"
	"unicode"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/moyai/area"
	ench "github.com/moyai-network/teams/moyai/enchantment"
)

func AddOrDrop(p *player.Player, it item.Stack) {
	if _, err := p.Inventory().AddItem(it); err != nil {
		Drop(p, it)
	}
}

func AddArmorOrDrop(p *player.Player, it item.Stack) {
	switch it.Item().(type) {
	case item.Boots:
		if !p.Armour().Boots().Empty() {
			AddOrDrop(p, p.Armour().Boots())
			break
		}
		p.Armour().SetBoots(it)
	case item.Leggings:
		if !p.Armour().Leggings().Empty() {
			AddOrDrop(p, p.Armour().Leggings())
			break
		}
		p.Armour().SetLeggings(it)
	case item.Chestplate:
		if !p.Armour().Chestplate().Empty() {
			AddOrDrop(p, p.Armour().Chestplate())
			break
		}
		p.Armour().SetChestplate(it)
	case item.Helmet:
		if !p.Armour().Helmet().Empty() {
			AddOrDrop(p, p.Armour().Helmet())
			break
		}
		p.Armour().SetHelmet(it)
	}
}

func Drop(p *player.Player, it item.Stack) {
	w, pos := p.World(), p.Position()
	et := entity.NewItem(it, pos)
	et.SetVelocity(mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1})
	w.AddEntity(et)

	if area.Spawn(w).Vec3WithinOrEqualFloorXZ(pos) {
		for _, e := range w.Entities() {
			if p, ok := e.(*player.Player); ok {
				p.HideEntity(et)
			}
		}
	}
}

func DropFromPosition(w *world.World, pos mgl64.Vec3, it item.Stack) {
	et := entity.NewItem(it, pos)
	et.SetVelocity(mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1})
	w.AddEntity(et)
}

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
	partnerItems []SpecialItemType
)

// NewSpecialItem returns a new special item of the type passed with the amount n.
func NewSpecialItem(typ SpecialItemType, n int) item.Stack {
	unbreaking := item.NewEnchantment(ench.Protection{}, 1)
	return item.NewStack(typ.Item(), n).WithValue(typ.Key(), true).WithCustomName(typ.Name()).WithLore(typ.Lore()...).WithEnchantments(unbreaking)
}

// SpecialItem returns the special item type of the item stack passed. If the item stack does not represent a
// special item, the second return value will be false.
func SpecialItem(i item.Stack) (SpecialItemType, bool) {
	for k := range i.Values() {
		for _, v := range specialItems {
			if k == v.Key() {
				return v, true
			}
		}
	}
	return NopType{}, false
}

// PartnerItem returns the partner item type of the item stack passed. If the item stack does not represent a
// partner item, the second return value will be false.
func PartnerItem(i item.Stack) (SpecialItemType, bool) {
	for k, _ := range i.Values() {
		for _, v := range partnerItems {
			if k == v.Key() {
				return v, true
			}
		}
	}
	return NopType{}, false
}

// SpecialItems returns all special items.
func SpecialItems() []SpecialItemType {
	return specialItems
}

// PartnerItems returns all partner items.
func PartnerItems() []SpecialItemType {
	return partnerItems
}

// RegisterSpecialItem registers a special item type to be used in the game.
func RegisterSpecialItem(typ SpecialItemType) {
	creative.RegisterItem(NewSpecialItem(typ, 1))
	specialItems = append(specialItems, typ)
}

// RegisterPartnerItem registers a partner item type to be used in the game.
func RegisterPartnerItem(typ SpecialItemType) {
	creative.RegisterItem(NewSpecialItem(typ, 1))
	partnerItems = append(partnerItems, typ)
}

func init() {
	RegisterPartnerItem(SwitcherBallType{})
	RegisterPartnerItem(NinjaStarType{})
	RegisterPartnerItem(ExoticBoneType{})
	RegisterPartnerItem(PearlDisablerType{})
	RegisterPartnerItem(FullInvisibilityType{})
	RegisterPartnerItem(EffectDisablerType{})
	RegisterPartnerItem(BeserkAbilityType{})
	RegisterPartnerItem(CloseCallType{})
	//RegisterPartnerItem(SigilType{})
	RegisterPartnerItem(TimeWarpType{})
	RegisterPartnerItem(StormBreakerType{})
	RegisterPartnerItem(FocusModeType{})
	RegisterPartnerItem(RocketType{})
	//RegisterSpecialItem(ScramblerType{})
	RegisterSpecialItem(PartnerPackageType{})
	creative.RegisterItem(NewCrowBar())
	creative.RegisterItem(NewPearlBow())

	RegisterSpecialItem(StaffRandomTeleportType{})
	RegisterSpecialItem(StaffFreezeBlockType{})
	RegisterSpecialItem(StaffKnockBackStickType{})
	RegisterSpecialItem(StaffTeleportStickType{})
	RegisterSpecialItem(StaffVanishType{})
	RegisterSpecialItem(StaffUnVanishType{})
}

// DisplayName returns the name of the item.
func DisplayName(i world.Item) string {
	var s strings.Builder

	if it, ok := i.(item.Sword); ok {
		s.WriteString(toolTierName(it.Tier) + " ")
	} else if it, ok := i.(item.Pickaxe); ok {
		s.WriteString(toolTierName(it.Tier) + " ")
	} else if it, ok := i.(item.Axe); ok {
		s.WriteString(toolTierName(it.Tier) + " ")
	} else if it, ok := i.(item.Shovel); ok {
		s.WriteString(toolTierName(it.Tier) + " ")
	} else if it, ok := i.(item.Hoe); ok {
		s.WriteString(toolTierName(it.Tier) + " ")
	} else if it, ok := i.(item.Helmet); ok {
		s.WriteString(armourTierName(it.Tier) + " ")
	} else if it, ok := i.(item.Chestplate); ok {
		s.WriteString(armourTierName(it.Tier) + " ")
	} else if it, ok := i.(item.Leggings); ok {
		s.WriteString(armourTierName(it.Tier) + " ")
	} else if it, ok := i.(item.Boots); ok {
		s.WriteString(armourTierName(it.Tier) + " ")
	}

	t := reflect.TypeOf(i)
	if t == nil {
		return ""
	}
	name := t.Name()

	for _, r := range name {
		if unicode.IsUpper(r) && !strings.HasPrefix(name, string(r)) {
			s.WriteRune(' ')
		}
		s.WriteRune(r)
	}
	return s.String()
}

func toolTierName(t item.ToolTier) string {
	switch t {
	case item.ToolTierDiamond:
		return "Diamond"
	case item.ToolTierGold:
		return "Golden"
	case item.ToolTierIron:
		return "Iron"
	case item.ToolTierStone:
		return "Stone"
	case item.ToolTierWood:
		return "Wooden"
	}
	return ""
}

func armourTierName(t item.ArmourTier) string {
	switch t.(type) {
	case item.ArmourTierDiamond:
		return "Diamond"
	case item.ArmourTierGold:
		return "Golden"
	case item.ArmourTierIron:
		return "Iron"
	case item.ArmourTierChain:
		return "Chainmail"
	case item.ArmourTierLeather:
		return "Leather"
	}
	return ""
}

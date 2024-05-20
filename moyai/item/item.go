package item

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/moyai/area"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	"math/rand"
	"reflect"
	"strings"
	"unicode"
)

func AddOrDrop(p *player.Player, it item.Stack) {
	if _, err := p.Inventory().AddItem(it); err != nil {
		Drop(p, it)
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
	unbreaking := item.NewEnchantment(ench.Protection{}, 1)
	return item.NewStack(typ.Item(), n).WithValue(typ.Key(), true).WithCustomName(typ.Name()).WithLore(typ.Lore()...).WithEnchantments(unbreaking)
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
	RegisterSpecialItem(NinjaStarType{})
	RegisterSpecialItem(ScramblerType{})
	RegisterSpecialItem(ExoticBoneType{})
	RegisterSpecialItem(PearlDisablerType{})
	RegisterSpecialItem(FullInvisibilityType{})
	RegisterSpecialItem(SigilType{})
	RegisterSpecialItem(TimeWarpType{})
	RegisterSpecialItem(StormBreakerType{})
	creative.RegisterItem(NewPartnerPackage(1))
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

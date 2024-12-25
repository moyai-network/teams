package class

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/ports"
	"reflect"
)

var (
	classes []ports.Class
)

func ResolveFromArmour(a [4]item.ArmourTier) ports.Class {
	for _, c := range classes {
		if compareTypes(a[0], c.Armour()[0]) && compareTypes(a[1], c.Armour()[1]) && compareTypes(a[2], c.Armour()[2]) && compareTypes(a[3], c.Armour()[3]) {
			return c
		}
	}
	return nil
}

func compareTypes(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

func Resolve(p *player.Player) ports.Class {
	a := p.Armour()
	helmet, ok := a.Helmet().Item().(item.Helmet)
	if !ok {
		return nil
	}
	chestplate, ok := a.Chestplate().Item().(item.Chestplate)
	if !ok {
		return nil
	}
	leggings, ok := a.Leggings().Item().(item.Leggings)
	if !ok {
		return nil
	}
	boots, ok := a.Boots().Item().(item.Boots)
	if !ok {
		return nil
	}

	return ResolveFromArmour([4]item.ArmourTier{helmet.Tier, chestplate.Tier, leggings.Tier, boots.Tier})
}

func Compare(a ports.Class, b ports.Class) bool {
	if a == b {
		return true
	}
	return false
}

func CompareAny(a ports.Class, b ...ports.Class) bool {
	for _, c := range b {
		if Compare(a, c) {
			return true
		}
	}
	return false
}

func Register(c ports.Class) {
	classes = append(classes, c)
}

func init() {
	Register(Archer{})
	Register(Bard{})
	Register(Miner{})
	Register(Rogue{})
	Register(Mage{})
}

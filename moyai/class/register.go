package class

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

var (
	classes []Class
)

func ResolveFromArmour(a [4]item.ArmourTier) Class {
	for _, c := range classes {
		if compareTierTypes(a[0], c.Armour()[0]) && compareTierTypes(a[1], c.Armour()[1]) && compareTierTypes(a[2], c.Armour()[2]) && compareTierTypes(a[3], c.Armour()[3]) {
			return c
		}
	}
	return nil
}

func compareTierTypes[a item.ArmourTier, b item.ArmourTier](t1 a, t2 b) bool {
	_, ok1 := any(t1).(b)
	_, ok2 := any(t2).(a)
	return ok1 && ok2
}

func Resolve(p *player.Player) Class {
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

func Compare(a Class, b Class) bool {
	if a == b {
		return true
	}
	return false
}

func CompareAny(a Class, b ...Class) bool {
	for _, c := range b {
		if Compare(a, c) {
			return true
		}
	}
	return false
}

func Register(c Class) {
	classes = append(classes, c)
}

func init() {
	Register(Archer{})
	Register(Bard{})
	Register(Miner{})
	Register(Rogue{})
	Register(Mage{})
}

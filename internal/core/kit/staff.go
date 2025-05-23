package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/core"
	item2 "github.com/moyai-network/teams/internal/core/item"
)

// Staff represents the staff mode kit.
type Staff struct{}

// Name ...
func (Staff) Name() string {
	return "Staff"
}

// Texture ...
func (Staff) Texture() string {
	return ""
}

// Items ...
func (Staff) Items(p *player.Player) [36]item.Stack {
	items := [36]item.Stack{}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return items
	}

	items[0] = item2.NewSpecialItem(item2.StaffFreezeBlockType{}, 1)
	items[2] = item2.NewSpecialItem(item2.StaffKnockBackStickType{}, 1)
	items[4] = item2.NewSpecialItem(item2.StaffRandomTeleportType{}, 1)
	items[6] = item2.NewSpecialItem(item2.StaffTeleportStickType{}, 1)

	if u.Vanished {
		items[8] = item2.NewSpecialItem(item2.StaffVanishType{}, 1)
	} else {
		items[8] = item2.NewSpecialItem(item2.StaffUnVanishType{}, 1)
	}

	return items
}

// Armour ...
func (Staff) Armour(*player.Player) [4]item.Stack {
	return [4]item.Stack{}
}

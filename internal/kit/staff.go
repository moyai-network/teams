package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/data"
	it "github.com/moyai-network/teams/internal/item"
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

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return items
	}

	items[0] = it.NewSpecialItem(it.StaffFreezeBlockType{}, 1)
	items[2] = it.NewSpecialItem(it.StaffKnockBackStickType{}, 1)
	items[4] = it.NewSpecialItem(it.StaffRandomTeleportType{}, 1)
	items[6] = it.NewSpecialItem(it.StaffTeleportStickType{}, 1)

	if u.Vanished {
		items[8] = it.NewSpecialItem(it.StaffVanishType{}, 1)
	} else {
		items[8] = it.NewSpecialItem(it.StaffUnVanishType{}, 1)
	}

	return items
}

// Armour ...
func (Staff) Armour(*player.Player) [4]item.Stack {
	return [4]item.Stack{}
}

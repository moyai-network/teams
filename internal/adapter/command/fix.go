package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	rls "github.com/moyai-network/teams/internal/core/roles"
)

// Fix is a command that allows the player to fix the item in their hand or another player's hand.
type Fix struct {
	Player cmd.Optional[[]cmd.Target] `cmd:"player"`
}

// FixAll is a command that allows the player to fix all items in their inventory or another player's hand.
type FixAll struct {
	Sub    cmd.SubCommand             `cmd:"all"`
	Player cmd.Optional[[]cmd.Target] `cmd:"player"`
}

// Run ...
func (f Fix) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p := s.(*player.Player)
	if t, ok := f.Player.Load(); ok {
		if len(t) > 1 {
			internal.Messagef(p, "command.targets.exceed")
			return
		}
		target, ok := t[0].(*player.Player)
		if !ok {
			internal.Messagef(p, "command.target.unknown")
			return
		}
		p = target
	}

	it, off := p.HeldItems()
	it = it.WithDurability(it.MaxDurability())
	p.SetHeldItems(it, off)

	internal.Messagef(p, "command.fix.success")
}

// Run ...
func (f FixAll) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p := s.(*player.Player)
	if t, ok := f.Player.Load(); ok {
		if len(t) > 1 {
			internal.Messagef(p, "command.targets.exceed")
			return
		}
		target, ok := t[0].(*player.Player)
		if !ok {
			internal.Messagef(p, "command.target.unknown")
			return
		}
		p = target
	}

	for i, it := range p.Inventory().Slots() {
		_ = p.Inventory().SetItem(i, it.WithDurability(it.MaxDurability()))
	}

	for i, it := range p.Armour().Items() {
		newIt := it.WithDurability(it.MaxDurability())
		switch i {
		case 0:
			p.Armour().SetHelmet(newIt)
		case 1:
			p.Armour().SetChestplate(newIt)
		case 2:
			p.Armour().SetLeggings(newIt)
		case 3:
			p.Armour().SetBoots(newIt)
		}
	}

	internal.Messagef(p, "command.fix.success")
}

// Allow ...
func (Fix) Allow(s cmd.Source) bool {
	return Allow(s, false, rls.Khufu())
}

// Allow ...
func (FixAll) Allow(s cmd.Source) bool {
	return Allow(s, false, rls.Ramses())
}

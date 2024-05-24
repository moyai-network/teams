package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/user"
)

// Fix is a command that allows the player to fix the item in their hand or another player's hand.
type Fix struct {
	Player cmd.Optional[[]cmd.Target] `cmd:"player"`
}

// FixAll is a command that allows the player to fix all items in their inventory or another player's hand.
type FixAll struct {
	Player cmd.Optional[[]cmd.Target] `cmd:"player"`
}

// Run ...
func (f Fix) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	if t, ok := f.Player.Load(); ok {
		if len(t) > 1 {
			user.Messagef(p, "command.targets.exceed")
			return
		}
		target, ok := t[0].(*player.Player)
		if !ok {
			user.Messagef(p, "command.target.unknown")
			return
		}
		p = target
	}

	it, _ := p.HeldItems()
	it.WithDurability(it.MaxDurability())

	user.Messagef(p, "command.fix.success")
}

// Run ...
func (f FixAll) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	if t, ok := f.Player.Load(); ok {
		if len(t) > 1 {
			user.Messagef(p, "command.targets.exceed")
			return
		}
		target, ok := t[0].(*player.Player)
		if !ok {
			user.Messagef(p, "command.target.unknown")
			return
		}
		p = target
	}

	for i, it := range p.Inventory().Items() {
		new := it.WithDurability(it.MaxDurability())
		p.Inventory().SetItem(i, new)
	}

	for i, it := range p.Armour().Items() {
		new := it.WithDurability(it.MaxDurability())
		switch i {
		case 0:
			p.Armour().SetHelmet(new)
		case 1:
			p.Armour().SetChestplate(new)
		case 2:
			p.Armour().SetLeggings(new)
		case 3:
			p.Armour().SetBoots(new)
		}
	}

	user.Messagef(p, "command.fix.success")
}

// Allow ...
func (Fix) Allow(s cmd.Source) bool {
	return allow(s, false, role.Donor1{})
}

// Allow ...
func (FixAll) Allow(s cmd.Source) bool {
	return allow(s, false, role.Donor2{})
}

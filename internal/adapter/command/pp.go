package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/item"
)

// PartnerPackage is a command that allows admins to give players partner packages
type PartnerPackage struct {
	adminAllower
	Targets []cmd.Target `cmd:"target"`
	Count   int          `cmd:"count"`
}

// PartnerPackageAll is a command that distributes partner packages to the server
type PartnerPackageAll struct {
	adminAllower
	Sub   cmd.SubCommand `cmd:"all"`
	Count int            `cmd:"count"`
}

type PartnerItemsRefresh struct {
	adminAllower
	Sub cmd.SubCommand `cmd:"refresh"`
}

// Run ...
func (pp PartnerPackage) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	t, ok := pp.Targets[0].(*player.Player)
	if !ok {
		internal.Messagef(p, "command.target.unknown")
		return
	}

	internal.Messagef(p, "command.partner_package.give.success", t.Name(), pp.Count)
	item.AddOrDrop(p, item.NewSpecialItem(item.PartnerPackageType{}, pp.Count))
}

// Run ...
func (pa PartnerPackageAll) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	var i int
	for t := range internal.Players(tx) {
		item.AddOrDrop(t, item.NewSpecialItem(item.PartnerPackageType{}, pa.Count))
		i++
	}

	internal.Broadcastf(tx, "command.partner_package.all.success", pa.Count, i)
}

func (PartnerItemsRefresh) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	item.RefreshPartnerItems()
}

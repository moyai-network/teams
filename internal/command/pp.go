package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal"
	it "github.com/moyai-network/teams/internal/item"
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
func (pp PartnerPackage) Run(s cmd.Source, o *cmd.Output) {
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
	it.AddOrDrop(p, it.NewSpecialItem(it.PartnerPackageType{}, pp.Count))
}

// Run ...
func (pa PartnerPackageAll) Run(s cmd.Source, o *cmd.Output) {
	var i int
	for _, t := range internal.Players() {
		it.AddOrDrop(t, it.NewSpecialItem(it.PartnerPackageType{}, pa.Count))
		i++
	}

	internal.Broadcastf("command.partner_package.all.success", pa.Count, i)
}

func (PartnerItemsRefresh) Run(s cmd.Source, o *cmd.Output) {
	it.RefreshPartnerItems()
}

package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
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

// Run ...
func (pp PartnerPackage) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	t, ok := pp.Targets[0].(*player.Player)
	if !ok {
		moyai.Messagef(p, "command.target.unknown")
		return
	}

	moyai.Messagef(p, "command.partner_package.give.success", t.Name(), pp.Count)
	it.AddOrDrop(p, it.NewPartnerPackage(pp.Count))
	t.Message("command.partner_package.give.received", pp.Count)
}

// Run ...
func (pa PartnerPackageAll) Run(s cmd.Source, o *cmd.Output) {
	for _, t := range moyai.Players() {
		it.AddOrDrop(t, it.NewPartnerPackage(pa.Count))
		t.Message("command.partner_package.give.received", pa.Count)
	}

	o.Print(text.Colourf("<green>Successfully gave partner packages to all online players.</green>"))
}

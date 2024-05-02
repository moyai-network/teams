package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// PartnerPackage is a command that allows admins to give players partner packages
type PartnerPackage struct {
	Targets []cmd.Target `cmd:"target"`
	Count   int          `cmd:"count"`
}

// PartnerPackageAll is a command that distributes partner packages to the server
type PartnerPackageAll struct {
	Sub   cmd.SubCommand `cmd:"all"`
	Count int            `cmd:"count"`
}

// Run ...
func (pp PartnerPackage) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := pp.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	t, ok := user.Lookup(p.Name())
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	o.Print(lang.Translatef(l, "command.partner_package.give.success", p.Name(), pp.Count))
	t.AddItemOrDrop(it.NewPartnerPackage(pp.Count))
	t.Message("command.partner_package.give.received", pp.Count)
}

// Run ...
func (pa PartnerPackageAll) Run(s cmd.Source, o *cmd.Output) {
	for _, t := range moyai.Server().Players() {
		//t.AddItemOrDrop(it.NewPartnerPackage(pa.Count))
		t.Message("command.partner_package.give.received", pa.Count)
	}

	o.Print(text.Colourf("<green>Successfully gave partner packages to all online players.</green>"))
}

// Allow ...
func (PartnerPackage) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

// Allow ...
func (PartnerPackageAll) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

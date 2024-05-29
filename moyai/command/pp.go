package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/role"
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
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	t, ok := pp.Targets[0].(*player.Player)
	if !ok {
		user.Messagef(p, "command.target.unknown")
		return
	}

	user.Messagef(p, "command.partner_package.give.success", t.Name(), pp.Count)
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

// Allow ...
func (PartnerPackage) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

// Allow ...
func (PartnerPackageAll) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

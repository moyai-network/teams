package command

import (
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/user"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// AliasOnline is a command used to check the alt accounts of an online player.
type AliasOnline struct {
	Targets []cmd.Target `cmd:"target"`
}

// AliasOffline is a command used to check the alt accounts of an offline player.
type AliasOffline struct {
	Target string `cmd:"target"`
}

// Run ...
func (a AliasOnline) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	if len(a.Targets) > 1 {
		user.Messagef(p, "command.targets.exceed")
		return
	}

	target, ok := a.Targets[0].(*player.Player)
	if !ok {
		user.Messagef(p, "command.target.unknown")
		return
	}

	u, err := data.LoadUserFromName(target.Name())
	if err != nil {
		user.Messagef(p, "command.target.unknown")
		return
	}

	usersIPs, _ := data.LoadUsersFromAddress(u.Address)
	ipNames := names(usersIPs, true)

	usersDID, _ := data.LoadUsersFromDeviceID(u.DeviceID)
	deviceNames := names(usersDID, true)

	usersSSID, _ := data.LoadUsersFromSelfSignedID(u.SelfSignedID)
	ssidNames := names(usersSSID, true)

	g := text.Colourf("<grey> - </grey>")
	user.Messagef(p, "command.alias.accounts",
		target.Name(), strings.Join(ipNames, g),
		target.Name(), strings.Join(deviceNames, g),
		target.Name(), strings.Join(ssidNames, g),
	)
}

// Run ...
func (a AliasOffline) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	u, err := data.LoadUserFromName(a.Target)
	if err != nil {
		o.Error(lang.Translate(l, "command.target.unknown"))
		return
	}

	usersIPs, _ := data.LoadUsersFromAddress(u.Address)
	ipNames := names(usersIPs, true)

	usersDID, _ := data.LoadUsersFromDeviceID(u.DeviceID)
	deviceNames := names(usersDID, true)

	usersSSID, _ := data.LoadUsersFromSelfSignedID(u.SelfSignedID)
	ssidNames := names(usersSSID, true)

	g := text.Colourf("<grey> - </grey>")
	o.Print(lang.Translatef(l, "command.alias.accounts",
		u.DisplayName, strings.Join(ipNames, g),
		u.DisplayName, strings.Join(deviceNames, g),
		u.DisplayName, strings.Join(ssidNames, g)),
	)
}

// Allow ...
func (AliasOnline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

// Allow ...
func (AliasOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

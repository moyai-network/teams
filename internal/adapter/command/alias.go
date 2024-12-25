package command

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/pkg/lang"
)

// AliasOnline is a command used to check the alt accounts of an online player.
type AliasOnline struct {
	adminAllower
	Targets []cmd.Target `cmd:"target"`
}

// AliasOffline is a command used to check the alt accounts of an offline player.
type AliasOffline struct {
	adminAllower
	Target string `cmd:"target"`
}

// Run ...
func (a AliasOnline) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	if len(a.Targets) > 1 {
		internal.Messagef(p, "command.targets.exceed")
		return
	}

	target, ok := a.Targets[0].(*player.Player)
	if !ok {
		internal.Messagef(p, "command.target.unknown")
		return
	}

	_, ok = core.UserRepository.FindByName(target.Name())
	if !ok {
		internal.Messagef(p, "command.target.unknown")
		return
	}

	/*usersIPs, _ := data.LoadUsersFromAddress(u.Address)
	ipNames := names(usersIPs, true)

	usersDID, _ := data.LoadUsersFromDeviceID(u.DeviceID)
	deviceNames := names(usersDID, true)

	usersSSID, _ := data.LoadUsersFromSelfSignedID(u.SelfSignedID)
	ssidNames := names(usersSSID, true)

	g := text.Colourf("<grey> - </grey>")
	internal.Messagef(p, "command.alias.accounts") /*target.Name(), strings.Join(ipNames, g),
	target.Name(), strings.Join(deviceNames, g),
	target.Name(), strings.Join(ssidNames, g),*/

}

// Run ...
func (a AliasOffline) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(s)
	_, ok := core.UserRepository.FindByName(a.Target)
	if !ok {
		o.Error(lang.Translate(l, "command.target.unknown"))
		return
	}

	/*usersIPs, _ := data.LoadUsersFromAddress(u.Address)
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
	)*/
}

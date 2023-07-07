package command

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.mongodb.org/mongo-driver/bson"
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
	l := locale(s)
	if len(a.Targets) > 1 {
		o.Error(lang.Translate(l, "command.targets.exceed"))
		return
	}
	target, ok := a.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translate(l, "command.target.unknown"))
		return
	}
	u, err := data.LoadUser(target.Name(), target.Handler().(*user.Handler).XUID())
	if err != nil {
		o.Error(lang.Translate(l, "command.target.unknown"))
		return
	}

	usersIPs, _ := data.LoadUsersCond(bson.M{"address": u.Address})
	ipNames := names(usersIPs, true)

	usersDID, _ := data.LoadUsersCond(bson.M{"did": u.DeviceID})
	deviceNames := names(usersDID, true)

	usersSSID, _ := data.LoadUsersCond(bson.M{"ssid": u.SelfSignedID})
	ssidNames := names(usersSSID, true)

	g := text.Colourf("<grey> - </grey>")
	o.Print(lang.Translatef(l, "command.alias.accounts",
		target.Name(), strings.Join(ipNames, g),
		target.Name(), strings.Join(deviceNames, g),
		target.Name(), strings.Join(ssidNames, g)),
	)
}

// Run ...
func (a AliasOffline) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	users, _ := data.LoadUsersCond(bson.M{"name": strings.ToLower(a.Target)})
	if len(users) == 0 {
		o.Error(lang.Translate(l, "command.target.unknown"))
		return
	}
	u := users[0]

	usersIPs, _ := data.LoadUsersCond(bson.M{"address": u.Address})
	ipNames := names(usersIPs, true)

	usersDID, _ := data.LoadUsersCond(bson.M{"did": u.DeviceID})
	deviceNames := names(usersDID, true)

	usersSSID, _ := data.LoadUsersCond(bson.M{"ssid": u.SelfSignedID})
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

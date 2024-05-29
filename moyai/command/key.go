package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/user"
)

// Key is a command that allows admins to give players keys.
type Key struct {
	Targets []cmd.Target `cmd:"target"`
	Count   int          `cmd:"count"`
}

// KeyAll is a command that distributes keys to the server
type KeyAll struct {
	Sub   cmd.SubCommand `cmd:"all"`
	Count int            `cmd:"count"`
}

// Run ...
func (k Key) Run(s cmd.Source, o *cmd.Output) {
	// l := locale(s)
	// p, ok := k.Targets[0].(*player.Player)
	// if !ok {
	// 	o.Error(lang.Translatef(l, "command.target.unknown"))
	// 	return
	// }
	// t, ok := user.Lookup(p.Name())
	// if !ok {
	// 	o.Error(lang.Translatef(l, "command.target.unknown"))
	// 	return
	// }
	// key type
}

// Run ...
func (k KeyAll) Run(s cmd.Source, o *cmd.Output) {
	for _, t := range moyai.Players() {
		//t.AddItemOrDrop(it.NewKey(it.KeyTypePharaoh, k.Count))
		t.Message("command.key.give.received", k.Count)
	}

	user.Alertf(s, "command.key.all.success", k.Count)
}

// Allow ...
func (Key) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

// Allow ...
func (KeyAll) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

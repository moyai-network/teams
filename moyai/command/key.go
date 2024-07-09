package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/user"
)

// Key is a command that allows admins to give players keys.
type Key struct {
	adminAllower
	Targets []cmd.Target `cmd:"target"`
	Count   int          `cmd:"count"`
	Key     key          `cmd:"key"`
}

// KeyAll is a command that distributes keys to the server
type KeyAll struct {
	adminAllower
	Sub   cmd.SubCommand `cmd:"all"`
	Count int            `cmd:"count"`
	Key   key            `cmd:"key"`
}

// Run ...
func (k Key) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := k.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	t, ok := user.Lookup(p.Name())
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	keyType := keyToType(string(k.Key))
	it.AddOrDrop(t, it.NewKey(keyType, k.Count))
	moyai.Alertf(s, "command.key.give.received", k.Count)
}

// Run ...
func (k KeyAll) Run(s cmd.Source, o *cmd.Output) {
	for _, t := range moyai.Players() {
		keyType := keyToType(string(k.Key))
		it.AddOrDrop(t, it.NewKey(keyType, k.Count))
	}

	moyai.Broadcastf("command.key.all.success", k.Count)
}

type key string

// Type ...
func (key) Type() string {
	return "key"
}

// Options ...
func (key) Options(s cmd.Source) []string {
	return []string{
		"koth",
		"pharaoh",
		"partner",
		"menes",
		"ramses",
		"conquest",
	}
}

func keyToType(k string) int {
	switch k {
	case "koth":
		return 0
	case "pharaoh":
		return 1
	case "partner":
		return 2
	case "menes":
		return 3
	case "ramses":
		return 4
	case "conquest":
		return 5
	default:
		panic("should never happen")
	}
}
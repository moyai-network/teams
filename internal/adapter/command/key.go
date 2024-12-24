package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/core/user"
	"github.com/moyai-network/teams/pkg/lang"
	"strings"
)

// Key is a command that allows admins to give players keys.
type Key struct {
	adminAllower
	Targets []cmd.Target `cmd:"target"`
	Key     key          `cmd:"key"`
	Count   int          `cmd:"count"`
}

// KeyAll is a command that distributes keys to the server
type KeyAll struct {
	adminAllower
	Sub   cmd.SubCommand `cmd:"all"`
	Key   key            `cmd:"key"`
	Count int            `cmd:"count"`
}

// Run ...
func (k Key) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(s)
	p, ok := k.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	t, ok := user.Lookup(tx, p.Name())
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	keyType, ok := keyToType(string(k.Key))
	if !ok {
		return
	}
	item.AddOrDrop(t, item.NewKey(keyType, k.Count))
	internal.Messagef(p, "command.key.give.received", k.Count)
}

// Run ...
func (k KeyAll) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	keyType, ok := keyToType(string(k.Key))
	if !ok {
		return
	}
	var i int
	for t := range internal.Players(tx) {
		item.AddOrDrop(t, item.NewKey(keyType, k.Count))
		i++
	}

	internal.Broadcastf(tx, "command.key.all.success", k.Count, strings.Title(string(k.Key)), i)
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
		"seasonal",
	}
}

func keyToType(k string) (item.KeyType, bool) {
	switch k {
	case "koth":
		return item.KeyTypeKOTH, true
	case "pharaoh":
		return item.KeyTypePharaoh, true
	case "partner":
		return item.KeyTypePartner, true
	case "menes":
		return item.KeyTypeMenes, true
	case "ramses":
		return item.KeyTypeRamses, true
	case "conquest":
		return item.KeyTypeConquest, true
	case "seasonal":
		return item.KeyTypeSeasonal, true
	default:
		return item.KeyType{}, false
	}
}

package command

import (
    "github.com/df-mc/dragonfly/server/cmd"
    "github.com/df-mc/dragonfly/server/player"
    "github.com/moyai-network/teams/internal/lang"
    "github.com/moyai-network/teams/moyai"
    it "github.com/moyai-network/teams/moyai/item"
    "github.com/moyai-network/teams/moyai/user"
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

    keyType, ok := keyToType(string(k.Key))
    if !ok {
        return
    }
    it.AddOrDrop(t, it.NewKey(keyType, k.Count))
    moyai.Messagef(p, "command.key.give.received", k.Count)
}

// Run ...
func (k KeyAll) Run(s cmd.Source, o *cmd.Output) {
    keyType, ok := keyToType(string(k.Key))
    if !ok {
        return
    }
    var i int
    for _, t := range moyai.Players() {
        it.AddOrDrop(t, it.NewKey(keyType, k.Count))
        i++
    }

    moyai.Broadcastf("command.key.all.success", k.Count, strings.Title(string(k.Key)), i)
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

func keyToType(k string) (it.KeyType, bool) {
    switch k {
    case "koth":
        return it.KeyTypeKOTH, true
    case "pharaoh":
        return it.KeyTypePharaoh, true
    case "partner":
        return it.KeyTypePartner, true
    case "menes":
        return it.KeyTypeMenes, true
    case "ramses":
        return it.KeyTypeRamses, true
    case "conquest":
        return it.KeyTypeConquest, true
    case "seasonal":
        return it.KeyTypeSeasonal, true
    default:
        return it.KeyType{}, false
    }
}

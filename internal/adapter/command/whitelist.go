package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/pkg/lang"
)

type WhiteListAdd struct {
	operatorAllower
	Sub    cmd.SubCommand `cmd:"add"`
	Target string         `cmd:"target"`
}

func (w WhiteListAdd) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	if len(w.Target) <= 0 {
		internal.Messagef(src, "invalid.username")
		return
	}
	u, ok := core.UserRepository.FindByName(w.Target)
	if !ok {
		internal.Messagef(src, "target.data.load.error", w.Target)
		return
	}
	if u.Whitelisted {
		internal.Messagef(src, "whitelist.already", u.DisplayName)
		return
	}

	u.Whitelisted = true
	core.UserRepository.Save(u)

	internal.Messagef(src, "whitelist.add", u.DisplayName)
}

type WhiteListRemove struct {
	operatorAllower
	Sub    cmd.SubCommand `cmd:"remove"`
	Target string         `cmd:"target"`
}

func (w WhiteListRemove) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	l := locale(src)
	if len(w.Target) <= 0 {
		out.Error(lang.Translatef(l, "invalid.username"))
		return
	}
	u, ok := core.UserRepository.FindByName(w.Target)
	if !ok {
		out.Error(lang.Translatef(l, "target.data.load.error", w.Target))
		return
	}
	if !u.Whitelisted {
		out.Error(lang.Translatef(l, "whitelist.already.not", u.DisplayName))
		return
	}

	u.Whitelisted = false
	core.UserRepository.Save(u)

	out.Print(lang.Translatef(l, "whitelist.remove", u.DisplayName))
}

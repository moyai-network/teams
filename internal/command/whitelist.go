package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/data"
	"github.com/moyai-network/teams/internal/role"
	"github.com/moyai-network/teams/internal/user"
	"github.com/moyai-network/teams/pkg/lang"
)

type WhiteListAdd struct {
	Sub    cmd.SubCommand `cmd:"add"`
	Target string         `cmd:"target"`
}

func (w WhiteListAdd) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	if len(w.Target) <= 0 {
		user.Messagef(p, "invalid.username")
		return
	}
	u, err := data.LoadUserFromName(w.Target)
	if err != nil {
		user.Messagef(p, "target.data.load.error", w.Target)
		return
	}
	if u.Whitelisted {
		user.Messagef(p, "whitelist.already", u.DisplayName)
		return
	}

	u.Whitelisted = true
	data.SaveUser(u)

	user.Messagef(p, "whitelist.add", u.DisplayName)
}

type WhiteListRemove struct {
	Sub    cmd.SubCommand `cmd:"remove"`
	Target string         `cmd:"target"`
}

func (w WhiteListRemove) Run(src cmd.Source, out *cmd.Output) {
	l := locale(src)
	if len(w.Target) <= 0 {
		out.Error(lang.Translatef(l, "invalid.username"))
		return
	}
	u, err := data.LoadUserFromName(w.Target)
	if err != nil {
		out.Error(lang.Translatef(l, "target.data.load.error", w.Target))
		return
	}
	if !u.Whitelisted {
		out.Error(lang.Translatef(l, "whitelist.already.not", u.DisplayName))
		return
	}

	u.Whitelisted = false
	data.SaveUser(u)

	out.Print(lang.Translatef(l, "whitelist.remove", u.DisplayName))
}

func (w WhiteListRemove) Allow(src cmd.Source) bool {
	return allow(src, true, role.Operator{})
}

func (w WhiteListAdd) Allow(src cmd.Source) bool {
	return allow(src, true, role.Operator{})
}

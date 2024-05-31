package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
)

type WhiteListAdd struct {
	operatorAllower
	Sub    cmd.SubCommand `cmd:"add"`
	Target string         `cmd:"target"`
}

func (w WhiteListAdd) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	if len(w.Target) <= 0 {
		moyai.Messagef(p, "invalid.username")
		return
	}
	u, err := data.LoadUserFromName(w.Target)
	if err != nil {
		moyai.Messagef(p, "target.data.load.error", w.Target)
		return
	}
	if u.Whitelisted {
		moyai.Messagef(p, "whitelist.already", u.DisplayName)
		return
	}

	u.Whitelisted = true
	data.SaveUser(u)

	moyai.Messagef(p, "whitelist.add", u.DisplayName)
}

type WhiteListRemove struct {
	operatorAllower
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

package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// PvpEnable is a command to enable PVP.
type PvpEnable struct {
	Sub cmd.SubCommand `cmd:"enable"`
}

// Run ...
func (c PvpEnable) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUser(p.Name(), p.Handler().(*user.Handler).XUID())
	if err != nil {
		return
	}
	if u.PVP.Active() {
		u.PVP.Reset()
		out.Print(text.Colourf("<green>You have enabled PVP!</green>"))
	} else {
		out.Error(text.Colourf("<red>You have already have enabled PVP</red>"))
	}
}

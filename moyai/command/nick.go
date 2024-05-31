package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
)

// Nick is a command that allows the player to change their nickname.
type Nick struct {
	adminAllower
	Name string `cmd:"name"`
}

// NickReset is a command that allows the player to reset their nickname.
type NickReset struct{ adminAllower }

// Run ...
func (n Nick) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	u.DisplayName = n.Name
	data.SaveUser(u)
	moyai.Messagef(p, "nick.changed", n.Name)
}

// Run ...
func (n NickReset) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	u.DisplayName = p.Name()
	data.SaveUser(u)
	moyai.Messagef(p, "nick.reset")
}

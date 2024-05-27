package command

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/menu"
)

// Kit is a command that allows players to select a kit.
type Kit struct{}

// KitReset is a command to reset a players kit.
type KitReset struct {
	adminAllower
	Sub    cmd.SubCommand           `cmd:"reset"`
	Target cmd.Optional[cmd.Target] `cmd:"target"`
}

// Run ...
func (Kit) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	inv.SendMenu(p, menu.NewKitsMenu(p))
}

// Run ...
func (k KitReset) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	t, ok := k.Target.Load()
	if !ok {
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			return
		}
		for _, kt := range u.Teams.Kits {
			kt.Reset()
		}
		data.SaveUser(u)
		return
	}
	tg, ok := t.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(tg.Name())
	if err != nil {
		return
	}
	for _, kt := range u.Teams.Kits {
		kt.Reset()
	}
	data.SaveUser(u)
}

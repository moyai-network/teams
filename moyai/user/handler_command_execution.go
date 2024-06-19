package user

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
)

var (
	spawn = []string{}
	combat = []string{}
	deathban = []string{
		"reclaim",
		"trim",
		"prizes",
		"logout",
		"ec",
	}
)

func (h *Handler) HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string) {
	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil {
		return
	}

	if u.Teams.DeathBan.Active() {
		for _, d := range deathban {
			names := []string{command.Name()}
			names = append(names, command.Aliases()...)
			for _, n := range names {
				if strings.EqualFold(d, n) {
					moyai.Messagef(h.p, "deathban.cooldown")
					ctx.Cancel()
				}
			}

		}
	}
}
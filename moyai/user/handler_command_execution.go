package user

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
)

var (
	spawn  = []string{}
	combat = []string{
		"ec",
		"logout",
		"pv",
	}
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
	names := []string{command.Name()}
	names = append(names, command.Aliases()...)

	if h.tagCombat.Active() && containsAny(combat, names...) {
		moyai.Messagef(h.p, "command.error.combat-tagged")
		ctx.Cancel()
	} else if u.Teams.DeathBan.Active() && containsAny(deathban, names...) {
		moyai.Messagef(h.p, "command.error.death-ban")
		ctx.Cancel()
	}
}

func containsAny(s []string, e ...string) bool {
	for _, a := range e {
		if contains(s, a) {
			return true
		}
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}

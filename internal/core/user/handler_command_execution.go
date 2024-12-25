package user

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/core"
	rls "github.com/moyai-network/teams/internal/core/roles"
	"slices"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/internal"
)

var (
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

func (h *Handler) HandleCommandExecution(ctx *player.Context, command cmd.Command, args []string) {
	p := ctx.Val()

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	names := append(command.Aliases(), command.Name())

	whisper := slices.Contains(names, "whisper")
	reply := slices.Contains(names, "reply")
	if whisper || reply {
		if time.Since(h.lastMessage.Load()) < internal.ChatCoolDown() && !u.Roles.Contains(rls.Admin()) {
			ctx.Cancel()
			internal.Messagef(p, "chat.cooldown", time.Until(h.lastMessage.Load().Add(internal.ChatCoolDown())).Seconds())
			return
		}
		h.lastMessage.Store(time.Now().Add(internal.ChatCoolDown()))
	}

	if h.tagCombat.Active() && containsAny(combat, names...) {
		internal.Messagef(p, "command.error.combat-tagged")
		ctx.Cancel()
	} else if u.Teams.DeathBan.Active() && containsAny(deathban, names...) {
		internal.Messagef(p, "command.error.death-ban")
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

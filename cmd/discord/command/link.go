package command

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/moyai-network/teams/moyai/data"
	"strings"
)

func (h *Handler) link(ctx context.Context, d cmdroute.CommandData) *api.InteractionResponseData {
	u, err := data.LinkUser(d.Options[0].String(), d.Event.Sender())
	if err != nil {
		return h.error("Failed to link: " + err.Error())
	}
	_ = h.s.AddRole(h.guildID, d.Event.Sender().ID, discord.RoleID(1209213713337548830), api.AddRoleData{AuditLogReason: "Linking"})
	_ = h.s.RemoveRole(h.guildID, d.Event.Sender().ID, discord.RoleID(1198584931459272806), "Linking")
	return h.success(fmt.Sprintf("Votre compte Minecraft **%s** est désormais lié à votre compte Discord !", strings.ReplaceAll(u.DisplayName, "_", " ")))
}

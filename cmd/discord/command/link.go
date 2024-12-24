package command

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/data"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
)

func (h *Handler) link(ctx context.Context, d cmdroute.CommandData) *api.InteractionResponseData {
	u, err := data.LinkUser(d.Options[0].String(), d.Event.Sender())
	if err != nil {
		return h.error("Failed to link: " + err.Error())
	}
	userID := d.Event.Sender().ID

	_ = h.s.ModifyMember(h.guildID, userID, api.ModifyMemberData{
		Nick: option.NewString(u.DisplayName),
	})
	_ = h.s.AddRole(h.guildID, userID, discord.RoleID(1255290630922436698), api.AddRoleData{AuditLogReason: "Linking"})

	for _, p := range internal.Players() {
		if p.Name() == u.DisplayName {
			internal.Messagef(p, "discord.linked", d.Event.Sender().Username)
		}
	}
	return h.success(fmt.Sprintf("Your MC Account (**%s**) has been linked to the Discord!", strings.ReplaceAll(u.DisplayName, "_", " ")))
}

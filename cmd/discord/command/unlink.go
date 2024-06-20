package command

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/moyai-network/teams/moyai/data"
	"strings"
)

func (h *Handler) unlink(ctx context.Context, d cmdroute.CommandData) *api.InteractionResponseData {
	u, err := data.LoadUserFromDiscordID(d.Event.SenderID().String())
	if err != nil {
		return h.error("Vous n'êtes pas lié à un compte")
	}

	err = data.UnlinkUser(u, h.s, h.guildID)
	if err != nil {
		return h.error("Une erreur est survenue lors du deliage de votre compte : " + err.Error())
	}

	return h.success(fmt.Sprintf("Votre compte Minecraft **%s** est désormais délié de votre compte Discord !", strings.ReplaceAll(u.DisplayName, "_", " ")))
}

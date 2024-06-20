package command

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

var commands = []api.CreateCommandData{
	{Name: "link", Description: "link", Options: discord.CommandOptions{
		discord.NewStringOption("code", "code", true),
	}},
	{Name: "unlink", Description: "unlink", Options: discord.CommandOptions{}},
}

type Handler struct {
	s       *state.State
	r       *cmdroute.Router
	guildID discord.GuildID
}

func NewHandler(r *cmdroute.Router, s *state.State, g discord.GuildID) *Handler {
	return &Handler{
		r:       r,
		s:       s,
		guildID: g,
	}
}

func (h *Handler) RegisterCommands() {
	h.r.AddFunc("link", h.link)
	h.r.AddFunc("unlink", h.unlink)

	app, _ := h.s.CurrentApplication()
	for _, c := range commands {
		_, _ = h.s.CreateGuildCommand(app.ID, h.guildID, c)
	}
}

func (h *Handler) error(msg string) *api.InteractionResponseData {
	return &api.InteractionResponseData{Embeds: &[]discord.Embed{
		coloredEmbedMessage(msg, discord.Color(0xAA0000)),
	}, Flags: discord.EphemeralMessage}
}

func (h *Handler) success(msg string) *api.InteractionResponseData {
	return &api.InteractionResponseData{Embeds: &[]discord.Embed{
		coloredEmbedMessage(msg, discord.Color(0x55FF55)),
	}, Flags: discord.EphemeralMessage}
}

func coloredEmbedMessage(msg string, color discord.Color) discord.Embed {
	return discord.Embed{
		Description: msg,
		Color:       color,
	}
}

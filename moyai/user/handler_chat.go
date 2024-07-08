package user

import (
	"regexp"
	"strings"
	"time"

	"github.com/moyai-network/teams/moyai/roles"

	"github.com/bedrock-gophers/role/role"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/tag"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// formatRegex is a regex used to clean color formatting on a string.
var formatRegex = regexp.MustCompile(`§[\da-gk-or]`)

// HandleChat ...
func (h *Handler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()
	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil {
		return
	}

	*message = emojis.Replace(*message)
	r := u.Roles.Highest()

	if !u.Teams.Mute.Expired() {
		h.p.Message(lang.Translatef(*u.Language, "user.message.mute"))
		return
	}
	tm, teamErr := data.LoadTeamFromMemberName(h.p.Name())
	msg := strings.TrimSpace(*message)
	if len(msg) <= 0 {
		return
	}
	msg = formatRegex.ReplaceAllString(msg, "")

	switch u.Teams.ChatType {
	case 0:
		if msg[0] == '!' && roles.Staff(r) {
			h.staffMessage(msg, r)
			return
		}
		h.globalMessage(msg, u, r, tm)
	case 1:
		if msg[0] == '!' && roles.Staff(r) {
			h.staffMessage(msg, r)
			return
		}

		if teamErr != nil {
			u.Teams.ChatType = 1
			data.SaveUser(u)
			h.globalMessage(msg, u, r, tm)
			return
		}
		for _, member := range tm.Members {
			if m, ok := Lookup(member.Name); ok {
				m.Message(text.Colourf("<dark-aqua>[<yellow>T</yellow>] %s: %s</dark-aqua>", h.p.Name(), msg))
			}
		}
	case 2:
		if msg[0] == '!' {
			h.globalMessage(msg, u, r, tm)
			return
		}
		h.staffMessage(msg, r)
	}
}

func (h *Handler) staffMessage(msg string, r role.Role) {
	for _, s := range moyai.Players() {
		if us, err := data.LoadUserOrCreate(s.Name(), s.XUID()); err == nil && roles.Staff(us.Roles.Highest()) {
			moyai.Messagef(s, "staff.chat", r.Name(), h.p.Name(), strings.TrimPrefix(msg, "!"))
		}
	}
}

func (h *Handler) globalMessage(msg string, u data.User, r role.Role, tm data.Team) {
	if !moyai.GlobalChatEnabled() {
		moyai.Messagef(h.p, "chat.global.muted")
		return
	}
	if time.Since(h.lastMessage.Load()) < moyai.ChatCoolDown() && !u.Roles.Contains(roles.Admin()) {
		moyai.Messagef(h.p, "chat.cooldown", time.Until(h.lastMessage.Load().Add(moyai.ChatCoolDown())).Seconds())
		return
	}
	h.lastMessage.Store(time.Now())
	displayName := u.DisplayName
	if t, ok := tag.ByName(u.Teams.Settings.Display.ActiveTag); ok {
		displayName = u.DisplayName + " " + t.Format()
	}

	highestRole := u.Roles.Highest()
	chatMessage := text.Colourf("%s<dark-grey>:</dark-grey> <white>%s</white>", highestRole.Coloured(displayName), msg)

	if len(tm.Name) > 0 {
		formatTeam := text.Colourf("<grey>[<green>%s</green>]</grey> %s", tm.DisplayName, chatMessage)
		formatEnemy := text.Colourf("<grey>[<red>%s</red>]</grey> %s", tm.DisplayName, chatMessage)

		for _, t := range moyai.Players() {
			if tm.Member(t.Name()) {
				t.Message(formatTeam)
			} else {
				t.Message(formatEnemy)
			}
		}
		chat.StdoutSubscriber{}.Message(formatEnemy)
	} else {
		_, _ = chat.Global.WriteString(chatMessage)
	}
}

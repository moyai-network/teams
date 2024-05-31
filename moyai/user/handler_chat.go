package user

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/tag"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"regexp"
	"strings"
	"time"
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
		h.p.Message(lang.Translatef(u.Language, "user.message.mute"))
		return
	}
	tm, teamErr := data.LoadTeamFromMemberName(h.p.Name())

	if msg := strings.TrimSpace(*message); len(msg) > 0 {
		msg = formatRegex.ReplaceAllString(msg, "")

		global := func() {
			if !moyai.GlobalChatEnabled() {
				moyai.Messagef(h.p, "chat.global.muted")
				return
			}
			if time.Since(h.lastMessage.Load()) < moyai.ChatCoolDown() && !u.Roles.Contains(role.Admin{}) {
				moyai.Messagef(h.p, "chat.cooldown", time.Until(h.lastMessage.Load().Add(moyai.ChatCoolDown())).Seconds())
				return
			}
			h.lastMessage.Store(time.Now())
			displayName := u.DisplayName
			if t, ok := tag.ByName(u.Teams.Settings.Display.ActiveTag); ok {
				displayName = u.DisplayName + " " + t.Format()
			}

			if teamErr == nil {

				formatTeam := text.Colourf("<grey>[<green>%s</green>]</grey> %s", tm.DisplayName, r.Chat(displayName, msg))
				formatEnemy := text.Colourf("<grey>[<red>%s</red>]</grey> %s", tm.DisplayName, r.Chat(displayName, msg))

				for _, t := range moyai.Players() {
					if tm.Member(t.Name()) {
						t.Message(formatTeam)
					} else {
						t.Message(formatEnemy)
					}
				}
				chat.StdoutSubscriber{}.Message(formatEnemy)
			} else {
				_, _ = chat.Global.WriteString(r.Chat(displayName, msg))
			}
		}

		staff := func() {
			for _, s := range moyai.Players() {
				if us, err := data.LoadUserOrCreate(s.Name(), s.XUID()); err == nil && role.Staff(us.Roles.Highest()) {
					moyai.Messagef(s, "staff.chat", r.Name(), h.p.Name(), strings.TrimPrefix(msg, "!"))
				}
			}
		}
		switch u.Teams.ChatType {
		case 0:
			if msg[0] == '!' && role.Staff(r) {
				staff()
				return
			}
			global()
		case 1:
			if msg[0] == '!' && role.Staff(r) {
				staff()
				return
			}

			if teamErr != nil {
				u.Teams.ChatType = 1
				data.SaveUser(u)
				global()
				return
			}
			for _, member := range tm.Members {
				if m, ok := Lookup(member.Name); ok {
					m.Message(text.Colourf("<dark-aqua>[<yellow>T</yellow>] %s: %s</dark-aqua>", h.p.Name(), msg))
				}
			}
		case 2:
			if msg[0] == '!' {
				global()
				return
			}
			staff()
		}
	}
}

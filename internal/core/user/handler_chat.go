package user

import (
	"github.com/df-mc/dragonfly/server/player"
	data2 "github.com/moyai-network/teams/internal/core/data"
	item2 "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/core/roles"
	"github.com/moyai-network/teams/internal/ports/model"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/bedrock-gophers/tag/tag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/bedrock-gophers/role/role"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var (
	// formatRegex is a regex used to clean color formatting on a string.
	formatRegex  = regexp.MustCompile(`§[\da-gk-or]`)
	englishCaser = cases.Title(language.English)
)

// HandleChat ...
func (h *Handler) HandleChat(ctx *player.Context, message *string) {
	p := ctx.Val()

	ctx.Cancel()
	u, err := data2.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	defer data2.SaveUser(u)

	if internal.ChatGameWord() != "" && *message == internal.ChatGameWord() {
		internal.Broadcastf(p.Tx(), "internal.broadcast.chatgame.guessed", p.Name(), internal.ChatGameWord())
		internal.SetChatGameWord("")
		if !u.Teams.DeathBan.Active() {
			item2.AddOrDrop(p, item2.NewKey(item2.KeyTypePharaoh, rand.Intn(10)+1))
		} else {
			u.Teams.DeathBan.Reduce(time.Minute * 5)
		}
		return
	}

	*message = emojis.Replace(*message)
	r := u.Roles.Highest()

	if !u.Teams.Mute.Expired() {
		p.Message(lang.Translatef(*u.Language, "user.message.mute"))
		return
	}
	tm, teamErr := data2.LoadTeamFromMemberName(p.Name())
	msg := strings.TrimSpace(*message)
	if len(msg) <= 0 {
		return
	}
	msg = formatRegex.ReplaceAllString(msg, "")

	if len(msg) <= 0 {
		return
	}

	switch u.Teams.ChatType {
	case 0:
		if msg[0] == '!' && roles.Staff(r) {
			h.staffMessage(p, msg, r)
			return
		}
		h.globalMessage(p, msg, u, r, tm)
	case 1:
		if msg[0] == '!' && roles.Staff(r) {
			h.staffMessage(p, msg, r)
			return
		}

		if teamErr != nil {
			u.Teams.ChatType = 1
			h.globalMessage(p, msg, u, r, tm)
			return
		}
		for _, member := range tm.Members {
			if m, ok := Lookup(p.Tx(), member.Name); ok {
				m.Message(text.Colourf("<dark-aqua>[<yellow>T</yellow>] %s: %s</dark-aqua>", p.Name(), msg))
			}
		}
	case 2:
		if msg[0] == '!' {
			h.globalMessage(p, msg, u, r, tm)
			return
		}
		h.staffMessage(p, msg, r)
	}
}

func (h *Handler) staffMessage(p *player.Player, msg string, r role.Role) {
	for s := range internal.Players(p.Tx()) {
		if us, err := data2.LoadUserOrCreate(s.Name(), s.XUID()); err == nil && roles.Staff(us.Roles.Highest()) {
			internal.Messagef(s, "staff.chat", r.Name(), p.Name(), strings.TrimPrefix(msg, "!"))
		}
	}
}

func (h *Handler) globalMessage(p *player.Player, msg string, u model.User, r role.Role, tm model.Team) {
	if !internal.GlobalChatEnabled() {
		internal.Messagef(p, "chat.global.muted")
		return
	}
	if time.Since(h.lastMessage.Load()) < internal.ChatCoolDown() && !u.Roles.Contains(roles.Admin()) {
		internal.Messagef(p, "chat.cooldown", time.Until(h.lastMessage.Load().Add(internal.ChatCoolDown())).Seconds())
		return
	}
	h.lastMessage.Store(time.Now())
	displayName := u.DisplayName
	if t, ok := tag.ByName(u.Teams.Settings.Display.ActiveTag); ok {
		displayName = u.DisplayName + " " + t.Format()
	}

	highestRole := u.Roles.Highest()
	chatMessage := text.Colourf("%s<dark-grey>:</dark-grey> <white>%s</white>", highestRole.Coloured(displayName), msg)

	if highestRole.Tier() > roles.Default().Tier() {
		roleFormat := text.Colourf("<dark-grey>[</dark-grey>%s<dark-grey>]</dark-grey>", highestRole.Coloured(englishCaser.String(r.Name())))
		chatMessage = roleFormat + " " + chatMessage
	}

	if len(tm.Name) > 0 {
		formatTeam := text.Colourf("<grey>[<green>%s</green>]</grey> %s", tm.DisplayName, chatMessage)
		formatEnemy := text.Colourf("<grey>[<red>%s</red>]</grey> %s", tm.DisplayName, chatMessage)

		for t := range internal.Players(p.Tx()) {
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

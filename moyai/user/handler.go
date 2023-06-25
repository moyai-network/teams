package user

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/class"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
	"strings"
	"time"
)

var (
	// tlds is a list of top level domains used for checking for advertisements.
	tlds = [...]string{".me", ".club", "www.", ".com", ".net", ".gg", ".cc", ".net", ".co", ".co.uk", ".ddns", ".ddns.net", ".cf", ".live", ".ml", ".gov", "http://", "https://", ",club", "www,", ",com", ",cc", ",net", ",gg", ",co", ",couk", ",ddns", ",ddns.net", ",cf", ",live", ",ml", ",gov", ",http://", "https://", "gg/"}
	// emojis is a map between emojis and their unicode representation.
	emojis = strings.NewReplacer(
		":l:", "\uE107",
		":skull:", "\uE105",
		":fire:", "\uE108",
		":eyes:", "\uE109",
		":clown:", "\uE10A",
		":100:", "\uE10B",
		":heart:", "\uE10C",
	)
)

type Handler struct {
	player.NopHandler
	s *session.Session
	p *player.Player
	u data.User

	cd struct {
		enderPearl *moose.CoolDown
	}
}

func NewHandler(p *player.Player) *Handler {
	ha := &Handler{
		p: p,
	}
	ha.cd.enderPearl = moose.NewCoolDown()

	s := player_session(p)
	u, _ := data.LoadUser(p)

	u.DeviceID = s.ClientData().DeviceID
	u.SelfSignedID = s.ClientData().SelfSignedID

	ha.u = u
	ha.s = s

	playersMu.Lock()
	players[p.XUID()] = p
	playersMu.Unlock()

	return ha
}

func (h *Handler) HandleItemUse(ctx *event.Context) {
	held, _ := h.p.HeldItems()
	switch held.Item().(type) {
	case item.EnderPearl:
		if cd := h.cd.enderPearl; cd.Active() {
			h.p.Message(lang.Translatef(h.p.Locale(), "user.cool-down", "Ender Pearl", cd.Remaining().Seconds()))
			ctx.Cancel()
		} else {
			cd.Set(15 * time.Second)
		}
	}
}

func (h *Handler) HandleChat(ctx *event.Context, msg *string) {
	*msg = emojis.Replace(*msg)
}

func (h *Handler) HandleQuit() {
	_ = data.SaveUser(h.u)

	playersMu.Lock()
	delete(players, h.p.XUID())
	playersMu.Unlock()
}

// setClass sets the class of the user.
func (h *ClassHandler) setClass(c moose.Class) {
	lastClass := h.class.Load()
	if lastClass != c {
		if class.CompareAny(c, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}) {
			addEffects(h.p, c.Effects()...)
		} else if class.CompareAny(lastClass, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}) {
			h.bardEnergy.Store(0)
			removeEffects(h.p, lastClass.Effects()...)
		}
		h.class.Store(c)
	}
}

package user

import (
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/class"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
	"math"
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

	pearlCooldown *moose.CoolDown

	class      atomic.Value[moose.Class]
	bardEnergy atomic.Value[float64]

	combatTag *moose.Tag
	archerTag *moose.Tag
}

func NewHandler(p *player.Player) *Handler {
	ha := &Handler{
		p: p,

		pearlCooldown: moose.NewCoolDown(),

		combatTag: moose.NewTag(nil, nil),
		archerTag: moose.NewTag(nil, nil),
	}

	s := player_session(p)
	u, _ := data.LoadUser(p)

	u.DeviceID = s.ClientData().DeviceID
	u.SelfSignedID = s.ClientData().SelfSignedID

	_ = data.SaveUser(u)
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
		if cd := h.pearlCooldown; cd.Active() {
			h.p.Message(lang.Translatef(h.p.Locale(), "user.cool-down", "Ender Pearl", cd.Remaining().Seconds()))
			ctx.Cancel()
		} else {
			cd.Set(15 * time.Second)
		}
	}
}

func (h *Handler) HandleHurt(ctx *event.Context, dmg *float64, _ *time.Duration, src world.DamageSource) {
	var target *player.Player

	switch s := src.(type) {
	case entity.AttackDamageSource:
		if t, ok := s.Attacker.(*player.Player); ok {
			target = t
		}
		if !canAttack(h.p, target) {
			ctx.Cancel()
			return
		}
	case entity.ProjectileDamageSource:
		if t, ok := s.Owner.(*player.Player); ok {
			target = t
		}
		if !canAttack(h.p, target) {
			ctx.Cancel()
			return
		}
		if ha := target.Handler().(*Handler); class.Compare(ha.class, class.Archer{}) {
			ha.archerTag.Set(time.Second * 10)

			dist := h.p.Position().Sub(target.Position()).Len()
			d := math.Round(dist)
			if d > 20 {
				d = 20
			}
			*dmg = (d / 10) * 2

			h.p.Message(lang.Translatef(h.p.Locale(), "archer.tag", math.Round(dist), *dmg/2))
		}
	}
	target.Handler().(*Handler).combatTag.Set(time.Second * 20)
	h.combatTag.Set(time.Second * 20)
}

func (h *Handler) HandleItemDamage(_ *event.Context, i item.Stack, n int) {
	dur := i.Durability()
	if _, ok := i.Item().(item.Armour); ok && dur != -1 && dur-n <= 0 {
		setClass(h.p, class.Resolve(h.p))
	}
}

func (h *Handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	t, ok := e.(*player.Player)
	if !ok {
		return
	}
	if !canAttack(h.p, t) {
		ctx.Cancel()
		return
	}
}

func (h *Handler) HandleChat(ctx *event.Context, msg *string) {
	*msg = emojis.Replace(*msg)
}

func (h *Handler) HandleQuit() {
	playersMu.Lock()
	delete(players, h.p.XUID())
	playersMu.Unlock()
}

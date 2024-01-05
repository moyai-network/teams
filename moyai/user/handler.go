package user

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/moyai-network/teams/moyai/kit"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/moyai-network/moose/crate"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/area"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/koth"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/restartfu/roman"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/text/language"

	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/class"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
	ench "github.com/moyai-network/teams/moyai/enchantment"
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

	xuid    string
	logTime time.Time

	sign cube.Pos

	pearl       *moose.CoolDown
	rogue       *moose.CoolDown
	goldenApple *moose.CoolDown

	factionCreate *moose.CoolDown

	itemUse   moose.MappedCoolDown[world.Item]
	bardItem  moose.MappedCoolDown[world.Item]
	strayItem moose.MappedCoolDown[world.Item]

	ability   *moose.CoolDown
	abilities moose.MappedCoolDown[it.SpecialItemType]

	waypoint *WayPoint

	armour atomic.Value[[4]item.Stack]
	class  atomic.Value[moose.Class]
	energy atomic.Value[float64]

	combat *moose.Tag
	archer *moose.Tag

	boneHits map[string]int
	bone     *moose.CoolDown

	scramblerHits map[string]int
	pearlDisabled bool

	logout *moose.Teleportation
	stuck  *moose.Teleportation
	home   *moose.Teleportation

	lastScoreBoard atomic.Value[*scoreboard.Scoreboard]
	area           atomic.Value[moose.NamedArea]

	lastAttackerName atomic.Value[string]
	lastAttackTime   atomic.Value[time.Time]

	lastMessage atomic.Value[time.Time]
	chatType    atomic.Value[int]

	wallBlocks   map[cube.Pos]float64
	wallBlocksMu sync.Mutex

	claimPos [2]mgl64.Vec2

	loggedOut bool
	logger    bool
	close     chan struct{}
	death     chan struct{}
}

func NewHandler(p *player.Player, xuid string) *Handler {
	ha := &Handler{
		p:    p,
		xuid: xuid,

		pearl:       moose.NewCoolDown(),
		rogue:       moose.NewCoolDown(),
		goldenApple: moose.NewCoolDown(),
		ability:     moose.NewCoolDown(),

		factionCreate: moose.NewCoolDown(),

		bone:     moose.NewCoolDown(),
		boneHits: map[string]int{},

		scramblerHits: map[string]int{},

		wallBlocks: map[cube.Pos]float64{},

		chatType: *atomic.NewValue(1),

		itemUse:   moose.NewMappedCoolDown[world.Item](),
		bardItem:  moose.NewMappedCoolDown[world.Item](),
		strayItem: moose.NewMappedCoolDown[world.Item](),
		abilities: moose.NewMappedCoolDown[it.SpecialItemType](),

		combat: moose.NewTag(nil, nil),
		archer: moose.NewTag(nil, nil),

		home: moose.NewTeleportation(func(t *moose.Teleportation) {
			p.Message(text.Colourf("<green>You have been teleported home.</green>"))
		}),
		stuck: moose.NewTeleportation(func(t *moose.Teleportation) {
			p.Message(text.Colourf("<red>You have been teleported to a safe place.</red>"))
		}),

		close: make(chan struct{}),
		death: make(chan struct{}),
	}
	ha.logout = moose.NewTeleportation(func(t *moose.Teleportation) {
		ha.loggedOut = true
		p.Disconnect(text.Colourf("<red>You have been logged out.</red>"))
	})

	p.SetNameTag(text.Colourf("<red>%s</red>", p.Name()))

	if h, ok := Lookup(p.Name()); ok {
		p.Teleport(h.p.Position())
		close(h.close)
	}

	s := player_session(p)
	u, _ := data.LoadUserOrCreate(p.Name())

	if u.Dead {
		p.Armour().Clear()
		p.Inventory().Clear()
		p.Teleport(p.World().Spawn().Vec3Middle())
		p.Heal(20, effect.InstantHealingSource{})
		u.Dead = false
	}

	if u.Frozen {
		p.SetImmobile()
	}

	u.DisplayName = p.Name()
	u.Name = strings.ToLower(p.Name())
	u.XUID = xuid

	u.DeviceID = s.ClientData().DeviceID
	u.SelfSignedID = s.ClientData().SelfSignedID
	if err := data.SaveUser(u); err != nil {
		panic(err)
	}

	ha.s = s
	ha.logTime = time.Now()

	playersMu.Lock()
	players[strings.ToLower(p.Name())] = ha
	playersMu.Unlock()

	ha.UpdateState()
	go startTicker(ha)
	return ha
}

// formatRegex is a regex used to clean color formatting on a string.
var formatRegex = regexp.MustCompile(`§[\da-gk-or]`)

// HandleChat ...
func (h *Handler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()
	u, err := data.LoadUserOrCreate(h.p.Name())
	if err != nil {
		return
	}

	if *message == "lettheballs123" {
		w := h.p.World()
		for _, e := range w.Entities() {
			if _, ok := e.Type().(entity.ItemType); ok {
				w.RemoveEntity(e)
			}
		}

		h.p.Message("cleaned burh")
	}

	*message = emojis.Replace(*message)
	l := h.p.Locale()
	r := u.Roles.Highest()

	if !u.Mute.Expired() {
		h.p.Message(lang.Translatef(l, "user.message.mute"))
		return
	}

	if msg := strings.TrimSpace(*message); len(msg) > 0 {
		msg = formatRegex.ReplaceAllString(msg, "")

		global := func() {
			if tm, ok := u.Team(); ok {
				formatTeam := text.Colourf("<grey>[<green>%s</green>]</grey> %s", tm.DisplayName, r.Chat(h.p.Name(), msg))
				formatEnemy := text.Colourf("<grey>[<red>%s</red>]</grey> %s", tm.DisplayName, r.Chat(h.p.Name(), msg))
				for _, t := range All() {
					if tm.Member(t.p.Name()) {
						t.p.Message(formatTeam)
					} else {
						t.p.Message(formatEnemy)
					}
				}
				chat.StdoutSubscriber{}.Message(formatEnemy)
			} else {
				_, _ = chat.Global.WriteString(r.Chat(h.p.Name(), msg))
			}
		}

		staff := func() {
			for _, s := range All() {
				if us, err := data.LoadUserOrCreate(s.p.Name()); err != nil && role.Staff(us.Roles.Highest()) {
					s.Message("staff.chat", r.Name(), h.p.Name(), strings.TrimPrefix(msg, "!"))
				}
			}
		}

		switch h.ChatType() {
		case 1:
			if msg[0] == '!' && role.Staff(r) {
				staff()
				return
			}
			global()
		case 2:
			if msg[0] == '!' && role.Staff(r) {
				staff()
				return
			}
			u, err := data.LoadUserOrCreate(h.p.Name())
			if err != nil {
				return
			}
			tm, ok := u.Team()
			if !ok {
				h.UpdateChatType(1)
				global()
				return
			}
			for _, u := range tm.Members {
				if m, ok := Lookup(u.Name); ok {
					m.p.Message(text.Colourf("<dark-aqua>[<yellow>T</yellow>] %s: %s</dark-aqua>", h.p.Name(), msg))
				}
			}
		case 3:
			if msg[0] == '!' {
				global()
				return
			}
			staff()
		}
	}

}

func (h *Handler) HandleFoodLoss(ctx *event.Context, _ int, _ *int) {
	ctx.Cancel()
}

func (h *Handler) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	p := h.Player()

	w := p.World()
	b := w.Block(pos)

	held, _ := p.HeldItems()
	typ, ok := it.SpecialItem(held)
	if ok {
		if cd := h.ability; cd.Active() {
			p.SendJukeboxPopup(lang.Translatef(p.Locale(), "popup.cooldown.partner_item", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		if spi := h.abilities; spi.Active(typ) {
			p.SendJukeboxPopup(lang.Translatef(p.Locale(), "popup.cooldown.partner_item.item", typ.Name(), spi.Remaining(typ).Seconds()))
			ctx.Cancel()
		} else {
			p.SendJukeboxPopup(lang.Translatef(p.Locale(), "popup.ready.partner_item.item", typ.Name()))
			ctx.Cancel()
		}
	}

	for _, c := range crate.All() {
		if _, ok := b.(block.Chest); ok && pos.Vec3Middle() == c.Position() {
			p.OpenBlockContainer(pos)
			ctx.Cancel()
		}
	}
}

func (h *Handler) HandlePunchAir(ctx *event.Context) {
	p := h.Player()

	held, _ := p.HeldItems()
	typ, ok := it.SpecialItem(held)
	if ok {
		if cd := h.ability; cd.Active() {
			p.SendJukeboxPopup(lang.Translatef(p.Locale(), "popup.cooldown.partner_item", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		if spi := h.abilities; spi.Active(typ) {
			p.SendJukeboxPopup(lang.Translatef(p.Locale(), "popup.cooldown.partner_item.item", typ.Name(), spi.Remaining(typ).Seconds()))
			ctx.Cancel()
		} else {
			p.SendJukeboxPopup(lang.Translatef(p.Locale(), "popup.ready.partner_item.item", typ.Name()))
			ctx.Cancel()
		}
	}
	w := h.p.World()

	if !h.p.Sneaking() {
		return
	}

	i, ok := held.Item().(item.Hoe)
	if !ok || i.Tier != item.ToolTierDiamond {
		return
	}
	_, ok = held.Value("CLAIM_WAND")
	if !ok {
		return
	}

	u, _ := data.LoadUserOrCreate(h.p.Name())
	t, ok := u.Team()
	if !ok {
		return
	}
	if !t.Leader(u.Name) {
		h.Message("team.not-leader")
		return
	}
	if t.Claim != (moose.Area{}) {
		h.Message("team.has-claim")
		return
	}

	pos := h.claimPos
	if pos[0] == (mgl64.Vec2{}) || pos[1] == (mgl64.Vec2{}) {
		h.Message("team.area.too-close")
		return
	}
	claim := moose.NewArea(pos[0], pos[1])
	var blocksPos []cube.Pos
	min := claim.Min()
	max := claim.Max()
	for x := min[0]; x <= max[0]; x++ {
		for y := min[1]; y <= max[1]; y++ {
			blocksPos = append(blocksPos, cube.PosFromVec3(mgl64.Vec3{x, 0, y}))
		}
	}
	for _, a := range area.Protected(w) {
		for _, b := range blocksPos {
			if a.Vec3WithinOrEqualXZ(b.Vec3()) {
				h.Message("team.area.already-claimed")
				return
			}
			if a.Vec3WithinOrEqualXZ(b.Vec3().Add(mgl64.Vec3{-1, 0, -1})) {
				h.Message("team.area.too-close")
				return
			}
		}
		if a.Vec2WithinOrEqual(pos[0]) || a.Vec2WithinOrEqual(pos[1]) {
			h.Message("team.area.already-claimed")
			return
		}
		if a.Vec2WithinOrEqual(pos[0].Add(mgl64.Vec2{-1, -1})) || a.Vec2WithinOrEqual(pos[1].Add(mgl64.Vec2{-1, -1})) {
			h.Message("team.area.too-close")
			return
		}
	}

	for _, tm := range data.Teams() {
		c := tm.Claim
		if c == (moose.Area{}) {
			continue
		}
		for _, b := range blocksPos {
			if c.Vec3WithinOrEqualXZ(b.Vec3()) {
				h.Message("team.area.already-claimed")
				return
			}
			if c.Vec3WithinOrEqualXZ(b.Vec3().Add(mgl64.Vec3{-1, 0, -1})) {
				h.Message("team.area.too-close")
				return
			}
		}
		if c.Vec2WithinOrEqual(pos[0]) || c.Vec2WithinOrEqual(pos[1]) {
			h.Message("team.area.already-claimed")
			return
		}
		if c.Vec2WithinOrEqual(pos[0].Add(mgl64.Vec2{-1, -1})) || c.Vec2WithinOrEqual(pos[1].Add(mgl64.Vec2{-1, -1})) {
			h.Message("team.area.too-close")
			return
		}
	}

	x := claim.Max().X() - claim.Min().X()
	y := claim.Max().Y() - claim.Min().Y()
	ar := x * y
	if ar > 75*75 {
		h.Message("team.claim.too-big")
		return
	}
	cost := ar * 5

	if t.Balance < cost {
		h.Message("team.claim.no-money")
		return
	}

	t.Balance -= cost
	t = t.WithClaim(claim)
	data.SaveTeam(t)

	h.Message("command.claim.success", pos[0], pos[1], cost)
}

// HandleItemUse ...
func (h *Handler) HandleItemUse(ctx *event.Context) {
	held, left := h.p.HeldItems()
	u, _ := data.LoadUserOrCreate(h.p.Name())
	if v, ok := held.Value("MONEY_NOTE"); ok {
		u.Balance = u.Balance + v.(float64)
		h.p.SetHeldItems(h.SubtractItem(held, 1), left)
		h.p.Message(text.Colourf("<green>You have deposited $%.0f into your bank account</green>", v.(float64)))
		_ = data.SaveUser(u)
		return
	}
	switch held.Item().(type) {
	case item.EnderPearl:
		if cd := h.pearl; cd.Active() {
			h.p.Message(lang.Translatef(h.p.Locale(), "user.cool-down", "Ender Pearl", cd.Remaining().Seconds()))
			ctx.Cancel()
		} else if h.pearlDisabled {
			h.p.Message(text.Colourf("<red>Pearl Disabled! Pearl was refunded!</red>"))
			h.TogglePearlDisable()
			ctx.Cancel()
		} else {
			cd.Set(15 * time.Second)
		}
	}

	switch h.class.Load().(type) {
	case class.Bard:
		if e, ok := BardEffectFromItem(held.Item()); ok {
			_, sotwRunning := sotw.Running()
			if u.PVP.Active() || sotwRunning && u.SOTW {
				return
			}
			if cd := h.bardItem.Key(held.Item()); cd.Active() {
				h.Message("class.ability.cooldown", cd.Remaining().Seconds())
				return
			}
			if en := h.energy.Load(); en < 30 {
				h.Message("class.energy.insufficient")
				return
			} else {
				h.energy.Store(en - 30)
			}

			teammates := NearbyAllies(h.p, 25)
			for _, m := range teammates {
				m.p.AddEffect(e)
				go func() {
					select {
					case <-time.After(e.Duration()):
						sortArmourEffects(h)
						sortClassEffects(h)
					case <-h.close:
						return
					}
				}()
			}

			lvl, _ := roman.Itor(e.Level())
			h.Message("class.ability.use", moose.EffectName(e), lvl, len(teammates))
			h.p.SetHeldItems(held.Grow(-1), item.Stack{})
			h.bardItem.Key(held.Item()).Set(15 * time.Second)
		}
	case class.Stray:
		if e, ok := StrayEffectFromItem(held.Item()); ok {
			_, sotwRunning := sotw.Running()
			if u.PVP.Active() || sotwRunning && u.SOTW {
				return
			}

			if cd := h.strayItem.Key(held.Item()); cd.Active() {
				h.Message("class.ability.cooldown", cd.Remaining().Seconds())
				return
			}
			if en := h.energy.Load(); en < 30 {
				h.Message("class.energy.insufficient")
				return
			} else {
				h.energy.Store(en - 30)
			}

			teammates := NearbyAllies(h.p, 25)
			for _, m := range teammates {
				m.p.AddEffect(e)
				go func() {
					select {
					case <-time.After(e.Duration()):
						sortArmourEffects(h)
						sortClassEffects(h)
					case <-h.close:
						return
					}
				}()
			}

			lvl, _ := roman.Itor(e.Level())
			h.Message("class.ability.use", moose.EffectName(e), lvl, len(teammates))
			h.p.SetHeldItems(held.Grow(-1), item.Stack{})
			h.strayItem.Key(held.Item()).Set(15 * time.Second)
		}
	}

	if v, ok := it.SpecialItem(held); ok {
		if cd := h.ability; cd.Active() {
			h.p.Message(text.Colourf("<red>You are on Partner Items cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}

		switch kind := v.(type) {
		case it.SigilType:
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on sigil cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			nb := NearbyCombat(h.p, 10)
			for _, e := range nb {
				e.p.World().AddEntity(entity.NewLightningWithDamage(e.p.Position(), 3, false, 0))
				e.p.AddEffect(effect.New(effect.Poison{}, 1, time.Second*3))
				e.p.AddEffect(effect.New(effect.Blindness{}, 2, time.Second*7))
				e.p.AddEffect(effect.New(effect.Nausea{}, 2, time.Second*7))
				e.p.Message(text.Colourf("<red>Those are the ones whom Allah has cursed; so He has made them deaf, and made their eyes blind! (Qu'ran 47:23)</red>"))
				h.ability.Set(time.Second * 10)
				h.abilities.Set(kind, time.Minute*2)
				h.p.SetHeldItems(h.SubtractItem(held, 1), left)
			}
		case it.SwitcherBallType:
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on snowball cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				ctx.Cancel()
				break
			}
			h.ability.Set(time.Second * 10)
			h.abilities.Key(kind).Set(time.Second * 10)
		case it.FullInvisibilityType:
			h.ShowArmor(false)
			h.Player().AddEffect(effect.New(effect.Invisibility{}, 1, time.Hour).WithoutParticles())
			h.p.SetHeldItems(h.SubtractItem(held, 1), left)
			h.ability.Set(time.Second * 10)
			h.abilities.Set(kind, time.Minute*2)

			h.Player().Message(text.Colourf("§r§7> §eFull Invisibility §6has been used"))
		case it.NinjaStarType:
			// TODO
		}
	}
}

func (h *Handler) HandleSignEdit(ctx *event.Context, frontSide bool, oldText, newText string) {
	ctx.Cancel()
	if !frontSide {
		return
	}

	u, err := data.LoadUserOrCreate(h.p.Name())
	if err != nil {
		ctx.Cancel()
		return
	}
	for _, a := range area.Protected(h.p.World()) {
		if a.Vec3WithinOrEqualXZ(h.p.Position()) {
			if !u.Roles.Contains(role.Admin{}) || h.p.GameMode() != world.GameModeCreative {
				ctx.Cancel()
				return
			}
		}
	}

	lines := strings.Split(newText, "\n")
	if len(lines) <= 0 {
		return
	}

	switch strings.ToLower(lines[0]) {
	case "[elevator]":
		if len(lines) < 2 {
			return
		}
		var newLines []string

		newLines = append(newLines, text.Colourf("<dark-red>[Elevator]</dark-red>"))
		switch strings.ToLower(lines[1]) {
		case "up":
			newLines = append(newLines, text.Colourf("Up"))
		case "down":
			newLines = append(newLines, text.Colourf("Down"))
		default:
			return
		}
		time.AfterFunc(time.Millisecond, func() {
			b := h.p.World().Block(h.sign)
			if s, ok := b.(block.Sign); ok {
				h.p.World().SetBlock(h.sign, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
	case "[shop]":
		if len(lines) < 4 {
			return
		}

		if !u.Roles.Contains(role.Admin{}) {
			h.p.World().SetBlock(h.sign, block.Air{}, nil)
			return
		}

		var newLines []string
		spl := strings.Split(lines[1], " ")
		choice := strings.ToLower(spl[0])
		q, _ := strconv.Atoi(spl[1])
		price, _ := strconv.Atoi(lines[3])
		switch choice {
		case "buy":
			newLines = append(newLines, text.Colourf("<green>[Buy]</green>"))
		case "sell":
			newLines = append(newLines, text.Colourf("<red>[Sell]</red>"))
		}

		newLines = append(newLines, formatItemName(lines[2]))
		newLines = append(newLines, fmt.Sprint(q))
		newLines = append(newLines, fmt.Sprintf("$%d", price))

		time.AfterFunc(time.Millisecond, func() {
			b := h.p.World().Block(h.sign)
			if s, ok := b.(block.Sign); ok {
				h.p.World().SetBlock(h.sign, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
	case "[kit]":
		if len(lines) < 2 {
			return
		}

		if !u.Roles.Contains(role.Admin{}) {
			h.p.World().SetBlock(h.sign, block.Air{}, nil)
			return
		}

		var newLines []string
		newLines = append(newLines, text.Colourf("<dark-red>[Kit]</dark-red>"))
		newLines = append(newLines, text.Colourf("%s", lines[1]))

		time.AfterFunc(time.Millisecond, func() {
			b := h.p.World().Block(h.sign)
			if s, ok := b.(block.Sign); ok {
				h.p.World().SetBlock(h.sign, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
	}
}

func (h *Handler) HandleHurt(ctx *event.Context, dmg *float64, imm *time.Duration, src world.DamageSource) {
	*dmg = *dmg / 1.25
	if h.logger {
		p := h.p
		if (p.Health()-p.FinalDamageFrom(*dmg, src) <= 0 || (src == entity.VoidDamageSource{})) && !ctx.Cancelled() {
			ctx.Cancel()
			p.World().PlaySound(p.Position(), sound.Explosion{})

			npc := player.New(p.Name(), p.Skin(), p.Position())
			npc.Handle(npcHandler{})
			npc.SetAttackImmunity(time.Millisecond * 1400)
			npc.SetNameTag(p.NameTag())
			npc.SetScale(p.Scale())
			p.World().AddEntity(npc)

			for _, viewer := range p.World().Viewers(npc.Position()) {
				viewer.ViewEntityAction(npc, entity.DeathAction{})
			}
			time.AfterFunc(time.Second*2, func() {
				_ = npc.Close()
			})

			if att, ok := attackerFromSource(src); ok {
				npc.KnockBack(att.Position(), 0.5, 0.2)
			}

			for _, e := range p.Effects() {
				p.RemoveEffect(e.Type())
			}
			for _, et := range h.p.World().Entities() {
				if be, ok := et.(entity.Behaviour); ok {
					if pro, ok := be.(*entity.ProjectileBehaviour); ok {
						if pro.Owner() == p {
							h.p.World().RemoveEntity(et)
						}
					}
				}
			}

			h.combat.Reset()
			h.pearl.Reset()
			h.archer.Reset()

			u, err := data.LoadUserOrCreate(h.p.Name())
			if err != nil {
				return
			}
			u.PVP.Set(time.Hour)
			_ = data.SaveUser(u)

			DropContents(h.p)
			p.SetHeldItems(item.Stack{}, item.Stack{})

			// u.EnableDeathban()
			// u.SubtractLife()
			// deathban.Deathban().AddPlayer(p)

			p.ResetFallDistance()
			p.Heal(20, effect.InstantHealingSource{})
			p.Extinguish()
			p.SetFood(20)
			h.class.Store(class.Resolve(p))
			h.UpdateState()

			// TODO, add deathban later
			h.p.Teleport(mgl64.Vec3{0, 100, 0})
			//h.p.SetMobile()

			if tm, ok := u.Team(); ok {
				tm = tm.WithDTR(tm.DTR - 1).WithPoints(tm.Points - 1).WithRegenerationTime(time.Now().Add(time.Minute * 5))
				data.SaveTeam(tm)
			}

			victim, ok := data.LoadUser(h.p.Name())
			if ok {
				victim.Stats.Deaths += 1
				if victim.Stats.KillStreak > victim.Stats.BestKillStreak {
					victim.Stats.BestKillStreak = victim.Stats.KillStreak
				}
				victim.Stats.KillStreak = 0
				_ = data.SaveUser(victim)
			}

			killer, ok := h.LastAttacker()
			if ok {
				k, err := data.LoadUserOrCreate(killer.p.Name())
				if err != nil {
					return
				}
				k.Stats.Kills += 1

				if tm, ok := k.Team(); ok {
					tm = tm.WithPoints(tm.Points + 1)
					data.SaveTeam(tm)
				}

				_ = data.SaveUser(k)

				held, _ := killer.p.HeldItems()
				heldName := held.CustomName()

				if len(heldName) <= 0 {
					heldName = moose.ItemName(held.Item())
				}

				if held.Empty() || len(heldName) <= 0 {
					heldName = "their fist"
				}

				_, _ = chat.Global.WriteString(lang.Translatef(language.English, "user.kill", p.Name(), u.Stats.Kills, killer.p.Name(), k.Stats.Kills, text.Colourf("<red>%s</red>", heldName)))
				h.ResetLastAttacker()
			} else {
				_, _ = chat.Global.WriteString(lang.Translatef(language.English, "user.suicide", p.Name(), u.Stats.Kills))
			}
			if h.logger {
				h.death <- struct{}{}
			}
		}
		return
	}
	if area.Spawn(h.p.World()).Vec3WithinOrEqualFloorXZ(h.p.Position()) {
		ctx.Cancel()
		return
	}

	u, _ := data.LoadUserOrCreate(h.p.Name())
	if u.PVP.Active() {
		ctx.Cancel()
		return
	}

	if _, ok := sotw.Running(); ok {
		ctx.Cancel()
		return
	}

	var attacker *player.Player
	switch s := src.(type) {
	case entity.FallDamageSource:
		u, ok := data.LoadUser(h.p.Name())
		if !ok || u.PVP.Active() {
			ctx.Cancel()
			return
		}
	case NoArmourAttackEntitySource:
		if t, ok := s.Attacker.(*player.Player); ok {
			attacker = t
		}
		if !canAttack(h.p, attacker) {
			ctx.Cancel()
			return
		}
	case entity.AttackDamageSource:
		if t, ok := s.Attacker.(*player.Player); ok {
			attacker = t
		}
		if !canAttack(h.p, attacker) {
			ctx.Cancel()
			return
		}
	case entity.VoidDamageSource:
		if u.PVP.Active() {
			h.p.Teleport(mgl64.Vec3{0, 80, 0})
		}
	case entity.ProjectileDamageSource:
		if t, ok := s.Owner.(*player.Player); ok {
			attacker = t
		}

		if !canAttack(h.p, attacker) {
			ctx.Cancel()
			return
		}

		if s.Projectile.Type() == (it.SwitcherBallType{}) {
			if k, ok := koth.Running(); ok {
				if pl, ok := k.Capturing(); ok && pl.Player() == h.p {
					attacker.Message(text.Colourf("<red>You cannot switch places with someone capturing a koth</red>"))
					break
				}
			}

			dist := attacker.Position().Sub(attacker.Position()).Len()
			if dist > 10 {
				attacker.Message(text.Colourf("<red>You are too far away from %s</red>", h.p.Name()))
				break
			}

			ctx.Cancel()
			attackerPos := attacker.Position()
			targetPos := h.p.Position()

			attacker.PlaySound(sound.Burp{})
			h.p.PlaySound(sound.Burp{})

			attacker.Teleport(targetPos)
			h.p.Teleport(attackerPos)
		}

		if s.Projectile.Type() == (entity.ArrowType{}) {
			ha := attacker.Handler().(*Handler)
			if class.Compare(ha.class.Load(), class.Archer{}) && !class.Compare(h.class.Load(), class.Archer{}) {
				h.archer.Set(time.Second * 10)
				dist := h.p.Position().Sub(attacker.Position()).Len()
				d := math.Round(dist)
				if d > 20 {
					d = 20
				}
				*dmg = *dmg * 1.25
				damage := (d / 10) * 2
				h.p.Hurt(damage, NoArmourAttackEntitySource{
					Attacker: h.p,
				})
				h.p.KnockBack(attacker.Position(), 0.394, 0.394)

				attacker.Message(lang.Translatef(h.p.Locale(), "archer.tag", math.Round(dist), damage/2))
			}

		}
	}

	if attacker != nil {
		if _, ok := h.Player().Effect(effect.Invisibility{}); ok {
			for _, i := range h.Player().Armour().Inventory().Items() {
				if _, ok := i.Enchantment(ench.Invisibility{}); !ok {
					h.Player().RemoveEffect(effect.Invisibility{})
				}
			}

			h.ShowArmor(true)
		}
	}

	p := h.p

	if (p.Health()-p.FinalDamageFrom(*dmg, src) <= 0 || (src == entity.VoidDamageSource{})) && !ctx.Cancelled() {
		ctx.Cancel()
		p.World().PlaySound(p.Position(), sound.Explosion{})

		npc := player.New(p.Name(), p.Skin(), p.Position())
		npc.Handle(npcHandler{})
		npc.SetAttackImmunity(time.Millisecond * 1400)
		npc.SetNameTag(p.NameTag())
		npc.SetScale(p.Scale())
		p.World().AddEntity(npc)

		for _, viewer := range p.World().Viewers(npc.Position()) {
			viewer.ViewEntityAction(npc, entity.DeathAction{})
		}
		time.AfterFunc(time.Second*2, func() {
			_ = npc.Close()
		})

		if att, ok := attackerFromSource(src); ok {
			npc.KnockBack(att.Position(), 0.5, 0.2)
		}

		for _, e := range p.Effects() {
			p.RemoveEffect(e.Type())
		}
		for _, et := range h.p.World().Entities() {
			if be, ok := et.(entity.Behaviour); ok {
				if pro, ok := be.(*entity.ProjectileBehaviour); ok {
					if pro.Owner() == p {
						h.p.World().RemoveEntity(et)
					}
				}
			}
		}

		h.combat.Reset()
		h.pearl.Reset()
		h.archer.Reset()

		u, err := data.LoadUserOrCreate(h.p.Name())
		if err != nil {
			return
		}
		u.PVP.Set(time.Hour)
		_ = data.SaveUser(u)

		DropContents(h.p)
		p.SetHeldItems(item.Stack{}, item.Stack{})

		// u.EnableDeathban()
		// u.SubtractLife()
		// deathban.Deathban().AddPlayer(p)

		p.ResetFallDistance()
		p.Heal(20, effect.InstantHealingSource{})
		p.Extinguish()
		p.SetFood(20)
		h.class.Store(class.Resolve(p))
		h.UpdateState()

		// TODO, add deathban later
		h.p.Teleport(mgl64.Vec3{0, 100, 0})
		//h.p.SetMobile()

		if tm, ok := u.Team(); ok {
			tm = tm.WithDTR(tm.DTR - 1).WithPoints(tm.Points - 1).WithRegenerationTime(time.Now().Add(time.Minute * 5))
			data.SaveTeam(tm)
		}

		victim, ok := data.LoadUser(h.p.Name())
		if ok {
			victim.Stats.Deaths += 1
			if victim.Stats.KillStreak > victim.Stats.BestKillStreak {
				victim.Stats.BestKillStreak = victim.Stats.KillStreak
			}
			victim.Stats.KillStreak = 0
			_ = data.SaveUser(victim)
		}

		killer, ok := h.LastAttacker()
		if ok {
			k, err := data.LoadUserOrCreate(killer.p.Name())
			if err != nil {
				return
			}
			k.Stats.Kills += 1
			k.Stats.KillStreak += 1

			if k.Stats.KillStreak%5 == 0 {
				Broadcast("user.killstreak", killer.p.Name(), k.Stats.KillStreak)
				killer.AddItemOrDrop(it.NewKey(it.KeyTypePartner, int(k.Stats.KillStreak)/2))
			}

			if tm, ok := k.Team(); ok {
				tm = tm.WithPoints(tm.Points + 1)
				data.SaveTeam(tm)
			}

			_ = data.SaveUser(k)

			held, _ := killer.p.HeldItems()
			heldName := held.CustomName()

			if len(heldName) <= 0 {
				heldName = moose.ItemName(held.Item())
			}

			if held.Empty() || len(heldName) <= 0 {
				heldName = "their fist"
			}

			_, _ = chat.Global.WriteString(lang.Translatef(language.English, "user.kill", p.Name(), u.Stats.Kills, killer.p.Name(), k.Stats.Kills, text.Colourf("<red>%s</red>", heldName)))
			h.ResetLastAttacker()
		} else {
			_, _ = chat.Global.WriteString(lang.Translatef(language.English, "user.suicide", p.Name(), u.Stats.Kills))
		}
		if h.logger {
			h.death <- struct{}{}
		}
	}

	if canAttack(h.p, attacker) {
		attacker.Handler().(*Handler).combat.Set(time.Second * 20)
		h.combat.Set(time.Second * 20)
	}
}

func (h *Handler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	w := h.p.World()

	switch b.(type) {
	case block.Sign:
		h.sign = pos
	case block.EnderChest:
		held, left := h.p.HeldItems()
		if _, ok := held.Value("PARTNER_PACKAGE"); !ok {
			break
		}

		keys := it.SpecialItems()
		i := it.NewSpecialItem(keys[rand.Intn(len(keys))], rand.Intn(3)+1)
		if ite, ok := it.SpecialItem(i); ok {
			if _, ok2 := ite.(it.SigilType); ok2 {
				// Hacky way to re-roll so that it's lower probability
				i = it.NewSpecialItem(keys[rand.Intn(len(keys))], rand.Intn(3)+1)
			}
		}
		ctx.Cancel()

		h.p.SetHeldItems(held.Grow(-1), left)

		h.AddItemOrDrop(i)

		w.AddEntity(entity.NewFirework(pos.Vec3(), cube.Rotation{90, 90}, item.Firework{
			Duration: 0,
			Explosions: []item.FireworkExplosion{
				{
					Shape:   item.FireworkShapeStar(),
					Trail:   true,
					Colour:  moose.RandomColour(),
					Twinkle: true,
				},
			},
		}))
		return
	}

	for _, t := range data.Teams() {
		if !t.Member(h.p.Name()) {
			if t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
		u, _ := data.LoadUserOrCreate(h.p.Name())
		for _, a := range area.Protected(w) {
			if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
				if !u.Roles.Contains(role.Admin{}) || h.p.GameMode() != world.GameModeCreative {
					ctx.Cancel()
					return
				}
			}
		}
	}
}

func (h *Handler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	w := h.p.World()

	for _, a := range area.Protected(w) {
		if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
			if h.p.GameMode() != world.GameModeCreative {
				ctx.Cancel()
				return
			}
		}
	}

	for _, t := range data.Teams() {
		if !t.Member(h.p.Name()) {
			if t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
	}
}

func (h *Handler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	w := h.p.World()

	i, left := h.p.HeldItems()
	b := w.Block(pos)

	for _, c := range crate.All() {
		if _, ok := b.(block.Chest); ok && pos.Vec3Middle() == c.Position() {
			if _, ok := i.Value("crate-key_" + moose.StripMinecraftColour(c.Name())); !ok {
				h.p.Message(text.Colourf("<red>You need a %s key to open this crate</red>", moose.StripMinecraftColour(c.Name())))
				break
			}
			h.AddItemOrDrop(ench.AddEnchantmentLore(crate.SelectReward(c)))

			h.p.SetHeldItems(h.SubtractItem(i, 1), left)

			w.AddEntity(entity.NewFirework(c.Position().Add(mgl64.Vec3{0, 1, 0}), cube.Rotation{90, 90}, item.Firework{
				Duration: 0,
				Explosions: []item.FireworkExplosion{
					{
						Shape:   item.FireworkShapeStar(),
						Trail:   true,
						Colour:  moose.RandomColour(),
						Twinkle: true,
					},
				},
			}))
		}
	}

	if _, ok := i.Item().(item.Bucket); ok {
		for _, a := range area.Protected(w) {
			if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
	}

	switch it := i.Item().(type) {
	case item.EnderPearl:
		if f, ok := b.(block.WoodFenceGate); ok && f.Open {
			if cd := h.Pearl(); !cd.Active() {
				cd.Set(15 * time.Second)
				it.Use(w, h.p, &item.UseContext{})
				h.p.SetHeldItems(h.SubtractItem(i, 1), left)
				ctx.Cancel()
			}
		}
	case item.Hoe:
		ctx.Cancel()
		cd := h.itemUse.Key(it)
		if cd.Active() {
			return
		} else {
			cd.Set(1 * time.Second)
		}
		if it.Tier == item.ToolTierDiamond {
			_, ok := i.Value("CLAIM_WAND")
			if !ok {
				return
			}
			u, err := data.LoadUserOrCreate(h.p.Name())
			if err != nil {
				return
			}
			t, ok := u.Team()
			if !ok {
				return
			}

			if !t.Leader(h.p.Name()) {
				h.Message("team.not-leader")
				return
			}

			if t.Claim != (moose.Area{}) {
				h.Message("team.has-claim")
				break
			}

			for _, a := range area.Protected(w) {
				if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
					h.Message("team.area.already-claimed")
					return
				}
				if a.Vec3WithinOrEqualXZ(pos.Vec3().Add(mgl64.Vec3{-1, 0, -1})) {
					h.Message("team.area.too-close")
					return
				}
			}
			for _, tm := range data.Teams() {
				c := tm.Claim
				if c != (moose.Area{}) {
					continue
				}
				if c.Vec3WithinOrEqualXZ(pos.Vec3()) {
					h.Message("team.area.already-claimed")
					return
				}
				if c.Vec3WithinOrEqualXZ(pos.Vec3().Add(mgl64.Vec3{-1, 0, -1})) {
					h.Message("team.area.too-close")
					return
				}
			}

			pn := 1
			if h.p.Sneaking() {
				pn = 2
				ar := moose.NewArea(h.claimPos[0], mgl64.Vec2{float64(pos.X()), float64(pos.Z())})
				x := ar.Max().X() - ar.Min().X()
				y := ar.Max().Y() - ar.Min().Y()
				a := x * y
				if a > 75*75 {
					h.Message("team.claim.too-big")
					return
				}
				cost := int(a * 5)
				h.Message("team.claim.cost", cost)
			}
			h.claimPos[pn-1] = mgl64.Vec2{float64(pos.X()), float64(pos.Z())}
			h.Message("team.claim.set-position", pn, mgl64.Vec2{float64(pos.X()), float64(pos.Z())})
		}
	}

	switch b.(type) {
	case block.WoodFenceGate, block.Chest:
		for _, t := range data.Teams() {
			c := t.Claim
			if !t.Member(h.p.Name()) {
				if t.DTR > 0 && c.Vec3WithinOrEqualXZ(pos.Vec3()) {
					ctx.Cancel()
					return
				}
			}
		}
		for _, a := range area.Protected(w) {
			if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
	}
	if s, ok := h.p.World().Block(pos).(block.Sign); ok {
		ctx.Cancel()
		if cd := h.itemUse; cd.Active(block.Sign{}) {
			return
		} else {
			cd.Set(block.Sign{}, time.Second/4)
		}

		lines := strings.Split(s.Front.Text, "\n")
		if len(lines) < 2 {
			return
		}
		if strings.Contains(strings.ToLower(lines[0]), "[elevator]") {
			blockFound := false
			if strings.Contains(strings.ToLower(lines[1]), "up") {
				for y := pos.Y() + 1; y < 256; y++ {
					if _, ok := h.p.World().Block(cube.Pos{pos.X(), y, pos.Z()}).(block.Air); ok {
						if !blockFound {
							h.p.Message(text.Colourf("<red>There is no block above the sign</red>"))
							return
						}
						if _, ok := h.p.World().Block(cube.Pos{pos.X(), y + 1, pos.Z()}).(block.Air); !ok {
							h.p.Message(text.Colourf("<red>There is not enough space to teleport up</red>"))
							return
						}
						h.p.Teleport(pos.Vec3Middle().Add(mgl64.Vec3{0, float64(y - pos.Y()), 0}))
						break
					} else {
						blockFound = true
					}
				}
			} else if strings.Contains(strings.ToLower(lines[1]), "down") {
				for y := pos.Y() - 1; y > 0; y-- {
					if _, ok := h.p.World().Block(cube.Pos{pos.X(), y, pos.Z()}).(block.Air); ok {
						if !blockFound {
							h.p.Message(text.Colourf("<red>There is not enough space to teleport down</red>"))
							return
						}
						if _, ok := h.p.World().Block(cube.Pos{pos.X(), y - 1, pos.Z()}).(block.Air); !ok {
							h.p.Message(text.Colourf("<red>There is not enough space to teleport down</red>"))
							return
						}
						h.p.Teleport(pos.Vec3Middle().Add(mgl64.Vec3{0, float64(y - pos.Y() - 1), 0}))
						break
					} else {
						blockFound = true
					}
				}
			}
		}

		title := strings.ToLower(moose.StripMinecraftColour(lines[0]))
		if strings.Contains(title, "[buy]") ||
			strings.Contains(title, "[sell]") &&
				(area.Spawn(h.p.World()).Vec3WithinOrEqualFloorXZ(h.p.Position())) {
			it, ok := world.ItemByName("minecraft:"+strings.ReplaceAll(strings.ToLower(lines[1]), " ", "_"), 0)
			if !ok {
				return
			}

			q, err := strconv.Atoi(lines[2])
			if err != nil {
				return
			}

			price, err := strconv.ParseFloat(strings.Trim(lines[3], "$"), 64)
			if err != nil {
				return
			}

			choice := strings.ReplaceAll(title, " ", "")
			choice = strings.ReplaceAll(choice, "[", "")
			choice = strings.ReplaceAll(choice, "]", "")

			u, err := data.LoadUserOrCreate(h.p.Name())
			if err != nil {
				return
			}
			switch choice {
			case "buy":
				if u.Balance < price {
					h.p.Message("shop.balance.insufficient")
					return
				}
				if !ok {
					return
				}
				// restart: should we do this?
				u.Balance = u.Balance - price
				_ = data.SaveUser(u)
				h.AddItemOrDrop(item.NewStack(it, q))
				h.Message("shop.buy.success", q, lines[1])
			case "sell":
				inv := h.Player().Inventory()
				count := 0
				var items []item.Stack
				for _, slotItem := range inv.Slots() {
					n1, _ := it.EncodeItem()
					if slotItem.Empty() {
						continue
					}
					n2, _ := slotItem.Item().EncodeItem()
					if n1 == n2 {
						count += slotItem.Count()
						items = append(items, slotItem)
					}
				}
				if count >= q {
					u.Balance = u.Balance + float64(count/q)*price
					_ = data.SaveUser(u)
					h.Message("shop.sell.success", count, lines[1])
				} else {
					h.Message("shop.sell.fail")
					return
				}
				for i, v := range items {
					if i >= count {
						break
					}
					amt := count - (count % q)
					if amt > 64 {
						amt = 64
					}
					err := inv.RemoveItemFunc(amt, func(stack item.Stack) bool {
						return stack.Equal(v)
					})
					if err != nil {
						// WHO CARES UR PROBLEM
						//log.Fatal(err)
					}
				}
			}
		} else if title == "[kit]" {
			key := moose.StripMinecraftColour(lines[1])
			u, err := data.LoadUserOrCreate(h.p.Name())
			if err != nil {
				return
			}
			cd := u.Kits.Key(key)
			if cd.Active() {
				h.Message("command.kit.cooldown", cd.Remaining().Round(time.Second))
				return
			} else {
				cd.Set(time.Minute)
			}
			switch strings.ToLower(moose.StripMinecraftColour(lines[1])) {
			case "diamond":
				kit.Apply(kit.Diamond{}, h.p)
			case "archer":
				kit.Apply(kit.Archer{}, h.p)
			case "bard":
				kit.Apply(kit.Bard{}, h.p)
			case "rogue":
				kit.Apply(kit.Rogue{}, h.p)
			case "stray":
				kit.Apply(kit.Stray{}, h.p)
			case "miner":
				kit.Apply(kit.Miner{}, h.p)
			case "builder":
				kit.Apply(kit.Builder{}, h.p)
			}
		}
	}
}

func (h *Handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, _ *bool) {
	*force, *height = 0.394, 0.394
	t, ok := e.(*player.Player)
	if !ok {
		return
	}

	if _, ok := h.Player().Effect(effect.Invisibility{}); ok {
		for _, i := range h.Player().Armour().Inventory().Items() {
			if _, ok := i.Enchantment(ench.Invisibility{}); !ok {
				h.Player().RemoveEffect(effect.Invisibility{})
			}
		}

		h.ShowArmor(true)
	}

	if !canAttack(h.p, t) {
		ctx.Cancel()
		return
	}
	if t.AttackImmune() {
		return
	}

	arm := h.p.Armour()
	for _, a := range arm.Slots() {
		for _, e := range a.Enchantments() {
			if att, ok := e.Type().(ench.AttackEnchantment); ok {
				att.AttackEntity(h.p, t)
			}
		}
	}

	held, left := h.p.HeldItems()
	if s, ok := held.Item().(item.Sword); ok && s.Tier == item.ToolTierGold && class.Compare(h.class.Load(), class.Rogue{}) && t.Rotation().Direction() == h.p.Rotation().Direction() {
		cd := h.rogue
		w := h.p.World()
		if cd.Active() {
			h.p.Message(lang.Translatef(h.p.Locale(), "user.cool-down", "Rogue", cd.Remaining().Seconds()))
		} else {
			ctx.Cancel()
			for i := 1; i <= 3; i++ {
				w.AddParticle(t.Position().Add(mgl64.Vec3{0, float64(i), 0}), particle.Dust{
					Colour: item.ColourRed().RGBA(),
				})
			}
			w.PlaySound(h.p.Position(), sound.ItemBreak{})
			t.Hurt(8, NoArmourAttackEntitySource{
				Attacker: h.p,
			})
			t.KnockBack(h.p.Position(), *force, *height)

			h.p.SetHeldItems(item.Stack{}, left)
			cd.Set(time.Second * 10)
		}
	}

	//u, err := data.LoadUserOrCreate(h.p.Name(), h.p.Handler().(*Handler).XUID())
	target, ok := Lookup(t.Name())
	typ, ok2 := it.SpecialItem(held)
	if ok && ok2 {
		if cd := h.ability; cd.Active() {
			h.p.Message(text.Colourf("<red>You are on Partner Items cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		switch kind := typ.(type) {
		case it.ExoticBoneType:
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on bone cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			target.AddBoneHit(t)
			if target.Boned() {
				target.Player().Message(text.Colourf("<red>You have been boned by %s</red>", h.p.Name()))
				h.p.Message(text.Colourf("<green>You have boned %s</green>", t.Name()))
				h.ability.Set(time.Second * 10)
				h.abilities.Set(kind, time.Minute)
				h.p.SetHeldItems(h.SubtractItem(held, 1), left)
				target.ResetBoneHits(h.p)
			} else {
				h.p.Message(text.Colourf("<green>You have hit %s with a bone %d times</green>", t.Name(), target.BoneHits(h.p)))
			}
		case it.ScramblerType:
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on scrambler cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			target.AddScramblerHit(h.p)
			if target.ScramblerHits(h.p) >= 3 {
				for i := 0; i <= 8; i++ {
					j := rand.Intn(8)
					it1, _ := h.p.Inventory().Item(i)
					it2, _ := h.p.Inventory().Item(j)
					target.Player().Inventory().SetItem(j, it1)
					target.Player().Inventory().SetItem(i, it2)
				}
				target.Player().Message(text.Colourf("<red>You have been scrambled by %s</red>", h.p.Name()))
				h.p.Message(text.Colourf("<green>You have scrambled %s</green>", t.Name()))
				h.ability.Set(time.Second * 10)
				h.abilities.Set(kind, time.Minute*2)
				target.ResetScramblerHits(h.p)
				h.p.SetHeldItems(h.SubtractItem(held, 1), left)
			}
		case it.PearlDisablerType:
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on pearl disabler cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			if !target.PearlDisabled() {
				target.Player().Message(text.Colourf("<red>You have been pearl disabled by %s</red>", h.p.Name()))
				h.p.Message(text.Colourf("<green>You have pearl disabled %s</green>", t.Name()))
				target.TogglePearlDisable()
				h.ability.Set(time.Second * 10)
				h.abilities.Set(kind, time.Minute)
				h.p.SetHeldItems(h.SubtractItem(held, 1), left)
			}
		}
	}

	th, ok := t.Handler().(*Handler)
	if !ok {
		return
	}

	if canAttack(h.p, t) {
		th.SetLastAttacker(h)
		th.combat.Set(time.Second * 20)
		h.combat.Set(time.Second * 20)
	}
}

func (h *Handler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if h.logger {
		return
	}
	u, _ := data.LoadUserOrCreate(h.p.Name())
	p := h.p
	w := p.World()

	if !newPos.ApproxEqual(p.Position()) {
		h.home.Cancel()
		h.logout.Cancel()
		h.stuck.Cancel()
	}

	if h.waypoint != nil && h.waypoint.active {
		h.UpdateWayPointPosition()
	}

	h.clearWall()
	cubePos := cube.PosFromVec3(newPos)

	if h.combat.Active() {
		a := area.Spawn(w)
		mul := moose.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
		if mul.Vec3WithinOrEqualFloorXZ(p.Position()) {
			h.sendWall(cubePos, area.Overworld.Spawn().Area, item.ColourRed())
		}
	}

	if u.PVP.Active() {
		for _, a := range data.Teams() {
			a := a.Claim
			if a != (moose.Area{}) && a.Vec3WithinOrEqualXZ(newPos) {
				ctx.Cancel()
				return
			}

			mul := moose.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
			if mul.Vec2WithinOrEqualFloor(mgl64.Vec2{p.Position().X(), p.Position().Z()}) && !area.Spawn(p.World()).Vec3WithinOrEqualFloorXZ(newPos) {
				h.sendWall(cubePos, a, item.ColourBlue())
			}
		}

		if newPos.Y() < 0 {
			h.p.Teleport(mgl64.Vec3{0, 100, 0})
		}
	}

	if _, ok := sotw.Running(); ok && u.SOTW {
		if newPos.Y() < 0 {
			h.p.Teleport(mgl64.Vec3{0, 100, 0})
		}
	}

	if area.Spawn(w).Vec3WithinOrEqualFloorXZ(newPos) && h.combat.Active() {
		ctx.Cancel()
		return
	}

	us, ok := Lookup(u.Name)
	if !ok {
		return
	}
	k, ok := koth.Running()
	if ok {
		r := u.Roles.Highest()
		if k.Area().Vec3WithinOrEqualFloorXZ(newPos) {
			// Need to handle for Y-axis cases because some koths are irregular
			if k == koth.Spiral {
				if newPos.Y() < 120 {
					return
				}
			} else if k == koth.Dragon {
				if newPos.Y() < 109 {
					return
				}
			} else if k == koth.Stairs {
				if newPos.Y() < 150 {
					return
				}
			}

			if u.PVP.Active() {
				return
			}
			if k.StartCapturing(us) {
				Broadcast("koth.capturing", k.Name(), r.Colour(u.Name))
			}
		} else {
			if k.StopCapturing(us) {
				Broadcast("koth.not.capturing", k.Name())
			}
		}
	}

	var areas []moose.NamedArea

	for _, tm := range data.Teams() {
		a := tm.Claim

		name := text.Colourf("<red>%s</red>", tm.DisplayName)
		if t, ok := u.Team(); ok && strings.EqualFold(t.Name, tm.Name) {
			name = text.Colourf("<green>%s</green>", tm.DisplayName)
		}
		areas = append(areas, moose.NewNamedArea(mgl64.Vec2{a.Min().X(), a.Min().Y()}, mgl64.Vec2{a.Max().X(), a.Max().Y()}, name))
	}

	ar := h.area.Load()
	for _, a := range append(area.Protected(w), areas...) {
		if a.Vec3WithinOrEqualFloorXZ(newPos) {
			if ar != a {
				if ar != (moose.NamedArea{}) {
					h.Message("area.leave", ar.Name())
				}
				h.area.Store(a)
				h.Message("area.enter", a.Name())
				return
			} else {
				return
			}
		}
	}

	if ar != area.Wilderness(w) {
		if ar != (moose.NamedArea{}) {
			h.Message("area.leave", ar.Name())

		}
		h.area.Store(area.Wilderness(w))
		h.Message("area.enter", area.Wilderness(w).Name())
	}
}

func (h *Handler) HandleQuit() {
	if h.logger {
		return
	}
	h.close <- struct{}{}
	p := h.p

	u, _ := data.LoadUserOrCreate(p.Name())
	u.PlayTime += time.Since(h.logTime)
	_ = data.SaveUser(u)

	tm, _ := u.Team()
	_, sotwRunning := sotw.Running()
	if !h.loggedOut && !tm.Claim.Vec3WithinOrEqualFloorXZ(p.Position()) && !area.Spawn(p.World()).Vec3WithinOrEqualFloorXZ(p.Position()) || ((sotwRunning && u.SOTW) || u.PVP.Active()) {
		arm := h.p.Armour()
		inv := h.p.Inventory()

		h.p = player.New(p.Name(), p.Skin(), p.Position())
		h.p.SetNameTag(text.Colourf("<red>%s</red> <grey>(LOGGER)</grey>", p.Name()))
		h.p.Handle(h)
		if p.Health() < 20 {
			h.p.Hurt(20-p.Health(), effect.InstantDamageSource{})
		}

		for j, i := range inv.Slots() {
			_ = h.p.Inventory().SetItem(j, i)
		}
		h.p.Armour().Set(arm.Helmet(), arm.Chestplate(), arm.Leggings(), arm.Boots())

		p.World().AddEntity(h.p)
		go func() {
			select {
			case <-time.After(time.Second * 30):
			case <-h.close:
			case <-h.death:
				u, ok := data.LoadUser(h.p.Name())
				if !ok {
					return
				}
				u.Dead = true
				u.Stats.Deaths = 0
				if u.Stats.KillStreak > u.Stats.BestKillStreak {
					u.Stats.BestKillStreak = u.Stats.KillStreak
				}
				u.Stats.KillStreak = 0
				if tm, ok := u.Team(); ok {
					tm = tm.WithDTR(tm.DTR - 1).WithPoints(tm.Points - 1).WithRegenerationTime(time.Now().Add(time.Minute * 5))
					data.SaveTeam(tm)
				}
				DropContents(h.p)
				_ = data.SaveUser(u)
			}
			playersMu.Lock()
			delete(players, h.p.Name())
			playersMu.Unlock()
			_ = h.p.Close()
		}()
		h.logger = true
		h.UpdateState()
		return
	}
	playersMu.Lock()
	delete(players, h.p.Name())
	playersMu.Unlock()
}

type npcHandler struct {
	player.NopHandler
}

func (npcHandler) HandleItemPickup(ctx *event.Context, _ *item.Stack) {
	ctx.Cancel()
}

type NoArmourAttackEntitySource struct {
	Attacker world.Entity
}

func (NoArmourAttackEntitySource) Fire() bool {
	return false
}

func (NoArmourAttackEntitySource) ReducedByArmour() bool {
	return false
}

func (NoArmourAttackEntitySource) ReducedByResistance() bool {
	return false
}

// attackerFromSource returns the Attacker from a DamageSource. If the source is not an entity false is
// returned.
func attackerFromSource(src world.DamageSource) (world.Entity, bool) {
	switch s := src.(type) {
	case entity.AttackDamageSource:
		return s.Attacker, true
	case NoArmourAttackEntitySource:
		return s.Attacker, true
	}
	return nil, false
}

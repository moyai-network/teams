package user

import (
	"fmt"
	"github.com/moyai-network/teams/moyai"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/moyai-network/teams/internal/cooldown"
	"github.com/moyai-network/teams/internal/effectutil"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/internal/unsafe"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/class"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/moyai-network/teams/moyai/crate"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/koth"
	"github.com/moyai-network/teams/moyai/process"
	"github.com/moyai-network/teams/moyai/role"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	it "github.com/moyai-network/teams/moyai/item"
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

	loggers = map[string]*Handler{}
)

type Handler struct {
	player.NopHandler
	s *session.Session
	p *player.Player

	logTime time.Time

	vanished atomic.Bool

	sign cube.Pos

	pearl       *cooldown.CoolDown
	rogue       *cooldown.CoolDown
	goldenApple *cooldown.CoolDown

	factionCreate *cooldown.CoolDown

	itemUse cooldown.MappedCoolDown[world.Item]

	archerRogueItem cooldown.MappedCoolDown[world.Item]
	bardItem        cooldown.MappedCoolDown[world.Item]
	strayItem       cooldown.MappedCoolDown[world.Item]

	ability   *cooldown.CoolDown
	abilities cooldown.MappedCoolDown[it.SpecialItemType]

	waypoint *WayPoint

	armour atomic.Value[[4]item.Stack]
	class  atomic.Value[class.Class]
	energy atomic.Value[float64]

	combat *cooldown.CoolDown
	archer *cooldown.CoolDown

	boneHits map[string]int
	bone     *cooldown.CoolDown

	scramblerHits map[string]int
	pearlDisabled bool

	lastPearlPos mgl64.Vec3
	lastHitBy    *player.Player

	logout *process.Process
	stuck  *process.Process
	home   *process.Process

	lastScoreBoard atomic.Value[*scoreboard.Scoreboard]
	area           atomic.Value[area.NamedArea]

	lastAttackerName atomic.Value[string]
	lastAttackTime   atomic.Value[time.Time]

	lastMessage atomic.Value[time.Time]

	wallBlocks   map[cube.Pos]float64
	wallBlocksMu sync.Mutex

	claimPos [2]mgl64.Vec2

	l language.Tag

	loggedOut bool
	logger    bool
	close     chan struct{}
	death     chan struct{}
}

func NewHandler(p *player.Player, xuid string) *Handler {
	h, ok := loggers[p.XUID()]
	if ok {
		_ = h.p.Close()
		h.logger = false
	}
	ha := &Handler{
		p: p,

		pearl:       cooldown.NewCoolDown(),
		rogue:       cooldown.NewCoolDown(),
		goldenApple: cooldown.NewCoolDown(),
		ability:     cooldown.NewCoolDown(),

		factionCreate: cooldown.NewCoolDown(),

		bone:     cooldown.NewCoolDown(),
		boneHits: map[string]int{},

		scramblerHits: map[string]int{},

		wallBlocks: map[cube.Pos]float64{},

		itemUse:         cooldown.NewMappedCoolDown[world.Item](),
		archerRogueItem: cooldown.NewMappedCoolDown[world.Item](),
		bardItem:        cooldown.NewMappedCoolDown[world.Item](),
		strayItem:       cooldown.NewMappedCoolDown[world.Item](),
		abilities:       cooldown.NewMappedCoolDown[it.SpecialItemType](),

		combat: cooldown.NewCoolDown(),
		archer: cooldown.NewCoolDown(),

		home: process.NewProcess(func(t *process.Process) {
			p.Message(text.Colourf("<green>You have been teleported home.</green>"))
		}),
		stuck: process.NewProcess(func(t *process.Process) {
			p.Message(text.Colourf("<red>You have been teleported to a safe place.</red>"))
		}),

		close: make(chan struct{}),
		death: make(chan struct{}),
	}
	ha.logout = process.NewProcess(func(t *process.Process) {
		ha.loggedOut = true
		p.Disconnect(text.Colourf("<red>You have been logged out.</red>"))
	})

	p.SetNameTag(text.Colourf("<red>%s</red>", p.Name()))
	UpdateState(p)

	// if p, ok := Lookup(p.Name()); ok {
	// 	p.Teleport(p.Position())
	// 	if h, ok := p.Handler().(*Handler); ok {
	// 		close(h.close)
	// 	}
	// }

	s := player_session(p)
	u, _ := data.LoadUserFromName(p.Name())

	if u.Teams.Dead {
		p.Armour().Clear()
		p.Inventory().Clear()
		p.Teleport(p.World().Spawn().Vec3Middle())
		p.Heal(20, effect.InstantHealingSource{})
		u.Teams.Dead = false
	}

	/*if u.Frozen {
		p.SetImmobile()
	}*/

	u.DisplayName = p.Name()
	u.Name = strings.ToLower(p.Name())
	u.XUID = xuid

	u.DeviceID = s.ClientData().DeviceID
	u.SelfSignedID = s.ClientData().SelfSignedID
	data.SaveUser(u)

	ha.s = s
	ha.logTime = time.Now()

	UpdateState(ha.p)
	go startTicker(ha)
	return ha
}

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

	/*if !u.Mute.Expired() {
		h.p.Message(lang.Translatef(l, "user.message.mute"))
		return
	}*/
	tm, teamErr := data.LoadTeamFromMemberName(h.p.Name())

	if msg := strings.TrimSpace(*message); len(msg) > 0 {
		msg = formatRegex.ReplaceAllString(msg, "")

		global := func() {
			if teamErr == nil {
				formatTeam := text.Colourf("<grey>[<green>%s</green>]</grey> %s", tm.DisplayName, r.Chat(h.p.Name(), msg))
				formatEnemy := text.Colourf("<grey>[<red>%s</red>]</grey> %s", tm.DisplayName, r.Chat(h.p.Name(), msg))
				for _, t := range moyai.Server().Players() {
					if tm.Member(t.Name()) {
						t.Message(formatTeam)
					} else {
						t.Message(formatEnemy)
					}
				}
				chat.StdoutSubscriber{}.Message(formatEnemy)
			} else {
				_, _ = chat.Global.WriteString(r.Chat(h.p.Name(), msg))
			}
		}

		staff := func() {
			for _, s := range moyai.Server().Players() {
				if us, err := data.LoadUserOrCreate(s.Name(), s.XUID()); err == nil && role.Staff(us.Roles.Highest()) {
					Messagef(s, "staff.chat", r.Name(), h.p.Name(), strings.TrimPrefix(msg, "!"))
				}
			}
		}
		fmt.Println(u.Teams.ChatType)
		h.lastMessage.Store(time.Now())
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

func (h *Handler) HandleFoodLoss(ctx *event.Context, _ int, _ *int) {
	ctx.Cancel()
}

func (h *Handler) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	p := h.Player()
	u, err := data.LoadUserFromXUID(h.p.XUID())
	if err != nil {
		return
	}

	w := p.World()
	b := w.Block(pos)

	held, _ := p.HeldItems()
	typ, ok := it.SpecialItem(held)
	if ok {
		if cd := h.ability; cd.Active() {
			p.SendJukeboxPopup(lang.Translatef(u.Language, "popup.cooldown.partner_item", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		if spi := h.abilities; spi.Active(typ) {
			p.SendJukeboxPopup(lang.Translatef(u.Language, "popup.cooldown.partner_item.item", typ.Name(), spi.Remaining(typ).Seconds()))
			ctx.Cancel()
		} else {
			p.SendJukeboxPopup(lang.Translatef(u.Language, "popup.ready.partner_item.item", typ.Name()))
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
	u, err := data.LoadUserFromXUID(h.p.XUID())
	if err != nil {
		return
	}

	held, _ := p.HeldItems()
	typ, ok := it.SpecialItem(held)
	if ok {
		if cd := h.ability; cd.Active() {
			p.SendJukeboxPopup(lang.Translatef(u.Language, "popup.cooldown.partner_item", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		if spi := h.abilities; spi.Active(typ) {
			p.SendJukeboxPopup(lang.Translatef(u.Language, "popup.cooldown.partner_item.item", typ.Name(), spi.Remaining(typ).Seconds()))
			ctx.Cancel()
		} else {
			p.SendJukeboxPopup(lang.Translatef(u.Language, "popup.ready.partner_item.item", typ.Name()))
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

	t, err := data.LoadTeamFromMemberName(h.p.Name())
	if err != nil {
		return
	}
	if !t.Leader(u.Name) {
		Messagef(h.p, "team.not-leader")
		return
	}
	if t.Claim != (area.Area{}) {
		Messagef(h.p, "team.has-claim")
		return
	}

	pos := h.claimPos
	if pos[0] == (mgl64.Vec2{}) || pos[1] == (mgl64.Vec2{}) {
		Messagef(h.p, "team.claim.select-before")
		return
	}
	claim := area.NewArea(pos[0], pos[1])
	var blocksPos []cube.Pos
	mn := claim.Min()
	mx := claim.Max()
	for x := mn[0]; x <= mx[0]; x++ {
		for y := mn[1]; y <= mx[1]; y++ {
			blocksPos = append(blocksPos, cube.PosFromVec3(mgl64.Vec3{x, 0, y}))
		}
	}
	for _, a := range area.Protected(w) {
		var threshold float64 = 1
		message := "team.area.too-close"
		for _, k := range area.KOTHs(h.p.World()) {
			if a.Area == k.Area {
				threshold = 25
				message = "team.area.too-close.koth"
			}
		}

		for _, b := range blocksPos {
			if a.Vec3WithinOrEqualXZ(b.Vec3()) {
				Messagef(h.p, "team.area.already-claimed")
				return
			}
			if areaTooClose(a.Area, vec3ToVec2(b.Vec3()), threshold) {
				Messagef(h.p, message)
				return
			}
		}
		if a.Vec2WithinOrEqual(pos[0]) || a.Vec2WithinOrEqual(pos[1]) {
			Messagef(h.p, "team.area.already-claimed")
			return
		}
		if areaTooClose(a.Area, pos[0], threshold) || areaTooClose(a.Area, pos[1], threshold) {
			Messagef(h.p, message)
			return
		}
	}

	teams, err := data.LoadAllTeams()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, tm := range teams {
		c := tm.Claim
		if c == (area.Area{}) {
			continue
		}
		for _, b := range blocksPos {
			if c.Vec3WithinOrEqualXZ(b.Vec3()) {
				Messagef(h.p, "team.area.already-claimed")
				return
			}
			if areaTooClose(c, vec3ToVec2(b.Vec3()), 1) {
				Messagef(h.p, "team.area.too-close")
				return
			}
		}
		if c.Vec2WithinOrEqual(pos[0]) || c.Vec2WithinOrEqual(pos[1]) {
			Messagef(h.p, "team.area.already-claimed")
			return
		}
		if areaTooClose(c, pos[0], 1) || areaTooClose(c, pos[1], 1) {
			Messagef(h.p, "team.area.too-close")
			return
		}
	}

	x := claim.Max().X() - claim.Min().X()
	y := claim.Max().Y() - claim.Min().Y()
	ar := x * y
	if ar > 75*75 {
		Messagef(h.p, "team.claim.too-big")
		return
	}
	cost := ar * 5

	if t.Balance < cost {
		Messagef(h.p, "team.claim.no-money")
		return
	}

	t.Balance -= cost
	t = t.WithClaim(claim)
	data.SaveTeam(t)

	Messagef(h.p, "command.claim.success", pos[0], pos[1], cost)
}

func vec3ToVec2(v mgl64.Vec3) mgl64.Vec2 {
	return mgl64.Vec2{v.X(), v.Z()}
}

func areaTooClose(area area.Area, pos mgl64.Vec2, threshold float64) bool {
	var vectors []mgl64.Vec2
	for x := -threshold; x <= threshold; x++ {
		for y := -threshold; y <= threshold; y++ {
			vectors = append(vectors, mgl64.Vec2{pos.X() + x, pos.Y() + y})
		}
	}

	for _, v := range vectors {
		if area.Vec2WithinOrEqual(v) {
			return true
		}
	}
	return false
}

// HandleItemUse ...
func (h *Handler) HandleItemUse(ctx *event.Context) {
	held, left := h.p.HeldItems()
	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil {
		return
	}
	if v, ok := held.Value("MONEY_NOTE"); ok {
		u.Teams.Balance = u.Teams.Balance + v.(float64)
		h.p.SetHeldItems(h.SubtractItem(held, 1), left)
		h.p.Message(text.Colourf("<green>You have deposited $%.0f into your bank account</green>", v.(float64)))
		data.SaveUser(u)
		return
	}
	switch held.Item().(type) {
	case item.EnderPearl:
		if cd := h.pearl; cd.Active() {
			Messagef(h.p, "user.cool-down", "Ender Pearl", cd.Remaining().Seconds())
			ctx.Cancel()
		} else if h.pearlDisabled {
			h.p.Message(text.Colourf("<red>Pearl Disabled! Pearl was refunded!</red>"))
			h.TogglePearlDisable()
			ctx.Cancel()
		} else {
			cd.Set(15 * time.Second)
			h.lastPearlPos = h.p.Position()
		}
	}

	switch h.class.Load().(type) {
	case class.Archer, class.Rogue:
		if e, ok := ArcherRogueEffectFromItem(held.Item()); ok {
			if cd := h.archerRogueItem.Key(held.Item()); cd.Active() {
				Messagef(h.p, "class.ability.cooldown", cd.Remaining().Seconds())
				return
			}
			h.p.AddEffect(e)
			go func() {
				select {
				case <-time.After(e.Duration()):
					sortArmourEffects(h)
					sortClassEffects(h)
				case <-h.close:
					return
				}
			}()
			h.archerRogueItem.Key(held.Item()).Set(60 * time.Second)
			h.p.SetHeldItems(held.Grow(-1), item.Stack{})
		}
	case class.Bard:
		if e, ok := BardEffectFromItem(held.Item()); ok {
			_, sotwRunning := sotw.Running()
			if u.Teams.PVP.Active() || sotwRunning && u.Teams.SOTW {
				return
			}
			if cd := h.bardItem.Key(held.Item()); cd.Active() {
				Messagef(h.p, "class.ability.cooldown", cd.Remaining().Seconds())
				return
			}
			if en := h.energy.Load(); en < 30 {
				Messagef(h.p, "class.energy.insufficient")
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
			Messagef(h.p, "class.ability.use", effectutil.EffectName(e), lvl, len(teammates))
			h.p.SetHeldItems(held.Grow(-1), item.Stack{})
			h.bardItem.Key(held.Item()).Set(15 * time.Second)
		}
	case class.Stray:
		if e, ok := StrayEffectFromItem(held.Item()); ok {
			_, sotwRunning := sotw.Running()
			if u.Teams.PVP.Active() || sotwRunning && u.Teams.SOTW {
				return
			}

			if cd := h.strayItem.Key(held.Item()); cd.Active() {
				Messagef(h.p, "class.ability.cooldown", cd.Remaining().Seconds())
				return
			}
			if en := h.energy.Load(); en < 30 {
				Messagef(h.p, "class.energy.insufficient")
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
			Messagef(h.p, "class.ability.use", effectutil.EffectName(e), lvl, len(teammates))
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
		case it.TimeWarpType:
			if h.lastPearlPos == (mgl64.Vec3{}) {
				h.p.Message(text.Colourf("<red>You do not have a last thrown pearl.</red>"))
				break
			}
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on timewarp cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			h.p.Message(text.Colourf("<green>Ongoing in 2 seconds...</green>"))
			h.ability.Set(time.Second * 10)
			h.abilities.Set(kind, time.Minute*1)
			time.AfterFunc(time.Second*2, func() {
				if h.p != nil {
					h.p.Teleport(h.lastPearlPos)
					h.lastPearlPos = mgl64.Vec3{}
				}
			})
		case it.SigilType:
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on sigil cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			nb := NearbyCombat(h.p, 10)
			for _, e := range nb {
				if e.p == h.p {
					continue
				}
				e.p.World().AddEntity(entity.NewLightning(e.p.Position()))
				e.p.Hurt(4, NoArmourAttackEntitySource{})
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
			if h.lastHitBy == nil {
				h.p.Message(text.Colourf("<red>No last hit found</red>"))
				break
			}
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on Ninja Star cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			h.p.Message(text.Colourf("<green>Ongoing to %s in 5 seconds...</red>", h.lastHitBy.Name()))
			h.lastHitBy.Message(text.Colourf("<red>%s is teleporting to your in 5 seconds...</red>", h.p.Name()))
			h.ability.Set(time.Second * 10)
			h.abilities.Set(kind, time.Minute*2)
			time.AfterFunc(time.Second*5, func() {
				if h.p != nil && h.lastHitBy != nil {
					h.p.Teleport(h.lastHitBy.Position())
				}
			})
		}
	}
}

func (h *Handler) HandleSignEdit(ctx *event.Context, frontSide bool, oldText, newText string) {
	ctx.Cancel()
	if !frontSide {
		return
	}

	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil {
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
	if h.archer.Active() {
		*dmg = *dmg + *dmg*0.15
	}
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
				if be, ok := et.Type().(entity.Behaviour); ok {
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

			u, err := data.LoadUserFromName(h.p.Name())
			if err != nil {
				return
			}
			u.Teams.PVP.Set(time.Hour)
			data.SaveUser(u)

			DropContents(h.p)
			p.SetHeldItems(item.Stack{}, item.Stack{})

			p.ResetFallDistance()
			p.Heal(20, effect.InstantHealingSource{})
			p.Extinguish()
			p.SetFood(20)
			h.class.Store(class.Resolve(p))
			UpdateState(h.p)

			// TODO, add deathban later
			h.p.Teleport(mgl64.Vec3{0, 100, 0})
			//h.p.SetMobile()

			if tm, err := data.LoadTeamFromMemberName(h.p.Name()); err == nil {
				tm = tm.WithDTR(tm.DTR - 1).WithPoints(tm.Points - 1).WithRegenerationTime(time.Now().Add(time.Minute * 5))
				data.SaveTeam(tm)
			}

			victim, err := data.LoadUserFromName(h.p.Name())
			if err == nil {
				victim.Teams.Stats.Deaths += 1
				if victim.Teams.Stats.KillStreak > victim.Teams.Stats.BestKillStreak {
					victim.Teams.Stats.BestKillStreak = victim.Teams.Stats.KillStreak
				}
				victim.Teams.Stats.KillStreak = 0
				data.SaveUser(victim)
			}

			killer, ok := h.LastAttacker()
			if ok {
				k, err := data.LoadUserFromName(killer.Name())
				if err != nil {
					return
				}
				k.Teams.Stats.Kills += 1

				if tm, err := data.LoadTeamFromMemberName(killer.Name()); err == nil {
					tm = tm.WithPoints(tm.Points + 1)
					data.SaveTeam(tm)
				}
				data.SaveUser(k)

				held, _ := killer.HeldItems()
				heldName := held.CustomName()

				if len(heldName) <= 0 {
					heldName = it.DisplayName(held.Item())
				}

				if held.Empty() || len(heldName) <= 0 {
					heldName = "their fist"
				}

				_, _ = chat.Global.WriteString(lang.Translatef(data.Language{}, "user.kill", p.Name(), u.Teams.Stats.Kills, killer.Name(), k.Teams.Stats.Kills, text.Colourf("<red>%s</red>", heldName)))
				h.ResetLastAttacker()
			} else {
				_, _ = chat.Global.WriteString(lang.Translatef(data.Language{}, "user.suicide", p.Name(), u.Teams.Stats.Kills))
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

	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil || u.Teams.PVP.Active() {
		ctx.Cancel()
		return
	}

	if u.Frozen {
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
		u, err := data.LoadUserFromName(h.p.Name())
		if err != nil || u.Teams.PVP.Active() {
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
		h.lastHitBy = attacker
	case entity.AttackDamageSource:
		if t, ok := s.Attacker.(*player.Player); ok {
			attacker = t
		}
		if !canAttack(h.p, attacker) {
			ctx.Cancel()
			return
		}
		h.lastHitBy = attacker
	case entity.VoidDamageSource:
		if u.Teams.PVP.Active() {
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
				if pl, ok := k.Capturing(); ok && pl == h.p {
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

		h.lastHitBy = attacker

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

				attacker.Message(lang.Translatef(data.Language{}, "archer.tag", math.Round(dist), damage/2))
			}

		}
	}

	if attacker != nil {
		//if _, ok := h.Player().Effect(effect.Invisibility{}); ok {
		//for _, i := range h.Player().Armour().Inventory().Items() {
		//if _, ok := i.Enchantment(ench.Invisibility{}); !ok {
		//	h.Player().RemoveEffect(effect.Invisibility{})
		//}
		//}

		//h.ShowArmor(true)
		//}
		percent := 0.90
		e, ok := attacker.Effect(effect.Strength{})
		if e.Level() > 1 {
			percent = 0.85
		}

		if ok {
			*dmg = *dmg * percent
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
					fmt.Println(pro.Owner().Position())
					if pro.Owner() == p {
						et.Close()
						h.p.World().RemoveEntity(et)
					}
				}
			}
		}

		h.combat.Reset()
		h.pearl.Reset()
		h.archer.Reset()

		victim, err := data.LoadUserFromName(h.p.Name())
		if err != nil {
			return
		}
		victim.Teams.PVP.Set(time.Hour)
		victim.Teams.Stats.Deaths += 1
		if victim.Teams.Stats.KillStreak > victim.Teams.Stats.BestKillStreak {
			victim.Teams.Stats.BestKillStreak = victim.Teams.Stats.KillStreak
		}
		victim.Teams.Stats.KillStreak = 0
		data.SaveUser(victim)

		DropContents(h.p)
		p.SetHeldItems(item.Stack{}, item.Stack{})

		//h.EnableDeathban()
		// u.SubtractLife()
		//deathban.Deathban().AddPlayer(p)

		p.ResetFallDistance()
		p.Heal(20, effect.InstantHealingSource{})
		p.Extinguish()
		p.SetFood(20)
		h.class.Store(class.Resolve(p))
		UpdateState(h.p)

		// TODO, add deathban later
		h.p.Teleport(mgl64.Vec3{0, 100, 0})
		//h.p.SetMobile()

		if tm, err := data.LoadTeamFromMemberName(h.p.Name()); err == nil {
			tm = tm.WithDTR(tm.DTR - 1).WithPoints(tm.Points - 1).WithRegenerationTime(time.Now().Add(time.Minute * 5))
			data.SaveTeam(tm)
		}

		killer, ok := h.LastAttacker()
		if ok {
			k, err := data.LoadUserFromName(killer.Name())
			if err != nil {
				return
			}
			k.Teams.Stats.Kills += 1
			k.Teams.Stats.KillStreak += 1

			if k.Teams.Stats.KillStreak%5 == 0 {
				Broadcast("user.killstreak", killer.Name(), k.Teams.Stats.KillStreak)
				it.AddOrDrop(killer, it.NewKey(it.KeyTypePartner, int(k.Teams.Stats.KillStreak)/2))
			}

			if tm, err := data.LoadTeamFromMemberName(killer.Name()); err == nil {
				tm = tm.WithPoints(tm.Points + 1)
				data.SaveTeam(tm)
			}
			data.SaveUser(k)

			held, _ := killer.HeldItems()
			heldName := held.CustomName()

			if len(heldName) <= 0 {
				heldName = it.DisplayName(held.Item())
			}

			if held.Empty() || len(heldName) <= 0 {
				heldName = "their fist"
			}

			_, _ = chat.Global.WriteString(lang.Translatef(data.Language{}, "user.kill", p.Name(), u.Teams.Stats.Kills, killer.Name(), k.Teams.Stats.Kills, text.Colourf("<red>%s</red>", heldName)))
			h.ResetLastAttacker()
		} else {
			_, _ = chat.Global.WriteString(lang.Translatef(data.Language{}, "user.suicide", p.Name(), u.Teams.Stats.Kills))
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

		it.AddOrDrop(h.p, i)

		w.AddEntity(entity.NewFirework(pos.Vec3(), cube.Rotation{90, 90}, item.Firework{
			Duration: 0,
			Explosions: []item.FireworkExplosion{
				{
					Shape:   item.FireworkShapeStar(),
					Trail:   true,
					Colour:  colour.RandomColour(),
					Twinkle: true,
				},
			},
		}))
		return
	}

	teams, _ := data.LoadAllTeams()
	for _, t := range teams {
		if !t.Member(h.p.Name()) {
			if t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
		u, _ := data.LoadUserFromName(h.p.Name())
		for _, a := range area.Protected(w) {
			if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
				if !u.Roles.Contains(role.Operator{}) || h.p.GameMode() != world.GameModeCreative {
					ctx.Cancel()
					return
				}
			}
		}
	}
}

func (h *Handler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	w := h.p.World()

	u, _ := data.LoadUserFromName(h.p.Name())
	for _, a := range area.Protected(w) {
		if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
			if !u.Roles.Contains(role.Operator{}) || h.p.GameMode() != world.GameModeCreative {
				ctx.Cancel()
				return
			}
		}
	}

	teams, _ := data.LoadAllTeams()
	for _, t := range teams {
		if !t.Member(h.p.Name()) {
			if t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
	}
}

func (h *Handler) HandleItemConsume(ctx *event.Context, i item.Stack) {
	switch i.Item().(type) {
	case item.GoldenApple:
		if cd := h.goldenApple; cd.Active() {
			h.p.Message(text.Colourf("<red>You are on golden apple cooldown.</red>"))
			ctx.Cancel()
		} else {
			cd.Set(time.Second * 30)
		}
	}
}

func (h *Handler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	w := h.p.World()

	i, left := h.p.HeldItems()
	b := w.Block(pos)

	for _, c := range crate.All() {
		if _, ok := b.(block.Chest); ok && pos.Vec3Middle() == c.Position() {
			if _, ok := i.Value("crate-key_" + colour.StripMinecraftColour(c.Name())); !ok {
				h.p.Message(text.Colourf("<red>You need a %s key to open this crate</red>", colour.StripMinecraftColour(c.Name())))
				break
			}
			it.AddOrDrop(h.p, ench.AddEnchantmentLore(crate.SelectReward(c)))

			h.p.SetHeldItems(h.SubtractItem(i, 1), left)

			w.AddEntity(entity.NewFirework(c.Position().Add(mgl64.Vec3{0, 1, 0}), cube.Rotation{90, 90}, item.Firework{
				Duration: 0,
				Explosions: []item.FireworkExplosion{
					{
						Shape:   item.FireworkShapeStar(),
						Trail:   true,
						Colour:  colour.RandomColour(),
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
			tm, err := data.LoadTeamFromMemberName(h.p.Name())
			if err != nil {
				return
			}

			if !tm.Leader(h.p.Name()) {
				Messagef(h.p, "team.not-leader")
				return
			}

			if tm.Claim != (area.Area{}) {
				Messagef(h.p, "team.has-claim")
				break
			}

			for _, a := range area.Protected(w) {
				var threshold float64 = 1
				message := "team.area.too-close"
				for _, k := range area.KOTHs(h.p.World()) {
					if a.Area == k.Area {
						threshold = 25
						message = "team.area.too-close.koth"
					}
				}

				if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
					Messagef(h.p, "team.area.already-claimed")
					return
				}
				if areaTooClose(a.Area, vec3ToVec2(pos.Vec3()), threshold) {
					Messagef(h.p, message)
					return
				}
			}

			teams, _ := data.LoadAllTeams()
			for _, t := range teams {
				c := t.Claim
				if c != (area.Area{}) {
					continue
				}
				if c.Vec3WithinOrEqualXZ(pos.Vec3()) {
					Messagef(h.p, "team.area.already-claimed")
					return
				}
				if areaTooClose(c, vec3ToVec2(pos.Vec3()), 1) {
					Messagef(h.p, "team.area.too-close")
					return
				}
			}

			pn := 1
			if h.p.Sneaking() {
				pn = 2
				ar := area.NewArea(h.claimPos[0], mgl64.Vec2{float64(pos.X()), float64(pos.Z())})
				x := ar.Max().X() - ar.Min().X()
				y := ar.Max().Y() - ar.Min().Y()
				a := x * y
				if a > 75*75 {
					Messagef(h.p, "team.claim.too-big")
					return
				}
				cost := int(a * 5)
				Messagef(h.p, "team.claim.cost", cost)
			}
			h.claimPos[pn-1] = mgl64.Vec2{float64(pos.X()), float64(pos.Z())}
			Messagef(h.p, "team.claim.set-position", pn, mgl64.Vec2{float64(pos.X()), float64(pos.Z())})
		}
	}

	switch b.(type) {
	case block.WoodFenceGate, block.Chest:
		teams, _ := data.LoadAllTeams()
		for _, t := range teams {
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

		title := strings.ToLower(colour.StripMinecraftColour(lines[0]))
		if strings.Contains(title, "[buy]") ||
			strings.Contains(title, "[sell]") &&
				(area.Spawn(h.p.World()).Vec3WithinOrEqualFloorXZ(h.p.Position())) {
			itm, ok := world.ItemByName("minecraft:"+strings.ReplaceAll(strings.ToLower(lines[1]), " ", "_"), 0)
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

			u, err := data.LoadUserFromName(h.p.Name())
			if err != nil {
				return
			}
			switch choice {
			case "buy":
				if u.Teams.Balance < price {
					h.p.Message("shop.balance.insufficient")
					return
				}
				if !ok {
					return
				}
				// restart: should we do this?
				u.Teams.Balance = u.Teams.Balance - price
				data.SaveUser(u)
				it.AddOrDrop(h.p, item.NewStack(itm, q))
				Messagef(h.p, "shop.buy.success", q, lines[1])
			case "sell":
				inv := h.Player().Inventory()
				count := 0
				var items []item.Stack
				for _, slotItem := range inv.Slots() {
					n1, _ := itm.EncodeItem()
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
					u.Teams.Balance = u.Teams.Balance + float64(count/q)*price
					data.SaveUser(u)
					Messagef(h.p, "shop.sell.success", count, lines[1])
				} else {
					Messagef(h.p, "shop.sell.fail")
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
		} /*else if title == "[kit]" {
			key := colour.StripMinecraftColour(lines[1])
			u, err := data.LoadUserFromName(h.p.Name())
			if err != nil {
				return
			}
			cd := u.Teams.Kits.Key(key)
			if cd.Active() {
				Messagef(h.p, "command.kit.cooldown", cd.Remaining().Round(time.Second))
				return
			} else {
				cd.Set(time.Minute)
			}
			switch strings.ToLower(colour.StripMinecraftColour(lines[1])) {
			case "diamond":
				kit2.Apply(kit2.Master{}, h.p)
			case "archer":
				kit2.Apply(kit2.Archer{}, h.p)
			case "bard":
				kit2.Apply(kit2.Bard{}, h.p)
			case "rogue":
				kit2.Apply(kit2.Rogue{}, h.p)
			case "stray":
				kit2.Apply(kit2.Stray{}, h.p)
			case "miner":
				kit2.Apply(kit2.Miner{}, h.p)
			case "builder":
				kit2.Apply(kit2.Builder{}, h.p)
			}
		}*/
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
			h.p.Message(lang.Translatef(data.Language{}, "user.cool-down", "Rogue", cd.Remaining().Seconds()))
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

	if !t.OnGround() {
		max, min := maxMin(t.Position().Y(), h.p.Position().Y())
		if max-min >= 2.5 {
			*height = 0.38 / 1.25
		}
	}

	//u, err := data.LoadUserOrCreate(h.p.Name(), h.p.Handler().(*Handler).XUID())
	target, ok := Lookup(t.Name())
	targetHandler, ok := t.Handler().(*Handler)
	if !ok {
		return
	}
	typ, ok2 := it.SpecialItem(held)
	if ok && ok2 {
		if cd := h.ability; cd.Active() {
			h.p.Message(text.Colourf("<red>You are on Partner Items cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
			return
		}
		switch kind := typ.(type) {
		case it.StormBreakerType:
			if cd := h.abilities.Key(it.StormBreakerType{}); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on storm breaker cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			h.p.World().PlaySound(h.p.Position(), sound.ItemBreak{})
			h.p.World().AddEntity(entity.NewLightning(h.p.Position()))
			h.abilities.Set(it.StormBreakerType{}, time.Minute*2)
			h.ability.Set(time.Second * 10)

			targetArmourHandler, ok := target.Armour().Inventory().Handler().(*ArmourHandler)
			if !ok {
				break
			}

			h.p.SetHeldItems(item.Stack{}, left)
			targetArmourHandler.stormBreak()
		case it.ExoticBoneType:
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on bone cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			targetHandler.AddBoneHit(t)
			if targetHandler.Boned() {
				target.Message(text.Colourf("<red>You have been boned by %s</red>", h.p.Name()))
				h.p.Message(text.Colourf("<green>You have boned %s</green>", t.Name()))
				h.ability.Set(time.Second * 10)
				h.abilities.Set(kind, time.Minute)
				h.p.SetHeldItems(h.SubtractItem(held, 1), left)
				targetHandler.ResetBoneHits(h.p)
			} else {
				h.p.Message(text.Colourf("<green>You have hit %s with a bone %d times</green>", t.Name(), targetHandler.BoneHits(h.p)))
			}
		case it.ScramblerType:
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on scrambler cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			targetHandler.AddScramblerHit(h.p)
			if targetHandler.ScramblerHits(h.p) >= 3 {
				var used []int
				for i := 0; i <= 7; i++ {
					j := rand.Intn(8)
					var u bool
					for _, v := range used {
						if v == j {
							u = true
						}
					}
					if u {
						continue
					}
					used = append(used, j)
					it1, _ := target.Inventory().Item(i)
					it2, _ := target.Inventory().Item(j)
					target.Inventory().SetItem(j, it1)
					target.Inventory().SetItem(i, it2)
				}
				target.Message(text.Colourf("<red>You have been scrambled by %s</red>", h.p.Name()))
				h.p.Message(text.Colourf("<green>You have scrambled %s</green>", t.Name()))
				h.ability.Set(time.Second * 10)
				h.abilities.Set(kind, time.Minute*2)
				targetHandler.ResetScramblerHits(h.p)
				h.p.SetHeldItems(h.SubtractItem(held, 1), left)
			}
		case it.PearlDisablerType:
			if cd := h.abilities.Key(kind); cd.Active() {
				h.p.Message(text.Colourf("<red>You are on pearl disabler cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
				break
			}
			if !targetHandler.PearlDisabled() {
				target.Message(text.Colourf("<red>You have been pearl disabled by %s</red>", h.p.Name()))
				targetHandler.pearl.Set(time.Second * 15)
				targetHandler.TogglePearlDisable()

				h.p.Message(text.Colourf("<green>You have pearl disabled %s</green>", t.Name()))
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

func maxMin(n, n2 float64) (max float64, min float64) {
	if n > n2 {
		return n, n2
	}
	return n2, n
}

func (h *Handler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if h.logger {
		return
	}
	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil {
		return
	}

	p := h.p
	w := p.World()

	if !newPos.ApproxFuncEqual(p.Position(), func(f float64, f2 float64) bool {
		return math.Abs(f-f2) < 0.03
	}) {
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
		mul := area.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
		if mul.Vec3WithinOrEqualFloorXZ(p.Position()) {
			h.sendWall(cubePos, area.Overworld.Spawn().Area, item.ColourRed())
		}
	}

	if u.Frozen {
		ctx.Cancel()
		return
	}

	if u.Teams.PVP.Active() {
		teams, _ := data.LoadAllTeams()
		for _, a := range teams {
			a := a.Claim
			if a != (area.Area{}) && a.Vec3WithinOrEqualXZ(newPos) {
				ctx.Cancel()
				return
			}

			mul := area.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
			if mul.Vec2WithinOrEqualFloor(mgl64.Vec2{p.Position().X(), p.Position().Z()}) && !area.Spawn(p.World()).Vec3WithinOrEqualFloorXZ(newPos) {
				h.sendWall(cubePos, a, item.ColourBlue())
			}
		}

		if newPos.Y() < 0 {
			h.p.Teleport(mgl64.Vec3{0, 100, 0})
		}
	}

	if _, ok := sotw.Running(); ok && u.Teams.SOTW {
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
			switch k {
			case koth.Cosmic:
				if newPos.Y() < 77 || newPos.Y() > 85 {
					if k.StopCapturing(us) {
						Broadcast("koth.not.capturing", r.Color(u.DisplayName), k.Name())
					}
					return
				}
			}

			if u.Teams.PVP.Active() {
				return
			}
			if k.StartCapturing(us) {
				Broadcast("koth.capturing", k.Name(), r.Color(u.DisplayName))
			}
		} else {
			if k.StopCapturing(us) {
				Broadcast("koth.not.capturing", r.Color(u.DisplayName), k.Name())
			}
		}
	}

	var areas []area.NamedArea

	teams, err := data.LoadAllTeams()
	if err != nil {
		fmt.Println(err)
	}
	for _, tm := range teams {
		a := tm.Claim

		name := text.Colourf("<red>%s</red>", tm.DisplayName)
		if t, err := data.LoadTeamFromMemberName(h.p.Name()); err == nil && t.Name == tm.Name {
			name = text.Colourf("<green>%s</green>", tm.DisplayName)
		}
		areas = append(areas, area.NewNamedArea(mgl64.Vec2{a.Min().X(), a.Min().Y()}, mgl64.Vec2{a.Max().X(), a.Max().Y()}, name))
	}

	ar := h.area.Load()
	for _, a := range append(area.Protected(w), areas...) {
		if a.Vec3WithinOrEqualFloorXZ(newPos) {
			if ar != a {
				if ar != (area.NamedArea{}) {
					Messagef(h.p, "area.leave", ar.Name())
				}
				h.area.Store(a)
				Messagef(h.p, "area.enter", a.Name())
				return
			} else {
				return
			}
		}
	}

	if ar != area.Wilderness(w) {
		if ar != (area.NamedArea{}) {
			Messagef(h.p, "area.leave", ar.Name())

		}
		h.area.Store(area.Wilderness(w))
		Messagef(h.p, "area.enter", area.Wilderness(w).Name())
	}
}

func (h *Handler) HandleItemDrop(ctx *event.Context, e world.Entity) {
	w := h.p.World()
	if h.area.Load() == area.Spawn(w) {
		for _, ent := range w.Entities() {
			if p, ok := ent.(*player.Player); ok {
				p.HideEntity(e)
			}
		}
	}
}

func (*Handler) HandleItemPickup(ctx *event.Context, i *item.Stack) {
	for _, sp := range it.SpecialItems() {
		if _, ok := i.Value(sp.Key()); ok {
			*i = it.NewSpecialItem(sp, i.Count())
		}
	}
}

func (h *Handler) HandleQuit() {
	if h.logger {
		return
	}
	h.close <- struct{}{}
	p := h.p

	u, _ := data.LoadUserFromName(p.Name())
	//u.PlayTime += time.Since(h.logTime)
	data.SaveUser(u)

	tm, _ := data.LoadTeamFromMemberName(p.Name())
	_, sotwRunning := sotw.Running()
	if !h.loggedOut && !tm.Claim.Vec3WithinOrEqualFloorXZ(p.Position()) && !area.Spawn(p.World()).Vec3WithinOrEqualFloorXZ(p.Position()) || ((sotwRunning && u.Teams.SOTW) || u.Teams.PVP.Active()) {
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
				u, err := data.LoadUserFromName(h.p.Name())
				if err != nil {
					return
				}
				u.Teams.Dead = true
				u.Teams.Stats.Deaths = 0
				if u.Teams.Stats.KillStreak > u.Teams.Stats.BestKillStreak {
					u.Teams.Stats.BestKillStreak = u.Teams.Stats.KillStreak
				}
				u.Teams.Stats.KillStreak = 0
				if tm, err := data.LoadTeamFromMemberName(h.p.Name()); err == nil {
					tm = tm.WithDTR(tm.DTR - 1).WithPoints(tm.Points - 1).WithRegenerationTime(time.Now().Add(time.Minute * 5))
					data.SaveTeam(tm)
				}
				DropContents(h.p)
				data.SaveUser(u)
			}
			_ = h.p.Close()
		}()
		h.logger = true
		UpdateState(h.p)

		loggers[p.XUID()] = h
		return
	}
}

func Online(p *player.Player) bool {
	return unsafe.Session(p) != session.Nop
}

func Broadcastf(key string, a ...interface{}) {
	for _, p := range moyai.Server().Players() {
		Messagef(p, key, a...)
	}
}

func Messagef(p *player.Player, key string, a ...interface{}) {
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		p.Message("An error occurred while loading your user data.")
		return
	}
	p.Message(lang.Translatef(u.Language, key, a...))
}

func UpdateState(p *player.Player) {
	for _, v := range p.World().Viewers(p.Position()) {
		v.ViewEntityState(p)
	}
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

package user

import (
	"fmt"
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
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

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

	itemUse   moose.MappedCoolDown[world.Item]
	bardItem  moose.MappedCoolDown[world.Item]
	strayItem moose.MappedCoolDown[world.Item]

	ability   *moose.CoolDown
	abilities moose.MappedCoolDown[it.SpecialItemType]

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

	close chan struct{}
}

func NewHandler(p *player.Player, xuid string) *Handler {
	ha := &Handler{
		p:    p,
		xuid: xuid,

		pearl:       moose.NewCoolDown(),
		rogue:       moose.NewCoolDown(),
		goldenApple: moose.NewCoolDown(),
		ability:     moose.NewCoolDown(),

		bone:     moose.NewCoolDown(),
		boneHits: map[string]int{},

		scramblerHits: map[string]int{},

		wallBlocks: map[cube.Pos]float64{},

		itemUse:   moose.NewMappedCoolDown[world.Item](),
		bardItem:  moose.NewMappedCoolDown[world.Item](),
		strayItem: moose.NewMappedCoolDown[world.Item](),
		abilities: moose.NewMappedCoolDown[it.SpecialItemType](),

		combat: moose.NewTag(nil, nil),
		archer: moose.NewTag(nil, nil),

		home: moose.NewTeleportation(func(t *moose.Teleportation) {
			p.Message(text.Colourf("<green>You have been teleported home.</green>"))
		}),
		logout: moose.NewTeleportation(func(t *moose.Teleportation) {
			p.Disconnect(text.Colourf("<red>You have been logged out.</red>"))
		}),
		stuck: moose.NewTeleportation(func(t *moose.Teleportation) {
			p.Message(text.Colourf("<red>You have been teleported to a safe place.</red>"))
		}),

		close: make(chan struct{}, 0),
	}

	s := player_session(p)
	u, _ := data.LoadUser(p.Name())

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

	go startTicker(ha)
	return ha
}

// formatRegex is a regex used to clean color formatting on a string.
var formatRegex = regexp.MustCompile(`§[\da-gk-or]`)

// HandleChat ...
func (h *Handler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()
	u, err := data.LoadUser(h.p.Name())
	if err != nil {
		return
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
}

// HandleItemUse ...
func (h *Handler) HandleItemUse(ctx *event.Context) {
	held, left := h.p.HeldItems()
	u, _ := data.LoadUser(h.p.Name())
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
			if u.PVP.Active() || u.SOTW {
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
			h.abilities.Key(kind).Set(time.Second * 10)
		case it.FullInvisibilityType:
			// Restart TODO
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

	u, err := data.LoadUser(h.p.Name())
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
	}
}

func (h *Handler) HandleHurt(ctx *event.Context, dmg *float64, imm *time.Duration, src world.DamageSource) {
	if area.Spawn(h.p.World()).Vec3WithinOrEqualFloorXZ(h.p.Position()) {
		ctx.Cancel()
		return
	}

	var target *player.Player
	switch s := src.(type) {
	case NoArmourAttackEntitySource:
		if t, ok := s.Attacker.(*player.Player); ok {
			target = t
		}
		fmt.Println("hurt")
		if !canAttack(h.p, target) {
			fmt.Println("hurt")
			ctx.Cancel()
			return
		}
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
		if ha := target.Handler().(*Handler); !class.Compare(h.class.Load(), class.Archer{}) && class.Compare(ha.class.Load(), class.Archer{}) && (s.Projectile.Type() == (entity.ArrowType{})) {
			ctx.Cancel()
			h.archer.Set(time.Second * 10)

			dist := h.p.Position().Sub(target.Position()).Len()
			d := math.Round(dist)
			if d > 20 {
				d = 20
			}
			damage := (d / 10) * 2
			h.p.Hurt(damage, NoArmourAttackEntitySource{
				Attacker: h.p,
			})
			h.p.KnockBack(target.Position(), 0.394, 0.394)

			target.Message(lang.Translatef(h.p.Locale(), "archer.tag", math.Round(dist), damage/2))
		}

		if s.Projectile.Type() == (it.SwitcherBallType{}) {
			if k, ok := koth.Running(); ok {
				if pl, ok := k.Capturing(); ok && pl.Player() == h.p {
					target.Message(text.Colourf("<red>You cannot switch places with someone capturing a koth</red>"))
					break
				}
			}

			dist := h.p.Position().Sub(target.Position()).Len()
			if dist > 10 {
				target.Message(text.Colourf("<red>You are too far away from %s</red>", h.p.Name()))
				break
			}

			ctx.Cancel()
			targetPos := target.Position()
			pos := h.p.Position()

			target.Teleport(pos)
			h.p.Teleport(targetPos)
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

		u, err := data.LoadUser(h.p.Name())
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
		h.p.Teleport(mgl64.Vec3{0, 67, 0})
		h.p.SetMobile()

		if tm, ok := u.Team(); ok {
			tm = tm.WithDTR(tm.DTR - 1).WithPoints(tm.Points - 1).WithRegenerationTime(time.Now().Add(time.Minute * 15))
			defer data.SaveTeam(tm)
		}

		killer, ok := h.LastAttacker()
		if ok {
			k, err := data.LoadUser(killer.p.Name())
			if err != nil {
				return
			}
			k.Stats.Kills += 1

			if tm, ok := k.Team(); ok {
				tm.WithPoints(tm.Points + 1)
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
	}

	if canAttack(h.p, target) {
		target.Handler().(*Handler).combat.Set(time.Second * 20)
		h.combat.Set(time.Second * 20)
	}
}

func (h *Handler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	w := h.p.World()

	switch b.(type) {
	case block.EnderChest:
		held, left := h.p.HeldItems()
		if _, ok := held.Value("PARTNER_PACKAGE"); !ok {
			break
		}

		keys := it.SpecialItems()
		i := it.NewSpecialItem(keys[rand.Intn(len(keys))], rand.Intn(3)+1)

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
		u, _ := data.LoadUser(h.p.Name())
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

	for _, t := range data.Teams() {
		if !t.Member(h.p.Name()) {
			if t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
		for _, a := range area.Protected(w) {
			if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
				if h.p.GameMode() != world.GameModeCreative {
					ctx.Cancel()
					return
				}
			}
		}
	}
}

func (h *Handler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	w := h.p.World()

	i, _ := h.p.HeldItems()
	if _, ok := i.Item().(item.Bucket); ok {
		for _, a := range area.Protected(w) {
			if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
	}

	u, _ := data.LoadUser(h.p.Name())
	switch it := i.Item().(type) {
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
	b := w.Block(pos)

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

			u, err := data.LoadUser(h.p.Name())
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
		}
	}
}

func (h *Handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, _ *bool) {
	*force, *height = 0.394, 0.394
	t, ok := e.(*player.Player)
	if !ok {
		return
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

	//u, err := data.LoadUser(h.p.Name(), h.p.Handler().(*Handler).XUID())
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
				inv := target.Player().Inventory()
				for i := 36; i <= 44; i++ {
					j := rand.Intn(i+1-36) + 36
					it1, _ := inv.Item(i)
					it2, _ := inv.Item(j)
					inv.SetItem(i, it1)
					inv.SetItem(j, it2)
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
	if ok {
		th.SetLastAttacker(h)
		th.combat.Set(time.Second * 20)
	}
	h.combat.Set(time.Second * 20)
}

func (h *Handler) HandleQuit() {
	close(h.close)
	p := h.p

	u, _ := data.LoadUser(p.Name())
	u.PlayTime += time.Since(h.logTime)
	_ = data.SaveUser(u)

	playersMu.Lock()
	delete(players, h.p.Name())
	playersMu.Unlock()
}

func (h *Handler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	u, _ := data.LoadUser(h.p.Name())
	p := h.p
	w := p.World()

	if !newPos.ApproxEqual(p.Position()) {
		h.home.Cancel()
		h.logout.Cancel()
		h.stuck.Cancel()
	}

	h.clearWall()
	cubePos := cube.PosFromVec3(newPos)

	if h.combat.Active() {
		a := area.Spawn(w)
		mul := moose.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
		if mul.Vec2WithinOrEqualFloor(mgl64.Vec2{p.Position().X(), p.Position().Z()}) {
			h.sendWall(cubePos, area.Overworld.Spawn().Area, item.ColourRed())
		}
	}

	if u.PVP.Active() {
		for _, a := range data.Teams() {
			a := a.Claim
			if a.Vec3WithinOrEqualXZ(newPos) {
				ctx.Cancel()
				return
			}

			mul := moose.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
			if mul.Vec2WithinOrEqualFloor(mgl64.Vec2{p.Position().X(), p.Position().Z()}) && !area.Spawn(p.World()).Vec3WithinOrEqualFloorXZ(newPos) {
				h.sendWall(cubePos, a, item.ColourBlue())
			}
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
		areas = append(areas, moose.NewNamedArea(mgl64.Vec2{a.Min().X(), a.Max().X()}, mgl64.Vec2{a.Min().Y(), a.Max().Y()}, name))
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

package user

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/restartfu/roman"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

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

	proxypacket "github.com/paroxity/portal/socket/packet"
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
	s       *session.Session
	p       *player.Player
	xuid    string
	logTime time.Time

	pearl       *moose.CoolDown
	rogue       *moose.CoolDown
	goldenApple *moose.CoolDown

	itemUse   moose.MappedCoolDown[world.Item]
	bardItem  moose.MappedCoolDown[world.Item]
	strayItem moose.MappedCoolDown[world.Item]

	abilities moose.MappedCoolDown[any]

	armour atomic.Value[[4]item.Stack]
	class  atomic.Value[moose.Class]
	energy atomic.Value[float64]

	combat *moose.Tag
	archer *moose.Tag

	logout *moose.Teleportation
	stuck  *moose.Teleportation
	home   *moose.Teleportation

	lastScoreBoard atomic.Value[*scoreboard.Scoreboard]
	area           atomic.Value[moose.NamedArea]

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

		itemUse:   moose.NewMappedCoolDown[world.Item](),
		bardItem:  moose.NewMappedCoolDown[world.Item](),
		strayItem: moose.NewMappedCoolDown[world.Item](),

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
	u, _ := data.LoadUser(p.Name(), xuid)

	u.DisplayName = p.Name()
	u.Name = strings.ToLower(p.Name())
	u.XUID = p.XUID()

	u.DeviceID = s.ClientData().DeviceID
	u.SelfSignedID = s.ClientData().SelfSignedID
	if err := data.SaveUser(u); err != nil {
		panic(err)
	}

	ha.s = s
	ha.logTime = time.Now()

	playersMu.Lock()
	players[p.XUID()] = ha
	playersXUID[p.Name()] = xuid
	playersMu.Unlock()

	go startTicker(ha)

	return ha
}

// formatRegex is a regex used to clean color formatting on a string.
var formatRegex = regexp.MustCompile(`§[\da-gk-or]`)

// HandleChat ...
func (h *Handler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()
	if *message == "allah" {
		moyai.Socket().WritePacket(&proxypacket.TransferRequest{
			PlayerUUID: h.Player().UUID(),
			Server:     "syn.lobby",
		})
		moyai.SocketUnlock()
		fmt.Println("aaa")
		return
	}
	u, err := data.LoadUser(h.p.Name(), h.p.Handler().(*Handler).XUID())
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
				if slices.ContainsFunc(tm.Members, func(member data.Member) bool {
					return member.XUID == t.p.Handler().(*Handler).XUID()
				}) {
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

// HandleItemUse ...
func (h *Handler) HandleItemUse(ctx *event.Context) {
	held, _ := h.p.HeldItems()
	switch held.Item().(type) {
	case item.EnderPearl:
		if cd := h.pearl; cd.Active() {
			h.p.Message(lang.Translatef(h.p.Locale(), "user.cool-down", "Ender Pearl", cd.Remaining().Seconds()))
			ctx.Cancel()
		} else {
			cd.Set(15 * time.Second)
		}
	}

	u, _ := data.LoadUser(h.p.Name(), h.p.Handler().(*Handler).XUID())
	switch h.class.Load().(type) {
	case class.Bard:
		if e, ok := BardEffectFromItem(held.Item()); ok {
			if u.PVP.Active() /*|| u.SOTW*/ {
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
			}

			lvl, _ := roman.Itor(e.Level())
			h.Message("class.ability.use", moose.EffectName(e), lvl, len(teammates))
			h.p.SetHeldItems(held.Grow(-1), item.Stack{})
			h.bardItem.Key(held.Item()).Set(15 * time.Second)
		}
	case class.Stray:
		if e, ok := StrayEffectFromItem(held.Item()); ok {
			if u.PVP.Active() /*|| u.SOTW*/ {
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
			}

			lvl, _ := roman.Itor(e.Level())
			h.Message("class.ability.use", moose.EffectName(e), lvl, len(teammates))
			h.p.SetHeldItems(held.Grow(-1), item.Stack{})
			h.strayItem.Key(held.Item()).Set(15 * time.Second)
		}
	}
}

func (h *Handler) HandleHurt(ctx *event.Context, dmg *float64, _ *time.Duration, src world.DamageSource) {
	var target *player.Player
	switch s := src.(type) {
	case NoArmourAttackEntitySource:
		if t, ok := s.Attacker.(*player.Player); ok {
			target = t
		}
		if !canAttack(h.p, target) {
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
			ha.archer.Set(time.Second * 10)

			dist := h.p.Position().Sub(target.Position()).Len()
			d := math.Round(dist)
			if d > 20 {
				d = 20
			}
			*dmg = (d / 10) * 2
			h.p.Hurt(*dmg, NoArmourAttackEntitySource{
				Attacker: h.p,
			})
			h.p.KnockBack(target.Position(), 0.394, 0.394)

			target.Message(lang.Translatef(h.p.Locale(), "archer.tag", math.Round(dist), *dmg/2))
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
		if !slices.ContainsFunc(t.Members, func(member data.Member) bool {
			return member.XUID == h.p.Handler().(*Handler).XUID()
		}) {
			if t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
		u, _ := data.LoadUser(h.p.Name(), h.p.Handler().(*Handler).XUID())
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
		if !slices.ContainsFunc(t.Members, func(member data.Member) bool {
			return member.XUID == h.p.Handler().(*Handler).XUID()
		}) {
			if t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
		u, _ := data.LoadUser(h.p.Name(), h.p.Handler().(*Handler).XUID())
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

	u, _ := data.LoadUser(h.p.Name(), h.p.Handler().(*Handler).XUID())
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

	t.Handler().(*Handler).combat.Set(time.Second * 20)
	h.combat.Set(time.Second * 20)
}

func (h *Handler) HandleQuit() {
	logrus.Info("DIsconnected")
	h.close <- struct{}{}
	p := h.p

	u, _ := data.LoadUser(p.Name(), playersXUID[p.Name()])
	u.PlayTime += time.Since(h.logTime)
	_ = data.SaveUser(u)

	playersMu.Lock()
	delete(players, h.p.Name())
	playersMu.Unlock()
}

func (h *Handler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	u, _ := data.LoadUser(h.p.Name(), h.p.Handler().(*Handler).XUID())
	p := h.p
	w := p.World()

	if !newPos.ApproxEqual(p.Position()) {
		h.home.Cancel()
		h.logout.Cancel()
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
			if mul.Vec2WithinOrEqualFloor(mgl64.Vec2{p.Position().X(), p.Position().Z()}) {
				h.sendWall(cubePos, a, item.ColourBlue())
			}
		}
	}

	if area.Spawn(w).Vec3WithinOrEqualFloorXZ(newPos) && h.combat.Active() {
		ctx.Cancel()
		return
	}

	/*k, ok := koth.Running()
	if ok {
		r := u.Roles().Highest()
		if k.Area().Vec3WithinOrEqualFloorXZ(newPos) {
			if k.StartCapturing(u, r.Colour(u.Name())) {
				user.Broadcast("koth.capturing", k.Name(), r.Colour(u.Name()))
			}
		} else {
			if k.StopCapturing(u) {
				user.Broadcast("koth.not.capturing", k.Name())
			}
		}
	}*/

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

func (h *Handler) Player() *player.Player {
	return h.p
}

func (h *Handler) XUID() string {
	return h.xuid
}

func (h *Handler) Message(key string, args ...interface{}) {
	h.p.Message(lang.Translatef(h.p.Locale(), key, args...))
}

func (h *Handler) Logout() *moose.Teleportation {
	return h.logout
}

// AddItemOrDrop adds an item to the user's inventory or drops it if the inventory is full.
func (h *Handler) AddItemOrDrop(it item.Stack) {
	if _, err := h.p.Inventory().AddItem(it); err != nil {
		h.DropItem(it)
	}
}

func (h *Handler) Combat() *moose.Tag {
	return h.combat
}

func (h *Handler) Pearl() *moose.CoolDown {
	return h.pearl
}

func (h *Handler) DropItem(it item.Stack) {
	p := h.p
	w, pos := p.World(), p.Position()
	ent := entity.NewItem(it, pos)
	ent.SetVelocity(mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1})
	w.AddEntity(ent)
}

func (h *Handler) sendWall(newPos cube.Pos, z moose.Area, color item.Colour) {
	areaMin := cube.Pos{int(z.Min().X()), 0, int(z.Min().Y())}
	areaMax := cube.Pos{int(z.Max().X()), 255, int(z.Max().Y())}
	wallBlock := block.StainedGlass{Colour: color}
	const wallLength, wallHeight = 15, 10

	if newPos.X() >= areaMin.X() && newPos.X() <= areaMax.X() { // edges of the top and bottom walls (relative to South)
		zCoord := areaMin.Z()
		if newPos.Z() >= areaMax.Z() {
			zCoord = areaMax.Z()
		}
		for horizontal := newPos.X() - wallLength; horizontal < newPos.X()+wallLength; horizontal++ {
			if horizontal >= areaMin.X() && horizontal <= areaMax.X() {
				for vertical := newPos.Y(); vertical < newPos.Y()+wallHeight; vertical++ {
					blockPos := cube.Pos{horizontal, vertical, zCoord}
					if blockReplaceable(h.p.World().Block(blockPos)) {
						h.s.ViewBlockUpdate(blockPos, wallBlock, 0)
						h.wallBlocksMu.Lock()
						h.wallBlocks[blockPos] = rand.Float64() * float64(rand.Intn(1)+1)
						h.wallBlocksMu.Unlock()
					}
				}
			}
		}
	}
	if newPos.Z() >= areaMin.Z() && newPos.Z() <= areaMax.Z() { // edges of the left and right walls (relative to South)
		xCoord := areaMin.X()
		if newPos.X() >= areaMax.X() {
			xCoord = areaMax.X()
		}
		for horizontal := newPos.Z() - wallLength; horizontal < newPos.Z()+wallLength; horizontal++ {
			if horizontal >= areaMin.Z() && horizontal <= areaMax.Z() {
				for vertical := newPos.Y(); vertical < newPos.Y()+wallHeight; vertical++ {
					blockPos := cube.Pos{xCoord, vertical, horizontal}
					if blockReplaceable(h.p.World().Block(blockPos)) {
						h.s.ViewBlockUpdate(blockPos, wallBlock, 0)
						h.wallBlocksMu.Lock()
						h.wallBlocks[blockPos] = rand.Float64() * float64(rand.Intn(1)+1)
						h.wallBlocksMu.Unlock()
					}
				}
			}
		}
	}
}

// clearWall clears the users walls or lowers the remaining duration.
func (h *Handler) clearWall() {
	h.wallBlocksMu.Lock()
	for p, duration := range h.wallBlocks {
		if duration <= 0 {
			delete(h.wallBlocks, p)
			h.s.ViewBlockUpdate(p, block.Air{}, 0)
			h.p.ShowParticle(p.Vec3(), particle.BlockForceField{})
			continue
		}
		h.wallBlocks[p] = duration - 0.1
	}
	h.wallBlocksMu.Unlock()
}

// blockReplaceable checks if the combat wall should replace a block.
func blockReplaceable(b world.Block) bool {
	_, air := b.(block.Air)
	_, doubleFlower := b.(block.DoubleFlower)
	_, flower := b.(block.Flower)
	_, tallGrass := b.(block.TallGrass)
	_, doubleTallGrass := b.(block.DoubleTallGrass)
	_, deadBush := b.(block.DeadBush)
	//_, cobweb := b.(block.Cobweb)
	//_, sapling := b.(block.Sapling)
	_, torch := b.(block.Torch)
	_, fire := b.(block.Fire)
	return air || tallGrass || deadBush || torch || fire || flower || doubleFlower || doubleTallGrass
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

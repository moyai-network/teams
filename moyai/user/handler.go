package user

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/google/uuid"

	"github.com/moyai-network/teams/moyai/roles"

	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/moyai-network/teams/moyai"

	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal/cooldown"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/internal/unsafe"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/class"
	"github.com/moyai-network/teams/moyai/conquest"
	"github.com/moyai-network/teams/moyai/crate"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/eotw"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/koth"
	"github.com/moyai-network/teams/moyai/process"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
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

	loggers  = map[string]*Handler{}
	loggerMu sync.Mutex
)

type Handler struct {
	player.NopHandler

	p    *player.Player
	uuid uuid.UUID

	logTime           time.Time
	claimSelectionPos [2]mgl64.Vec3
	waypoint          *WayPoint
	energy            atomic.Value[float64]

	wallBlocks   map[cube.Pos]float64
	wallBlocksMu sync.Mutex

	lastPlacedSignPos cube.Pos
	lastArmour        atomic.Value[[4]item.Stack]
	lastClass         atomic.Value[class.Class]
	lastScoreBoard    atomic.Value[*scoreboard.Scoreboard]
	lastArea          atomic.Value[area.NamedArea]
	lastAttackerName  atomic.Value[string]
	lastAttackTime    atomic.Value[time.Time]
	lastPearlPos      mgl64.Vec3
	lastMessage       atomic.Value[time.Time]

	tagCombat                 *cooldown.CoolDown
	tagArcher                 *cooldown.CoolDown
	coolDownComboAbility      *cooldown.CoolDown
	coolDownVampireAbility    *cooldown.CoolDown
	coolDownBonedEffect       *cooldown.CoolDown
	coolDownEffectDisabled    *cooldown.CoolDown
	coolDownFocusMode         *cooldown.CoolDown
	coolDownPearl             *cooldown.CoolDown
	coolDownBackStab          *cooldown.CoolDown
	coolDownGoldenApple       *cooldown.CoolDown
	coolDownItemUse           *cooldown.CoolDown
	coolDownArcherRogueItem   cooldown.MappedCoolDown[world.Item]
	coolDownBardItem          cooldown.MappedCoolDown[world.Item]
	coolDownMageItem          cooldown.MappedCoolDown[world.Item]
	coolDownGlobalAbilities   *cooldown.CoolDown
	coolDownSpecificAbilities cooldown.MappedCoolDown[it.SpecialItemType]
	processLogout             *process.Process
	processStuck              *process.Process
	processHome               *process.Process
	processCamp               *process.Process

	gracefulLogout bool
	logger         bool

	close chan struct{}
	death chan struct{}
}

func NewHandler(p *player.Player, xuid string) *Handler {
	if h, ok := logger(p); ok {
		if h.p.World().Dimension() == world.End {
			moyai.End().AddEntity(p)
			unsafe.WritePacket(p, &packet.PlayerFog{
				Stack: []string{"minecraft:fog_the_end"},
			})
			<-time.After(time.Second)
		} else if h.p.World().Dimension() == world.Nether {
			unsafe.WritePacket(p, &packet.PlayerFog{
				Stack: []string{"minecraft:fog_hell"},
			})
			moyai.Nether().AddEntity(p)
			<-time.After(time.Second)
		} else if h.p.World() == moyai.Deathban() {
			<-time.After(time.Second)
			moyai.Deathban().AddEntity(p)
			p.SetGameMode(world.GameModeSurvival)
		}
		p.Teleport(h.p.Position())
		currentHealth := h.p.Health()
		p.Hurt(20-currentHealth, NoArmourAttackEntitySource{})
		_ = h.p.Close()
	}

	if p.World().Dimension() == world.End {
		unsafe.WritePacket(p, &packet.PlayerFog{
			Stack: []string{"minecraft:fog_the_end"},
		})
		moyai.End().AddEntity(p)
	} else if p.World().Dimension() == world.Nether {
		unsafe.WritePacket(p, &packet.PlayerFog{
			Stack: []string{"minecraft:fog_hell"},
		})
		moyai.Nether().AddEntity(p)
	}

	h := &Handler{
		p:          p,
		uuid:       p.UUID(),
		wallBlocks: map[cube.Pos]float64{},

		tagCombat:                 cooldown.NewCoolDown(),
		tagArcher:                 cooldown.NewCoolDown(),
		coolDownPearl:             cooldown.NewCoolDown(),
		coolDownBackStab:          cooldown.NewCoolDown(),
		coolDownGoldenApple:       cooldown.NewCoolDown(),
		coolDownGlobalAbilities:   cooldown.NewCoolDown(),
		coolDownBonedEffect:       cooldown.NewCoolDown(),
		coolDownEffectDisabled:    cooldown.NewCoolDown(),
		coolDownFocusMode:         cooldown.NewCoolDown(),
		coolDownItemUse:           cooldown.NewCoolDown(),
		coolDownComboAbility:      cooldown.NewCoolDown(),
		coolDownVampireAbility:    cooldown.NewCoolDown(),
		coolDownArcherRogueItem:   cooldown.NewMappedCoolDown[world.Item](),
		coolDownBardItem:          cooldown.NewMappedCoolDown[world.Item](),
		coolDownMageItem:          cooldown.NewMappedCoolDown[world.Item](),
		coolDownSpecificAbilities: cooldown.NewMappedCoolDown[it.SpecialItemType](),
		processHome: process.NewProcess(func(t *process.Process) {
			p.Message(text.Colourf("<green>You have been teleported home.</green>"))
		}),
		processStuck: process.NewProcess(func(t *process.Process) {
			p.Message(text.Colourf("<red>You have been teleported to a safe place.</red>"))
		}),

		close: make(chan struct{}),
		death: make(chan struct{}),
	}

	h.processLogout = process.NewProcess(func(t *process.Process) {
		h.gracefulLogout = true
		p.Disconnect(text.Colourf("<red>You have been logged out.</red>"))
	})

	p.SetNameTag(text.Colourf("<red>%s</red>", p.Name()))
	UpdateState(p)

	s := unsafe.Session(p)
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return nil
	}

	u.StaffMode = false
	if u.Teams.DeathBan.Active() {
		moyai.Deathban().AddEntity(p)
		p.Teleport(mgl64.Vec3{5, 13, 44})
	} else {
		if u.Teams.DeathBanned {
			u.Teams.DeathBanned = false
			u.Teams.PVP.Set(time.Hour + (time.Millisecond * 500))
			if !u.Teams.PVP.Paused() {
				u.Teams.PVP.TogglePause()
			}

			moyai.Overworld().AddEntity(p)
			p.Teleport(mgl64.Vec3{0, 80, 0})
		}
	}

	u.DisplayName = p.Name()
	u.Name = strings.ToLower(p.Name())
	u.XUID = xuid
	u.DeviceID = s.ClientData().DeviceID
	u.SelfSignedID = s.ClientData().SelfSignedID

	if !u.Roles.Contains(roles.Default()) {
		u.Roles.Add(roles.Default())
	}

	u.Roles.Add(roles.Pharaoh())
	p.Message(lang.Translatef(*u.Language, "discord.message"))
	h.handleBoosterRole(u)

	data.SaveUser(u)
	if u.Frozen {
		p.SetImmobile()
	}

	h.updateCurrentArea(p.Position(), u)
	h.updateKOTHState(p.Position(), u)
	UpdateVanishState(p, u)

	h.logTime = time.Now()
	UpdateState(h.p)
	go startTicker(h)
	return h
}

func (h *Handler) handleBoosterRole(u data.User) {
	p := h.p

	if len(u.DiscordID) > 0 {
		userID, _ := strconv.Atoi(u.DiscordID)
		rl, err := moyai.DiscordState().MemberRoles(discord.GuildID(1111055709300342826), discord.UserID(userID))
		if err == nil && slices.ContainsFunc(rl, func(d discord.Role) bool {
			return discord.RoleID(1113243316805447830) == d.ID
		}) {
			{
				p.Message(text.Colourf("<green>Thank you for being a Nitro Booster!</green>"))
				u.Roles.Add(roles.Nitro())
				return
			}
		}
	}
	if u.Roles.Contains(roles.Nitro()) {
		p.Message(text.Colourf("<red>You are no longer a Nitro Booster.</red>"))
		u.Roles.Remove(roles.Nitro())
	}
}

func (h *Handler) HandleFoodLoss(ctx *event.Context, _ int, _ *int) {
	ctx.Cancel()
}

func (h *Handler) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	p := h.p
	u, err := data.LoadUserFromXUID(h.p.XUID())
	if err != nil {
		return
	}

	w := p.World()
	b := w.Block(pos)

	if _, ok := b.(block.ItemFrame); ok {
		for _, a := range area.Protected(w) {
			if a.Vec3WithinOrEqualXZ(pos.Vec3()) && h.p.GameMode() != world.GameModeCreative {
				ctx.Cancel()
				return
			}
		}
	}

	held, _ := p.HeldItems()
	typ, ok := it.PartnerItem(held)
	if ok {
		if cd := h.coolDownGlobalAbilities; cd.Active() {
			p.SendJukeboxPopup(lang.Translatef(*u.Language, "popup.cooldown.partner_item", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		if spi := h.coolDownSpecificAbilities; spi.Active(typ) {
			p.SendJukeboxPopup(lang.Translatef(*u.Language, "popup.cooldown.partner_item.item", typ.Name(), spi.Remaining(typ).Seconds()))
			ctx.Cancel()
		} else {
			p.SendJukeboxPopup(lang.Translatef(*u.Language, "popup.ready.partner_item.item", typ.Name()))
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

func vec3ToVec2(v mgl64.Vec3) mgl64.Vec2 {
	return mgl64.Vec2{v.X(), v.Z()}
}

func (h *Handler) HandleSignEdit(ctx *event.Context, frontSide bool, oldText, newText string) {
	//ctx.Cancel()
	if !frontSide {
		return
	}

	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil {
		return
	}

	teams, _ := data.LoadAllTeams()
	if posWithinProtectedArea(h.p, h.lastPlacedSignPos, teams) {
		return
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
			b := h.p.World().Block(h.lastPlacedSignPos)
			if s, ok := b.(block.Sign); ok {
				h.p.World().SetBlock(h.lastPlacedSignPos, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
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

		if !u.Roles.Contains(roles.Admin()) {
			h.p.World().SetBlock(h.lastPlacedSignPos, block.Air{}, nil)
			return
		}

		var newLines []string
		spl := strings.Split(lines[1], " ")
		if len(spl) < 2 {
			return
		}
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
			b := h.p.World().Block(h.lastPlacedSignPos)
			if s, ok := b.(block.Sign); ok {
				h.p.World().SetBlock(h.lastPlacedSignPos, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
	case "[kit]":
		return // disabled for HCF
		if len(lines) < 2 {
			return
		}

		if !u.Roles.Contains(roles.Admin()) {
			h.p.World().SetBlock(h.lastPlacedSignPos, block.Air{}, nil)
			return
		}

		var newLines []string
		newLines = append(newLines, text.Colourf("<dark-red>[Kit]</dark-red>"))
		newLines = append(newLines, text.Colourf("%s", lines[1]))

		time.AfterFunc(time.Millisecond, func() {
			b := h.p.World().Block(h.lastPlacedSignPos)
			if s, ok := b.(block.Sign); ok {
				h.p.World().SetBlock(h.lastPlacedSignPos, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
	}
}

func (h *Handler) HandleHurt(ctx *event.Context, dmg *float64, imm *time.Duration, src world.DamageSource) {
	p := h.p
	*dmg = *dmg / 1.25
	if h.tagArcher.Active() {
		*dmg = *dmg + *dmg*0.25
	} else if h.coolDownFocusMode.Active() &&
		!class.Compare(h.lastClass.Load(), class.Archer{}) &&
		!class.Compare(h.lastClass.Load(), class.Mage{}) &&
		!class.Compare(h.lastClass.Load(), class.Bard{}) &&
		!class.Compare(h.lastClass.Load(), class.Rogue{}) {
		*dmg = *dmg + *dmg*0.25
	}

	u, err := data.LoadUserFromName(h.p.Name())
	if area.Spawn(h.p.World()).Vec3WithinOrEqualFloorXZ(h.p.Position()) && h.p.World() != moyai.Deathban() {
		ctx.Cancel()
		return
	}

	if area.Deathban.Spawn().Vec3WithinOrEqualFloorXZ(h.p.Position()) {
		ctx.Cancel()
		return
	}

	if err != nil || (u.Teams.PVP.Active() && !u.Teams.DeathBan.Active()) {
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
		if err != nil || (u.Teams.PVP.Active() && !u.Teams.DeathBan.Active()) {
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
					moyai.Messagef(attacker, "snowball.koth")
					break
				}
			}

			if ok := conquest.Running(); ok {
				for _, c := range conquest.All() {
					if pl, ok := c.Capturing(); ok && pl == h.p {
						moyai.Messagef(attacker, "snowball.koth")
						break
					}
				}
			}

			dist := attacker.Position().Sub(attacker.Position()).Len()
			if dist > 10 {
				moyai.Messagef(attacker, "snowball.far")
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
			h.setLastAttacker(ha)
			if class.Compare(ha.lastClass.Load(), class.Archer{}) && !class.Compare(h.lastClass.Load(), class.Archer{}) {
				h.tagArcher.Set(time.Second * 10)
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
		h.ShowArmor(true)

		percent := 0.90
		e, ok := attacker.Effect(effect.Strength{})
		if e.Level() > 1 {
			percent = 0.80
		}

		if ok {
			*dmg = *dmg * percent
		}
	}

	if (p.Health()-p.FinalDamageFrom(*dmg, src) <= 0 || (src == entity.VoidDamageSource{})) && !ctx.Cancelled() {
		ctx.Cancel()
		h.kill(src)

		killer, ok := h.lastAttacker()
		if ok {
			k, err := data.LoadUserFromName(killer.Name())
			if err != nil {
				return
			}
			if k.Teams.DeathBan.Active() {
				k.Teams.DeathBan.Set(k.Teams.DeathBan.Remaining() - time.Minute*2)
				return
			}
			k.Teams.Stats.Kills += 1
			k.Teams.Stats.KillStreak += 1

			if k.Teams.Stats.KillStreak%5 == 0 {
				moyai.Broadcastf("user.killstreak", killer.Name(), k.Teams.Stats.KillStreak)
				it.AddOrDrop(killer, it.NewKey(it.KeyTypePartner, int(k.Teams.Stats.KillStreak)/2))
			}

			if tm, err := data.LoadTeamFromMemberName(killer.Name()); err == nil {
				tm = tm.WithPoints(tm.Points + 1)
				if conquest.Running() {
					for _, k := range area.KOTHs(h.p.World()) {
						if k.Name() == "Conquest" && k.Vec3WithinOrEqualXZ(h.p.Position()) {
							conquest.IncreaseTeamPoints(tm, 15)
							if otherTm, err := data.LoadTeamFromMemberName(killer.Name()); err == nil {
								conquest.IncreaseTeamPoints(otherTm, -15)
							}
						}
					}
				}

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
			h.resetLastAttacker()
		} else {
			_, _ = chat.Global.WriteString(lang.Translatef(data.Language{}, "user.suicide", p.Name(), u.Teams.Stats.Kills))
		}
	}

	if canAttack(h.p, attacker) {
		attacker.Handler().(*Handler).tagCombat.Set(time.Second * 20)
		h.tagCombat.Set(time.Second * 20)

		if attacker.Handler().(*Handler).coolDownVampireAbility.Active() {
			attacker.Heal(*dmg*0.5, effect.RegenerationHealingSource{})
		}
	}
}

func (h *Handler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	if h.coolDownBonedEffect.Active() {
		moyai.Messagef(h.p, "bone.interact")
		ctx.Cancel()
		return
	}
	w := h.p.World()
	teams, _ := data.LoadAllTeams()

	if posWithinProtectedArea(h.p, pos, teams) {
		ctx.Cancel()
	}

	switch bl := b.(type) {
	case block.TNT, it.TripwireHook:
		ctx.Cancel()
	case block.Chest:
		for _, dir := range []cube.Direction{bl.Facing.RotateLeft(), bl.Facing.RotateRight()} {
			sidePos := pos.Side(dir.Face())
			for _, t := range teams {
				if !t.Member(h.p.Name()) {
					c := w.Block(sidePos)
					_, eotw := eotw.Running()
					if _, ok := c.(block.Chest); ok && !eotw && t.DTR > 0 && t.Claim.Vec3WithinOrEqualXZ(sidePos.Vec3()) {
						ctx.Cancel()
					}
				}
			}
		}
	case block.Sign:
		h.lastPlacedSignPos = pos
	case block.EnderChest:
		held, _ := h.p.HeldItems()
		if _, ok := held.Value("partner_package"); !ok {
			break
		}
		if typ, ok := it.SpecialItem(held); ok {
			if _, ok := typ.(it.PartnerPackageType); ok {
				ctx.Cancel()
				return
			}
		}
	}
}

func maxMin(n, n2 float64) (max float64, min float64) {
	if n > n2 {
		return n, n2
	}
	return n2, n
}

func (h *Handler) HandleQuit() {
	if h.logger {
		return
	}
	h.close <- struct{}{}
	p := h.p

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	u.PlayTime += time.Since(h.logTime)
	if !u.Teams.PVP.Paused() {
		u.Teams.PVP.TogglePause()
	}
	if u.StaffMode {
		restorePlayerData(h.p)
	}

	data.SaveUser(u)
	data.FlushUser(u)

	_, sotwRunning := sotw.Running()
	if !h.gracefulLogout && h.p.GameMode() != world.GameModeCreative && !u.Teams.PVP.Active() {
		if sotwRunning && u.Teams.SOTW {
			return
		}
		if area.Spawn(p.World()).Vec3WithinOrEqualFloorXZ(p.Position()) && p.World().Dimension() != world.End {
			return
		}
		arm := h.p.Armour()
		inv := h.p.Inventory()

		h.p = player.New(p.Name(), p.Skin(), p.Position())
		rot := p.Rotation()
		unsafe.Rotate(h.p, rot[0], rot[1])
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
				u, err = data.LoadUserFromName(h.p.Name())
				if err != nil {
					return
				}
				data.SaveUser(u)
				data.FlushUser(u)
				break
			case <-h.close:
				break
			case <-h.death:
				break
			}
			_ = h.p.Close()
		}()
		UpdateState(h.p)

		h.p.Handle(h)
		h.p.Armour().Handle(arm.Inventory().Handler())

		setLogger(p, h)
		return
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

func restorePlayerData(p *player.Player) {
	p.Inventory().Handle(inventory.NopHandler{})
	dat, err := moyai.LoadPlayerData(p.UUID())
	if err != nil {
		return
	}
	p.Teleport(dat.Position)
	p.SetGameMode(dat.GameMode)

	newInv := dat.Inventory
	p.Inventory().Clear()
	p.Armour().Clear()
	p.Armour().Set(newInv.Helmet, newInv.Chestplate, newInv.Leggings, newInv.Boots)
	for slot, itm := range newInv.Items {
		if itm.Empty() {
			continue
		}
		_ = p.Inventory().SetItem(slot, itm)
	}
}

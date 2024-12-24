package user

import (
	"github.com/moyai-network/teams/internal/core/area"
	"github.com/moyai-network/teams/internal/core/data"
	it "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/core/roles"
	"github.com/moyai-network/teams/internal/core/user/class"
	"github.com/moyai-network/teams/internal/model"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bedrock-gophers/unsafe/unsafe"

	"github.com/google/uuid"

	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/moyai-network/teams/internal"

	"github.com/bedrock-gophers/cooldown/cooldown"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/pkg/lang"

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
	uuid uuid.UUID

	logTime           time.Time
	claimSelectionPos [2]mgl64.Vec3
	waypoint          *WayPoint
	energy            atomic.Value[float64]

	wallBlocks   map[cube.Pos]float64
	wallBlocksMu sync.Mutex

	lastArmour       atomic.Value[[4]item.Stack]
	lastClass        atomic.Value[class.Class]
	lastScoreBoard   atomic.Value[*scoreboard.Scoreboard]
	lastArea         atomic.Value[area.NamedArea]
	lastAttackerName atomic.Value[string]
	lastAttackTime   atomic.Value[time.Time]
	lastPearlPos     mgl64.Vec3
	lastMessage      atomic.Value[time.Time]

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
	processLogout             *Process
	processStuck              *Process
	processHome               *Process
	processCamp               *Process

	gracefulLogout bool
	logger         bool

	close chan struct{}
	death chan struct{}
}

func NewHandler(p *player.Player, xuid string) (*Handler, error) {
	sendFog(p)
	/*if h, ok := logger(p); ok {
		if p.World() == internal.Deathban() {
			<-time.After(time.Second)
			internal.Deathban().AddEntity(p)
			p.SetGameMode(world.GameModeSurvival)
		}
		p.Teleport(p.Position())
		currentHealth := p.Health()
		p.Hurt(20-currentHealth, NoArmourAttackEntitySource{})
		_ = p.Close()
	}*/

	h := &Handler{
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
		processHome: NewProcess(func(t *Process) {
			p.Message(text.Colourf("<green>You have been teleported home.</green>"))
		}),
		processStuck: NewProcess(func(t *Process) {
			p.Message(text.Colourf("<red>You have been teleported to a safe place.</red>"))
		}),

		close: make(chan struct{}),
		death: make(chan struct{}),
	}

	h.processLogout = NewProcess(func(t *Process) {
		h.gracefulLogout = true
		p.Disconnect(text.Colourf("<red>You have been logged out.</red>"))
	})

	p.SetNameTag(text.Colourf("<red>%s</red>", p.Name()))
	UpdateState(p)

	s := unsafe.Session(p)
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return nil, err
	}

	u.StaffMode = false
	if u.Teams.DeathBan.Active() {
		internal.Deathban().Exec(func(tx *world.Tx) {
			tx.AddEntity(p.H())
		})
		p.Teleport(mgl64.Vec3{5, 13, 44})
	} else {
		if u.Teams.DeathBanned {
			u.Teams.DeathBanned = false
			u.Teams.PVP.Set(time.Hour + (time.Millisecond * 500))
			if !u.Teams.PVP.Paused() {
				u.Teams.PVP.TogglePause()
			}

			internal.Deathban().Exec(func(tx *world.Tx) {
				tx.AddEntity(p.H())
			})
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

	p.Message(lang.Translatef(*u.Language, "discord.message"))
	h.handleBoosterRole(p, u)

	data.SaveUser(u)
	if u.Frozen {
		p.SetImmobile()
	}

	h.updateCurrentArea(p, p.Position(), u)
	h.updateKOTHState(p, p.Position(), u)
	UpdateVanishState(p, u)

	h.logTime = time.Now()
	UpdateState(p)
	go startTicker(p, h)
	return h, nil
}

func sendFog(p *player.Player) {
	var stack []string
	switch p.Tx().World().Dimension() {
	case world.End:
		stack = []string{"minecraft:fog_the_end"}
	case world.Nether:
		stack = []string{"minecraft:fog_hell"}
	}
	unsafe.WritePacket(p, &packet.PlayerFog{
		Stack: stack,
	})
}

func (h *Handler) handleBoosterRole(p *player.Player, u model.User) {
	if len(u.DiscordID) > 0 {
		userID, _ := strconv.Atoi(u.DiscordID)
		rl, err := internal.DiscordState().MemberRoles(discord.GuildID(1111055709300342826), discord.UserID(userID))
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

func vec3ToVec2(v mgl64.Vec3) mgl64.Vec2 {
	return mgl64.Vec2{v.X(), v.Z()}
}

func maxMin(n, n2 float64) (max float64, min float64) {
	if n > n2 {
		return n, n2
	}
	return n2, n
}

type npcHandler struct {
	player.NopHandler
}

func (npcHandler) HandleItemPickup(ctx *player.Context, _ *item.Stack) {
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
	panic("to fix")
	/*p.Inventory().Handle(inventory.NopHandler{})
	dat, err := internal.LoadPlayerData(p.UUID())
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
	}*/
}

package user

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/core"
	class2 "github.com/moyai-network/teams/internal/core/class"
	"github.com/moyai-network/teams/internal/core/colour"
	item2 "github.com/moyai-network/teams/internal/core/item"
	kit2 "github.com/moyai-network/teams/internal/core/kit"
	"github.com/moyai-network/teams/internal/core/koth"
	"github.com/moyai-network/teams/internal/core/menu"
	"github.com/moyai-network/teams/internal/core/sotw"
	"github.com/moyai-network/teams/internal/model"
	"math/rand"
	"slices"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/pkg/effectutil"
	"github.com/restartfu/roman"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// HandleItemUse handles the use of items by the player.
func (h *Handler) HandleItemUse(ctx *player.Context) {
	p := ctx.Val()

	held, left := p.HeldItems()
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}

	switch held.Item().(type) {
	case item.EnderPearl:
		h.handlePearlUse(ctx)
	}

	// Deposit money note into bank account.
	if v, ok := held.Value("MONEY_NOTE"); ok {
		u.Teams.Balance += v.(float64)
		p.SetHeldItems(subtractItem(p, held, 1), left)
		p.Message(text.Colourf("<green>You have deposited $%.0f into your bank account</green>", v.(float64)))
		core.UserRepository.Save(u)
		return
	}

	_, sotwRunning := sotw.Running()
	// Handle specific item types based on player's class.
	switch h.lastClass.Load().(type) {
	case class2.Archer, class2.Rogue:
		h.handleArcherRogueItemUse(ctx, held, left)
	case class2.Bard:
		h.handleBardItemUse(ctx, held, sotwRunning, u)
	case class2.Mage:
		h.handleMageItemUse(p, held, sotwRunning, u)
	}

	// Handle special items.
	if v, ok := item2.SpecialItem(held); ok {
		h.handleSpecialItemUse(ctx, v, held, left)
	} else if v, ok = item2.PartnerItem(held); ok {
		if h.lastArea.Load().Name == koth.Citadel.Name() {
			internal.Messagef(p, "item.use.citadel.disabled")
			ctx.Cancel()
			return
		}
		if h.coolDownGlobalAbilities.Active() {
			internal.Messagef(p, "partner_item.cooldown", h.coolDownGlobalAbilities.Remaining().Seconds())
			ctx.Cancel()
		} else {
			h.handleSpecialItemUse(ctx, v, held, left)
		}
	}
}

// handlePartnerPackage uses the partner package item.
func (h *Handler) handlePartnerPackage(ctx *player.Context, held item.Stack, left item.Stack, pos cube.Pos) {
	p := ctx.Val()

	keys := item2.PartnerItems()
	i := item2.NewSpecialItem(keys[rand.Intn(len(keys))], rand.Intn(3)+1)
	if ite, ok := item2.SpecialItem(i); ok {
		if _, ok2 := ite.(item2.SigilType); ok2 {
			// Hacky way to re-roll so that it's lower probability
			i = item2.NewSpecialItem(keys[rand.Intn(len(keys))], rand.Intn(3)+1)
		}
		if _, ok2 := ite.(item2.StormBreakerType); ok2 {
			// Hacky way to re-roll so that it's lower probability
			i = item2.NewSpecialItem(keys[rand.Intn(len(keys))], rand.Intn(3)+1)
		}
	}

	ctx.Cancel()
	if held.Count() == 1 {
		p.SetHeldItems(item.Stack{}, left)
	} else {
		p.SetHeldItems(held.Grow(-1), left)
	}
	item2.AddOrDrop(p, i)

	opts := world.EntitySpawnOpts{
		Position: pos.Vec3(),
		Rotation: cube.Rotation{90, 90},
	}
	p.Tx().AddEntity(entity.NewFirework(opts, item.Firework{
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

// handleArcherRogueItemUse handles the use of items by Archer or Rogue class.
func (h *Handler) handleArcherRogueItemUse(ctx *player.Context, held item.Stack, left item.Stack) {
	p := ctx.Val()

	// Get the corresponding effect for the item.
	if _, ok := held.Item().(item.Firework); ok {
		return
	}
	if e, ok := ArcherRogueEffectFromItem(held.Item()); ok {
		// Check cooldown for the item.
		if cd := h.coolDownArcherRogueItem.Key(held.Item()); cd.Active() {
			internal.Messagef(p, "class.ability.cooldown", cd.Remaining().Seconds())
			return
		}

		// Apply the effect.
		p.AddEffect(e)
		h.coolDownArcherRogueItem.Key(held.Item()).Set(60 * time.Second)
		p.SetHeldItems(held.Grow(-1), item.Stack{})
	}
}

// handleBardItemUse handles the use of items by Bard class.
func (h *Handler) handleBardItemUse(ctx *player.Context, held item.Stack, sotwRunning bool, u model.User) {
	p := ctx.Val()

	// Ignore if the item is a chest.
	if _, ok := held.Item().(block.Chest); ok {
		return
	}

	if _, ok := held.Item().(item.Firework); ok {
		return
	}

	// Get the corresponding effect for the item.
	if e, ok := BardEffectFromItem(held.Item()); ok {
		// Check PvP and SOTW status.
		if u.Teams.PVP.Active() || sotwRunning && u.Teams.SOTW {
			return
		}

		// Check cooldown for the item.
		if cd := h.coolDownBardItem.Key(held.Item()); cd.Active() {
			internal.Messagef(p, "class.ability.cooldown", cd.Remaining().Seconds())
			return
		}

		// Check energy level.
		if en := h.energy.Load(); en < 30 {
			internal.Messagef(p, "class.energy.insufficient")
			return
		} else {
			h.energy.Store(en - 30)
		}

		// Apply the effect to nearby allies.
		teammates := nearbyAllies(p, 30)
		for _, m := range teammates {
			m.AddEffect(e)
		}

		// Notify player about ability usage.
		lvl, _ := roman.Itor(e.Level())
		internal.Messagef(p, "bard.ability.use", effectutil.EffectName(e), lvl, len(teammates))
		p.SetHeldItems(held.Grow(-1), item.Stack{})
		h.coolDownBardItem.Key(held.Item()).Set(15 * time.Second)
	}
}

// handleMageItemUse handles the use of items by Mage class.
func (h *Handler) handleMageItemUse(p *player.Player, held item.Stack, sotwRunning bool, u model.User) {
	// Get the corresponding effect for the item.
	if e, ok := MageEffectFromItem(held.Item()); ok {
		// Check PvP and SOTW status.
		if u.Teams.PVP.Active() || sotwRunning && u.Teams.SOTW {
			return
		}

		// Check cooldown for the item.
		if cd := h.coolDownMageItem.Key(held.Item()); cd.Active() {
			internal.Messagef(p, "class.ability.cooldown", cd.Remaining().Seconds())
			return
		}

		// Check energy level.
		if en := h.energy.Load(); en < 30 {
			internal.Messagef(p, "class.energy.insufficient")
			return
		} else {
			h.energy.Store(en - 30)
		}

		// Apply the effect to nearby enemies.
		enemies := nearbyEnemies(p, 25)
		for _, m := range enemies {
			m.AddEffect(e)
		}

		// Notify player about ability usage.
		lvl, _ := roman.Itor(e.Level())
		internal.Messagef(p, "mage.ability.use", effectutil.EffectName(e), lvl, len(enemies))
		p.SetHeldItems(held.Grow(-1), item.Stack{})
		h.coolDownMageItem.Key(held.Item()).Set(15 * time.Second)
	}
}

// handleSpecialItemUse handles the use of special items.
func (h *Handler) handleSpecialItemUse(ctx *player.Context, v item2.SpecialItemType, held item.Stack, left item.Stack) {
	p := ctx.Val()

	// Handle specific abilities for special items.
	switch kind := v.(type) {
	case item2.TimeWarpType:
		h.handleTimeWarpAbility(p, kind)
	case item2.SwitcherBallType:
		h.handleSwitcherBallAbility(ctx, kind)
	case item2.FullInvisibilityType:
		h.handleFullInvisibilityAbility(p, kind, held, left)
	case item2.CloseCallType:
		ctx.Cancel()
		h.handleCloseCallAbility(p, kind, held, left)
	case item2.BeserkAbilityType:
		h.handleBeserkAbility(p, kind, held, left)
	case item2.NinjaStarType:
		h.handleNinjaStarAbility(p, kind)
	case item2.FocusModeType:
		h.handleFocusModeAbility(p, kind, held, left)
	case item2.RocketType:
		h.handleRocketAbility(p, kind)
	case item2.VampireAbilityType:
		h.handleVampireAbility(p, kind, held, left)
	case item2.AbilityDisablerType:
		h.handleAbilityDisablerAbility(p, kind, held, left)
	case item2.StrengthPowderType:
		h.handleStrengthPowderAbility(p, kind, held, left)
	case item2.TankIngotType:
		h.handleTankIngotAbility(p, kind, held, left)
	case item2.RageBrickType:
		h.handleRageBrickAbility(p, kind, held, left)
	case item2.ComboAbilityType:
		h.handleComboAbility(p, kind, held, left)
	case item2.PartnerPackageType:
		ctx.Cancel()
		h.handlePartnerPackage(ctx, held, left, cube.PosFromVec3(p.Position()))
	case item2.StaffRandomTeleportType:
		h.handleRandomTeleport(ctx)
	case item2.StaffUnVanishType, item2.StaffVanishType:
		p.ExecuteCommand("/vanish")
		p.Inventory().Clear()
		p.Armour().Clear()
		// TODO: potentially save the player's inventory and armour to restore it later.
		kit2.Apply(kit2.Staff{}, p)
	case item2.StaffTeleportStickType:
		menu.SendPlayerListTeleportMenu(p)
	}
}

// handleRandomTeleport handles the random teleport staff item.
func (h *Handler) handleRandomTeleport(ctx *player.Context) {
	p := ctx.Val()

	players := slices.Collect(internal.Players(p.Tx()))
	if len(players) == 1 {
		return
	}
	target := players[rand.Intn(len(players))]
	if target == p {
		h.handleRandomTeleport(ctx)
		return
	}
	target.Tx().AddEntity(p.H())
	p.Teleport(target.Position())
	internal.Messagef(p, "command.teleport.self", target.Name())
}

// handlePearlUse handles the use of ender pearls.`
func (h *Handler) handlePearlUse(ctx *player.Context) {
	p := ctx.Val()

	if cd := h.coolDownPearl; cd.Active() {
		internal.Messagef(p, "pearl.cooldown", cd.Remaining().Seconds())
		ctx.Cancel()
		return
	}

	if h.lastArea.Load().Name == koth.Citadel.Name() {
		internal.Messagef(p, "item.use.citadel.disabled")
		ctx.Cancel()
		return
	}

	h.coolDownPearl.Set(time.Second * 15)
	h.lastPearlPos = p.Position()

	//go h.handlePearlExperienceBar()
}

/*func (h *Handler) handlePearlExperienceBar() {
	t := time.NewTicker(time.Millisecond * 50)
	for range t.C {
		if !h.coolDownPearl.Active() {
			p.RemoveExperience(math.MaxInt)
			t.Stop()
			return
		}
		p.SetExperienceLevel(int(h.coolDownPearl.Remaining().Seconds()))
		p := float64(h.coolDownPearl.Remaining()) / float64(time.Second*15)
		if p > 1 {
			p = 1
		}

		if p < 0 {
			p = 0
		}

		p.SetExperienceProgress(p)
	}
}*/

// handleTimeWarpAbility handles the Time Warp ability.
func (h *Handler) handleTimeWarpAbility(p *player.Player, kind item2.TimeWarpType) {
	if h.lastPearlPos == (mgl64.Vec3{}) || !h.coolDownPearl.Active() {
		p.Message(text.Colourf("<red>You do not have a last thrown pearl or it has expired.</red>"))
		return
	}
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Time Warp", cd.Remaining().Seconds())
		return
	}
	internal.Messagef(p, "partner_item.used", "Time Warp")
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*1)
	time.AfterFunc(time.Second*2, func() {
		if p != nil {
			p.Teleport(h.lastPearlPos)
			h.lastPearlPos = mgl64.Vec3{}
		}
	})
}

// handleSwitcherBallAbility handles the Switcher Ball ability.
func (h *Handler) handleSwitcherBallAbility(ctx *player.Context, kind item2.SwitcherBallType) {
	p := ctx.Val()
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Switcher Ball", cd.Remaining().Seconds())
		ctx.Cancel()
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Key(kind).Set(time.Second * 10)
}

// handleFullInvisibilityAbility handles the Full Invisibility ability.
func (h *Handler) handleFullInvisibilityAbility(p *player.Player, kind item2.FullInvisibilityType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Full Invisibility", cd.Remaining().Seconds())
		return
	}
	h.ShowArmor(p, false)
	p.AddEffect(effect.New(effect.Invisibility, 1, time.Second*60))
	p.SetHeldItems(subtractItem(p, held, 1), left)
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)
	internal.Messagef(p, "partner_item.used", "Full Invisibility")
}

// handleCloseCallAbility handles the Close Call ability.
func (h *Handler) handleCloseCallAbility(p *player.Player, kind item2.CloseCallType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Close Call", cd.Remaining().Seconds())
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*3)

	p.SetHeldItems(subtractItem(p, held, 1), left)

	if p.Health() <= 8 {
		p.AddEffect(effect.New(effect.Regeneration, 6, time.Second*5))
		p.AddEffect(effect.New(effect.Strength, 2, time.Second*5))
		p.Message()
		internal.Messagef(p, "partner_item.used", "Close Call")
	}
}

// handleBeserkAbility handles the Beserk Ability.
func (h *Handler) handleBeserkAbility(p *player.Player, kind item2.BeserkAbilityType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Beserk Ability", cd.Remaining().Seconds())
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*5)

	p.SetHeldItems(subtractItem(p, held, 1), left)

	p.AddEffect(effect.New(effect.Strength, 2, time.Second*12))
	p.AddEffect(effect.New(effect.Resistance, 3, time.Second*12))
	p.AddEffect(effect.New(effect.Regeneration, 3, time.Second*12))
	internal.Messagef(p, "partner_item.used", "Beserk Ability")
}

// handleNinjaStarAbility handles the Ninja Star ability.
func (h *Handler) handleNinjaStarAbility(p *player.Player, kind item2.NinjaStarType) {
	lastAttacker, ok := h.lastAttacker(p.Tx())
	if !ok {
		internal.Messagef(p, "partner_item.last_hit")
		return
	}
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Ninja Star", cd.Remaining().Seconds())
		return
	}

	internal.Messagef(p, "partner_item.used.teleporting", "Ninja Star", lastAttacker.Name())
	internal.Messagef(lastAttacker, "ninja_star.teleporting", p.Name())
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	time.AfterFunc(time.Second*5, func() {
		lastAttacker, ok = Lookup(p.Tx(), h.lastAttackerName.Load())
		if p != nil && ok {
			p.Teleport(lastAttacker.Position())
		}
	})
}

// handleFocusModeAbility handles the Focus Mode ability.
func (h *Handler) handleFocusModeAbility(p *player.Player, kind item2.FocusModeType, held item.Stack, left item.Stack) {
	lastAttacker, ok := h.lastAttacker(p.Tx())
	if !ok {
		p.Message(text.Colourf("<red>No last valid hit found</red>"))
		return
	}
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Focus Mode", cd.Remaining().Seconds())
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	if t, ok := lastAttacker.Handler().(*Handler); ok {
		t.coolDownFocusMode.Set(time.Second * 10)
		p.SetHeldItems(subtractItem(p, held, 1), left)
		//internal.Messagef(p, "partner_item.used.on", "Focus Mode", t.Name())
		//internal.Messagef(t.p, "focus_mode.used", p)
	}
}

// handleFocusModeAbility handles the Rocket ability
func (h *Handler) handleRocketAbility(p *player.Player, kind item2.RocketType) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Rocket", cd.Remaining().Seconds())
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	p.SetVelocity(mgl64.Vec3{0, 2.5})
	p.Tx().PlaySound(p.Position(), sound.FireworkLaunch{})
}

// handleAbilityDisablerAbility handles the ability disabler ability
func (h *Handler) handleAbilityDisablerAbility(p *player.Player, kind item2.AbilityDisablerType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Ability Disabler", cd.Remaining().Seconds())
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	p.SetHeldItems(subtractItem(p, held, 1), left)

	enemies := nearbyHurtable(p, 15)
	for _, e := range enemies {
		//e.coolDownGlobalAbilities.Set(time.Second * 10)
		internal.Messagef(e, "ability_disabler.used", p.Name())
	}

	internal.Messagef(p, "partner_item.used", "Ability Disabler")
}

// handleAbilityDisablerAbility handles the ability disabler ability
func (h *Handler) handleStrengthPowderAbility(p *player.Player, kind item2.StrengthPowderType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Strength Powder", cd.Remaining().Seconds())
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	p.SetHeldItems(subtractItem(p, held, 1), left)

	p.AddEffect(effect.New(effect.Strength, 2, time.Second*7))
	internal.Messagef(p, "partner_item.used", "Strength Powder")
}

func (h *Handler) handleTankIngotAbility(p *player.Player, kind item2.TankIngotType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Tank Ingot", cd.Remaining().Seconds())
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	p.SetHeldItems(subtractItem(p, held, 1), left)

	p.AddEffect(effect.New(effect.Resistance, 3, time.Second*7))
	internal.Messagef(p, "partner_item.used", "Tank Ingot")
}

func (h *Handler) handleRageBrickAbility(p *player.Player, kind item2.RageBrickType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Rage Brick", cd.Remaining().Seconds())
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	p.SetHeldItems(subtractItem(p, held, 1), left)

	enemies := nearbyHurtable(p, 15)
	lasts := time.Second * time.Duration(len(enemies))
	p.AddEffect(effect.New(effect.Strength, 2, lasts))
	internal.Messagef(p, "partner_item.used", "Rage Brick")
}

func (h *Handler) handleComboAbility(p *player.Player, kind item2.ComboAbilityType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Combo Ability", cd.Remaining().Seconds())
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)
	h.coolDownComboAbility.Set(time.Second * 10)

	p.SetHeldItems(subtractItem(p, held, 1), left)

	internal.Messagef(p, "partner_item.used", "Combo Ability")
}

func (h *Handler) handleVampireAbility(p *player.Player, kind item2.VampireAbilityType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		internal.Messagef(p, "partner_item.item.cooldown", "Vampire Ability", cd.Remaining().Seconds())
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)
	h.coolDownVampireAbility.Set(time.Second * 10)

	p.SetHeldItems(subtractItem(p, held, 1), left)

	p.AddEffect(effect.New(effect.Haste, 2, time.Second*7))
	internal.Messagef(p, "partner_item.used", "Vampire Ability")
}

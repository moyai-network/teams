package user

import (
	"math/rand"
	"time"

	"github.com/moyai-network/teams/moyai/koth"

	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/moyai/kit"
	"github.com/moyai-network/teams/moyai/menu"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal/effectutil"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/class"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/restartfu/roman"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// HandleItemUse handles the use of items by the player.
func (h *Handler) HandleItemUse(ctx *event.Context) {
	held, left := h.p.HeldItems()
	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil {
		return
	}

	switch held.Item().(type) {
	case item.EnderPearl:
		h.handlePearlUse(ctx)
	}

	// Deposit money note into bank account.
	if v, ok := held.Value("MONEY_NOTE"); ok {
		u.Teams.Balance += v.(float64)
		h.p.SetHeldItems(substractItem(h.p, held, 1), left)
		h.p.Message(text.Colourf("<green>You have deposited $%.0f into your bank account</green>", v.(float64)))
		data.SaveUser(u)
		return
	}

	_, sotwRunning := sotw.Running()
	// Handle specific item types based on player's class.
	switch h.lastClass.Load().(type) {
	case class.Archer, class.Rogue:
		h.handleArcherRogueItemUse(held, left)
	case class.Bard:
		h.handleBardItemUse(held, sotwRunning, u)
	case class.Mage:
		h.handleMageItemUse(held, sotwRunning, u)
	}

	// Handle special items.
	if v, ok := it.SpecialItem(held); ok {
		h.handleSpecialItemUse(v, held, left, ctx)
	} else if v, ok = it.PartnerItem(held); ok {
		if h.lastArea.Load().Name() == koth.Citadel.Name() {
			moyai.Messagef(h.p, "item.use.citadel.disabled")
			ctx.Cancel()
			return
		}
		h.handleSpecialItemUse(v, held, left, ctx)
	}
}

// handlePartnerPackage uses the partner package item.
func (h *Handler) handlePartnerPackage(ctx *event.Context, held item.Stack, left item.Stack, pos cube.Pos) {
	w := h.p.World()

	keys := it.PartnerItems()
	i := it.NewSpecialItem(keys[rand.Intn(len(keys))], rand.Intn(3)+1)
	if ite, ok := it.SpecialItem(i); ok {
		if _, ok2 := ite.(it.SigilType); ok2 {
			// Hacky way to re-roll so that it's lower probability
			i = it.NewSpecialItem(keys[rand.Intn(len(keys))], rand.Intn(3)+1)
		}
		if _, ok2 := ite.(it.StormBreakerType); ok2 {
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
}

// handleArcherRogueItemUse handles the use of items by Archer or Rogue class.
func (h *Handler) handleArcherRogueItemUse(held item.Stack, left item.Stack) {
	// Get the corresponding effect for the item.
	if e, ok := ArcherRogueEffectFromItem(held.Item()); ok {
		// Check cooldown for the item.
		if cd := h.coolDownArcherRogueItem.Key(held.Item()); cd.Active() {
			moyai.Messagef(h.p, "class.ability.cooldown", cd.Remaining().Seconds())
			return
		}

		// Apply the effect.
		h.p.AddEffect(e)
		h.coolDownArcherRogueItem.Key(held.Item()).Set(60 * time.Second)
		h.p.SetHeldItems(held.Grow(-1), item.Stack{})
	}
}

// handleBardItemUse handles the use of items by Bard class.
func (h *Handler) handleBardItemUse(held item.Stack, sotwRunning bool, u data.User) {
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
			moyai.Messagef(h.p, "class.ability.cooldown", cd.Remaining().Seconds())
			return
		}

		// Check energy level.
		if en := h.energy.Load(); en < 30 {
			moyai.Messagef(h.p, "class.energy.insufficient")
			return
		} else {
			h.energy.Store(en - 30)
		}

		// Apply the effect to nearby allies.
		teammates := nearbyAllies(h.p, 25)
		for _, m := range teammates {
			m.p.AddEffect(e)
		}

		// Notify player about ability usage.
		lvl, _ := roman.Itor(e.Level())
		moyai.Messagef(h.p, "bard.ability.use", effectutil.EffectName(e), lvl, len(teammates))
		h.p.SetHeldItems(held.Grow(-1), item.Stack{})
		h.coolDownBardItem.Key(held.Item()).Set(15 * time.Second)
	}
}

// handleMageItemUse handles the use of items by Mage class.
func (h *Handler) handleMageItemUse(held item.Stack, sotwRunning bool, u data.User) {
	// Get the corresponding effect for the item.
	if e, ok := MageEffectFromItem(held.Item()); ok {
		// Check PvP and SOTW status.
		if u.Teams.PVP.Active() || sotwRunning && u.Teams.SOTW {
			return
		}

		// Check cooldown for the item.
		if cd := h.coolDownMageItem.Key(held.Item()); cd.Active() {
			moyai.Messagef(h.p, "class.ability.cooldown", cd.Remaining().Seconds())
			return
		}

		// Check energy level.
		if en := h.energy.Load(); en < 30 {
			moyai.Messagef(h.p, "class.energy.insufficient")
			return
		} else {
			h.energy.Store(en - 30)
		}

		// Apply the effect to nearby enemies.
		enemies := nearbyEnemies(h.p, 25)
		for _, m := range enemies {
			m.p.AddEffect(e)
		}

		// Notify player about ability usage.
		lvl, _ := roman.Itor(e.Level())
		moyai.Messagef(h.p, "mage.ability.use", effectutil.EffectName(e), lvl, len(enemies))
		h.p.SetHeldItems(held.Grow(-1), item.Stack{})
		h.coolDownMageItem.Key(held.Item()).Set(15 * time.Second)
	}
}

// handleSpecialItemUse handles the use of special items.
func (h *Handler) handleSpecialItemUse(v it.SpecialItemType, held item.Stack, left item.Stack, ctx *event.Context) {
	// Handle specific abilities for special items.
	switch kind := v.(type) {
	case it.TimeWarpType:
		h.handleTimeWarpAbility(kind)
	case it.SwitcherBallType:
		h.handleSwitcherBallAbility(kind, ctx)
	case it.FullInvisibilityType:
		h.handleFullInvisibilityAbility(kind, held, left)
	case it.CloseCallType:
		ctx.Cancel()
		h.handleCloseCallAbility(kind, held, left)
	case it.BeserkAbilityType:
		h.handleBeserkAbility(kind, held, left)
	case it.NinjaStarType:
		h.handleNinjaStarAbility(kind)
	case it.FocusModeType:
		h.handleFocusModeAbility(kind, held, left)
	case it.RocketType:
		h.handleRocketAbility(kind)
	case it.VampireAbilityType:
		h.handleVampireAbility(kind, held, left)
	case it.AbilityDisablerType:
		h.handleAbilityDisablerAbility(kind, held, left)
	case it.StrengthPowderType:
		h.handleStrengthPowderAbility(kind, held, left)
	case it.TankIngotType:
		h.handleTankIngotAbility(kind, held, left)
	case it.RageBrickType:
		h.handleRageBrickAbility(kind, held, left)
	case it.ComboAbilityType:
		h.handleComboAbility(kind, held, left)
	case it.PartnerPackageType:
		ctx.Cancel()
		h.handlePartnerPackage(ctx, held, left, cube.PosFromVec3(h.p.Position()))
	case it.StaffRandomTeleportType:
		h.handleRandomTeleport()
	case it.StaffUnVanishType, it.StaffVanishType:
		h.p.ExecuteCommand("/vanish")
		h.p.Inventory().Clear()
		h.p.Armour().Clear()
		// TODO: potentially save the player's inventory and armour to restore it later.
		kit.Apply(kit.Staff{}, h.p)
	case it.StaffTeleportStickType:
		menu.SendPlayerListTeleportMenu(h.p)
	}
}

// handleRandomTeleport handles the random teleport staff item.
func (h *Handler) handleRandomTeleport() {
	players := moyai.Players()
	if len(players) == 1 {
		return
	}
	target := players[rand.Intn(len(players))]
	if target == h.p {
		h.handleRandomTeleport()
		return
	}
	target.World().AddEntity(h.p)
	h.p.Teleport(target.Position())
	moyai.Messagef(h.p, "command.teleport.self", target.Name())
}

// handlePearlUse handles the use of ender pearls.`
func (h *Handler) handlePearlUse(ctx *event.Context) {
	if cd := h.coolDownPearl; cd.Active() {
		h.p.Message(text.Colourf("<red>You are on ender pearl cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		ctx.Cancel()
		return
	}

	if h.lastArea.Load().Name() == koth.Citadel.Name() {
		moyai.Messagef(h.p, "item.use.citadel.disabled")
		ctx.Cancel()
		return
	}

	h.coolDownPearl.Set(time.Second * 15)
	h.lastPearlPos = h.p.Position()
}

// handleTimeWarpAbility handles the Time Warp ability.
func (h *Handler) handleTimeWarpAbility(kind it.TimeWarpType) {
	if h.lastPearlPos == (mgl64.Vec3{}) || !h.coolDownPearl.Active() {
		h.p.Message(text.Colourf("<red>You do not have a last thrown pearl or it has expired.</red>"))
		return
	}
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Time Warp cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}
	h.p.Message(text.Colourf("<green>Ongoing in 2 seconds...</green>"))
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*1)
	time.AfterFunc(time.Second*2, func() {
		if h.p != nil {
			h.p.Teleport(h.lastPearlPos)
			h.lastPearlPos = mgl64.Vec3{}
		}
	})
}

// handleSwitcherBallAbility handles the Switcher Ball ability.
func (h *Handler) handleSwitcherBallAbility(kind it.SwitcherBallType, ctx *event.Context) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on snowball cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		ctx.Cancel()
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Key(kind).Set(time.Second * 10)
}

// handleFullInvisibilityAbility handles the Full Invisibility ability.
func (h *Handler) handleFullInvisibilityAbility(kind it.FullInvisibilityType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Full Invisibility cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}
	h.ShowArmor(false)
	h.p.AddEffect(effect.New(effect.Invisibility{}, 1, time.Hour).WithoutParticles())
	h.p.SetHeldItems(substractItem(h.p, held, 1), left)
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)
	h.p.Message(text.Colourf("§r§7> §eFull Invisibility §6has been used"))
}

// handleCloseCallAbility handles the Close Call ability.
func (h *Handler) handleCloseCallAbility(kind it.CloseCallType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Close Call cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*3)

	h.p.SetHeldItems(substractItem(h.p, held, 1), left)

	if h.p.Health() <= 6 {
		h.p.AddEffect(effect.New(effect.Regeneration{}, 6, time.Second*5))
		h.p.AddEffect(effect.New(effect.Strength{}, 2, time.Second*5))
		h.p.Message(text.Colourf("<green>Close Call has been used.</green>"))
	} else {
		h.p.Message(text.Colourf("<red>You are not below 3 hearts; Close call has been wasted.</red>"))
	}
}

// handleBeserkAbility handles the Beserk Ability.
func (h *Handler) handleBeserkAbility(kind it.BeserkAbilityType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Beserk Ability cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*5)

	h.p.SetHeldItems(substractItem(h.p, held, 1), left)

	h.p.AddEffect(effect.New(effect.Strength{}, 2, time.Second*12))
	h.p.AddEffect(effect.New(effect.Resistance{}, 3, time.Second*12))
	h.p.AddEffect(effect.New(effect.Regeneration{}, 3, time.Second*12))
	h.p.Message(text.Colourf("<green>Beserk Ability has been used.</green>"))
}

// handleNinjaStarAbility handles the Ninja Star ability.
func (h *Handler) handleNinjaStarAbility(kind it.NinjaStarType) {
	lastAttacker, ok := h.lastAttacker()
	if !ok {
		h.p.Message(text.Colourf("<red>No last valid hit found</red>"))
		return
	}
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Ninja Star cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}

	h.p.Message(text.Colourf("<red>Ongoing to %s in 5 seconds...</red>", lastAttacker.Name()))
	lastAttacker.Message(text.Colourf("<red>%s is teleporting to your in 5 seconds...</red>", h.p.Name()))
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	time.AfterFunc(time.Second*5, func() {
		lastAttacker, ok = Lookup(h.lastAttackerName.Load())
		if h.p != nil && ok {
			h.p.Teleport(lastAttacker.Position())
		}
	})
}

// handleFocusModeAbility handles the Focus Mode ability.
func (h *Handler) handleFocusModeAbility(kind it.FocusModeType, held item.Stack, left item.Stack) {
	lastAttacker, ok := h.lastAttacker()
	if !ok {
		h.p.Message(text.Colourf("<red>No last valid hit found</red>"))
		return
	}
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Focus Mode cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	if t, ok := lastAttacker.Handler().(*Handler); ok {
		t.coolDownFocusMode.Set(time.Second * 10)
		h.p.SetHeldItems(substractItem(h.p, held, 1), left)
		h.p.Message(text.Colourf("<green>Focus Mode has been used on %s.</green>", t.p.Name()))
		t.p.Message(text.Colourf("<green>%s has used Focus Mode on you.</green>", h.p.Name()))
	}
}

// handleFocusModeAbility handles the Rocket ability
func (h *Handler) handleRocketAbility(kind it.RocketType) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Rocket cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	h.p.SetVelocity(mgl64.Vec3{0, 2.5})
	h.p.World().PlaySound(h.p.Position(), sound.FireworkLaunch{})
}

// handleAbilityDisablerAbility handles the ability disabler ability
func (h *Handler) handleAbilityDisablerAbility(kind it.AbilityDisablerType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Ability Disabler cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	h.p.SetHeldItems(substractItem(h.p, held, 1), left)

	enemies := nearbyHurtable(h.p, 15)
	for _, e := range enemies {
		e.coolDownGlobalAbilities.Set(time.Second * 10)
		e.p.Message(text.Colourf("<red>You have been ability disabled for 10 seconds by %s.</red>", h.p.Name()))
	}
}

// handleAbilityDisablerAbility handles the ability disabler ability
func (h *Handler) handleStrengthPowderAbility(kind it.StrengthPowderType,  held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Strength Powder cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	h.p.SetHeldItems(substractItem(h.p, held, 1), left)

	h.p.AddEffect(effect.New(effect.Strength{}, 2, time.Second*7))
	h.p.Message(text.Colourf("§r§7> §eStrength Powder §6has been used"))
}

func (h *Handler) handleTankIngotAbility(kind it.TankIngotType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Tank Ingot cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	h.p.SetHeldItems(substractItem(h.p, held, 1), left)

	h.p.AddEffect(effect.New(effect.Resistance{}, 2, time.Second*7))
	h.p.Message(text.Colourf("§r§7> §eTank Ingot §6has been used"))
}

func (h *Handler) handleRageBrickAbility(kind it.RageBrickType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Rage Brick cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)

	h.p.SetHeldItems(substractItem(h.p, held, 1), left)

	enemies := nearbyHurtable(h.p, 15)
	lasts := time.Second * time.Duration(len(enemies))
	h.p.AddEffect(effect.New(effect.Strength{}, 2, lasts))
	h.p.Message(text.Colourf("§r§7> §eRage Brick §6has been used"))
}

func (h *Handler) handleComboAbility(kind it.ComboAbilityType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Combo Ability cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)
	h.coolDownComboAbility.Set(time.Second * 10)

	h.p.SetHeldItems(substractItem(h.p, held, 1), left)

	h.p.Message(text.Colourf("§r§7> §eCombo Ability §6has been used"))
}

func (h *Handler) handleVampireAbility(kind it.VampireAbilityType, held item.Stack, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(kind); cd.Active() {
		h.p.Message(text.Colourf("<red>You are on Vampire Ability cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
		return
	}
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(kind, time.Minute*2)
	h.coolDownVampireAbility.Set(time.Second * 10)

	h.p.SetHeldItems(substractItem(h.p, held, 1), left)

	h.p.AddEffect(effect.New(effect.Haste{}, 2, time.Second*7))
	h.p.Message(text.Colourf("§r§7> §eVampire Ability §6has been used"))
}
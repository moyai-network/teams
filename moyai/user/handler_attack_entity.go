package user

import (
	"time"

	"github.com/bedrock-gophers/knockback/knockback"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/class"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/moyai-network/teams/moyai/data"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/menu"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// HandleAttackEntity handles the attack on an entity.
func (h *Handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, _ *bool) {
	knockback.ApplyForce(force)
	knockback.ApplyHeight(height)

	targetPlayer, ok := e.(*player.Player)
	if !ok || h.handleSpecialCases(targetPlayer, ctx) {
		return
	}

	h.ShowArmor(true)

	if !canAttack(h.p, targetPlayer) || targetPlayer.AttackImmune() {
		ctx.Cancel()
		return
	}

	h.adjustKnockbackHeight(targetPlayer, height)
	h.applyAttackEnchantments(targetPlayer)
	h.handleRogueBackstab(ctx, targetPlayer, *force, *height)
	h.handleSpecialItemAbility(ctx, targetPlayer, force, height)
	h.setLastAttackerForCombat(targetPlayer)
	h.applyComboAbility()
}

// handleRogueBackstab handles the backstab ability for the Rogue class.
func (h *Handler) handleRogueBackstab(ctx *event.Context, t *player.Player, force, height float64) {
	held, left := h.p.HeldItems()
	if s, ok := held.Item().(item.Sword); ok && s.Tier == item.ToolTierGold && class.Compare(h.lastClass.Load(), class.Rogue{}) && t.Rotation().Direction() == h.p.Rotation().Direction() {
		cd := h.coolDownBackStab
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
			t.KnockBack(h.p.Position(), force, height)

			h.p.AddEffect(effect.New(effect.Slowness{}, 3, time.Second*5))
			h.p.SetHeldItems(item.Stack{}, left)
			cd.Set(time.Second * 10)
		}
	}
}

// applyAttackEnchantments applies attack enchantments from the attacker's armor.
func (h *Handler) applyAttackEnchantments(targetPlayer *player.Player) {
	arm := h.p.Armour()
	for _, a := range arm.Slots() {
		for _, e := range a.Enchantments() {
			if att, ok := e.Type().(ench.AttackEnchantment); ok {
				att.AttackEntity(h.p, targetPlayer)
			}
		}
	}
}

// handleSpecialItemAbility handles special abilities if the attacker is holding a special item.
func (h *Handler) handleSpecialItemAbility(ctx *event.Context, targetPlayer *player.Player, force, height *float64) {
	heldItem, left := h.p.HeldItems()
	if specialItemType, ok := it.PartnerItem(heldItem); ok {
		switch specialItemType.(type) {
		case it.PearlDisablerType:
			h.handlePearlDisablerAbility(targetPlayer, heldItem, left)
		case it.ExoticBoneType:
			h.handleExoticBoneAbility(targetPlayer, heldItem, left)
		case it.StormBreakerType:
			h.handleStormBreakerAbility(targetPlayer, left)
		case it.EffectDisablerType:
			h.handleEffectDisablerAbility(targetPlayer, left)
		}
	} else if specialItemType, ok = it.SpecialItem(heldItem); ok {
		switch specialItemType.(type) {
		case it.StaffKnockBackStickType:
			*force = *force * 1.5
		case it.StaffFreezeBlockType:
			h.p.ExecuteCommand("/freeze \"" + targetPlayer.Name() + "\"")
			ctx.Cancel()
		}
	}
}

// adjustKnockbackHeight adjusts the knockback height based on the target's airborne status.
func (h *Handler) adjustKnockbackHeight(targetPlayer *player.Player, height *float64) {
	if !targetPlayer.OnGround() {
		mx, mn := maxMin(targetPlayer.Position().Y(), h.p.Position().Y())
		if mx-mn >= 2.5 {
			*height = 0.38 / 1.25
		}
	}
}

// setLastAttackerForCombat sets the last attacker for combat tagging.
func (h *Handler) setLastAttackerForCombat(targetPlayer *player.Player) {
	targetHandler, ok := targetPlayer.Handler().(*Handler)
	if !ok {
		return
	}

	if canAttack(h.p, targetPlayer) {
		targetHandler.setLastAttacker(h)
		h.tagCombat.Set(time.Second * 30)
		targetHandler.tagCombat.Set(time.Second * 30)
	}
}

// canAttack returns true if the given players can attack each other.
func canAttack(pl, target *player.Player) bool {
	if target == nil || pl == nil || target.GameMode() == world.GameModeCreative {
		return false
	}
	w := pl.World()
	if ((area.Spawn(w).Vec3WithinOrEqualFloorXZ(pl.Position())) || area.Spawn(w).Vec3WithinOrEqualFloorXZ(target.Position())) && w != moyai.Deathban() {
		return false
	}

	u, err := data.LoadUserFromName(pl.Name())
	t, err2 := data.LoadUserFromName(target.Name())

	if err != nil || err2 != nil {
		return true
	}

	_, sotwRunning := sotw.Running()
	if u.Teams.PVP.Active() && !u.Teams.DeathBan.Active() || t.Teams.PVP.Active() && !t.Teams.DeathBan.Active() {
		return false
	}

	if sotwRunning && (u.Teams.SOTW && !u.Teams.DeathBan.Active() || t.Teams.SOTW && !t.Teams.DeathBan.Active()) {
		return false
	}

	if area.Deathban.Spawn().Vec3WithinOrEqualFloorXZ(pl.Position()) || area.Deathban.Spawn().Vec3WithinOrEqualFloorXZ(target.Position()) {
		return false
	}

	tm, err := data.LoadTeamFromMemberName(pl.Name())
	if err != nil {
		return true
	}

	if tm.Member(target.Name()) {
		return false
	}

	return true
}

// handlePearlDisablerAbility handles the Pearl Disabler special ability.
func (h *Handler) handlePearlDisablerAbility(targetPlayer *player.Player, held, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(it.PearlDisablerType{}); cd.Active() {
		moyai.Messagef(h.p, "pearl_disabler.cooldown", cd.Remaining().Seconds())
		return
	}

	moyai.Messagef(targetPlayer, "pearl_disabler.target", h.p.Name())
	targetHandler, ok := targetPlayer.Handler().(*Handler)
	if !ok {
		return
	}
	targetHandler.coolDownPearl.Set(time.Second * 15)

	moyai.Messagef(h.p, "pearl_disabler.user", targetPlayer.Name())
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(it.PearlDisablerType{}, time.Minute)
	h.p.SetHeldItems(subtractItem(h.p, held, 1), left)
}

// handleExoticBoneAbility handles the Exotic Bone special ability.
func (h *Handler) handleExoticBoneAbility(targetPlayer *player.Player, held, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(it.ExoticBoneType{}); cd.Active() {
		moyai.Messagef(h.p, "bone.cooldown", cd.Remaining().Seconds())
		return
	}

	moyai.Messagef(targetPlayer, "bone.target", h.p.Name())
	moyai.Messagef(h.p, "bone.user", targetPlayer.Name())

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(it.ExoticBoneType{}, time.Minute)
	h.p.SetHeldItems(subtractItem(h.p, held, 1), left)

	targetHandler, ok := targetPlayer.Handler().(*Handler)
	if !ok {
		return
	}
	targetHandler.coolDownBonedEffect.Set(time.Second * 10)
}

// handleStormBreakerAbility handles the Storm Breaker special ability.
func (h *Handler) handleStormBreakerAbility(targetPlayer *player.Player, left item.Stack) {
	targetHandler, ok := targetPlayer.Handler().(*Handler)
	if !ok || targetHandler.lastClass.Load() != nil || h.checkCoolDownSpecific(it.StormBreakerType{}) {
		return
	}

	h.triggerStormBreakerEffect(targetPlayer, left)
}

// handleEffectDisablerAbility handles the Effect Disabler special ability.
func (h *Handler) handleEffectDisablerAbility(targetPlayer *player.Player, left item.Stack) {
	if h.checkCoolDownSpecific(it.EffectDisablerType{}) {
		h.p.Message(text.Colourf("<red>You are on Effect Disabler cooldown for %.1f seconds</red>", h.coolDownSpecificAbilities.Key(it.EffectDisablerType{}).Remaining().Seconds()))
		return
	}

	h.applyEffectDisablerEffect(targetPlayer, left)
}

// checkCoolDownSpecific checks if the specific ability cooldown is active.
func (h *Handler) checkCoolDownSpecific(abilityType it.SpecialItemType) bool {
	cd := h.coolDownSpecificAbilities.Key(abilityType)
	return cd.Active()
}

// triggerStormBreakerEffect triggers the Storm Breaker special ability effect.
func (h *Handler) triggerStormBreakerEffect(target *player.Player, left item.Stack) {
	h.p.World().PlaySound(h.p.Position(), sound.ItemBreak{})
	h.p.World().AddEntity(h.p.World().EntityRegistry().Config().Lightning(h.p.Position()))
	h.coolDownSpecificAbilities.Set(it.StormBreakerType{}, time.Minute*2)
	h.coolDownGlobalAbilities.Set(time.Second * 10)

	targetArmourHandler, ok := h.targetPlayerArmourHandler(target)
	if !ok {
		return
	}

	h.p.SetHeldItems(item.Stack{}, left)
	targetArmourHandler.stormBreak()
}

// targetPlayerArmourHandler retrieves the armour handler of the target player.
func (h *Handler) targetPlayerArmourHandler(p *player.Player) (*ArmourHandler, bool) {
	targetArmour := p.Armour().Inventory()
	targetArmourHandler, ok := targetArmour.Handler().(*ArmourHandler)
	return targetArmourHandler, ok
}

// applyEffectDisablerEffect applies the Effect Disabler special ability effect.
func (h *Handler) applyEffectDisablerEffect(targetPlayer *player.Player, left item.Stack) {
	h.coolDownSpecificAbilities.Set(it.EffectDisablerType{}, time.Minute*2)
	h.coolDownGlobalAbilities.Set(time.Second * 10)

	if targetHandler, ok := targetPlayer.Handler().(*Handler); ok {
		targetHandler.clearEffects()
		targetHandler.coolDownEffectDisabled.Set(time.Second * 10)
		h.p.Message(text.Colourf("<red>You have effect disabled %s for 10 seconds</red>", targetPlayer.Name()))
		targetPlayer.Message(text.Colourf("<red>You have been effect disabled for 10 seconds by %s</red>", h.p.Name()))
	}
}

// applyComboAbility applies an additional second of strength if has combo ability active.
func (h *Handler) applyComboAbility() {
	// check if strength effect last 15 seconds
	effects := h.p.Effects()
	currentDuration := time.Duration(0)
	for _, e := range effects {
		if e.Type() == (effect.Strength{}) {
			currentDuration = e.Duration()
			break
		}
	}

	if currentDuration >= time.Second*15 {
		return
	}

	if h.coolDownComboAbility.Active() {
		h.p.AddEffect(effect.New(effect.Strength{}, 1, time.Second+currentDuration))
	}
}

// handleSpecialCases handles special cases like kits menu.
func (h *Handler) handleSpecialCases(targetPlayer *player.Player, ctx *event.Context) bool {
	if colour.StripMinecraftColour(targetPlayer.Name()) == "Click to use kits" {
		if menu, ok := menu.NewKitsMenu(h.p); ok {
			inv.SendMenu(h.p, menu)
		}
		ctx.Cancel()
		return true
	}
	return false
}

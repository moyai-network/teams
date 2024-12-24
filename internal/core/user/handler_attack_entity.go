package user

import (
	"github.com/moyai-network/teams/internal/core/area"
	"github.com/moyai-network/teams/internal/core/colour"
	data2 "github.com/moyai-network/teams/internal/core/data"
	ench "github.com/moyai-network/teams/internal/core/enchantment"
	item2 "github.com/moyai-network/teams/internal/core/item"
	menu2 "github.com/moyai-network/teams/internal/core/menu"
	"github.com/moyai-network/teams/internal/core/sotw"
	"github.com/moyai-network/teams/internal/core/user/class"
	"time"

	"github.com/bedrock-gophers/knockback/knockback"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// HandleAttackEntity handles the attack on an entity.
func (h *Handler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, _ *bool) {
	knockback.ApplyForce(force)
	knockback.ApplyHeight(height)

	targetPlayer, ok := e.(*player.Player)
	p := ctx.Val()
	if !ok || h.handleSpecialCases(p, targetPlayer, ctx) {
		return
	}

	h.ShowArmor(p, true)

	if !canAttack(p, targetPlayer) {
		ctx.Cancel()
		return
	}
	h.adjustKnockbackHeight(p, targetPlayer, height)
	h.applyAttackEnchantments(ctx, targetPlayer)
	h.handleRogueBackstab(ctx, targetPlayer, *force, *height)
	h.handleSpecialItemAbility(ctx, targetPlayer, force, height)
	h.setLastAttackerForCombat(p, targetPlayer)
	h.applyComboAbility(p)
}

// handleRogueBackstab handles the backstab ability for the Rogue class.
func (h *Handler) handleRogueBackstab(ctx *player.Context, t *player.Player, force, height float64) {
	p := ctx.Val()

	held, left := p.HeldItems()
	if s, ok := held.Item().(item.Sword); ok && s.Tier == item.ToolTierGold && class.Compare(h.lastClass.Load(), class.Rogue{}) && t.Rotation().Direction() == p.Rotation().Direction() {
		cd := h.coolDownBackStab
		if cd.Active() {
			p.Message(lang.Translatef(data2.Language{}, "user.cool-down", "Rogue", cd.Remaining().Seconds()))
		} else {
			ctx.Cancel()
			for i := 1; i <= 3; i++ {
				p.Tx().AddParticle(t.Position().Add(mgl64.Vec3{0, float64(i), 0}), particle.Dust{
					Colour: item.ColourRed().RGBA(),
				})
			}
			p.Tx().PlaySound(p.Position(), sound.ItemBreak{})
			t.Hurt(8, NoArmourAttackEntitySource{
				Attacker: p,
			})
			t.KnockBack(p.Position(), force, height)

			p.AddEffect(effect.New(effect.Slowness, 3, time.Second*5))
			p.SetHeldItems(item.Stack{}, left)
			cd.Set(time.Second * 10)
		}
	}
}

// applyAttackEnchantments applies attack enchantments from the attacker's armor.
func (h *Handler) applyAttackEnchantments(ctx *player.Context, targetPlayer *player.Player) {
	p := ctx.Val()

	arm := p.Armour()
	for _, a := range arm.Slots() {
		for _, e := range a.Enchantments() {
			if att, ok := e.Type().(ench.AttackEnchantment); ok {
				att.AttackEntity(p, targetPlayer)
			}
		}
	}
}

// handleSpecialItemAbility handles special abilities if the attacker is holding a special item.
func (h *Handler) handleSpecialItemAbility(ctx *player.Context, targetPlayer *player.Player, force, height *float64) {
	p := ctx.Val()
	heldItem, left := p.HeldItems()
	if specialItemType, ok := item2.PartnerItem(heldItem); ok {
		switch specialItemType.(type) {
		case item2.PearlDisablerType:
			h.handlePearlDisablerAbility(p, targetPlayer, heldItem, left)
		case item2.ExoticBoneType:
			h.handleExoticBoneAbility(p, targetPlayer, heldItem, left)
		case item2.StormBreakerType:
			h.handleStormBreakerAbility(p, targetPlayer, left)
		case item2.EffectDisablerType:
			h.handleEffectDisablerAbility(p, targetPlayer, left)
		}
	} else if specialItemType, ok = item2.SpecialItem(heldItem); ok {
		switch specialItemType.(type) {
		case item2.StaffKnockBackStickType:
			*force = *force * 1.5
		case item2.StaffFreezeBlockType:
			p.ExecuteCommand("/freeze \"" + targetPlayer.Name() + "\"")
			ctx.Cancel()
		}
	}
}

// adjustKnockbackHeight adjusts the knockback height based on the target's airborne status.
func (h *Handler) adjustKnockbackHeight(p *player.Player, targetPlayer *player.Player, height *float64) {
	if !targetPlayer.OnGround() {
		mx, mn := maxMin(targetPlayer.Position().Y(), p.Position().Y())
		if mx-mn >= 2.5 {
			*height = 0.38 / 1.25
		}
	}
}

// setLastAttackerForCombat sets the last attacker for combat tagging.
func (h *Handler) setLastAttackerForCombat(p *player.Player, targetPlayer *player.Player) {
	targetHandler, ok := targetPlayer.Handler().(*Handler)
	if !ok {
		return
	}

	if canAttack(p, targetPlayer) {
		targetHandler.setLastAttacker(p)
		h.tagCombat.Set(time.Second * 30)
		targetHandler.tagCombat.Set(time.Second * 30)
	}
}

// canAttack returns true if the given players can attack each other.
func canAttack(pl, target *player.Player) bool {
	if target == nil || pl == nil || target.GameMode() == world.GameModeCreative {
		return false
	}
	w := pl.Tx().World()
	if ((area.Spawn(w).Vec3WithinOrEqualFloorXZ(pl.Position())) || area.Spawn(w).Vec3WithinOrEqualFloorXZ(target.Position())) && w != internal.Deathban() {
		return false
	}

	u, err := data2.LoadUserFromName(pl.Name())
	t, err2 := data2.LoadUserFromName(target.Name())

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

	if u.Teams.DeathBanned || t.Teams.DeathBanned {
		if area.Deathban.Spawn().Vec3WithinOrEqualFloorXZ(pl.Position()) || area.Deathban.Spawn().Vec3WithinOrEqualFloorXZ(target.Position()) {
			return false
		}
	}

	tm, err := data2.LoadTeamFromMemberName(pl.Name())
	if err != nil {
		return true
	}

	if tm.Member(target.Name()) {
		return false
	}

	return true
}

// handlePearlDisablerAbility handles the Pearl Disabler special ability.
func (h *Handler) handlePearlDisablerAbility(p *player.Player, targetPlayer *player.Player, held, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(item2.PearlDisablerType{}); cd.Active() {
		internal.Messagef(p, "pearl_disabler.cooldown", cd.Remaining().Seconds())
		return
	}

	internal.Messagef(targetPlayer, "pearl_disabler.target", p.Name())
	targetHandler, ok := targetPlayer.Handler().(*Handler)
	if !ok {
		return
	}
	targetHandler.coolDownPearl.Set(time.Second * 15)

	internal.Messagef(p, "pearl_disabler.user", targetPlayer.Name())
	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(item2.PearlDisablerType{}, time.Minute)
	p.SetHeldItems(subtractItem(p, held, 1), left)
}

// handleExoticBoneAbility handles the Exotic Bone special ability.
func (h *Handler) handleExoticBoneAbility(p *player.Player, targetPlayer *player.Player, held, left item.Stack) {
	if cd := h.coolDownSpecificAbilities.Key(item2.ExoticBoneType{}); cd.Active() {
		internal.Messagef(p, "bone.cooldown", cd.Remaining().Seconds())
		return
	}

	internal.Messagef(targetPlayer, "bone.target", p.Name())
	internal.Messagef(p, "bone.user", targetPlayer.Name())

	h.coolDownGlobalAbilities.Set(time.Second * 10)
	h.coolDownSpecificAbilities.Set(item2.ExoticBoneType{}, time.Minute)
	p.SetHeldItems(subtractItem(p, held, 1), left)

	targetHandler, ok := targetPlayer.Handler().(*Handler)
	if !ok {
		return
	}
	targetHandler.coolDownBonedEffect.Set(time.Second * 10)
}

// handleStormBreakerAbility handles the Storm Breaker special ability.
func (h *Handler) handleStormBreakerAbility(p *player.Player, targetPlayer *player.Player, left item.Stack) {
	targetHandler, ok := targetPlayer.Handler().(*Handler)
	if !ok || targetHandler.lastClass.Load() != nil || h.checkCoolDownSpecific(item2.StormBreakerType{}) {
		return
	}

	h.triggerStormBreakerEffect(p, targetPlayer, left)
}

// handleEffectDisablerAbility handles the Effect Disabler special ability.
func (h *Handler) handleEffectDisablerAbility(p *player.Player, targetPlayer *player.Player, left item.Stack) {
	if h.checkCoolDownSpecific(item2.EffectDisablerType{}) {
		p.Message(text.Colourf("<red>You are on Effect Disabler cooldown for %.1f seconds</red>", h.coolDownSpecificAbilities.Key(item2.EffectDisablerType{}).Remaining().Seconds()))
		return
	}

	h.applyEffectDisablerEffect(p, targetPlayer, left)
}

// checkCoolDownSpecific checks if the specific ability cooldown is active.
func (h *Handler) checkCoolDownSpecific(abilityType item2.SpecialItemType) bool {
	cd := h.coolDownSpecificAbilities.Key(abilityType)
	return cd.Active()
}

// triggerStormBreakerEffect triggers the Storm Breaker special ability effect.
func (h *Handler) triggerStormBreakerEffect(p *player.Player, target *player.Player, left item.Stack) {
	p.Tx().PlaySound(p.Position(), sound.ItemBreak{})
	p.Tx().AddEntity(p.Tx().World().EntityRegistry().Config().Lightning(world.EntitySpawnOpts{Position: p.Position()}))
	h.coolDownSpecificAbilities.Set(item2.StormBreakerType{}, time.Minute*2)
	h.coolDownGlobalAbilities.Set(time.Second * 10)

	targetArmourHandler, ok := h.targetPlayerArmourHandler(target)
	if !ok {
		return
	}

	p.SetHeldItems(item.Stack{}, left)
	targetArmourHandler.stormBreak()
}

// targetPlayerArmourHandler retrieves the armour handler of the target player.
func (h *Handler) targetPlayerArmourHandler(p *player.Player) (*ArmourHandler, bool) {
	targetArmour := p.Armour().Inventory()
	targetArmourHandler, ok := targetArmour.Handler().(*ArmourHandler)
	return targetArmourHandler, ok
}

// applyEffectDisablerEffect applies the Effect Disabler special ability effect.
func (h *Handler) applyEffectDisablerEffect(p *player.Player, targetPlayer *player.Player, left item.Stack) {
	h.coolDownSpecificAbilities.Set(item2.EffectDisablerType{}, time.Minute*2)
	h.coolDownGlobalAbilities.Set(time.Second * 10)

	if targetHandler, ok := targetPlayer.Handler().(*Handler); ok {
		targetHandler.clearEffects(p)
		targetHandler.coolDownEffectDisabled.Set(time.Second * 10)
		p.Message(text.Colourf("<red>You have effect disabled %s for 10 seconds</red>", targetPlayer.Name()))
		targetPlayer.Message(text.Colourf("<red>You have been effect disabled for 10 seconds by %s</red>", p.Name()))
	}
}

// applyComboAbility applies an additional second of strength if has combo ability active.
func (h *Handler) applyComboAbility(p *player.Player) {
	// check if strength effect last 15 seconds
	effects := p.Effects()
	currentDuration := time.Duration(0)
	for _, e := range effects {
		if e.Type() == (effect.Strength) {
			currentDuration = e.Duration()
			break
		}
	}

	if currentDuration >= time.Second*15 {
		return
	}

	if h.coolDownComboAbility.Active() {
		p.AddEffect(effect.New(effect.Strength, 1, time.Second+currentDuration))
	}
}

// handleSpecialCases handles special cases like kits menu.
func (h *Handler) handleSpecialCases(p *player.Player, targetPlayer *player.Player, ctx *player.Context) bool {
	if colour.StripMinecraftColour(targetPlayer.Name()) == "Click to use kits" {
		if mn, ok := menu2.NewKitsMenu(p); ok {
			inv.SendMenu(p, mn)
		}
		ctx.Cancel()
		return true
	} else if colour.StripMinecraftColour(targetPlayer.Name()) == "Block Shop" {
		inv.SendMenu(p, menu2.NewBlocksMenu(p))
		ctx.Cancel()
		return true
	}
	return false
}

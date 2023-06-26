package user

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/class"
	ench "github.com/moyai-network/teams/moyai/enchantment"
)

type NopArmourTier struct{}

func (NopArmourTier) BaseDurability() float64      { return 0 }
func (NopArmourTier) Toughness() float64           { return 0 }
func (NopArmourTier) KnockBackResistance() float64 { return 0 }
func (NopArmourTier) EnchantmentValue() int        { return 0 }
func (NopArmourTier) Name() string                 { return "" }

type ArmourHandler struct {
	p *player.Player
}

func NewArmourHandler(p *player.Player) *ArmourHandler {
	return &ArmourHandler{p: p}
}

func (a *ArmourHandler) HandleTake(_ *event.Context, _ int, it item.Stack) {
	setClass(a.p, nil)

	var effects []effect.Effect

	for _, e := range it.Enchantments() {
		if enc, ok := e.Type().(ench.EffectEnchantment); ok {
			effects = append(effects, enc.Effect())
		}
	}

	for _, e := range effects {
		typ := e.Type()
		if hasEffectLevel(a.p, e) {
			a.p.RemoveEffect(typ)
		}
	}
}
func (a *ArmourHandler) HandlePlace(_ *event.Context, _ int, it item.Stack) {
	var helmetTier item.ArmourTier = NopArmourTier{}
	var chestplateTier item.ArmourTier = NopArmourTier{}
	var leggingsTier item.ArmourTier = NopArmourTier{}
	var bootsTier item.ArmourTier = NopArmourTier{}

	arm := a.p.Armour()
	helmet, ok := arm.Helmet().Item().(item.Helmet)
	if ok {
		helmetTier = helmet.Tier
	}
	chestplate, ok := arm.Chestplate().Item().(item.Chestplate)
	if ok {
		chestplateTier = chestplate.Tier
	}
	leggings, ok := arm.Leggings().Item().(item.Leggings)
	if ok {
		leggingsTier = leggings.Tier
	}
	boots, ok := arm.Boots().Item().(item.Boots)
	if ok {
		bootsTier = boots.Tier
	}
	switch it := it.Item().(type) {
	case item.Helmet:
		helmetTier = it.Tier
	case item.Chestplate:
		chestplateTier = it.Tier
	case item.Leggings:
		leggingsTier = it.Tier
	case item.Boots:
		bootsTier = it.Tier
	}
	newArmour := [4]item.ArmourTier{helmetTier, chestplateTier, leggingsTier, bootsTier}
	setClass(a.p, class.ResolveFromArmour(newArmour))

	var effects []effect.Effect

	for _, e := range it.Enchantments() {
		if enc, ok := e.Type().(ench.EffectEnchantment); ok {
			effects = append(effects, enc.Effect())
		}
	}

	for _, e := range effects {
		a.p.AddEffect(e)
	}
}
func (a *ArmourHandler) HandleDrop(_ *event.Context, _ int, it item.Stack) {
	setClass(a.p, nil)

	var effects []effect.Effect

	for _, e := range it.Enchantments() {
		if enc, ok := e.Type().(ench.EffectEnchantment); ok {
			effects = append(effects, enc.Effect())
		}
	}

	for _, e := range effects {
		typ := e.Type()
		if hasEffectLevel(a.p, e) {
			a.p.RemoveEffect(typ)
		}
	}
}

// hasEffectLevel returns whether the user has the effect or not.
func hasEffectLevel(p *player.Player, e effect.Effect) bool {
	for _, ef := range p.Effects() {
		if e.Type() == ef.Type() && e.Level() == ef.Level() {
			return true
		}
	}
	return false
}

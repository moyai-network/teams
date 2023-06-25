package user

import (
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/class"
	ench "github.com/moyai-network/teams/moyai/enchantment"
)

type NopArmourTier struct{}

func (NopArmourTier) BaseDurability() float64      { return 0 }
func (NopArmourTier) Toughness() float64           { return 0 }
func (NopArmourTier) KnockBackResistance() float64 { return 0 }
func (NopArmourTier) EnchantmentValue() int        { return 0 }
func (NopArmourTier) Name() string                 { return "" }

type ClassHandler struct {
	p *player.Player

	class      atomic.Value[moose.Class]
	bardEnergy atomic.Value[float64]
}

func NewClassHandler(p *player.Player) *ClassHandler {
	return &ClassHandler{p: p}
}

func (h *ClassHandler) HandleTake(_ *event.Context, _ int, it item.Stack) {
	h.setClass(nil)

	var effects []effect.Effect

	for _, e := range it.Enchantments() {
		if enc, ok := e.Type().(ench.EffectEnchantment); ok {
			effects = append(effects, enc.Effect())
		}
	}

	for _, e := range effects {
		typ := e.Type()
		if hasEffectLevel(h.p, e) {
			h.p.RemoveEffect(typ)
		}
	}
}
func (h *ClassHandler) HandlePlace(_ *event.Context, _ int, it item.Stack) {
	var helmetTier item.ArmourTier = NopArmourTier{}
	var chestplateTier item.ArmourTier = NopArmourTier{}
	var leggingsTier item.ArmourTier = NopArmourTier{}
	var bootsTier item.ArmourTier = NopArmourTier{}

	a := h.p.Armour()
	helmet, ok := a.Helmet().Item().(item.Helmet)
	if ok {
		helmetTier = helmet.Tier
	}
	chestplate, ok := a.Chestplate().Item().(item.Chestplate)
	if ok {
		chestplateTier = chestplate.Tier
	}
	leggings, ok := a.Leggings().Item().(item.Leggings)
	if ok {
		leggingsTier = leggings.Tier
	}
	boots, ok := a.Boots().Item().(item.Boots)
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
	h.setClass(class.ResolveFromArmour(newArmour))

	var effects []effect.Effect

	for _, e := range it.Enchantments() {
		if enc, ok := e.Type().(ench.EffectEnchantment); ok {
			effects = append(effects, enc.Effect())
		}
	}

	for _, e := range effects {
		h.p.AddEffect(e)
	}
}
func (h *ClassHandler) HandleDrop(_ *event.Context, _ int, it item.Stack) {
	h.setClass(nil)

	var effects []effect.Effect

	for _, e := range it.Enchantments() {
		if enc, ok := e.Type().(ench.EffectEnchantment); ok {
			effects = append(effects, enc.Effect())
		}
	}

	for _, e := range effects {
		typ := e.Type()
		if hasEffectLevel(h.p, e) {
			h.p.RemoveEffect(typ)
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

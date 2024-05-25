package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/role"
)

// Trim is a command that allows the player to add trims to their armor
type Trim struct {
	Template template `cmd:"template"`
	Material material `cmd:"material"`
}

// Trim is a command that allows the player to clear trims to their armor
type TrimClear struct {
	Sub      cmd.SubCommand       `cmd:"clear"`
}

// Run ...
func (t Trim) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	arm := p.Armour()
	trim := item.ArmourTrim{
		Template: nameToTemplate(string(t.Template)),
		Material: item.ArmourTrimMaterialFromString(string(t.Material)),
	}

	arm.Set(p.Armour().Helmet().WithArmourTrim(trim), p.Armour().Chestplate().WithArmourTrim(trim), p.Armour().Leggings().WithArmourTrim(trim), p.Armour().Boots().WithArmourTrim(trim))
}

// Allow ...
func (Trim) Allow(src cmd.Source) bool {
	return allow(src, false, role.Donor1{})
}

// Run ...
func (TrimClear) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	arm := p.Armour()

	arm.Set(arm.Helmet().WithArmourTrim(item.ArmourTrim{}), arm.Chestplate().WithArmourTrim(item.ArmourTrim{}), arm.Leggings().WithArmourTrim(item.ArmourTrim{}), arm.Boots().WithArmourTrim(item.ArmourTrim{}))
}
type template string
type material string

// Type ...
func (template) Type() string {
	return "template"
}

// Options ...
func (template) Options(s cmd.Source) []string {
	return []string{
		"netherite_upgrade",
		"sentry",
		"vex",
		"wild",
		"coast",
		"dune",
		"wayfinder",
		"raiser",
		"shaper",
		"host",
		"ward",
		"silence",
		"tide",
		"snout",
		"rib",
		"eye",
		"spire",
		"flow",
		"bolt",
	}
}

// Type ...
func (material) Type() string {
	return "material"
}

func nameToTemplate(name string) item.ArmourSmithingTemplate {
	switch name {
	case "netherite_upgrade":
		return item.TemplateNetheriteUpgrade()
	case "sentry":
		return item.TemplateSentry()
	case "vex":
		return item.TemplateVex()
	case "wild":
		return item.TemplateWild()
	case "coast":
		return item.TemplateCoast()
	case "dune":
		return item.TemplateDune()
	case "wayfinder":
		return item.TemplateWayFinder()
	case "raiser":
		return item.TemplateRaiser()
	case "shaper":
		return item.TemplateShaper()
	case "host":
		return item.TemplateHost()
	case "ward":
		return item.TemplateWard()
	case "silence":
		return item.TemplateSilence()
	case "tide":
		return item.TemplateTide()
	case "snout":
		return item.TemplateSnout()
	case "rib":
		return item.TemplateRib()
	case "eye":
		return item.TemplateEye()
	case "spire":
		return item.TemplateSpire()
	case "flow":
		return item.TemplateFlow()
	case "bolt":
		return item.TemplateBolt()
	}
	panic("should not happen")
}
// Options ...
func (material) Options(s cmd.Source) []string {
	return []string{
		"diamond",
		"emerald",
		"copper",
		"gold",
		"iron",
		"netherite",
		"amethyst",
		"lapis",
		"quartz",
	}
}
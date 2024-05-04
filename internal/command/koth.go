package command

import (
	"github.com/moyai-network/teams/internal/colour"
	"github.com/moyai-network/teams/internal/data"
	"github.com/moyai-network/teams/internal/koth"
	"github.com/moyai-network/teams/internal/role"
	"github.com/moyai-network/teams/moyai"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// KothList is a command that lists all KOTHs.
type KothList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

// KothStart is a command that starts a KOTH.
type KothStart struct {
	Sub  cmd.SubCommand `cmd:"start"`
	KOTH kothList       `cmd:"koth"`
}

// KothStop is a command that stops a KOTH.
type KothStop struct {
	Sub cmd.SubCommand `cmd:"stop"`
}

// Run ...
func (KothList) Run(s cmd.Source, o *cmd.Output) {
	all := []string{
		text.Colourf("<yellow>KOTHs</yellow><grey>:</grey>"),
	}
	for _, k := range koth.All() {
		coords := k.Coordinates()
		all = append(all, text.Colourf("%s<grey>:</grey> <yellow>%0.f, %0.f</yellow>", k.Name(), coords.X(), coords.Y()))
	}
	o.Print(strings.Join(all, "\n"))
}

// Run ...
func (k KothStart) Run(s cmd.Source, o *cmd.Output) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	if p, ok := s.(*player.Player); ok {
		if u, err := data.LoadUserFromName(p.Name()); err == nil {
			r := u.Roles.Highest()
			name = r.Color(p.Name())
		}
	}
	if _, ok := koth.Running(); ok {
		o.Print(text.Colourf("<red>A KOTH is already running.</red>"))
		return
	}
	ko, ok := koth.Lookup(string(k.KOTH))
	if !ok {
		o.Print(text.Colourf("<red>Invalid KOTH.</red>"))
		return
	}
	ko.Start()

	for _, u := range moyai.Server().Players() {
		if ko.Area().Vec3WithinOrEqualXZ(u.Position()) {
			//ko.StartCapturing(u)
		}
	}

	coords := ko.Coordinates()
	user.Broadcastf("koth.start", name, ko.Name(), coords.X(), coords.Y())
}

// Run ...
func (KothStop) Run(s cmd.Source, o *cmd.Output) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	if p, ok := s.(*player.Player); ok {
		if u, err := data.LoadUserFromName(p.Name()); err == nil {
			r := u.Roles.Highest()
			name = r.Color(p.Name())
		}
	}
	if k, ok := koth.Running(); !ok {
		o.Print(text.Colourf("<red>No KOTH is running.</red>"))
		return
	} else {
		k.Stop()
		user.Broadcastf("koth.stop", name, k.Name())
	}
}

// Allow ...
func (KothStart) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

// Allow ...
func (KothStop) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

type (
	kothList string
)

// Type ...
func (kothList) Type() string {
	return "koth_list"
}

// Options ...
func (kothList) Options(cmd.Source) []string {
	var opts []string
	for _, k := range koth.All() {
		opts = append(opts, colour.StripMinecraftColour(strings.ToLower(k.Name())))
	}
	return opts
}

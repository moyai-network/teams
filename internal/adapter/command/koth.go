package command

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/core/colour"
	"github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/eotw"
	"github.com/moyai-network/teams/internal/core/koth"
	rls "github.com/moyai-network/teams/internal/core/roles"
	"github.com/moyai-network/teams/internal/core/sotw"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/hako/durafmt"
	"github.com/moyai-network/teams/internal"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// KothList is a command that lists all KOTHs.
type KothList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

// KothStart is a command that starts a KOTH.
type KothStart struct {
	donor1Allower
	Sub  cmd.SubCommand `cmd:"start"`
	KOTH kothList       `cmd:"koth"`
}

// KothStop is a command that stops a KOTH.
type KothStop struct {
	adminAllower
	Sub cmd.SubCommand `cmd:"stop"`
}

// Run ...
func (KothList) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
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
func (k KothStart) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	_, sotwRunning := sotw.Running()
	if sotwRunning && !u.Roles.Contains(rls.Operator()) {
		internal.Messagef(p, "command.koth.sotw")
		return
	}
	_, eotwRunning := eotw.Running()
	if eotwRunning && !u.Roles.Contains(rls.Operator()) {
		internal.Messagef(p, "command.koth.eotw")
		return
	}

	r := u.Roles.Highest()
	name = r.Coloured(p.Name())
	if u.Teams.KOTHStart.Active() {
		internal.Messagef(p, "command.koth.cooldown", durafmt.ParseShort(u.Teams.KOTHStart.Remaining()).LimitFirstN(2))
		return
	}

	if _, ok := koth.Running(); ok {
		internal.Messagef(p, "command.koth.running")
		return
	}

	ko, ok := koth.Lookup(string(k.KOTH))
	if !ok {
		internal.Messagef(p, "command.koth.invalid")
		return
	}

	ko.Start()
	if !u.Roles.Contains(rls.Operator(), rls.Admin(), rls.Manager()) {
		u.Teams.KOTHStart.Set(time.Hour * 48)
	}

	for u := range internal.Players(tx) {
		if ko.Area().Vec3WithinOrEqualXZ(u.Position()) {
			ko.StartCapturing(u)
		}
	}

	coords := ko.Coordinates()
	internal.Broadcastf(tx, "koth.start", name, ko.Name(), coords.X(), coords.Y())
	var st string
	if ko == koth.Citadel {
		st = fmt.Sprintf(`
 §e█████████§r
 §e█████████§r
 §e█§6█§e█§6█§e█§6█§e█§6█§e█§r
 §e█§6███████§e█§r
 §e█§6█§b█§6█§b█§6█§b█§6█§e█§r §e%s§r
 §e█§6███████§e█§r §6can be contested now!§r
 §e█████████§r
 §e█████████§r
 §e█████████§r

`, ko.Name())
	} else {
		st = fmt.Sprintf(`
 §7█████████§r
 §7██§4█§7███§4█§7██§r
 §7██§4█§7██§4█§7███§r
 §7██§4███§7████§r
 §7██§4█§7██§4█§7███ §e%s KOTH§r
 §7██§4█§7███§4█§7██§r §6can be contested now!§r
 §7██§4█§7███§4█§7██§r
 §7█████████§r
`, ko.Name())
	}

	p.Message(text.Colourf(st))
}

// Run ...
func (KothStop) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	p, ok := s.(*player.Player)
	if ok {
		if u, err := data.LoadUserFromName(p.Name()); err == nil {
			r := u.Roles.Highest()
			name = r.Coloured(p.Name())
		}
	}
	if k, ok := koth.Running(); !ok {
		internal.Messagef(p, "command.koth.not.running")
		return
	} else {
		k.Stop()
		internal.Broadcastf(tx, "koth.stop", name, k.Name())
	}
}

type (
	kothList string
)

// Type ...
func (kothList) Type() string {
	return "koth_list"
}

// Options ...
func (kothList) Options(src cmd.Source) []string {
	p, playerSrc := src.(*player.Player)
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return []string{}
	}

	var opts []string
	for _, k := range koth.All() {
		if k == koth.Citadel && (playerSrc && !u.Roles.Contains(rls.Operator(), rls.Admin(), rls.Manager())) {
			continue
		}
		opts = append(opts, colour.StripMinecraftColour(strings.ToLower(k.Name())))
	}
	return opts
}

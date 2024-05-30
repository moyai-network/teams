package command

import (
	"fmt"

	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/conquest"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// ConquestStart is a command that starts a KOTH.
type ConquestStart struct {
	Sub  cmd.SubCommand `cmd:"start"`
}

// ConquestStop is a command that stops a KOTH.
type ConquestStop struct {
	Sub cmd.SubCommand `cmd:"stop"`
}

// Run ...
func (k ConquestStart) Run(s cmd.Source, o *cmd.Output) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	p, ok := s.(*player.Player)
	if ok {
		if u, err := data.LoadUserFromName(p.Name()); err == nil {
			r := u.Roles.Highest()
			name = r.Color(p.Name())
		}
	}
	if conquest.Running() {
		user.Messagef(p, "command.koth.running")
		return
	}
	conquest.Start()

	areas := conquest.All()
	for _, u := range moyai.Players() {
		for _, a := range areas {
			if a.Area().Vec3WithinOrEqualXZ(u.Position()) {
				a.StartCapturing(u)
			}
		}
	}

	user.Broadcastf("koth.start", name, "Conquest", -100.0, -500.0)
	st := fmt.Sprintf(`
§e█████████§r
§e█████████§r
§e█§6█§e█§6█§e█§6█§e█§6█§e█§r
§e█§6███████§e█§r
§e█§6█§b█§6█§b█§6█§b█§6█§e█§r §eConquest§r
§e█§6███████§e█§r §6can be contested now!§r
§e█████████§r
§e█████████§r
§e█████████§r

`)

	p.Message(text.Colourf(st))
}

// Run ...
func (ConquestStop) Run(s cmd.Source, o *cmd.Output) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	p, ok := s.(*player.Player)
	if ok {
		if u, err := data.LoadUserFromName(p.Name()); err == nil {
			r := u.Roles.Highest()
			name = r.Color(p.Name())
		}
	}
	if !conquest.Running() {
		user.Messagef(p, "command.koth.not.running")
	} else {
		conquest.Stop()
		user.Broadcastf("koth.stop", name, "Conquest")
	}
}

// Allow ...
func (ConquestStart) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

// Allow ...
func (ConquestStop) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

package command

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/core"
	rls "github.com/moyai-network/teams/internal/core/roles"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/pkg/lang"
)

// TeleportToPos is a command that teleports the user to a position.
type TeleportToPos struct {
	adminAllower
	Position mgl64.Vec3 `cmd:"destination"`
}

// TeleportToTarget is a command that teleports the user to another player.
type TeleportToTarget struct {
	trialAllower
	Targets []cmd.Target `cmd:"destination"`
}

// TeleportTargetsToTarget is a command that teleports player(s) to another player.
type TeleportTargetsToTarget struct {
	adminAllower
	Targets  []cmd.Target `cmd:"target"`
	Position []cmd.Target `cmd:"destination"`
}

// TeleportTargetsToPos is a command that teleports player(s) to a position.
type TeleportTargetsToPos struct {
	adminAllower
	Targets  []cmd.Target `cmd:"target"`
	Position mgl64.Vec3   `cmd:"destination"`
}

// Run ...
func (t TeleportToPos) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p := s.(*player.Player)
	p.Teleport(t.Position)
	internal.Messagef(p, "command.teleport.self", t.Position)
}

// Run ...
func (tp TeleportToTarget) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	t, ok := tp.Targets[0].(*player.Player)
	if !ok {
		internal.Messagef(p, "command.target.unknown")
		return
	}
	if p.Tx().World() != t.Tx().World() {
		p.Tx().AddEntity(t.H())
	}
	p.Teleport(t.Position())
	internal.Messagef(p, "command.teleport.self", t.Name())
}

// Run ...
func (tp TeleportTargetsToTarget) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		// Somehow left midway through the process, so just return.
		return
	}

	if len(tp.Targets) > 1 && !u.Roles.Contains(rls.Operator()) {
		o.Print(lang.Translatef(l, "command.teleport.operator"))
		return
	}

	if len(tp.Position) > 1 {
		o.Print(lang.Translatef(l, "command.teleport.multiple"))
		return
	}
	t, ok := tp.Position[0].(*player.Player)
	if !ok {
		o.Print(lang.Translatef(l, "command.target.unknown"))
		return
	}
	o.Print(lang.Translatef(l, "command.teleport.target", teleportTargets(tp.Targets, t.Position(), t), t.Name()))
}

// Run ...
func (t TeleportTargetsToPos) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		// Somehow left midway through the process, so just return.
		return
	}

	if len(t.Targets) > 1 && !u.Roles.Contains(rls.Operator()) {
		o.Print(lang.Translatef(l, "command.teleport.operator"))
		return
	}
	o.Print(lang.Translatef(l, "command.teleport.target", teleportTargets(t.Targets, t.Position, p)), t.Position)
}

// teleportTargets teleports a list of targets to a specified position and world. If the world is nil, it will only
// teleport to the position. If the position is empty, it will only teleport to the world of the player. It returns the
// players affected in a readable string.
func teleportTargets(targets []cmd.Target, destination mgl64.Vec3, t *player.Player) string {
	for _, target := range targets {
		if p, ok := target.(*player.Player); ok {
			if p.Tx().World() != t.Tx().World() {
				t.Tx().AddEntity(p.H())
			}
			p.Teleport(destination)
		}
	}
	if l := len(targets); l > 1 {
		return fmt.Sprintf("%d players", l)
	}
	return targets[0].(cmd.NamedTarget).Name()
}

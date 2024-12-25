package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/sotw"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// SOTWStart is a command to start SOTW.
type SOTWStart struct {
	managerAllower
	Sub cmd.SubCommand `cmd:"start"`
}

// SOTWEnd is a command to end SOTW.
type SOTWEnd struct {
	managerAllower
	Sub cmd.SubCommand `cmd:"end"`
}

// SOTWDisable is a command to disable SOTW.
type SOTWDisable struct {
	Sub cmd.SubCommand `cmd:"disable"`
}

// Run ...
func (c SOTWStart) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	if _, ok := sotw.Running(); ok {
		o.Print(text.Colourf("<red>SOTW is already running!</red>"))
		return
	}
	sotw.Start()

	users := core.UserRepository.FindAll()
	for u := range users {
		u.Teams.SOTW = true
		core.UserRepository.Save(u)
	}
	internal.Broadcastf(tx, "sotw.commenced")
}

// Run ...
func (c SOTWEnd) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	if _, ok := sotw.Running(); !ok {
		o.Print(lang.Translatef(locale(s), "command.sotw.not.running"))
		return
	}
	sotw.End()

	users := core.UserRepository.FindAll()
	for u := range users {
		u.Teams.SOTW = false
		core.UserRepository.Save(u)
	}
	internal.Broadcastf(tx, "sotw.ended")
}

// Run ...
func (c SOTWDisable) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	if !u.Teams.SOTW {
		internal.Messagef(p, "sotw.disabled.already")
		return
	}
	internal.Messagef(p, "sotw.disabled")

	u.Teams.SOTW = false
	core.UserRepository.Save(u)
}

// Allow ...
func (SOTWDisable) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

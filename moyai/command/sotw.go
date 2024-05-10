package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// SOTWStart is a command to start SOTW.
type SOTWStart struct {
	Sub cmd.SubCommand `cmd:"start"`
}

// SOTWEnd is a command to end SOTW.
type SOTWEnd struct {
	Sub cmd.SubCommand `cmd:"end"`
}

// SOTWDisable is a command to disable SOTW.
type SOTWDisable struct {
	Sub cmd.SubCommand `cmd:"disable"`
}

// Run ...
func (c SOTWStart) Run(s cmd.Source, o *cmd.Output) {
	if _, ok := sotw.Running(); ok {
		o.Print(text.Colourf("<red>SOTW is already running!</red>"))
		return
	}
	sotw.Start()

	users, err := data.LoadAllUsers()
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		u.Teams.SOTW = true
		data.SaveUser(u)
	}
	user.Broadcast("sotw.commenced")
}

// Run ...
func (c SOTWEnd) Run(s cmd.Source, o *cmd.Output) {
	if _, ok := sotw.Running(); !ok {
		o.Print(text.Colourf("<red>SOTW is not running!</red>"))
		return
	}
	sotw.End()

	users, err := data.LoadAllUsers()
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		u.Teams.SOTW = false
		data.SaveUser(u)
	}
	user.Broadcast("sotw.ended")
}

// Run ...
func (c SOTWDisable) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	if !u.Teams.SOTW {
		user.Messagef(p, "sotw.disabled.already")
		return
	}
	user.Messagef(p, "sotw.disabled")

	u.Teams.SOTW = false
	data.SaveUser(u)
}

// Allow ...
func (SOTWStart) Allow(s cmd.Source) bool {
	return allow(s, true, role.Manager{})
}

// Allow ...
func (SOTWEnd) Allow(s cmd.Source) bool {
	return allow(s, true, role.Manager{})
}

// Allow ...
func (SOTWDisable) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

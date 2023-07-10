package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.mongodb.org/mongo-driver/bson"
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

	offline, err := data.LoadUsersCond(bson.M{})
	if err != nil {
		panic(err)
	}
	for _, u := range offline {
		u.SOTW = true
		_ = data.SaveUser(u)
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

	offline, err := data.LoadUsersCond(bson.M{})
	if err != nil {
		panic(err)
	}
	for _, u := range offline {
		u.SOTW = false
		err = data.SaveUser(u)
		if err != nil {
			panic(err)
		}
	}
	user.Broadcast("sotw.ended")
}

// Run ...
func (c SOTWDisable) Run(s cmd.Source, o *cmd.Output) {
	h, ok := user.Lookup(s.(*player.Player).Name())
	if !ok {
		return
	}
	u, err := data.LoadUserOrCreate(h.Player().Name())
	if err != nil {
		return
	}
	if !u.SOTW {
		h.Message("sotw.disabled.already")
		return
	}
	h.Message("sotw.disabled")

	u.SOTW = false
	_ = data.SaveUser(u)
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

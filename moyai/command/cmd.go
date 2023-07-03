package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/data"
	"golang.org/x/text/language"
)

// locale returns the locale of a cmd.Source.
func locale(s cmd.Source) language.Tag {
	if p, ok := s.(*player.Player); ok {
		return p.Locale()
	}
	return language.English
}

// allow is a helper function for command allowers. It allows users to easily check for the specified roles.
func allow(src cmd.Source, console bool, roles ...moose.Role) bool {
	p, ok := src.(*player.Player)
	if !ok {
		return console
	}
	if roles == nil {
		return true
	}
	u, err := data.LoadUser(p.Name(), p.XUID())
	if err != nil {
		return false
	}
	return ok && u.Roles.Contains(append(roles, role.Operator{})...)
}

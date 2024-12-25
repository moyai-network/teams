package command

import (
	"github.com/bedrock-gophers/role/role"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/core"
	rls "github.com/moyai-network/teams/internal/core/roles"
	model2 "github.com/moyai-network/teams/internal/model"
	"golang.org/x/text/language"
)

// locale returns the locale of a cmd.Source.
func locale(s cmd.Source) model2.Language {
	l := model2.Language{
		Tag: language.English,
	}

	if p, ok := s.(*player.Player); ok {
		u, ok := core.UserRepository.FindByName(p.Name())
		if !ok {
			return l
		}
		return *u.Language
	}
	return l
}

// Allow is a helper function for command allowers. It allows users to easily check for the specified roles.
func Allow(src cmd.Source, console bool, roles ...role.Role) bool {
	p, ok := src.(*player.Player)
	if !ok {
		return console
	}
	if roles == nil {
		return true
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return false
	}
	return u.Roles.Contains(append(roles, rls.Operator())...)
}

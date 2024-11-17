package command

import (
	"github.com/bedrock-gophers/role/role"
	"github.com/df-mc/dragonfly/server/cmd"
	rls "github.com/moyai-network/teams/internal/roles"
)

// adminAllower is an allower that allows all users with the admin role to execute a command.
type adminAllower struct{}

// Allow ...
func (adminAllower) Allow(s cmd.Source) bool {
	return Allow(s, true, rls.Admin())
}

// donor1Allower is an allower that allows all users with the donor1 role to execute a command.
type donor1Allower struct{}

// Allow ...
func (donor1Allower) Allow(s cmd.Source) bool {
	return Allow(s, true, rls.Khufu())
}

// trialAllower is an allower that allows all users with the trial role to execute a command.
type trialAllower struct{}

// Allow ...
func (trialAllower) Allow(s cmd.Source) bool {
	return Allow(s, true, rls.Trial())
}

// modAllower is an allower that allows all users with the mod role to execute a command.
type modAllower struct{}

// Allow ...
func (modAllower) Allow(s cmd.Source) bool {
	return Allow(s, true, rls.Mod())
}

// managerAllower is an allower that allows all users with the manager role to execute a command.
type managerAllower struct{}

// Allow ...
func (managerAllower) Allow(s cmd.Source) bool {
	return Allow(s, true, rls.Manager())
}

// operatorAllower is an allower that allows all users with the operator role to execute a command.
type operatorAllower struct{}

// Allow ...
func (operatorAllower) Allow(s cmd.Source) bool {
	return Allow(s, true, []role.Role{}...)
}

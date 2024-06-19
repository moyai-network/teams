package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/moyai/role"
)

// adminAllower is an allower that allows all users with the admin role to execute a command.
type adminAllower struct{}

// Allow ...
func (adminAllower) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

// donor1Allower is an allower that allows all users with the donor1 role to execute a command.
type donor1Allower struct{}

// Allow ...
func (donor1Allower) Allow(s cmd.Source) bool {
	return allow(s, true, role.Khufu{})
}

// trialAllower is an allower that allows all users with the trial role to execute a command.
type trialAllower struct{}

// Allow ...
func (trialAllower) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

// modAllower is an allower that allows all users with the mod role to execute a command.
type modAllower struct{}

// Allow ...
func (modAllower) Allow(s cmd.Source) bool {
	return allow(s, true, role.Mod{})
}

// managerAllower is an allower that allows all users with the manager role to execute a command.
type managerAllower struct{}

// Allow ...
func (managerAllower) Allow(s cmd.Source) bool {
	return allow(s, true, role.Manager{})
}

// operatorAllower is an allower that allows all users with the operator role to execute a command.
type operatorAllower struct{}

// Allow ...
func (operatorAllower) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

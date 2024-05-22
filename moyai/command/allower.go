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

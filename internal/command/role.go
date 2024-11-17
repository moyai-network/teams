package command

import (
	"github.com/bedrock-gophers/role/role"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/data"
	rls "github.com/moyai-network/teams/internal/roles"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/moyai-network/teams/pkg/timeutil"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/hako/durafmt"
)

// RoleAdd is a command to add a role to a player.
type RoleAdd struct {
	managerAllower
	Sub      cmd.SubCommand       `cmd:"add"`
	Targets  []cmd.Target         `cmd:"target"`
	Role     roles                `cmd:"role"`
	Duration cmd.Optional[string] `cmd:"duration"`
}

// RoleRemove is a command to remove a role from a player.
type RoleRemove struct {
	managerAllower
	Sub     cmd.SubCommand `cmd:"remove"`
	Targets []cmd.Target   `cmd:"target"`
	Role    roles          `cmd:"role"`
}

// RoleAddOffline is a command to remove a role from an offline user.
type RoleAddOffline struct {
	managerAllower
	Sub      cmd.SubCommand       `cmd:"add"`
	Target   string               `cmd:"target"`
	Role     roles                `cmd:"role"`
	Duration cmd.Optional[string] `cmd:"duration"`
}

// RoleRemoveOffline is a command to remove a role from an offline user.
type RoleRemoveOffline struct {
	managerAllower
	Sub    cmd.SubCommand `cmd:"remove"`
	Target string         `cmd:"target"`
	Role   roles          `cmd:"role"`
}

// RoleList is a command to list all users with a role.
type RoleList struct {
	managerAllower
	Sub  cmd.SubCommand `cmd:"list"`
	Role roles          `cmd:"role"`
}

// Run ...
func (r RoleAdd) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	if len(r.Targets) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	tg, ok := r.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	c := RoleAddOffline{
		Target:   tg.Name(),
		Role:     r.Role,
		Duration: r.Duration,
	}
	c.Run(s, o)
}

// Run ...
func (r RoleRemove) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	if len(r.Targets) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	tg, ok := r.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	c := RoleRemoveOffline{
		Target: tg.Name(),
		Role:   r.Role,
	}
	c.Run(s, o)
}

// Run ...
func (a RoleAddOffline) Run(src cmd.Source, _ *cmd.Output) {
	t, err := data.LoadUserOrCreate(a.Target, "")
	if err != nil {
		internal.Messagef(src, "command.target.unknown")
		return
	}

	r, _ := role.ByName(string(a.Role))
	if p, ok := src.(*player.Player); ok {
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
		if u.Roles.Highest().Tier() < r.Tier() {
			internal.Messagef(src, "command.role.higher")
			return
		}
		if t.Roles.Contains(rls.Operator()) {
			internal.Messagef(src, "command.role.modify.other")
			return
		}
	}

	if t.Roles.Contains(r) {
		internal.Messagef(src, "command.role.has", r.Name())
		return
	}
	t.Roles.Add(r)
	d := "infinity and beyond"
	duration, hasDuration := a.Duration.Load()
	if hasDuration {
		duration, err := timeutil.ParseDuration(duration)
		if err != nil {
			internal.Messagef(src, "command.duration.invalid")
			return
		}
		d = durafmt.ParseShort(duration).String()
		t.Roles.Expire(r, time.Now().Add(duration))
	}
	data.SaveUser(t)

	internal.Alertf(src, "staff.alert.role.add", r.Name(), t.DisplayName, d)
	internal.Messagef(src, "command.role.add", r.Name(), t.DisplayName, d)
}

// Run ...
func (d RoleRemoveOffline) Run(src cmd.Source, _ *cmd.Output) {
	t, err := data.LoadUserFromName(d.Target)
	if err != nil {
		internal.Messagef(src, "command.target.unknown")
		return
	}

	r, _ := role.ByName(string(d.Role))
	if p, ok := src.(*player.Player); ok {
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
		if u.Roles.Highest().Tier() < r.Tier() {
			internal.Messagef(src, "command.role.higher")
			return
		}
		if t.Roles.Contains(rls.Operator()) {
			internal.Messagef(src, "command.role.modify.other")
			return
		}
	}

	if !t.Roles.Contains(r) {
		internal.Messagef(src, "command.role.missing", r.Name())
		return
	}
	t.Roles.Remove(r)
	data.SaveUser(t)

	internal.Alertf(src, "staff.alert.role.remove", r.Name(), t.DisplayName)
	internal.Messagef(src, "command.role.remove", r.Name(), t.DisplayName)
}

// Run ...
func (r RoleList) Run(src cmd.Source, _ *cmd.Output) {
	rl, ok := role.ByName(string(r.Role))
	if !ok {
		return
	}

	users, err := data.LoadUsersFromRole(rl)
	if err != nil {
		panic(err)
	}
	if len(users) <= 0 {
		internal.Messagef(src, "command.role.list.empty")
		return
	}

	var usernames []string
	for _, u := range users {
		usernames = append(usernames, u.DisplayName)
	}
	internal.Messagef(src, "command.role.list", r.Role, len(users), strings.Join(usernames, ", "))
}

type (
	roles string
)

// Type ...
func (roles) Type() string {
	return "role"
}

// Options ...
func (roles) Options(src cmd.Source) (roles []string) {
	p, disallow := src.(*player.Player)
	if disallow {
		u, err := data.LoadUserFromName(p.Name())
		if err == nil {
			disallow = !u.Roles.Contains(rls.Operator())
		}
	}
	for _, r := range role.All() {
		if r == rls.Operator() && disallow {
			continue
		}
		roles = append(roles, r.Name())
	}
	return roles
}

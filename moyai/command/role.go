package command

import (
	"github.com/bedrock-gophers/role/role"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/internal/timeutil"
	"github.com/moyai-network/teams/moyai/data"
	rls "github.com/moyai-network/teams/moyai/roles"
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
func (a RoleAddOffline) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	t, err := data.LoadUserFromName(a.Target)
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	r, _ := role.ByName(string(a.Role))
	if p, ok := s.(*player.Player); ok {
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
		if u.Roles.Highest().Tier() < r.Tier() {
			o.Error(lang.Translatef(l, "command.role.higher"))
			return
		}
		if t.Roles.Contains(rls.Operator()) {
			o.Error(lang.Translatef(l, "command.role.modify.other"))
			return
		}
	}

	if t.Roles.Contains(r) {
		o.Error(lang.Translatef(l, "command.role.has", r.Name()))
		return
	}
	t.Roles.Add(r)
	d := "infinity and beyond"
	duration, hasDuration := a.Duration.Load()
	if hasDuration {
		duration, err := timeutil.ParseDuration(duration)
		if err != nil {
			o.Error(lang.Translatef(l, "command.duration.invalid"))
			return
		}
		d = durafmt.ParseShort(duration).String()
		t.Roles.Expire(r, time.Now().Add(duration))
	}
	data.SaveUser(t)

	//user.Alert(s, "staff.alert.role.add", r.Name(), t.DisplayName, d)
	o.Print(lang.Translatef(l, "command.role.add", r.Name(), t.DisplayName, d))
}

// Run ...
func (d RoleRemoveOffline) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	t, err := data.LoadUserFromName(d.Target)
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	r, _ := role.ByName(string(d.Role))
	if p, ok := s.(*player.Player); ok {
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
		if u.Roles.Highest().Tier() < r.Tier() {
			o.Error(lang.Translatef(l, "command.role.higher"))
			return
		}
		if t.Roles.Contains(rls.Operator()) {
			o.Error(lang.Translatef(l, "command.role.modify.other"))
			return
		}
	}

	if !t.Roles.Contains(r) {
		o.Error(lang.Translatef(l, "command.role.missing", r.Name()))
		return
	}
	t.Roles.Remove(r)
	data.SaveUser(t)

	//user.Alert(s, "staff.alert.role.remove", r.Name(), t.DisplayName)
	o.Print(lang.Translatef(l, "command.role.remove", r.Name(), t.DisplayName))
}

// Run ...
func (r RoleList) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	rl, ok := role.ByName(string(r.Role))
	if !ok {
		return
	}

	users, err := data.LoadUsersFromRole(rl)
	if err != nil {
		panic(err)
	}
	if len(users) <= 0 {
		o.Error(lang.Translatef(l, "command.role.list.empty"))
		return
	}

	var usernames []string
	for _, u := range users {
		usernames = append(usernames, u.DisplayName)
	}
	o.Print(lang.Translatef(l, "command.role.list", r.Role, len(users), strings.Join(usernames, ", ")))
}

type (
	roles string
)

// Type ...
func (roles) Type() string {
	return "role"
}

// Options ...
func (roles) Options(s cmd.Source) (roles []string) {
	p, disallow := s.(*player.Player)
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

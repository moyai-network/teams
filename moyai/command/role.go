package command

import (
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/hako/durafmt"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/data"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/user"
	"go.mongodb.org/mongo-driver/bson"
)

// RoleAdd is a command to add a role to a player.
type RoleAdd struct {
	Sub      cmd.SubCommand       `cmd:"add"`
	Targets  []cmd.Target         `cmd:"target"`
	Role     roles                `cmd:"role"`
	Duration cmd.Optional[string] `cmd:"duration"`
}

// RoleRemove is a command to remove a role from a player.
type RoleRemove struct {
	Sub     cmd.SubCommand `cmd:"remove"`
	Targets []cmd.Target   `cmd:"target"`
	Role    roles          `cmd:"role"`
}

// RoleAddOffline is a command to remove a role from an offline user.
type RoleAddOffline struct {
	Sub      cmd.SubCommand       `cmd:"add"`
	Target   string               `cmd:"target"`
	Role     roles                `cmd:"role"`
	Duration cmd.Optional[string] `cmd:"duration"`
}

// RoleRemoveOffline is a command to remove a role from an offline user.
type RoleRemoveOffline struct {
	Sub    cmd.SubCommand `cmd:"remove"`
	Target string         `cmd:"target"`
	Role   roles          `cmd:"role"`
}

// RoleList is a command to list all users with a role.
type RoleList struct {
	Sub  cmd.SubCommand `cmd:"list"`
	Role roles          `cmd:"role"`
}

// Run ...
func (a RoleAdd) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	if len(a.Targets) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	p, isPlayer := s.(*player.Player)
	if !isPlayer && len(a.Targets) < 1 {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	var t data.User
	var err error

	if len(a.Targets) > 0 {
		tP, ok := a.Targets[0].(*player.Player)
		if !ok {
			o.Error(lang.Translatef(l, "command.target.unknown"))
			return
		}
		t, err = data.LoadUserOrCreate(tP.Name())
		if err != nil {
			o.Error(lang.Translatef(l, "command.target.unknown"))
			return
		}
		if isPlayer && tP != p && t.Roles.Contains(role.Operator{}) {
			o.Error(lang.Translatef(l, "command.role.modify.other"))
			return
		}
	} else {
		t, err = data.LoadUserOrCreate(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
	}

	r, _ := role.ByName(string(a.Role))
	if isPlayer {
		u, err := data.LoadUserOrCreate(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
		if !u.Roles.Contains(role.Operator{}) {
			if strings.EqualFold(p.Handler().(*user.Handler).XUID(), t.XUID) {
				o.Error(lang.Translatef(l, "command.role.modify.self"))
				return
			}
			if role.Tier(u.Roles.Highest()) < role.Tier(r) {
				o.Error(lang.Translatef(l, "command.role.higher"))
				return
			}
		}
	}

	duration, hasDuration := a.Duration.Load()
	if t.Roles.Contains(r) {
		e, ok := t.Roles.Expiration(r)
		if !ok {
			o.Error(lang.Translatef(l, "command.role.has", r.Name()))
			return
		}
		if hasDuration {
			duration, err := moose.ParseDuration(duration)
			if err != nil {
				o.Error(lang.Translatef(l, "command.duration.invalid"))
				return
			}
			if e.After(time.Now().Add(duration)) {
				o.Error(lang.Translatef(l, "command.role.has", r.Name()))
				return
			}
		}
		t.Roles.Remove(r)
	}
	t.Roles.Add(r)
	d := "infinity and beyond"
	if hasDuration {
		duration, err := moose.ParseDuration(duration)
		if err != nil {
			o.Error(lang.Translatef(l, "command.duration.invalid"))
			return
		}
		d = durafmt.ParseShort(duration).String()
		t.Roles.Expire(r, time.Now().Add(duration))
	}
	_ = data.SaveUser(t)

	user.Alert(s, "staff.alert.role.add", r.Name(), t.DisplayName, d)
	o.Print(lang.Translatef(l, "command.role.add", r.Name(), t.DisplayName, d))
}

// Run ...
func (d RoleRemove) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	if len(d.Targets) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	p, isPlayer := s.(*player.Player)
	if !isPlayer && len(d.Targets) < 1 {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	var t data.User
	var err error

	if len(d.Targets) > 0 {
		tP, ok := d.Targets[0].(*player.Player)
		if !ok {
			o.Error(lang.Translatef(l, "command.target.unknown"))
			return
		}
		t, err = data.LoadUserOrCreate(tP.Name())
		if err != nil {
			o.Error(lang.Translatef(l, "command.target.unknown"))
			return
		}
		if isPlayer && tP != p && t.Roles.Contains(role.Operator{}) {
			o.Error(lang.Translatef(l, "command.role.modify.other"))
			return
		}
	} else {
		t, err = data.LoadUserOrCreate(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
	}

	r, _ := role.ByName(string(d.Role))
	if isPlayer {
		u, err := data.LoadUserOrCreate(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
		if !u.Roles.Contains(role.Operator{}) {
			if strings.EqualFold(p.Handler().(*user.Handler).XUID(), t.XUID) {
				o.Error(lang.Translatef(l, "command.role.modify.self"))
				return
			}
			if role.Tier(u.Roles.Highest()) < role.Tier(r) {
				o.Error(lang.Translatef(l, "command.role.higher"))
				return
			}
		}
	}

	if !t.Roles.Contains(r) {
		o.Error(lang.Translatef(l, "command.role.missing", r.Name()))
		return
	}
	t.Roles.Remove(r)

	user.Alert(s, "staff.alert.role.remove", r.Name(), t.DisplayName)
	o.Print(lang.Translatef(l, "command.role.remove", r.Name(), t.DisplayName))
}

// Run ...
func (a RoleAddOffline) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	t, err := data.LoadUserOrCreate(a.Target)
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	r, _ := role.ByName(string(a.Role))
	if p, ok := s.(*player.Player); ok {
		u, err := data.LoadUserOrCreate(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
		if role.Tier(u.Roles.Highest()) < role.Tier(r) {
			o.Error(lang.Translatef(l, "command.role.higher"))
			return
		}
		if t.Roles.Contains(role.Operator{}) {
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
		duration, err := moose.ParseDuration(duration)
		if err != nil {
			o.Error(lang.Translatef(l, "command.duration.invalid"))
			return
		}
		d = durafmt.ParseShort(duration).String()
		t.Roles.Expire(r, time.Now().Add(duration))
	}
	_ = data.SaveUser(t)

	user.Alert(s, "staff.alert.role.add", r.Name(), t.DisplayName, d)
	o.Print(lang.Translatef(l, "command.role.add", r.Name(), t.DisplayName, d))
}

// Run ...
func (d RoleRemoveOffline) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	t, err := data.LoadUserOrCreate(d.Target)
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	r, _ := role.ByName(string(d.Role))
	if p, ok := s.(*player.Player); ok {
		u, err := data.LoadUserOrCreate(p.Name())
		if err != nil {
			// The user somehow left in the middle of this, so just stop in our tracks.
			return
		}
		if role.Tier(u.Roles.Highest()) < role.Tier(r) {
			o.Error(lang.Translatef(l, "command.role.higher"))
			return
		}
		if t.Roles.Contains(role.Operator{}) {
			o.Error(lang.Translatef(l, "command.role.modify.other"))
			return
		}
	}

	if !t.Roles.Contains(r) {
		o.Error(lang.Translatef(l, "command.role.missing", r.Name()))
		return
	}
	t.Roles.Remove(r)
	_ = data.SaveUser(t)

	user.Alert(s, "staff.alert.role.remove", r.Name(), t.DisplayName)
	o.Print(lang.Translatef(l, "command.role.remove", r.Name(), t.DisplayName))
}

// Run ...
func (r RoleList) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	users, err := data.LoadUsersCond(bson.M{"roles": r.Role})
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

// Allow ...
func (RoleAdd) Allow(s cmd.Source) bool {
	return allow(s, true, role.Manager{})
}

// Allow ...
func (RoleRemove) Allow(s cmd.Source) bool {
	return allow(s, true, role.Manager{})
}

// Allow ...
func (RoleAddOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Manager{})
}

// Allow ...
func (RoleRemoveOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Manager{})
}

// Allow ...
func (RoleList) Allow(s cmd.Source) bool {
	return allow(s, true, role.Manager{})
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
		u, err := data.LoadUserOrCreate(p.Name())
		if err == nil {
			disallow = !u.Roles.Contains(role.Operator{})
		}
	}
	for _, r := range role.All() {
		if _, ok := r.(role.Operator); ok && disallow {
			continue
		}
		roles = append(roles, r.Name())
	}
	return roles
}

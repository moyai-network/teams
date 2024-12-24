package command

/*import (
	"strings"
	"time"

	"github.com/internal-network/moose"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/internal-network/moose/data"
	"github.com/internal-network/moose/lang"
	"github.com/internal-network/moose/role"
	"github.com/internal-network/teams/internal/form"
	"github.com/internal-network/teams/internal/user"
	"go.repository.org/repository-driver/bson"
)

// BlacklistForm is a command that is used to blacklist a player through a punishment form.
type BlacklistForm struct{}

// BlacklistList is a command that outputs a list of blacklisted players.
type BlacklistList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

// BlacklistInfoOffline is a command that displays the blacklist information of an offline player.
type BlacklistInfoOffline struct {
	Sub    cmd.SubCommand `cmd:"info"`
	Target string         `cmd:"target"`
}

// BlacklistLiftOffline is a command that is used to lift the blacklist of an offline player.
type BlacklistLiftOffline struct {
	Sub    cmd.SubCommand `cmd:"lift"`
	Target string         `cmd:"target"`
}

// Blacklist is a command that is used to blacklist an online player.
type Blacklist struct {
	Targets []cmd.Target              `cmd:"target"`
	Reason  cmd.Optional[cmd.Varargs] `cmd:"reason"`
}

// BlacklistOffline is a command that is used to blacklist an offline player.
type BlacklistOffline struct {
	Target string                    `cmd:"target"`
	Reason cmd.Optional[cmd.Varargs] `cmd:"reason"`
}

// Run ...
func (b BlacklistList) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(s)
	users, err := data.LoadUsersCond(bson.M{
		"$and": bson.A{
			bson.M{
				"punishments.ban.permanent": true,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		o.Error(lang.Translate(l, "command.blacklist.none"))
		return
	}
	o.Print(lang.Translatef(l, "command.blacklist.list", len(users), strings.Join(names(users, false), ", ")))
}

// Run ...
func (b BlacklistInfoOffline) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(s)
	u, _ := data.LoadUserOrCreate(b.Target)
	if u.Ban.Expired() || !u.Ban.Permanent {
		o.Error(lang.Translate(l, "command.blacklist.not"))
		return
	}
	o.Print(lang.Translatef(l, "punishment.details",
		u.DisplayName,
		u.Ban.Reason,
		"Permanent",
		u.Ban.Staff,
		u.Ban.Occurrence.Format("01/02/2006"),
	))
}

// Run ...
func (b BlacklistLiftOffline) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(src)
	u, err := data.LoadUserOrCreate(b.Target)
	if err != nil {
		o.Error(lang.Translate(l, "command.target.unknown"))
		return
	}
	if u.Ban.Expired() {
		o.Error(lang.Translate(l, "command.blacklist.not"))
		return
	}
	u.Ban = moose.Punishment{}
	data.SaveUser(u)

	user.Alert(src, "staff.alert.unblacklist", u.DisplayName)
	//webhook.SendPunishment(s.Name(), u.DisplayName(), "", "Unblacklist")
	o.Print(lang.Translatef(l, "command.blacklist.lift", u.DisplayName))
}

// Run ...
func (BlacklistForm) Run(s cmd.Source, _ *cmd.Output, tx *world.Tx) {
	p := s.(*player.Player)
	p.SendForm(form.NewBlacklist())
}

// Run ...
func (b Blacklist) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(src)
	s := src.(cmd.NamedTarget)
	if len(b.Targets) > 1 {
		o.Error(lang.Translate(l, "command.targets.exceed"))
		return
	}
	t, ok := b.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translate(l, "command.target.unknown"))
		return
	}
	if t == s {
		o.Error(lang.Translate(l, "command.blacklist.self"))
		return
	}
	u, err := data.LoadUserOrCreate(t.Name())
	if err != nil {
		o.Error(lang.Translate(l, "command.target.unknown"))
		return
	}
	if u.Roles.Contains(role.Operator{}) {
		o.Error(lang.Translate(l, "command.blacklist.operator"))
		return
	}

	reason := strings.TrimSpace(string(b.Reason.LoadOr("")))
	if len(reason) == 0 {
		reason = "None"
	}
	u.Ban = moose.Punishment{
		Staff:      s.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Permanent:  true,
	}
	data.SaveUser(u)

	t.Disconnect(strings.Join([]string{
		lang.Translate(l, "user.blacklist.header"),
		lang.Translatef(l, "user.blacklist.description", reason),
	}, "\n"))

	user.Alert(src, "staff.alert.blacklist", t.Name())
	user.Broadcast("command.blacklist.broadcast", s.Name(), t.Name())
	//webhook.SendPunishment(s.Name(), t.Name(), reason, "Blacklist")
	o.Print(lang.Translatef(l, "command.blacklist.success", t.Name(), reason))
}

// Run ...
func (b BlacklistOffline) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(src)
	s := src.(cmd.NamedTarget)
	if s.Name() == b.Target {
		o.Error(lang.Translate(l, "command.blacklist.self"))
		return
	}
	u, err := data.LoadUserOrCreate(b.Target)
	if err != nil {
		o.Error(lang.Translate(l, "command.target.unknown"))
		return
	}
	if u.Roles.Contains(role.Operator{}) {
		o.Error(lang.Translate(l, "command.blacklist.operator"))
		return
	}
	if !u.Ban.Expired() && u.Ban.Permanent {
		o.Error(lang.Translate(l, "command.blacklist.already"))
		return
	}
	reason := strings.TrimSpace(string(b.Reason.LoadOr("")))
	if len(reason) == 0 {
		reason = "None"
	}
	u.Ban = moose.Punishment{
		Staff:      s.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Permanent:  true,
	}
	data.SaveUser(u)

	user.Alert(src, "staff.alert.blacklist", u.DisplayName)
	user.Broadcast("command.blacklist.broadcast", s.Name(), u.DisplayName)
	//webhook.SendPunishment(s.Name(), u.DisplayName(), reason, "Blacklist")
	o.Print(lang.Translatef(l, "command.blacklist.success", u.DisplayName, reason))
}

// Allow ...
func (BlacklistList) Allow(s cmd.Source) bool {
	return Allow(s, true, role.Manager{})
}

// Allow ...
func (BlacklistInfoOffline) Allow(s cmd.Source) bool {
	return Allow(s, true, role.Manager{})
}

// Allow ...
func (BlacklistForm) Allow(s cmd.Source) bool {
	return Allow(s, false, role.Manager{})
}

// Allow ...
func (Blacklist) Allow(s cmd.Source) bool {
	return Allow(s, true, role.Manager{})
}

// Allow ...
func (BlacklistOffline) Allow(s cmd.Source) bool {
	return Allow(s, true, role.Manager{})
}

// Allow ...
func (BlacklistLiftOffline) Allow(s cmd.Source) bool {
	return Allow(s, true, role.Manager{})
}
*/

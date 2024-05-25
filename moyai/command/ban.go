package command

import (
	"strings"
	"time"

	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/internal/punishment"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/hako/durafmt"
	"github.com/moyai-network/teams/moyai/user"
)

// BanForm is a command that is used to ban a player through a punishment form.
type BanForm struct{}

// BanList is a command that outputs a list of banned players.
type BanList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

// BanInfoOffline is a command that displays the ban information of an offline player.
type BanInfoOffline struct {
	Sub    cmd.SubCommand `cmd:"info"`
	Target string         `cmd:"target"`
}

// BanLiftOffline is a command that is used to lift the ban of an offline player.
type BanLiftOffline struct {
	Sub    cmd.SubCommand `cmd:"lift"`
	Target string         `cmd:"target"`
}

// Ban is a command that is used to ban an online player.
type Ban struct {
	Targets []cmd.Target `cmd:"target"`
	Reason  banReason    `cmd:"reason"`
}

// BanOffline is a command that is used to ban an offline player.
type BanOffline struct {
	Target string    `cmd:"target"`
	Reason banReason `cmd:"reason"`
}

// Run ...
func (BanList) Run(s cmd.Source, o *cmd.Output) {
	// l := locale(s)
	// users, err := data.LoadUsersCond(
	// 	bson.M{
	// 		"$and": bson.A{
	// 			bson.M{
	// 				"punishments.ban.expiration": bson.M{"$ne": time.Time{}},
	// 			}, bson.M{
	// 				"punishments.ban.expiration": bson.M{"$gt": time.Now()},
	// 			},
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// if len(users) == 0 {
	// 	o.Error(lang.Translatef(l, "command.ban.none"))
	// 	return
	// }
	// o.Print(lang.Translatef(l, "command.ban.list", len(users), strings.Join(names(users, false), ", ")))
	// }
	// o.Print(lang.Translatef(l, "punishment.details",
	// 	u.DisplayName,
	// 	u.Ban.Reason,
	// 	durafmt.ParseShort(u.Ban.Remaining()),
	// 	u.Ban.Staff,
	// 	u.Ban.Occurrence.Format("01/02/2006"),
	// ))
}

// Run ...
func (b BanLiftOffline) Run(src cmd.Source, o *cmd.Output) {
	l := locale(src)
	u, err := data.LoadUserFromName(b.Target)
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Teams.Ban.Expired() || u.Teams.Ban.Permanent {
		o.Error(lang.Translatef(l, "command.ban.not"))
		return
	}
	u.Teams.Ban = punishment.Punishment{}
	data.SaveUser(u)

	user.Alertf(src, "staff.alert.unban", u.DisplayName)
	//webhook.SendPunishment(s.Name(), u.DisplayName(), "", "Unban")
	o.Print(lang.Translatef(l, "command.ban.lift", u.DisplayName))
}

// Run ...
func (BanForm) Run(s cmd.Source, _ *cmd.Output) {
	//p := s.(*player.Player)
	//p.SendForm(form.NewBan())
}

// Run ...
func (b Ban) Run(src cmd.Source, o *cmd.Output) {
	l := locale(src)
	s := src.(cmd.NamedTarget)
	if len(b.Targets) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	t, ok := b.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if t == src {
		o.Error(lang.Translatef(l, "command.ban.self"))
		return
	}
	u, err := data.LoadUserFromName(t.Name())
	if err != nil {
		return
	}
	if u.Roles.Contains(role.Operator{}) {
		o.Error(lang.Translatef(l, "command.ban.operator"))
		return
	}
	reason, length := parseBanReason(b.Reason)
	u.Teams.Ban = punishment.Punishment{
		Staff:      s.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Expiration: time.Now().Add(length),
	}
	data.SaveUser(u)

	t.Disconnect(strings.Join([]string{
		lang.Translatef(l, "user.ban.header"),
		lang.Translatef(l, "user.ban.description", reason, durafmt.ParseShort(length)),
	}, "\n"))

	user.Alertf(src, "staff.alert.ban", t.Name(), reason)
	user.Broadcastf("command.ban.broadcast", s.Name(), t.Name(), reason)
	//webhook.SendPunishment(s.Name(), t.Name(), reason, "Ban")
	o.Print(lang.Translatef(l, "command.ban.success", t.Name(), reason))
}

// Run ...
func (b BanOffline) Run(src cmd.Source, o *cmd.Output) {
	l := locale(src)
	s := src.(cmd.NamedTarget)
	if strings.EqualFold(s.Name(), b.Target) {
		o.Error(lang.Translatef(l, "command.ban.self"))
		return
	}
	u, err := data.LoadUserFromName(b.Target)
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Roles.Contains(role.Operator{}) {
		o.Error(lang.Translatef(l, "command.ban.operator"))
		return
	}
	if !u.Teams.Ban.Expired() {
		o.Error(lang.Translatef(l, "command.ban.already"))
		return
	}

	reason, length := parseBanReason(b.Reason)
	u.Teams.Ban = punishment.Punishment{
		Staff:      s.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Expiration: time.Now().Add(length),
	}
	data.SaveUser(u)

	user.Alertf(src, "staff.alert.ban", u.DisplayName, reason)
	user.Broadcastf("command.ban.broadcast", s.Name(), u.DisplayName, reason)
	//webhook.SendPunishment(s.Name(), u.DisplayName(), reason, "Ban")
	o.Print(lang.Translatef(l, "command.ban.success", u.DisplayName, reason))
}

// Allow ...
func (BanList) Allow(s cmd.Source) bool {
	return allow(s, true, role.Mod{})
}

// Allow ...
func (BanInfoOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Mod{})
}

// Allow ...
func (BanForm) Allow(s cmd.Source) bool {
	return allow(s, false, role.Mod{})
}

// Allow ...
func (Ban) Allow(s cmd.Source) bool {
	return allow(s, true, role.Mod{})
}

// Allow ...
func (BanOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Mod{})
}

// Allow ...
func (BanLiftOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Mod{})
}

type banReason string

// Type ...
func (banReason) Type() string {
	return "banReason"
}

// Options ...
func (banReason) Options(cmd.Source) []string {
	return []string{
		"advantage",
		"allying",
		"glitching",
		"hostage",
		"exploitation",
		"abuse",
		"skin",
		"advertisement",
		"evasion",
	}
}

// parseBanReason returns the formatted BanReason and ban duration.
func parseBanReason(r banReason) (string, time.Duration) {
	switch r {
	case "advantage":
		return "Unfair Advantage", time.Hour * 24 * 30
	case "exploitation":
		return "Exploitation", time.Hour * 24 * 9
	case "abuse":
		return "Permission Abuse", time.Hour * 24 * 30
	case "skin":
		return "Invalid Skin", time.Hour * 24 * 3
	case "evasion":
		return "Evasion", time.Hour * 24 * 120
	case "advertisement":
		return "Advertisement", time.Hour * 24 * 6
	case "glitching":
		return "Glitching", time.Hour
	}
	return "Unknown", time.Hour * 24
}

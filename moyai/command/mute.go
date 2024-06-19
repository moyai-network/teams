package command

import (
	"github.com/moyai-network/teams/moyai"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/hako/durafmt"

	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/internal/punishment"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
)

// MuteForm is a command that is used to mute an online player through a punishment form.
type MuteForm struct{ trialAllower }

// MuteList is a command that outputs a list of muted players.
type MuteList struct {
	trialAllower
	Sub cmd.SubCommand `cmd:"list"`
}

// MuteInfo is a command that displays the mute information of an online player.
type MuteInfo struct {
	trialAllower
	Sub     cmd.SubCommand `cmd:"info"`
	Targets []cmd.Target   `cmd:"target"`
}

// MuteInfoOffline is a command that displays the mute information of an offline player.
type MuteInfoOffline struct {
	trialAllower
	Sub    cmd.SubCommand `cmd:"info"`
	Target string         `cmd:"target"`
}

// MuteLift is a command that is used to lift the mute of an online player.
type MuteLift struct {
	modAllower
	Sub     cmd.SubCommand `cmd:"lift"`
	Targets []cmd.Target   `cmd:"target"`
}

// MuteLiftOffline is a command that is used to lift the mute of an offline player.
type MuteLiftOffline struct {
	modAllower
	Sub    cmd.SubCommand `cmd:"lift"`
	Target string         `cmd:"target"`
}

// Mute is a command that is used to mute an online player.
type Mute struct {
	trialAllower
	Targets []cmd.Target `cmd:"target"`
	Reason  muteReason   `cmd:"reason"`
}

// MuteOffline is a command that is used to mute an offline player.
type MuteOffline struct {
	trialAllower
	Target string     `cmd:"target"`
	Reason muteReason `cmd:"reason"`
}

// Run ...
func (MuteList) Run(s cmd.Source, o *cmd.Output) {
	// l := locale(s)
	// users, err := data.LoadUsersCond(
	// 	bson.M{
	// 		"$and": bson.A{
	// 			bson.M{
	// 				"punishments.mute.expiration": bson.M{"$ne": time.Time{}},
	// 			}, bson.M{
	// 				"punishments.mute.expiration": bson.M{"$gt": time.Now()},
	// 			},
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// if len(users) == 0 {
	// 	o.Error(lang.Translatef(l, "command.mute.none"))
	// 	return
	// }
	// o.Print(lang.Translatef(l, "command.mute.list", len(users), strings.Join(names(users, true), ", ")))
}

// Run ...
func (m MuteInfo) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := m.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Teams.Mute.Expired() {
		o.Error(lang.Translatef(l, "command.mute.not"))
		return
	}
	mute := u.Teams.Mute
	o.Print(lang.Translatef(l, "punishment.details",
		p.Name(),
		mute.Reason,
		durafmt.Parse(mute.Remaining()),
		mute.Staff,
		mute.Occurrence.Format("01/02/2006"),
	))
}

// Run ...
func (m MuteInfoOffline) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	u, err := data.LoadUserFromName(m.Target)
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Teams.Mute.Expired() {
		o.Error(lang.Translatef(l, "command.mute.not"))
		return
	}
	o.Print(lang.Translatef(l, "punishment.details",
		u.DisplayName,
		u.Teams.Mute.Reason,
		durafmt.Parse(u.Teams.Mute.Remaining()),
		u.Teams.Mute.Staff,
		u.Teams.Mute.Occurrence.Format("01/02/2006"),
	))
}

// Run ...
func (m MuteLift) Run(src cmd.Source, out *cmd.Output) {
	l := locale(src)
	p, ok := m.Targets[0].(*player.Player)
	if !ok {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Teams.Mute.Expired() {
		out.Error(lang.Translatef(l, "command.mute.not"))
		return
	}
	u.Teams.Mute = punishment.Punishment{}
	data.SaveUser(u)

	moyai.Alertf(src, "staff.alert.unmute", p.Name())
	//webhook.SendPunishment(s.Name(), u.DisplayName(), "", "Unmute")
	out.Print(lang.Translatef(l, "command.mute.lift", p.Name()))
}

// Run ...
func (m MuteLiftOffline) Run(src cmd.Source, out *cmd.Output) {
	l := locale(src)
	u, err := data.LoadUserFromName(m.Target)
	if err != nil {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Teams.Mute.Expired() {
		out.Error(lang.Translatef(l, "command.mute.not"))
		return
	}
	u.Teams.Mute = punishment.Punishment{}
	data.SaveUser(u)

	moyai.Alertf(src, "staff.alert.unmute", u.DisplayName)
	//webhook.SendPunishment(src.Name(), u.DisplayName(), "", "Unmute")
	out.Print(lang.Translatef(l, "command.mute.lift", u.DisplayName))
}

// Run ...
func (m MuteForm) Run(s cmd.Source, _ *cmd.Output) {
	// p := s.(*player.Player)
	// p.SendForm(form.NewMute(p))
}

// Run ...
func (m Mute) Run(src cmd.Source, out *cmd.Output) {
	l := locale(src)
	if len(m.Targets) > 1 {
		out.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	t, ok := m.Targets[0].(*player.Player)
	if !ok {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if t == src {
		out.Error(lang.Translatef(l, "command.mute.self"))
		return
	}
	u, err := data.LoadUserFromName(t.Name())
	if err != nil {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Roles.Contains(role.Operator{}) {
		out.Error(lang.Translatef(l, "command.mute.operator"))
		return
	}
	if !u.Teams.Mute.Expired() {
		out.Error(lang.Translatef(l, "command.mute.already"))
		return
	}
	sn := src.(cmd.NamedTarget)
	reason, length := parseMuteReason(m.Reason)
	u.Teams.Mute = punishment.Punishment{
		Staff:      sn.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Expiration: time.Now().Add(length),
	}
	data.SaveUser(u)

	moyai.Alertf(src, "staff.alert.mute", t.Name(), reason)
	//webhook.SendPunishment(src.Name(), t.Name(), reason, "Mute")
	out.Print(lang.Translatef(l, "command.mute.success", t.Name(), reason))
}

// Run ...
func (m MuteOffline) Run(src cmd.Source, out *cmd.Output) {
	l := locale(src)
	sn := src.(cmd.NamedTarget)

	u, err := data.LoadUserFromName(m.Target)
	if err != nil {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	if strings.EqualFold(u.Name, m.Target) {
		out.Error(lang.Translatef(l, "command.mute.self"))
		return
	}

	if u.Roles.Contains(role.Operator{}) {
		out.Error(lang.Translatef(l, "command.mute.operator"))
		return
	}
	if !u.Teams.Mute.Expired() {
		out.Error(lang.Translatef(l, "command.mute.already"))
		return
	}

	reason, length := parseMuteReason(m.Reason)
	u.Teams.Mute = punishment.Punishment{
		Staff:      sn.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Expiration: time.Now().Add(length),
	}
	data.SaveUser(u)

	moyai.Alertf(src, "staff.alert.mute", u.DisplayName, reason)
	//webhook.SendPunishment(s.Name(), u.DisplayName(), reason, "Mute")
	out.Print(lang.Translatef(l, "command.mute.success", u.DisplayName, reason))
}

type (
	muteReason string
)

// Type ...
func (muteReason) Type() string {
	return "muteReason"
}

// Options ...
func (muteReason) Options(cmd.Source) []string {
	return []string{
		"spam",
		"toxic",
		"advertising",
		"threats",
	}
}

// parseMuteReason returns the formatted muteReason and mute duration.
func parseMuteReason(r muteReason) (string, time.Duration) {
	switch r {
	case "spam":
		return "Spam", time.Hour * 6
	case "toxic":
		return "Toxicity", time.Hour * 9
	case "advertising":
		return "Advertising", time.Hour * 24 * 3
	case "threats":
		return "Threats", time.Hour * 24 * 5
	}
	panic("should never happen")
}

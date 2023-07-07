package command

import (
	"regexp"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
	"golang.org/x/exp/slices"
)

var regex = regexp.MustCompile("^[a-zA-Z0-9]*$")

// TeamCreate is the command used to create teams.
type TeamCreate struct {
	Sub  cmd.SubCommand `cmd:"create"`
	Name string         `cmd:"name"`
}

// TeamInvite is the command used to invite players to teams.
type TeamInvite struct {
	Sub    cmd.SubCommand `cmd:"invite"`
	Target []cmd.Target   `cmd:"target"`
}

// TeamJoin is the command used to join teams.
type TeamJoin struct {
	Sub  cmd.SubCommand `cmd:"join"`
	Team teamInvitation `cmd:"team"`
}

// Run ...
func (t TeamCreate) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, _ := data.LoadUser(p.Name(), p.Handler().(*user.Handler).XUID())
	if _, ok = u.Team(); ok {
		out.Error(lang.Translatef(p.Locale(), "team.create.already"))
		return
	}
	t.Name = strings.TrimSpace(t.Name)

	if !regex.MatchString(t.Name) {
		out.Error(lang.Translatef(p.Locale(), "team.create.name.invalid"))
		return
	}
	if len(t.Name) < 3 {
		out.Error(lang.Translatef(p.Locale(), "team.create.name.short"))
		return
	}
	if len(t.Name) > 15 {
		out.Error(lang.Translatef(p.Locale(), "team.create.name.long"))
		return
	}

	if _, ok := data.LoadTeam(t.Name); ok {
		out.Error(lang.Translatef(p.Locale(), "team.create.exists"))
		return
	}
	tm := data.DefaultTeam(t.Name).WithMembers(data.DefaultMember(p.Handler().(*user.Handler).XUID(), p.Name()).WithRank(3))
	data.SaveTeam(tm)

	out.Print(lang.Translatef(p.Locale(), "team.create.success", tm.DisplayName))
	user.Broadcast("team.create.success.broadcast", p.Name(), tm.DisplayName)
}

// Run ...
func (t TeamInvite) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	target, ok := t.Target[0].(*player.Player)
	if !ok {
		return
	}
	if target == p {
		out.Error(lang.Translatef(p.Locale(), "team.invite.self"))
		return
	}
	u, _ := data.LoadUser(p.Name(), p.Handler().(*user.Handler).XUID())
	tm, ok := u.Team()
	if !ok {
		out.Error(lang.Translatef(p.Locale(), "user.team-less"))
		return
	}

	if slices.ContainsFunc(tm.Members, func(member data.Member) bool {
		return strings.EqualFold(target.Name(), member.Name)
	}) {
		out.Error(lang.Translatef(p.Locale(), "team.invite.member", target.Name()))
		return
	}

	tg, _ := data.LoadUser(target.Name(), target.XUID())

	if tg.Invitations.Active(tm.Name) {
		out.Error(lang.Translatef(p.Locale(), "team.invite.already", target.Name()))
		return
	}
	tg.Invitations.Set(tm.Name, time.Minute*5)

	_ = data.SaveUser(tg)

	user.BroadcastTeam(tm, "team.invite.success.broadcast", target.Name())
	target.Message(lang.Translatef(target.Locale(), "team.invite.target", tm.DisplayName))
}

// Run ...
func (t TeamJoin) Run(src cmd.Source, out *cmd.Output) {
	p := src.(*player.Player)
	l := locale(src)

	u, _ := data.LoadUser(p.Name(), p.Handler().(*user.Handler).XUID())
	if _, ok := u.Team(); ok {
		out.Error(lang.Translatef(l, "team.join.error"))
		return
	}

	tm, ok := data.LoadTeam(string(t.Team))
	if !ok {
		// TODO: error message
		return
	}
	tm = tm.WithMembers(append(tm.Members, data.DefaultMember(p.Handler().(*user.Handler).XUID(), p.Name()))...)

	data.SaveTeam(tm)

	out.Print(lang.Translatef(l, "team.join.target", tm.DisplayName))
	user.BroadcastTeam(tm, "team.join.broadcast", p.Name())
}

// Allow ...
func (TeamCreate) Allow(src cmd.Source) bool {
	return allow(src, false)
}

// Allow ...
func (TeamInvite) Allow(src cmd.Source) bool {
	return allow(src, false)
}

// Allow ...
func (TeamJoin) Allow(src cmd.Source) bool {
	return allow(src, false)
}

type (
	// teamInvitation represents the type used as command arguments, listing all the team invitations the user has.
	teamInvitation string
)

// Type ...
func (teamInvitation) Type() string {
	return "team_invitation"
}

// Options ...
func (teamInvitation) Options(src cmd.Source) (options []string) {
	p := src.(*player.Player)
	u, err := data.LoadUser(p.Name(), "")
	if err != nil {
		return
	}
	for t, i := range u.Invitations {
		if i.Active() {
			options = append(options, t)
		}
	}
	return
}

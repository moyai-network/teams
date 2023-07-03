package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
	"time"
)

var regex = regexp.MustCompile("^[a-zA-Z0-9]*$")

type TeamCreate struct {
	Sub  cmd.SubCommand `cmd:"create"`
	Name string         `cmd:"name"`
}

type TeamInvite struct {
	Sub    cmd.SubCommand `cmd:"invite"`
	Target []cmd.Target   `cmd:"target"`
}

func (t TeamCreate) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	if _, ok = data.LoadUserTeam(p.Name()); ok {
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

	if data.TeamExists(t.Name) {
		out.Error(lang.Translatef(p.Locale(), "team.create.exists"))
		return
	}
	tm := data.DefaultTeam(t.Name).WithMembers(data.DefaultMember(p.XUID(), p.Name()).WithRank(3))
	_ = data.SaveTeam(tm)

	out.Print(lang.Translatef(p.Locale(), "team.create.success", tm.DisplayName))
	user.Broadcast("team.create.success.broadcast", p.Name(), tm.DisplayName)
}

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
	tm, ok := data.LoadUserTeam(p.Name())
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

	u, _ := data.LoadUser(target.Name(), target.XUID())

	if u.Invitations.Active(tm.Name) {
		out.Error(lang.Translatef(p.Locale(), "team.invite.already", target.Name()))
		return
	}
	u.Invitations.Set(tm.Name, time.Minute*5)

	_ = data.SaveUser(u)

	for _, m := range tm.Members {
		pl, ok := user.Lookup(m.XUID)
		if ok {
			pl.Message(lang.Translatef(pl.Locale(), "team.invite.success.broadcast", target.Name()))
		}
	}
	target.Message(lang.Translatef(target.Locale(), "team.invite.target", tm.DisplayName))
}

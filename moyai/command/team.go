package command

import (
	"regexp"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server"
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

// TeamInformation is a command used to get information about a team.
type TeamInformation struct {
	Sub  cmd.SubCommand         `cmd:"info"`
	Name cmd.Optional[teamName] `optional:"" cmd:"name"`
	srv  *server.Server
}

func NewTeamInformation(srv *server.Server) TeamInformation {
	return TeamInformation{srv: srv}
}

// TeamWho is a command used to get information about a team.
type TeamWho struct {
	Sub  cmd.SubCommand           `cmd:"who"`
	Name cmd.Optional[teamMember] `optional:"" cmd:"name"`
	srv  *server.Server
}

func NewTeamWho(srv *server.Server) TeamWho {
	return TeamWho{srv: srv}
}

// TeamLeave is a command used to leave a team.
type TeamLeave struct {
	Sub cmd.SubCommand `cmd:"leave"`
}

// TeamKick is a command used to kick a player from a team.
type TeamKick struct {
	Sub    cmd.SubCommand `cmd:"kick"`
	Member member         `cmd:"member"`
}

// TeamPromote is a command used to promote a player in a team.
type TeamPromote struct {
	Sub    cmd.SubCommand `cmd:"promote"`
	Member member         `cmd:"member"`
}

// TeamDemote is a command used to demote a player in a team.
type TeamDemote struct {
	Sub    cmd.SubCommand `cmd:"demote"`
	Member member         `cmd:"member"`
}

// TeamTop is a command used to get the top teams.
type TeamTop struct {
	Sub cmd.SubCommand `cmd:"top"`
}

// TeamClaim is a command used to claim land for a team.
type TeamClaim struct {
	Sub cmd.SubCommand `cmd:"claim"`
}

// TeamUnClaim is a command used to unclaim land for a team.
type TeamUnClaim struct {
	Sub cmd.SubCommand `cmd:"unclaim"`
}

// TeamSetHome is a command used to set a team's home.
type TeamSetHome struct {
	Sub cmd.SubCommand `cmd:"sethome"`
}

// TeamHome is a command used to teleport to a team's home.
type TeamHome struct {
	Sub cmd.SubCommand `cmd:"home"`
}

// TeamList is a command used to list teams.
type TeamList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

// TeamFocusTeam is a command used to focus on a team.
type TeamFocusTeam struct {
	Sub  cmd.SubCommand `cmd:"focus"`
	Name teamName       `cmd:"name"`
}

// TeamFocusPlayer is a command used to focus on a player.
type TeamFocusPlayer struct {
	Sub    cmd.SubCommand `cmd:"focus"`
	Target []cmd.Target   `cmd:"target"`
}

// TeamUnFocus is a command used to unfocus on a team.
type TeamUnFocus struct {
	Sub cmd.SubCommand `cmd:"unfocus"`
}

// TeamChat is a command used to chat in a team.
type TeamChat struct {
	Sub cmd.SubCommand `cmd:"chat"`
}

// TeamWithdraw is a command used to withdraw money from a team.
type TeamWithdraw struct {
	Sub     cmd.SubCommand `cmd:"withdraw"`
	Balance float64        `cmd:"balance"`
}

// TeamDeposit is a command used to deposit money into a team.
type TeamDeposit struct {
	Sub     cmd.SubCommand `cmd:"deposit"`
	Balance float64        `cmd:"balance"`
}

// TeamWithdrawAll is a command used to withdraw all the money from a team.
type TeamWithdrawAll struct {
	Sub cmd.SubCommand `cmd:"withdraw"`
	All cmd.SubCommand `cmd:"all"`
}

// TeamDepositAll is a command used to deposit all of a user's money into a team.
type TeamDepositAll struct {
	Sub cmd.SubCommand `cmd:"deposit"`
	All cmd.SubCommand `cmd:"all"`
}

// TeamStuck is a command to teleport to a safe position
type TeamStuck struct {
	Sub cmd.SubCommand `cmd:"stuck"`
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

	if tm.Frozen() {
		out.Error(lang.Translatef(p.Locale(), "command.team.dtr"))
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
	tm = tm.WithDTR(tm.DTR + 1)

	data.SaveTeam(tm)

	out.Print(lang.Translatef(l, "team.join.target", tm.DisplayName))
	user.BroadcastTeam(tm, "team.join.broadcast", p.Name())
}

// Run ...
func (t TeamInformation) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	n, _ := t.Name.Load()
	name := string(n)
	u, err := data.LoadUser(sourceName, "")
	if err != nil {
		return
	}
	if strings.TrimSpace(name) == "" {
		tm, ok := u.Team()
		if !ok {
			o.Error(lang.Translatef(l, "user.team-less"))
			return
		}
		o.Print(tm.Information(t.srv))
		return
	}
	var anyFound bool

	tm, ok := data.LoadTeam(name)
	if ok {
		o.Print(tm.Information(t.srv))
		anyFound = true
	}
	tm, ok = u.Team()
	if ok {
		o.Print(tm.Information(t.srv))
		anyFound = true
	}
	if !anyFound {
		o.Error(lang.Translatef(l, "command.team.info.not.found", name))
		return
	}
}

// Run ...
func (t TeamWho) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	n, _ := t.Name.Load()
	name := string(n)
	u, err := data.LoadUser(sourceName, "")
	if err != nil {
		return
	}
	if strings.TrimSpace(name) == "" {
		tm, ok := u.Team()
		if !ok {
			o.Error(lang.Translatef(l, "user.team-less"))
			return
		}
		o.Print(tm.Information(t.srv))
		return
	}
	var anyFound bool

	tm, ok := data.LoadTeam(name)
	if ok {
		o.Print(tm.Information(t.srv))
		anyFound = true
	}
	tm, ok = u.Team()
	if ok {
		o.Print(tm.Information(t.srv))
		anyFound = true
	}
	if !anyFound {
		o.Error(lang.Translatef(l, "command.team.info.not.found", name))
		return
	}
}

// Run ...
func (t TeamLeave) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	src := s.(cmd.NamedTarget).Name()
	u, err := data.LoadUser(src, "")
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}
	if !tm.Leader(u.Name) {
		o.Error(lang.Translatef(l, "command.team.leave.leader"))
		return
	}

	if tm.Frozen() {
		o.Error(lang.Translatef(l, "command.team..dtr"))
		return
	}

	players := tm.Members
	p := s.(*player.Player)
	if !ok {
		return
	}
	tm = tm.WithoutMember(data.DefaultMember(p.Handler().(*user.Handler).XUID(), p.Name()))
	for _, m := range tm.Members {
		if _, ok := user.Lookup(m.Name); ok {
			//mem.UpdateState()
		}
	}
	data.SaveTeam(tm)
	for _, m := range players {
		if u, ok := user.Lookup(m.Name); ok {
			u.Player().Message(lang.Translatef(l, "command.team.leave.user.left", p.Name()))
		}
	}
}

// Run ...
func (t TeamKick) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	u, err := data.LoadUser(sourceName, "")
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}
	if !tm.Leader(sourceName) && !tm.Captain(sourceName) {
		o.Error(lang.Translatef(l, "command.team.kick.missing.permission"))
		return
	}
	if string(t.Member) == sourceName {
		o.Error(lang.Translatef(l, "command.team.kick.self"))
		return
	}
	if tm.Leader(string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.kick.leader"))
		return
	}
	if tm.Captain(sourceName) && tm.Captain(string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.kick.captain"))
		return
	}
	if tm.Member(string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.kick.not.found", string(t.Member)))
		return
	}

	if tm.Frozen() {
		o.Error(lang.Translatef(l, "command.team..dtr"))
		return
	}

	us, ok := user.Lookup(string(t.Member))
	if ok {
		tm = tm.WithoutMember(data.DefaultMember(us.XUID(), us.Player().Name()))
		us.Message(lang.Translatef(l, "command.team.kick.user.kicked"))
	}
	for _, m := range tm.Members {
		if _, ok := user.Lookup(m.Name); ok {
			//mem.UpdateState()
		}
	}
	for _, m := range tm.Members {
		if u, ok := user.Lookup(m.Name); ok {
			u.Player().Message(lang.Translatef(l, "command.team.kick.user.kicked", string(t.Member)))
		}
	}

	data.SaveTeam(tm)
	if err != nil {
		o.Error(lang.Translatef(l, "team.save.fail"))
		return
	}
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

// Allow ...
func (TeamInformation) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamWho) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamLeave) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamKick) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamPromote) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamDemote) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamTop) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamClaim) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamUnClaim) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamSetHome) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamHome) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamList) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamFocusTeam) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamUnFocus) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamFocusPlayer) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamChat) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamWithdraw) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamDeposit) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamWithdrawAll) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (TeamDepositAll) Allow(s cmd.Source) bool {
	return allow(s, false)
}

type (
	teamInvitation string
	member         string
	teamName       string
	teamMember     string
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

// Type ...
func (member) Type() string {
	return "member"
}

// Options ...
func (member) Options(src cmd.Source) []string {
	p := src.(*player.Player)
	u, err := data.LoadUser(p.Name(), "")
	if err != nil {
		return nil
	}

	var members []string
	if t, ok := u.Team(); ok {
		for _, m := range t.Members {
			if !strings.EqualFold(m.Name, p.Name()) {
				members = append(members, m.DisplayName)
			}
		}
	}

	return members
}

// Type ...
func (teamName) Type() string {
	return "team_name"
}

// Options ...
func (teamName) Options(cmd.Source) []string {
	var teams []string
	for _, tm := range data.Teams() {
		teams = append(teams, tm.DisplayName)
	}
	return teams
}

// Type ...
func (teamMember) Type() string {
	return "team_member"
}

// Options ...
func (teamMember) Options(src cmd.Source) []string {
	var members []string
	for _, tm := range data.Teams() {
		for _, m := range tm.Members {
			members = append(members, m.DisplayName)
		}
	}

	return members
}

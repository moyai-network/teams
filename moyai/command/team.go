package command

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var regex = regexp.MustCompile("^[a-zA-Z0-9]*$")

// TeamCreate is the command used to create teams.
type TeamCreate struct {
	Sub  cmd.SubCommand `cmd:"create"`
	Name string         `cmd:"name"`
}

// TeamDisband is a command used to disband a team.
type TeamDisband struct {
	Sub cmd.SubCommand `cmd:"disband"`
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

	u, _ := data.LoadUser(p.Name())
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
	u, _ := data.LoadUser(p.Name())
	tm, ok := u.Team()
	if !ok {
		out.Error(lang.Translatef(p.Locale(), "user.team-less"))
		return
	}

	if tm.Frozen() {
		out.Error(lang.Translatef(p.Locale(), "command.team.dtr"))
		return
	}

	if tm.Member(target.Name()) {
		out.Error(lang.Translatef(p.Locale(), "team.invite.member", target.Name()))
		return
	}

	tg, _ := data.LoadUser(target.Name())

	if _, ok := tg.Team(); ok {
		out.Error(lang.Translatef(p.Locale(), "team.invite.has-team", target.Name()))
		return
	}

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

	u, _ := data.LoadUser(p.Name())
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
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	n, _ := t.Name.Load()
	name := string(n)
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	if strings.TrimSpace(name) == "" {
		tm, ok := u.Team()
		if !ok {
			o.Error(lang.Translatef(l, "user.team-less"))
			return
		}
		o.Print(teamInformationFormat(tm, t.srv))
		return
	}
	var anyFound bool

	tm, ok := data.LoadTeam(name)
	if ok {
		o.Print(teamInformationFormat(tm, t.srv))
		anyFound = true
	}
	tm, ok = u.Team()
	if ok {
		o.Print(teamInformationFormat(tm, t.srv))
		anyFound = true
	}
	if !anyFound {
		o.Error(lang.Translatef(l, "command.team.info.not.found", name))
		return
	}
}

// Run ...
func (t TeamDisband) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	name := p.Name()
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}
	if !tm.Leader(name) {
		o.Error(lang.Translatef(l, "command.team.disband.must.leader"))
		return
	}
	if tm.Frozen() {
		o.Error(lang.Translatef(l, "command.team..dtr"))
		return
	}
	players := tm.Members
	data.DisbandTeam(tm)
	for _, m := range players {
		mem, ok := user.Lookup(m.Name)
		if ok {
			mem.UpdateState()
			mem.Player().Message(lang.Translatef(l, "command.team.disband.disbanded", name))
		}
	}
}

// Run ...
func (t TeamWho) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	n, _ := t.Name.Load()
	name := string(n)
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	if strings.TrimSpace(name) == "" {
		tm, ok := u.Team()
		if !ok {
			o.Error(lang.Translatef(l, "user.team-less"))
			return
		}
		o.Print(teamInformationFormat(tm, t.srv))
		return
	}
	var anyFound bool

	tm, ok := data.LoadTeam(name)
	if ok {
		o.Print(teamInformationFormat(tm, t.srv))
		anyFound = true
	}
	tm, ok = u.Team()
	if ok {
		o.Print(teamInformationFormat(tm, t.srv))
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
	u, err := data.LoadUser(src)
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
		if mem, ok := user.Lookup(m.Name); ok {
			mem.UpdateState()
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
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}
	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		o.Error(lang.Translatef(l, "command.team.kick.missing.permission"))
		return
	}
	if strings.EqualFold(p.Name(), string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.kick.self"))
		return
	}
	if tm.Leader(string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.kick.leader"))
		return
	}
	if tm.Captain(p.Name()) && tm.Captain(string(t.Member)) {
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
		if mem, ok := user.Lookup(m.Name); ok {
			mem.UpdateState()
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

// Run ...
func (t TeamPromote) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}
	if !tm.Leader(p.Name()) || !tm.Captain(p.Name()) {
		o.Error(lang.Translatef(l, "command.team.promote.missing.permission"))
		return
	}
	if strings.EqualFold(p.Name(), string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.promote.self"))
		return
	}
	if tm.Leader(p.Name()) {
		o.Error(lang.Translatef(l, "command.team.promote.leader"))
		return
	}
	if tm.Captain(p.Name()) && tm.Captain(string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.promote.captain"))
		return
	}
	if !tm.Member(string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.member.not.found", string(t.Member)))
		return
	}
	tm = tm.Promote(string(t.Member))
	data.SaveTeam(tm)
	if err != nil {
		o.Error(lang.Translatef(l, "team.save.fail"))
		return
	}
	rankName := "Captain"
	if tm.Leader(string(t.Member)) {
		rankName = "Leader"
	}
	for _, m := range tm.Members {
		if u, ok := user.Lookup(m.Name); ok {
			u.Message("command.team.promote.user.promoted", string(t.Member), rankName)
		}
	}
}

// Run ...
func (t TeamDemote) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}
	if !tm.Leader(p.Name()) {
		o.Error(lang.Translatef(l, "command.team.demote.missing.permission"))
		return
	}
	if strings.EqualFold(string(t.Member), p.Name()) {
		o.Error(lang.Translatef(l, "command.team.demote.self"))
		return
	}
	if !tm.Leader(string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.demote.leader"))
		return
	}
	if !tm.Member(string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.member.not.found", string(t.Member)))
		return
	}
	if !tm.Captain(string(t.Member)) && !tm.Captain(string(t.Member)) {
		o.Error(lang.Translatef(l, "command.team.demote.member"))
		return
	}
	tm.Demote(string(t.Member))
	data.SaveTeam(tm)
	if err != nil {
		o.Error(lang.Translatef(l, "team.save.fail"))
		return
	}
	for _, m := range tm.Members {
		if u, ok := user.Lookup(m.Name); ok {
			u.Message("command.team.demote.user.demoted", string(t.Member), "Member")
		}
	}
}

// Run ...
func (t TeamTop) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	teams := data.Teams()

	if len(teams) == 0 {
		o.Error("There are no teams.")
		return
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Points > teams[j].Points
	})

	var top string
	top += text.Colourf("        <yellow>Top Teams</yellow>\n")
	top += "\uE000\n"
	userTeam, ok := u.Team()
	for i, tm := range teams {
		if i > 9 {
			break
		}
		if ok && userTeam.Name == tm.Name {
			top += text.Colourf(" <grey>%d. <green>%s</green> (<yellow>%d</yellow>)</grey>\n", i+1, tm.DisplayName, tm.Points)
		} else {
			top += text.Colourf(" <grey>%d. <red>%s</red> (<yellow>%d</yellow>)</grey>\n", i+1, tm.DisplayName, tm.Points)
		}
	}
	top += "\uE000"
	p.Message(top)
}

// Run ...
func (t TeamClaim) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error("You are not in a team.")
	}
	if cl := tm.Claim; cl != (moose.Area{}) {
		o.Error("Your team already has a claim.")
		return
	}
	_, _ = p.Inventory().AddItem(item.NewStack(item.Hoe{Tier: item.ToolTierDiamond}, 1).WithValue("CLAIM_WAND", true))
}

// Run ...
func (t TeamUnClaim) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error("You are not in a team.")
	}
	if tm.Leader(u.Name) {
		o.Error("You are not the team leader.")
		return
	}
	tm = tm.WithClaim(moose.Area{}).WithHome(mgl64.Vec3{})
	data.SaveTeam(tm)
	p.Message(lang.Translate(p.Locale(), "command.unclaim.success"))
}

// Run ...
func (t TeamSetHome) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error("You are not in a team.")
	}
	cl := tm.Claim
	if cl == (moose.Area{}) {
		o.Error("Your team does not have a claim.")
		return
	}
	if !cl.Vec3WithinOrEqualXZ(p.Position()) {
		o.Error("You are not within your team's claim.")
		return
	}
	if !tm.Leader(u.Name) || !tm.Captain(u.Name) {
		o.Error("You are not the team leader or captain.")
		return
	}
	tm = tm.WithHome(p.Position())
	data.SaveTeam(tm)
	o.Print("Your team's home has been set.")
}

// Run ...
func (t TeamHome) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	us, ok := user.Lookup(p.Name())
	if err != nil || !ok {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error("You are not in a team.")
	}
	if u.PVP.Active() {
		o.Error("You cannot teleport while in combat.")
		return
	}

	if us.Home().Teleporting() {
		o.Error("You are already teleporting.")
		return
	}

	h := tm.Home
	if h == (mgl64.Vec3{}) {
		o.Error("Your team does not have a home.")
		return
	}
	if area.Spawn(p.World()).Vec3WithinOrEqualXZ(p.Position()) {
		p.Teleport(h)
		return
	}

	dur := time.Second * 10
	for _, tm := range data.Teams() {
		if tm.Claim.Vec3WithinOrEqualXZ(p.Position()) {
			dur = time.Second * 20
			break
		}
	}
	us.Home().Teleport(p, dur, h)
}

// Run ...
func (t TeamList) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	teams := data.Teams()
	if len(teams) == 0 {
		o.Error("There are no teams.")
		return
	}
	sort.Slice(teams, func(i, j int) bool {
		return user.TeamOnlineCount(teams[i]) > user.TeamOnlineCount(teams[j])
	})
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].DTR < teams[j].DTR
	})

	for _, tm := range teams {
		if tm.DTR <= 0 {
			teams = append(teams[:0], teams[1:]...)
		}
	}

	var list string
	list += text.Colourf("        <yellow>Team List</yellow>\n")
	list += "\uE000\n"
	userTeam, ok := u.Team()

	for i, tm := range teams {
		if i > 9 {
			break
		}

		dtr := text.Colourf("<green>%.1f■</green>", tm.DTR)
		if tm.DTR < 5 {
			dtr = text.Colourf("<yellow>%.1f■</yellow>", tm.DTR)
		}
		if tm.DTR <= 0 {
			dtr = text.Colourf("<red>%.1f■</red>", tm.DTR)
		}
		if ok && userTeam.Name == tm.Name {
			list += text.Colourf(" <grey>%d. <green>%s</green> (<green>%d/%d</green>)</grey> %s <yellow>DTR</yellow>\n", i+1, tm.DisplayName, user.TeamOnlineCount(tm), len(tm.Members), dtr)
		} else {
			list += text.Colourf(" <grey>%d. <red>%s</red> (<green>%d/%d</green>)</grey> %s <yellow>DTR</yellow>\n", i+1, tm.DisplayName, user.TeamOnlineCount(tm), len(tm.Members), dtr)
		}
	}
	list += "\uE000"
	p.Message(list)
}

// Run ...
func (t TeamFocusTeam) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error("You are not in a team.")
		return
	}
	targetTeam, ok := data.LoadTeam(string(t.Name))
	if !ok {
		o.Error("That team does not exist.")
		return
	}

	if tm.Name == targetTeam.Name {
		o.Error("You cannot focus your own team.")
		return
	}
	// tm.FocusTeam(targetTeam)

	// us, ok := user.Lookup(p.Name())
	// if !ok {
	// 	return
	// }
	// members, _ := u.Focusing()
	// for _, m := range members {
	// 	p, ok := user.LookupName(m.Name())
	// 	if !ok {
	// 		continue
	// 	}
	// 	p.UpdateState()
	// }

	// tm.Broadcast("command.team.focus", targetTeam.DisplayName())
}

// Run ...
func (t TeamUnFocus) Run(s cmd.Source, o *cmd.Output) {
	// p, ok := s.(*player.Player)
	// if !ok {
	// 	return
	// }
	// u, err := data.LoadUser(p.Name(), p.Handler().(*user.Handler).XUID())
	// if err != nil {
	// 	return
	// }
	// tm, ok := u.Team()
	// if !ok {
	// 	o.Error("You are not in a team.")
	// 	return
	// }
	// if !ok {
	// 	o.Error("You are not in a team.")
	// 	return
	// }

	// var name string

	// focused, okTeam := tm.FocusedTeam()
	// focusedPlayer, okPlayer := tm.FocusedPlayer()

	// if okTeam {
	// 	name = focused.DisplayName()
	// } else if okPlayer {
	// 	name = focusedPlayer
	// } else {
	// 	o.Error("Your team is not focusing another team.")
	// 	return
	// }
	// members, _ := u.Focusing()
	// tm.UnFocus()

	// for _, m := range members {
	// 	p, ok := user.LookupName(m.Name())
	// 	if !ok {
	// 		continue
	// 	}
	// 	p.UpdateState()
	// }

	// tm.Broadcast("command.team.unfocus", name)
}

// Run ...
func (t TeamFocusPlayer) Run(s cmd.Source, o *cmd.Output) {
	// l := locale(s)
	// p, ok := s.(*player.Player)
	// if !ok {
	// 	return
	// }
	// u, ok := user.Lookup(p)
	// if !ok {
	// 	return
	// }
	// tm, ok := u.Team()
	// if !ok {
	// 	o.Error("You are not in a team.")
	// 	return
	// }
	// if len(t.Target) > 1 {
	// 	o.Error(lang.Translatef(l, "command.targets.exceed"))
	// 	return
	// }
	// trg := t.Target[0]
	// target, ok := trg.(*player.Player)
	// if !ok {
	// 	o.Error("You must target a player.")
	// 	return
	// }
	// targetUser, ok := user.Lookup(target)
	// if !ok {
	// 	o.Error("That player is not online.")
	// 	return
	// }

	// if targetUser == u {
	// 	o.Error("You cannot focus yourself.")
	// 	return
	// }

	// if _, ok := tm.Member(targetUser.Name()); ok {
	// 	o.Error("You cannot focus a member of your team.")
	// 	return
	// }

	// tm.FocusPlayer(targetUser.Name())
	// targetUser.UpdateState()
	// tm.Broadcast("command.team.focus", target.Name())
}

// Run ...
func (t TeamChat) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	_, ok = u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}

	us, ok := user.Lookup(p.Name())
	if !ok {
		return
	}
	switch us.ChatType() {
	case 2:
		us.UpdateChatType(1)
	case 1:
		us.UpdateChatType(2)
	}
}

// Run ...
func (t TeamWithdraw) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}

	if !tm.Leader(p.Name()) || !tm.Captain(p.Name()) {
		o.Error("You cannot withdraw any balance from your team.")
		return
	}

	amt := t.Balance
	if amt < 1 {
		o.Error("You must provide a minimum balance of $1.")
		return
	}

	if amt > tm.Balance {
		o.Errorf("Your team does not have a balance of $%.2f.", amt)
		return
	}

	tm = tm.WithBalance(tm.Balance - amt)
	u.Balance = u.Balance + amt

	data.SaveTeam(tm)
	data.SaveUser(u)
	o.Print(text.Colourf("<green>You withdrew $%.2f from %s.</green>", amt, tm.DisplayName))
}

// Run ...
func (t TeamDeposit) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}

	amt := t.Balance
	if amt < 1 {
		o.Error("You must provide a minimum balance of $1.")
		return
	}

	if amt > u.Balance {
		o.Errorf("You do not have a balance of $%.2f.", amt)
		return
	}

	tm = tm.WithBalance(tm.Balance + amt)
	u.Balance = u.Balance - amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	o.Print(text.Colourf("<green>You deposited $%d into %s.</green>", int(amt), tm.DisplayName))
}

// Run ...
func (t TeamWithdrawAll) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}

	if !tm.Leader(p.Name()) || !tm.Captain(p.Name()) {
		o.Error("You cannot withdraw any balance from your team.")
		return
	}

	amt := tm.Balance
	if amt < 1 {
		o.Error("Your team's balance is lower than $1.")
		return
	}

	tm = tm.WithBalance(tm.Balance - amt)
	u.Balance = u.Balance + amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	o.Print(text.Colourf("<green>You withdrew $%d from %s.</green>", amt, tm.Name))
}

// Run ...
func (t TeamDepositAll) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}
	tm, ok := u.Team()
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}
	if !ok {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}

	amt := u.Balance
	if amt < 1 {
		o.Error("Your balance is lower than $1.")
		return
	}

	tm = tm.WithBalance(tm.Balance + amt)
	u.Balance = u.Balance - amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	o.Print(text.Colourf("<green>You deposited $%d into %s.</green>", amt, tm.Name))
}

// Run ...
func (t TeamStuck) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	us, ok := user.Lookup(p.Name())
	if err != nil || !ok {
		return
	}
	pos := safePosition(p, u, 24)
	if pos == (cube.Pos{}) {
		us.Message("command.team.stuck.no-safe")
		return
	}

	if us.Stuck().Teleporting() {
		o.Error("You are already stucking, step-bro.")
		return
	}

	us.Message("command.team.stuck.teleporting", pos.X(), pos.Y(), pos.Z(), 30)
	us.Stuck().Teleport(p, time.Second*30, mgl64.Vec3{
		float64(pos.X()),
		float64(pos.Y()),
		float64(pos.Z()),
	})
}

func safePosition(p *player.Player, u data.User, radius int) cube.Pos {
	pos := cube.PosFromVec3(p.Position())
	minX := pos.X() - radius
	maxX := pos.X() + radius
	minZ := pos.Z() - radius
	maxZ := pos.Z() + radius

	for x := minX; x < maxX; x++ {
		for z := minZ; z < maxZ; z++ {
			at := pos.Add(cube.Pos{x, 0, z})
			for _, tm := range data.Teams() {
				if tm.Claim != (moose.Area{}) {
					if tm.Claim.Vec3WithinOrEqualXZ(at.Vec3Centre()) {
						if t, ok := u.Team(); ok && t.Name == tm.Name {
							y := p.World().Range().Max()
							for y > pos.Y() {
								y--
								b := p.World().Block(cube.Pos{x, y, z})
								if b != (block.Air{}) {
									return cube.Pos{x, y, z}
								}
							}
						}
					}
				}
			}

			for _, area := range append(area.Protected(p.World()), area.Wilderness(p.World())) {
				if area.Vec3WithinOrEqualXZ(at.Vec3Centre()) {
					y := p.World().Range().Max()
					for y > pos.Y() {
						y--
						b := p.World().Block(cube.Pos{x, y, z})
						if b != (block.Air{}) {
							return cube.Pos{x, y, z}
						}
					}
				}
			}
		}
	}
	return cube.Pos{}
}

// Allow ...
func (TeamCreate) Allow(src cmd.Source) bool {
	return allow(src, false)
}

// Allow ...
func (TeamDisband) Allow(s cmd.Source) bool {
	return allow(s, false)
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
	u, err := data.LoadUser(p.Name())
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
	u, err := data.LoadUser(p.Name())
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

// teamInformationFormat returns a formatted string containing the information of the faction.
func teamInformationFormat(t data.Team, srv *server.Server) string {
	var formattedRegenerationTime string
	var formattedDtr string
	var formattedLeader string
	var formattedCaptains []string
	var formattedMembers []string
	if time.Now().Before(t.RegenerationTime) {
		formattedRegenerationTime = text.Colourf("\n <yellow>Time Until Regen</yellow> <blue>%s</blue>", time.Until(t.RegenerationTime).Round(time.Second))
	}
	formattedDtr = t.DTRString()
	var onlineCount int
	for _, p := range t.Members {
		_, ok := srv.PlayerByName(p.DisplayName)
		if ok {
			if t.Leader(p.Name) {
				formattedLeader = text.Colourf("<green>%s</green>", p.DisplayName)
			} else if t.Captain(p.Name) {
				formattedCaptains = append(formattedCaptains, text.Colourf("<green>%s</green>", p.DisplayName))
			} else {
				formattedMembers = append(formattedMembers, text.Colourf("<green>%s</green>", p.DisplayName))
			}
			onlineCount++
		} else {
			if t.Leader(p.Name) {
				formattedLeader = text.Colourf("<grey>%s</grey>", p.DisplayName)
			} else if t.Captain(p.Name) {
				formattedCaptains = append(formattedCaptains, text.Colourf("<grey>%s</grey>", p.DisplayName))
			} else {
				formattedMembers = append(formattedMembers, text.Colourf("<grey>%s</grey>", p.DisplayName))
			}
		}
	}
	if len(formattedCaptains) == 0 {
		formattedCaptains = []string{"None"}
	}
	if len(formattedMembers) == 0 {
		formattedMembers = []string{"None"}
	}
	var home string
	h := t.Home
	if h.X() == 0 && h.Y() == 0 && h.Z() == 0 {
		home = "not set"
	} else {
		home = fmt.Sprintf("%.0f, %.0f, %.0f", h.X(), h.Y(), h.Z())
	}
	return text.Colourf(
		"\uE000\n <blue>%s</blue> <grey>[%d/%d]</grey> <dark-aqua>-</dark-aqua> <yellow>HQ:</yellow> %s\n "+
			"<yellow>Leader: </yellow>%s\n "+
			"<yellow>Captains: </yellow>%s\n "+
			"<yellow>Members: </yellow>%s\n "+
			"<yellow>Balance: </yellow><blue>$%2.f</blue>\n "+
			"<yellow>Points: </yellow><blue>%d</blue>\n "+
			"<yellow>Deaths until Raidable: </yellow>%s%s\n\uE000", t.DisplayName, onlineCount, len(t.Members), home, formattedLeader, strings.Join(formattedCaptains, ", "), strings.Join(formattedMembers, ", "), t.Balance, t.Points, formattedDtr, formattedRegenerationTime)
}

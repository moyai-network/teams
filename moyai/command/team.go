package command

import (
	"fmt"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/internal/timeutil"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/team"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
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
}

// TeamWho is a command used to get information about a team.
type TeamWho struct {
	Sub  cmd.SubCommand           `cmd:"who"`
	Name cmd.Optional[teamMember] `optional:"" cmd:"name"`
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

type TeamDelete struct {
	Sub  cmd.SubCommand `cmd:"delete"`
	Name teamName       `cmd:"name"`
}

type TeamSetDTR struct {
	Sub  cmd.SubCommand `cmd:"setdtr"`
	Name teamName       `cmd:"name"`
	DTR  float64        `cmd:"dtr"`
}

func (t TeamSetDTR) Run(s cmd.Source, o *cmd.Output) {
	tm, err := data.LoadTeamFromName(strings.ToLower(string(t.Name)))
	if err != nil {
		o.Error("Invalid Team.")
		return
	}

	tm = tm.WithDTR(t.DTR)
	data.SaveTeam(tm)
	o.Printf("Successfully set DTR to %v", t.DTR)
}

func (t TeamDelete) Run(s cmd.Source, o *cmd.Output) {
	tm, err := data.LoadTeamFromName(strings.ToLower(string(t.Name)))
	if err != nil {
		o.Error("Invalid Team.")
		return
	}

	players := tm.Members
	data.DisbandTeam(tm)
	o.Print("Disbanded faction.")
	for _, m := range players {
		if mem, ok := user.Lookup(m.Name); ok {
			user.UpdateState(mem)
			user.Messagef(mem, "command.team.disband.disbanded", "MANAGEMENT")
		}
	}
}

/*func (t TeamMap) Run(s cmd.Source, o *cmd.Output) {
	p, _ := s.(*player.Player)
	posX := int(p.Position().X())
	posZ := int(p.Position().Z())

	i := 0
	rows := 10
	disp := make([][]string, rows)
	for i := range disp {
		disp[i] = make([]string, 10)
		for j := range disp[i] {
			disp[i][j] = text.Colourf("<grey>█</grey>")
		}
	}
	areas := area.Protected(p.World())
	for _, t := range data.Teams() {
		areas = append(areas, moose.NewNamedArea(t.Claim.Max(), t.Claim.Min(), t.Name))
	}

	for x := posX - 5; x < posX+5; x++ {
		for z := posZ - 5; z < posZ+5; z++ {
			found := false
			for _, a := range areas {
				if a.Vec2WithinOrEqualFloor(mgl64.Vec2{float64(x), float64(z)}) {
					found = true
					if found {
						if slices.Contains(area.Roads(p.World()), a) {
							disp[i][z-(posZ-5)] = text.Colourf("<black>█</black>")
						} else if a == area.Spawn(p.World()) {
							disp[i][z-(posZ-5)] = text.Colourf("<green>█</green>")
						} else if slices.Contains(area.KOTHs(p.World()), a) {
							disp[i][z-(posZ-5)] = text.Colourf("<dark-red>█</dark-red>")
						} else if a == area.WarZone(p.World()) {
							if disp[i][z-(posZ-5)] == text.Colourf("<grey>█</grey>") {
								disp[i][z-(posZ-5)] = text.Colourf("<red>█</red>")
							}
						} else {
							if a != area.Wilderness(p.World()) {
								disp[i][z-(posZ-5)] = text.Colourf("<aqua>█</aqua>")

							}
						}
					} else {
						disp[i][z-(posZ-5)] = text.Colourf("<grey>█</grey>")
					}
				}
			}
		}
		i++
	}
	disp[5][5] = text.Colourf("<gold>✪</gold>")
	var b strings.Builder
	for _, row := range disp {
		b.WriteString(strings.Join(row, ""))
		b.WriteString("\n")
	}
	p.Message("\uE000\uE000")
	p.Message(b.String())
	p.Message("\uE000\uE000")
}

func createEmptyMap(size int) [][]rune {
	// Create an empty map with the specified size
	myMap := make([][]rune, size)
	for i := range myMap {
		myMap[i] = make([]rune, size)
		for j := range myMap[i] {
			myMap[i][j] = '.'
		}
	}
	return myMap
}

func markClaim(myMap [][]rune, x1, z1, x2, z2 int) {
	// Mark the claimed area on the map
	for x := min(x1, x2); x <= max(x1, x2); x++ {
		for z := min(z1, z2); z <= max(z1, z2); z++ {
			myMap[z][x] = 'X'
		}
	}
}

func displayMap(myMap [][]rune, size, centerX, centerZ int) {
	// Display the map with only the nearby claims
	startX := max(0, centerX-size/2)
	startZ := max(0, centerZ-size/2)
	endX := min(size-1, centerX+size/2)
	endZ := min(size-1, centerZ+size/2)

	for z := startZ; z <= endZ; z++ {
		for x := startX; x <= endX; x++ {
			fmt.Print(string(myMap[z][x]))
		}
		fmt.Println()
	}
}*/

// Run ...
func (t TeamCreate) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err == nil {
		user.Messagef(p, "team.create.already")
		return
	}
	t.Name = strings.TrimSpace(t.Name)

	if !regex.MatchString(t.Name) {
		user.Messagef(p, "team.create.name.invalid")
		return
	}
	if len(t.Name) < 3 {
		user.Messagef(p, "team.create.name.short")
		return
	}
	if len(t.Name) > 15 {
		user.Messagef(p, "team.create.name.long")
		return
	}

	if u.Teams.Create.Active() {
		user.Messagef(p, "command.team.create.cooldown", timeutil.FormatDuration(u.Teams.Create.Remaining()))
		return
	}

	if _, err = data.LoadTeamFromName(t.Name); err == nil {
		user.Messagef(p, "team.create.exists")
		return
	}

	tm = data.DefaultTeam(t.Name).WithMembers(data.DefaultMember(p.XUID(), p.Name()).WithRank(3))
	data.SaveTeam(tm)

	u.Teams.Create.Set(time.Minute * 3)
	data.SaveUser(u)

	user.Messagef(p, "team.create.success", tm.DisplayName)
	user.Broadcastf("team.create.success.broadcast", p.Name(), tm.DisplayName)
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
		user.Messagef(p, "team.invite.self")
		return
	}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}

	if tm.Frozen() {
		user.Messagef(p, "command.team.dtr")
		return
	}

	if len(tm.Members) == 6 {
		user.Messagef(p, "command.team.join.full")
		return
	}

	if tm.Member(target.Name()) {
		user.Messagef(p, "team.invite.member", target.Name())
		return
	}

	targetUser, err := data.LoadUserFromName(target.Name())
	if err != nil {
		return
	}

	_, err = data.LoadTeamFromMemberName(target.Name())
	if err == nil {
		user.Messagef(p, "team.invite.has-team", target.Name())

	}

	if targetUser.Teams.Invitations.Active(tm.Name) {
		user.Messagef(p, "team.invite.already", target.Name())
		return
	}

	targetUser.Teams.Invitations.Set(tm.Name, time.Minute*5)
	data.SaveUser(targetUser)

	team.Broadcastf(tm, "team.invite.success.broadcast", target.Name())
	user.Messagef(p, "team.invite.target", tm.DisplayName)
}

// Run ...
func (t TeamJoin) Run(src cmd.Source, out *cmd.Output) {
	p := src.(*player.Player)

	_, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "team.join.error")
		return
	}

	tm, err := data.LoadTeamFromName(string(t.Team))
	if err != nil {
		user.Messagef(p, "command.team.not.found")
		return
	}

	if tm.Frozen() {
		user.Messagef(p, "command.team.dtr")
		return
	}

	if len(tm.Members) == 6 {
		user.Messagef(p, "command.team.join.full")
		return
	}

	tm = tm.WithMembers(append(tm.Members, data.DefaultMember(p.XUID(), p.Name()))...)
	tm = tm.WithDTR(tm.DTR + 1)
	data.SaveTeam(tm)

	team.Broadcastf(tm, "team.join.broadcast", p.Name())
}

// Run ...
func (t TeamInformation) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	n, _ := t.Name.Load()
	name := string(n)
	if strings.TrimSpace(name) == "" {
		tm, err := data.LoadTeamFromMemberName(p.Name())
		if err != nil {
			user.Messagef(p, "user.team-less")
			return
		}
		o.Print(teamInformationFormat(tm))
		return
	}
	var anyFound bool

	tm, err := data.LoadTeamFromName(strings.ToLower(name))
	if err == nil {
		o.Print(teamInformationFormat(tm))
		anyFound = true
	}

	tm, err = data.LoadTeamFromMemberName(strings.ToLower(name))
	if err == nil {
		o.Print(teamInformationFormat(tm))
		anyFound = true
	}

	if !anyFound {
		user.Messagef(p, "command.team.info.not.found", name)
		return
	}
}

// Run ...
func (t TeamWho) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	n, _ := t.Name.Load()
	name := string(n)
	if strings.TrimSpace(name) == "" {
		tm, err := data.LoadTeamFromMemberName(p.Name())
		if err != nil {
			user.Messagef(p, "user.team-less")
			return
		}
		o.Print(teamInformationFormat(tm))
		return
	}
	var anyFound bool

	tm, err := data.LoadTeamFromName(strings.ToLower(name))
	if err == nil {
		o.Print(teamInformationFormat(tm))
		anyFound = true
	}

	tm, err = data.LoadTeamFromMemberName(strings.ToLower(name))
	if err == nil {
		o.Print(teamInformationFormat(tm))
		anyFound = true
	}

	if !anyFound {
		user.Messagef(p, "command.team.info.not.found", name)
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

	name := p.Name()
	tm, err := data.LoadTeamFromMemberName(name)
	if err != nil {
		o.Error(lang.Translatef(l, "user.team-less"))
		return
	}

	if !tm.Leader(name) {
		o.Error(lang.Translatef(l, "command.team.disband.must.leader"))
		return
	}

	if tm.Frozen() {
		o.Error(lang.Translatef(l, "command.team.dtr"))
		return
	}

	players := tm.Members
	data.DisbandTeam(tm)

	for _, m := range players {
		if mem, ok := user.Lookup(m.Name); ok {
			user.UpdateState(mem)
			user.Messagef(mem, "command.team.disband.disbanded")
		}
	}
}

// Run ...
func (t TeamLeave) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}

	if tm.Leader(p.Name()) {
		o.Error(lang.Translatef(l, "command.team.leave.leader"))
		return
	}

	if tm.Frozen() {
		o.Error(lang.Translatef(l, "command.team.dtr"))
		return
	}

	tm = tm.WithoutMember(p.Name())
	for _, m := range tm.Members {
		if mem, ok := user.Lookup(m.Name); ok {
			user.UpdateState(mem)
			user.Messagef(mem, "command.team.leave.user.left")
		}
	}
	data.SaveTeam(tm)
}

// Run ...
func (t TeamKick) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}
	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		user.Messagef(p, "command.team.kick.missing.permission")
		return
	}
	if strings.EqualFold(p.Name(), string(t.Member)) {
		user.Messagef(p, "command.team.kick.self")
		return
	}
	if tm.Leader(string(t.Member)) {
		user.Messagef(p, "command.team.kick.leader")
		return
	}
	if tm.Captain(p.Name()) && tm.Captain(string(t.Member)) {
		user.Messagef(p, "command.team.kick.captain")
		return
	}
	if !tm.Member(string(t.Member)) {
		user.Messagef(p, "command.team.kick.not.found", string(t.Member))
		return
	}

	if tm.Frozen() {
		user.Messagef(p, "command.team.dtr")
		return
	}

	us, ok := user.Lookup(string(t.Member))
	if ok {
		user.Messagef(us, "command.team.kick.user.kicked", tm.DisplayName)
	}
	tm = tm.WithoutMember(string(t.Member))
	tm = tm.WithDTR(tm.MaxDTR())
	for _, m := range tm.Members {
		if mem, ok := user.Lookup(m.Name); ok {
			user.UpdateState(mem)
			user.Messagef(mem, "command.team.kick.kicked")
		}
	}

	data.SaveTeam(tm)
}

// Run ...
func (t TeamPromote) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}
	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		user.Messagef(p, "command.team.promote.missing.permission")
		return
	}
	if strings.EqualFold(p.Name(), string(t.Member)) {
		user.Messagef(p, "command.team.promote.self")
		return
	}
	if tm.Leader(p.Name()) {
		user.Messagef(p, "command.team.promote.leader")
		return
	}
	if tm.Captain(p.Name()) && tm.Captain(string(t.Member)) {
		user.Messagef(p, "command.team.promote.captain")
		return
	}
	if !tm.Member(string(t.Member)) {
		user.Messagef(p, "command.team.member.not.found", string(t.Member))
		return
	}

	tm = tm.Promote(string(t.Member))
	data.SaveTeam(tm)
	if err != nil {
		user.Messagef(p, "team.save.fail")
		return
	}
	rankName := "Captain"
	if tm.Leader(string(t.Member)) {
		rankName = "Leader"
	}
	team.Broadcastf(tm, "command.team.promote.user.promoted", string(t.Member), rankName)
}

// Run ...
func (t TeamDemote) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
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
	team.Broadcastf(tm, "command.team.demote.user.demoted", string(t.Member), "Member")
}

// Run ...
func (t TeamTop) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	teams, err := data.LoadAllTeams()
	if err != nil {
		o.Error("Failed to load teams. Contact an administrator.")
	}

	if len(teams) == 0 {
		user.Messagef(p, "command.team.top.none")
		return
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Points > teams[j].Points
	})

	var top string
	top += text.Colourf("        <yellow>Top Teams</yellow>\n")
	top += "\uE000\n"
	userTeam, err := data.LoadTeamFromMemberName(p.Name())
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
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}
	if cl := tm.Claim; cl != (area.Area{}) {
		user.Messagef(p, "team.has-claim")
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
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}
	if !tm.Leader(p.Name()) {
		user.Messagef(p, "team.not-leader")
		return
	}
	tm = tm.WithClaim(area.Area{}).WithHome(mgl64.Vec3{})
	data.SaveTeam(tm)
	user.Messagef(p, "command.unclaim.success")
}

// Run ...
func (t TeamSetHome) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}
	cl := tm.Claim
	if cl == (area.Area{}) {
		user.Messagef(p, "team.claim.none")
		return
	}
	if !cl.Vec3WithinOrEqualXZ(p.Position()) {
		user.Messagef(p, "team.claim.not-within")
		return
	}
	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		user.Messagef(p, "team.not-leader-or-captain")
		return
	}
	tm = tm.WithHome(p.Position())
	data.SaveTeam(tm)
	user.Messagef(p, "command.team.home.set")
}

// Run ...
func (t TeamHome) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		o.Error("Failed to load user handler. Contact an administrator.")
		return
	}
	if h.Combat().Active() {
		user.Messagef(p, "user.teleport.combat")
		return
	}

	if h.Home().Ongoing() {
		user.Messagef(p, "user.already.teleporting")
		return
	}

	hm := tm.Home
	if hm == (mgl64.Vec3{}) {
		user.Messagef(p, "command.team.home.none")
		return
	}
	if area.Spawn(p.World()).Vec3WithinOrEqualXZ(p.Position()) {
		p.Teleport(hm)
		return
	}

	dur := time.Second * 10
	teams, err := data.LoadAllTeams()
	if err != nil {
		o.Error("Failed to load teams. Contact an administrator.")
	}

	for _, tm := range teams {
		if tm.Claim.Vec3WithinOrEqualXZ(p.Position()) {
			dur = time.Second * 20
			break
		}
	}
	h.Home().Teleport(p, dur, hm)
}

// Run ...
func (t TeamList) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	teams, err := data.LoadAllTeams()
	if err != nil {
		o.Error("Failed to load teams. Contact an administrator.")
	}

	sort.Slice(teams, func(i, j int) bool {
		return len(team.OnlineMembers(teams[i])) > len(team.OnlineMembers(teams[j]))
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
	userTeam, err := data.LoadTeamFromMemberName(p.Name())

	for i, tm := range teams {
		if i > 9 {
			break
		}

		onlineCount := len(team.OnlineMembers(tm))
		dtr := tm.DTRString()
		if err == nil && userTeam.Name == tm.Name {
			list += text.Colourf(" <grey>%d. <green>%s</green> (<green>%d/%d</green>)</grey> %s\n", i+1, tm.DisplayName, onlineCount, len(tm.Members), dtr)
		} else {
			list += text.Colourf(" <grey>%d. <red>%s</red> (<green>%d/%d</green>)</grey> %s\n", i+1, tm.DisplayName, onlineCount, len(tm.Members), dtr)
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
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}
	targetTeam, err := data.LoadTeamFromName(strings.ToLower(string(t.Name)))
	if err != nil {
		user.Messagef(p, "command.team.not.found", string(t.Name))
		return
	}

	if tm.Name == targetTeam.Name {
		user.Messagef(p, "command.team.focus.self")
		return
	}
	tm = tm.WithTeamFocus(targetTeam)
	data.SaveTeam(tm)

	for _, m := range team.OnlineMembers(tm) {
		user.UpdateState(m)
	}

	team.Broadcastf(tm, "command.team.focus", targetTeam.DisplayName)
}

// Run ...
func (t TeamUnFocus) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}
	focus := tm.Focus

	if focus.Type() == data.FocusTypeNone() {
		user.Messagef(p, "command.team.focus.none")
		return
	}

	tm = tm.WithoutFocus()
	data.SaveTeam(tm)

	for _, m := range team.OnlineMembers(tm) {
		user.UpdateState(m)
	}

	team.Broadcastf(tm, "command.team.unfocus", focus.Value())
}

// Run ...
func (t TeamFocusPlayer) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}
	if len(t.Target) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	target, ok := t.Target[0].(*player.Player)
	if !ok {
		o.Error("You must target a player.")
		return
	}

	if strings.EqualFold(target.Name(), p.Name()) {
		o.Error("You cannot focus yourself.")
		return
	}

	if tm.Member(target.Name()) {
		o.Error("You cannot focus a member of your team.")
		return
	}

	tm = tm.WithPlayerFocus(target.Name())
	user.UpdateState(target)
	team.Broadcastf(tm, "command.team.focus", target.Name())
}

// Run ...
func (t TeamChat) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	switch u.Teams.ChatType {
	case 1:
		u.Teams.ChatType = 0
		o.Print(text.Colourf("<green>Switched to global chat.</green>"))
	case 0:
		u.Teams.ChatType = 1
		o.Print(text.Colourf("<green>Switched to faction chat.</green>"))
	}
	data.SaveUser(u)
}

// Run ...
func (t TeamWithdraw) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}

	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
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
	u.Teams.Balance = u.Teams.Balance + amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	o.Print(text.Colourf("<green>You withdrew $%.2f from %s.</green>", amt, tm.DisplayName))
}

// Run ...
func (t TeamDeposit) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}

	amt := t.Balance
	if amt < 1 {
		o.Error("You must provide a minimum balance of $1.")
		return
	}

	if amt > u.Teams.Balance {
		o.Errorf("You do not have a balance of $%.2f.", amt)
		return
	}

	tm = tm.WithBalance(tm.Balance + amt)
	u.Teams.Balance = u.Teams.Balance - amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	o.Print(text.Colourf("<green>You deposited $%d into %s.</green>", int(amt), tm.DisplayName))
}

// Run ...
func (t TeamWithdrawAll) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}

	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		o.Error("You cannot withdraw any balance from your team.")
		return
	}

	amt := tm.Balance
	if amt < 1 {
		o.Error("Your team's balance is lower than $1.")
		return
	}

	tm = tm.WithBalance(tm.Balance - amt)
	u.Teams.Balance = u.Teams.Balance + amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	o.Print(text.Colourf("<green>You withdrew $%d from %s.</green>", int(amt), tm.Name))
}

// Run ...
func (t TeamDepositAll) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}

	amt := u.Teams.Balance
	if amt < 1 {
		o.Error("Your balance is lower than $1.")
		return
	}

	tm = tm.WithBalance(tm.Balance + amt)
	u.Teams.Balance = u.Teams.Balance - amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	o.Print(text.Colourf("<green>You deposited $%d into %s.</green>", int(amt), tm.Name))
}

// Run ...
func (t TeamStuck) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	pos := safePosition(p, 24)
	if pos == (cube.Pos{}) {
		user.Messagef(p, "command.team.stuck.no-safe")
		return
	}

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	if h.Stuck().Ongoing() {
		o.Error("You are already stucking, step-bro.")
		return
	}

	user.Messagef(p, "command.team.stuck.teleporting", pos.X(), pos.Y(), pos.Z(), 30)
	h.Stuck().Teleport(p, time.Second*30, mgl64.Vec3{
		float64(pos.X()),
		float64(pos.Y()),
		float64(pos.Z()),
	})
}

func safePosition(p *player.Player, radius int) cube.Pos {
	pos := cube.PosFromVec3(p.Position())
	minX := pos.X() - radius
	maxX := pos.X() + radius
	minZ := pos.Z() - radius
	maxZ := pos.Z() + radius

	teams, err := data.LoadAllTeams()
	if err != nil {
		return cube.Pos{}
	}
	for x := minX; x < maxX; x++ {
		for z := minZ; z < maxZ; z++ {
			at := pos.Add(cube.Pos{x, 0, z})
			for _, tm := range teams {
				if tm.Claim != (area.Area{}) {
					if tm.Claim.Vec3WithinOrEqualXZ(at.Vec3Centre()) {
						if t, err := data.LoadTeamFromMemberName(p.Name()); err == nil && t.Name == tm.Name {
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

func (TeamDelete) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}

func (TeamSetDTR) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
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
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	for t, i := range u.Teams.Invitations {
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
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		return nil
	}

	var members []string
	for _, m := range tm.Members {
		if !strings.EqualFold(m.Name, p.Name()) {
			members = append(members, m.DisplayName)
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
	tms, _ := data.LoadAllTeams()

	for _, tm := range tms {
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
	tms, err := data.LoadAllTeams()
	if err != nil {
		return members
	}

	for _, tm := range tms {
		for _, m := range tm.Members {
			members = append(members, m.DisplayName)
		}
	}

	return members
}

// teamInformationFormat returns a formatted string containing the information of the faction.
func teamInformationFormat(t data.Team) string {
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
		u, _ := data.LoadUserFromXUID(p.XUID)
		_, ok := user.Lookup(p.DisplayName)
		if ok {
			onlineCount++
		}
		format := formatMember(p.DisplayName, u.Teams.Stats.Kills, ok)

		if t.Leader(p.Name) {
			formattedLeader = format
		} else if t.Captain(p.Name) {
			formattedCaptains = append(formattedCaptains, format)
		} else {
			formattedMembers = append(formattedMembers, format)
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

func formatMember(name string, kills int, online bool) string {
	if online {
		return text.Colourf("<green>%s</green><dark-red>[%d]</dark-red>", name, kills)
	}
	return text.Colourf("<grey>%s</grey><dark-red>[%d]</dark-red>", name, kills)
}

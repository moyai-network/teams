package command

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/moyai-network/teams/moyai"

	"github.com/moyai-network/teams/moyai/colour"

	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/internal/timeutil"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/team"

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
	Name teamName       `cmd:"team"`
}

// TeamFocusPlayer is a command used to focus on a player.
type TeamFocusPlayer struct {
	Sub    cmd.SubCommand `cmd:"focus"`
	Target []cmd.Target   `cmd:"player"`
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

// TeamW is a alias for the TeamWithdraw command.
type TeamW struct {
	Sub     cmd.SubCommand `cmd:"w"`
	Balance float64        `cmd:"balance"`
}

// TeamDeposit is a command used to deposit money into a team.
type TeamDeposit struct {
	Sub     cmd.SubCommand `cmd:"deposit"`
	Balance float64        `cmd:"balance"`
}

// TeamD is a alias for the TeamDeposit command.
type TeamD struct {
	Sub     cmd.SubCommand `cmd:"d"`
	Balance float64        `cmd:"balance"`
}

// TeamWithdrawAll is a command used to withdraw all the money from a team.
type TeamWithdrawAll struct {
	Sub cmd.SubCommand `cmd:"withdraw"`
	All cmd.SubCommand `cmd:"all"`
}

// TeamWAll is an alias for the TeamWithdrawAll command.
type TeamWAll struct {
	Sub cmd.SubCommand `cmd:"w"`
	All cmd.SubCommand `cmd:"all"`
}

// TeamDepositAll is a command used to deposit all of a user's money into a team.
type TeamDepositAll struct {
	Sub cmd.SubCommand `cmd:"deposit"`
	All cmd.SubCommand `cmd:"all"`
}

// TeamDAll is an alias for the TeamDepositAll command.
type TeamDAll struct {
	Sub cmd.SubCommand `cmd:"d"`
	All cmd.SubCommand `cmd:"all"`
}

// TeamStuck is a command to teleport to a safe position
type TeamStuck struct {
	Sub cmd.SubCommand `cmd:"stuck"`
}

type TeamDelete struct {
	adminAllower
	Sub  cmd.SubCommand `cmd:"delete"`
	Name teamName       `cmd:"name"`
}

type TeamSetDTR struct {
	adminAllower
	Sub  cmd.SubCommand `cmd:"setdtr"`
	Name teamName       `cmd:"name"`
	DTR  float64        `cmd:"dtr"`
}

type TeamMap struct {
	Sub cmd.SubCommand `cmd:"map"`
}

type TeamClearMap struct {
	Sub cmd.SubCommand `cmd:"clearmap"`
}

// TeamRally is a command that enables waypoint to rally.
type TeamRally struct {
	Sub cmd.SubCommand `cmd:"rally"`
}

// TeamUnRally is a command that disable waypoint to rally.
type TeamUnRally struct {
	Sub cmd.SubCommand `cmd:"unrally"`
}

// TeamRename is a command to rename your team.
type TeamRename struct {
	Sub  cmd.SubCommand `cmd:"rename"`
	Name string         `cmd:"name"`
}

// TeamCamp is a command to teleport close to a base.
type TeamCamp struct {
	Sub  cmd.SubCommand `cmd:"camp"`
	Team teamName
}

// TeamIncrementDTR is a command to increment DTR
type TeamIncrementDTR struct {
	adminAllower
	Sub  cmd.SubCommand `cmd:"incdtr"`
	Name teamName
}

// TeamDecrementDTR is a command to decrement DTR
type TeamDecrementDTR struct {
	adminAllower
	Sub  cmd.SubCommand `cmd:"decdtr"`
	Name teamName
}

// TeamResetRegen is a command to reset regeneration
type TeamResetRegen struct {
	adminAllower
	Sub  cmd.SubCommand `cmd:"resetregen"`
	Name teamName
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
	src, _ := s.(cmd.NamedTarget)
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
			moyai.Messagef(mem, src.Name(), "command.team.disband.disbanded", "MANAGEMENT")
		}
	}
}

func (t TeamMap) Run(src cmd.Source, _ *cmd.Output) {
	p, _ := src.(*player.Player)
	h, _ := p.Handler().(*user.Handler)

	areas := make([]area.NamedArea, 0)
	teams, err := data.LoadAllTeams()
	if err != nil {
		moyai.Messagef(p, "command.team.load.fail")
		return
	}
	for _, t := range teams {
		areas = append(areas, area.NewNamedArea(t.Claim.Max(), t.Claim.Min(), t.Name))
	}

	// check all areas, if the player is within 50 blocks of one, send pillars
	for _, a := range areas {
		pos0 := cube.Pos{
			int(a.Max()[0]),
			int(p.Position().Y()),
			int(a.Max()[1]),
		}

		pos1 := cube.Pos{
			int(a.Min()[0]),
			int(p.Position().Y()),
			int(a.Min()[1]),
		}
		pos2 := cube.Pos{pos0.X(), pos0.Y(), pos1.Z()}
		pos3 := cube.Pos{pos1.X(), pos0.Y(), pos0.Z()}
		h.SendAirPillar(pos0)
		h.SendAirPillar(pos1)
		h.SendAirPillar(pos2)
		h.SendAirPillar(pos3)
		playerPos := p.Position()
		distX0 := math.Abs(playerPos.X() - float64(pos0.X()))
		distZ0 := math.Abs(playerPos.Z() - float64(pos0.Z()))
		distX1 := math.Abs(playerPos.X() - float64(pos1.X()))
		distZ1 := math.Abs(playerPos.Z() - float64(pos1.Z()))

		// TODO: FIX THIS DISTANCE CALCULATION
		if distX0 < 30 || distZ0 < 30 || distX1 < 30 || distZ1 < 30 {
			h.SendClaimPillar(pos0)
			h.SendClaimPillar(pos1)
			h.SendClaimPillar(pos2)
			h.SendClaimPillar(pos3)
			least := math.Min(distX0, distZ0)
			least = math.Min(least, math.Min(distX1, distZ1))
			if a.Vec2WithinOrEqualFloor(mgl64.Vec2{float64(playerPos.X()), float64(playerPos.Z())}) {
				least = 0
			}
			moyai.Messagef(p, "command.team.map.display", a.Name(), int(least))
		}
	}
}

func (t TeamClearMap) Run(src cmd.Source, _ *cmd.Output) {
	p, _ := src.(*player.Player)
	h, _ := p.Handler().(*user.Handler)

	areas := make([]area.NamedArea, 0)
	teams, err := data.LoadAllTeams()
	if err != nil {
		moyai.Messagef(p, "command.team.load.fail")
		return
	}
	for _, t := range teams {
		areas = append(areas, area.NewNamedArea(t.Claim.Max(), t.Claim.Min(), t.Name))
	}

	// check all areas, if the player is within 50 blocks of one, send pillars
	for _, a := range areas {
		pos0 := cube.Pos{
			int(a.Max()[0]),
			int(p.Position().Y()),
			int(a.Max()[1]),
		}

		pos1 := cube.Pos{
			int(a.Min()[0]),
			int(p.Position().Y()),
			int(a.Min()[1]),
		}
		pos2 := cube.Pos{pos0.X(), pos0.Y(), pos1.Z()}
		pos3 := cube.Pos{pos1.X(), pos0.Y(), pos0.Z()}
		h.SendAirPillar(pos0)
		h.SendAirPillar(pos1)
		h.SendAirPillar(pos2)
		h.SendAirPillar(pos3)
	}
}

// Run executes the TeamCreate command.
func (t TeamCreate) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	if teamExists(p) {
		moyai.Messagef(p, "team.create.already")
		return
	}

	t.Name = colour.StripMinecraftColour(t.Name)
	if !validateTeamName(p, t.Name) {
		return
	}

	if u.Teams.Create.Active() {
		moyai.Messagef(p, "command.team.create.cooldown", timeutil.FormatDuration(u.Teams.Create.Remaining()))
		return
	}

	if _, err = data.LoadTeamFromName(t.Name); err == nil {
		moyai.Messagef(p, "team.create.exists")
		return
	}

	tm := createTeam(p, t.Name)
	u.Teams.Create.Set(time.Minute * 3)
	data.SaveUser(u)

	moyai.Messagef(p, "team.create.success", tm.DisplayName)
	moyai.Broadcastf("team.create.success.broadcast", p.Name(), tm.DisplayName)
}

func teamExists(p *player.Player) bool {
	_, err := data.LoadTeamFromMemberName(p.Name())
	return err == nil
}

func createTeam(p *player.Player, name string) data.Team {
	tm := data.DefaultTeam(name).WithMembers(data.DefaultMember(p.XUID(), p.Name()).WithRank(3))
	data.SaveTeam(tm)
	return tm
}

func validateTeamName(p *player.Player, name string) bool {
	name = strings.TrimSpace(name)

	if !regex.MatchString(name) {
		moyai.Messagef(p, "team.create.name.invalid")
		return false
	}
	if len(name) < 3 {
		moyai.Messagef(p, "team.create.name.short")
		return false
	}
	if len(name) > 15 {
		moyai.Messagef(p, "team.create.name.long")
		return false
	}
	return true
}

func (t TeamInvite) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	target, ok := t.Target[0].(*player.Player)
	if !ok || target == p {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	if tm.Frozen() {
		moyai.Messagef(p, "command.team.dtr")
		return
	}

	if len(tm.Members) == 7 {
		moyai.Messagef(p, "command.team.join.full")
		return
	}

	if tm.Member(target.Name()) {
		moyai.Messagef(p, "team.invite.member", target.Name())
		return
	}

	targetUser, err := data.LoadUserFromName(target.Name())
	if err != nil {
		return
	}

	_, err = data.LoadTeamFromMemberName(target.Name())
	if err == nil {
		moyai.Messagef(p, "team.invite.has-team")
		return
	}

	if targetUser.Teams.Invitations.Active(tm.Name) {
		moyai.Messagef(p, "team.invite.already", target.Name())
		return
	}

	targetUser.Teams.Invitations.Set(tm.Name, time.Minute*5)
	data.SaveUser(targetUser)

	team.Broadcastf(tm, "team.invite.success.broadcast", target.Name())
	moyai.Messagef(target, "team.invite.target", tm.DisplayName)
}

// Run ...
func (t TeamJoin) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	// Check if player is already in a team
	_, err := data.LoadTeamFromMemberName(p.Name())
	if err == nil {
		moyai.Messagef(p, "team.join.error")
		return
	}

	// Load the team to join
	tm, err := data.LoadTeamFromName(string(t.Team))
	if err != nil {
		moyai.Messagef(p, "command.team.not.found")
		return
	}

	// Check if the team is frozen or at full capacity
	if tm.Frozen() {
		moyai.Messagef(p, "command.team.dtr")
		return
	}

	if len(tm.Members) == 7 {
		moyai.Messagef(p, "command.team.join.full")
		return
	}

	// Load user data and reset any existing invitations
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	u.Teams.Invitations.Reset(tm.Name)
	data.SaveUser(u)

	// Add player to the team and update team DTR
	tm = tm.WithMembers(append(tm.Members, data.DefaultMember(p.XUID(), p.Name()))...)
	tm = tm.WithDTR(tm.DTR + 1.01)
	data.SaveTeam(tm)

	// Broadcast team join event
	team.Broadcastf(tm, "team.member.join", p.Name(), tm.DTR)
}

// Run ...
func (t TeamInformation) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	n, _ := t.Name.Load()
	name := strings.TrimSpace(string(n))

	var anyFound bool

	// Check if the specified name is empty
	if name == "" {
		tm, err := data.LoadTeamFromMemberName(p.Name())
		if err != nil {
			moyai.Messagef(p, "user.team-less")
			return
		}

		out.Print(teamInformationFormat(tm))
		return
	}

	// Attempt to load team by name or member name
	tm, err := data.LoadTeamFromName(name)
	if err == nil {
		out.Print(teamInformationFormat(tm))
		anyFound = true
	}

	tm, err = data.LoadTeamFromMemberName(name)
	if err == nil {
		out.Print(teamInformationFormat(tm))
		anyFound = true
	}

	// If no team was found, send a message to the player
	if !anyFound {
		moyai.Messagef(p, "command.team.info.not.found", name)
		return
	}
}

// Run ...
func (t TeamWho) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	n, _ := t.Name.Load()
	name := strings.TrimSpace(string(n))

	var anyFound bool

	// Check if the specified name is empty
	if name == "" {
		tm, err := data.LoadTeamFromMemberName(p.Name())
		if err != nil {
			moyai.Messagef(p, "user.team-less")
			return
		}
		out.Print(teamInformationFormat(tm))
		return
	}

	// Attempt to load team by name or member name
	tm, err := data.LoadTeamFromName(strings.ToLower(name))
	if err == nil {
		out.Print(teamInformationFormat(tm))
		anyFound = true
	}

	tm, err = data.LoadTeamFromMemberName(strings.ToLower(name))
	if err == nil {
		out.Print(teamInformationFormat(tm))
		anyFound = true
	}

	// If no team was found, send a message to the player
	if !anyFound {
		moyai.Messagef(p, "command.team.info.not.found", name)
		return
	}
}

func (t TeamRally) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	online := team.OnlineMembers(tm)
	playerPos := p.Position()

	for _, o := range online {
		if h, ok := o.Handler().(*user.Handler); ok {
			h.SetWayPoint(user.NewWayPoint("Rally", playerPos))
			moyai.Messagef(o, "command.team.rallying", p.Name(), int(playerPos.X()), int(playerPos.Y()), int(playerPos.Z()))
		}
	}
}

func (t TeamUnRally) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	online := team.OnlineMembers(tm)

	for _, o := range online {
		if h, ok := o.Handler().(*user.Handler); ok {
			h.RemoveWaypoint()
		}
	}
}

// Run ...
func (t TeamDisband) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	name := p.Name()
	tm, err := data.LoadTeamFromMemberName(name)
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	if !tm.Leader(name) {
		moyai.Messagef(p, "command.team.disband.must.leader")
		return
	}

	if tm.Frozen() {
		moyai.Messagef(p, "command.team.dtr")
		return
	}

	team.Broadcastf(tm, "command.team.disband.disbanded", p.Name())
	data.DisbandTeam(tm)
}

func (t TeamLeave) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	if tm.Leader(p.Name()) {
		moyai.Messagef(p, "command.team.leave.leader")
		return
	}

	if tm.Frozen() {
		moyai.Messagef(p, "command.team.dtr")
		return
	}

	tm = tm.WithoutMember(p.Name())
	tm = tm.WithDTR(tm.DTR - 1.01)
	team.Broadcastf(tm, "team.member.leave", p.Name(), tm.DTR)
	data.SaveTeam(tm)
}

func (t TeamKick) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	// Load the team of the player issuing the command
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	// Check if the player has permission to kick
	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		moyai.Messagef(p, "command.team.kick.missing.permission")
		return
	}

	// Check if the target member is the player themselves
	if strings.EqualFold(p.Name(), string(t.Member)) {
		moyai.Messagef(p, "command.team.kick.self")
		return
	}

	// Check if the target member is a leader or captain
	if tm.Leader(string(t.Member)) {
		moyai.Messagef(p, "command.team.kick.leader")
		return
	}
	if tm.Captain(p.Name()) && tm.Captain(string(t.Member)) {
		moyai.Messagef(p, "command.team.kick.captain")
		return
	}

	// Check if the target member exists in the team
	if !tm.Member(string(t.Member)) {
		moyai.Messagef(p, "command.team.kick.not.found", string(t.Member))
		return
	}

	// Check if the team is frozen
	if tm.Frozen() {
		moyai.Messagef(p, "command.team.dtr")
		return
	}

	// Kick the member from the team
	us, ok := user.Lookup(string(t.Member))
	if ok {
		moyai.Messagef(us, "command.team.kick.user.kicked", tm.DisplayName)
	}

	tm = tm.WithoutMember(string(t.Member))
	tm = tm.WithDTR(tm.DTR - 1.01)
	data.SaveTeam(tm)

	// Update state and notify kicked member
	for _, m := range tm.Members {
		if mem, ok := user.Lookup(m.Name); ok {
			user.UpdateState(mem)
			moyai.Messagef(mem, "command.team.kick.kicked")
		}
	}
}

func (t TeamPromote) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	// Load the team of the player issuing the command
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	// Check if the player has permission to promote
	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		moyai.Messagef(p, "command.team.promote.missing.permission")
		return
	}

	// Check if the target member is the player themselves
	if strings.EqualFold(p.Name(), string(t.Member)) {
		moyai.Messagef(p, "command.team.promote.self")
		return
	}

	// Check if the player is trying to promote a leader
	if tm.Leader(string(t.Member)) {
		moyai.Messagef(p, "command.team.promote.leader")
		return
	}

	// Check if the player is trying to promote a captain who is already a captain
	if tm.Captain(p.Name()) && tm.Captain(string(t.Member)) {
		moyai.Messagef(p, "command.team.promote.captain")
		return
	}

	// Check if the target member exists in the team
	if !tm.Member(string(t.Member)) {
		moyai.Messagef(p, "command.team.member.not.found", string(t.Member))
		return
	}

	// Promote the member and save the team
	tm = tm.Promote(string(t.Member))
	data.SaveTeam(tm)

	// Determine the rank name for broadcasting
	rankName := "Captain"
	if tm.Leader(string(t.Member)) {
		rankName = "Leader"
	}

	// Broadcast promotion message to the team
	team.Broadcastf(tm, "command.team.promote.user.promoted", string(t.Member), rankName)
}

func (t TeamDemote) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	// Load the team of the player issuing the command
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	// Check if the player has permission to demote
	if !tm.Leader(p.Name()) {
		moyai.Messagef(p, "command.team.demote.missing.permission")
		return
	}

	// Check if the target member is the player themselves
	if strings.EqualFold(string(t.Member), p.Name()) {
		moyai.Messagef(p, "command.team.demote.self")
		return
	}

	// Check if the target member is already a leader
	if !tm.Leader(string(t.Member)) {
		moyai.Messagef(p, "command.team.demote.leader")
		return
	}

	// Check if the target member exists in the team
	if !tm.Member(string(t.Member)) {
		moyai.Messagef(p, "command.team.member.not.found", string(t.Member))
		return
	}

	// Check if the target member is neither captain nor leader
	if !tm.Captain(string(t.Member)) && !tm.Leader(string(t.Member)) {
		moyai.Messagef(p, "command.team.demote.member")
		return
	}

	// Demote the member and save the team
	tm.Demote(string(t.Member))
	data.SaveTeam(tm)

	// Broadcast demotion message to the team
	team.Broadcastf(tm, "command.team.demote.user.demoted", string(t.Member), "Member")
}

// Run ...
func (t TeamTop) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	teams, err := data.LoadAllTeams()
	if err != nil {
		moyai.Messagef(p, "command.team.load.fail")
	}

	if len(teams) == 0 {
		moyai.Messagef(p, "command.team.top.none")
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
		if userTeam.Name == tm.Name {
			top += text.Colourf(" <grey>%d. <green>%s</green> (<yellow>%d</yellow>)</grey>\n", i+1, tm.DisplayName, tm.Points)
		} else {
			top += text.Colourf(" <grey>%d. <red>%s</red> (<yellow>%d</yellow>)</grey>\n", i+1, tm.DisplayName, tm.Points)
		}
	}
	top += "\uE000"
	p.Message(top)
}

// Run ...
func (t TeamClaim) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		moyai.Messagef(p, "command.team.promote.missing.permission")
		return
	}

	if tm.Claim != (area.Area{}) {
		moyai.Messagef(p, "team.has-claim")
		return
	}

	_, _ = p.Inventory().AddItem(item.NewStack(item.Hoe{Tier: item.ToolTierDiamond}, 1).WithValue("CLAIM_WAND", true).WithLore(
		text.Colourf("<green>1. <yellow>Right click one position</yellow></green>"),
		text.Colourf("<green>2. <yellow>Right click another position while shifting</yellow></green>"),
		text.Colourf("<green>3. <yellow>Left click the air to confirm the claim</yellow></green>"),
	))
}

// Run ...
func (t TeamUnClaim) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	if !tm.Leader(p.Name()) {
		moyai.Messagef(p, "team.not-leader")
		return
	}

	tm = tm.WithClaim(area.Area{}).WithHome(mgl64.Vec3{})
	data.SaveTeam(tm)
	moyai.Messagef(p, "command.unclaim.success")
}

// Run ...
func (t TeamSetHome) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	cl := tm.Claim
	if cl == (area.Area{}) {
		moyai.Messagef(p, "team.claim.none")
		return
	}

	if !cl.Vec3WithinOrEqualXZ(p.Position()) {
		moyai.Messagef(p, "team.claim.not-within")
		return
	}

	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		moyai.Messagef(p, "team.not-leader-or-captain")
		return
	}

	tm = tm.WithHome(p.Position())
	data.SaveTeam(tm)
	moyai.Messagef(p, "command.team.home.set")
}

// Run ...
func (t TeamHome) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	if h.Combat().Active() {
		moyai.Messagef(p, "user.teleport.combat")
		return
	}

	if h.Home().Ongoing() {
		moyai.Messagef(p, "user.already.teleporting")
		return
	}

	if tm.Home == (mgl64.Vec3{}) {
		moyai.Messagef(p, "command.team.home.none")
		return
	}

	if area.Spawn(p.World()).Vec3WithinOrEqualXZ(p.Position()) {
		if p.World() != moyai.Overworld() {
			moyai.Overworld().AddEntity(p)
		}
		p.Teleport(tm.Home)
		return
	}

	// Adjust teleport duration based on nearby claimed areas
	dur := time.Second * 10
	teams, err := data.LoadAllTeams()
	if err != nil {
		moyai.Messagef(p, "command.team.load.fail")
		return
	}

	for _, otherTeam := range teams {
		if otherTeam.Claim.Vec3WithinOrEqualXZ(p.Position()) && otherTeam.Name != tm.Name {
			dur = time.Second * 20
			break
		}
	}

	h.Home().Teleport(p, dur, tm.Home, moyai.Overworld())
}

// Run ...
func (t TeamList) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	teams, err := data.LoadAllTeams()
	if err != nil {
		moyai.Messagef(p, "command.team.load.fail")
		return
	}

	// Sort teams by online member count descending, then by DTR ascending
	sort.Slice(teams, func(i, j int) bool {
		if len(team.OnlineMembers(teams[i])) != len(team.OnlineMembers(teams[j])) {
			return len(team.OnlineMembers(teams[i])) > len(team.OnlineMembers(teams[j]))
		}
		return teams[i].DTR < teams[j].DTR
	})

	// Filter out teams with DTR <= 0
	filteredTeams := make([]data.Team, 0, len(teams))
	for _, tm := range teams {
		if tm.DTR > 0 {
			filteredTeams = append(filteredTeams, tm)
		}
	}

	var list strings.Builder
	list.WriteString(text.Colourf("        <yellow>Team List</yellow>\n"))
	list.WriteString("\uE000\n")

	userTeam, _ := data.LoadTeamFromMemberName(p.Name())

	for i, tm := range filteredTeams {
		if i > 9 {
			break
		}

		onlineCount := len(team.OnlineMembers(tm))
		dtr := tm.DTRString()

		var lineFormat string
		if userTeam.Name != "" && userTeam.Name == tm.Name {
			lineFormat = " <grey>%d. <green>%src</green> (<green>%d/%d</green>) %s</grey>\n"
		} else {
			lineFormat = " <grey>%d. <red>%src</red> (<green>%d/%d</green>) %s</grey>\n"
		}

		list.WriteString(text.Colourf(lineFormat, i+1, tm.DisplayName, onlineCount, len(tm.Members), dtr))
	}

	list.WriteString("\uE000")
	p.Message(list.String())
}

// Run ...
func (t TeamFocusTeam) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}
	targetTeam, err := data.LoadTeamFromName(strings.ToLower(string(t.Name)))
	if err != nil {
		moyai.Messagef(p, "command.team.not.found", string(t.Name))
		return
	}

	if tm.Name == targetTeam.Name {
		moyai.Messagef(p, "command.team.focus.self")
		return
	}
	tm = tm.WithTeamFocus(targetTeam)
	data.SaveTeam(tm)

	for _, m := range team.OnlineMembers(targetTeam) {
		user.UpdateState(m)
		// if targetTeam.Home != (mgl64.Vec3{}) {
		// 	if h, ok := m.Handler().(*user.Handler); ok {
		// 		h.SetWayPoint(user.NewWayPoint(targetTeam.DisplayName, targetTeam.Home))
		// 	}
		// }
	}

	team.Broadcastf(tm, "command.team.focus", targetTeam.DisplayName)
}

// Run ...
func (t TeamUnFocus) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}
	focus := tm.Focus

	if focus.Kind == data.FocusTypeNone {
		moyai.Messagef(p, "command.team.focus.none")
		return
	}

	tm = tm.WithoutFocus()
	data.SaveTeam(tm)

	if focus.Kind == data.FocusTypeTeam {
		targetTeam, err := data.LoadTeamFromName(focus.Value)
		if err == nil {
			for _, m := range team.OnlineMembers(targetTeam) {
				user.UpdateState(m)
				// if h, ok := m.Handler().(*user.Handler); ok {
				// 	h.RemoveWaypoint()
				// }
			}
		}
	} else if focus.Kind == data.FocusTypePlayer {
		if mem, ok := user.Lookup(focus.Value); ok {
			user.UpdateState(mem)
		}
	}

	tm = tm.WithoutFocus()
	data.SaveTeam(tm)

	for _, m := range team.OnlineMembers(tm) {
		user.UpdateState(m)
	}

	team.Broadcastf(tm, "command.team.unfocus", focus.Value)
}

// Run ...
func (t TeamFocusPlayer) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}
	if len(t.Target) > 1 {
		moyai.Messagef(p, "command.targets.exceed")
		return
	}
	target, ok := t.Target[0].(*player.Player)
	if !ok {
		moyai.Messagef(p, "command.team.focus.player")
		return
	}

	if strings.EqualFold(target.Name(), p.Name()) {
		moyai.Messagef(p, "command.team.focus.yourself")
		return
	}

	if tm.Member(target.Name()) {
		moyai.Messagef(p, "command.team.focus.member")
		return
	}

	tm = tm.WithPlayerFocus(target.Name())
	data.SaveTeam(tm)
	user.UpdateState(target)
	team.Broadcastf(tm, "command.team.focus", target.Name())
}

// Run ...
func (t TeamChat) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
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
		p.Message(lang.Translatef(*u.Language, "command.team.chat.global"))
	case 0:
		u.Teams.ChatType = 1
		p.Message(lang.Translatef(*u.Language, "command.team.chat.team"))
	}
	data.SaveUser(u)
}

// Run ...
func (t TeamWithdraw) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		moyai.Messagef(p, "command.team.withdraw.permission")
		return
	}

	amt := t.Balance
	if amt < 1 {
		moyai.Messagef(p, "command.team.withdraw.minimum")
		return
	}

	if amt > tm.Balance {
		moyai.Messagef(p, "command.team.withdraw.insufficient", amt)
		return
	}

	tm = tm.WithBalance(tm.Balance - amt)
	u.Teams.Balance = u.Teams.Balance + amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	moyai.Messagef(p, "command.team.withdraw.success", int(amt), tm.DisplayName)
}

func (t TeamW) Run(s cmd.Source, o *cmd.Output) {
	TeamWithdraw{
		Balance: t.Balance,
	}.Run(s, o)
}

// Run ...
func (t TeamDeposit) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}

	amt := t.Balance
	if amt < 1 {
		moyai.Messagef(p, "command.team.deposit.minimum")
		return
	}

	if amt > u.Teams.Balance {
		moyai.Messagef(p, "command.team.deposit.insufficient", amt)
		return
	}

	tm = tm.WithBalance(tm.Balance + amt)
	u.Teams.Balance = u.Teams.Balance - amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	moyai.Messagef(p, "command.team.deposit.success", int(amt), tm.DisplayName)
}

func (t TeamD) Run(s cmd.Source, o *cmd.Output) {
	TeamDeposit{
		Balance: t.Balance,
	}.Run(s, o)
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
		moyai.Messagef(p, "user.team-less")
		return
	}

	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		moyai.Messagef(p, "command.team.withdraw.permission")
		return
	}

	amt := tm.Balance
	if amt < 1 {
		moyai.Messagef(p, "command.team.withdraw.minimum")
		return
	}

	tm = tm.WithBalance(tm.Balance - amt)
	u.Teams.Balance = u.Teams.Balance + amt

	data.SaveTeam(tm)
	data.SaveUser(u)

	moyai.Messagef(p, "command.team.withdraw.success", int(amt), tm.Name)
}

func (t TeamWAll) Run(s cmd.Source, o *cmd.Output) {
	TeamWithdrawAll{
		All: t.All,
	}.Run(s, o)
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
		moyai.Messagef(p, "user.team-less")
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

func (t TeamDAll) Run(s cmd.Source, o *cmd.Output) {
	TeamDepositAll{
		All: t.All,
	}.Run(s, o)
}

// Run ...
func (t TeamStuck) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	if p.World() != moyai.Overworld() {
		return
	}
	pos := safePosition(p, cube.PosFromVec3(p.Position()), 24)
	if pos == (cube.Pos{}) {
		moyai.Messagef(p, "command.team.stuck.no-safe")
		return
	}

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	if h.Stuck().Ongoing() {
		o.Error("You are already in the stuck process.")
		return
	}

	moyai.Messagef(p, "command.team.stuck.teleporting", pos.X(), pos.Y(), pos.Z(), 30)
	h.Stuck().Teleport(p, time.Second*30, mgl64.Vec3{
		float64(pos.X()),
		float64(pos.Y()),
		float64(pos.Z()),
	}, moyai.Overworld())
}

// Run ...
func (t TeamCamp) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	if p.World() != moyai.Overworld() {
		return
	}

	tm, err := data.LoadTeamFromName(string(t.Team))
	if err != nil {
		moyai.Messagef(p, "command.team.not.found", t.Team)
		return
	}

	if tm.Home == (mgl64.Vec3{}) {
		moyai.Messagef(p, "command.team.homeless", tm.DisplayName)
		return
	}
	pos := safePosition(p, cube.PosFromVec3(tm.Home), 50)

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	if h.CampOngoing() {
		o.Error("You are already in the camp process.")
		return
	}

	h.BeginCamp(tm, pos)
}

func safePosition(p *player.Player, pos cube.Pos, radius int) cube.Pos {
	w := p.World()

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
							y := w.Range().Max()
							for y > pos.Y() {
								y--
								b := w.Block(cube.Pos{x, y, z})
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

// Run ...
func (t TeamIncrementDTR) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromName(string(t.Name))
	if err != nil {
		moyai.Messagef(p, "command.team.not.found", t.Name)
		return
	}
	tm = tm.WithDTR(tm.DTR + 1.00)
	data.SaveTeam(tm)

	p.Message(text.Colourf("<green>Successfully incremented DTR by 1.00.</green>"))
}

// Run ...
func (t TeamDecrementDTR) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromName(string(t.Name))
	if err != nil {
		moyai.Messagef(p, "command.team.not.found", t.Name)
		return
	}
	tm = tm.WithDTR(tm.DTR - 1.00)
	tm = tm.WithLastDeath(time.Now())
	data.SaveTeam(tm)

	p.Message(text.Colourf("<green>Successfully decremented DTR by 1.00.</green>"))
}

// Run ...
func (t TeamResetRegen) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromName(string(t.Name))
	if err != nil {
		moyai.Messagef(p, "command.team.not.found", t.Name)
		return
	}
	tm = tm.WithLastDeath(time.Time{})
	data.SaveTeam(tm)

	p.Message(text.Colourf("<green>Successfully reset team regeneration.</green>"))
}

// Run ...
func (t TeamRename) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.team-less")
		return
	}
	if !tm.Leader(p.Name()) && !tm.Captain(p.Name()) {
		moyai.Messagef(p, "command.team.promote.missing.permission")
		return
	}
	if tm.Renamed {
		moyai.Messagef(p, "team.rename.already")
		return
	}

	t.Name = colour.StripMinecraftColour(t.Name)
	if !validateTeamName(p, t.Name) {
		return
	}
	tm = tm.WithRename(t.Name)

	for _, m := range tm.Members {
		if mem, ok := user.Lookup(m.Name); ok {
			user.UpdateState(mem)
		}
	}

	moyai.Messagef(p, "team.rename.success", tm.DisplayName)
	team.Broadcastf(tm, "team.rename.success.broadcast", p.Name(), tm.DisplayName)
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
			members = append(members, m.Name)
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
		teams = append(teams, tm.Name)
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
			members = append(members, m.Name)
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
	regenerationTime := t.LastDeath.Add(time.Minute * 15)
	if time.Now().Before(regenerationTime) {
		formattedRegenerationTime = text.Colourf("\n <yellow>Time Until Regen</yellow> <blue>%s</blue>", time.Until(regenerationTime).Round(time.Second))
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
			"<yellow>Points: </yellow><red>%d</red>\n "+
			"<yellow>Koth Captures: <red>%d</red>\n "+
			"<yellow>Deaths Until Raidable: </yellow>%s%s\n\uE000", t.DisplayName, onlineCount, len(t.Members), home, formattedLeader, strings.Join(formattedCaptains, ", "), strings.Join(formattedMembers, ", "), t.Balance, t.Points, t.KOTHWins, formattedDtr, formattedRegenerationTime)
}

func formatMember(name string, kills int, online bool) string {
	if online {
		return text.Colourf("<green>%s</green><dark-red>[%d]</dark-red>", name, kills)
	}
	return text.Colourf("<grey>%s</grey><dark-red>[%d]</dark-red>", name, kills)
}

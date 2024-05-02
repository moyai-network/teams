package data

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal/area"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	teamCollection *mongo.Collection

	teamMu sync.Mutex
	teams  = map[string]Team{}
)

func teamCached(f func(Team) bool) (Team, bool) {
	teamMu.Lock()
	defer teamMu.Unlock()
	for _, t := range teams {
		if f(t) {
			return t, true
		}
	}
	return Team{}, false
}

func saveTeamData(t Team) {
	filter := bson.M{"name": bson.M{"$eq": t.Name}}
	update := bson.M{"$set": t}

	res, _ := teamCollection.UpdateOne(ctx(), filter, update)

	if res.MatchedCount == 0 {
		_, _ = teamCollection.InsertOne(ctx(), t)
	}
}

func SaveTeam(t Team) {
	teamMu.Lock()
	teams[t.Name] = t
	teamMu.Unlock()
}

type Team struct {
	// Name is the identifier for the team data.
	Name string
	// DisplayName is the display name for the team.
	DisplayName string
	// Members is a slice of all the members in the team.
	Members []Member
	// DTR is the amount of deaths required before the teams goes raidable.
	DTR float64
	// Home is the home point for the team.
	Home mgl64.Vec3
	// RegenerationTime is the time until the team can start regenerating their DTR.
	RegenerationTime time.Time
	// Points is the amount of points the team has.
	Points int
	// Balance is the amount of money the team has.
	Balance float64
	// Claim is the claim area of the team.
	Claim area.Area
	// Focus is the focus information for a team
	Focus Focus
}

// DefaultTeam returns a team with default values
func DefaultTeam(name string) Team {
	return Team{
		Name:        strings.ToLower(name),
		DTR:         1.01,
		DisplayName: name,
	}
}

// WithMembers returns the team with the given members.
func (t Team) WithMembers(m ...Member) Team {
	t.Members = m
	return t
}

// WithoutMember returns the team without the given member
func (t Team) WithoutMember(name string) Team {
	for i, m := range t.Members {
		if strings.EqualFold(name, m.Name) {
			t.Members = slices.Delete(t.Members, i, i+1)
		}
	}
	return t
}

// Promote promotes a member of the faction.
func (t Team) Promote(name string) Team {
	var m Member
	var l Member
	for _, me := range t.Members {
		if me.Name == name {
			m = me
		}
		if t.Leader(me.Name) {
			l = me
		}
	}
	switch m.Rank {
	case 1:
		m.Rank = 2
	case 2:
		l.Rank = 2
		m.Rank = 3
	}

	return t
}

// Demote demotes a member of the faction.
func (t Team) Demote(name string) Team {
	var m Member
	for _, me := range t.Members {
		if me.Name == name {
			m = me
		}
	}
	if m.Rank == 2 {
		m.Rank = 1
	}

	return t
}

// WithClaim returns the team with the given claim area.
func (t Team) WithClaim(claim area.Area) Team {
	t.Claim = claim
	return t
}

// WithPoints returns the team with the given amount of points
func (t Team) WithPoints(points int) Team {
	t.Points = points
	return t
}

// WithBalance returns the team with the given balance.
func (t Team) WithBalance(bal float64) Team {
	t.Balance = bal
	return t
}

// WithHome returns the team with the given home coordinate.
func (t Team) WithHome(home mgl64.Vec3) Team {
	t.Home = home
	return t
}

// WithRegenerationTime returns the team with the given regeneration time.
func (t Team) WithRegenerationTime(regen time.Time) Team {
	t.RegenerationTime = regen
	return t
}

// Frozen returns whether the team is frozen.
func (t Team) Frozen() bool {
	return time.Now().Before(t.RegenerationTime)
}

// WithDTR returns the team with the given dtr.
func (t Team) WithDTR(dtr float64) Team {
	t.DTR = dtr
	return t
}

// MaxDTR returns the max DTR of the faction.
func (t Team) MaxDTR() float64 {
	dtr := 1.01 * float64(len(t.Members))
	return math.Round(dtr*100) / 100
}

// DTRString returns the DTR string of the faction
func (t Team) DTRString() string {
	if eq(t.DTR, t.MaxDTR()) {
		return text.Colourf("<green>%.2f%s</green>", t.DTR, t.DTRDot())
	}
	if t.DTR < 0 {
		return text.Colourf("<redstone>%.2f%s</redstone>", t.DTR, t.DTRDot())
	}
	return text.Colourf("<yellow>%.2f%s</yellow>", t.DTR, t.DTRDot())
}

// DTRDot returns the DTR dot of the faction.
func (t Team) DTRDot() string {
	if eq(t.DTR, t.MaxDTR()) {
		return "<green>■</green>"
	}
	if t.DTR < 0 {
		return "<redstone>■</redstone>"
	}
	return "<yellow>■</yellow>"
}

// Leader returns whether the given username is the one of the leader.
func (t Team) Leader(name string) bool {
	for _, m := range t.Members {
		if strings.EqualFold(m.Name, name) && m.Rank == 3 {
			return true
		}
	}
	return false
}

// Captain returns whether the given username is the one of the captain.
func (t Team) Captain(name string) bool {
	for _, m := range t.Members {
		if strings.EqualFold(m.Name, name) && m.Rank == 2 {
			return true
		}
	}
	return false
}

// Member returns whether the given username is in the team.
func (t Team) Member(name string) bool {
	for _, m := range t.Members {
		if strings.EqualFold(m.Name, name) {
			return true
		}
	}
	return false
}

// WithTeamFocus returns the team with the given team as the focus.
func (t Team) WithTeamFocus(tm Team) Team {
	t.Focus.focusType = FocusTypeTeam()
	t.Focus.value = tm.Name
	return t
}

// WithPlayerFocus returns the team with the given player as the focus.
func (t Team) WithPlayerFocus(name string) Team {
	t.Focus.focusType = FocusTypePlayer()
	t.Focus.value = name
	return t
}

// WithoutFocus returns the team without a focus.
func (t Team) WithoutFocus() Team {
	t.Focus.focusType = FocusTypeNone()
	t.Focus.value = ""
	return t
}

// Member represents a team member.
type Member struct {
	Name        string
	DisplayName string
	XUID        string
	Rank        int
}

// DefaultMember returns a default team member.
func DefaultMember(xuid, name string) Member {
	return Member{
		Name:        strings.ToLower(name),
		DisplayName: name,
		XUID:        xuid,
		Rank:        1,
	}
}

// WithRank returns a team member with the given rank.
func (m Member) WithRank(n int) Member {
	m.Rank = n
	return m
}

func decodeSingleTeamFromFilter(filter any) (Team, error) {
	return decodeSingleTeamResult(userCollection.FindOne(ctx(), filter))
}

func decodeSingleTeamResult(result *mongo.SingleResult) (Team, error) {
	var t Team

	err := result.Decode(&t)
	if err != nil {
		return Team{}, err
	}

	teamMu.Lock()
	teams[t.Name] = t
	teamMu.Unlock()

	return updatedRegeneration(t), nil
}

func LoadAllTeams() ([]Team, error) {
	var tms []Team

	result, err := teamCollection.Find(ctx(), bson.M{})
	err = result.Decode(&tms)
	if err != nil {
		return []Team{}, err
	}
	var mappedTeams map[string]Team

	for _, t := range tms {
		mappedTeams[t.Name] = updatedRegeneration(t)
	}

	for _, t := range mappedTeams {
		if tm, ok := teamCached(func(team Team) bool {
			return t.Name == team.Name
		}); ok {
			mappedTeams[t.Name] = updatedRegeneration(tm)
		}
	}

	return tms, nil
}

// LoadTeamFromName loads a team using the given name.
func LoadTeamFromName(name string) (Team, error) {
	name = strings.ToLower(name)

	if t, ok := teamCached(func(t Team) bool {
		return t.Name == name
	}); ok {
		return updatedRegeneration(t), nil
	}

	return decodeSingleTeamFromFilter(bson.M{"name": bson.M{"$eq": name}})
}

// LoadTeamFromMemberName loads a team using the given member name.
func LoadTeamFromMemberName(name string) (Team, error) {
	if t, ok := teamCached(func(t Team) bool {
		for _, m := range t.Members {
			if strings.EqualFold(m.Name, name) {
				return true
			}
		}
		return false
	}); ok {
		return updatedRegeneration(t), nil
	}

	return decodeSingleTeamFromFilter(bson.M{"members.name": bson.M{"$eq": name}})
}

func updatedRegeneration(t Team) Team {
	if t.RegenerationTime.Before(time.Now()) && t.DTR < t.MaxDTR() {
		t = t.WithDTR(t.MaxDTR())
	}
	return t
}

// DisbandTeam disbands the given team.
func DisbandTeam(t Team) {
	teamMu.Lock()
	delete(teams, t.Name)
	teamMu.Unlock()

	filter := bson.M{"name": bson.M{"$eq": t.Name}}
	_, _ = teamCollection.DeleteOne(ctx(), filter)
}

// Focus is the focus information for a team
type Focus struct {
	focusType FocusType // 0:Player ; 1: Team
	value     string    // XUID: Player ; Name: Team
}

// FocusType is the type of focus.
type FocusType struct {
	n int
}

// FocusTypeNone is a type for when not focusing on anyone.
func FocusTypeNone() FocusType {
	return FocusType{0}
}

// FocusTypePlayer is a type for when focusing one specific player.
func FocusTypePlayer() FocusType {
	return FocusType{1}
}

// FocusTypeTeam is a type for when focusing another team.
func FocusTypeTeam() FocusType {
	return FocusType{2}
}

// Value returns the string value associated with the focus.
func (f *Focus) Value() string {
	return f.value
}

// Type returns the type of focus.
func (f *Focus) Type() FocusType {
	return f.focusType
}

// focusData is a struct used for encoding/decoding focus data.
type focusData struct {
	Kind  int
	Value string
}

// UnmarshalBSON ...
func (f *Focus) UnmarshalBSON(b []byte) error {
	var d focusData
	err := bson.Unmarshal(b, &d)
	switch d.Kind {
	case 0:
		f.focusType = FocusTypeNone()
	case 1:
		f.focusType = FocusTypePlayer()
	case 2:
		f.focusType = FocusTypeTeam()
	}
	f.value = d.Value
	return err
}

// MarshalBSON ...
func (f *Focus) MarshalBSON() ([]byte, error) {
	d := focusData{
		Kind:  f.focusType.n,
		Value: f.value,
	}
	return bson.Marshal(d)
}

func eq(a, b float64) bool {
	return math.Abs(a-b) <= 1e-5
}

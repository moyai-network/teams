package data

import (
	"context"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	teamCollection *mongo.Collection
	teamsMu        sync.Mutex
	teams          = map[string]Team{}
)

func init() {
	var tms []Team
	cursor, err := db.Collection("teams").Find(ctx(), bson.M{})
	if err != nil {
		panic(err)
	}

	err = cursor.All(context.Background(), &tms)
	if err != nil {
		panic(err)
	}

	teamsMu.Lock()
	for _, t := range tms {
		teams[t.Name] = t
	}
	teamsMu.Unlock()
}

func Teams() []Team {
	teamsMu.Lock()
	defer teamsMu.Unlock()

	return maps.Values(teams)
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
	Claim moose.Area
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
func (t Team) WithoutMember(m Member) Team {
	if i := slices.Index(t.Members, m); i != -1 {
		slices.Delete(t.Members, i, i+1)
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
func (t Team) WithClaim(claim moose.Area) Team {
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
	dtr := 1.1 * float64(len(t.Members))
	return math.Round(dtr*100) / 100
}

// DTRString returns the DTR string of the faction
func (t Team) DTRString() string {
	if t.DTR == t.MaxDTR() {
		return text.Colourf("<green>%.2f%s</green>", t.DTR, t.DTRDot())
	}
	if t.DTR < 0 {
		return text.Colourf("<redstone>%.2f%s</redstone>", t.DTR, t.DTRDot())
	}
	return text.Colourf("<yellow>%.2f%s</yellow>", t.DTR, t.DTRDot())
}

// DTRDot returns the DTR dot of the faction.
func (t Team) DTRDot() string {
	if t.DTR == t.MaxDTR() {
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

// LoadTeam loads a team using the given name.
func LoadTeam(name string) (Team, bool) {
	teamsMu.Lock()
	t, ok := teams[name]
	teamsMu.Unlock()
	return t, ok
}

// DisbandTeam disbands the given team.
func DisbandTeam(t Team) {
	teamsMu.Lock()
	delete(teams, t.Name)
	teamsMu.Unlock()
}

// SaveTeam saves the given team.
func SaveTeam(t Team) {
	teamsMu.Lock()
	teams[t.Name] = t
	teamsMu.Unlock()
}

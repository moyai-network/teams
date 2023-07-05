package data

import (
	"context"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"strings"
	"sync"
	"time"
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

// Focus is the focus information for a team
type Focus struct {
	Kind  int    // 0:Player ; 1: Team
	Value string // XUID: Player ; Name: Team
}

func DefaultTeam(name string) Team {
	return Team{
		Name:        strings.ToLower(name),
		DisplayName: name,
		Focus:       Focus{Kind: -1},
	}
}

func (t Team) WithMembers(m ...Member) Team {
	t.Members = m
	return t
}

func (t Team) WithoutMember(m Member) Team {
	if i := slices.Index(t.Members, m); i != -1 {
		slices.Delete(t.Members, i, i+1)
	}
	return t
}

func (t Team) WithClaim(claim moose.Area) Team {
	t.Claim = claim
	return t
}

func (t Team) WithPoints(points int) Team {
	t.Points = points
	return t
}

func (t Team) WithBalance(bal float64) Team {
	t.Balance = bal
	return t
}

func (t Team) WithHome(home mgl64.Vec3) Team {
	t.Home = home
	return t
}

func (t Team) WithRegenerationTime(regen time.Time) Team {
	t.RegenerationTime = regen
	return t
}

func (t Team) WithDTR(dtr float64) Team {
	t.DTR = dtr
	return t
}

type Member struct {
	Name        string
	DisplayName string
	XUID        string
	Rank        int
}

func DefaultMember(xuid, name string) Member {
	return Member{
		Name:        strings.ToLower(name),
		DisplayName: name,
		XUID:        xuid,
		Rank:        1,
	}
}

func (m Member) WithRank(n int) Member {
	m.Rank = n
	return m
}

func LoadTeam(name string) (Team, bool) {
	teamsMu.Lock()
	t, ok := teams[name]
	teamsMu.Unlock()
	return t, ok
}

func DisbandTeam(t Team) {
	teamsMu.Lock()
	delete(teams, t.Name)
	teamsMu.Unlock()
}

func SaveTeam(t Team) {
	teamsMu.Lock()
	teams[t.Name] = t
	teamsMu.Unlock()
}

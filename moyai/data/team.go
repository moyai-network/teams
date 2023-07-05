package data

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
	"strings"
	"time"
)

var (
	teamCollection *mongo.Collection
)

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

func LoadUserTeam(name string) (Team, bool) {
	filter := bson.M{"members.name": strings.ToLower(name)}
	result := teamCollection.FindOne(ctx(), filter)
	if err := result.Err(); err != nil {
		return Team{}, false
	}
	var data Team
	err := result.Decode(&data)
	if err != nil {
		return Team{}, false
	}
	return data, true
}

func TeamExists(name string) bool {
	filter := bson.M{"name": bson.M{"$eq": strings.ToLower(name)}}
	result := teamCollection.FindOne(ctx(), filter)
	if err := result.Err(); err != nil {
		return false
	}
	var data Team
	err := result.Decode(&data)
	if err != nil {
		return false
	}
	return true
}

func LoadTeam(name string) (Team, error) {
	filter := bson.M{"name": bson.M{"$eq": strings.ToLower(name)}}
	result := teamCollection.FindOne(ctx(), filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return DefaultTeam(name), nil
		}
		return Team{}, err
	}
	var data Team
	err := result.Decode(&data)
	if err != nil {
		return Team{}, err
	}
	return data, nil
}

func SaveTeam(t Team) error {
	filter := bson.M{"name": bson.M{"$eq": t.Name}}
	update := bson.M{"$set": t}

	res, err := teamCollection.UpdateOne(ctx(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		_, err = teamCollection.InsertOne(ctx(), t)
		return err
	}
	return nil
}

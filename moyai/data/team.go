package data

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
	"strings"
)

var (
	teamCollection *mongo.Collection
)

type Team struct {
	// Name is the identifier for the team data.
	Name string
	// DisplayName is the display name for the team.
	DisplayName string
	// Balance is the amount of money the team has.
	Balance float64
	// Members is a slice of all the members in the team.
	Members []Member
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

type Member struct {
	Name        string
	DisplayName string
	Rank        int
}

func DefaultMember(name string) Member {
	return Member{
		Name:        strings.ToLower(name),
		DisplayName: name,
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

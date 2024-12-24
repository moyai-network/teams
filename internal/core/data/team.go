package data

import (
	"github.com/moyai-network/teams/internal/ports/model"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/maps"
)

var (
	teamCollection *mongo.Collection

	teamMu sync.Mutex
	teams  = map[string]model.Team{}
)

func init() {
	var tms []model.Team

	result, err := teamCollection.Find(ctx(), bson.M{})
	if err != nil {
		panic(err)
	}
	err = result.All(ctx(), &tms)
	if err != nil {
		panic(err)
	}

	for _, t := range tms {
		teams[t.Name] = updatedRegeneration(t)
	}
}

func teamCached(f func(model.Team) bool) (model.Team, bool) {
	teamMu.Lock()
	tms := teams
	teamMu.Unlock()
	for _, t := range tms {
		if f(t) {
			return t, true
		}
	}
	return model.Team{}, false
}

func saveTeamData(t model.Team) error {
	filter := bson.M{"name": bson.M{"$eq": t.Name}}
	update := bson.M{"$set": t}

	res, err := teamCollection.UpdateOne(ctx(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		_, _ = teamCollection.InsertOne(ctx(), t)
	}
	return nil
}

func SaveTeam(t model.Team) {
	teamMu.Lock()
	teams[t.Name] = t
	teamMu.Unlock()

	go func() {
		err := saveTeamData(t)
		if err != nil {
			log.Println("Error saving team data:", err)
		}
	}()
}

func decodeSingleTeamFromFilter(filter any) (model.Team, error) {
	return decodeSingleTeamResult(teamCollection.FindOne(ctx(), filter))
}

func decodeSingleTeamResult(result *mongo.SingleResult) (model.Team, error) {
	var t model.Team

	err := result.Decode(&t)
	if err != nil {
		return model.Team{}, err
	}

	teamMu.Lock()
	teams[t.Name] = t
	teamMu.Unlock()

	return updatedRegeneration(t), nil
}

func LoadAllTeams() ([]model.Team, error) {
	return maps.Values(teams), nil
}

// LoadTeamFromName loads a team using the given name.
func LoadTeamFromName(name string) (model.Team, error) {
	name = strings.ToLower(name)

	if t, ok := teamCached(func(t model.Team) bool {
		return t.Name == name
	}); ok {
		return updatedRegeneration(t), nil
	}

	// TEAMS ARE NOW ALWAYS CACHED
	//return decodeSingleTeamFromFilter(bson.M{"name": bson.M{"$eq": name}})
	return model.Team{}, mongo.ErrNoDocuments
}

// LoadTeamFromMemberName loads a team using the given member name.
func LoadTeamFromMemberName(name string) (model.Team, error) {
	name = strings.ToLower(name)
	if t, ok := teamCached(func(t model.Team) bool {
		for _, m := range t.Members {
			if name == m.Name {
				return true
			}
		}
		return false
	}); ok {
		return updatedRegeneration(t), nil
	}

	// TEAMS ARE NOW ALWAYS CACHED
	//return decodeSingleTeamFromFilter(bson.M{"members.name": bson.M{"$eq": name}})
	return model.Team{}, mongo.ErrNoDocuments
}

// updatedRegeneration returns the team with the regeneration updated.
func updatedRegeneration(t model.Team) model.Team {
	since := time.Since(t.LastDeath)

	if eq(t.DTR, t.MaxDTR()) {
		return t
	}

	if since <= time.Minute*15 {
		return t
	}
	since = since - time.Minute*15

	prog := float64(since-time.Minute*2) / float64(time.Minute*3)
	t.DTR = t.DTR + prog
	if t.DTR > t.MaxDTR() {
		t.DTR = t.MaxDTR()
	}
	return t
}

// DisbandTeam disbands the given team.
func DisbandTeam(t model.Team) {
	teamMu.Lock()
	delete(teams, t.Name)
	teamMu.Unlock()

	filter := bson.M{"name": bson.M{"$eq": t.Name}}
	_, _ = teamCollection.DeleteOne(ctx(), filter)
}

func eq(a, b float64) bool {
	return math.Abs(a-b) <= 1e-5
}

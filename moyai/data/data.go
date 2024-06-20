package data

import (
	"context"
	"github.com/moyai-network/teams/internal/cooldown"
	"github.com/moyai-network/teams/internal/sets"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FlushCache() {
	userMu.Lock()
	defer userMu.Unlock()
	for _, u := range users {
		if time.Since(u.lastSaved) > time.Minute*2 {
			delete(users, u.XUID)
		}
	}
	teamMu.Lock()
	defer teamMu.Unlock()
	for _, t := range teams {
		err := saveTeamData(t)
		if err != nil {
			log.Println("Error saving team data:", err)
			return
		}
	}
}

func Reset() {
	userMu.Lock()
	defer userMu.Unlock()
	users = map[string]User{}
	teamMu.Lock()
	defer teamMu.Unlock()
	teams = map[string]Team{}

	// Reset the database.
	_, err := userCollection.DeleteMany(ctx(), bson.M{})
	if err != nil {
		log.Println("Error deleting users:", err)
	}
	_, err = teamCollection.DeleteMany(ctx(), bson.M{})
	if err != nil {
		log.Println("Error deleting teams:", err)
	}
}

func PartialReset() {
	userMu.Lock()
	defer userMu.Unlock()
	users = map[string]User{}
	teamMu.Lock()
	defer teamMu.Unlock()
	teams = map[string]Team{}

	// Reset the database.
	usrs, err := LoadAllUsers()
	if err != nil {
		log.Println("Error loading users:", err)
	}
	for _, u := range usrs {
		u.PlayTime = 0
		u.Vanished = false
		u.StaffMode = false
		u.Frozen = false
		u.Teams.Lives = 0
		u.Teams.ChatType = 0
		u.Teams.Dead = false
		u.Teams.Reclaimed = false
		u.Teams.SOTW = true
		u.Teams.Invitations = cooldown.NewMappedCoolDown[string]()
		u.Teams.Kits = cooldown.NewMappedCoolDown[string]()
		u.Teams.DeathBan = cooldown.NewCoolDown()
		u.Teams.Report = cooldown.NewCoolDown()
		u.Teams.Refill = cooldown.NewCoolDown()
		u.Teams.PVP = cooldown.NewCoolDown()
		u.Teams.PVP.Set(time.Hour + time.Second)
		u.Teams.Create = cooldown.NewCoolDown()
		u.Teams.Stats = Stats{}
		u.Teams.ClaimedRewards = sets.New[int]()
		u.Teams.DeathInventory = &Inventory{}

		err = saveUserData(u)
		if err != nil {
			log.Println("Error saving user data:", err)

		}
	}

	_, err = teamCollection.DeleteMany(ctx(), bson.M{})
	if err != nil {
		log.Println("Error deleting teams:", err)
	}
}

// ctx returns a context.Context.
func ctx() context.Context {
	return context.Background()
}

// db is the Upper database session.
var db *mongo.Database

// init creates the Upper database connection.
func init() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017").SetServerAPIOptions(serverAPI))
	if err != nil {
		panic(err)
	}
	db = client.Database("teams")

	userCollection = db.Collection("users")
	teamCollection = db.Collection("teams")
}

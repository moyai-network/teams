package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FlushCache() {
	userMu.Lock()
	defer userMu.Unlock()
	for _, u := range users {
		err := saveUserData(u)
		if err != nil {
			log.Println("Error saving user data:", err)
			return
		}
		delete(users, u.XUID)
	}
	teamMu.Lock()
	defer teamMu.Unlock()
	for _, t := range teams {
		t.ConquestPoints = 0
		err := saveTeamData(t)
		if err != nil {
			log.Println("Error saving team data:", err)
			return
		}
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

package data

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://hi:iKrlmjwZ87eWrqZM@cluster0.dkysi.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI))
	if err != nil {
		panic(err)
	}
	db = client.Database("teams")

	userCollection = db.Collection("users")
}

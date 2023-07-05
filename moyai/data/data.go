package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

func ctx() context.Context {
	return context.Background()
}

// db is the Upper database session.
var db *mongo.Database

// init creates the Upper database connection.
func init() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost"))
	if err != nil {
		panic(err)
	}
	db = client.Database("moyai")

	userCollection = db.Collection("users")
	teamCollection = db.Collection("teams")
}

// Close closes and saves the data.
func Close() error {
	usersMu.Lock()
	teamsMu.Lock()
	defer usersMu.Unlock()
	defer teamsMu.Unlock()

	for _, u := range users {
		filter := bson.M{"$or": []bson.M{{"name": strings.ToLower(u.Name)}, {"xuid": u.XUID}}}
		update := bson.M{"$set": u}

		res, err := userCollection.UpdateOne(ctx(), filter, update)
		if err != nil {
			return err
		}

		if res.MatchedCount == 0 {
			_, err = userCollection.InsertOne(ctx(), u)
			return err
		}
	}

	for _, t := range teams {
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
	}
	return nil
}

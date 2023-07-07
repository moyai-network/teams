package data

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
func (f Focus) Value() string {
	return f.value
}

// Type returns the type of focus.
func (f Focus) Type() FocusType {
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
func (f Focus) MarshalBSON() ([]byte, error) {
	d := focusData{
		Kind:  f.focusType.n,
		Value: f.value,
	}
	return bson.Marshal(d)
}

// ctx returns a context.Context.
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

func ResetTeams() {
	for _, t := range Teams() {
		DisbandTeam(t)
	}
}

func ResetUsers() {
	users := db.Collection("users")
	_, _ = users.DeleteMany(context.Background(), bson.M{})
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

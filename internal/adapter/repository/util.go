package repository

import (
	"context"
	"github.com/moyai-network/teams/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"time"
)

func saveObject(col *mongo.Collection, name string, v any) error {
	filter := bson.M{"name": bson.M{"$eq": name}}
	update := bson.M{"$set": v}

	res, err := col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		_, err = col.InsertOne(context.Background(), v)
		if err != nil {
			return err
		}
	}
	return nil
}
func deleteObject(col *mongo.Collection, name string) error {
	filter := bson.M{"name": bson.M{"$eq": name}}
	_, err := col.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func eq(a, b float64) bool {
	return math.Abs(a-b) <= 1e-5
}

func updatedTeamRegeneration(t model.Team) model.Team {
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

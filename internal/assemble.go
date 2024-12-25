package internal

import (
	"context"
	"github.com/moyai-network/teams/internal/adapter/repository"
	"github.com/moyai-network/teams/internal/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Assemble() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017").SetServerAPIOptions(serverAPI))
	if err != nil {
		panic(err)
	}

	db := client.Database("odju")
	core.TeamRepository, err = repository.NewTeamRepository(db.Collection("teams"))
	if err != nil {
		panic(err)
	}
	core.UserRepository, err = repository.NewUserRepository(db.Collection("users"))
	if err != nil {
		panic(err)
	}
}

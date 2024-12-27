package internal

import (
	"context"
	"github.com/bedrock-gophers/knockback/knockback"
	"github.com/bedrock-gophers/role/role"
	"github.com/bedrock-gophers/tag/tag"
	"github.com/moyai-network/teams/internal/adapter/repository"
	"github.com/moyai-network/teams/internal/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Assemble() {
	err := knockback.Load("assets/knockback.json")
	if err != nil {
		panic(err)
	}
	err = role.Load("assets/roles")
	if err != nil {
		panic(err)
	}
	err = tag.Load("assets/tags")
	if err != nil {
		panic(err)
	}

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

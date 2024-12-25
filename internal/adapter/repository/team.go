package repository

import (
	"context"
	"github.com/moyai-network/teams/internal/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"iter"
	"slices"
	"strings"
	"sync"
)

type TeamRepository struct {
	collection *mongo.Collection
	teams      map[string]model.Team
	sync.Mutex
}

func NewTeamRepository(collection *mongo.Collection) (*TeamRepository, error) {
	repo := &TeamRepository{
		collection: collection,
		teams:      make(map[string]model.Team),
	}

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var teams []model.Team
	if err = cursor.All(context.Background(), &teams); err != nil {
		return nil, err
	}

	for _, team := range teams {
		repo.teams[team.Name] = team
	}

	return repo, nil
}

func (u *TeamRepository) FindByMemberName(name string) (model.Team, bool) {
	for _, team := range u.teams {
		for _, member := range team.Members {
			if strings.EqualFold(member.DisplayName, name) {
				return updatedTeamRegeneration(team), true
			}
		}
	}
	return model.Team{}, false
}

func (u *TeamRepository) FindByName(name string) (model.Team, bool) {
	u.Lock()
	defer u.Unlock()
	tm, ok := u.teams[strings.ToLower(name)]
	return updatedTeamRegeneration(tm), ok
}

func (u *TeamRepository) FindAll() iter.Seq[model.Team] {
	u.Lock()
	defer u.Unlock()

	var teams []model.Team
	for _, team := range u.teams {
		teams = append(teams, updatedTeamRegeneration(team))
	}
	return slices.Values(teams)
}

func (u *TeamRepository) Save(team model.Team) {
	u.Lock()
	defer u.Unlock()
	u.teams[team.Name] = team

	go func() {
		err := saveObject(u.collection, team.Name, team)
		if err != nil {
			logrus.Errorf("Mongo insert: %s", err)
		}
	}()
}

func (u *TeamRepository) Delete(team model.Team) {
	u.Lock()
	defer u.Unlock()
	delete(u.teams, team.Name)

	go func() {
		err := deleteObject(u.collection, team.Name)
		if err != nil {
			logrus.Errorf("Mongo delete: %s", err)
		}
	}()
}

package repository

import (
	"context"
	"fmt"
	"github.com/moyai-network/teams/internal/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"iter"
	"maps"
	"strings"
	"sync"
)

type UserRepository struct {
	collection *mongo.Collection
	users      map[string]model.User
	sync.Mutex
}

func NewUserRepository(collection *mongo.Collection) (*UserRepository, error) {
	repo := &UserRepository{
		collection: collection,
		users:      make(map[string]model.User),
	}

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var users []model.User
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	for _, u := range users {
		repo.users[u.Name] = u
		fmt.Println(u.DisplayName, u.Roles.All())
	}
	return repo, nil
}

func (u *UserRepository) FindByName(name string) (model.User, bool) {
	u.Lock()
	defer u.Unlock()
	user, ok := u.users[strings.ToLower(name)]
	return user, ok
}

func (u *UserRepository) FindAll() iter.Seq[model.User] {
	u.Lock()
	defer u.Unlock()
	return maps.Values(u.users)
}

func (u *UserRepository) Save(user model.User) {
	u.Lock()
	defer u.Unlock()
	u.users[user.Name] = user

	go func() {
		err := saveObject(u.collection, user.Name, user)
		if err != nil {
			logrus.Errorf("Mongo insert: %s", err)
		}
	}()
}

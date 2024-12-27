package repository

import (
	"context"
	"github.com/moyai-network/teams/internal/model"
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

	saveChan chan savable
}

func NewUserRepository(collection *mongo.Collection) (*UserRepository, error) {
	repo := &UserRepository{
		collection: collection,
		users:      make(map[string]model.User),
		saveChan:   make(chan savable),
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
	}

	go startSaveWorker(repo.collection, repo.saveChan)
	return repo, nil
}

func (u *UserRepository) FindByName(name string) (model.User, bool) {
	u.Lock()
	usr, ok := u.users[strings.ToLower(name)]
	u.Unlock()
	return usr, ok
}

func (u *UserRepository) FindAll() iter.Seq[model.User] {
	u.Lock()
	defer u.Unlock()
	return maps.Values(u.users)
}

func (u *UserRepository) Save(user model.User) {
	u.Lock()
	u.users[user.Name] = user
	u.Unlock()
	u.saveChan <- newSavable(user.Name, user)
}

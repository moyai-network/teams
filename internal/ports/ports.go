package ports

import (
	"github.com/moyai-network/teams/internal/ports/model"
	"iter"
)

type UserRepository interface {
	FindByName(name string) (model.User, bool)
	FindAll() iter.Seq[model.User]
	Save(model.User)
}

type TeamRepository interface {
	FindByMemberName(name string) (model.Team, bool)
	FindByName(name string) (model.Team, bool)
	FindAll() iter.Seq[model.Team]
	Save(model.Team)
	Delete(model.Team)
}

package ports

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
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

// Crate represents a crate utilized to Reward users.
type Crate interface {
	// Name returns the name of the crate.
	Name() string
	// Position returns the position of the crate.
	Position() mgl64.Vec3
	// Facing returns the facing of the crate.
	Facing() cube.Face
	// Rewards returns the rewards associated with the crate.
	Rewards() []model.Reward
}

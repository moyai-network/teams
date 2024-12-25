package ports

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"
	model2 "github.com/moyai-network/teams/internal/model"
	"iter"
)

type UserRepository interface {
	FindByName(name string) (model2.User, bool)
	FindAll() iter.Seq[model2.User]
	Save(model2.User)
}

type TeamRepository interface {
	FindByMemberName(name string) (model2.Team, bool)
	FindByName(name string) (model2.Team, bool)
	FindAll() iter.Seq[model2.Team]
	Save(model2.Team)
	Delete(model2.Team)
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
	Rewards() []model2.Reward
}

type Class interface {
	Armour() [4]item.ArmourTier
	Effects() []effect.Effect
}

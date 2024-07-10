package crate

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/moyai/colour"
)

// Crate represents a crate utilized to Reward users.
type Crate interface {
	// Name returns the name of the crate.
	Name() string
	// Position returns the position of the crate.
	Position() mgl64.Vec3
	// Facing returns the facing of the crate.
	Facing() cube.Face
	// Rewards returns the rewards associated with the crate.
	Rewards() []Reward
}

// Reward represents a crate reward.
type Reward struct {
	stack  item.Stack
	chance int
}

// NewReward returns a new reward.
func NewReward(stack item.Stack, chance int) Reward {
	return Reward{stack: stack, chance: chance}
}

// Stack returns the stack representing the reward.
func (r Reward) Stack() item.Stack {
	return r.stack
}

// Chance returns the chance of the reward being given.
func (r Reward) Chance() int {
	return r.chance
}

var (
	// crates contains all registered moose.Crate implementations.
	crates []Crate
	// cratesRewards contains all registered moose.Crate implementations along with their rewards.
	cratesRewards = map[Crate][]item.Stack{}
	// cratesByName contains all registered moose.Crate implementations indexed by their name.
	cratesByName = map[string]Crate{}
)

func All() []Crate {
	return crates
}

func ByName(s string) (Crate, bool) {
	c, ok := cratesByName[s]
	return c, ok
}

func SelectReward(c Crate) item.Stack {
	r, _ := cratesRewards[c]
	return r[rand.Intn(len(r))]
}

func Register(c Crate) {
	crates = append(crates, c)
	cratesByName[colour.StripMinecraftColour(c.Name())] = c

	var stacks []item.Stack
	for _, r := range c.Rewards() {
		if r.stack.Empty() {
			continue
		}
		for i := 0; i < r.Chance(); i++ {
			stacks = append(stacks, r.Stack())
		}
	}
	cratesRewards[c] = stacks
}

func init() {
	Register(conquest{})
	Register(koth{})
	Register(partner{})
	Register(pharaoh{})
	Register(menes{})
	Register(ramses{})
}

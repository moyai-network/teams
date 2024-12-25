package crate

import (
	"github.com/moyai-network/teams/internal/core/colour"
	"github.com/moyai-network/teams/internal/ports"
	"math/rand"

	"github.com/df-mc/dragonfly/server/item"
)

var (
	// crates contains all registered moose.Crate implementations.
	crates []ports.Crate
	// cratesRewards contains all registered moose.Crate implementations along with their rewards.
	cratesRewards = map[ports.Crate][]item.Stack{}
	// cratesByName contains all registered moose.Crate implementations indexed by their name.
	cratesByName = map[string]ports.Crate{}
)

func All() []ports.Crate {
	return crates
}

func ByName(s string) (ports.Crate, bool) {
	c, ok := cratesByName[s]
	return c, ok
}

func SelectReward(c ports.Crate) item.Stack {
	r := cratesRewards[c]
	return r[rand.Intn(len(r))]
}

func Register(c ports.Crate) {
	crates = append(crates, c)
	cratesByName[colour.StripMinecraftColour(c.Name())] = c

	var stacks []item.Stack
	for _, r := range c.Rewards() {
		if r.Stack().Empty() {
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
	Register(seasonal{})
	Register(partner{})
	Register(pharaoh{})
	Register(menes{})
	Register(ramses{})
}

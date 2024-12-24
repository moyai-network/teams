package model

import "github.com/df-mc/dragonfly/server/item"

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

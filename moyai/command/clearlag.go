package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/role"
)

// ClearLag clears all entitys on the floor.
type ClearLag struct{}

// Run ...
func (c ClearLag) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	var itemCount int64

	for _, e := range p.World().Entities() {
		if _, ok := e.Type().(entity.ItemType); ok {
			err := e.Close()
			if err != nil {
				continue
			}
			p.World().RemoveEntity(e)
			itemCount++
		}
	}

	moyai.Messagef(p, "command.clearlag", itemCount)
}

// Allow ...
func (c ClearLag) Allow(s cmd.Source) bool {
	return allow(s, true, role.Operator{})
}

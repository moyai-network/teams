package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

// ClearLag clears all entitys on the floor.
type ClearLag struct{ operatorAllower }

// Run ...
func (c ClearLag) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	/*p, ok := s.(*player.Player)
	if !ok {
		return
	}

	var itemCount int64

	for e := range tx.Entities() {
		if _, ok := e.H().Type().(entity.ItemBehaviour); ok {
			err := e.Close()
			if err != nil {
				continue
			}
			tx.RemoveEntity(e)
			itemCount++
		}
	}

	internal.Messagef(p, "command.clearlag", itemCount)*/
	panic("todo")
}

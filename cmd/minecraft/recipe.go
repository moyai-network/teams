package minecraft

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/recipe"
)

func init() {
	for _, w := range block.WoodTypes() {
		plank := item.NewStack(block.Planks{
			Wood: w,
		}, 1)

		recipe.Register(recipe.NewShaped([]recipe.Item{
			item.Stack{}, plank, plank,
			item.Stack{}, plank, plank,
			item.Stack{}, plank, plank,
		}, item.NewStack(block.WoodDoor{
			Wood: w}, 3), recipe.NewShape(3, 3), "crafting_table"))
	}

}

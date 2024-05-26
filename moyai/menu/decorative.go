package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

type Decorative struct{}

func NewDecorativeMenu(p *player.Player) inv.Menu {
	m := inv.NewMenu(Decorative{}, "<gold>Decorative</gold>", inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	return m.WithStacks(stacks...)
}

func (Decorative) Submit(p *player.Player, i item.Stack) {

}

package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/hako/durafmt"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

func cornerFilledStack() []item.Stack {
	stack := item.NewStack(block.StainedGlassPane{Colour: item.ColourPink()}, 1).WithCustomName(text.Colourf("<aqua>Moyai</aqua>")).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 1))

	return []item.Stack{
		0: stack,
		1: stack,
		9: stack,

		7:  stack,
		8:  stack,
		17: stack,

		36: stack,
		45: stack,
		46: stack,

		52: stack,
		53: stack,
		44: stack,
	}
}

type Prizes struct {
	open bool
}

func SendPrizesMenu(p *player.Player) {
	u, err := data.LoadUserFromXUID(p.XUID())
	if err != nil {
		return
	}
	h, ok := p.Handler().(userHandler)
	if !ok {
		return
	}

	sub := &Prizes{}
	m := inv.NewMenu(sub, "Play Time Rewards", inv.ContainerChest{DoubleChest: true})
	stacks := cornerFilledStack()

	sub.open = true
	go func() {
		for sub.open {
			playtime := u.PlayTime + h.LogTime()

			stacks[13] = item.NewStack(item.Compass{}, 1).WithCustomName(text.Colourf("<yellow>Playtime:</yellow>")).WithLore(text.Colourf("<purple>%s</purple>", durafmt.Parse(playtime).LimitFirstN(2)))
			// REWARD INDEXES:
			// 19
			// 20
			// 21
			// 22
			// 23
			// 29
			// 31

			m = m.WithStacks(stacks...)

			inv.SendMenu(p, m)
			<-time.After(time.Second)
		}
	}()
}

func (*Prizes) Submit(p *player.Player, i item.Stack) {

}
func (pr *Prizes) Close(p *player.Player) {
	pr.open = false
}

type userHandler interface {
	player.Handler
	LogTime() time.Duration
}

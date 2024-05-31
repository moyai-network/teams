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

			stacks[20] = formatRewardItem(1, time.Hour, playtime, false)
			stacks[21] = formatRewardItem(1, time.Hour*3, playtime, false)
			stacks[22] = formatRewardItem(1, time.Hour*5, playtime, false)
			stacks[23] = formatRewardItem(1, time.Hour*7, playtime, false)
			stacks[24] = formatRewardItem(1, time.Hour*10, playtime, false)
			stacks[30] = formatRewardItem(1, time.Hour*15, playtime, false)
			stacks[32] = formatRewardItem(1, time.Hour*24, playtime, false)

			m = m.WithStacks(stacks...)

			inv.SendMenu(p, m)
			<-time.After(time.Second)
		}
	}()
}

func formatRewardItem(n int, requiredPlayTime time.Duration, playtime time.Duration, claimed bool) item.Stack {
	col := item.ColourGrey()
	if claimed {
		col = item.ColourGreen()
	}

	it := item.NewStack(item.Dye{Colour: col}, 1).WithCustomName(text.Colourf("<red>Reward #%d</red>", n))
	if claimed {
		return it.WithLore(text.Colourf("<green>Claimed</green>"))
	}
	lores := []string{
		text.Colourf("<yellow>%s of play time</yellow>", durafmt.Parse(requiredPlayTime).LimitFirstN(1)),
	}

	if playtime >= requiredPlayTime {
		lores = append(lores, text.Colourf("</grey>You may claim this prize now</green>"))
	} else {
		lores = append(lores, text.Colourf("<grey>You may claim this prize in</grey> <yellow>%s<yellow>", durafmt.Parse(requiredPlayTime-playtime).LimitFirstN(3)))
	}

	return it.WithLore(lores...)
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

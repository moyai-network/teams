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
			rew := u.Teams.ClaimedRewards

			stacks[13] = item.NewStack(item.Compass{}, 1).WithCustomName(text.Colourf("<yellow>Playtime:</yellow>")).WithLore(text.Colourf("<purple>%s</purple>", durafmt.Parse(playtime).LimitFirstN(2)))

			stacks[20] = formatRewardItem(1, time.Hour, playtime, rew.Contains(1))
			stacks[21] = formatRewardItem(2, time.Hour*3, playtime, rew.Contains(2))
			stacks[22] = formatRewardItem(3, time.Hour*5, playtime, rew.Contains(3))
			stacks[23] = formatRewardItem(4, time.Hour*7, playtime, rew.Contains(4))
			stacks[24] = formatRewardItem(5, time.Hour*10, playtime, rew.Contains(5))
			stacks[30] = formatRewardItem(6, time.Hour*15, playtime, rew.Contains(6))
			stacks[32] = formatRewardItem(7, time.Hour*24, playtime, rew.Contains(7))

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

	it := item.NewStack(item.Dye{Colour: col}, 1).WithCustomName(text.Colourf("<red>Reward #%d</red>", n)).WithValue("index", n)
	if claimed {
		return it.WithLore(text.Colourf("<green>Claimed</green>"))
	}
	lores := []string{
		text.Colourf("<yellow>%s of play time</yellow>", durafmt.Parse(requiredPlayTime).LimitFirstN(1)),
	}

	if playtime >= requiredPlayTime {
		lores = append(lores, text.Colourf("</grey>You may claim this prize now</green>"))
	} else {
		limit := 3
		dur := requiredPlayTime - playtime
		if dur < time.Hour {
			limit = 2
		}
		lores = append(lores, text.Colourf("<grey>You may claim this prize in</grey> <yellow>%s<yellow>", durafmt.Parse(requiredPlayTime-playtime).LimitFirstN(limit)))
	}

	return it.WithLore(lores...)
}

func (*Prizes) Submit(p *player.Player, it item.Stack) {
	dye, ok := it.Item().(item.Dye)
	if !ok {
		return
	}
	if dye.Colour != item.ColourGreen() {
		return
	}
	i, ok := it.Value("index")
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	in, ok := i.(int)
	if !ok {
		return
	}
	if u.Teams.ClaimedRewards.Contains(in) {
		return
	}

	u.Teams.ClaimedRewards.Add(in)
}
func (pr *Prizes) Close(p *player.Player) {
	pr.open = false
}

type userHandler interface {
	player.Handler
	LogTime() time.Duration
}

package menu

import (
	"github.com/moyai-network/teams/internal/core"
	item2 "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/model"
	"time"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/hako/durafmt"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type reward struct {
	slot             int
	requiredPlayTime time.Duration
	rewards          []item.Stack
}

var (
	advancedRewards = []item.Stack{
		item2.NewSpecialItem(item2.PartnerPackageType{}, 25),
		item2.NewKey(item2.KeyTypePharaoh, 25),
	}

	rewards = []reward{
		// 1 hour
		{slot: 4, requiredPlayTime: time.Hour, rewards: []item.Stack{
			item2.NewSpecialItem(item2.PartnerPackageType{}, 5),
			item2.NewKey(item2.KeyTypePharaoh, 5),
		}},
		// 3 hours
		{slot: 15, requiredPlayTime: time.Hour * 3, rewards: []item.Stack{
			item2.NewSpecialItem(item2.PartnerPackageType{}, 10),
			item2.NewKey(item2.KeyTypePharaoh, 10),
		}},
		// 5 hours
		{slot: 34, requiredPlayTime: time.Hour * 5, rewards: []item.Stack{
			item2.NewSpecialItem(item2.PartnerPackageType{}, 15),
			item2.NewKey(item2.KeyTypePharaoh, 15),
		}},
		// 7 hours
		{slot: 41, requiredPlayTime: time.Hour * 7, rewards: []item.Stack{
			item2.NewSpecialItem(item2.PartnerPackageType{}, 20),
			item2.NewKey(item2.KeyTypePharaoh, 20),
		}},
		// 10 hours
		{slot: 39, requiredPlayTime: time.Hour * 10, rewards: advancedRewards},
		// 15 hours
		{slot: 28, requiredPlayTime: time.Hour * 15, rewards: append(advancedRewards, []item.Stack{
			item2.NewKey(item2.KeyTypePartner, 5),
		}...)},
		// 24 hours
		{slot: 11, requiredPlayTime: time.Hour * 24, rewards: append(advancedRewards, []item.Stack{
			item2.NewKey(item2.KeyTypePartner, 5),
		}...)},
	}
)

type Prizes struct {
	close chan struct{}
}

func SendPrizesMenu(p *player.Player) {
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}

	pr := &Prizes{
		close: make(chan struct{}),
	}
	go func() {
		pr.sendPrizesMenu(u, p)
		for {
			select {
			case <-time.After(time.Second):
				pr.sendPrizesMenu(u, p)
			case <-pr.close:
				return
			}
		}
	}()
}

func (pr *Prizes) sendPrizesMenu(u model.User, p *player.Player) {
	h, ok := p.Handler().(userHandler)
	if !ok {
		return
	}
	m := inv.NewMenu(pr, "Play Time Rewards", inv.ContainerChest{DoubleChest: true})
	stacks := make([]item.Stack, 54)

	playtime := u.PlayTime + h.LogTime()
	rew := u.Teams.ClaimedRewards

	stacks[22] = item.NewStack(item.Clock{}, 1).
		WithCustomName(text.Colourf("<yellow>Playtime:</yellow>")).
		WithLore(text.Colourf("<purple>%s</purple>", durafmt.Parse(playtime).LimitFirstN(2)))

	var claimedIndex int
	for i, r := range rewards {
		index := i + 1
		claimed := rew.Contains(index)
		if claimed {
			claimedIndex = index
		}
		stacks[r.slot] = formatRewardItem(index, r.requiredPlayTime, playtime, claimed)
	}

	nextReward := rewards[claimedIndex]
	cornerColors := colourFromTimeLeft(nextReward.requiredPlayTime - playtime)
	fillCorners(stacks, cornerColors)

	m = m.WithStacks(stacks...)
	inv.SendMenu(p, m)
}

func (pr *Prizes) Submit(p *player.Player, stack item.Stack) {
	dye, ok := stack.Item().(item.Dye)
	if !ok {
		return
	}
	if dye.Colour != item.ColourLime() {
		return
	}
	_, ok = stack.Value("claimable")
	if !ok {
		return
	}
	i, ok := stack.Value("index")
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
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
	for _, r := range rewards[in-1].rewards {
		item2.AddOrDrop(p, r)
	}

	core.UserRepository.Save(u)
	pr.sendPrizesMenu(u, p)
}

func colourFromTimeLeft(timeLeft time.Duration) item.Colour {
	col := item.ColourRed()
	if timeLeft <= time.Hour && timeLeft > time.Minute*45 {
		col = item.ColourOrange()
	} else if timeLeft <= time.Minute*45 && timeLeft > time.Minute*30 {
		col = item.ColourYellow()
	} else if timeLeft <= time.Minute*30 && timeLeft > time.Minute*15 {
		col = item.ColourGreen()
	} else if timeLeft <= time.Minute*15 {
		col = item.ColourLime()
	}
	return col
}

func fillCorners(stacks []item.Stack, col item.Colour) {
	stack := item.NewStack(block.StainedGlass{Colour: col}, 1).WithCustomName(text.Colourf("<aqua>Moyai</aqua>")).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking, 1))

	stacks[0] = stack
	stacks[1] = stack
	stacks[9] = stack

	stacks[7] = stack
	stacks[8] = stack
	stacks[17] = stack

	stacks[36] = stack
	stacks[45] = stack
	stacks[46] = stack

	stacks[52] = stack
	stacks[53] = stack
	stacks[44] = stack
}

func formatRewardItem(n int, requiredPlayTime time.Duration, playtime time.Duration, claimed bool) item.Stack {
	col := colourFromTimeLeft(requiredPlayTime - playtime)
	itm := item.NewStack(item.Dye{Colour: col}, 1).WithCustomName(text.Colourf("<red>Reward #%d</red>", n)).WithValue("index", n)
	if claimed {
		itm = item.NewStack(item.FireworkStar{
			FireworkExplosion: item.FireworkExplosion{
				Colour: item.ColourGrey(),
			},
		}, 1)
		return itm.WithLore(text.Colourf("<green>Claimed</green>"))
	}
	lores := []string{
		text.Colourf("<yellow>%s of play time</yellow>", durafmt.Parse(requiredPlayTime).LimitFirstN(1)),
	}

	if playtime >= requiredPlayTime {
		lores = append(lores, text.Colourf("<green>You may claim this prize now</green>"))
		itm = itm.WithValue("claimable", true)
	} else {
		limit := 3
		dur := requiredPlayTime - playtime
		if dur < time.Hour {
			limit = 2
		} else if dur < time.Minute {
			limit = 1
		}
		lores = append(lores, text.Colourf("<grey>You may claim this prize in</grey> <yellow>%s<yellow>", durafmt.Parse(requiredPlayTime-playtime).LimitFirstN(limit)))
	}

	return itm.WithLore(lores...)
}

func (pr *Prizes) Close(p *player.Player) {
	pr.close <- struct{}{}
}

type userHandler interface {
	player.Handler
	LogTime() time.Duration
}

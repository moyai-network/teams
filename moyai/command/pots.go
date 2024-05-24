package command

import (
	"strings"
	"time"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hako/durafmt"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
)

type Pots struct{}

func (Pots) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	if u.Teams.Refill.Active() {
		user.Messagef(p, "user.refill.cooldown", durafmt.Parse(u.Teams.Refill.Remaining()).LimitFirstN(2))
		return
	}

	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		user.Messagef(p, "user.team-less")
		return
	}
	if tm.Claim == (area.Area{}) {
		user.Messagef(p, "team.claim.none")
		return
	}

	if !tm.Claim.Vec3WithinOrEqualFloorXZ(p.Position()) {
		user.Messagef(p, "team.claim.not-within")
		return
	}

	teams, err := data.LoadAllTeams()
	if err != nil {
		return
	}

	for _, t := range teams {
		if strings.EqualFold(t.Name, tm.Name) {
			placePotionChests(p)
			u.Teams.Refill.Set(time.Hour * 4)
			break
		}
	}
}

func placePotionChests(p *player.Player) {
	pos := cube.PosFromVec3(p.Position().Add(mgl64.Vec3{1, 1, 2}))

	placePairedChests(
		pos.Add(cube.Pos{1, 0, 0}),
		pos,
		p, item.SplashPotion{Type: potion.StrongHealing()},
	)
	placePairedChests(
		pos.Add(cube.Pos{1, 1, 0}),
		pos.Add(cube.Pos{0, 1, 0}),
		p, item.Potion{Type: potion.LongSwiftness()},
	)
	placePairedChests(
		pos.Add(cube.Pos{1, 2, 0}),
		pos.Add(cube.Pos{0, 2, 0}),
		p, item.Potion{Type: potion.LongFireResistance()},
	)
}

func placePairedChests(a, b cube.Pos, p *player.Player, it world.Item) {
	w := p.World()
	che := block.NewChest()
	pair := block.NewChest()
	fillInventory(che.Inventory(), item.NewStack(it, 1))
	fillInventory(pair.Inventory(), item.NewStack(it, 1))

	w.SetBlock(a, che, nil)
	w.SetBlock(b, pair, nil)
	ch, pair, ok := che.Pair(w, a, b)
	if ok {
		w.SetBlock(a, ch, nil)
		w.SetBlock(b, pair, nil)
	}
}

func fillInventory(in *inventory.Inventory, it item.Stack) {
	for i := 0; i < in.Size(); i++ {
		in.SetItem(i, it)
	}
}

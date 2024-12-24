package command

import (
	"github.com/moyai-network/teams/internal/core/area"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"strings"
	"time"
	_ "unsafe"

	"github.com/moyai-network/teams/internal"

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
)

type Pots struct{}

func (Pots) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, err := data2.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	if u.Teams.Refill.Active() {
		internal.Messagef(p, "user.refill.cooldown", durafmt.Parse(u.Teams.Refill.Remaining()).LimitFirstN(2))
		return
	}

	tm, err := data2.LoadTeamFromMemberName(p.Name())
	if err != nil {
		internal.Messagef(p, "user.team-less")
		return
	}
	if tm.Claim == (area.Area{}) {
		internal.Messagef(p, "team.claim.none")
		return
	}

	if !tm.Claim.Vec3WithinOrEqualFloorXZ(p.Position()) {
		internal.Messagef(p, "team.claim.not-within")
		return
	}

	teams, err := data2.LoadAllTeams()
	if err != nil {
		return
	}

	for _, t := range teams {
		if strings.EqualFold(t.Name, tm.Name) {
			placePotionChests(tx, p)
			u.Teams.Refill.Set(time.Hour * 4)
			break
		}
	}
}

func placePotionChests(tx *world.Tx, p *player.Player) {
	pos := cube.PosFromVec3(p.Position().Add(mgl64.Vec3{1, 1, 2}))

	placePairedChests(tx,
		pos.Add(cube.Pos{1, 0, 0}),
		pos,
		p, item.SplashPotion{Type: potion.StrongHealing()},
	)
	placePairedChests(tx,
		pos.Add(cube.Pos{1, 1, 0}),
		pos.Add(cube.Pos{0, 1, 0}),
		p, item.Potion{Type: potion.LongSwiftness()},
	)
	placePairedChests(tx,
		pos.Add(cube.Pos{1, 2, 0}),
		pos.Add(cube.Pos{0, 2, 0}),
		p, item.Potion{Type: potion.LongFireResistance()},
	)
}

func placePairedChests(tx *world.Tx, a, b cube.Pos, p *player.Player, it world.Item) {
	che := block.NewChest()
	pair := block.NewChest()
	fillInventory(che.Inventory(tx, a), item.NewStack(it, 1))
	fillInventory(pair.Inventory(tx, b), item.NewStack(it, 1))

	tx.SetBlock(a, che, nil)
	tx.SetBlock(b, pair, nil)
	ch, pair, ok := chest_pair(tx, a, b)
	if ok {
		tx.SetBlock(a, ch, nil)
		tx.SetBlock(b, pair, nil)
	}
}

func fillInventory(in *inventory.Inventory, it item.Stack) {
	for i := 0; i < in.Size(); i++ {
		in.SetItem(i, it)
	}
}

// noinspection ALL
//
//go:linkname chest_pair github.com/df-mc/dragonfly/server/block.(*Chest).pair
func chest_pair(tx *world.Tx, pos, pairPos cube.Pos) (ch, pair block.Chest, ok bool)

package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/roles"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"slices"
)

type PlayerListTeleport struct {
	index int
}

func SendPlayerListTeleportMenu(p *player.Player) {
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	pl := &PlayerListTeleport{}
	pl.sendPlayerListTeleportMenu(p.Tx(), u, p)
}

func (pl *PlayerListTeleport) sendPlayerListTeleportMenu(tx *world.Tx, u data.User, p *player.Player) {
	rl := u.Roles.Highest()
	if !roles.Staff(rl) && !u.Roles.Contains(roles.Operator()) {
		return
	}

	var index = pl.index
	if internal.PlayerCount() < (index + 1*18) {
		index = 0
	}

	players := slices.Collect(internal.Players(tx))
	if index == 0 {
		if internal.PlayerCount() > 18 {
			players = players[:18]
		}
	} else {
		players = players[index*18 : (index+1)*18]
	}

	m := inv.NewMenu(pl, "Teleport to Player", inv.ContainerChest{})
	var stacks = make([]item.Stack, 27)

	for i, trg := range players {
		if i > 18 {
			break
		}
		stacks[i] = item.NewStack(block.Skull{Type: block.PlayerHead()}, 1).WithCustomName(text.Colourf("<aqua>%s</aqua>", trg.Name())).WithValue("player_name", trg.Name())
	}

	stacks[21] = item.NewStack(item.Arrow{}, 1).WithCustomName(text.Colourf("<aqua>Previous</aqua>")).WithValue("previous", true)
	stacks[23] = item.NewStack(item.Arrow{}, 1).WithCustomName(text.Colourf("<aqua>Next</aqua>")).WithValue("next", true)

	inv.SendMenu(p, m.WithStacks(stacks...))
}

func (pl *PlayerListTeleport) Submit(p *player.Player, stack item.Stack) {
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	switch stack.Item().(type) {
	case item.Arrow:
		if _, ok := stack.Value("previous"); ok {
			if pl.index < 0 {
				return
			}
			pl.index--
			pl.sendPlayerListTeleportMenu(p.Tx(), u, p)
			return
		}
		if _, ok := stack.Value("next"); ok {
			if internal.PlayerCount() < (pl.index+1)*18 {
				return
			}
			pl.index++
			pl.sendPlayerListTeleportMenu(p.Tx(), u, p)
			return
		}
	case block.Skull:
		v, ok := stack.Value("player_name")
		if !ok {
			break
		}
		for trg := range internal.Players(p.Tx()) {
			if trg.Name() == v {
				trg.Teleport(p.Position())
				p.Tx().AddEntity(trg.H())
				break
			}
		}
	default:
		return
	}
	inv.CloseContainer(p)
}

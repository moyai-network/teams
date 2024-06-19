package minecraft

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
)

type worldHandler struct {
	world.NopHandler
	w *world.World
}

func (w *worldHandler) HandleLiquidFlow(ctx *event.Context, from, into cube.Pos, liquid world.Liquid, replaced world.Block) {
	teams, _ := data.LoadAllTeams()
	var initialTeam data.Team
	var nextTeam data.Team

	for _, t := range teams {
		if len(initialTeam.Name) != 0 && len(nextTeam.Name) != 0 {
			break
		}
		if t.Claim == (area.Area{}) {
			continue
		}
		if t.Claim.Vec3WithinOrEqualXZ(from.Vec3()) {
			initialTeam = t
		}
		if t.Claim.Vec3WithinOrEqualXZ(into.Vec3()) {
			nextTeam = t
		}
	}

	if len(nextTeam.Name) == 0 || nextTeam.Name != initialTeam.Name {
		ctx.Cancel()
		return
	}
}

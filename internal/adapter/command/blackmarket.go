package command

import (
	"github.com/df-mc/dragonfly/server/world"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/internal"
)

type BlackMarket struct{}

func (BlackMarket) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
}

func (BlackMarket) Allow(src cmd.Source) bool {
	return Allow(src, false) && time.Since(internal.LastBlackMarket()) < time.Minute*10
}

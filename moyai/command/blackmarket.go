package command

import (
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/moyai"
)

type BlackMarket struct{}

func (BlackMarket) Run(src cmd.Source, out *cmd.Output) {
}

func (BlackMarket) Allow(src cmd.Source) bool {
	return Allow(src, false) && time.Since(moyai.LastBlackMarket()) < time.Minute*10
}

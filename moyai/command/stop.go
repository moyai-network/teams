package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"os"
	"time"
)

type Stop struct{}

func (Stop) Run(src cmd.Source, out *cmd.Output) {
	time.Sleep(time.Millisecond * 500)
	data.FlushCache()
	_ = moyai.Server().Close()
	os.Exit(0)
}

func (Stop) Allow(src cmd.Source) bool {
	return allow(src, true, role.Operator{})
}

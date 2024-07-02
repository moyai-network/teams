package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/moyai/data"
	rls "github.com/moyai-network/teams/moyai/roles"
)

type DataReset struct {
	Kind dataKind `name:"kind"`
}

func (d DataReset) Run(src cmd.Source, _ *cmd.Output) {
	switch d.Kind {
	case "all":
		data.Reset()
	case "partial":
		data.PartialReset()
	}
}

// Allow ...
func (DataReset) Allow(src cmd.Source) bool {
	return allow(src, true, rls.Operator())
}

type dataKind string

// Type ...
func (dataKind) Type() string {
	return "data_kind"
}

// Options ...
func (dataKind) Options(cmd.Source) []string {
	return []string{"all", "partial"}
}

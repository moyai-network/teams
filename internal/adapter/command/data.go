package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	rls "github.com/moyai-network/teams/internal/core/roles"
)

type DataReset struct {
	Kind dataKind `name:"kind"`
}

func (d DataReset) Run(src cmd.Source, _ *cmd.Output, tx *world.Tx) {
	switch d.Kind {
	case "all":
		//data.Reset()
	case "partial":
		//data.PartialReset()
	}
}

// Allow ...
func (DataReset) Allow(src cmd.Source) bool {
	return Allow(src, true, rls.Operator())
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

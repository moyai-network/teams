package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/data"
)

type DataReset struct {
	Kind dataKind `name:"kind"`
}

func (d DataReset) Run(src cmd.Source, _ *cmd.Output) {
	switch d.Kind {
	case "users":
		data.ResetUsers()
	case "teams":
		data.ResetTeams()
	}
}

// Allow ...
func (DataReset) Allow(src cmd.Source) bool {
	return allow(src, true, role.Operator{})
}

type dataKind string

// Type ...
func (dataKind) Type() string {
	return "data_kind"
}

// Options ...
func (dataKind) Options(cmd.Source) []string {
	return []string{"users", "teams"}
}

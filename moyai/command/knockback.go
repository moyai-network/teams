package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/moyai"
)

// Knockback is a command to change the server knockback
type Knockback struct {
	operatorAllower
	Force  	 float64          `cmd:"force"`
	Height   float64          `cmd:"height"`
}


func (k Knockback) Run(src cmd.Source, out *cmd.Output) {
	moyai.SetForce(k.Force)
	moyai.SetHeight(k.Height)
	out.Printf("Set KB to (%.2f, %.2f)", k.Force, k.Height)
}
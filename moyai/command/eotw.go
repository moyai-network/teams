package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/eotw"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// EOTWStart is a command to start EOTW.
type EOTWStart struct {
	managerAllower
	Sub cmd.SubCommand `cmd:"start"`
}

// EOTWEnd is a command to end EOTW.
type EOTWEnd struct {
	managerAllower
	Sub cmd.SubCommand `cmd:"end"`
}

// Run ...
func (c EOTWStart) Run(s cmd.Source, o *cmd.Output) {
	if _, ok := eotw.Running(); ok {
		o.Print(text.Colourf("<red>EOTW is already running!</red>"))
		return
	}
	eotw.Start()

	moyai.Broadcastf("eotw.commenced")
}

// Run ...
func (c EOTWEnd) Run(s cmd.Source, o *cmd.Output) {
	if _, ok := eotw.Running(); !ok {
		o.Print(lang.Translatef(locale(s), "command.eotw.not.running"))
		return
	}
	eotw.End()

	moyai.Broadcastf("eotw.ended")
}
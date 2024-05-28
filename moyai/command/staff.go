package command

import "github.com/df-mc/dragonfly/server/cmd"

type StaffMode struct {
	adminAllower
	Sub cmd.SubCommand `cmd:"mode"`
}

func (sm StaffMode) Run(s cmd.Source, o *cmd.Output) {
	o.Print("Staff mode is not yet implemented.")
}

package command

import "github.com/df-mc/dragonfly/server/cmd"

type LastInv struct{
	adminAllower
	Target []cmd.Target `cmd:"target"`
}

func (i LastInv)  Run(src cmd.Source, o *cmd.Output) {

}

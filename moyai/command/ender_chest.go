package command

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/unsafe"
)

// Enderchest is a command to open a players enderchest
type Enderchest struct {}

// Run ...
func (e Enderchest) Run(src cmd.Source, out *cmd.Output) {
	p := src.(*player.Player)
	pos := cube.PosFromVec3(p.Rotation().Vec3().Mul(-2).Add(p.Position()))
	s := unsafe.Session(p)
	s.ViewBlockUpdate(pos, block.NewEnderChest(), 0)
	s.ViewBlockUpdate(pos.Add(cube.Pos{0, 1}), block.Air{}, 0)
	p.OpenBlockContainer(pos)
}

type MySubmittable struct{}

func (m MySubmittable) Submit(p *player.Player, it item.Stack) {
	fmt.Println("Submitted", it)
}

func (m MySubmittable) Close(p *player.Player) {
	fmt.Println("Closed")
}
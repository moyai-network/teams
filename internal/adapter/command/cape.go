package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/cape"
	rls "github.com/moyai-network/teams/internal/core/roles"
	"github.com/samber/lo"
)

type Cape struct {
	Cape capes `cmd:"capes"`
}

// Run ...
func (a Cape) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(pl.Name())
	if !ok {
		return
	}
	c, ok := cape.ByName(string(a.Cape))
	if !ok {
		return
	}
	sk := pl.Skin()
	sk.Cape = c.Cape()
	pl.SetSkin(sk)
	u.Teams.Settings.Advanced.Cape = c.Name()
	core.UserRepository.Save(u)
	internal.Messagef(pl, "cape.selected", c.Name())
}

type (
	capes string
)

// Type ...
func (capes) Type() string {
	return "cape"
}

// Options ...
func (capes) Options(s cmd.Source) (capes []string) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}

	all := cape.All()
	if rls.Premium(u.Roles.Highest()) {
		names := lo.Map(all, func(c cape.Cape, _ int) string {
			return c.Name()
		})
		capes = append(capes, names...)
	} else {
		for _, c := range all {
			if !c.Premium() {
				capes = append(capes, c.Name())
			}
		}
	}
	return
}

package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
)

type Lives struct{}

func (Lives) Run(src cmd.Source, _ *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}

	internal.Messagef(p, "command.lives", u.Teams.Lives)
}

type LivesGiveOnline struct {
	Sub    cmd.SubCommand `cmd:"give"`
	Target []cmd.Target
	Count  int
}

func (l LivesGiveOnline) Run(src cmd.Source, _ *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	tg, ok := l.Target[0].(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	target, ok := core.UserRepository.FindByName(tg.Name())
	if !ok {
		internal.Messagef(p, "command.target.unknown", tg.Name())
		return
	}
	if l.Count <= 0 {
		return
	}

	lives := u.Teams.Lives
	if lives < l.Count {
		internal.Messagef(p, "command.lives.not-enough", l.Count, tg.Name())
		return
	}

	u.Teams.Lives -= l.Count
	target.Teams.Lives += l.Count
	core.UserRepository.Save(u)
	core.UserRepository.Save(target)

	internal.Messagef(p, "command.lives.give.sender", l.Count, tg.Name())
	internal.Messagef(tg, "command.lives.give.receiver", p.Name(), l.Count)
}

type LivesGiveOffline struct {
	Sub    cmd.SubCommand `cmd:"give"`
	Target string
	Count  int
}

func (l LivesGiveOffline) Run(src cmd.Source, _ *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	target, ok := core.UserRepository.FindByName(l.Target)
	if !ok {
		internal.Messagef(p, "command.target.unknown", l.Target)
		return
	}
	if l.Count <= 0 {
		return
	}

	lives := u.Teams.Lives
	if lives < l.Count {
		internal.Messagef(p, "command.lives.not-enough", l.Count, l.Target)
		return
	}

	u.Teams.Lives -= l.Count
	target.Teams.Lives += l.Count
	core.UserRepository.Save(u)
	core.UserRepository.Save(target)

	internal.Messagef(p, "command.lives.give.sender", l.Count, l.Target)
}

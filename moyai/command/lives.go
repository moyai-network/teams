package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
)

type Lives struct{}

func (Lives) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromXUID(p.XUID())
	if err != nil {
		return
	}

	moyai.Messagef(p, "command.lives", u.Teams.Lives)
}

type LivesGiveOnline struct {
	Target cmd.Optional[[]cmd.Target]
	Count  int
}

func (l LivesGiveOnline) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	t, ok := l.Target.Load()
	if !ok {
		return
	}
	tg, ok := t[0].(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	target, err := data.LoadUserFromName(tg.Name())
	if err != nil {
		moyai.Messagef(p, "command.target.unknown", tg.Name())
		return
	}
	if l.Count <= 0 {
		return
	}

	lives := u.Teams.Lives
	if lives < l.Count {
		moyai.Messagef(p, "command.lives.not-enough", l.Count, tg.Name())
		return
	}

	u.Teams.Lives -= l.Count
	target.Teams.Lives += l.Count
	data.SaveUser(u)
	data.SaveUser(target)

	moyai.Messagef(p, "command.lives.give.sender", l.Count, tg.Name())
	moyai.Messagef(tg, "command.lives.give.receiver", p.Name(), l.Count)
}

type LivesGiveOffline struct {
	Target string
	Count  int
}

func (l LivesGiveOffline) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	target, err := data.LoadUserFromName(l.Target)
	if err != nil {
		moyai.Messagef(p, "command.target.unknown", l.Target)
		return
	}
	if l.Count <= 0 {
		return
	}

	lives := u.Teams.Lives
	if lives < l.Count {
		moyai.Messagef(p, "command.lives.not-enough", l.Count, l.Target)
		return
	}

	u.Teams.Lives -= l.Count
	target.Teams.Lives += l.Count
	data.SaveUser(u)
	data.SaveUser(target)

	moyai.Messagef(p, "command.lives.give.sender", l.Count, l.Target)
}

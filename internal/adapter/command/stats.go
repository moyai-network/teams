package command

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/model"
	"math"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/internal"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type StatsOnlineCommand struct {
	Targets cmd.Optional[[]cmd.Target] `cmd:"target"`
}

type StatsOfflineCommand struct {
	Target string `cmd:"target"`
}

func (t StatsOnlineCommand) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	targets := t.Targets.LoadOr(nil)

	if len(targets) > 1 {
		internal.Messagef(s, "command.targets.exceed")
		return
	}

	var name string

	if len(targets) == 1 {
		target, ok := targets[0].(cmd.NamedTarget)
		if !ok {
			internal.Messagef(s, "command.target.unknown", "")
			return
		}
		name = target.Name()
	} else {
		s, ok := s.(cmd.NamedTarget)
		if !ok {
			return
		}
		name = s.Name()
	}

	u, ok := core.UserRepository.FindByName(name)
	if !ok {
		return
	}

	o.Print(userStatsFormat(u))
}

func (s StatsOfflineCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	u, ok := core.UserRepository.FindByName(s.Target)
	if !ok {
		internal.Messagef(src, "command.target.unknown", s.Target)
		return
	}

	o.Print(userStatsFormat(u))
}

func userStatsFormat(u model.User) string {
	kills := u.Teams.Stats.Kills
	deaths := u.Teams.Stats.Deaths
	if deaths == 0 {
		deaths = 1
	}
	kdr := math.Round(float64(kills)/float64(deaths)*100) / 100
	return text.Colourf("<red>%s's Stats:</red>\n\uE000\n<red>Kills:</red> <yellow>%d</yellow>\n<red>Deaths:</red> <yellow>%d</yellow>\n<red>KDR:</red> <yellow>%.2f</yellow>\n<red>Killstreak:</red> <yellow>%d</yellow>\n<red>Best Killstreak:</red> <yellow>%d</yellow>\n\uE000", u.DisplayName, u.Teams.Stats.Kills, u.Teams.Stats.Deaths, kdr, u.Teams.Stats.KillStreak, u.Teams.Stats.BestKillStreak)
}

package command

import (
	"math"
	"sort"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func orderedUsersByKills() []data.User {
	usrs, _ := data.LoadAllUsers()
	sort.SliceStable(usrs, func(i, j int) bool {
		if usrs[i].Teams.Stats.Kills != usrs[j].Teams.Stats.Kills {
			return usrs[i].Teams.Stats.Kills > usrs[j].Teams.Stats.Kills
		}
		return usrs[i].DisplayName < usrs[j].DisplayName
	})
	return usrs
}

func orderedUsersByDeaths() []data.User {
	usrs, _ := data.LoadAllUsers()
	sort.SliceStable(usrs, func(i, j int) bool {
		if usrs[i].Teams.Stats.Deaths != usrs[j].Teams.Stats.Deaths {
			return usrs[i].Teams.Stats.Deaths > usrs[j].Teams.Stats.Deaths
		}
		return usrs[i].DisplayName < usrs[j].DisplayName
	})
	return usrs
}

func orderedUsersByKDR() []data.User {
	usrs, _ := data.LoadAllUsers()
	sort.SliceStable(usrs, func(i, j int) bool {
		iDeath := usrs[i].Teams.Stats.Deaths
		if iDeath == 0 {
			iDeath = 1
		}
		jDeath := usrs[j].Teams.Stats.Deaths
		if jDeath == 0 {
			jDeath = 1
		}
		iKdr := math.Round(float64(usrs[i].Teams.Stats.Kills)/float64(iDeath)*100) / 100
		jKdr := math.Round(float64(usrs[j].Teams.Stats.Kills)/float64(jDeath)*100) / 100
		if iKdr != jKdr {
			return iKdr > jKdr
		}
		return usrs[i].DisplayName < usrs[j].DisplayName
	})
	return usrs
}

func orderedUsersByKillStreaks() []data.User {
	usrs, _ := data.LoadAllUsers()
	sort.SliceStable(usrs, func(i, j int) bool {
		if usrs[i].Teams.Stats.KillStreak != usrs[j].Teams.Stats.KillStreak {
			return usrs[i].Teams.Stats.KillStreak > usrs[j].Teams.Stats.KillStreak
		}
		return usrs[i].DisplayName < usrs[j].DisplayName
	})
	return usrs
}

type LeaderboardKills struct {
	Sub cmd.SubCommand `cmd:"kills"`
}

type LeaderboardDeaths struct {
	Sub cmd.SubCommand `cmd:"deaths"`
}

type LeaderboardKillStreaks struct {
	Sub cmd.SubCommand `cmd:"killstreaks"`
}

type LeaderboardKDR struct {
	Sub cmd.SubCommand `cmd:"kdr"`
}

func (l LeaderboardKills) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	usrs := orderedUsersByKills()
	var top string
	top += text.Colourf("        <yellow>Top Kills</yellow>\n")
	top += "\uE000\n"
	for i, u := range usrs {
		if i > 9 {
			break
		}
		top += text.Colourf(" <grey>%d. <red>%s</red> (<yellow>%d</yellow>)</grey>\n", i+1, u.DisplayName, u.Teams.Stats.Kills)
	}
	top += "\uE000"
	p.Message(top)
}

func (l LeaderboardDeaths) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	usrs := orderedUsersByDeaths()
	var top string
	top += text.Colourf("        <yellow>Top Deaths</yellow>\n")
	top += "\uE000\n"
	for i, u := range usrs {
		if i > 9 {
			break
		}
		top += text.Colourf(" <grey>%d. <red>%s</red> (<yellow>%d</yellow>)</grey>\n", i+1, u.DisplayName, u.Teams.Stats.Deaths)
	}
	top += "\uE000"
	p.Message(top)
}

func (l LeaderboardKillStreaks) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	usrs := orderedUsersByKillStreaks()
	var top string
	top += text.Colourf("        <yellow>Top Kill Streaks</yellow>\n")
	top += "\uE000\n"
	for i, u := range usrs {
		if i > 9 {
			break
		}
		top += text.Colourf(" <grey>%d. <red>%s</red> (<yellow>%d</yellow>)</grey>\n", i+1, u.DisplayName, u.Teams.Stats.KillStreak)
	}
	top += "\uE000"
	p.Message(top)
}

func (l LeaderboardKDR) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	usrs := orderedUsersByKDR()
	var top string
	top += text.Colourf("        <yellow>Top KDR</yellow>\n")
	top += "\uE000\n"
	for i, u := range usrs {
		if i > 9 {
			break
		}
		deaths := u.Teams.Stats.Deaths
		if deaths == 0 {
			deaths = 1
		}
		kdr := math.Round(float64(u.Teams.Stats.Kills)/float64(deaths)*100) / 100
		top += text.Colourf(" <grey>%d. <red>%s</red> (<yellow>%.2f</yellow>)</grey>\n", i+1, u.DisplayName, kdr)
	}
	top += "\uE000"
	p.Message(top)
}

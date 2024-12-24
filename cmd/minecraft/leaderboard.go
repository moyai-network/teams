package minecraft

import (
	"github.com/df-mc/dragonfly/server/world"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/icza/abcsort"
	"github.com/moyai-network/teams/internal"
	"github.com/restartfu/roman"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func startLeaderboard() {
	killsLeaderboard := entity.NewText(formattedUserLeaderBoard("KILLS", func(u data2.User) int {
		return u.Teams.Stats.Kills
	}), cube.Pos{9, 71, 11}.Vec3Middle())

	kdrLeaderboard := entity.NewText(formattedUserLeaderBoard("KDR", func(u data2.User) float64 {
		kills := u.Teams.Stats.Kills
		deaths := u.Teams.Stats.Deaths
		if deaths == 0 {
			deaths = 1
		}
		return math.Round(float64(kills)/float64(deaths)*100) / 100
	}), cube.Pos{13, 71, 10}.Vec3Middle())

	deathsLeaderboard := entity.NewText(formattedUserLeaderBoard("DEATHS", func(u data2.User) int {
		return u.Teams.Stats.Deaths
	}), cube.Pos{17, 71, 8}.Vec3Middle())

	topPointsLeaderboard := entity.NewText(formattedTeamLeaderBoard("POINTS", func(t data2.Team) int {
		return t.Points
	}), cube.Pos{-9, 71, 11}.Vec3Middle())

	topKOTHLeaderboard := entity.NewText(formattedTeamLeaderBoard("KOTH CAPTURES", func(t data2.Team) int {
		return t.KOTHWins
	}), cube.Pos{-13, 71, 10}.Vec3Middle())

	topKillsLeaderboard := entity.NewText(formattedTeamLeaderBoard("KILLS", func(t data2.Team) int {
		kills := 0
		for _, m := range t.Members {
			u, err := data2.LoadUserFromName(m.Name)
			if err != nil {
				continue
			}
			kills += u.Teams.Stats.Kills
		}
		return kills
	}), cube.Pos{-17, 71, 8}.Vec3Middle())

	internal.Overworld().Exec(func(tx *world.Tx) {
		tx.AddEntity(killsLeaderboard)
		tx.AddEntity(kdrLeaderboard)
		tx.AddEntity(deathsLeaderboard)
		tx.AddEntity(topPointsLeaderboard)
		tx.AddEntity(topKOTHLeaderboard)
		tx.AddEntity(topKillsLeaderboard)
	})

	t := time.NewTicker(time.Minute)

	for range t.C {
		/*killsLeaderboard.SetNameTag(formattedUserLeaderBoard("KILLS", func(u data.User) int {
			return u.Teams.Stats.Kills
		}))

		kdrLeaderboard.SetNameTag(formattedUserLeaderBoard("KDR", func(u data.User) float64 {
			kills := u.Teams.Stats.Kills
			deaths := u.Teams.Stats.Deaths
			if deaths == 0 {
				deaths = 1
			}
			return math.Round(float64(kills)/float64(deaths)*100) / 100
		}))

		deathsLeaderboard.SetNameTag(formattedUserLeaderBoard("DEATHS", func(u data.User) int {
			return u.Teams.Stats.Deaths
		}))

		topPointsLeaderboard.SetNameTag(formattedTeamLeaderBoard("POINTS", func(t data.Team) int {
			return t.Points
		}))

		topKOTHLeaderboard.SetNameTag(formattedTeamLeaderBoard("KOTH CAPTURES", func(t data.Team) int {
			return t.KOTHWins
		}))

		topKillsLeaderboard.SetNameTag(formattedTeamLeaderBoard("KILLS", func(t data.Team) int {
			kills := 0
			for _, m := range t.Members {
				u, err := data.LoadUserFromName(m.Name)
				if err != nil {
					continue
				}
				kills += u.Teams.Stats.Kills
			}
			return kills
		}))*/
	}
}

func formattedTeamLeaderBoard[T int | float64](name string, value func(u data2.Team) T) string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><red>TOP TEAM %v</red></bold>\n", strings.ToUpper(name)))
	teams, err := data2.LoadAllTeams()
	if err != nil {
		return sb.String()
	}

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(teams, func(i int) string {
		return teams[i].Name
	})

	slices.SortFunc(teams, func(a, b data2.Team) int {
		if value(a) == value(b) {
			return 0
		}
		if value(a) > value(b) {
			return -1
		}
		return 1
	})

	for i := 0; i < 10; i++ {
		if len(teams) < i+1 {
			break
		}
		leader := teams[i]
		name := leader.DisplayName

		position, _ := roman.Itor(i + 1)
		sb.WriteString(text.Colourf(
			"<grey>%v.</grey> <white>%v</white> <dark-grey>-</dark-grey> <grey>%v</grey>\n",
			position,
			name,
			value(leader),
		))
	}
	return sb.String()
}

func formattedUserLeaderBoard[T int | float64](name string, value func(u data2.User) T) string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><red>TOP %v</red></bold>\n", strings.ToUpper(name)))
	users, err := data2.LoadAllUsers()
	if err != nil {
		return sb.String()
	}

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(users, func(i int) string {
		return users[i].Name
	})

	slices.SortFunc(users, func(a, b data2.User) int {
		if value(a) == value(b) {
			return 0
		}
		if value(a) > value(b) {
			return -1
		}
		return 1
	})

	for i := 0; i < 10; i++ {
		if len(users) < i+1 {
			break
		}
		leader := users[i]
		name := leader.DisplayName

		position, _ := roman.Itor(i + 1)
		sb.WriteString(text.Colourf(
			"<grey>%v.</grey> <white>%v</white> <dark-grey>-</dark-grey> <grey>%v</grey>\n",
			position,
			name,
			value(leader),
		))
	}
	return sb.String()
}

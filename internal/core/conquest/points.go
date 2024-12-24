package conquest

import (
	"github.com/moyai-network/teams/internal/core/data"
	"sort"
	"sync"
)

var (
	// points is a map that stores the points of each team. The key of the map is the name of the team, and the
	// value is the amount of points the team has.
	points   = map[string]int{}
	pointsMu sync.Mutex
)

// resetPoints resets the points of all teams to 0.
func resetPoints() {
	pointsMu.Lock()
	points = map[string]int{}
	pointsMu.Unlock()
}

// increaseTeamPoints increases the points of a team by n.
func IncreaseTeamPoints(team data.Team, n int) {
	pointsMu.Lock()
	points[team.Name] += n
	pointsMu.Unlock()
}

// OrderedTeamsByPoints returns a slice of all teams ordered by the amount of points they have. The team with the
// most points will be at the start of the slice.
func OrderedTeamsByPoints() []data.Team {
	pointsMu.Lock()
	tms, _ := data.LoadAllTeams()
    sort.SliceStable(tms, func(i, j int) bool {
        if points[tms[i].Name] != points[tms[j].Name] {
            return points[tms[i].Name] > points[tms[j].Name]
        }
        return tms[i].Name < tms[j].Name
    })
	pointsMu.Unlock()
	return tms
}

// LookupTeamPoints returns the amount of points a team has.
func LookupTeamPoints(team data.Team) int {
	pointsMu.Lock()
	defer pointsMu.Unlock()
	pts, ok := points[team.Name]
	if !ok {
		return 0
	}
	return pts
}

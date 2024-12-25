package model

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
	"strings"
	"time"
)

type Team struct {
	// Name is the identifier for the team data.
	Name string
	// DisplayName is the display name for the team.
	DisplayName string
	// Members is a slice of all the members in the team.
	Members []Member
	// DTR is the amount of deaths required before the teams goes raidable.
	DTR float64
	// Home is the home point for the team.
	Home mgl64.Vec3
	// LastDeath is the last time someone died in the team
	LastDeath time.Time
	// Points is the amount of points the team has.
	Points int
	// KOTHWins is the amount of KOTH wins the team has.
	KOTHWins int
	// ConquestWins is the amount of Conquest wins the team has.
	ConquestWins int
	// Balance is the amount of money the team has.
	Balance float64
	// Claim is the claim area of the team.
	Claim Area
	// Focus is the focus information for a team
	Focus Focus
	// Renamed is whether the team has been renamed.
	Renamed bool
}

// DefaultTeam returns a team with default values
func DefaultTeam(name string) Team {
	return Team{
		Name:        strings.ToLower(name),
		DTR:         1.01,
		DisplayName: name,
	}
}

// WithRename renames the team.
func (t Team) WithRename(name string) Team {
	t.DisplayName = name
	t.Name = strings.ToLower(name)
	t.Renamed = true
	return t
}

// WithMembers returns the team with the given members.
func (t Team) WithMembers(m ...Member) Team {
	t.Members = m
	return t
}

// WithoutMember returns the team without the given member
func (t Team) WithoutMember(name string) Team {
	for i, m := range t.Members {
		if strings.EqualFold(name, m.Name) {
			t.Members = append(t.Members[:i], t.Members[i+1:]...)
			break
		}
	}
	return t
}

// Promote promotes a member of the faction.
func (t Team) Promote(name string) Team {
	var m Member
	var l Member
	for _, me := range t.Members {
		if me.Name == name {
			m = me
		}
		if t.Leader(me.Name) {
			l = me
		}
	}
	switch m.Rank {
	case 1:
		m.Rank = 2
	case 2:
		l.Rank = 2
		m.Rank = 3
	}

	return t
}

// Demote demotes a member of the faction.
func (t Team) Demote(name string) Team {
	var m Member
	for _, me := range t.Members {
		if me.Name == name {
			m = me
		}
	}
	if m.Rank == 2 {
		m.Rank = 1
	}

	return t
}

// WithClaim returns the team with the given claim area.
func (t Team) WithClaim(claim Area) Team {
	t.Claim = claim
	return t
}

// WithPoints returns the team with the given amount of points
func (t Team) WithPoints(points int) Team {
	t.Points = points
	if t.Points < 0 {
		t.Points = 0
	}
	return t
}

// WithBalance returns the team with the given balance.
func (t Team) WithBalance(bal float64) Team {
	t.Balance = bal
	return t
}

// WithHome returns the team with the given home coordinate.
func (t Team) WithHome(home mgl64.Vec3) Team {
	t.Home = home
	return t
}

func (t Team) WithLastDeath(lastDeath time.Time) Team {
	t.LastDeath = lastDeath
	return t
}

// Frozen returns whether the team is frozen.
func (t Team) Frozen() bool {
	return time.Since(t.LastDeath) < time.Minute*15
}

// Regenerating returns whether the team is regenerating.
func (t Team) Regenerating() bool {
	return !eq(t.DTR, t.MaxDTR()) && !t.Frozen()
}

// WithDTR returns the team with the given dtr.
func (t Team) WithDTR(dtr float64) Team {
	t.DTR = dtr
	return t
}

// MaxDTR returns the max DTR of the faction.
func (t Team) MaxDTR() float64 {
	dtr := 1.01 * float64(len(t.Members))
	return math.Round(dtr*100) / 100
}

// DTRString returns the DTR string of the faction
func (t Team) DTRString() string {
	if t.Regenerating() {
		return text.Colourf("<purple>%.2f%s</purple>", t.DTR, t.DTRDot())
	}
	if t.DTR >= t.MaxDTR() {
		return text.Colourf("<green>%.2f%s</green>", t.DTR, t.DTRDot())
	}
	if t.DTR < 0 {
		return text.Colourf("<red>%.2f%s</red>", t.DTR, t.DTRDot())
	}
	return text.Colourf("<yellow>%.2f%s</yellow>", t.DTR, t.DTRDot())
}

// DTRDot returns the DTR dot of the faction.
func (t Team) DTRDot() string {
	if t.Regenerating() {
		return "<purple>■</purple>"
	}
	if eq(t.DTR, t.MaxDTR()) {
		return "<green>■</green>"
	}
	return "<red>■</red>"
}

// Leader returns whether the given username is the one of the leader.
func (t Team) Leader(name string) bool {
	for _, m := range t.Members {
		if strings.EqualFold(m.Name, name) && m.Rank == 3 {
			return true
		}
	}
	return false
}

// Captain returns whether the given username is the one of the captain.
func (t Team) Captain(name string) bool {
	for _, m := range t.Members {
		if strings.EqualFold(m.Name, name) && m.Rank == 2 {
			return true
		}
	}
	return false
}

// Member returns whether the given username is in the team.
func (t Team) Member(name string) bool {
	for _, m := range t.Members {
		if strings.EqualFold(m.Name, name) {
			return true
		}
	}
	return false
}

// WithTeamFocus returns the team with the given team as the focus.
func (t Team) WithTeamFocus(tm Team) Team {
	t.Focus.Kind = FocusTypeTeam
	t.Focus.Value = tm.DisplayName
	return t
}

// WithPlayerFocus returns the team with the given player as the focus.
func (t Team) WithPlayerFocus(name string) Team {
	t.Focus.Kind = FocusTypePlayer
	t.Focus.Value = name
	return t
}

// WithoutFocus returns the team without a focus.
func (t Team) WithoutFocus() Team {
	t.Focus.Kind = FocusTypeNone
	t.Focus.Value = ""
	return t
}

// Member represents a team member.
type Member struct {
	Name        string
	DisplayName string
	XUID        string
	Rank        int
}

// DefaultMember returns a default team member.
func DefaultMember(xuid, name string) Member {
	return Member{
		Name:        strings.ToLower(name),
		DisplayName: name,
		XUID:        xuid,
		Rank:        1,
	}
}

// WithRank returns a team member with the given rank.
func (m Member) WithRank(n int) Member {
	m.Rank = n
	return m
}

// Focus is the focus information for a team
type Focus struct {
	Kind  int    // 0:Player ; 1: Team
	Value string // XUID: Player ; Name: Team
}

var (
	FocusTypeNone   = 0
	FocusTypePlayer = 1
	FocusTypeTeam   = 2
)

func eq(a, b float64) bool {
	return math.Abs(a-b) <= 1e-5
}

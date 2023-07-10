package form

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/moyai-network/moose"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/moose/role"
	"golang.org/x/exp/maps"
)

// mute is a form that allows a user to issue a mute.
type mute struct {
	// Reason is a dropdown that allows a user to select a mute reason.
	Reason form.Dropdown
	// OnlinePlayer is a dropdown that allows a user to select an online player.
	OnlinePlayer form.Dropdown
	// OfflinePlayer is an input field that allows a user to enter an offline player.
	OfflinePlayer form.Input
	// online is a list of online players' XUIDs indexed by their names.
	online map[string]string
	// p is the player that is using the form.
	p *player.Player
}

// NewMute creates a new form to issue a mute.
func NewMute(p *player.Player) form.Form {
	online := make(map[string]string)
	for _, u := range user.All() {
		online[u.Player().Name()] = u.Player().Name()
	}
	names := [...]string{"Steve Harvey", "Elon Musk", "Bill Gates", "Mark Zuckerberg", "Jeff Bezos", "Warren Buffet", "Larry Page", "Sergey Brin", "Larry Ellison", "Tim Cook", "Steve Ballmer", "Daniel Larson", "Steve"}
	list := maps.Keys(online)
	sort.Strings(list)
	return form.New(mute{
		Reason:        form.NewDropdown("Reason", []string{"Spam", "Toxicity", "Advertisement"}, 0),
		OnlinePlayer:  form.NewDropdown("Online Player", list, 0),
		OfflinePlayer: form.NewInput("Offline Player", "", names[rand.Intn(len(names)-1)]),
		online:        online,
		p:             p,
	}, "Mute")
}

// Submit ...
func (m mute) Submit(s form.Submitter) {
	p := s.(*player.Player)
	u, err := data.LoadUserOrCreate(p.Name())
	if err != nil {
		// User somehow left midway through the form.
		return
	}

	h, ok := user.Lookup(p.Name())
	if !ok {
		// User somehow left midway through the form.
		return
	}

	if !u.Roles.Contains(role.Trial{}, role.Operator{}) {
		// In case the user's role was removed while the form was open.
		return
	}
	var length time.Duration
	reason := m.Reason.Options[m.Reason.Value()]
	switch reason {
	case "Spam":
		length = time.Hour * 6
	case "Toxicity":
		length = time.Hour * 9
	case "Advertising":
		length = time.Hour * 24 * 3
	default:
		panic("should never happen")
	}

	mu := moose.Punishment{
		Staff:      m.p.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Expiration: time.Now().Add(length),
	}
	if offlineName := strings.TrimSpace(m.OfflinePlayer.Value()); offlineName != "" {
		if strings.EqualFold(offlineName, m.p.Name()) {
			h.Message("command.mute.self")
			return
		}
		t, err := data.LoadUserOrCreate(offlineName)
		if err != nil {
			h.Message("command.target.unknown")
			return
		}
		if t.Roles.Contains(role.Operator{}) {
			h.Message("command.mute.operator")
			return
		}
		if !t.Mute.Expired() {
			h.Message("command.mute.already")
			return
		}
		t.Mute = mu
		_ = data.SaveUser(t)

		user.Alert(m.p, "staff.alert.mute", t.DisplayName, reason)
		//webhook.SendPunishment(m.p.Name(), t.DisplayName(), reason, "Mute")
		h.Message("command.mute.success", t.DisplayName, reason)
		return
	}
	t, err := data.LoadUserOrCreate(m.online[m.OnlinePlayer.Options[m.OnlinePlayer.Value()]])
	if err != nil {
		h.Message("command.target.unknown")
		return
	}
	if t.Roles.Contains(role.Operator{}) {
		h.Message("command.mute.operator")
		return
	}
	if !t.Mute.Expired() {
		h.Message("command.mute.already")
		return
	}
	t.Mute = mu
	_ = data.SaveUser(t) // Save in case of a server crash or anything that may cause the data to not get saved.

	user.Alert(m.p, "staff.alert.mute", t.Name, reason)
	//webhook.SendPunishment(m.p.Name(), t.Player().Name(), reason, "Mute")
	h.Message("command.mute.success", t.Name, reason)
}

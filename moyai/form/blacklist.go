package form

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/hako/durafmt"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"golang.org/x/exp/maps"
)

// blacklist is a form that allows a user to issue a blacklist.
type blacklist struct {
	// Reason is a dropdown that allows a user to select a blacklist reason.
	Reason form.Input
	// OnlinePlayer is a dropdown that allows a user to select an online player.
	OnlinePlayer form.Dropdown
	// OfflinePlayer is an input field that allows a user to enter an offline player.
	OfflinePlayer form.Input
	// online is a list of online players' XUIDs indexed by their names.
	online map[string]string
}

// NewBlacklist creates a new form to issue a blacklist.
func NewBlacklist() form.Form {
	online := make(map[string]string)
	for _, u := range user.All() {
		online[u.Player().Name()] = u.Player().Name()
	}
	names := [...]string{"Steve Harvey", "Elon Musk", "Bill Gates", "Mark Zuckerberg", "Jeff Bezos", "Warren Buffet", "Larry Page", "Sergey Brin", "Larry Ellison", "Tim Cook", "Steve Ballmer", "Daniel Larson", "Steve"}
	list := maps.Keys(online)
	sort.Strings(list)
	return form.New(blacklist{
		Reason:        form.NewInput("Reason", "", "Enter a reason for the blacklist."),
		OnlinePlayer:  form.NewDropdown("Online Player", list, 0),
		OfflinePlayer: form.NewInput("Offline Player", "", names[rand.Intn(len(names)-1)]),
		online:        online,
	}, "Blacklist")
}

// Submit ...
func (b blacklist) Submit(s form.Submitter) {
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
	reason := b.Reason.Value()

	punishment := moose.Punishment{
		Staff:      p.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Permanent:  true,
	}
	var name string
	if offlineName := strings.TrimSpace(b.OfflinePlayer.Value()); offlineName != "" {
		if strings.EqualFold(offlineName, p.Name()) {
			h.Message("command.blacklist.self")
			return
		}
		t, err := data.LoadUserOrCreate(offlineName)
		if err != nil {
			h.Message("command.target.unknown")
			return
		}
		if t.Roles.Contains(role.Operator{}) {
			h.Message("command.blacklist.operator")
			return
		}
		if !t.Ban.Expired() {
			h.Message("command.blacklist.already")
			return
		}
		t.Ban = punishment
		_ = data.SaveUser(t)

		name = t.DisplayName
	} else {
		t, err := data.LoadUserOrCreate(b.online[b.OnlinePlayer.Options[b.OnlinePlayer.Value()]])
		if err != nil {
			h.Message("command.target.unknown")
			return
		}
		if t.Roles.Contains(role.Operator{}) {
			h.Message("command.blacklist.operator`")
			return
		}

		tH, ok := user.Lookup(t.Name)
		t.Ban = punishment

		if ok {
			tH.Player().Disconnect(strings.Join([]string{
				lang.Translatef(tH.Player().Locale(), "user.blacklist.header"),
				lang.Translatef(tH.Player().Locale(), "user.blacklist.description", reason, durafmt.ParseShort(length)),
			}, "\n"))
		}
		name = t.DisplayName

		_ = data.SaveUser(t)
	}

	_ = data.SaveUser(u) // Save in case of a server crash or anything that may cause the data to not get saved.

	user.Alert(p, "staff.alert.blacklist", name, reason)
	user.Broadcast("command.blacklist.broadcast", p.Name(), name, reason)
	//webhook.SendPunishment(p.Name(), name, reason, "Ban")
	h.Message("command.blacklist.success", name, reason)
}

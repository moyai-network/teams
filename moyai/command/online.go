package command

import (
	"fmt"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/moose/data"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/user"
)

// Online is a command that displays the number of players online and their names.
type Online struct{}

// Run ...
func (Online) Run(s cmd.Source, o *cmd.Output) {
	var users []string
	for _, u := range user.All() {
		d, err := data.LoadUserOrCreate(u.Player().Name())
		if err != nil {
			o.Print(lang.Translatef(locale(s), "target.data.load.error", u.Player().Name()))
			return
		}
		name := d.Name
		if name != u.Player().Name() {
			name += fmt.Sprintf("(%s)", u.Player().Name())
		}
		highest := d.Roles.Highest()
		users = append(users, highest.Colour(u.Player().Name()))
	}
	o.Printf(lang.Translatef(locale(s), "command.online.users", len(users), strings.Join(users, ", ")))
}

// Allow ...
func (Online) Allow(s cmd.Source) bool {
	return allow(s, true)
}

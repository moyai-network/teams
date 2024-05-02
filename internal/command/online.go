package command

import (
	"fmt"
	"github.com/moyai-network/teams/moyai"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/moose/data"
	"github.com/moyai-network/moose/lang"
)

// Online is a command that displays the number of players online and their names.
type Online struct{}

// Run ...
func (Online) Run(s cmd.Source, o *cmd.Output) {
	var users []string
	for _, u := range moyai.Server().Players() {
		d, err := data.LoadUserOrCreate(u.Name())
		if err != nil {
			o.Print(lang.Translatef(locale(s), "target.data.load.error", u.Name()))
			return
		}
		name := d.Name
		if name != u.Name() {
			name += fmt.Sprintf("(%s)", u.Name())
		}
		highest := d.Roles.Highest()
		users = append(users, highest.Colour(u.Name()))
	}
	o.Printf(lang.Translatef(locale(s), "command.online.users", len(users), strings.Join(users, ", ")))
}

// Allow ...
func (Online) Allow(s cmd.Source) bool {
	return allow(s, true)
}

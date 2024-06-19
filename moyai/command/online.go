package command

import (
	"fmt"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
)

// Online is a command that displays the number of players online and their names.
type Online struct{}

// Run ...
func (Online) Run(s cmd.Source, o *cmd.Output) {
	var users []string
	for _, u := range moyai.Players() {
		d, err := data.LoadUserFromName(u.Name())
		if err != nil {
			o.Print(lang.Translatef(locale(s), "target.data.load.error", u.Name()))
			return
		}
		name := d.Name
		if name != u.Name() {
			name += fmt.Sprintf("(%s)", u.Name())
		}
		highest := d.Roles.Highest()
		users = append(users, highest.Color(u.Name()))
	}
	o.Printf(lang.Translatef(locale(s), "command.online.users", len(users), strings.Join(users, ", ")))
}

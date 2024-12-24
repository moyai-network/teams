package command

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/pkg/lang"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
)

// Online is a command that displays the number of players online and their names.
type Online struct{}

// Run ...
func (Online) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	var users []string
	for u := range internal.Players(tx) {
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
		users = append(users, highest.Coloured(u.Name()))
	}
	o.Printf(lang.Translatef(locale(s), "command.online.users", len(users), strings.Join(users, ", ")))
}

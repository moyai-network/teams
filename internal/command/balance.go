package command

import (
	"github.com/moyai-network/teams/internal/data"
	"github.com/moyai-network/teams/internal/role"
	"github.com/moyai-network/teams/internal/user"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type Balance struct{}

func (Balance) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	user.Messagef(p, "command.balance.self", u.Teams.Balance)
}

type BalancePayOnline struct {
	Sub    cmd.SubCommand `cmd:"pay"`
	Target []cmd.Target   `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalancePayOnline) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	t, ok := b.Target[0].(*player.Player)
	if !ok {
		return
	}

	if t == p {
		user.Messagef(p, "command.pay.self")
		return
	}

	target, err := data.LoadUserFromName(t.Name())
	if err != nil {
		return
	}

	if b.Amount < 0 {
		user.Messagef(p, "command.pay.negative")
		return
	}
	if u.Teams.Balance < b.Amount {
		user.Messagef(p, "command.pay.insufficient")
		return
	}

	u.Teams.Balance -= b.Amount
	target.Teams.Balance += b.Amount

	data.SaveUser(u)
	data.SaveUser(target)

	user.Messagef(t, "command.pay.receiver", u.Roles.Highest().Color(p.Name()), 0)
	user.Messagef(p, "command.pay.sender", target.Roles.Highest().Color(t.Name()), 0)
}

type BalancePayOffline struct {
	Sub    cmd.SubCommand `cmd:"pay"`
	Target string         `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalancePayOffline) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	if b.Amount < 0 {
		user.Messagef(p, "command.pay.negative")
		return
	}

	if u.Teams.Balance < b.Amount {
		user.Messagef(p, "command.pay.insufficient")
		return
	}

	if strings.EqualFold(b.Target, p.Name()) {
		user.Messagef(p, "command.pay.self")
		return
	}

	t, err := data.LoadUserFromName(p.Name())
	if err != nil {
		out.Error("Unexpected error occurred. Please contact an administrator.")
		return
	}

	u.Teams.Balance -= b.Amount
	t.Teams.Balance += b.Amount

	data.SaveUser(u)
	data.SaveUser(t)

	user.Messagef(p, "command.pay.sender", t.Roles.Highest().Color(t.DisplayName), b.Amount)
}

type BalanceAdd struct {
	Sub    cmd.SubCommand `cmd:"add"`
	Target []cmd.Target   `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalanceAdd) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	t, ok := b.Target[0].(*player.Player)
	if !ok {
		return
	}

	if b.Amount < 0 {
		user.Messagef(p, "command.add.negative")
		return
	}

	target, err := data.LoadUserFromName(t.Name())
	if err != nil {
		return
	}

	target.Teams.Balance += b.Amount

	data.SaveUser(target)

	user.Messagef(t, "command.add.receiver", u.Roles.Highest().Color(p.Name()), b.Amount)
	user.Messagef(p, "command.add.sender", target.Roles.Highest().Color(t.Name()), b.Amount)
}

type BalanceAddOffline struct {
	Sub    cmd.SubCommand `cmd:"add"`
	Target string         `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalanceAddOffline) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	if b.Amount < 0 {
		user.Messagef(p, "command.add.negative")
		return
	}

	t, err := data.LoadUserFromName(b.Target)
	if err != nil {
		out.Error("Unexpected error occurred. Please contact an administrator.")
		return
	}

	t.Teams.Balance += b.Amount

	data.SaveUser(t)

	user.Messagef(p, "command.add.sender", t.Roles.Highest().Color(t.DisplayName), b.Amount)
}

func (BalanceAdd) Allow(src cmd.Source) bool {
	return allow(src, true, role.Operator{})
}

package command

import (
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/data"
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

	internal.Messagef(p, "command.balance.self", u.Teams.Balance)
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
		internal.Messagef(p, "command.pay.self")
		return
	}

	target, err := data.LoadUserFromName(t.Name())
	if err != nil {
		return
	}

	if b.Amount < 0 {
		internal.Messagef(p, "command.pay.negative")
		return
	}
	if u.Teams.Balance < b.Amount {
		internal.Messagef(p, "command.pay.insufficient")
		return
	}

	u.Teams.Balance -= b.Amount
	target.Teams.Balance += b.Amount

	data.SaveUser(u)
	data.SaveUser(target)

	internal.Messagef(t, "command.add.receiver", u.Roles.Highest().Coloured(p.Name()), 0)
	internal.Messagef(p, "command.add.sender", target.Roles.Highest().Coloured(t.Name()), 0)
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
		internal.Messagef(p, "command.pay.negative")
		return
	}

	if u.Teams.Balance < b.Amount {
		internal.Messagef(p, "command.pay.insufficient")
		return
	}

	if strings.EqualFold(b.Target, p.Name()) {
		internal.Messagef(p, "command.pay.self")
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

	internal.Messagef(p, "command.add.sender", t.Roles.Highest().Coloured(t.DisplayName), b.Amount)
}

type BalanceAdd struct {
	operatorAllower
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
		internal.Messagef(p, "command.add.negative")
		return
	}

	target, err := data.LoadUserFromName(t.Name())
	if err != nil {
		return
	}

	target.Teams.Balance += b.Amount

	data.SaveUser(target)

	internal.Messagef(t, "command.add.receiver", u.Roles.Highest().Coloured(p.Name()), b.Amount)
	internal.Messagef(p, "command.add.sender", target.Roles.Highest().Coloured(t.Name()), b.Amount)
}

type BalanceAddOffline struct {
	operatorAllower
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
		internal.Messagef(p, "command.add.negative")
		return
	}

	t, err := data.LoadUserFromName(b.Target)
	if err != nil {
		out.Error("Unexpected error occurred. Please contact an administrator.")
		return
	}

	t.Teams.Balance += b.Amount

	data.SaveUser(t)

	internal.Messagef(p, "command.add.sender", t.Roles.Highest().Coloured(t.DisplayName), b.Amount)
}

package command

import (
	"github.com/moyai-network/moose/lang"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/data"
)

type Balance struct{}

func (Balance) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserOrCreate(p.Name())
	if err != nil {
		return
	}

	out.Print(lang.Translatef(u.Language(), "command.balance.self", u.GameMode.Teams.Balance))
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
	u, err := data.LoadUserOrCreate(p.Name())
	if err != nil {
		return
	}

	t, ok := b.Target[0].(*player.Player)
	if !ok {
		return
	}

	if t == p {
		out.Print(lang.Translatef(u.Language(), "command.pay.self"))
		return
	}

	target, err := data.LoadUserOrCreate(t.Name())
	if err != nil {
		return
	}

	if b.Amount < 0 {
		out.Print(lang.Translatef(u.Language(), "command.pay.negative"))
		return
	}
	if u.GameMode.Teams.Balance < b.Amount {
		out.Print(lang.Translatef(u.Language(), "command.pay.insufficient"))
		return
	}

	u.GameMode.Teams.Balance -= b.Amount
	target.GameMode.Teams.Balance += b.Amount

	_ = data.SaveUser(u)
	_ = data.SaveUser(target)

	t.Message(lang.Translatef(target.Language(), "command.pay.receiver", u.Roles.Highest().Colour(p.Name()), 0))
	out.Print(lang.Translatef(u.Language(), "command.pay.sender", target.Roles.Highest().Colour(t.Name()), 0))
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
	u, err := data.LoadUserOrCreate(p.Name())
	if err != nil {
		return
	}

	if b.Amount < 0 {
		out.Print(lang.Translatef(u.Language(), "command.pay.negative"))
		return
	}

	if u.GameMode.Teams.Balance < b.Amount {
		out.Print(lang.Translatef(u.Language(), "command.pay.insufficient"))
		return
	}

	if strings.EqualFold(b.Target, p.Name()) {
		out.Print(lang.Translatef(u.Language(), "command.pay.self"))
		return
	}

	t, err := data.LoadUserOrCreate(p.Name())
	if err != nil {
		out.Error("Unexpected error occurred. Please contact an administrator.")
		return
	}

	u.GameMode.Teams.Balance -= b.Amount
	t.GameMode.Teams.Balance += b.Amount

	_ = data.SaveUser(u)
	_ = data.SaveUser(t)

	out.Print(lang.Translatef(u.Language(), "command.pay.sender", t.Roles.Highest().Colour(t.DisplayName), b.Amount))
}

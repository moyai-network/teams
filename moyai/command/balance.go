package command

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Balance struct{}

func (Balance) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}

	p.Message(text.Colourf("<green>Your balance is $%2.f.</green>", u.Balance))
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
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}

	t, ok := b.Target[0].(*player.Player)
	if !ok {
		return
	}

	if t == p {
		out.Error("You cannot pay yourself.")
		return
	}

	target, err := data.LoadUser(t.Name())
	if err != nil {
		return
	}

	if b.Amount < 0 {
		out.Error("You cannot pay a negative amount.")
		return
	}
	if u.Balance < b.Amount {
		out.Error("You do not have enough money.")
		return
	}

	u.Balance -= b.Amount
	target.Balance += b.Amount

	_ = data.SaveUser(u)
	_ = data.SaveUser(target)

	t.Message(text.Colourf("<green>%s has paid you $%2.f.</green>", u.Roles.Highest().Colour(p.Name()), 0))
	p.Message(text.Colourf("<green>You have paid %s $%2.f.</green>", target.Roles.Highest().Colour(t.Name()), 0))
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
	u, err := data.LoadUser(p.Name())
	if err != nil {
		return
	}

	if b.Amount < 0 {
		out.Error("You cannot pay a negative amount.")
		return
	}

	if u.Balance < b.Amount {
		out.Error("You do not have enough money.")
		return
	}

	if strings.EqualFold(b.Target, p.Name()) {
		out.Error("You cannot pay yourself.")
		return
	}

	t, err := data.LoadUser(p.Name())
	if err != nil {
		out.Error("user has never joined the server")
		return
	}

	u.Balance -= b.Amount
	t.Balance += b.Amount

	_ = data.SaveUser(u)
	_ = data.SaveUser(t)

	p.Message(text.Colourf("<green>You have paid</green> %s <green>$%2.f.</green>", t.Roles.Highest().Colour(t.DisplayName), b.Amount))
}

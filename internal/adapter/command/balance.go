package command

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type Balance struct{}

func (Balance) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}

	internal.Messagef(p, "command.balance.self", u.Teams.Balance)
}

type BalancePayOnline struct {
	Sub    cmd.SubCommand `cmd:"pay"`
	Target []cmd.Target   `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalancePayOnline) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
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

	target, ok := core.UserRepository.FindByName(t.Name())
	if !ok {
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

	core.UserRepository.Save(u)
	core.UserRepository.Save(target)

	internal.Messagef(t, "command.add.receiver", u.Roles.Highest().Coloured(p.Name()), 0)
	internal.Messagef(p, "command.add.sender", target.Roles.Highest().Coloured(t.Name()), 0)
}

type BalancePayOffline struct {
	Sub    cmd.SubCommand `cmd:"pay"`
	Target string         `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalancePayOffline) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
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

	t, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		out.Error("Unexpected error occurred. Please contact an administrator.")
		return
	}

	u.Teams.Balance -= b.Amount
	t.Teams.Balance += b.Amount

	core.UserRepository.Save(u)
	core.UserRepository.Save(t)

	internal.Messagef(p, "command.add.sender", t.Roles.Highest().Coloured(t.DisplayName), b.Amount)
}

type BalanceAdd struct {
	operatorAllower
	Sub    cmd.SubCommand `cmd:"add"`
	Target []cmd.Target   `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalanceAdd) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
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

	target, ok := core.UserRepository.FindByName(t.Name())
	if !ok {
		return
	}

	target.Teams.Balance += b.Amount

	core.UserRepository.Save(target)

	internal.Messagef(t, "command.add.receiver", u.Roles.Highest().Coloured(p.Name()), b.Amount)
	internal.Messagef(p, "command.add.sender", target.Roles.Highest().Coloured(t.Name()), b.Amount)
}

type BalanceAddOffline struct {
	operatorAllower
	Sub    cmd.SubCommand `cmd:"add"`
	Target string         `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalanceAddOffline) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	if b.Amount < 0 {
		internal.Messagef(p, "command.add.negative")
		return
	}

	t, ok := core.UserRepository.FindByName(b.Target)
	if !ok {
		out.Error("Unexpected error occurred. Please contact an administrator.")
		return
	}

	t.Teams.Balance += b.Amount

	core.UserRepository.Save(t)

	internal.Messagef(p, "command.add.sender", t.Roles.Highest().Coloured(t.DisplayName), b.Amount)
}

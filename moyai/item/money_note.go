package item

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func NewMoneyNote(val float64, n int) item.Stack {
	return item.NewStack(item.Paper{}, n).WithValue("MONEY_NOTE", val).WithCustomName(text.Colourf("§r<amethyst>$%.0f Money Note</amethyst>", val))
}

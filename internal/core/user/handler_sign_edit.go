package user

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/core"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/roles"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) HandleSignEdit(ctx *player.Context, pos cube.Pos, frontSide bool, oldText, newText string) {
	p := ctx.Val()
	if !frontSide {
		return
	}

	u, err := data2.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	teams := core.TeamRepository.FindAll()
	if posWithinProtectedArea(p, pos, teams) {
		return
	}

	lines := strings.Split(newText, "\n")
	if len(lines) <= 0 {
		return
	}

	switch formatRegex.ReplaceAllString(strings.ToLower(lines[0]), "") {
	case "[elevator]":
		if len(lines) < 2 {
			return
		}
		var newLines []string
		newLines = append(newLines, text.Colourf("<dark-red>[Elevator]</dark-red>"))
		switch strings.ToLower(lines[1]) {
		case "up":
			newLines = append(newLines, text.Colourf("Up"))
		case "down":
			newLines = append(newLines, text.Colourf("Down"))
		default:
			return
		}
		time.AfterFunc(time.Millisecond, func() {
			b := p.Tx().Block(pos)
			if s, ok := b.(block.Sign); ok {
				p.Tx().SetBlock(pos, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
	case "[buy]", "[sell]":
		if u.Roles.Contains(roles.Operator()) {
			return
		}
		time.AfterFunc(time.Millisecond, func() {
			b := p.Tx().Block(pos)
			if s, ok := b.(block.Sign); ok {
				p.Tx().SetBlock(pos, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       "",
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
		return
	case "[shop]":
		if len(lines) < 4 {
			return
		}

		if !u.Roles.Contains(roles.Admin()) {
			p.Tx().SetBlock(pos, block.Air{}, nil)
			return
		}

		var newLines []string
		spl := strings.Split(lines[1], " ")
		if len(spl) < 2 {
			return
		}
		choice := strings.ToLower(spl[0])
		q, _ := strconv.Atoi(spl[1])
		price, _ := strconv.Atoi(lines[3])
		switch choice {
		case "buy":
			newLines = append(newLines, text.Colourf("<green>[Buy]</green>"))
		case "sell":
			newLines = append(newLines, text.Colourf("<red>[Sell]</red>"))
		}

		newLines = append(newLines, formatItemName(lines[2]))
		newLines = append(newLines, fmt.Sprint(q))
		newLines = append(newLines, fmt.Sprintf("$%d", price))

		time.AfterFunc(time.Millisecond, func() {
			b := p.Tx().Block(pos)
			if s, ok := b.(block.Sign); ok {
				p.Tx().SetBlock(pos, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
		/*case "[kit]":
		return // disabled for HCF
		if len(lines) < 2 {
			return
		}

		if !u.Roles.Contains(roles.Admin()) {
			p.World().SetBlock(pos, block.Air{}, nil)
			return
		}

		var newLines []string
		newLines = append(newLines, text.Colourf("<dark-red>[Kit]</dark-red>"))
		newLines = append(newLines, text.Colourf("%s", lines[1]))

		time.AfterFunc(time.Millisecond, func() {
			b := p.World().Block(pos)
			if s, ok := b.(block.Sign); ok {
				p.World().SetBlock(pos, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})*/
	}
}

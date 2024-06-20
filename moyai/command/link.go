package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"math/rand"
	"strings"
)

type Link struct{}

func (Link) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromXUID(p.XUID())
	if err != nil {
		moyai.Messagef(p, "user.data.load.error")
		return
	}

	if u.DiscordID != "" {
		moyai.Messagef(p, "command.link.already")
		return
	}

	code := generateCode()
	u.LinkCode = code
	data.SaveUser(u)
	moyai.Messagef(p, "command.link.code", code)
}

var codeChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func generateCode() string {
	var s strings.Builder
	for i := 0; i < 6; i++ {
		s.WriteByte(codeChars[rand.Intn(len(codeChars))])
	}
	_, err := data.LoadUserFromCode(s.String())
	for err == nil {
		s.Reset()
		for i := 0; i < 6; i++ {
			s.WriteByte(codeChars[rand.Intn(len(codeChars))])
		}
		_, err = data.LoadUserFromCode(s.String())
	}

	return s.String()
}

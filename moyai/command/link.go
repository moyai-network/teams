package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"math/rand"
	"strings"
)

type Unlink struct{}

func (Unlink) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.data.load.error")
		return
	}

	err = data.UnlinkUser(u, moyai.DiscordState(), discord.GuildID(1111055709300342826))
	if err != nil {
		return
	}
	moyai.Messagef(p, "command.unlink.done")
}

type Link struct{}

func (Link) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		moyai.Messagef(p, "user.data.load.error")
		return
	}

	if u.DiscordID != "" {
		moyai.Messagef(p, "command.link.already")
		return
	}
	code := u.LinkCode
	if code == "" {
		code = generateCode()
		u.LinkCode = code
		data.SaveUser(u)
	}

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

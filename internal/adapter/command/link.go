package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/data"
	"math/rand"
	"strings"
)

type Unlink struct{}

func (Unlink) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		internal.Messagef(p, "user.data.load.error")
		return
	}

	err = data.UnlinkUser(u, internal.DiscordState(), discord.GuildID(1111055709300342826))
	if err != nil {
		return
	}
	internal.Messagef(p, "command.unlink.done")
}

type Link struct{}

func (Link) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		internal.Messagef(p, "user.data.load.error")
		return
	}

	if u.DiscordID != "" {
		internal.Messagef(p, "command.link.already")
		return
	}
	code := u.LinkCode
	if code == "" {
		code = generateCode()
		u.LinkCode = code
		data.SaveUser(u)
	}

	internal.Messagef(p, "command.link.code", code)
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

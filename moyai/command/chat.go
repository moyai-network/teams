package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/user"
	"time"
)

// ChatMute is a command that is used to mute the global chat.
type ChatMute struct {
	adminAllower
	Sub cmd.SubCommand `cmd:"mute"`
}

// ChatUnMute is a command that is used to unmute the global chat.
type ChatUnMute struct {
	adminAllower
	Sub cmd.SubCommand `cmd:"unmute"`
}

// ChatCoolDown is a command that is used to set the global chat cooldown.
type ChatCoolDown struct {
	adminAllower
	Sub      cmd.SubCommand `cmd:"cooldown"`
	CoolDown int            `cmd:"seconds"`
}

// Run ...
func (c ChatMute) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	if !moyai.GlobalChatEnabled() {
		user.Messagef(p, "chat.global.already-muted")
		return
	}
	moyai.ToggleGlobalChat()

	user.Messagef(p, "chat.global.mute")
}

// Run ...
func (c ChatUnMute) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	if moyai.GlobalChatEnabled() {
		user.Messagef(p, "chat.global.already-unmuted")
		return
	}
	moyai.ToggleGlobalChat()

	user.Messagef(p, "chat.global.unmute")

}

// Run ...
func (c ChatCoolDown) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	moyai.UpdateChatCoolDown(time.Duration(c.CoolDown) * time.Second)
	user.Messagef(p, "chat.cooldown.updated", c.CoolDown)
}

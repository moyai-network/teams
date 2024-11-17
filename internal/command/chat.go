package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal"
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

	if !internal.GlobalChatEnabled() {
		internal.Messagef(p, "chat.global.already-muted")
		return
	}
	internal.ToggleGlobalChat()

	internal.Messagef(p, "chat.global.mute")
}

// Run ...
func (c ChatUnMute) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	if internal.GlobalChatEnabled() {
		internal.Messagef(p, "chat.global.already-unmuted")
		return
	}
	internal.ToggleGlobalChat()

	internal.Messagef(p, "chat.global.unmute")

}

// Run ...
func (c ChatCoolDown) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	internal.UpdateChatCoolDown(time.Duration(c.CoolDown) * time.Second)
	internal.Messagef(p, "chat.cooldown.updated", c.CoolDown)
}

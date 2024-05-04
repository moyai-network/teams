package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/role"
	"github.com/moyai-network/teams/internal/user"
	"github.com/moyai-network/teams/moyai"
)

// vanishGameMode is the game mode used by vanished players.
type vanishGameMode struct{}

func (vanishGameMode) AllowsEditing() bool      { return true }
func (vanishGameMode) AllowsTakingDamage() bool { return false }
func (vanishGameMode) CreativeInventory() bool  { return false }
func (vanishGameMode) HasCollision() bool       { return false }
func (vanishGameMode) AllowsFlying() bool       { return true }
func (vanishGameMode) AllowsInteraction() bool  { return true }
func (vanishGameMode) Visible() bool            { return true }

// Vanish is a command that hides a staff from everyone else.
type Vanish struct{}

// Run ...
func (Vanish) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	if h.Vanished() {
		//user.Alert(s, "staff.alert.vanish.off")
		p.SetGameMode(world.GameModeSurvival)
		user.Messagef(p, "command.vanish.disabled")
	} else {
		//user.Alert(s, "staff.alert.vanish.on")
		p.SetGameMode(world.GameModeSpectator)
		user.Messagef(p, "command.vanish.enabled")
	}
	for _, t := range moyai.Server().Players() {
		if !h.Vanished() {
			t.HideEntity(p)
			continue
		}
		t.ShowEntity(p)
	}
	h.ToggleVanish()
}

// Allow ...
func (Vanish) Allow(s cmd.Source) bool {
	return allow(s, false, role.Trial{})
}

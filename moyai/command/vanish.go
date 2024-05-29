package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/user"
)

// vanishGameMode is the game mode used by vanished players.
type vanishGameMode struct {
	lastMode world.GameMode
}

func (vanishGameMode) AllowsEditing() bool      { return true }
func (vanishGameMode) AllowsTakingDamage() bool { return false }
func (vanishGameMode) CreativeInventory() bool  { return true }
func (vanishGameMode) HasCollision() bool       { return true }
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

	u, err := data.LoadUserFromXUID(p.XUID())
	if err != nil {
		return
	}
	mode := p.GameMode()

	if u.Vanished {
		//user.Alertf(s, "staff.alert.vanish.off")
		vanishMode, ok := mode.(vanishGameMode)
		if !ok {
			return
		}
		p.SetGameMode(vanishMode.lastMode)
		user.Messagef(p, "command.vanish.disabled")
	} else {
		//user.Alertf(s, "staff.alert.vanish.on")
		p.SetGameMode(vanishGameMode{lastMode: mode})
		user.Messagef(p, "command.vanish.enabled")
	}

	user.ToggleVanish(p, u)
}

// Allow ...
func (Vanish) Allow(s cmd.Source) bool {
	return allow(s, false, role.Trial{})
}

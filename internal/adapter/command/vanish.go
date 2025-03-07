package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/user"
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
type Vanish struct {
	trialAllower
}

// Run ...
func (Vanish) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	mode := p.GameMode()

	if u.Vanished {
		//internal.Alertf(s, "staff.alert.vanish.off")
		internal.Messagef(p, "command.vanish.disabled")
		vanishMode, ok := mode.(vanishGameMode)
		if ok {
			p.SetGameMode(vanishMode.lastMode)
		}
	} else {
		//internal.Alertf(s, "staff.alert.vanish.on")
		p.SetGameMode(vanishGameMode{lastMode: mode})
		internal.Messagef(p, "command.vanish.enabled")
	}

	u.Vanished = !u.Vanished
	core.UserRepository.Save(u)
	user.UpdateVanishState(p, u)
}

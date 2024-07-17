package user

import (
	"time"

	"github.com/bedrock-gophers/cooldown/cooldown"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/process"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Logout is a process that handles the logout of a player.
func (h *Handler) Logout() *process.Process {
	return h.processLogout
}

// Stuck is a process that handles the stuck command.
func (h *Handler) Stuck() *process.Process {
	return h.processStuck
}

// BeginCamp is a process that handles the camp command
func (h *Handler) BeginCamp(tm data.Team, pos cube.Pos) {
	h.processCamp = process.NewProcess(func(t *process.Process) {
		h.p.Message(text.Colourf("<green>You have been teleported close to %s's home.</green>", tm.DisplayName))
	})
	h.processCamp.Teleport(h.p, time.Second*45, pos.Vec3(), moyai.Overworld())
}

// CampOngoing returns true if the camp process is ongoing
func (h *Handler) CampOngoing() bool {
	if h.processCamp == nil {
		return false
	}
	return h.processCamp.Ongoing()
}

// Home is a process that handles the home command.
func (h *Handler) Home() *process.Process {
	return h.processHome
}

// Combat is a cooldown that handles the combat cooldown.
func (h *Handler) Combat() *cooldown.CoolDown {
	return h.tagCombat
}

// Pearl is a cooldown that handles the ender pearl cooldown.
func (h *Handler) Pearl() *cooldown.CoolDown {
	return h.coolDownPearl
}

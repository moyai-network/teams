package user

import (
	"github.com/moyai-network/teams/internal/cooldown"
	"github.com/moyai-network/teams/moyai/process"
)

// Logout is a process that handles the processLogout of a player.
func (h *Handler) Logout() *process.Process {
	return h.processLogout
}

// Stuck is a process that handles the processStuck command.
func (h *Handler) Stuck() *process.Process {
	return h.processStuck
}

// Home is a process that handles the processHome command.
func (h *Handler) Home() *process.Process {
	return h.processHome
}

// Combat is a cooldown that handles the tagCombat cooldown.
func (h *Handler) Combat() *cooldown.CoolDown {
	return h.tagCombat
}

// Pearl is a cooldown that handles the ender coolDownPearl cooldown.
func (h *Handler) Pearl() *cooldown.CoolDown {
	return h.coolDownPearl
}

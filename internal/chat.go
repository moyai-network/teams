package internal

import "time"

var (
	// globalChatEnabled is a boolean that determines whether the global chat is enabled.
	globalChatEnabled = true
	// globalChatCoolDown is an integer that determines the global chat cooldown.
	globalChatCoolDown = 3 * time.Second
)

// ToggleGlobalChat toggles the global chat mute.
func ToggleGlobalChat() {
	globalChatEnabled = !globalChatEnabled
}

// UpdateChatCoolDown sets the global chat cooldown.
func UpdateChatCoolDown(seconds time.Duration) {
	globalChatCoolDown = seconds
}

// GlobalChatEnabled returns whether the global chat is enabled.
func GlobalChatEnabled() bool {
	return globalChatEnabled
}

// ChatCoolDown returns the global chat cooldown.
func ChatCoolDown() time.Duration {
	return globalChatCoolDown
}

package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Manager represents the role specification for the manager role.
type Manager struct{}

// Name returns the name of the role.
func (Manager) Name() string {
	return "manager"
}

// Chat returns the formatted chat message using the name and message provided.
func (Manager) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-red>Manager</dark-red>]</grey> <dark-red>%s</dark-red><dark-grey>:</dark-grey> <dark-red>%s</dark-red>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Manager) Color(name string) string {
	return text.Colourf("<dark-red>%s</dark-red>", name)
}

// Inherits returns the role that this role inherits from.
func (Manager) Inherits() Role {
	return Admin{}
}

package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Mod represents the role specification for the mod role.
type Mod struct{}

// Name returns the name of the role.
func (Mod) Name() string {
	return "mod"
}

// Chat returns the formatted chat message using the name and message provided.
func (Mod) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-green>Mod</dark-green>]</grey> <dark-green>%s</dark-green><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Mod) Color(name string) string {
	return text.Colourf("<dark-green>%s</dark-green>", name)
}

// Inherits returns the role that this role inherits from.
func (Mod) Inherits() Role {
	return Trial{}
}

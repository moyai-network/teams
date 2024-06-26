package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Pharaoh represents the role specification for the Pharaoh role.
type Pharaoh struct{}

// Name returns the name of the role.
func (Pharaoh) Name() string {
	return "pharaoh"
}

// Chat returns the formatted chat message using the name and message provided.
func (Pharaoh) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-red>Pharaoh</dark-red>]</grey> <dark-red>%s</dark-red><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Pharaoh) Color(name string) string {
	return text.Colourf("<dark-red>%s</dark-red>", name)
}

// Inherits returns the role that this role inherits from.
func (Pharaoh) Inherits() Role {
	return Menes{}
}

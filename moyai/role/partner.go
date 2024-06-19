package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Partner represents the role specification for the Partner role.
type Partner struct{}

// Name returns the name of the role.
func (Partner) Name() string {
	return "partner"
}

// Chat returns the formatted chat message using the name and message provided.
func (Partner) Chat(name, message string) string {
	return text.Colourf("<grey>[<aqua><i>Partner</i></aqua>]</grey> <aqua>%s</aqua><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Partner) Color(name string) string {
	return text.Colourf("<aqua>%s</aqua>", name)
}

// Inherits returns the role that this role inherits from.
func (Partner) Inherits() Role {
	return Menes{}
}

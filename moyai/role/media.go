package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Media represents the role specification for the Media role.
type Media struct{}

// Name returns the name of the role.
func (Media) Name() string {
	return "media"
}

// Chat returns the formatted chat message using the name and message provided.
func (Media) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-aqua><i>Media</i></dark-aqua>]</grey> <dark-aqua>%s</dark-aqua><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Media) Color(name string) string {
	return text.Colourf("<dark-aqua>%s</dark-aqua>", name)
}

// Inherits returns the role that this role inherits from.
func (Media) Inherits() Role {
	return Pharaoh{}
}

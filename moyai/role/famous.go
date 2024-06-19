package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Famous represents the role specification for the Famous role.
type Famous struct{}

// Name returns the name of the role.
func (Famous) Name() string {
	return "famous"
}

// Chat returns the formatted chat message using the name and message provided.
func (Famous) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-purple><i>Famous</i></dark-purple>]</grey> <dark-purple>%s</dark-purple><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Famous) Color(name string) string {
	return text.Colourf("<dark-purple>%s</dark-purple>", name)
}

// Inherits returns the role that this role inherits from.
func (Famous) Inherits() Role {
	return Pharaoh{}
}

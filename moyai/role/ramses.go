package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Ramses represents the role specification for the Ramses role.
type Ramses struct{}

// Name returns the name of the role.
func (Ramses) Name() string {
	return "ramses"
}

// Chat returns the formatted chat message using the name and message provided.
func (Ramses) Chat(name, message string) string {
	return text.Colourf("<grey>[<gold>Ramses</gold>]</grey> <gold>%s</gold><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Ramses) Color(name string) string {
	return text.Colourf("<gold>%s</gold>", name)
}

// Inherits returns the role that this role inherits from.
func (Ramses) Inherits() Role {
	return Khufu{}
}

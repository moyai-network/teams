package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Menes represents the role specification for the Menes role.
type Menes struct{}

// Name returns the name of the role.
func (Menes) Name() string {
	return "Menes"
}

// Chat returns the formatted chat message using the name and message provided.
func (Menes) Chat(name, message string) string {
	return text.Colourf("<grey>[<purple>Menes</purple>]</grey> <purple>%s</purple><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Menes) Color(name string) string {
	return text.Colourf("<purple>%s</purple>", name)
}

// Inherits returns the role that this role inherits from.
func (Menes) Inherits() Role {
	return Ramses{}
}

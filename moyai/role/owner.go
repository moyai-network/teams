package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Owner represents the role specification for the owner role.
type Owner struct{}

// Name returns the name of the role.
func (Owner) Name() string {
	return "owner"
}

// Chat returns the formatted chat message using the name and message provided.
func (Owner) Chat(name, message string) string {
	return text.Colourf("<grey>[<black>Owner</black>]</grey> <black>%s</black><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Owner) Color(name string) string {
	return text.Colourf("<black>%s</black>", name)
}

// Inherits returns the role that this role inherits from.
func (Owner) Inherits() Role {
	return Admin{}
}

package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Donor4 represents the role specification for the Donor4 role.
type Donor4 struct{}

// Name returns the name of the role.
func (Donor4) Name() string {
	return "Donor4"
}

// Chat returns the formatted chat message using the name and message provided.
func (Donor4) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-green>Donor4</dark-green>]</grey> <dark-green>%s</dark-green><dark-grey>:</dark-grey> <dark-green>%s</dark-green>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Donor4) Color(name string) string {
	return text.Colourf("<dark-green>%s</dark-green>", name)
}

// Inherits returns the role that this role inherits from.
func (Donor4) Inherits() Role {
	return Donor3{}
}

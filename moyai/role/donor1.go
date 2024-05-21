package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Donor1 represents the role specification for the Donor1 role.
type Donor1 struct{}

// Name returns the name of the role.
func (Donor1) Name() string {
	return "Donor1"
}

// Chat returns the formatted chat message using the name and message provided.
func (Donor1) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-green>Donor1</dark-green>]</grey> <dark-green>%s</dark-green><dark-grey>:</dark-grey> <dark-green>%s</dark-green>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Donor1) Color(name string) string {
	return text.Colourf("<dark-green>%s</dark-green>", name)
}
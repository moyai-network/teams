package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Donor2 represents the role specification for the Donor2 role.
type Donor2 struct{}

// Name returns the name of the role.
func (Donor2) Name() string {
	return "Donor2"
}

// Chat returns the formatted chat message using the name and message provided.
func (Donor2) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-green>Donor2</dark-green>]</grey> <dark-green>%s</dark-green><dark-grey>:</dark-grey> <dark-green>%s</dark-green>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Donor2) Color(name string) string {
	return text.Colourf("<dark-green>%s</dark-green>", name)
}
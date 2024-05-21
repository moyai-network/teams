package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Donor3 represents the role specification for the Donor3 role.
type Donor3 struct{}

// Name returns the name of the role.
func (Donor3) Name() string {
	return "Donor3"
}

// Chat returns the formatted chat message using the name and message provided.
func (Donor3) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-green>Donor3</dark-green>]</grey> <dark-green>%s</dark-green><dark-grey>:</dark-grey> <dark-green>%s</dark-green>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Donor3) Color(name string) string {
	return text.Colourf("<dark-green>%s</dark-green>", name)
}
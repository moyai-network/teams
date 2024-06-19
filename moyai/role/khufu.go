package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Khufu represents the role specification for the Khufu role.
type Khufu struct{}

// Name returns the name of the role.
func (Khufu) Name() string {
	return "Khufu"
}

// Chat returns the formatted chat message using the name and message provided.
func (Khufu) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-blue>Khufu</dark-blue>]</grey> <dark-blue>%s</dark-blue><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Khufu) Color(name string) string {
	return text.Colourf("<dark-blue>%s</dark-blue>", name)
}
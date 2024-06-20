package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Nitro represents the role specification for the nitro role.
type Nitro struct{}

// Name returns the name of the role.
func (Nitro) Name() string {
	return "nitro"
}

// Chat returns the formatted chat message using the name and message provided.
func (Nitro) Chat(name, message string) string {
	return text.Colourf("<grey>[<purple>Nitro</purple>]</grey> <purple>%s</purple><grey>:</grey> <white>%s</white>", name, message)
}

// Colour returns the formatted name-Colour using the name provided.
func (Nitro) Color(name string) string {
	return text.Colourf("<purple>%s</purple>", name)
}
package role

import "github.com/sandertv/gophertunnel/minecraft/text"

// Trial represents the role specification for the trial role.
type Trial struct{}

// Name returns the name of the role.
func (Trial) Name() string {
	return "trial"
}

// Chat returns the formatted chat message using the name and message provided.
func (Trial) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-yellow>Trial</dark-yellow>]</grey> <dark-yellow>%s</dark-yellow><dark-grey>:</dark-grey> <white>%s</white>", name, message)
}

// Color returns the formatted name-Colour using the name provided.
func (Trial) Color(name string) string {
	return text.Colourf("<dark-yellow>%s</dark-yellow>", name)
}

func (Trial) Inherits() Role {
	return nil
}

package konsole

import (
	"github.com/bedrock-gophers/konsole/konsole"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type formatter struct {
	konsole.NopFormatter
}

func (formatter) FormatMessage(s string) string {
	return text.Colourf("<purple>[CONSOLE]: %s</purple>", s)
}

func (formatter) FormatAlert(s string) string {
	return text.Colourf("<b><red>[</red><yellow>ALERT<red>]:</red> <yellow>%s</yellow></b>", s)
}

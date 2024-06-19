package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type Canada struct{}

func (Canada) Name() string {
	return "canada"
}

func (Canada) Format() string {
	return text.Colourf("<grey>[<red>CA</red><white>NA</white><red>DA</red>]</grey>")
}

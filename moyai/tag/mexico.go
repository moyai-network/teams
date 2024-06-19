package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type Mexico struct{}

func (Mexico) Name() string {
	return "mexico"
}

func (Mexico) Format() string {
	return text.Colourf("<grey>[<green>ME</green><white>XI</white><red>CO</red>]</grey>")
}
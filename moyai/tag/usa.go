package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type USA struct{}

func (USA) Name() string {
	return "usa"
}

func (USA) Format() string {
	return text.Colourf("<grey>[<red>U</red><white>S</white><blue>A</blue>]</grey>")
}
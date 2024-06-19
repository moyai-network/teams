package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type Russia struct{}

func (Russia) Name() string {
	return "russia"
}

func (Russia) Format() string {
	return text.Colourf("<grey>[<white>RU</white><blue>SS</blue><red>IA</red>]</grey>")
}

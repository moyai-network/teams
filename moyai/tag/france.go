package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type France struct{}

func (France) Name() string {
	return "france"
}

func (France) Format() string {
	return text.Colourf("<grey>[<blue>FR</blue><white>AN</white><red>CE</red>]</grey>")
}

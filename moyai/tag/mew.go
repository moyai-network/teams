package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type Mew struct{}

func (Mew) Name() string {
	return "mew"
}

func (Mew) Format() string {
	return text.Colourf("<dark-purple>mew</dark-purple>")
}
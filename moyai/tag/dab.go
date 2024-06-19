package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type Dab struct{}

func (Dab) Name() string {
	return "dab"
}

func (Dab) Format() string {
	return text.Colourf("<green><o/</green>")
}
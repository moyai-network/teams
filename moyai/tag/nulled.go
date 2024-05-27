package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type Nulled struct{}

func (Nulled) Name() string {
	return "nulled"
}

func (Nulled) Format() string {
	return text.Colourf("<dark-red>NULLED</dark-red>")
}
package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type PVP struct{}

func (PVP) Name() string {
	return "pvp"
}

func (PVP) Format() string {
	return text.Colourf("<dark-aqua>PVP</dark-aqua>")
}
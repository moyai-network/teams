package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type UWU struct{}

func (UWU) Name() string {
	return "uwu"
}

func (UWU) Format() string {
	return text.Colourf("<purple>uWu</purple>")
}
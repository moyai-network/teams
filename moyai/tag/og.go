package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type OG struct{}

func (OG) Name() string {
	return "og"
}

func (OG) Format() string {
	return text.Colourf("<white>O</white><gold>G</gold>")
}
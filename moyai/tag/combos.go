package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type Combos struct{}

func (Combos) Name() string {
	return "combos"
}

func (Combos) Format() string {
	return text.Colourf("<gold>Combos</gold>")
}
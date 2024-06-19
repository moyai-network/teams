package tag

import "github.com/sandertv/gophertunnel/minecraft/text"

type Akhi struct{}

func (Akhi) Name() string {
	return "akhi"
}

func (Akhi) Format() string {
	return text.Colourf("<dark-green>Akhi</dark-green>")
}
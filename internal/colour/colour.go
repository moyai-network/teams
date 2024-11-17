package colour

import (
	"github.com/df-mc/dragonfly/server/item"
	"math/rand"
	"regexp"
)

var (
	stripRegex = regexp.MustCompile(`§[\da-gk-or]`)

	colours = []item.Colour{
		item.ColourWhite(),
		item.ColourOrange(),
		item.ColourMagenta(),
		item.ColourLightBlue(),
		item.ColourYellow(),
		item.ColourLime(),
		item.ColourPink(),
		item.ColourGrey(),
		item.ColourLightGrey(),
		item.ColourCyan(),
		item.ColourPurple(),
		item.ColourBlue(),
		item.ColourBrown(),
		item.ColourGreen(),
		item.ColourRed(),
		item.ColourBlack(),
	}
)

func Random() item.Colour {
	return colours[rand.Intn(len(colours))]
}

func Strip(s string) string {
	return stripRegex.ReplaceAllString(s, "")
}

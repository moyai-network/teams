package cape

import (
	"image/png"
	"os"

	"github.com/df-mc/dragonfly/server/player/skin"
)

// Cape stores data about a cosmetic cape, such as it's name, image, etc.
type Cape struct {
	name string
	plus bool
	cape skin.Cape
}

// NewCape creates a new cape with the given name and image. The plus parameter determines whether the cape can only
// be accessed by players with the Plus role.
func NewCape(name string, path string, plus bool) Cape {
	return Cape{
		name: name,
		plus: plus,
		cape: read("assets/capes/" + path),
	}
}

// Name returns the name of the cape.
func (c Cape) Name() string {
	return c.name
}

// Premium returns true if the cape can only be accessed by players with the Plus role.
func (c Cape) Premium() bool {
	return c.plus
}

// Cape returns the image data of the cape.
func (c Cape) Cape() skin.Cape {
	return c.cape
}

// read performs a read on the path provided and returns a dragonfly cape.
func read(path string) skin.Cape {
	f, _ := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	defer f.Close()
	i, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	c := skin.NewCape(i.Bounds().Max.X, i.Bounds().Max.Y)
	for y := 0; y < i.Bounds().Max.Y; y++ {
		for x := 0; x < i.Bounds().Max.X; x++ {
			color := i.At(x, y)
			r, g, b, a := color.RGBA()
			i := x*4 + i.Bounds().Max.X*y*4
			c.Pix[i], c.Pix[i+1], c.Pix[i+2], c.Pix[i+3] = uint8(r), uint8(g), uint8(b), uint8(a)
		}
	}
	return c
}

var (
	// capes contains all registered capes.
	capes []Cape
	// capesByName maps a cape's name to the cape itself.
	capesByName = map[string]Cape{}
)

// All returns all registered capes.
func All() []Cape {
	return capes
}

// Register registers the cape provided.
func Register(cape Cape) {
	capes = append(capes, cape)
	capesByName[cape.Name()] = cape
}

// ByName returns the cape with the given name. If no cape with the given name is registered, the second return value is
// false.
func ByName(name string) (Cape, bool) {
	cape, ok := capesByName[name]
	return cape, ok
}

// init registers all capes in the capes folder.
func init() {
	Register(NewCape("Vasar Series 1", "regular/vasar_series_one.png", false))
	Register(NewCape("Vasar Series 2", "regular/vasar_series_two.png", false))

	Register(NewCape("Portal", "plus/portal.png", true))
	Register(NewCape("Spicy", "plus/spicy.png", true))
	Register(NewCape("Happy", "plus/happy.png", true))
	Register(NewCape("Sad", "plus/sad.png", true))
	Register(NewCape("Fold", "plus/fold.png", true))

	Register(NewCape("Optifine Red", "plus/optifine_red.png", true))
	Register(NewCape("Optifine Pink", "plus/optifine_pink.png", true))
	Register(NewCape("Optifine Cyan", "plus/optifine_cyan.png", true))
	Register(NewCape("Optifine Green", "plus/optifine_green.png", true))
	Register(NewCape("Optifine Dark", "plus/optifine_dark.png", true))

	Register(NewCape("Vlone", "plus/vlone.png", true))
	Register(NewCape("Crow", "plus/crow.png", true))
	Register(NewCape("AK", "plus/ak.png", true))
	Register(NewCape("Evil", "plus/evil.png", true))
	Register(NewCape("Drag", "plus/drag.png", true))
	Register(NewCape("Bland", "plus/bland.png", true))
	Register(NewCape("Pumpkin", "plus/pumpkin.png", true))
	Register(NewCape("Wave", "plus/wave.png", true))
	Register(NewCape("CBA", "plus/cba.png", true))
	Register(NewCape("Cookie", "plus/cookie.png", true))
}

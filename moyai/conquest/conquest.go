package conquest

import (
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func Start() {
	running = true
}

func Stop() {
	running = false
}

var running bool

var (
	Blue = &Conquest{
		name:        text.Colourf("<blue>Blue</blue>"),
		cancel:      make(chan struct{}),
		area:        area.NewArea(mgl64.Vec2{89, 436}, mgl64.Vec2{-95, 430}),
		coordinates: mgl64.Vec2{500, 500},
		duration:    time.Second * 30,
	}
	Green = &Conquest{
		name:        text.Colourf("<green>Green</green>"),
		cancel:      make(chan struct{}),
		area:        area.NewArea(mgl64.Vec2{-222, 431}, mgl64.Vec2{-228, 437}),
		coordinates: mgl64.Vec2{500, 500},
		duration:    time.Second * 30,
	}
	Yellow = &Conquest{
		name:        text.Colourf("<yellow>Yellow</yellow>"),
		cancel:      make(chan struct{}),
		area:        area.NewArea(mgl64.Vec2{-221, 564}, mgl64.Vec2{-227, 570}),
		coordinates: mgl64.Vec2{500, 500},
		duration:    time.Second * 30,
	}
	Red = &Conquest{
		name:        text.Colourf("<red>Red</red>"),
		cancel:      make(chan struct{}),
		area:        area.NewArea(mgl64.Vec2{-94, 563}, mgl64.Vec2{-88, 569}),
		coordinates: mgl64.Vec2{500, 500},
		duration:    time.Second * 30,
	}
)

func All() []*Conquest {
	return []*Conquest{Blue, Green, Yellow, Red}
}

func Running() bool {
	return running
}

func Lookup(name string) (*Conquest, bool) {
	for _, c := range All() {
		if strings.EqualFold(colour.StripMinecraftColour(c.Name()), name) {
			return c, true
		}
	}
	return nil, false
}

// Conquest refers to a cap for the conquest nether event.
type Conquest struct {
	name        string
	capturing   *player.Player
	time        time.Time
	cancel      chan struct{}
	area        area.Area
	coordinates mgl64.Vec2
	duration    time.Duration
}

// Name returns the name of the Conquest.
func (c *Conquest) Name() string {
	return c.name
}

func (c *Conquest) start() {
	c.capturing = nil
	c.cancel = make(chan struct{})
}

func (c *Conquest) stop() {
	c.capturing = nil
	c.time = time.Time{}
	close(c.cancel)
}

func (c *Conquest) PlayerCapturing(p *player.Player) bool  {
	return c.capturing == p
}

func (c *Conquest) Capturing() (*player.Player, bool) {
	return c.capturing, c.capturing != nil
}

func (c *Conquest) StartCapturing(p *player.Player) bool {
	if c.capturing != nil || !running {
		return false
	}
	t := c.duration
	c.time = time.Now().Add(t)
	go func() {
		select {
		case <-time.After(t):
			c.capturing = nil
			//TODO: c.running = false

			// u, err := data.LoadUserFromName(p.Name())
			// if err != nil {
			// 	c.StopCapturing(p)
			// 	return
			// }
			// tm, err := data.LoadTeamFromMemberName(u.Name)
			// if err != nil {
			// 	c.StopCapturing(p)
			// 	return
			// }
			// TODO: Add conquest points
		case <-c.cancel:
			c.capturing = nil
			return
		}
	}()
	c.capturing = p
	return true
}

// StopCapturing stops the capturing of the KOTH.
func (c *Conquest) StopCapturing(p *player.Player) bool {
	if !running {
		return false
	}
	if c.capturing == p {
		c.capturing = nil
		c.cancel <- struct{}{}
		return true
	}
	return false
}

// Time returns the time at which the KOTH will be captured.
func (c *Conquest) Time() time.Time {
	return c.time
}

// Area returns the area of the KOTH.
func (c *Conquest) Area() area.Area {
	return c.area
}

// Coordinates returns the coordinates of the KOTH.
func (c *Conquest) Coordinates() mgl64.Vec2 {
	return c.coordinates
}

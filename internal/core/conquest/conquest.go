package conquest

import (
	"github.com/moyai-network/teams/internal/core/area"
	"github.com/moyai-network/teams/internal/core/colour"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/ports/model"
	"strings"
	"sync/atomic"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var (
	// running is true if the Conquest event is running.
	running atomic.Bool

	Red = &Conquest{
		name:        text.Colourf("<red>Red Zone</red>"),
		cancel:      make(chan struct{}),
		area:        area.NewArea(mgl64.Vec2{-94, 563}, mgl64.Vec2{-88, 569}),
		coordinates: mgl64.Vec2{500, 500},
		duration:    time.Second * 30,
	}
	Blue = &Conquest{
		name:        text.Colourf("<blue>Blue Zone</blue>"),
		cancel:      make(chan struct{}),
		area:        area.NewArea(mgl64.Vec2{-89, 436}, mgl64.Vec2{-95, 430}),
		coordinates: mgl64.Vec2{500, 500},
		duration:    time.Second * 30,
	}
	Green = &Conquest{
		name:        text.Colourf("<green>Green Zone</green>"),
		cancel:      make(chan struct{}),
		area:        area.NewArea(mgl64.Vec2{-222, 431}, mgl64.Vec2{-228, 437}),
		coordinates: mgl64.Vec2{500, 500},
		duration:    time.Second * 30,
	}
	Yellow = &Conquest{
		name:        text.Colourf("<yellow>Yellow Zone</yellow>"),
		cancel:      make(chan struct{}),
		area:        area.NewArea(mgl64.Vec2{-221, 564}, mgl64.Vec2{-227, 570}),
		coordinates: mgl64.Vec2{500, 500},
		duration:    time.Second * 30,
	}
)

// Start starts all Conquest events.
func Start() {
	running.Store(true)
	for _, c := range All() {
		c.start()
	}
}

// Stop stops all Conquest events.
func Stop() {
	running.Store(false)
	for _, c := range All() {
		c.stop()
	}
}

// All returns all Conquest events.
func All() []*Conquest {
	return []*Conquest{Red, Blue, Green, Yellow}
}

// Running returns true if the Conquest event is running.
func Running() bool {
	return running.Load()
}

// Lookup returns a Conquest by its name. The second return value is true if the Conquest was found, and false
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

// Duration returns the duration of the Conquest.
func (c *Conquest) Duration() time.Duration {
	return c.duration
}

// start starts the capturing of the Conquest.
func (c *Conquest) start() {
	c.capturing = nil
	c.cancel = make(chan struct{})
}

// stop stops the capturing of the Conquest.
func (c *Conquest) stop() {
	c.capturing = nil
	c.time = time.Time{}
	close(c.cancel)
}

// PlayerCapturing returns true if the player is capturing the Conquest.
func (c *Conquest) PlayerCapturing(p *player.Player) bool {
	return c.capturing == p
}

// Capturing returns the player capturing the Conquest, if any.
func (c *Conquest) Capturing() (*player.Player, bool) {
	return c.capturing, c.capturing != nil
}

// StartCapturing starts the capturing of the Conquest.
func (c *Conquest) StartCapturing(p *player.Player) bool {
	if c.capturing != nil || !running.Load() {
		return false
	}
	t := c.duration
	c.time = time.Now().Add(t)
	go func() {
		select {
		case <-time.After(t):
			c.capturing = nil
			u, err := data2.LoadUserFromName(p.Name())
			if err != nil {
				c.StopCapturing(p)
				return
			}
			tm, err := data2.LoadTeamFromMemberName(u.Name)
			if err != nil {
				c.StopCapturing(p)
				return
			}
			IncreaseTeamPoints(tm, 10)
			_, _ = chat.Global.WriteString(lang.Translatef(model.Language{}, "conquest.captured", c.Name(), u.Roles.Highest().Coloured(u.DisplayName)))
			c.StopCapturing(p)

			pts := LookupTeamPoints(tm)
			if pts >= 150 {
				_, _ = chat.Global.WriteString(lang.Translatef(model.Language{}, "conquest.won", tm.Name, pts))
				for _, m := range tm.Members {
					for p := range internal.Players(p.Tx()) {
						if p.Name() == m.DisplayName {
							item.AddOrDrop(p, item.NewKey(item.KeyTypeConquest, 2))
						}
					}
				}
				tm.Points += 20
				data2.SaveTeam(tm)
				resetPoints()
				Stop()
			}

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
	if !running.Load() {
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

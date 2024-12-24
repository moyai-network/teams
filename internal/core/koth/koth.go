package koth

import (
	"github.com/moyai-network/teams/internal/core/area"
	"github.com/moyai-network/teams/internal/core/colour"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/ports/model"
	"math/rand"
	"strings"
	"time"

	"github.com/moyai-network/teams/pkg/lang"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func init() {
	t := time.NewTicker(time.Hour * 4)
	n := rand.Intn(len(All()[:2]))
	k := All()[n]
	_, _ = chat.Global.WriteString(text.Colourf("%s is now active!", k.name))
	k.Start()
	go func() {
		for range t.C {
			if _, ok := Running(); ok {
				continue
			}
			n := rand.Intn(len(All()))
			k := All()[n]
			_, _ = chat.Global.WriteString(text.Colourf("%s is now active!", k.name))
			k.Start()
		}
	}()
}

var (
	Garden = &KOTH{
		name:        text.Colourf("<dark-green>Garden</dark-green>"),
		dimension:   world.Overworld,
		area:        area.NewArea(mgl64.Vec2{-497, -503}, mgl64.Vec2{-503, -497}),
		cancel:      make(chan struct{}),
		coordinates: mgl64.Vec2{-500, -500},
		duration:    time.Minute * 5,
	}
	Oasis = &KOTH{
		name:        text.Colourf("<red>Oasis</red>"),
		dimension:   world.Overworld,
		area:        area.NewArea(mgl64.Vec2{503, 497}, mgl64.Vec2{497, 503}),
		cancel:      make(chan struct{}),
		coordinates: mgl64.Vec2{500, 500},
		duration:    time.Minute * 5,
	}
	Shrine = &KOTH{
		name:        text.Colourf("<gold>Mesa</gold>"),
		dimension:   world.Overworld,
		area:        area.NewArea(mgl64.Vec2{503, -503}, mgl64.Vec2{497, -497}),
		cancel:      make(chan struct{}),
		coordinates: mgl64.Vec2{500, -500},
		duration:    time.Minute * 5,
	}
	Kingdom = &KOTH{
		name:        text.Colourf("<aqua>Kingdom</aqua>"),
		dimension:   world.Overworld,
		area:        area.NewArea(mgl64.Vec2{-503, 497}, mgl64.Vec2{-497, 503}),
		cancel:      make(chan struct{}),
		coordinates: mgl64.Vec2{-500, 500},
		duration:    time.Minute * 5,
	}
	Nether = &KOTH{
		name:        text.Colourf("<red>Nether</red>"),
		dimension:   world.Nether,
		area:        area.NewArea(mgl64.Vec2{492, -104}, mgl64.Vec2{498, -110}),
		cancel:      make(chan struct{}),
		coordinates: mgl64.Vec2{500, -100},
		duration:    time.Minute * 5,
	}
	Citadel = &KOTH{
		name:        text.Colourf("<dark-red>Citadel</dark-red>"),
		dimension:   world.Nether,
		area:        area.NewArea(mgl64.Vec2{180, -504}, mgl64.Vec2{188, -496}),
		cancel:      make(chan struct{}),
		coordinates: mgl64.Vec2{200, -500},
		duration:    time.Minute * 10,
	}
	End = &KOTH{
		name:        text.Colourf("<purple>End</purple>"),
		dimension:   world.End,
		area:        area.NewArea(mgl64.Vec2{-11, 87}, mgl64.Vec2{-17, 93}),
		cancel:      make(chan struct{}),
		coordinates: mgl64.Vec2{-14, 90},
		duration:    time.Minute * 5,
	}
)

// All returns all KOTHs.
func All() []*KOTH {
	return []*KOTH{Garden, Oasis, Shrine, Kingdom, Citadel, End, Nether}
}

// Running returns true if the KOTH passed is currently running.
func Running() (*KOTH, bool) {
	for _, k := range All() {
		if k.running {
			return k, true
		}
	}
	return nil, false
}

// Lookup returns a KOTH by its name.
func Lookup(name string) (*KOTH, bool) {
	for _, k := range All() {
		if strings.EqualFold(colour.StripMinecraftColour(k.Name()), name) {
			return k, true
		}
	}
	return nil, false
}

// KOTH represents a King of the Hill event.
type KOTH struct {
	name        string
	dimension   world.Dimension
	capturing   *player.Player
	running     bool
	time        time.Time
	cancel      chan struct{}
	area        area.Area
	coordinates mgl64.Vec2
	duration    time.Duration
}

// Name returns the name of the KOTH.
func (k *KOTH) Name() string {
	return k.name
}

// Dimension returns the dimension of the KOTH.
func (k *KOTH) Dimension() world.Dimension {
	return k.dimension
}

// Start starts the KOTH.
func (k *KOTH) Start() {
	k.running = true
	k.capturing = nil
	k.cancel = make(chan struct{})
}

// Duration returns the duration of the KOTH.
func (k *KOTH) Duration() time.Duration {
	return k.duration
}

// Stop stops the KOTH.
func (k *KOTH) Stop() {
	k.running = false
	k.capturing = nil
	k.time = time.Time{}
	close(k.cancel)
}

// PlayerCapturing returns true if the player passed is currently capturing the KOTH.
func (k *KOTH) PlayerCapturing(p *player.Player) bool {
	return k.capturing == p
}

// Capturing returns the player that is currently capturing the KOTH, if any.
func (k *KOTH) Capturing() (*player.Player, bool) {
	return k.capturing, k.capturing != nil
}

// StartCapturing starts the capturing of the KOTH.
func (k *KOTH) StartCapturing(p *player.Player) bool {
	if k.capturing != nil || !k.running {
		return false
	}
	t := k.Duration()
	k.time = time.Now().Add(t)
	go func() {
		select {
		case <-time.After(t):
			k.capturing = nil
			k.running = false

			u, err := data2.LoadUserFromName(p.Name())
			if err != nil {
				k.StopCapturing(p)
				return
			}
			tm, err := data2.LoadTeamFromMemberName(u.Name)
			if err != nil {
				k.StopCapturing(p)
				return
			}
			tm.Points += 10
			tm.KOTHWins++
			data2.SaveTeam(tm)

			_, _ = chat.Global.WriteString(lang.Translatef(model.Language{}, "koth.captured", k.Name(), u.Roles.Highest().Coloured(u.DisplayName)))
			item.AddOrDrop(p, item.NewKey(item.KeyTypeKOTH, 2))
		case <-k.cancel:
			k.capturing = nil
			return
		}
	}()
	k.capturing = p
	return true
}

// StopCapturing stops the capturing of the KOTH.
func (k *KOTH) StopCapturing(p *player.Player) bool {
	if !k.running {
		return false
	}
	if k.capturing == p {
		k.capturing = nil
		k.cancel <- struct{}{}
		return true
	}
	return false
}

// Time returns the time at which the KOTH will be captured.
func (k *KOTH) Time() time.Time {
	return k.time
}

// Area returns the area of the KOTH.
func (k *KOTH) Area() area.Area {
	return k.area
}

// Coordinates returns the coordinates of the KOTH.
func (k *KOTH) Coordinates() mgl64.Vec2 {
	return k.coordinates
}

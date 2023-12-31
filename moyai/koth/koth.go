package koth

import (
	"math/rand"
	"strings"
	"time"

	"github.com/moyai-network/moose/crate"
	"github.com/moyai-network/teams/moyai/data"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func init() {
	t := time.NewTicker(time.Hour * 4)
	n := rand.Intn(len(All()))
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

type User interface {
	Player() *player.Player
	AddItemOrDrop(it item.Stack)
}

var (
	Broadcast = func(format string, a ...interface{}) {}
	Spiral    = &KOTH{
		name:        text.Colourf("<gold>Spiral</gold>"),
		area:        moose.NewArea(mgl64.Vec2{159, 98}, mgl64.Vec2{165, 104}),
		cancel:      make(chan struct{}, 0),
		coordinates: mgl64.Vec2{162, 101},
	}
	Circle = &KOTH{
		name:        text.Colourf("<dark-green>Circle</dark-green>"),
		area:        moose.NewArea(mgl64.Vec2{-88, 114}, mgl64.Vec2{-118, 85}),
		cancel:      make(chan struct{}, 0),
		coordinates: mgl64.Vec2{-103, 100},
	}
	Dragon = &KOTH{
		name:        text.Colourf("<amethyst>Dragon</amethyst>"),
		area:        moose.NewArea(mgl64.Vec2{63, 124}, mgl64.Vec2{57, 118}),
		cancel:      make(chan struct{}, 0),
		coordinates: mgl64.Vec2{60, 121},
	}
	Stairs = &KOTH{
		name:        text.Colourf("<aqua>Stairs</aqua>"),
		area:        moose.NewArea(mgl64.Vec2{0, 182}, mgl64.Vec2{-4, 178}),
		cancel:      make(chan struct{}, 0),
		coordinates: mgl64.Vec2{-2, -180},
	}
)

// All returns all KOTHs.
func All() []*KOTH {
	return []*KOTH{Spiral, Circle, Dragon, Stairs}
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
		if strings.EqualFold(moose.StripMinecraftColour(k.Name()), name) {
			return k, true
		}
	}
	return nil, false
}

// KOTH represents a King of the Hill event.
type KOTH struct {
	name        string
	capturing   User
	running     bool
	time        time.Time
	cancel      chan struct{}
	area        moose.Area
	coordinates mgl64.Vec2
}

// Name returns the name of the KOTH.
func (k *KOTH) Name() string {
	return k.name
}

// Start starts the KOTH.
func (k *KOTH) Start() {
	k.running = true
	k.capturing = nil
	k.cancel = make(chan struct{}, 0)
}

// Stop stops the KOTH.
func (k *KOTH) Stop() {
	k.running = false
	k.capturing = nil
	k.time = time.Time{}
	close(k.cancel)
}

// IsCapturing returns true if the player passed is currently capturing the KOTH.
func (k *KOTH) IsCapturing(u User) bool {
	return k.capturing == u
}

// Capturing returns the player that is currently capturing the KOTH, if any.
func (k *KOTH) Capturing() (User, bool) {
	return k.capturing, k.capturing != nil
}

// StartCapturing starts the capturing of the KOTH.
func (k *KOTH) StartCapturing(p User) bool {
	if k.capturing != nil || !k.running {
		return false
	}
	k.time = time.Now().Add(300 * time.Second)
	go func() {
		select {
		case <-time.After(300 * time.Second):
			k.capturing = nil
			k.running = false

			u, err := data.LoadUserOrCreate(p.Player().Name())
			if err != nil {
				k.StopCapturing(p)
				return
			}
			if t, ok := u.Team(); ok {
				t = t.WithPoints(10)
				data.SaveTeam(t)
			}

			Broadcast("koth.captured", k.Name(), u.Roles.Highest().Colour(u.DisplayName))
			kothCrate, ok := crate.ByName("KOTH")
			if !ok {
				panic("no crate found by the name of KOTH")
			}
			p.AddItemOrDrop(item.NewStack(item.Dye{Colour: item.ColourRed()}, 3).WithValue(kothCrate.Name(), true).WithCustomName(text.Colourf("<red>KOTH Crate Key</red>")))
		case <-k.cancel:
			k.capturing = nil
			return
		}
	}()
	k.capturing = p
	return true
}

// StopCapturing stops the capturing of the KOTH.
func (k *KOTH) StopCapturing(u User) bool {
	if !k.running {
		return false
	}
	if k.capturing == u {
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
func (k *KOTH) Area() moose.Area {
	return k.area
}

// Coordinates returns the coordinates of the KOTH.
func (k *KOTH) Coordinates() mgl64.Vec2 {
	return k.coordinates
}

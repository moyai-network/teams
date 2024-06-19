package cooldown

import (
	"github.com/df-mc/atomic"
	"github.com/rcrowley/go-bson"
	"time"
)

// CoolDown represents a time cooldown.
type CoolDown struct {
	expiration atomic.Value[time.Time]

	paused           atomic.Bool
	remainingAtPause atomic.Value[time.Duration]
}

// NewCoolDown returns a new process.
func NewCoolDown() *CoolDown {
	return &CoolDown{}
}

// TogglePause toggles the pause state of the cooldown.
func (c *CoolDown) TogglePause() {
	if c == nil {
		return
	}
	if !c.paused.Load() {
		c.remainingAtPause.Store(c.Remaining())
	} else {
		c.expiration = *atomic.NewValue(time.Now().Add(c.remainingAtPause.Load()))
	}

	c.paused.Toggle()
}

// Paused returns true if the cooldown is paused.
func (c *CoolDown) Paused() bool {
	if c == nil {
		return false
	}
	return c.paused.Load()
}

// Set sets the player the cooldown.
func (c *CoolDown) Set(dur time.Duration) {
	if c == nil {
		return
	}
	if c.paused.Load() {
		c.remainingAtPause.Store(dur)
		return
	}
	c.expiration = *atomic.NewValue(time.Now().Add(dur))
}

// Active returns true if the cooldown is currently active.
func (c *CoolDown) Active() bool {
	if c == nil {
		return false
	}
	return c.expiration.Load().After(time.Now())
}

func (c *CoolDown) Remaining() time.Duration {
	if c == nil {
		return 0
	}
	if c.paused.Load() {
		return c.remainingAtPause.Load()
	}
	return time.Until(c.expiration.Load())
}

// Reset resets the cooldown.
func (c *CoolDown) Reset() {
	if c == nil {
		return
	}
	c.paused.Store(false)
	c.remainingAtPause.Store(0)
	c.expiration = *atomic.NewValue(time.Time{})
}

type coolDownData struct {
	Duration         time.Duration
	Paused           bool
	RemainingAtPause time.Duration
}

// UnmarshalBSON ...
func (c *CoolDown) UnmarshalBSON(b []byte) error {
	if c == nil {
		return nil
	}
	d := coolDownData{}
	err := bson.Unmarshal(b, &d)
	c.expiration = *atomic.NewValue(time.Now().Add(d.Duration))
	c.paused.Store(d.Paused)
	c.remainingAtPause.Store(d.RemainingAtPause)
	return err
}

// MarshalBSON ...
func (c *CoolDown) MarshalBSON() ([]byte, error) {
	if c == nil {
		return bson.Marshal(&coolDownData{})
	}
	d := coolDownData{
		Duration:         time.Until(c.expiration.Load()),
		Paused:           c.paused.Load(),
		RemainingAtPause: c.remainingAtPause.Load(),
	}
	return bson.Marshal(d)
}

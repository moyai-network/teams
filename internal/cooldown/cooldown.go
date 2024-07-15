package cooldown

import (
	"github.com/df-mc/atomic"
	"github.com/rcrowley/go-bson"
	"time"
)

// CoolDown represents a time cooldown.
type CoolDown struct {
	expiration       atomic.Value[time.Time] // time.Time
	paused           atomic.Bool
	remainingAtPause atomic.Value[time.Duration] // time.Duration
}

// NewCoolDown returns a new CoolDown instance.
func NewCoolDown() *CoolDown {
	return &CoolDown{}
}

// TogglePause toggles the pause state of the cooldown.
func (c *CoolDown) TogglePause() {
	if c == nil {
		return
	}

	currentPaused := c.paused.Load()
	c.paused.Store(!currentPaused)

	if currentPaused { // If currently paused, resume
		remaining := c.remainingAtPause.Load()
		c.expiration.Store(time.Now().Add(remaining))
		c.remainingAtPause.Store(0)
	} else { // If currently active, pause
		remaining := time.Until(c.expiration.Load())
		c.remainingAtPause.Store(remaining)
		c.expiration.Store(time.Time{}) // Clear expiration on pause
	}
}

// Paused returns true if the cooldown is paused.
func (c *CoolDown) Paused() bool {
	if c == nil {
		return false
	}
	return c.paused.Load()
}

// Set sets the cooldown duration.
func (c *CoolDown) Set(dur time.Duration) {
	if c == nil {
		return
	}
	if c.paused.Load() {
		c.remainingAtPause.Store(dur)
		return
	}

	c.expiration.Store(time.Now().Add(dur))
}

// Active returns true if the cooldown is currently active.
func (c *CoolDown) Active() bool {
	if c == nil {
		return false
	}
	if c.paused.Load() {
		return c.remainingAtPause.Load() > 0
	}
	return c.expiration.Load().After(time.Now())
}

// Remaining returns the remaining cooldown duration.
func (c *CoolDown) Remaining() time.Duration {
	if c == nil {
		return 0
	}
	if c.paused.Load() {
		return c.remainingAtPause.Load()
	}
	exp := c.expiration.Load()
	return time.Until(exp)
}

// Reset resets the cooldown.
func (c *CoolDown) Reset() {
	if c == nil {
		return
	}
	c.paused.Store(false)
	c.remainingAtPause.Store(0)
	c.expiration.Store(time.Time{}) // Clear expiration
}

type coolDownData struct {
	Expiration       time.Time
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
	c.expiration = *atomic.NewValue(d.Expiration)
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
		Expiration:       c.expiration.Load(),
		Paused:           c.paused.Load(),
		RemainingAtPause: c.remainingAtPause.Load(),
	}
	return bson.Marshal(d)
}

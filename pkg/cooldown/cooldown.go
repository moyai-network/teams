package cooldown

import (
	"github.com/df-mc/atomic"
	"github.com/rcrowley/go-bson"
	"time"
)

type Func func(c *CoolDown)

// CoolDown represents a cool-down.
type CoolDown struct {
	expiration atomic.Value[time.Time]

	start Func
	end   Func

	c chan struct{}
}

// NewCoolDown returns a new cool-down.
func NewCoolDown(start, end Func) *CoolDown {
	return &CoolDown{
		start: start,
		end:   end,

		c: make(chan struct{}),
	}
}

// Active returns true if the CoolDown is active.
func (c *CoolDown) Active() bool {
	return c.expiration.Load().After(time.Now())
}

// Remaining returns the remaining time of the cool-down.
func (c *CoolDown) Remaining() time.Duration {
	return time.Until(c.expiration.Load())
}

// Set adds a duration to the Tag.
func (c *CoolDown) Set(d time.Duration) {
	if c.start != nil {
		c.start(c)
	}
	c.Cancel()

	go func() {
		select {
		case <-time.After(d):
			if c.end != nil {
				c.end(c)
			}
		case <-c.c:
			return
		}
	}()
	c.expiration.Store(time.Now().Add(d))
}

// Reset resets the Tag.
func (c *CoolDown) Reset() {
	if c.end != nil {
		c.end(c)
	}

	c.Cancel()
	c.expiration.Store(time.Time{})
}

// C returns the channel of the Tag.
func (c *CoolDown) C() <-chan struct{} {
	return c.c
}

// Cancel cancels the Tag.
func (c *CoolDown) Cancel() {
	if !c.Active() {
		return
	}
	c.expiration.Store(time.Time{})
	c.c <- struct{}{}
}

type coolDownData struct {
	Duration time.Duration
}

// UnmarshalBSON ...
func (c *CoolDown) UnmarshalBSON(b []byte) error {
	d := coolDownData{}
	err := bson.Unmarshal(b, &d)
	c.expiration = *atomic.NewValue(time.Now().Add(d.Duration))
	return err
}

// MarshalBSON ...
func (c *CoolDown) MarshalBSON() ([]byte, error) {
	d := coolDownData{Duration: time.Until(c.expiration.Load())}
	return bson.Marshal(d)
}

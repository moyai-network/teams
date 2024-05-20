package cooldown

import (
	"github.com/df-mc/atomic"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/rcrowley/go-bson"
	"time"
)

// CoolDown represents a time cooldown.
type CoolDown struct {
	expiration atomic.Value[time.Time]
	pos        mgl64.Vec3
}

// NewCoolDown returns a new process.
func NewCoolDown() *CoolDown {
	return &CoolDown{}
}

// Set sets the player the cooldown.
func (c *CoolDown) Set(dur time.Duration) {
	c.expiration = *atomic.NewValue(time.Now().Add(dur))
}

// Active returns true if the cooldown is currently active.
func (c *CoolDown) Active() bool {
	return c.expiration.Load().After(time.Now())
}

func (c *CoolDown) Remaining() time.Duration {
	return time.Until(c.expiration.Load())
}

// Reset resets the cooldown.
func (c *CoolDown) Reset() {
	c.expiration = *atomic.NewValue(time.Time{})
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

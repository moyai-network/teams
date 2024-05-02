package process

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

// Func is a function called when a process is performed.
type Func func(t *Process)

// Process represents a time Process.
type Process struct {
	expiration time.Time
	pos        mgl64.Vec3
	ongoing    bool

	f Func
	c chan struct{}
}

// NewProcess returns a new process.
func NewProcess(f Func) *Process {
	return &Process{
		f: f,
		c: make(chan struct{}),
	}
}

// Teleport teleports the player to the position after the duration has passed.
func (pr *Process) Teleport(p *player.Player, dur time.Duration, pos mgl64.Vec3) {
	pr.expiration = time.Now().Add(dur)
	pr.c = make(chan struct{})
	pr.pos = pos
	pr.ongoing = true

	go func() {
		select {
		case <-time.After(dur):
			if pr.f != nil {
				pr.f(pr)
			}
			p.Teleport(pos)
			pr.ongoing = false
		case <-pr.c:
			pr.ongoing = false
			return
		}
	}()
}

// Teleporting returns true if the player is currently ongoing.
func (pr *Process) Teleporting() bool {
	return pr.ongoing
}

// Expired returns true if the teleportation has expired.
func (pr *Process) Expired() bool {
	return time.Now().After(pr.expiration)
}

// Expiration returns the expiration time of the teleportation.
func (pr *Process) Expiration() time.Time {
	return pr.expiration
}

// Pos returns the position the player will be teleported to.
func (pr *Process) Pos() mgl64.Vec3 {
	return pr.pos
}

// C returns the channel that is closed when the teleportation is cancelled.
func (pr *Process) C() <-chan struct{} {
	return pr.c
}

// Cancel cancels the teleportation.
func (pr *Process) Cancel() {
	if pr.ongoing {
		pr.c <- struct{}{}
	}
}

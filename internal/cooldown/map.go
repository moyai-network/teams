package cooldown

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// MappedCoolDown represents a cool-down mapped to a key.
type MappedCoolDown[T comparable] map[T]*CoolDown

// NewMappedCoolDown returns a new mapped cool-down.
func NewMappedCoolDown[T comparable]() MappedCoolDown[T] {
	return make(map[T]*CoolDown)
}

// Active returns true if the cool-down is active.
func (m MappedCoolDown[T]) Active(key T) bool {
	if m == nil {
		m = make(MappedCoolDown[T])
	}
	coolDown, ok := m[key]
	return ok && coolDown.Active()
}

// Set sets the cool-down.
func (m MappedCoolDown[T]) Set(key T, d time.Duration) {
	if m == nil {
		m = make(MappedCoolDown[T])
	}
	coolDown := m.Key(key)
	coolDown.Set(d)
	m[key] = coolDown
}

// Key returns the cool-down for the key.
func (m MappedCoolDown[T]) Key(key T) *CoolDown {
	if m == nil {
		m = make(MappedCoolDown[T])
	}
	coolDown, ok := m[key]
	if !ok {
		newCD := NewCoolDown()
		m[key] = newCD
		return newCD
	}
	return coolDown
}

// Reset resets the cool-down.
func (m MappedCoolDown[T]) Reset(key T) {
	if m == nil {
		m = make(MappedCoolDown[T])
	}
	delete(m, key)
}

// Remaining returns the remaining time of the cool-down.
func (m MappedCoolDown[T]) Remaining(key T) time.Duration {
	if m == nil {
		m = make(MappedCoolDown[T])
	}
	coolDown, ok := m[key]
	if !ok {
		return 0
	}
	return coolDown.Remaining()
}

// All returns all cool-downs.
func (m MappedCoolDown[T]) All() (coolDowns []*CoolDown) {
	if m == nil {
		m = make(MappedCoolDown[T])
	}
	for _, coolDown := range m {
		coolDowns = append(coolDowns, coolDown)
	}
	return coolDowns
}

// MarshalBSON ...
func (m MappedCoolDown[T]) MarshalBSON() ([]byte, error) {
	if m == nil {
		m = make(MappedCoolDown[T])
	}
	d := map[T]time.Time{}
	for k, cd := range m {
		d[k] = cd.expiration.Load()
	}
	return bson.Marshal(d)
}

// UnmarshalBSON ...
func (m MappedCoolDown[T]) UnmarshalBSON(b []byte) error {
	if m == nil {
		m = make(MappedCoolDown[T])
	}
	d := map[T]time.Time{}
	err := bson.Unmarshal(b, &d)
	if err != nil {
		return err
	}

	for k, cd := range d {
		m.Set(k, time.Until(cd))
	}
	return nil
}

package entity

/*
// NewLightning creates a new lightning entity.
func NewLightning(pos mgl64.Vec3) *entity.Ent {
	state := &lightningState{
		state:    2,
		lifetime: rand.Intn(4) + 1,
	}
	conf := lightningConf
	conf.Tick = state.tick
	return entity.Config{Behaviour: conf.New()}.New(entity.LightningType{}, pos)
}

var lightningConf = entity.StationaryBehaviourConfig{SpawnSounds: []world.Sound{sound.Explosion{}, sound.Thunder{}}}

// lightningState holds the state of a lightning entity.
type lightningState struct {
	state, lifetime int
}

// tick carries out lightning logic, dealing damage and setting blocks/entities
// on fire when appropriate.
func (s *lightningState) tick(e *entity.Ent) {
	if s.state--; s.state < 0 {
		if s.lifetime == 0 {
			_ = e.Close()
		} else if s.state < -rand.Intn(10) {
			s.lifetime--
			s.state = 1
		}
	}
}
*/

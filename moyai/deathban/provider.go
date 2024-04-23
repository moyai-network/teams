package deathban

import (
	kit2 "github.com/moyai-network/teams/internal/kit"
	"sync"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/df-mc/goleveldb/leveldb/opt"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose/sets"
)

func init() {
	p, err := mcdb.Config{
		Compression: opt.DefaultCompression,
		ReadOnly:    true,
	}.Open("assets/kek")
	if err != nil {
		return
	}

	w := world.Config{
		ReadOnly: true,
		Provider: p,
		Entities: entity.DefaultRegistry,
	}.New()

	w.SetSpawn(cube.Pos{0, 40})

	deathban = new(w)
}

// deathbans maps between a *player.Player and a *Provider.
var deathbans sync.Map

// LookupProvider looks up the *Provider of the *player.Player passed.
func Lookup(p *player.Player) (*Provider, bool) {
	h, ok := deathbans.Load(p)
	if ok {
		return h.(*Provider), ok
	}
	return nil, false
}

// lobby ...
var deathban *Provider

// Lobby ...
func Deathban() *Provider {
	return deathban
}

// Provider is the provider for the deathban.
type Provider struct {
	w *world.World

	playerMu sync.Mutex
	players  sets.Set[*player.Player]
}

// new creates a new lobby provider.
func new(w *world.World) *Provider {
	deathban = &Provider{
		w:       w,
		players: make(sets.Set[*player.Player]),
	}
	return deathban
}

// AddPlayer ...
func (s *Provider) AddPlayer(p *player.Player) {
	deathbans.Store(p, s)

	kit2.Apply(kit2.Diamond{}, p)

	s.playerMu.Lock()
	s.players.Add(p)
	s.playerMu.Unlock()

	if p.World() != s.w {
		s.w.AddEntity(p)
	}
	p.Teleport(mgl64.Vec3{0, 40})

	p.Move(mgl64.Vec3{}, 180, -360)
}

// RemovePlayer ...
func (s *Provider) RemovePlayer(p *player.Player, force bool) {
	deathbans.Delete(p)
	p.Inventory().Handle(nil)

	s.playerMu.Lock()
	s.players.Delete(p)
	s.playerMu.Unlock()
}

// Players ...
func (s *Provider) Players() []*player.Player {
	s.playerMu.Lock()
	defer s.playerMu.Unlock()
	return s.players.Values()
}

// PlayerCount ...
func (s *Provider) PlayerCount() int {
	s.playerMu.Lock()
	defer s.playerMu.Unlock()
	return len(s.players)
}

// World ...
func (s *Provider) World() *world.World {
	return s.w
}

package user

import (
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode"
	_ "unsafe"

	"github.com/moyai-network/teams/internal/cooldown"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/process"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/moyai/sotw"

	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
)

// LookupRuntimeID ...
func LookupRuntimeID(p *player.Player, rid uint64) (*player.Player, bool) {
	h, ok := p.Handler().(*Handler)
	if !ok {
		return nil, false
	}
	for _, t := range moyai.Server().Players() {
		if session_entityRuntimeID(h.s, t) == rid {
			return t, true
		}
	}
	return nil, false
}

// Lookup looks up the Handler of a name passed.
func Lookup(name string) (*player.Player, bool) {
	for _, t := range moyai.Server().Players() {
		if strings.EqualFold(name, t.Name()) {
			return t, true
		}
	}
	return nil, false
}

// Broadcast broadcasts a message to every user using that user's locale.
func Broadcast(key string, args ...any) {
	for _, p := range moyai.Server().Players() {
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			continue
		}
		p.Message(lang.Translatef(u.Language, key, args...))
	}
}
func (h *Handler) Player() *player.Player {
	return h.p
}

func (h *Handler) Logout() *process.Process {
	return h.logout
}

func (h *Handler) Stuck() *process.Process {
	return h.stuck
}

func (h *Handler) Home() *process.Process {
	return h.home
}

// SubtractItem subtracts d from the count of the item stack passed and returns it, if the player is in
// survival or adventure mode.
func (h *Handler) SubtractItem(s item.Stack, d int) item.Stack {
	if !h.p.GameMode().CreativeInventory() && d != 0 {
		return s.Grow(-d)
	}
	return s
}

func (h *Handler) Combat() *cooldown.CoolDown {
	return h.combat
}

func (h *Handler) Pearl() *cooldown.CoolDown {
	return h.pearl
}

func (h *Handler) FactionCreate() *cooldown.CoolDown {
	return h.factionCreate
}

// Vanished returns whether the user is vanished or not.
func (h *Handler) Vanished() bool {
	return h.vanished.Load()
}

// ToggleVanish toggles the user's vanish state.
func (h *Handler) ToggleVanish() {
	h.vanished.Toggle()
}

func (h *Handler) SetLastPearlPos(pos mgl64.Vec3) {
	h.lastPearlPos = pos
}

func (h *Handler) SetLastHitBy(p *player.Player) {
	h.lastHitBy = p
}

func (h *Handler) LastPearlPos() mgl64.Vec3 {
	return h.lastPearlPos
}

func (h *Handler) LastHitBy() *player.Player {
	return h.lastHitBy
}

func (h *Handler) DropItem(it item.Stack) {
	p := h.p
	w, pos := p.World(), p.Position()
	et := entity.NewItem(it, pos)
	et.SetVelocity(mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1})
	w.AddEntity(et)
}

// Boned returns whether the user has been boned.
func (h *Handler) Boned() bool {
	return h.bone.Active()
}

// BoneHits returns the number of bone hits of the user.
func (h *Handler) BoneHits(p *player.Player) int {
	hits, ok := h.boneHits[p.Name()]
	if !ok {
		return 0
	}
	return hits
}

// AddBoneHit adds a bone hit to the user.
func (h *Handler) AddBoneHit(p *player.Player) {
	h.boneHits[p.Name()]++
	if h.boneHits[p.Name()] >= 3 {
		h.ResetBoneHits(p)
		h.bone.Set(15 * time.Second)
	}
}

// ResetBoneHits resets the bone hits of the user.
func (h *Handler) ResetBoneHits(p *player.Player) {
	h.boneHits[p.Name()] = 0
}

// ScramblerHits returns the number of scrambler hits of the user.
func (h *Handler) ScramblerHits(p *player.Player) int {
	hits, ok := h.scramblerHits[p.Name()]
	if !ok {
		return 0
	}
	return hits
}

// AddScramblerHit adds a scrambler hit to the user.
func (h *Handler) AddScramblerHit(p *player.Player) {
	h.scramblerHits[p.Name()] = h.scramblerHits[p.Name()] + 1
}

// ResetScramblerHits resets the scrambler hits of the user.
func (h *Handler) ResetScramblerHits(p *player.Player) {
	h.scramblerHits[p.Name()] = 0
}

// PearlDisabled returns whether the user is pearl disabled.
func (h *Handler) PearlDisabled() bool {
	return h.pearlDisabled
}

// TogglePearlDisable toggles the pearl disabler
func (h *Handler) TogglePearlDisable() {
	h.pearlDisabled = !h.pearlDisabled
}

// CanSendMessage returns true if the user can send a message.
func (h *Handler) CanSendMessage() bool {
	return time.Since(h.lastMessage.Load()) > time.Second*1
}

// LastAttacker returns the last attacker of the user.
func (h *Handler) LastAttacker() (*player.Player, bool) {
	if time.Since(h.lastAttackTime.Load()) > 15*time.Second {
		return nil, false
	}
	name := h.lastAttackerName.Load()
	if len(name) == 0 {
		return nil, false
	}
	return Lookup(name)
}

// SetLastAttacker sets the last attacker of the user.
func (h *Handler) SetLastAttacker(t *Handler) {
	h.lastAttackerName.Store(t.p.Name())
	h.lastAttackTime.Store(time.Now())
}

// ResetLastAttacker resets the last attacker of the user.
func (h *Handler) ResetLastAttacker() {
	h.lastAttackerName.Store("")
	h.lastAttackTime.Store(time.Time{})
}

// ShowArmor displays or removes players armor visibility from other players.
func (h *Handler) ShowArmor(visible bool) {
	p := h.Player()

	air := item.NewStack(block.Air{}, 1)

	helmet := item.NewStack(block.Air{}, 1)
	if !p.Armour().Helmet().Equal(air) && visible {
		helmet = p.Armour().Helmet()
	}

	chestplate := item.NewStack(block.Air{}, 1)
	if !p.Armour().Chestplate().Equal(air) && visible {
		chestplate = p.Armour().Chestplate()
	}

	leggings := item.NewStack(block.Air{}, 1)
	if !p.Armour().Leggings().Equal(air) && visible {
		leggings = p.Armour().Leggings()
	}

	boots := item.NewStack(block.Air{}, 1)
	if !p.Armour().Boots().Equal(air) && visible {
		boots = p.Armour().Boots()
	}

	for _, pl := range moyai.Server().Players() {
		if t, err := data.LoadTeamFromName(p.Name()); err == nil {
			// maybe add an option eventually, so we can use this for staff mode and other stuff IDK?
			if !t.Member(pl.Name()) {
				s := player_session(pl)
				session_writePacket(s, &packet.MobArmourEquipment{
					EntityRuntimeID: session_entityRuntimeID(s, p),
					Helmet:          instanceFromItem(s, helmet),
					Chestplate:      instanceFromItem(s, chestplate),
					Leggings:        instanceFromItem(s, leggings),
					Boots:           instanceFromItem(s, boots),
				})
			}
		}
	}
}

func (h *Handler) sendWall(newPos cube.Pos, z area.Area, color item.Colour) {
	areaMin := cube.Pos{int(z.Min().X()), 0, int(z.Min().Y())}
	areaMax := cube.Pos{int(z.Max().X()), 255, int(z.Max().Y())}
	wallBlock := block.StainedGlass{Colour: color}
	const wallLength, wallHeight = 15, 10

	if newPos.X() >= areaMin.X() && newPos.X() <= areaMax.X() { // edges of the top and bottom walls (relative to South)
		zCoord := areaMin.Z()
		if newPos.Z() >= areaMax.Z() {
			zCoord = areaMax.Z()
		}
		for horizontal := newPos.X() - wallLength; horizontal < newPos.X()+wallLength; horizontal++ {
			if horizontal >= areaMin.X() && horizontal <= areaMax.X() {
				for vertical := newPos.Y(); vertical < newPos.Y()+wallHeight; vertical++ {
					blockPos := cube.Pos{horizontal, vertical, zCoord}
					if blockReplaceable(h.p.World().Block(blockPos)) {
						h.s.ViewBlockUpdate(blockPos, wallBlock, 0)
						h.wallBlocksMu.Lock()
						h.wallBlocks[blockPos] = rand.Float64() * float64(rand.Intn(1)+1)
						h.wallBlocksMu.Unlock()
					}
				}
			}
		}
	}
	if newPos.Z() >= areaMin.Z() && newPos.Z() <= areaMax.Z() { // edges of the left and right walls (relative to South)
		xCoord := areaMin.X()
		if newPos.X() >= areaMax.X() {
			xCoord = areaMax.X()
		}
		for horizontal := newPos.Z() - wallLength; horizontal < newPos.Z()+wallLength; horizontal++ {
			if horizontal >= areaMin.Z() && horizontal <= areaMax.Z() {
				for vertical := newPos.Y(); vertical < newPos.Y()+wallHeight; vertical++ {
					blockPos := cube.Pos{xCoord, vertical, horizontal}
					if blockReplaceable(h.p.World().Block(blockPos)) {
						h.s.ViewBlockUpdate(blockPos, wallBlock, 0)
						h.wallBlocksMu.Lock()
						h.wallBlocks[blockPos] = rand.Float64() * float64(rand.Intn(1)+1)
						h.wallBlocksMu.Unlock()
					}
				}
			}
		}
	}
}

func formatItemName(s string) string {
	split := strings.Split(s, "_")
	for i, str := range split {
		upperCasesPrefix := unicode.ToUpper(rune(str[0]))
		split[i] = string(upperCasesPrefix) + str[1:]
	}
	return strings.Join(split, " ")
}

// clearWall clears the users walls or lowers the remaining duration.
func (h *Handler) clearWall() {
	h.wallBlocksMu.Lock()
	for p, duration := range h.wallBlocks {
		if duration <= 0 {
			delete(h.wallBlocks, p)
			h.s.ViewBlockUpdate(p, block.Air{}, 0)
			h.p.ShowParticle(p.Vec3(), particle.BlockForceField{})
			continue
		}
		h.wallBlocks[p] = duration - 0.1
	}
	h.wallBlocksMu.Unlock()
}

// viewers returns a list of all viewers of the Player.
func (h *Handler) viewers() []world.Viewer {
	viewers := h.p.World().Viewers(h.p.Position())
	for _, v := range viewers {
		if v == h.s {
			return viewers
		}
	}
	return append(viewers, h.s)
}

// blockReplaceable checks if the combat wall should replace a block.
func blockReplaceable(b world.Block) bool {
	_, air := b.(block.Air)
	_, doubleFlower := b.(block.DoubleFlower)
	_, flower := b.(block.Flower)
	_, tallGrass := b.(block.TallGrass)
	_, doubleTallGrass := b.(block.DoubleTallGrass)
	_, deadBush := b.(block.DeadBush)
	//_, cobweb := b.(block.Cobweb)
	//_, sapling := b.(block.Sapling)
	_, torch := b.(block.Torch)
	_, fire := b.(block.Fire)
	return air || tallGrass || deadBush || torch || fire || flower || doubleFlower || doubleTallGrass
}

// DropContents drops the contents of the user.
func DropContents(p *player.Player) {
	drop_contents(p)
}

// addEffects adds a list of effects to the user.
func addEffects(p *player.Player, effects ...effect.Effect) {
	for _, e := range effects {
		p.AddEffect(e)
	}
}

// removeEffects removes a list of effects from the user.
func removeEffects(p *player.Player, effects ...effect.Effect) {
	for _, e := range effects {
		p.RemoveEffect(e.Type())
	}
}

// hasEffectLevel returns whether the user has the effect or not.
func hasEffectLevel(p *player.Player, e effect.Effect) bool {
	for _, ef := range p.Effects() {
		if e.Type() == ef.Type() && e.Level() == ef.Level() {
			return true
		}
	}
	return false
}

// canAttack returns true if the given players can attack each other.
func canAttack(pl, target *player.Player) bool {
	if target == nil || pl == nil {
		return false
	}
	w := pl.World()
	if area.Spawn(w).Vec3WithinOrEqualFloorXZ(pl.Position()) || area.Spawn(w).Vec3WithinOrEqualFloorXZ(target.Position()) {
		return false
	}

	u, _ := data.LoadUserFromName(pl.Name())
	t, _ := data.LoadUserFromName(target.Name())

	_, sotwRunning := sotw.Running()
	if (u.Teams.PVP.Active() || t.Teams.PVP.Active()) || (sotwRunning && (u.Teams.SOTW || t.Teams.SOTW)) {
		return false
	}

	tm, err := data.LoadTeamFromMemberName(pl.Name())
	if err != nil {
		return true
	}

	return !tm.Member(target.Name())
}

// Nearby returns the nearby users of a certain distance from the user
func Nearby(p *player.Player, dist float64) []*Handler {
	var pl []*Handler
	for _, e := range p.World().Entities() {
		if e.Position().ApproxFuncEqual(p.Position(), func(f float64, f2 float64) bool {
			return math.Max(f, f2)-math.Min(f, f2) < dist
		}) {
			if target, ok := e.(*player.Player); ok && target != p {
				if h, ok := target.Handler().(*Handler); ok {
					pl = append(pl, h)
				}
			}
		}
	}
	return pl
}

// NearbyEnemies returns the nearby enemies of a certain distance from the user
func NearbyEnemies(p *player.Player, dist float64) []*Handler {
	var pl []*Handler
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		return Nearby(p, dist)
	}

	for _, target := range Nearby(p, dist) {
		if !tm.Member(target.p.Name()) {
			pl = append(pl, target)
		}
	}

	return pl
}

// NearbyAllies returns the nearby allies of a certain distance from the user
func NearbyAllies(p *player.Player, dist float64) []*Handler {
	h, ok := p.Handler().(*Handler)
	if !ok {
		return []*Handler{}

	}
	pl := []*Handler{h}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		return pl
	}

	for _, target := range Nearby(p, dist) {
		if tm.Member(target.p.Name()) {
			pl = append(pl, target)
		}
	}

	return pl
}

// NearbyCombat returns the nearby faction members, enemies, no faction players (basically anyone not on timer) of a certain distance from the user
// Nearby returns the nearby users of a certain distance from the user
func NearbyCombat(p *player.Player, dist float64) []*Handler {
	var pl []*Handler

	for _, target := range Nearby(p, dist) {
		t, _ := data.LoadUserFromName(target.p.Name())
		if !t.Teams.PVP.Active() {
			pl = append(pl, target)
		}
	}

	return pl
}

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session

// noinspection ALL
//
//go:linkname drop_contents github.com/df-mc/dragonfly/server/player.(*Player).dropContents
func drop_contents(*player.Player)

// noinspection ALL
//
//go:linkname instanceFromItem github.com/df-mc/dragonfly/server/session.instanceFromItem
func instanceFromItem(*session.Session, item.Stack) protocol.ItemInstance

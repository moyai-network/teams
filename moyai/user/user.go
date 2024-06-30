package user

import (
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode"
	_ "unsafe"

	"github.com/moyai-network/teams/moyai/conquest"
	"github.com/moyai-network/teams/moyai/koth"

	"github.com/df-mc/dragonfly/server/item/inventory"

	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/unsafe"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/class"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// lookupRuntimeID ...
func lookupRuntimeID(p *player.Player, rid uint64) (*player.Player, bool) {
	h, ok := p.Handler().(*Handler)
	if !ok {
		return nil, false
	}
	for _, t := range moyai.Players() {
		if session_entityRuntimeID(unsafe.Session(h.p), t) == rid {
			return t, true
		}
	}
	return nil, false
}

// Lookup looks up the Handler of a name passed.
func Lookup(name string) (*player.Player, bool) {
	for _, t := range moyai.Players() {
		if strings.EqualFold(name, t.Name()) {
			return t, true
		}
	}
	return nil, false
}

func UpdateState(p *player.Player) {
	for _, v := range p.World().Viewers(p.Position()) {
		v.ViewEntityState(p)
	}
}

func hideVanished(p *player.Player) {
	for _, t := range moyai.Players() {
		u, err := data.LoadUserFromName(t.Name())
		if err != nil {
			continue
		}
		if u.Vanished {
			p.HideEntity(t)
		}
	}
}

func showVanished(p *player.Player) {
	for _, t := range moyai.Players() {
		u, err := data.LoadUserFromName(t.Name())
		if err != nil {
			continue
		}
		if u.Vanished {
			p.ShowEntity(t)
		}
	}
}

// UpdateVanishState vanishes the user.
func UpdateVanishState(p *player.Player, u data.User) {
	if u.Vanished {
		showVanished(p)
	} else {
		hideVanished(p)
	}

	for _, t := range moyai.Players() {
		target, err := data.LoadUserFromName(p.Name())
		if err != nil {
			continue
		}
		if !target.Vanished && u.Vanished {
			t.HideEntity(p)
		} else if !target.Vanished && !u.Vanished {
			t.ShowEntity(p)
		}
	}
}

// LogTime returns the log time of the user.
func (h *Handler) LogTime() time.Duration {
	return time.Since(h.logTime)
}

// clearOwnedEntities clears all entities owned by the user.
func (h *Handler) clearOwnedEntities() {
	for _, et := range h.p.World().Entities() {
		if be, ok := et.(entity.Behaviour); ok {
			if pro, ok := be.(*entity.ProjectileBehaviour); ok {
				if pro.Owner() == h.p {
					_ = et.Close()
					h.p.World().RemoveEntity(et)
				}
			}
		}
	}
}

// clearEffects clears all effects from the user.
func (h *Handler) clearEffects() {
	for _, ef := range h.p.Effects() {
		h.p.RemoveEffect(ef.Type())
	}
}

// resetCoolDowns resets all cooldowns of the user.
func (h *Handler) resetCoolDowns() {
	h.tagCombat.Reset()
	h.coolDownPearl.Reset()
	h.coolDownGoldenApple.Reset()
}

// spawnDeathNPC spawns a death NPC at the user's position.
func (h *Handler) spawnDeathNPC(src world.DamageSource) {
	p := h.p

	npc := player.New(p.Name(), p.Skin(), p.Position())
	npc.Handle(npcHandler{})
	npc.SetAttackImmunity(time.Millisecond * 1400)
	npc.SetNameTag(p.NameTag())
	npc.SetScale(p.Scale())
	p.World().AddEntity(npc)

	for _, viewer := range p.World().Viewers(npc.Position()) {
		viewer.ViewEntityAction(npc, entity.DeathAction{})
	}
	time.AfterFunc(time.Second*2, func() {
		_ = npc.Close()
	})

	if att, ok := attackerFromSource(src); ok {
		npc.KnockBack(att.Position(), 0.5, 0.2)
	}
}

// kill handles the death of the user.
func (h *Handler) kill(src world.DamageSource) {
	p := h.p
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	p.World().PlaySound(p.Position(), sound.Explosion{})
	if !u.Teams.DeathBan.After(time.Now()) {
		h.handleTeamMemberDeath()
	}

	if k, ok := koth.Running(); ok {
		k.StopCapturing(p)
	}

	if conquest.Running() {
		for _, c := range conquest.All() {
			c.StopCapturing(p)
		}
	}

	h.stopCapturing()
	h.incrementDeath()
	h.issueDeathban()
	h.cancelStormBreak()
	h.spawnDeathNPC(src)
	h.clearEffects()
	h.clearOwnedEntities()
	h.resetCoolDowns()
	unsafe.Session(h.p).EmptyUIInventory()

	DropContents(p)
	p.SetHeldItems(item.Stack{}, item.Stack{})
	p.ResetFallDistance()
	p.Heal(20, effect.InstantHealingSource{})
	p.Extinguish()
	p.SetFood(20)

	h.lastClass.Store(class.Resolve(p))
	UpdateState(p)

	sortArmourEffects(h)
	sortClassEffects(h)
	moyai.Deathban().AddEntity(p)
	unsafe.WritePacket(h.p, &packet.PlayerFog{
		Stack: []string{"minecraft:fog_default"},
	})
	p.Teleport(mgl64.Vec3{5, 13, 44})
}

// stopCapturing stops the user from capturing a koth or a point.
func (h *Handler) stopCapturing() {
	k, ok := koth.Running()
	if ok {
		k.StopCapturing(h.p)
	}
	ok = conquest.Running()
	if ok {
		for _, c := range conquest.All() {
			c.StopCapturing(h.p)
		}
	}
}

// cancelStormBreak cancels the storm breaker effect.
func (h *Handler) cancelStormBreak() {
	p := h.p
	armourHandler, ok := p.Armour().Inventory().Handler().(*ArmourHandler)
	if ok {
		armourHandler.stormBreakCancel()
	}
}

// incrementDeath increments the death count of the user.
func (h *Handler) incrementDeath() {
	victim, err := data.LoadUserFromName(h.p.Name())
	if err != nil {
		return
	}
	victim.Teams.Stats.Deaths += 1
	if victim.Teams.Stats.KillStreak > victim.Teams.Stats.BestKillStreak {
		victim.Teams.Stats.BestKillStreak = victim.Teams.Stats.KillStreak
	}
	victim.Teams.Stats.KillStreak = 0

	held, off := h.p.HeldItems()
	*victim.Teams.DeathInventory = inventoryData(held, off, h.p.Armour(), h.p.Inventory())
	data.SaveUser(victim)
}

func inventoryData(held, off item.Stack, a *inventory.Armour, i *inventory.Inventory) data.Inventory {
	return data.Inventory{
		MainHandSlot: 0,
		OffHand:      off,
		Items:        i.Slots(),
		Boots:        a.Boots(),
		Leggings:     a.Leggings(),
		Chestplate:   a.Chestplate(),
		Helmet:       a.Helmet(),
	}
}

// issueDeathban issues a deathban for the user.
func (h *Handler) issueDeathban() {
	u, err := data.LoadUserFromName(h.p.Name())
	if err != nil {
		return
	}
	u.Teams.DeathBan = time.Now().Add(time.Minute * 20)
	u.Teams.DeathBanned = true

	data.SaveUser(u)
}

// handleTeamMemberDeath handles the death of a team member.
func (h *Handler) handleTeamMemberDeath() {
	if tm, err := data.LoadTeamFromMemberName(h.p.Name()); err == nil {
		tm = tm.WithDTR(tm.DTR - 1).WithLastDeath(time.Now())
		if tm.Points > 0 {
			tm = tm.WithPoints(tm.Points - 1)
		}
		data.SaveTeam(tm)
	}
}

// canMessage returns true if the user can send a message.
func (h *Handler) canMessage() bool {
	return time.Since(h.lastMessage.Load()) > time.Second*1
}

// lastAttacker returns the last attacker of the user.
func (h *Handler) lastAttacker() (*player.Player, bool) {
	if time.Since(h.lastAttackTime.Load()) > 15*time.Second {
		return nil, false
	}
	name := h.lastAttackerName.Load()
	if len(name) == 0 {
		return nil, false
	}
	return Lookup(name)
}

// setLastAttacker sets the last attacker of the user.
func (h *Handler) setLastAttacker(t *Handler) {
	h.lastAttackerName.Store(t.p.Name())
	h.lastAttackTime.Store(time.Now())
}

// resetLastAttacker resets the last attacker of the user.
func (h *Handler) resetLastAttacker() {
	h.lastAttackerName.Store("")
	h.lastAttackTime.Store(time.Time{})
}

// ShowArmor displays or removes players armor visibility from other players.
func (h *Handler) ShowArmor(visible bool) {
	p := h.p

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

	tm, tmErr := data.LoadTeamFromMemberName(p.Name())
	for _, pl := range moyai.Players() {
		s := unsafe.Session(pl)
		if tmErr == nil {
			if !tm.Member(pl.Name()) {
				h.updateArmour(s, helmet, chestplate, leggings, boots)
			}
			continue
		}
		h.updateArmour(s, helmet, chestplate, leggings, boots)
	}
}

func (h *Handler) updateArmour(s *session.Session, helmet item.Stack, chestplate item.Stack, leggings item.Stack, boots item.Stack) {
	unsafe.WritePacket(s, &packet.MobArmourEquipment{
		EntityRuntimeID: session_entityRuntimeID(s, h.p),
		Helmet:          instanceFromItem(s, helmet),
		Chestplate:      instanceFromItem(s, chestplate),
		Leggings:        instanceFromItem(s, leggings),
		Boots:           instanceFromItem(s, boots),
	})
}

func (h *Handler) SendClaimPillar(pos cube.Pos) {
	for y := pos.Y(); y <= pos.Y()+25; y++ {
		delta := y - pos.Y()
		var b world.Block
		if delta == 3 || delta == 4 {
			b = block.Diamond{}
		} else {
			b = block.Glass{}
		}
		h.viewBlockUpdate(cube.Pos{pos.X(), y, pos.Z()}, b, 0)
	}
}

func (h *Handler) SendAirPillar(pos cube.Pos) {
	for y := pos.Y(); y <= pos.Y()+25; y++ {
		h.viewBlockUpdate(cube.Pos{pos.X(), y, pos.Z()}, block.Air{}, 0)
	}
}

// sendWall sends a wall to the user.
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
						h.viewBlockUpdate(blockPos, wallBlock, 0)
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
						h.viewBlockUpdate(blockPos, wallBlock, 0)
						h.wallBlocksMu.Lock()
						h.wallBlocks[blockPos] = rand.Float64() * float64(rand.Intn(1)+1)
						h.wallBlocksMu.Unlock()
					}
				}
			}
		}
	}
}

// formatItemName formats the item name to a more readable format.
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
			h.viewBlockUpdate(p, h.p.World().Block(p), 0)
			//h.p.ShowParticle(p.Vec3(), particle.BlockForceField{})
			continue
		}
		h.wallBlocks[p] = duration - 0.1
	}
	h.wallBlocksMu.Unlock()
}

// viewBlockUpdate updates the block at the position passed for the user.
func (h *Handler) viewBlockUpdate(p cube.Pos, b world.Block, layer int) {
	s := unsafe.Session(h.p)
	s.ViewBlockUpdate(p, b, layer)
}

// viewers returns a list of all viewers of the Player.
func (h *Handler) viewers() []world.Viewer {
	viewers := h.p.World().Viewers(h.p.Position())
	s := unsafe.Session(h.p)

	for _, v := range viewers {
		if v == s {
			return viewers
		}
	}
	return append(viewers, s)
}

// blockReplaceable checks if the tagCombat wall should replace a block.
func blockReplaceable(b world.Block) bool {
	_, air := b.(block.Air)
	_, doubleFlower := b.(block.DoubleFlower)
	_, flower := b.(block.Flower)
	_, tallGrass := b.(block.DoubleTallGrass)
	_, doubleTallGrass := b.(block.DoubleTallGrass)
	_, deadBush := b.(block.DeadBush)
	//_, cobweb := b.(block.Cobweb)
	//_, sapling := b.(block.Sapling)
	_, torch := b.(block.Torch)
	_, fire := b.(block.Fire)
	return air || tallGrass || deadBush || torch || fire || flower || doubleFlower || doubleTallGrass
}

// substractItem subtracts d from the count of the item stack passed and returns it, if the player is in
// survival or adventure mode.
func substractItem(p *player.Player, s item.Stack, d int) item.Stack {
	if !p.GameMode().CreativeInventory() && d != 0 {
		return s.Grow(-d)
	}
	return s
}

func setLogger(p *player.Player, l *Handler) {
	l.logger = true

	loggerMu.Lock()
	loggers[p.XUID()] = l
	loggerMu.Unlock()
}

func logger(p *player.Player) (*Handler, bool) {
	loggerMu.Lock()
	l, ok := loggers[p.XUID()]
	loggerMu.Unlock()
	return l, ok
}

// PlayTime returns the play time of the user.
func PlayTime(p *player.Player) time.Duration {
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return 0
	}
	h, ok := p.Handler().(*Handler)
	if ok {
		u.PlayTime += time.Since(h.logTime)
	}
	return u.PlayTime
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

// nearbyPlayers returns the nearby users of a certain distance from the user
func nearbyPlayers(p *player.Player, dist float64) []*Handler {
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

// nearbyEnemies returns the nearby enemies of a certain distance from the user
func nearbyEnemies(p *player.Player, dist float64) []*Handler {
	var pl []*Handler
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		return nearbyPlayers(p, dist)
	}

	for _, target := range nearbyPlayers(p, dist) {
		if !tm.Member(target.p.Name()) {
			pl = append(pl, target)
		}
	}

	return pl
}

// nearbyAllies returns the nearby allies of a certain distance from the user
func nearbyAllies(p *player.Player, dist float64) []*Handler {
	h, ok := p.Handler().(*Handler)
	if !ok {
		return []*Handler{}

	}
	pl := []*Handler{h}
	tm, err := data.LoadTeamFromMemberName(p.Name())
	if err != nil {
		return pl
	}

	for _, target := range nearbyPlayers(p, dist) {
		if tm.Member(target.p.Name()) {
			pl = append(pl, target)
		}
	}

	return pl
}

// nearbyHurtable returns the nearby faction members, enemies, no faction players (basically anyone not on timer) of a certain distance from the user
// nearbyPlayers returns the nearby users of a certain distance from the user
func nearbyHurtable(p *player.Player, dist float64) []*Handler {
	var pl []*Handler

	for _, target := range nearbyPlayers(p, dist) {
		t, _ := data.LoadUserFromName(target.p.Name())
		if !t.Teams.PVP.Active() {
			pl = append(pl, target)
		}
	}

	return pl
}

// noinspection ALL
//
//go:linkname DropContents github.com/df-mc/dragonfly/server/player.(*Player).dropContents
func DropContents(*player.Player)

// noinspection ALL
//
//go:linkname instanceFromItem github.com/df-mc/dragonfly/server/session.instanceFromItem
func instanceFromItem(*session.Session, item.Stack) protocol.ItemInstance

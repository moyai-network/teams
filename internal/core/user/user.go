package user

import (
	"fmt"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/area"
	"github.com/moyai-network/teams/internal/core/conquest"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/eotw"
	"github.com/moyai-network/teams/internal/core/koth"
	"github.com/moyai-network/teams/internal/core/user/class"
	model3 "github.com/moyai-network/teams/internal/model"
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode"
	_ "unsafe"

	"github.com/bedrock-gophers/unsafe/unsafe"

	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	model2 "github.com/df-mc/dragonfly/server/block/model"
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
	for t := range internal.Players(nil) {
		if session_entityRuntimeID(unsafe.Session(p), t) == rid {
			return t, true
		}
	}
	return nil, false
}

// Lookup looks up the Handler of a name passed.
func Lookup(tx *world.Tx, name string) (*player.Player, bool) {
	for t := range internal.Players(tx) {
		if strings.EqualFold(name, t.Name()) {
			return t, true
		}
	}
	return nil, false
}

func UpdateState(p *player.Player) {
	for _, v := range p.Tx().Viewers(p.Position()) {
		v.ViewEntityState(p)
	}
}

func hideVanished(p *player.Player) {
	for t := range internal.Players(p.Tx()) {
		u, err := data2.LoadUserFromName(t.Name())
		if err != nil {
			continue
		}
		if u.Vanished {
			p.HideEntity(t)
		}
	}
}

func showVanished(p *player.Player) {
	for t := range internal.Players(p.Tx()) {
		u, err := data2.LoadUserFromName(t.Name())
		if err != nil {
			continue
		}
		if u.Vanished {
			p.ShowEntity(t)
		}
	}
}

// UpdateVanishState vanishes the user.
func UpdateVanishState(p *player.Player, u model3.User) {
	if u.Vanished {
		showVanished(p)
	} else {
		hideVanished(p)
	}

	for t := range internal.Players(p.Tx()) {
		target, err := data2.LoadUserFromName(p.Name())
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
func (h *Handler) clearOwnedEntities(p *player.Player) {
	for et := range p.Tx().Entities() {
		if ent, ok := et.(*entity.Ent); ok {
			if pro, ok := ent.Behaviour().(*entity.ProjectileBehaviour); ok {
				if pro.Owner() == p.H() {
					p.Tx().RemoveEntity(et)
					_ = et.Close()
				}
			}
		}
	}
}

// clearEffects clears all effects from the user.
func (h *Handler) clearEffects(p *player.Player) {
	for _, ef := range p.Effects() {
		p.RemoveEffect(ef.Type())
	}
}

// resetCoolDowns resets all cooldowns of the user.
func (h *Handler) resetCoolDowns() {
	h.tagCombat.Reset()
	h.coolDownPearl.Reset()
	h.coolDownGoldenApple.Reset()
}

// spawnDeathNPC spawns a death NPC at the user's position.
func (h *Handler) spawnDeathNPC(p *player.Player, src world.DamageSource) {
	/*npc := player.New(p.Name(), p.Skin(), p.Position())
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
	}*/
}

// kill handles the death of the user.
func (h *Handler) kill(p *player.Player, src world.DamageSource) {
	u, err := data2.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	p.Tx().PlaySound(p.Position(), sound.Explosion{})
	if !u.Teams.DeathBan.Active() {
		h.handleTeamMemberDeath(p)
	}

	if k, ok := koth.Running(); ok {
		k.StopCapturing(p)
	}

	if conquest.Running() {
		for _, c := range conquest.All() {
			c.StopCapturing(p)
		}
	}

	h.stopCapturing(p)
	h.incrementDeath(p)
	h.issueDeathban(p)
	h.cancelStormBreak(p)
	h.spawnDeathNPC(p, src)
	h.clearEffects(p)
	h.clearOwnedEntities(p)
	h.resetCoolDowns()
	//unsafe.Session(p).EmptyUIInventory()

	DropContents(p)
	p.Inventory().Clear()
	p.Armour().Clear()
	p.SetHeldItems(item.Stack{}, item.Stack{})
	p.ResetFallDistance()
	p.Heal(20, effect.InstantHealingSource{})
	p.Extinguish()
	p.SetFood(20)

	h.lastClass.Store(class.Resolve(p))
	UpdateState(p)

	sortArmourEffects(p, h)
	sortClassEffects(p, h)

	if _, ok := eotw.Running(); ok {
		p.Disconnect(lang.Translatef(*u.Language, "death.eotw"))
	}

	internal.Deathban().Exec(func(tx *world.Tx) {
		tx.AddEntity(p.H())
	})
	unsafe.WritePacket(p, &packet.PlayerFog{
		Stack: []string{"minecraft:fog_default"},
	})

	p.Teleport(mgl64.Vec3{5, 13, 44})
	h.savePlayerData(p)
	if h.logger {
		h.death <- struct{}{}
	}
}

func (h *Handler) savePlayerData(p *player.Player) {
	dat, err := internal.LoadPlayerData(h.uuid)
	if err != nil {
		fmt.Println("error loading player data: ", err)
		return
	}
	var i, a = p.Inventory(), p.Armour()
	if i == nil || a == nil {
		return
	}
	dat.Inventory = p.Inventory()
	dat.Position = p.Position()

	err = internal.PlayerProvider().Save(h.uuid, dat, p.Tx().World())
	if err != nil {
		fmt.Println("error saving player data: ", err)
		return
	}
}

// stopCapturing stops the user from capturing a koth or a point.
func (h *Handler) stopCapturing(p *player.Player) {
	k, ok := koth.Running()
	if ok {
		k.StopCapturing(p)
	}
	ok = conquest.Running()
	if ok {
		for _, c := range conquest.All() {
			c.StopCapturing(p)
		}
	}
}

// cancelStormBreak cancels the storm breaker effect.
func (h *Handler) cancelStormBreak(p *player.Player) {
	armourHandler, ok := p.Armour().Inventory().Handler().(*ArmourHandler)
	if ok {
		armourHandler.stormBreakCancel()
	}
}

// incrementDeath increments the death count of the user.
func (h *Handler) incrementDeath(p *player.Player) {
	victim, err := data2.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	victim.Teams.Stats.Deaths += 1
	if victim.Teams.Stats.KillStreak > victim.Teams.Stats.BestKillStreak {
		victim.Teams.Stats.BestKillStreak = victim.Teams.Stats.KillStreak
	}
	victim.Teams.Stats.KillStreak = 0

	held, off := p.HeldItems()
	*victim.Teams.DeathInventory = inventoryData(held, off, p.Armour(), p.Inventory())
	data2.SaveUser(victim)
}

func inventoryData(held, off item.Stack, a *inventory.Armour, i *inventory.Inventory) model3.Inventory {
	return model3.Inventory{
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
func (h *Handler) issueDeathban(p *player.Player) {
	u, err := data2.LoadUserFromName(p.Name())
	if err != nil || u.Teams.DeathBan.Active() {
		return
	}

	u.Teams.DeathBan.Set(time.Minute * 20)
	u.Teams.DeathBanned = true

	data2.SaveUser(u)
}

// handleTeamMemberDeath handles the death of a team member.
func (h *Handler) handleTeamMemberDeath(p *player.Player) {
	if tm, ok := core.TeamRepository.FindByMemberName(p.Name()); ok {
		tm = tm.WithDTR(tm.DTR - 1).WithLastDeath(time.Now())
		if tm.Points > 0 {
			tm = tm.WithPoints(tm.Points - 1)
		}
		core.TeamRepository.Save(tm)

		for _, member := range tm.Members {
			if m, ok := Lookup(p.Tx(), member.Name); ok {
				u, _ := data2.LoadUserFromName(m.Name())
				m.Message(lang.Translatef(*u.Language, "team.member.death", p.Name(), tm.DTR))
			}
		}
	}
}

// canMessage returns true if the user can send a message.
func (h *Handler) canMessage() bool {
	return time.Since(h.lastMessage.Load()) > time.Second*1
}

// lastAttacker returns the last attacker of the user.
func (h *Handler) lastAttacker(tx *world.Tx) (*player.Player, bool) {
	if time.Since(h.lastAttackTime.Load()) > 15*time.Second {
		return nil, false
	}
	name := h.lastAttackerName.Load()
	if len(name) == 0 {
		return nil, false
	}
	return Lookup(tx, name)
}

// setLastAttacker sets the last attacker of the user.
func (h *Handler) setLastAttacker(p *player.Player) {
	h.lastAttackerName.Store(p.Name())
	h.lastAttackTime.Store(time.Now())
}

// resetLastAttacker resets the last attacker of the user.
func (h *Handler) resetLastAttacker() {
	h.lastAttackerName.Store("")
	h.lastAttackTime.Store(time.Time{})
}

// ShowArmor displays or removes players armor visibility from other players.
func (h *Handler) ShowArmor(p *player.Player, visible bool) {
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

	tm, teamFound := core.TeamRepository.FindByMemberName(p.Name())
	for pl := range internal.Players(p.Tx()) {
		s := unsafe.Session(pl)
		if teamFound {
			if !tm.Member(pl.Name()) {
				h.updateArmour(p, s, helmet, chestplate, leggings, boots)
			}
			continue
		}
		h.updateArmour(p, s, helmet, chestplate, leggings, boots)
	}
}

func (h *Handler) updateArmour(p *player.Player, s *session.Session, helmet item.Stack, chestplate item.Stack, leggings item.Stack, boots item.Stack) {
	unsafe.WritePacket(s, &packet.MobArmourEquipment{
		EntityRuntimeID: session_entityRuntimeID(s, p),
		Helmet:          instanceFromItem(s, helmet),
		Chestplate:      instanceFromItem(s, chestplate),
		Leggings:        instanceFromItem(s, leggings),
		Boots:           instanceFromItem(s, boots),
	})
}

func (h *Handler) SendClaimPillar(p *player.Player, pos cube.Pos) {
	for y := pos.Y(); y <= pos.Y()+50; y++ {
		delta := y - pos.Y()
		var b world.Block
		if delta%4 == 0 {
			b = block.Diamond{}
		} else {
			b = block.Glass{}
		}
		h.viewBlockUpdate(p, cube.Pos{pos.X(), y, pos.Z()}, b, 0)
	}
}

func (h *Handler) SendAirPillar(p *player.Player, pos cube.Pos) {
	for y := pos.Y(); y <= pos.Y()+50; y++ {
		h.viewBlockUpdate(p, cube.Pos{pos.X(), y, pos.Z()}, block.Air{}, 0)
	}
}

// revertMovement reverts the user's movement.
func (h *Handler) revertMovement(p *player.Player) {
	latency := p.Latency() * 2
	pos := p.Position()
	time.AfterFunc(latency, func() {
		if p == nil {
			return
		}
		g := p.OnGround()
		p.Teleport(pos)
		if !g {
			w := p.Tx()
			x := int(pos.X())
			z := int(pos.Z())
			for y := int(pos.Y()) - 1; y > 0; y-- {
				b := w.Block(cube.Pos{x, y, z})
				if _, ok := b.Model().(model2.Solid); ok {
					new := pos
					new.Add(mgl64.Vec3{0, float64(y + 1), 0})
					p.Teleport(new)
				} else {
					continue
				}
			}
		}
	})
}

// sendWall sends a wall to the user.
func (h *Handler) sendWall(p *player.Player, newPos cube.Pos, z area.Area, color item.Colour) {
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
					if blockReplaceable(p.Tx().Block(blockPos)) {
						h.viewBlockUpdate(p, blockPos, wallBlock, 0)
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
					if blockReplaceable(p.Tx().Block(blockPos)) {
						h.viewBlockUpdate(p, blockPos, wallBlock, 0)
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
func (h *Handler) clearWall(pl *player.Player) {
	h.wallBlocksMu.Lock()
	for p, duration := range h.wallBlocks {
		if duration <= 0 {
			delete(h.wallBlocks, p)
			h.viewBlockUpdate(pl, p, pl.Tx().Block(p), 0)
			//p.ShowParticle(p.Vec3(), particle.BlockForceField{})
			continue
		}
		h.wallBlocks[p] = duration - 0.1
	}
	h.wallBlocksMu.Unlock()
}

// viewBlockUpdate updates the block at the position passed for the user.
func (h *Handler) viewBlockUpdate(p *player.Player, pos cube.Pos, b world.Block, layer int) {
	s := unsafe.Session(p)
	s.ViewBlockUpdate(pos, b, layer)
}

// viewers returns a list of all viewers of the Player.
func (h *Handler) viewers(p *player.Player) []world.Viewer {
	viewers := p.Tx().Viewers(p.Position())
	s := unsafe.Session(p)

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
func subtractItem(p *player.Player, s item.Stack, d int) item.Stack {
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
	u, err := data2.LoadUserFromName(p.Name())
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
		if !hasEffectLevel(p, e) {
			p.AddEffect(e)
		}
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
func nearbyPlayers(p *player.Player, dist float64) []*player.Player {
	var pl []*player.Player
	for e := range p.Tx().Entities() {
		if e.Position().ApproxFuncEqual(p.Position(), func(f float64, f2 float64) bool {
			return math.Max(f, f2)-math.Min(f, f2) < dist
		}) {
			if target, ok := e.(*player.Player); ok && target != p {
				pl = append(pl, target)
			}
		}
	}
	return pl
}

// nearbyEnemies returns the nearby enemies of a certain distance from the user
func nearbyEnemies(p *player.Player, dist float64) []*player.Player {
	var pl []*player.Player
	tm, ok := core.TeamRepository.FindByMemberName(p.Name())
	if !ok {
		return nearbyPlayers(p, dist)
	}

	for _, target := range nearbyPlayers(p, dist) {
		if !tm.Member(target.Name()) {
			pl = append(pl, target)
		}
	}

	return pl
}

// nearbyAllies returns the nearby allies of a certain distance from the user
func nearbyAllies(p *player.Player, dist float64) []*player.Player {
	pl := []*player.Player{p}
	tm, ok := core.TeamRepository.FindByMemberName(p.Name())
	if !ok {
		return pl
	}

	for _, target := range nearbyPlayers(p, dist) {
		if tm.Member(target.Name()) {
			pl = append(pl, p)
		}
	}

	return pl
}

// nearbyHurtable returns the nearby faction members, enemies, no faction players (basically anyone not on timer) of a certain distance from the user
// nearbyPlayers returns the nearby users of a certain distance from the user
func nearbyHurtable(p *player.Player, dist float64) []*player.Player {
	var pl []*player.Player

	for _, target := range nearbyPlayers(p, dist) {
		t, _ := data2.LoadUserFromName(target.Name())
		if !t.Teams.PVP.Active() {
			pl = append(pl, target)
		}
	}

	return pl
}

// noinspection ALL
//
//go:linkname DropContents github.com/df-mc/dragonfly/server/player.(*Player).dropItems
func DropContents(*player.Player)

// noinspection ALL
//
//go:linkname instanceFromItem github.com/df-mc/dragonfly/server/session.instanceFromItem
func instanceFromItem(*session.Session, item.Stack) protocol.ItemInstance

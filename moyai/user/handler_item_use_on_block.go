package user

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	blck "github.com/moyai-network/teams/moyai/block"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/moyai-network/teams/moyai/crate"
	"github.com/moyai-network/teams/moyai/data"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	"github.com/moyai-network/teams/moyai/eotw"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/kit"
)

// HandleItemUseOnBlock ...
func (h *Handler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	w := h.p.World()
	held, left := h.p.HeldItems()
	b := w.Block(pos)

	if _, ok := held.Item().(world.Block); !ok {
		if h.coolDownItemUse.Active() {
			ctx.Cancel()
			return
		}
		h.coolDownItemUse.Set(time.Second / 10)
	} 

	c, crateFound := resolveCrateFromPosition(pos, b)
	if crateFound {
		ctx.Cancel()
		if !canOpenCrate(held, c) {
			moyai.Messagef(h.p, "crate.key.require", colour.StripMinecraftColour(c.Name()))
			return
		}
		openCrate(h.p, w, held, left, c)
		return
	}

	tm, teamErr := data.LoadTeamFromMemberName(h.p.Name())
	teams, _ := data.LoadAllTeams()
	posWithinProtected := posWithinProtectedArea(h.p, pos, teams)

	if _, ok := b.(block.ItemFrame); ok && posWithinProtected {
		ctx.Cancel()
		return
	}

	switch itm := held.Item().(type) {
	case item.Firework:
		ctx.Cancel()
	case item.Bucket:
		if posWithinProtected {
			ctx.Cancel()
			break
		}
	case item.EnderPearl:
		h.handlePearlUseOnBlock(ctx, itm, pos)
	case item.Hoe:
		ctx.Cancel()
		if itm.Tier != item.ToolTierDiamond {
			break
		}
		_, crowbar := held.Value("CROWBAR")
		if crowbar {
			if _, portal := b.(blck.PortalFrame); !portal {
				break
			}
			if posWithinProtected {
				moyai.Messagef(h.p, "team.claim.not-within")
				break
			}

			w.SetBlock(pos, block.Air{}, nil)
			it.DropFromPosition(w, pos.Vec3Middle(), item.NewStack(blck.PortalFrame{}, 1))
			break
		}

		_, ok := held.Value("CLAIM_WAND")
		if !ok {
			break
		}
		if teamErr != nil {
			break
		}

		if h.p.World() != moyai.Overworld() {
			break
		}

		if !tm.Leader(h.p.Name()) {
			moyai.Messagef(h.p, "team.not-leader")
			break
		}

		if tm.Claim != (area.Area{}) {
			moyai.Messagef(h.p, "team.has-claim")
			break
		}

		for _, a := range area.Protected(w) {
			var threshold float64 = 1
			message := "team.area.too-close"
			for _, k := range area.KOTHs(h.p.World()) {
				if a.Area == k.Area {
					threshold = 100
					message = "team.area.too-close.koth"
				}
			}

			if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
				moyai.Messagef(h.p, "team.area.already-claimed")
				return
			}
			if areaTooClose(a.Area, vec3ToVec2(pos.Vec3()), threshold) {
				moyai.Messagef(h.p, message)
				return
			}
		}

		for _, t := range teams {
			c := t.Claim
			if c != (area.Area{}) {
				continue
			}
			if c.Vec3WithinOrEqualXZ(pos.Vec3()) {
				moyai.Messagef(h.p, "team.area.already-claimed")
				return
			}
			if areaTooClose(c, vec3ToVec2(pos.Vec3()), 1) {
				moyai.Messagef(h.p, "team.area.too-close")
				return
			}
		}

		pn := 1
		if h.p.Sneaking() {
			pn = 2
			ar := area.NewArea(mgl64.Vec2{h.claimSelectionPos[0].X(), h.claimSelectionPos[0].Z()}, mgl64.Vec2{float64(pos.X()), float64(pos.Z())})
			x := ar.Max().X() - ar.Min().X()
			y := ar.Max().Y() - ar.Min().Y()
			a := x * y
			if a > 75*75 {
				moyai.Messagef(h.p, "team.claim.too-big")
				return
			}
			cost := int(a * 5)
			moyai.Messagef(h.p, "team.claim.cost", cost)
		}
		if h.claimSelectionPos[pn-1] != (mgl64.Vec3{}) {
			h.SendAirPillar(cube.PosFromVec3(h.claimSelectionPos[pn-1]))
		}
		h.claimSelectionPos[pn-1] = mgl64.Vec3{float64(pos.X()), float64(pos.Y()), float64(pos.Z())}
		h.SendClaimPillar(pos)
		moyai.Messagef(h.p, "team.claim.set-position", pn, mgl64.Vec2{float64(pos.X()), float64(pos.Z())})
	default:
		if !h.p.Sneaking() {
			break
		}
		useItem(h.p, held, left)
	}

	if _, ok := b.(block.Chest); ok && area.WarZone(w).Vec3WithinOrEqualXZ(pos.Vec3()) && !area.Spawn(w).Vec3WithinOrEqualXZ(pos.Vec3()) {
		return
	}

	switch bl := b.(type) {
	case block.Anvil:
		ctx.Cancel()
	case block.WoodFenceGate, block.Chest, block.WoodTrapdoor, block.WoodDoor, block.ItemFrame, block.Hopper:
		if posWithinProtected {
			h.revertMovement()
			ctx.Cancel()
			return
		}
		if h.coolDownBonedEffect.Active() {
			moyai.Messagef(h.p, "user.interaction.boned")
			ctx.Cancel()
			return
		}
	case block.Sign:
		ctx.Cancel()

		lines := strings.Split(bl.Front.Text, "\n")
		if len(lines) < 2 {
			return
		}

		title := strings.ToLower(colour.StripMinecraftColour(lines[0]))
		body := strings.ToLower(colour.StripMinecraftColour(lines[1]))

		choice := strings.ReplaceAll(title, " ", "")
		choice = strings.ReplaceAll(choice, "[", "")
		choice = strings.ReplaceAll(choice, "]", "")

		switch choice {
		case "elevator":
			handleElevatorSignInteraction(h.p, body, pos)
			return
		case "buy", "sell":
			handleShopSignInteraction(h.p, choice, lines)
			return
		case "kit":
			key := colour.StripMinecraftColour(lines[1])
			u, err := data.LoadUserFromName(h.p.Name())
			if err != nil {
				return
			}

			if u.Teams.DeathBan.Active() && key == "deathban" {
				kit.Apply(kit.Diamond{}, h.p)
			}
			return
		}
		if body == "have lives?" {
			u, err := data.LoadUserFromName(h.p.Name())
			if err != nil {
				return
			}
			if u.Teams.Lives <= 0 {
				moyai.Messagef(h.p, "lives.none")
				return
			}
			if u.Teams.DeathBan.Active() {
				u.Teams.DeathBan.Reset()
				u.Teams.DeathBanned = false
				u.Teams.Lives -= 1
				u.Teams.PVP.Set(time.Hour + (time.Millisecond * 500))
				fmt.Println("HandleItemUseOnBlock: pvp timer set to an hour for", h.p.Name())
				if !u.Teams.PVP.Paused() {
					u.Teams.PVP.TogglePause()
				}
				data.SaveUser(u)
				moyai.Overworld().AddEntity(h.p)
				h.p.Armour().Clear()
				h.p.Inventory().Clear()
				h.p.Teleport(mgl64.Vec3{0, 80})
			}
		}
	}
}

func useItem(p *player.Player, held, left item.Stack) {
	usable, ok := held.Item().(item.Usable)
	if ok {
		ctx := &item.UseContext{}
		usable.Use(p.World(), p, ctx)
		handleUseContext(p, ctx)
	}
}

func handleElevatorSignInteraction(p *player.Player, body string, pos cube.Pos) {
	blockFound := false
	switch body {
	case "up":
		for y := pos.Y() + 1; y < 256; y++ {
			if _, ok := p.World().Block(cube.Pos{pos.X(), y, pos.Z()}).(block.Air); ok {
				if !blockFound {
					moyai.Messagef(p, "elevator.no-block")
					return
				}
				if _, ok := p.World().Block(cube.Pos{pos.X(), y + 1, pos.Z()}).(block.Air); !ok {
					moyai.Messagef(p, "elevator.no-space")
					return
				}
				p.Teleport(pos.Vec3Middle().Add(mgl64.Vec3{0, float64(y - pos.Y()), 0}))
				break
			} else {
				blockFound = true
			}
		}
	case "down":
		for y := pos.Y() - 1; y > 0; y-- {
			if _, ok := p.World().Block(cube.Pos{pos.X(), y, pos.Z()}).(block.Air); ok {
				if !blockFound {
					moyai.Messagef(p, "elevator.no-space")
					return
				}
				if _, ok := p.World().Block(cube.Pos{pos.X(), y - 2, pos.Z()}).(block.Air); ok {
					moyai.Messagef(p, "elevator.no-space")
					return
				}
				p.Teleport(pos.Vec3Middle().Add(mgl64.Vec3{0, float64(y - pos.Y() - 1), 0}))
				break
			} else {
				blockFound = true
			}
		}
	}
}

// handleShopSignInteraction handles the interaction of a player with a shop sign. If the player interacts with
// a shop sign, the player is able to buy or sell items from the shop sign. The player's balance is updated
// accordingly.
func handleShopSignInteraction(p *player.Player, choice string, lines []string) {
	if !area.Spawn(p.World()).Vec3WithinOrEqualFloorXZ(p.Position()) {
		return
	}
	var stack item.Stack
	q, err := strconv.Atoi(lines[2])
	if err != nil {
		return
	}

	itm, vanillaItem := world.ItemByName("minecraft:"+strings.ReplaceAll(strings.ToLower(lines[1]), " ", "_"), 0)
	if lines[1] == "Crowbar" {
		stack = it.NewCrowBar()
	} else if vanillaItem {
		stack = item.NewStack(itm, q)
	} else {
		return
	}

	price, err := strconv.ParseFloat(strings.Trim(lines[3], "$"), 64)
	if err != nil {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	switch choice {
	case "buy":
		handleBuy(p, u, stack, price, lines[1])
	case "sell":
		handleSell(p, u, itm, q, price, lines[1])
	}
}

// handleBuy handles the purchase of an item from a shop sign. If the player can buy the item, the item is
// added to the player's inventory and the player's balance is reduced by the price of the item. If the player
// cannot buy the item, the player is sent a message informing them that they do not have enough balance.
func handleBuy(p *player.Player, u data.User, stack item.Stack, price float64, itemName string) {
	if u.Teams.Balance < price {
		moyai.Messagef(p, "shop.balance.insufficient")
		return
	}
	u.Teams.Balance = u.Teams.Balance - price
	data.SaveUser(u)
	it.AddOrDrop(p, stack)
	moyai.Messagef(p, "shop.buy.success", stack.Count(), itemName)
}

// handleSell handles the selling of an item to a shop sign. If the player can sell the item, the player's
// balance is increased by the price of the item. If the player cannot sell the item, the player is sent a
// message informing them that they cannot sell the item.
func handleSell(p *player.Player, u data.User, itm world.Item, q int, price float64, itemName string) {
	inv := p.Inventory()
	count := 0
	var items []item.Stack
	for _, slotItem := range inv.Slots() {
		n1, _ := itm.EncodeItem()
		if slotItem.Empty() {
			continue
		}
		n2, _ := slotItem.Item().EncodeItem()
		if n1 == n2 {
			count += slotItem.Count()
			items = append(items, slotItem)
		}
	}
	if count >= q {
		u.Teams.Balance = u.Teams.Balance + float64(count/q)*price
		data.SaveUser(u)
		moyai.Messagef(p, "shop.sell.success", count, itemName)
	} else {
		moyai.Messagef(p, "shop.sell.fail")
		return
	}
	for i, v := range items {
		if i >= count {
			break
		}
		amt := count - (count % q)
		if amt > 64 {
			amt = 64
		}
		_ = inv.RemoveItemFunc(amt, func(stack item.Stack) bool {
			return stack.Equal(v)
		})
	}
}

// resolveCrateFromPosition resolves a crate from a position and block. If a crate is found, the crate and
// true are returned. If no crate is found, nil and false are returned.
func resolveCrateFromPosition(pos cube.Pos, b world.Block) (crate.Crate, bool) {
	for _, c := range crate.All() {
		if _, ok := b.(block.Chest); ok && pos.Vec3Middle() == c.Position() {
			return c, true
		}
	}
	return nil, false
}

// canOpenCrate checks if a player can open a crate with the held item stack and the crate passed. If the
// player can open the crate, true is returned. If the player cannot open the crate, false is returned.
func canOpenCrate(held item.Stack, c crate.Crate) bool {
	_, ok := held.Value("crate-key_" + colour.StripMinecraftColour(c.Name()))
	return ok
}

// openCrate opens a crate for a player, removing a key from the player's inventory and giving the player
// the reward from the crate. The player is also sent a firework to celebrate the opening of the crate.
func openCrate(p *player.Player, w *world.World, held, left item.Stack, c crate.Crate) {
	it.AddOrDrop(p, ench.AddEnchantmentLore(crate.SelectReward(c)))
	p.SetHeldItems(subtractItem(p, held, 1), left)

	w.AddEntity(entity.NewFirework(c.Position().Add(mgl64.Vec3{0, 1, 0}), cube.Rotation{90, 90}, item.Firework{
		Duration: 0,
		Explosions: []item.FireworkExplosion{
			{
				Shape:   item.FireworkShapeStar(),
				Trail:   true,
				Colour:  colour.RandomColour(),
				Twinkle: true,
			},
		},
	}))
}

// posWithinProtectedArea checks if a position is within a protected area. If the position is within a protected
// area, true is returned. If the position is not within a protected area, false is returned. The player passed
// is used to check if the player is a member of a team that has a claim in the area.
func posWithinProtectedArea(p *player.Player, pos cube.Pos, teams []data.Team) bool {
	if p.GameMode() == world.GameModeCreative {
		return false
	}
	w := p.World()

	var posWithinProtected bool
	for _, a := range area.Protected(w) {
		if a.Vec3WithinOrEqualXZ(pos.Vec3()) {
			posWithinProtected = true
			break
		}
	}
	if posWithinProtected {
		return posWithinProtected
	}

	_, eotwRunning := eotw.Running()
	if eotwRunning {
		return false
	}

	for _, t := range teams {
		cl := t.Claim
		_, eotwRunning := eotw.Running()
		if !t.Member(p.Name()) && cl.Vec3WithinOrEqualXZ(pos.Vec3()) && t.DTR > 0 && !eotwRunning {
			posWithinProtected = true
			break
		}
	}
	return posWithinProtected
}

// handlePearlUseOnBlock handles the use of an ender pearl on a block. If the ender pearl is used on a block
// successfully, true is returned. If the ender pearl is not used on a block, false is returned.
func (h *Handler) handlePearlUseOnBlock(ctx *event.Context, pearl item.EnderPearl, pos cube.Pos) {
	p, w := h.p, h.p.World()

	if f, ok := w.Block(pos).(block.WoodFenceGate); ok && f.Open {
		h.handlePearlUse(ctx)
		if ctx.Cancelled() {
			*ctx = *event.C()
			return
		}
		useCtx := &item.UseContext{}
		pearl.Use(w, p, useCtx)
		handleUseContext(p, useCtx)
		ctx.Cancel()
	}
}

func areaTooClose(area area.Area, pos mgl64.Vec2, threshold float64) bool {
	var vectors []mgl64.Vec2
	for x := -threshold; x <= threshold; x++ {
		for y := -threshold; y <= threshold; y++ {
			vectors = append(vectors, mgl64.Vec2{pos.X() + x, pos.Y() + y})
		}
	}

	for _, v := range vectors {
		if area.Vec2WithinOrEqual(v) {
			return true
		}
	}
	return false
}

// noinspection ALL
//
//go:linkname handleUseContext github.com/df-mc/dragonfly/server/player.(*Player).handleUseContext
func handleUseContext(p *player.Player, ctx *item.UseContext)

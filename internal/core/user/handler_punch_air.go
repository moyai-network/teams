package user

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core/area"
	data2 "github.com/moyai-network/teams/internal/core/data"
	it "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/ports/model"
	"github.com/moyai-network/teams/pkg/lang"
)

// HandlePunchAir ...
func (h *Handler) HandlePunchAir(ctx *player.Context) {
	p := ctx.Val()
	u, err := data2.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	held, _ := p.HeldItems()
	typ, ok := it.PartnerItem(held)
	if ok {
		if cd := h.coolDownGlobalAbilities; cd.Active() {
			p.SendJukeboxPopup(lang.Translatef(*u.Language, "popup.cooldown.partner_item", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		if spi := h.coolDownSpecificAbilities; spi.Active(typ) {
			p.SendJukeboxPopup(lang.Translatef(*u.Language, "popup.cooldown.partner_item.item", typ.Name(), spi.Remaining(typ).Seconds()))
		} else {
			p.SendJukeboxPopup(lang.Translatef(*u.Language, "popup.ready.partner_item.item", typ.Name()))
		}
		ctx.Cancel()
		return
	}

	if !p.Sneaking() {
		return
	}

	i, ok := held.Item().(item.Hoe)
	if !ok || i.Tier != item.ToolTierDiamond {
		return
	}
	_, ok = held.Value("CLAIM_WAND")
	if !ok {
		return
	}

	t, err := data2.LoadTeamFromMemberName(p.Name())
	if err != nil {
		return
	}

	if !t.Leader(u.Name) {
		internal.Messagef(p, "team.not-leader")
		return
	}
	if t.Claim != (area.Area{}) {
		internal.Messagef(p, "team.has-claim")
		return
	}

	handleClaimSelection(p, h, t)
}

// handleClaimSelection handles the selection of a claim area.
func handleClaimSelection(p *player.Player, h *Handler, tm model.Team) {
	pos := h.claimSelectionPos
	if pos[0] == (mgl64.Vec3{}) || pos[1] == (mgl64.Vec3{}) {
		internal.Messagef(p, "team.claim.select-before")
		return
	}

	// Create a new area claim.
	claim := area.NewArea(mgl64.Vec2{pos[0].X(), pos[0].Z()}, mgl64.Vec2{pos[1].X(), pos[1].Z()})
	if checkExistingClaims(p, h, claim) {
		return
	}

	// Check if the claim area is too big.
	if ar := calculateArea(claim); ar > 75*75 {
		internal.Messagef(p, "team.claim.too-big")
		return
	}

	// Check if the claim area is too small.
	if ar := calculateArea(claim); ar < 5*5 {
		internal.Messagef(p, "team.claim.too-small")
		return
	}

	// Calculate the cost of claiming the area.
	cost := calculateCost(claim)
	if !checkBalance(tm, cost) {
		internal.Messagef(p, "team.claim.no-money")
		return
	}

	// Clear blocks in the claimed area.
	clearBlocksInRange(p, h, pos[0])
	clearBlocksInRange(p, h, pos[1])

	// Set the claim area for the team and save it.
	tm.Claim = claim
	data2.SaveTeam(tm)

	internal.Messagef(p, "command.claim.success", pos[0], pos[1], cost)
}

// checkExistingClaims checks if the new claim overlaps with existing claims.
func checkExistingClaims(p *player.Player, h *Handler, claim area.Area) bool {
	blocksPos := generateBlocksPos(claim)
	w := p.Tx().World()
	for _, a := range area.Protected(w) {
		var threshold float64 = 1
		if checkAreaOverlap(p, h, a.Area, blocksPos, threshold) {
			return true
		}
	}

	teams, err := data2.LoadAllTeams()
	if err != nil {
		return true
	}
	for _, tm := range teams {
		c := tm.Claim
		if c == (area.Area{}) {
			continue
		}

		if checkAreaOverlap(p, h, c, blocksPos, 1) {
			return true
		}
	}

	return false
}

// checkAreaOverlap checks if the new claim overlaps with an existing claim.
func checkAreaOverlap(p *player.Player, h *Handler, existingClaim area.Area, blocksPos []cube.Pos, threshold float64) bool {
	pos := h.claimSelectionPos
	p0 := mgl64.Vec2{pos[0].X(), pos[0].Z()}
	p1 := mgl64.Vec2{pos[1].X(), pos[1].Z()}

	// Check if new claim overlaps with existing claim.
	for _, b := range blocksPos {
		if existingClaim.Vec3WithinOrEqualXZ(b.Vec3()) {
			internal.Messagef(p, "team.area.already-claimed")
			return true
		}
	}

	// Check if corners of new claim are within existing claim.
	if existingClaim.Vec2WithinOrEqual(p0) || existingClaim.Vec2WithinOrEqual(p1) ||
		areaTooClose(existingClaim, p0, threshold) || areaTooClose(existingClaim, p1, threshold) {
		internal.Messagef(p, "team.area.already-claimed")
		return true
	}

	return false
}

// generateBlocksPos generates the positions of blocks within the claim area.
func generateBlocksPos(claim area.Area) []cube.Pos {
	var blocksPos []cube.Pos
	mn := claim.Min()
	mx := claim.Max()
	for x := mn[0]; x <= mx[0]; x++ {
		for y := mn[1]; y <= mx[1]; y++ {
			blocksPos = append(blocksPos, cube.PosFromVec3(mgl64.Vec3{x, 0, y}))
		}
	}

	return blocksPos
}

// calculateArea calculates the area of the claim.
func calculateArea(claim area.Area) float64 {
	x := claim.Max().X() - claim.Min().X()
	y := claim.Max().Y() - claim.Min().Y()
	return x * y
}

// calculateCost calculates the cost of claiming the area.
func calculateCost(claim area.Area) int {
	ar := calculateArea(claim)
	return int(ar * 5)
}

// checkBalance checks if the team has enough balance to claim the area.
func checkBalance(tm model.Team, cost int) bool {
	return int(tm.Balance) >= cost
}

// clearBlocksInRange clears blocks within the specified range.
func clearBlocksInRange(p *player.Player, h *Handler, pos mgl64.Vec3) {
	for y := pos.Y(); y <= pos.Y()+25; y++ {
		h.viewBlockUpdate(p, cube.Pos{int(pos.X()), int(y), int(pos.Z())}, block.Air{}, 0)
	}
}

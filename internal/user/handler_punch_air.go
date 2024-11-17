package user

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/area"
	"github.com/moyai-network/teams/internal/data"
	it "github.com/moyai-network/teams/internal/item"
	"github.com/moyai-network/teams/pkg/lang"
)

// HandlePunchAir ...
func (h *Handler) HandlePunchAir(ctx *event.Context) {
	p := h.p
	u, err := data.LoadUserFromName(p.Name())
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

	t, err := data.LoadTeamFromMemberName(p.Name())
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

	handleClaimSelection(h, t)
}

// handleClaimSelection handles the selection of a claim area.
func handleClaimSelection(h *Handler, tm data.Team) {
	pos := h.claimSelectionPos
	if pos[0] == (mgl64.Vec3{}) || pos[1] == (mgl64.Vec3{}) {
		internal.Messagef(h.p, "team.claim.select-before")
		return
	}

	// Create a new area claim.
	claim := area.NewArea(mgl64.Vec2{pos[0].X(), pos[0].Z()}, mgl64.Vec2{pos[1].X(), pos[1].Z()})
	if checkExistingClaims(h, claim) {
		return
	}

	// Check if the claim area is too big.
	if ar := calculateArea(claim); ar > 75*75 {
		internal.Messagef(h.p, "team.claim.too-big")
		return
	}

	// Check if the claim area is too small.
	if ar := calculateArea(claim); ar < 5*5 {
		internal.Messagef(h.p, "team.claim.too-small")
		return
	}

	// Calculate the cost of claiming the area.
	cost := calculateCost(claim)
	if !checkBalance(tm, cost) {
		internal.Messagef(h.p, "team.claim.no-money")
		return
	}

	// Clear blocks in the claimed area.
	clearBlocksInRange(h, pos[0])
	clearBlocksInRange(h, pos[1])

	// Set the claim area for the team and save it.
	tm.Claim = claim
	data.SaveTeam(tm)

	internal.Messagef(h.p, "command.claim.success", pos[0], pos[1], cost)
}

// checkExistingClaims checks if the new claim overlaps with existing claims.
func checkExistingClaims(h *Handler, claim area.Area) bool {
	blocksPos := generateBlocksPos(claim)
	w := h.p.World()
	for _, a := range area.Protected(w) {
		var threshold float64 = 1
		message := "team.area.too-close"
		for _, k := range area.KOTHs(w) {
			if a.Area == k.Area {
				threshold = 100
				message = "team.area.too-close.koth"
			}
		}

		if checkAreaOverlap(h, a.Area, blocksPos, threshold, message) {
			return true
		}
	}

	teams, err := data.LoadAllTeams()
	if err != nil {
		return true
	}
	for _, tm := range teams {
		c := tm.Claim
		if c == (area.Area{}) {
			continue
		}

		if checkAreaOverlap(h, c, blocksPos, 1, "team.area.already-claimed") {
			return true
		}
	}

	return false
}

// checkAreaOverlap checks if the new claim overlaps with an existing claim.
func checkAreaOverlap(h *Handler, existingClaim area.Area, blocksPos []cube.Pos, threshold float64, message string) bool {
	pos := h.claimSelectionPos
	p0 := mgl64.Vec2{pos[0].X(), pos[0].Z()}
	p1 := mgl64.Vec2{pos[1].X(), pos[1].Z()}

	// Check if new claim overlaps with existing claim.
	for _, b := range blocksPos {
		if existingClaim.Vec3WithinOrEqualXZ(b.Vec3()) {
			internal.Messagef(h.p, "team.area.already-claimed")
			return true
		}
	}

	// Check if corners of new claim are within existing claim.
	if existingClaim.Vec2WithinOrEqual(p0) || existingClaim.Vec2WithinOrEqual(p1) ||
		areaTooClose(existingClaim, p0, threshold) || areaTooClose(existingClaim, p1, threshold) {
		internal.Messagef(h.p, "team.area.already-claimed")
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
func checkBalance(tm data.Team, cost int) bool {
	return int(tm.Balance) >= cost
}

// clearBlocksInRange clears blocks within the specified range.
func clearBlocksInRange(h *Handler, pos mgl64.Vec3) {
	for y := pos.Y(); y <= pos.Y()+25; y++ {
		h.viewBlockUpdate(cube.Pos{int(pos.X()), int(y), int(pos.Z())}, block.Air{}, 0)
	}
}

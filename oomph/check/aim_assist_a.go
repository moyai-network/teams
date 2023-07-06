package check

import (
	"math"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// AimAssistA checks for invalid headYaw to yaw patterns.
type AimAssistA struct {
	basic
}

// NewAimAssistA creates a new AimAssistA check.
func NewAimAssistA() *AimAssistA {
	return &AimAssistA{}
}

// Name ...
func (*AimAssistA) Name() (string, string) {
	return "AimAssistA", "A"
}

// Description ...
func (*AimAssistA) Description() string {
	return "This checks for invalid headYaw to yaw patterns."
}

// MaxViolations ...
func (*AimAssistA) MaxViolations() float64 {
	return 15
}

// Process ...
func (a *AimAssistA) Process(p Processor, pk packet.Packet) bool {
	if pk, ok := pk.(*packet.PlayerAuthInput); ok {
		var yaw float64
		if pk.Yaw > 0 {
			yaw = 0
		} else {
			yaw = 360
		}

		expHeadYaw := math.Mod(yaw+float64(pk.Yaw), 360)
		diff := math.Mod(math.Abs(expHeadYaw-float64(pk.Yaw)), 360)
		roundedDiff := math.Round(diff*10000) / 10000

		if diff > 5e-5 && roundedDiff != 360 && pk.HeadYaw > 0 {
			a.buffer++
			if a.buffer >= 3 {
				p.Flag(a, a.violationAfterTicks(p.ClientTick(), 40), map[string]any{
					"diff": roundedDiff,
				})
			}
		} else if pk.HeadYaw < 0 {
			expectedHeadYaw := math.Mod(float64(pk.HeadYaw), 180)
			diff := math.Mod(math.Abs(expectedHeadYaw-float64(pk.HeadYaw)), 360)
			roundedDiff := math.Round(diff*10000) / 10000

			if diff > 5e-5 && roundedDiff != 360 {
				a.buffer++
				if a.buffer >= 3 {
					p.Flag(a, a.violationAfterTicks(p.ClientTick(), 40), map[string]any{
						"diff": roundedDiff,
					})
				}
			} else {
				a.buffer = math.Max(float64(a.buffer)-0.025, 0)
				a.AddViolation(0.1)
			}
		} else {
			a.buffer = math.Max(float64(a.buffer)-0.025, 0)
			a.AddViolation(0.1)
		}
	}

	return false
}

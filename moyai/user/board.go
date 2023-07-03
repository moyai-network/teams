package user

import (
	"fmt"
	"time"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"

	// "github.com/moyai-network/hcf/moyai/koth"
	// "github.com/moyai-network/hcf/moyai/sotw"
	"github.com/moyai-network/moose/lang"
)

// sendBoard sends all user boards
func sendBoard(p *player.Player) {
	t := time.NewTicker(50 * time.Millisecond)
	l := p.Locale()
	for {
		select {
		case <-t.C:

			sb := scoreboard.New(lang.Translatef(l, "scoreboard.title"))
			_, _ = sb.WriteString("§r\uE000")
			sb.RemovePadding()

			_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.sotw", "1:30:00"))
			_, _ = sb.WriteString("\uE000")
			_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))

			p.RemoveScoreboard()
			p.SendScoreboard(sb)
		}
	}
}

func compareLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, l := range a {
		if l != b[i] {
			return false
		}
	}
	return true
}

func parseDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

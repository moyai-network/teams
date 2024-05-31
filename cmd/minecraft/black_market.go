package minecraft

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/moyai"
	"math/rand"
	"time"
)

// tickBlackMarket runs a ticker that checks every 15 minutes if the black market should be opened. The black
// market is opened randomly with a 25% chance every 15 minutes, but only if the last black market was opened
// more than an hour ago.
func tickBlackMarket(srv *server.Server) {
	t := time.NewTicker(time.Minute * 15)
	defer t.Stop()

	for range t.C {
		if time.Since(moyai.LastBlackMarket()) < time.Hour {
			continue
		}

		if rand.Intn(4) == 0 {
			moyai.SetLastBlackMarket(time.Now())
			for _, p := range srv.Players() {
				p.PlaySound(sound.BarrelOpen{})
				p.PlaySound(sound.FireworkHugeBlast{})
				p.PlaySound(sound.FireworkLaunch{})
				p.PlaySound(sound.Note{})
				moyai.Broadcastf("blackmarket.opened")
			}
		}
	}
}

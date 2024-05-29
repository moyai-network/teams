package moyai

import (
	"time"

	"github.com/df-mc/dragonfly/server/player"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
)

var (
	// srv is the server instance of the Moyai server.
	srv *server.Server
	// end is the world of the End dimension.
	end *world.World
	// nether is the world of the Nether dimension.
	nether *world.World

	// lastBlackMarket is the time at which the last black market was opened.
	lastBlackMarket time.Time
	// blackMarketOpened is the time at which the black market was last opened.
	blackMarketOpened time.Time
)

func Overworld() *world.World {
	return srv.World()
}

func End() *world.World {
	return end
}

func Nether() *world.World {
	return nether
}

func Players() []*player.Player {
	return srv.Players()
}

func NewServer(config server.Config) *server.Server {
	srv = config.New()
	return srv
}

func LastBlackMarket() time.Time {
	return lastBlackMarket
}

func SetLastBlackMarket(t time.Time) {
	lastBlackMarket = t
}

func BlackMarketOpened() time.Time {
	return blackMarketOpened
}

func SetBlackMarketOpened(t time.Time) {
	blackMarketOpened = t
}

func ConfigureDimensions(reg world.EntityRegistry, netherFolder, endFolder string) (*world.World, *world.World) {
	endProv, err := mcdb.Open(endFolder)
	if err != nil {
		panic(err)
	}
	end = world.Config{
		Provider: endProv,
		Dim:      world.End,
		Entities: reg,
	}.New()

	netherProv, err := mcdb.Open(netherFolder)
	if err != nil {
		panic(err)
	}

	nether = world.Config{
		Provider: netherProv,
		Dim:      world.Nether,
		Entities: reg,
	}.New()
	return nether, end
}

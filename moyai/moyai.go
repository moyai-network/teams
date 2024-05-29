package moyai

import (
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
)

var srv *server.Server
var end *world.World
var lastBlackMarket time.Time
var blackMarketOpened time.Time

func Server() *server.Server {
	return srv
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

func End() *world.World {
	return end
}

func ConfigureEnd(reg world.EntityRegistry) *world.World {
	prov, err := mcdb.Open("assets/end")
	if err != nil {
		panic(err)
	}
	end = world.Config{
		Provider: prov,
		Dim:      world.End,
		Entities: reg,
	}.New()
	return end
}

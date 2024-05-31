package moyai

import (
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/df-mc/dragonfly/server/player"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
)

func init() {
	go func() {
		for Overworld() == nil {
			<-time.After(time.Millisecond)
			continue
		}
		tickAirDrop(Overworld())
	}()
}

var (
	// srv is the server instance of the Moyai server.
	srv *server.Server
	// end is the world of the End dimension.
	end *world.World
	// nether is the world of the Nether dimension.
	nether *world.World
	// deathban is the world of the Deathban arena.
	deathban *world.World

	// lastBlackMarket is the time at which the last black market was opened.
	lastBlackMarket time.Time
	// blackMarketOpened is the time at which the black market was last opened.
	blackMarketOpened time.Time
)

func Worlds() []*world.World {
	return []*world.World{Overworld(), Nether(), End()}
}

func Overworld() *world.World {
	if srv == nil {
		return nil
	}
	return srv.World()
}

func End() *world.World {
	return end
}

func Nether() *world.World {
	return nether
}

func Deathban() *world.World {
	return deathban
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

func ConfigureDeathban(reg world.EntityRegistry, folder string) *world.World {
	deathbanProv, err := mcdb.Open(folder)
	if err != nil {
		panic(err)
	}
	deathban = world.Config{
		Provider:        deathbanProv,
		Dim:             world.Overworld,
		Entities:        reg,
		RandomTickSpeed: -1,
	}.New()

	return deathban
}

func Close() {
	destroyAirDrop(srv.World(), lastDropPos)

	time.Sleep(time.Millisecond * 500)
	data.FlushCache()

	sotw.Save()
	err := Nether().Close()
	if err != nil {
		logrus.Fatalln("close nether: %v", err)
	}
	End().Close()
	if err := srv.Close(); err != nil {
		logrus.Fatalln("close server: %v", err)
	}
}

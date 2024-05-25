package moyai

import (
	"time"

	"github.com/df-mc/dragonfly/server"
)

var srv *server.Server
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
package team

import (
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
)

func Broadcast(t data.Team, key string, args ...interface{}) {
	for _, m := range t.Members {
		if p, ok := user.Lookup(m.XUID); ok {
			p.Message(lang.Translatef(p.Locale(), key, args))
		}
	}
}

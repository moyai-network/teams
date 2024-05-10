package moyai

import (
	"fmt"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"net"
)

type Allower struct {
	whitelisted bool
}

func NewAllower(whitelisted bool) *Allower {
	return &Allower{
		whitelisted: whitelisted,
	}
}

func (a *Allower) Allow(addr net.Addr, d login.IdentityData, c login.ClientData) (string, bool) {
	u, err := data.LoadUserOrCreate(d.DisplayName, d.XUID)
	if err != nil {
		fmt.Println(err)
		return lang.Translatef(u.Language, "user.data.load.error"), false
	}
	/*if !strings.HasPrefix(addr.String(), "127.0.0.1") {
		return text.Colourf("<red>Please connect via the main hub: moyai.pro:19132</red>"), false
	}*/
	if a.whitelisted && !u.Whitelisted {
		return lang.Translatef(u.Language, "moyai.whitelisted"), false
	}
	return "", true
}

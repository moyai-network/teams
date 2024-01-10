package moyai

import (
	"fmt"
	"net"

	"github.com/moyai-network/moose/data"
	"github.com/moyai-network/moose/lang"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"golang.org/x/text/language"
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
	u, err := data.LoadUserOrCreate(d.DisplayName)
	if err != nil {
		fmt.Println(err)
		return lang.Translatef(language.English, "user.data.load.error"), false
	}
	if a.whitelisted && !u.Whitelisted {
		return lang.Translatef(language.English, "moyai.whitelisted"), false
	}
	return "", true
}

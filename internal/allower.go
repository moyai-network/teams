package internal

import (
	"fmt"
	data2 "github.com/moyai-network/teams/internal/core/data"
	"github.com/moyai-network/teams/internal/core/eotw"
	model2 "github.com/moyai-network/teams/internal/model"
	"net"
	"strings"

	"github.com/hako/durafmt"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/unickorn/strutils"
)

type Allower struct {
	whitelisted bool
}

func NewAllower(whitelisted bool) *Allower {
	return &Allower{
		whitelisted: whitelisted,
	}
}

func (a *Allower) Allow(_ net.Addr, d login.IdentityData, _ login.ClientData) (string, bool) {
	u, err := data2.LoadUserOrCreate(d.DisplayName, d.XUID)
	if err != nil {
		fmt.Printf("Failed to load user: %v\n", err)
		return lang.Translatef(model2.Language{}, "user.data.load.error"), false
	}

	if a.whitelisted && !u.Whitelisted {
		return lang.Translatef(*u.Language, "internal.whitelisted"), false
	}
	if _, ok := eotw.Running(); ok {
		return lang.Translatef(*u.Language, "internal.eotw"), false
	}
	if !u.Teams.Ban.Expired() {
		description := lang.Translatef(*u.Language, "user.ban.description", strings.TrimSpace(u.Teams.Ban.Reason), durafmt.ParseShort(u.Teams.Ban.Remaining()))
		if strings.EqualFold(u.Name, d.DisplayName) {
			return strutils.CenterLine(lang.Translatef(*u.Language, "user.ban.header") + "\n" + description), false
		}
	}
	var users []model2.User
	ssid, err := data2.LoadUsersFromSelfSignedID(u.SelfSignedID)
	if err == nil {
		users = append(users, ssid...)
	}
	did, err := data2.LoadUsersFromDeviceID(u.DeviceID)
	if err == nil {
		users = append(users, did...)
	}
	for _, u := range users {
		if !u.Teams.Ban.Expired() {
			if u.Teams.Ban.Permanent {
				description := lang.Translatef(*u.Language, "user.blacklist.description", strings.TrimSpace(u.Teams.Ban.Reason))
				if strings.EqualFold(u.Name, d.DisplayName) {
					return strutils.CenterLine(lang.Translatef(*u.Language, "user.blacklist.header") + "\n" + description), false
				}
				return strutils.CenterLine(lang.Translatef(*u.Language, "user.blacklist.header.alt") + "\n" + description), false
			}
			description := lang.Translatef(*u.Language, "user.ban.description", strings.TrimSpace(u.Teams.Ban.Reason), durafmt.ParseShort(u.Teams.Ban.Remaining()))
			if strings.EqualFold(u.Name, d.DisplayName) {
				return strutils.CenterLine(lang.Translatef(*u.Language, "user.ban.header") + "\n" + description), false
			}
			return strutils.CenterLine(lang.Translatef(*u.Language, "user.ban.header.alt") + "\n" + description), false
		}
	}

	return "", true
}

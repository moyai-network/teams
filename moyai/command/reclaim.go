package command

import (
	"strings"
	"unicode"

	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Reclaim is a command that allows players to reclaim their partner package.
type Reclaim struct{}

// ReclaimReset is a command that allows admins to reset the reclaim cooldown.
type ReclaimReset struct {
	Sub cmd.SubCommand `cmd:"reset"`
}

func (Reclaim) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	h, ok := user.Lookup(p.Name())
	if !ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	if u.Teams.Reclaimed {
		moyai.Messagef(p, "user.reclaimed")
		return
	}
	u.Teams.Reclaimed = true

	highest := u.Roles.Highest()
	for _, r := range u.Roles.All() {
		if r == (role.Operator{}) {
			continue
		}

		var items []item.Stack
		var lives int

		switch r {
		case role.Default{}:
			items = append(items, it.NewPartnerPackage(3))
			items = append(items, it.NewKey(it.KeyTypePartner, 3))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 1))
			items = append(items, it.NewKey(it.KeyTypeMenes, 2))
			items = append(items, it.NewKey(it.KeyTypeRamses, 3))
		case role.Trial{}:
			items = append(items, it.NewPartnerPackage(3))
			items = append(items, it.NewKey(it.KeyTypePartner, 4))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 2))
			items = append(items, it.NewKey(it.KeyTypeMenes, 3))
			items = append(items, it.NewKey(it.KeyTypeRamses, 4))
			lives = 3
		/*case role.Khufu{}:
			items = append(items, it.NewPartnerPackage(4))
			items = append(items, it.NewKey(it.KeyTypeRamses, 10))
			items = append(items, it.NewKey(it.KeyTypeMenes, 5))
			items = append(items, it.NewKey(it.KeyTypePartner, 3))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 1))
			lives = 30
		case role.Ramses{}:
			items = append(items, it.NewPartnerPackage(5))
			items = append(items, it.NewKey(it.KeyTypeRamses, 15))
			items = append(items, it.NewKey(it.KeyTypeMenes, 10))
			items = append(items, it.NewKey(it.KeyTypePartner, 6))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 2))
			lives = 45
		case role.Menes{}, role.Mod{}, role.Trial{}:
			items = append(items, it.NewPartnerPackage(7))
			items = append(items, it.NewKey(it.KeyTypeRamses, 20))
			items = append(items, it.NewKey(it.KeyTypeMenes, 15))
			items = append(items, it.NewKey(it.KeyTypePartner, 12))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 4))
			lives = 60
		case role.dono{}:
			items = append(items, it.NewPartnerPackage(9))
			items = append(items, it.NewKey(it.KeyTypeRamses, 30))
			items = append(items, it.NewKey(it.KeyTypeMenes, 20))
			items = append(items, it.NewKey(it.KeyTypePartner, 18))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 8))
			lives = 75*/
		case role.Donor1{}, role.Manager{}, role.Admin{}, role.Owner{}:
			items = append(items, it.NewPartnerPackage(10))
			items = append(items, it.NewKey(it.KeyTypeRamses, 30))
			items = append(items, it.NewKey(it.KeyTypeMenes, 20))
			items = append(items, it.NewKey(it.KeyTypePartner, 20))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 10))
			lives = 85
		}
		for _, i := range items {
			it.AddOrDrop(h, i)
		}

		u.Teams.Lives += lives

		var itemNames []string
		for _, i := range items {
			itemNames = append(itemNames, text.Colourf("<red>%dx</red> %s", i.Count(), i.CustomName()))
		}
		nm := []rune(r.Name())
		moyai.Broadcastf("user.reclaim", highest.Color(p.Name()), r.Color(string(append([]rune{unicode.ToUpper(nm[0])}, nm[1:]...))), strings.Join(itemNames, ", "), lives)
	}
	data.SaveUser(u)
}

// Run ...
func (ReclaimReset) Run(_ cmd.Source, _ *cmd.Output) {
	for _, p := range moyai.Players() {
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			continue
		}

		u.Teams.Reclaimed = false
		data.SaveUser(u)
	}
}

// Allow ...
func (ReclaimReset) Allow(s cmd.Source) bool {
	return allow(s, true, role.Operator{})
}

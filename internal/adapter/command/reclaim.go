package command

import (
	"github.com/df-mc/dragonfly/server/world"
	data2 "github.com/moyai-network/teams/internal/core/data"
	item2 "github.com/moyai-network/teams/internal/core/item"
	rls "github.com/moyai-network/teams/internal/core/roles"
	"github.com/moyai-network/teams/internal/core/user"
	"strings"
	"unicode"

	"github.com/moyai-network/teams/pkg/lang"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Reclaim is a command that allows players to reclaim their partner package.
type Reclaim struct{}

// ReclaimReset is a command that allows admins to reset the reclaim cooldown.
type ReclaimReset struct {
	Sub     cmd.SubCommand             `cmd:"reset"`
	Targets cmd.Optional[[]cmd.Target] `cmd:"target"`
}

func (Reclaim) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	h, ok := user.Lookup(tx, p.Name())
	if !ok {
		return
	}

	u, err := data2.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	if u.Teams.Reclaimed {
		internal.Messagef(p, "user.reclaimed")
		return
	}
	u.Teams.Reclaimed = true

	highest := u.Roles.Highest()
	for _, r := range u.Roles.All() {
		if r == rls.Operator() || r == rls.Voter() || r == rls.Nitro() {
			continue
		}

		var items []item.Stack
		var lives int

		switch r {
		case rls.Default():
			items = append(items, item2.NewSpecialItem(item2.PartnerPackageType{}, 3))
			items = append(items, item2.NewKey(item2.KeyTypePartner, 3))
			items = append(items, item2.NewKey(item2.KeyTypePharaoh, 1))
			items = append(items, item2.NewKey(item2.KeyTypeMenes, 2))
			items = append(items, item2.NewKey(item2.KeyTypeRamses, 3))
		case rls.Trial():
			items = append(items, item2.NewSpecialItem(item2.PartnerPackageType{}, 3))
			items = append(items, item2.NewKey(item2.KeyTypePartner, 4))
			items = append(items, item2.NewKey(item2.KeyTypePharaoh, 2))
			items = append(items, item2.NewKey(item2.KeyTypeMenes, 3))
			items = append(items, item2.NewKey(item2.KeyTypeRamses, 4))
			lives = 3
		case rls.Khufu(), rls.Media():
			items = append(items, item2.NewSpecialItem(item2.PartnerPackageType{}, 4))
			items = append(items, item2.NewKey(item2.KeyTypeRamses, 10))
			items = append(items, item2.NewKey(item2.KeyTypeMenes, 5))
			items = append(items, item2.NewKey(item2.KeyTypePartner, 3))
			items = append(items, item2.NewKey(item2.KeyTypePharaoh, 1))
			lives = 30
		case rls.Ramses(), rls.Famous():
			items = append(items, item2.NewSpecialItem(item2.PartnerPackageType{}, 5))
			items = append(items, item2.NewKey(item2.KeyTypeRamses, 15))
			items = append(items, item2.NewKey(item2.KeyTypeMenes, 10))
			items = append(items, item2.NewKey(item2.KeyTypePartner, 6))
			items = append(items, item2.NewKey(item2.KeyTypePharaoh, 2))
			lives = 45
		case rls.Menes(), rls.Mod(), rls.Partner(), rls.Trial():
			items = append(items, item2.NewSpecialItem(item2.PartnerPackageType{}, 7))
			items = append(items, item2.NewKey(item2.KeyTypeRamses, 20))
			items = append(items, item2.NewKey(item2.KeyTypeMenes, 15))
			items = append(items, item2.NewKey(item2.KeyTypePartner, 12))
			items = append(items, item2.NewKey(item2.KeyTypePharaoh, 4))
			lives = 60
		case rls.Pharaoh():
			items = append(items, item2.NewSpecialItem(item2.PartnerPackageType{}, 9))
			items = append(items, item2.NewKey(item2.KeyTypeRamses, 30))
			items = append(items, item2.NewKey(item2.KeyTypeMenes, 20))
			items = append(items, item2.NewKey(item2.KeyTypePartner, 18))
			items = append(items, item2.NewKey(item2.KeyTypePharaoh, 8))
			lives = 75
		case rls.Manager(), rls.Admin(), rls.Owner():
			items = append(items, item2.NewSpecialItem(item2.PartnerPackageType{}, 10))
			items = append(items, item2.NewKey(item2.KeyTypeRamses, 30))
			items = append(items, item2.NewKey(item2.KeyTypeMenes, 20))
			items = append(items, item2.NewKey(item2.KeyTypePartner, 20))
			items = append(items, item2.NewKey(item2.KeyTypePharaoh, 10))
			lives = 85
		}
		for _, i := range items {
			item2.AddOrDrop(h, i)
		}

		u.Teams.Lives += lives

		var itemNames []string
		for _, i := range items {
			itemNames = append(itemNames, text.Colourf("<red>%dx</red> %s", i.Count(), i.CustomName()))
		}
		nm := []rune(r.Name())
		internal.Broadcastf(tx, "user.reclaim", highest.Coloured(p.Name()), r.Coloured(string(append([]rune{unicode.ToUpper(nm[0])}, nm[1:]...))), strings.Join(itemNames, ", "), lives)
	}
	data2.SaveUser(u)
}

// Run ...
func (r ReclaimReset) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	targets := r.Targets.LoadOr(nil)
	if len(targets) > 1 {
		o.Error(lang.Translatef(data2.Language{}, "command.targets.exceed"))
		return
	}
	if len(targets) == 1 {
		target, ok := targets[0].(*player.Player)
		if !ok {
			o.Error(lang.Translatef(data2.Language{}, "command.target.unknown"))
			return
		}

		u, err := data2.LoadUserFromName(target.Name())
		if err != nil {
			o.Error(lang.Translatef(data2.Language{}, "command.target.unknown"))
			return
		}

		u.Teams.Reclaimed = false
		data2.SaveUser(u)
		return
	}

	if p, ok := s.(*player.Player); ok {
		u, err := data2.LoadUserFromName(p.Name())
		if err != nil {
			o.Error(lang.Translatef(data2.Language{}, "command.target.unknown"))
			return
		}

		u.Teams.Reclaimed = false
		data2.SaveUser(u)
		return
	}
}

// Allow ...
func (ReclaimReset) Allow(s cmd.Source) bool {
	return Allow(s, true, rls.Operator())
}

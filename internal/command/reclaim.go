package command

import (
	"strings"
	"unicode"

	rls "github.com/moyai-network/teams/internal/roles"
	"github.com/moyai-network/teams/pkg/lang"

	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/data"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	it "github.com/moyai-network/teams/internal/item"
	"github.com/moyai-network/teams/internal/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Reclaim is a command that allows players to reclaim their partner package.
type Reclaim struct{}

// ReclaimReset is a command that allows admins to reset the reclaim cooldown.
type ReclaimReset struct {
	Sub     cmd.SubCommand             `cmd:"reset"`
	Targets cmd.Optional[[]cmd.Target] `cmd:"target"`
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
			items = append(items, it.NewSpecialItem(it.PartnerPackageType{}, 3))
			items = append(items, it.NewKey(it.KeyTypePartner, 3))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 1))
			items = append(items, it.NewKey(it.KeyTypeMenes, 2))
			items = append(items, it.NewKey(it.KeyTypeRamses, 3))
		case rls.Trial():
			items = append(items, it.NewSpecialItem(it.PartnerPackageType{}, 3))
			items = append(items, it.NewKey(it.KeyTypePartner, 4))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 2))
			items = append(items, it.NewKey(it.KeyTypeMenes, 3))
			items = append(items, it.NewKey(it.KeyTypeRamses, 4))
			lives = 3
		case rls.Khufu(), rls.Media():
			items = append(items, it.NewSpecialItem(it.PartnerPackageType{}, 4))
			items = append(items, it.NewKey(it.KeyTypeRamses, 10))
			items = append(items, it.NewKey(it.KeyTypeMenes, 5))
			items = append(items, it.NewKey(it.KeyTypePartner, 3))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 1))
			lives = 30
		case rls.Ramses(), rls.Famous():
			items = append(items, it.NewSpecialItem(it.PartnerPackageType{}, 5))
			items = append(items, it.NewKey(it.KeyTypeRamses, 15))
			items = append(items, it.NewKey(it.KeyTypeMenes, 10))
			items = append(items, it.NewKey(it.KeyTypePartner, 6))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 2))
			lives = 45
		case rls.Menes(), rls.Mod(), rls.Partner(), rls.Trial():
			items = append(items, it.NewSpecialItem(it.PartnerPackageType{}, 7))
			items = append(items, it.NewKey(it.KeyTypeRamses, 20))
			items = append(items, it.NewKey(it.KeyTypeMenes, 15))
			items = append(items, it.NewKey(it.KeyTypePartner, 12))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 4))
			lives = 60
		case rls.Pharaoh():
			items = append(items, it.NewSpecialItem(it.PartnerPackageType{}, 9))
			items = append(items, it.NewKey(it.KeyTypeRamses, 30))
			items = append(items, it.NewKey(it.KeyTypeMenes, 20))
			items = append(items, it.NewKey(it.KeyTypePartner, 18))
			items = append(items, it.NewKey(it.KeyTypePharaoh, 8))
			lives = 75
		case rls.Manager(), rls.Admin(), rls.Owner():
			items = append(items, it.NewSpecialItem(it.PartnerPackageType{}, 10))
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
		internal.Broadcastf("user.reclaim", highest.Coloured(p.Name()), r.Coloured(string(append([]rune{unicode.ToUpper(nm[0])}, nm[1:]...))), strings.Join(itemNames, ", "), lives)
	}
	data.SaveUser(u)
}

// Run ...
func (r ReclaimReset) Run(s cmd.Source, o *cmd.Output) {
	targets := r.Targets.LoadOr(nil)
	if len(targets) > 1 {
		o.Error(lang.Translatef(data.Language{}, "command.targets.exceed"))
		return
	}
	if len(targets) == 1 {
		target, ok := targets[0].(*player.Player)
		if !ok {
			o.Error(lang.Translatef(data.Language{}, "command.target.unknown"))
			return
		}

		u, err := data.LoadUserFromName(target.Name())
		if err != nil {
			o.Error(lang.Translatef(data.Language{}, "command.target.unknown"))
			return
		}

		u.Teams.Reclaimed = false
		data.SaveUser(u)
		return
	}

	if p, ok := s.(*player.Player); ok {
		u, err := data.LoadUserFromName(p.Name())
		if err != nil {
			o.Error(lang.Translatef(data.Language{}, "command.target.unknown"))
			return
		}

		u.Teams.Reclaimed = false
		data.SaveUser(u)
		return
	}
}

// Allow ...
func (ReclaimReset) Allow(s cmd.Source) bool {
	return Allow(s, true, rls.Operator())
}

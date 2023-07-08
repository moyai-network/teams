package form

import (
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/kit"
	"github.com/moyai-network/teams/moyai/user"

	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/role"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Kits struct{}

func NewKitForm(p *player.Player) form.Menu {
	f := form.NewMenu(Kits{}, "Kits")
	u, _ := data.LoadUser(p.Name())
	for _, k := range kit.All() {
		t := k.Name()
		if !u.Roles.Contains(role.Wraith{}) && k == (kit.Master{}) {
			t = text.Colourf("<red>%s</red>", t)
		}
		cd := u.Kits.Key(t)
		if cd.Active() {
			t += text.Colourf("\n<red>%s</red>", cd.Remaining().Round(time.Second))
		}
		f = f.WithButtons(form.NewButton(t, k.Texture()))
	}
	return f
}

func (k Kits) Submit(s form.Submitter, pressed form.Button) {
	p := s.(*player.Player)
	h, ok := user.Lookup(p.Name())
	if !ok {
		return
	}

	u, _ := data.LoadUser(p.Name())

	if h.Combat().Active() {
		h.Message("command.kit.tagged")
		return
	}
	name := strings.Split(moose.StripMinecraftColour(pressed.Text), "\n")[0]
	cd := u.Kits.Key(name)
	if cd.Active() {
		h.Message("command.kit.cooldown", cd.Remaining().Round(time.Second))
		return
	} else {
		cd.Set(10 * time.Minute)
	}
	switch name {
	case "Archer":
		kit.Apply(kit.Archer{}, p)
	case "Master":
		if !u.Roles.Contains(role.Wraith{}, role.Revenant{}) {
			p.Message(text.Colourf("<red>You must be a Wraith to use this kit.</red>"))
			return
		}
		kit.Apply(kit.Master{}, p)
	case "Bard":
		kit.Apply(kit.Bard{}, p)
	case "Rogue":
		kit.Apply(kit.Rogue{}, p)
	case "Builder":
		kit.Apply(kit.Builder{}, p)
	case "Diamond":
		kit.Apply(kit.Diamond{}, p)
	case "Miner":
		kit.Apply(kit.Miner{}, p)
	case "Stray":
		kit.Apply(kit.Stray{}, p)
	case "Refill":
		kit.Apply(kit.Refill{}, p)
	}
}

package form

/*type Kits struct{}

func NewKitForm(p *player.Player) form.Menu {
	f := form.NewMenu(Kits{}, "Kits")
	u, _ := data.LoadUserOrCreate(p.Name())
	for _, k := range kit2.All() {
		t := k.Name()
		//if !u.Roles.Contains(role.Wraith{}) && k == (kit.Diamond{}) {
		//	t = text.Colourf("<red>%s</red>", t)
		//}
		cd := u.GameMode.Teams.Kits.Key(t)
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

	u, _ := data.LoadUserOrCreate(p.Name())

	if h.Combat().Active() {
		h.Message("command.kit.tagged")
		return
	}

	name := strings.Split(moose.StripMinecraftColour(pressed.Text), "\n")[0]
	cd := u.GameMode.Teams.Kits.Key(name)
	if cd.Active() {
		h.Message("command.kit.cooldown", cd.Remaining().Round(time.Second))
		return
	} else {
		cd.Set(time.Minute)
	}
	switch name {
	case "Archer":
		kit2.Apply(kit2.Archer{}, p)
	case "Master":
	if !u.Roles.Contains(role.Wraith{}, role.Revenant{}) {
		p.Message(text.Colourf("<red>You must be a Wraith to use this kit.</red>"))
		return
	}
	kit.Apply(kit.Master{}, p)
	case "Bard":
		kit2.Apply(kit2.Bard{}, p)
	case "Rogue":
		kit2.Apply(kit2.Rogue{}, p)
	case "Builder":
		kit2.Apply(kit2.Builder{}, p)
	case "Diamond":
		kit2.Apply(kit2.Diamond{}, p)
	case "Miner":
		kit2.Apply(kit2.Miner{}, p)
	case "Stray":
		kit2.Apply(kit2.Stray{}, p)
	case "Refill":
		kit2.Apply(kit2.Refill{}, p)
	}
}
*/

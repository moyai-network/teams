package item

/*func NewPearlBow() item.Stack {
	return item.NewStack(PearlBow{}, 1).WithCustomName(text.Colourf("<gold>Pearl Bow</gold>"))
}

// PearlBow is a ranged weapon that fires arrows.
type PearlBow struct{}

func (PearlBow) Name() string {
	return text.Colourf("<yellow>Pearl Bow</yellow>")
}

func (PearlBow) Item() world.Item {
	return item.Bow{}
}

func (PearlBow) Lore() []string {
	return []string{text.Colourf("<grey>Use to pearl away far distances.</grey>")}
}

func (PearlBow) Key() string {
	return "pearl_bow"
}

// MaxCount always returns 1.
func (PearlBow) MaxCount() int {
	return 1
}

// DurabilityInfo ...
func (PearlBow) DurabilityInfo() item.DurabilityInfo {
	return item.DurabilityInfo{
		MaxDurability: 385,
	}
}

// Release ...
func (PearlBow) Release(releaser item.Releaser, duration time.Duration, ctx *item.UseContext) {
	creative := releaser.GameMode().CreativeInventory()
	ticks := duration.Milliseconds() / 50
	if ticks < 3 {
		// The player must hold the PearlBow for at least three ticks.
		return
	}

	d := float64(ticks) / 20
	force := math.Min((d*d+d*2)/3, 1)
	if force < 0.1 {
		// The force must be at least 0.1.
		return
	}

	arrow, ok := ctx.FirstFunc(func(stack item.Stack) bool {
		_, ok := stack.Item().(item.EnderPearl)
		return ok
	})
	if !ok && !creative {
		// No arrows in inventory and not in creative mode.
		return
	}

	rot := releaser.Rotation()
	rot = cube.Rotation{-rot[0], -rot[1]}
	if rot[0] > 180 {
		rot[0] = 360 - rot[0]
	}

	consume := !creative

	create := releaser.World().EntityRegistry().Config().EnderPearl
	projectile := create(eyePosition(releaser), releaser.Rotation().Vec3().Mul(force*3), releaser)

	ctx.DamageItem(1)
	if consume {
		ctx.Consume(arrow.Grow(-arrow.Count() + 1))
	}

	releaser.PlaySound(sound.BowShoot{})
	releaser.World().AddEntity(projectile)
}


func eyePosition(e world.Entity) mgl64.Vec3 {
	pos := e.Position()
	if eyed, ok := e.(interface{ EyeHeight() float64 }); ok {
		pos = pos.Add(mgl64.Vec3{0, eyed.EyeHeight()})
	}
	return pos
}


// EnchantmentValue ...
func (PearlBow) EnchantmentValue() int {
	return 1
}

// Requirements returns the required items to release this item.
func (PearlBow) Requirements() []item.Stack {
	return []item.Stack{item.NewStack(item.EnderPearl{}, 1)}
}

// EncodeItem ...
func (PearlBow) EncodeItem() (name string, meta int16) {
	return "minecraft:bow", 0
}
*/

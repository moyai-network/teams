package kit

/*// Diamond represents the Diamond kit.
type Diamond struct{}

// Name ...
func (Diamond) Name() string {
	return "Diamond"
}

func (Diamond) Texture() string {
	return "textures/items/diamond_helmet"
}

// Items ...
func (Diamond) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, 1)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 3; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}
	items[2] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
	items[8] = item.NewStack(item.GoldenApple{}, 32)
	items[16] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
	items[17] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
	items[25] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
	items[26] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
	items[34] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
	items[35] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
	return items
}

// Armour ...
func (Diamond) Armour(*player.Player) [4]item.Stack {
	protection := item.NewEnchantment(ench.Protection{}, 1)
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)
	featherFalling := item.NewEnchantment(enchantment.FeatherFalling{}, 4)

	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, unbreaking),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, unbreaking),
		item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, unbreaking),
		item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, unbreaking, featherFalling),
	}
}*/

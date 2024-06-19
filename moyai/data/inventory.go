package data

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"go.mongodb.org/mongo-driver/bson"
)

// Inventory is a struct that contains all data of the player inventories.
type Inventory struct {
	// Items contains all the items in the player's main inventory.
	// This excludes armor and offhand.
	Items []item.Stack
	// Boots, Leggings, Chestplate, Helmet are armor pieces that belong to the slot corresponding to the name.
	Boots      item.Stack
	Leggings   item.Stack
	Chestplate item.Stack
	Helmet     item.Stack
	// OffHand is what the player is carrying in their non-main hand, like a shield or arrows.
	OffHand item.Stack
	// MainHandSlot saves the slot in the hotbar that the player is currently switched to.
	// Should be between 0-8.
	MainHandSlot uint32
}

// Apply applies the inventory to the player's inventory, armor and held items.
func (i Inventory) Apply(p *player.Player) {
	inv, arm := p.Inventory(), p.Armour()
	inv.Clear()
	arm.Clear()

	for in, st := range i.Items {
		if st.Empty() {
			continue
		}
		_ = inv.SetItem(in, st)
	}
	arm.SetHelmet(i.Helmet)
	arm.SetChestplate(i.Chestplate)
	arm.SetLeggings(i.Leggings)
	arm.SetBoots(i.Boots)

	held, _ := p.HeldItems()
	p.SetHeldItems(held, i.OffHand)
}

// MarshalBSON ...
func (i Inventory) MarshalBSON() ([]byte, error) {
	jsonInventoryData := invToData(i)
	return bson.Marshal(jsonInventoryData)
}

// UnmarshalBSON ...
func (i Inventory) UnmarshalBSON(b []byte) error {
	var jsonInventoryData inventoryData
	if err := bson.Unmarshal(b, &jsonInventoryData); err != nil {
		return err
	}
	return dataToInv(jsonInventoryData, &i)
}

func dataToInv(data inventoryData, inv *Inventory) error {
	inv.Items = make([]item.Stack, len(data.Items))
	for i, d := range data.Items {
		stack, err := d.toStack()
		if err != nil {
			return err
		}
		inv.Items[i] = stack
	}
	inv.Boots, _ = data.Boots.toStack()
	inv.Leggings, _ = data.Leggings.toStack()
	inv.Chestplate, _ = data.Chestplate.toStack()
	inv.Helmet, _ = data.Helmet.toStack()
	inv.OffHand, _ = data.OffHand.toStack()
	inv.MainHandSlot = data.MainHandSlot
	return nil
}

func invToData(inv Inventory) inventoryData {
	data := inventoryData{
		Items:        make([]itemData, len(inv.Items)),
		Boots:        stackToData(inv.Boots),
		Leggings:     stackToData(inv.Leggings),
		Chestplate:   stackToData(inv.Chestplate),
		Helmet:       stackToData(inv.Helmet),
		OffHand:      stackToData(inv.OffHand),
		MainHandSlot: inv.MainHandSlot,
	}
	for i, stack := range inv.Items {
		data.Items[i] = stackToData(stack)
	}
	return data
}

func stackToData(stack item.Stack) itemData {
	if stack.Empty() {
		return itemData{}
	}
	name, meta := stack.Item().EncodeItem()
	return itemData{
		Name:         name,
		Meta:         meta,
		Count:        stack.Count(),
		CustomName:   stack.CustomName(),
		Lore:         stack.Lore(),
		Damage:       stack.Durability(),
		AnvilCost:    stack.AnvilCost(),
		Data:         stack.Values(),
		Enchantments: enchantsToData(stack.Enchantments()),
	}
}

func enchantsToData(enchants []item.Enchantment) []enchantmentData {
	data := make([]enchantmentData, len(enchants))
	for i, ench := range enchants {
		data[i] = enchantmentData{
			Name:  ench.Type().Name(),
			Level: ench.Level(),
		}
	}
	return data
}

type inventoryData struct {
	Items        []itemData
	Boots        itemData
	Leggings     itemData
	Chestplate   itemData
	Helmet       itemData
	OffHand      itemData
	MainHandSlot uint32
}

type itemData struct {
	Name  string
	Meta  int16
	Count int

	CustomName   string
	Lore         []string
	Damage       int
	AnvilCost    int
	Data         map[string]any
	Enchantments []enchantmentData
}

func (i itemData) toStack() (item.Stack, error) {
	it, ok := world.ItemByName(i.Name, i.Meta)
	if !ok {
		return item.Stack{}, nil
	}
	stack := item.NewStack(it, i.Count)
	if len(i.CustomName) > 0 {
		stack = stack.WithCustomName(i.CustomName)
	}
	if len(i.Lore) > 0 {
		stack = stack.WithLore(i.Lore...)
	}
	stack = stack.WithDurability(i.Damage)
	stack = stack.WithAnvilCost(i.AnvilCost)
	for key, value := range i.Data {
		stack = stack.WithValue(key, value)
	}

	for _, ench := range i.Enchantments {
		en := ench.toEnchantment()
		if en != nil {
			stack = stack.WithEnchantments(item.NewEnchantment(en, ench.Level))
		}
	}

	return stack, nil
}

type enchantmentData struct {
	Name  string
	Level int
}

func (e enchantmentData) toEnchantment() item.EnchantmentType {
	for _, en := range item.Enchantments() {
		if en.Name() == e.Name {
			return en
		}
	}
	return nil
}

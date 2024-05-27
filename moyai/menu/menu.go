package menu

import (
	_ "unsafe"

	"github.com/bedrock-gophers/nbtconv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/unsafe"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func glassFilledStack(size int) []item.Stack {
	var stacks = make([]item.Stack, size)
	for i := 0; i < size; i++ {
		stacks[i] = item.NewStack(block.StainedGlassPane{Colour: item.ColourPink()}, 1).WithCustomName(text.Colourf("<aqua>Moyai</aqua>")).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 1))
	}
	return stacks
}

func updateInventory(p *player.Player) {
	inv := p.Inventory()
	arm := p.Armour()
	if s := unsafe.Session(p); s != session.Nop {
		for i := 0; i < 36; i++ {
			st, _ := inv.Item(i)
			viewSlotChange(s, i, st, protocol.WindowIDInventory)
		}

		for i, st := range arm.Slots() {
			viewSlotChange(s, i, st, protocol.WindowIDArmour)
		}
	}
}


func viewSlotChange(s *session.Session, slot int, it item.Stack, windowID uint32) {
	unsafe.WritePacket(s, &packet.InventorySlot{
		WindowID: windowID,
		Slot:     uint32(slot),
		NewItem:  instanceFromItem(it),
	})
}

// instanceFromItem converts an item.Stack to its network ItemInstance representation.
func instanceFromItem(it item.Stack) protocol.ItemInstance {
	return protocol.ItemInstance{
		StackNetworkID: item_id(it),
		Stack:          stackFromItem(it),
	}
}

type glint struct{}

func (glint) Name() string {
	return ""
}
func (glint) MaxLevel() int {
	return 1
}
func (glint) Cost(level int) (int, int) {
	return 0, 0
}
func (glint) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}
func (glint) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return true
}
func (glint) CompatibleWithItem(i world.Item) bool {
	return true
}


// noinspection ALL
//
//go:linkname item_id github.com/df-mc/dragonfly/server/item.id
func item_id(s item.Stack) int32

func formatBool(b bool) string {
	if b {
		return "<green>Yes</green>"
	} else {
		return "<red>No</red>"
	}
}

// stackFromItem converts an item.Stack to its network ItemStack representation.
func stackFromItem(it item.Stack) protocol.ItemStack {
	if it.Empty() {
		return protocol.ItemStack{}
	}

	var blockRuntimeID uint32
	if b, ok := it.Item().(world.Block); ok {
		blockRuntimeID = world.BlockRuntimeID(b)
	}

	rid, meta, _ := world.ItemRuntimeID(it.Item())

	return protocol.ItemStack{
		ItemType: protocol.ItemType{
			NetworkID:     rid,
			MetadataValue: uint32(meta),
		},
		HasNetworkID:   true,
		Count:          uint16(it.Count()),
		BlockRuntimeID: int32(blockRuntimeID),
		NBTData:        nbtconv.WriteItem(it, false),
	}
}
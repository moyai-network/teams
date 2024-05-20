package kit

import (
	"github.com/bedrock-gophers/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	_ "unsafe"
)

// Kit contains all the items, armour, and effects obtained by a kit.
type Kit interface {
	// Name is the name of the kit.
	Name() string
	// Texture is the texture of the kit.
	Texture() string
	// Items returns the items provided by the kit.
	Items(*player.Player) (items [36]item.Stack)
	// Armour contains the armour applied by using the kit.
	// The item stacks are ordered helmet, chestplate, leggings, and then boots.
	Armour(*player.Player) [4]item.Stack
}

func All() []Kit {
	return []Kit{
		Miner{},
		Builder{},
		Archer{},
		Bard{},
		Stray{},
		Rogue{},
		Master{},
	}
}

// Apply ...
func Apply(kit Kit, p *player.Player) {
	inv := p.Inventory()
	armour := kit.Armour(p)
	for slot, itm := range kit.Items(p) {
		if itm.Empty() {
			continue
		}
		itm = ench.AddEnchantmentLore(itm)
		if inv.Slots()[slot].Item() != nil {
			it.Drop(p, itm)
		} else {
			_ = inv.SetItem(slot, itm)
		}
	}
	arm := p.Armour()
	for slot, itm := range armour {
		if itm.Empty() {
			continue
		}
		itm = ench.AddEnchantmentLore(itm)
		if arm.Slots()[slot].Item() != nil {
			it.Drop(p, itm)
		} else {
			switch slot {
			case 0:
				arm.SetHelmet(itm)
				arm.Inventory().Handler().HandlePlace(nil, 0, itm)
			case 1:
				arm.SetChestplate(itm)
				arm.Inventory().Handler().HandlePlace(nil, 1, itm)
			case 2:
				arm.SetLeggings(itm)
				arm.Inventory().Handler().HandlePlace(nil, 2, itm)
			case 3:
				arm.SetBoots(itm)
				arm.Inventory().Handler().HandlePlace(nil, 3, itm)
			}
		}
	}
	if s := player_session(p); s != session.Nop {
		for i := 0; i < 36; i++ {
			st, _ := inv.Item(i)
			viewSlotChange(s, i, st, protocol.WindowIDInventory)
		}

		for i, st := range arm.Slots() {
			viewSlotChange(s, i, st, protocol.WindowIDArmour)
		}
	}
}

// viewSlotChange ...
func viewSlotChange(s *session.Session, slot int, it item.Stack, windowID uint32) {
	session_writePacket(s, &packet.InventorySlot{
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

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session

// noinspection ALL
//
//go:linkname item_id github.com/df-mc/dragonfly/server/item.id
func item_id(s item.Stack) int32

// noinspection ALL
//
//go:linkname session_writePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
func session_writePacket(*session.Session, packet.Packet)

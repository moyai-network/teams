package minecraft

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/moyai"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func placeText(c moyai.Config) {
	w := moyai.Overworld()
	spawn := mgl64.Vec3{0.5, 68, 5}
	crate := mgl64.Vec3{5.5, 68, 15.5}
	kit := mgl64.Vec3{-7, 67, 38.5}
	shop := mgl64.Vec3{13, 62, 60}
	adv := mgl64.Vec3{0, 62, 81}
	for _, e := range moyai.Overworld().Entities() {
		if _, ok := e.Type().(entity.TextType); ok {
			e.Close()
		}
	}
	for _, e := range moyai.Overworld().Entities() {
		if _, ok := e.Type().(entity.TextType); ok {
			e.Close()
		}
	}
	count := 1
	for i := 0; i < count; i++ {
		for _, e := range []*entity.Ent{
			entity.NewText(text.Colourf("<b><red>MoyaiHCF</red></b>"), mgl64.Vec3{spawn.X(), spawn.Y() + 2.5, spawn.Z()}),
			entity.NewText(text.Colourf("<grey>Season %v began on %v.</grey>", c.Moyai.Season, c.Moyai.Start), mgl64.Vec3{spawn.X(), spawn.Y() + 2, spawn.Z()}),
			entity.NewText(text.Colourf("<grey>It will conclude on %v.</grey>", c.Moyai.End), mgl64.Vec3{spawn.X(), spawn.Y() + 1.6, spawn.Z()}),
			entity.NewText(text.Colourf("<red>Store</red>: https://store.moyai.club"), mgl64.Vec3{spawn.X(), spawn.Y() + 1.1, spawn.Z()}),
			entity.NewText(text.Colourf("<red>Discord</red>: discord.moyai.club"), mgl64.Vec3{spawn.X(), spawn.Y() + 0.5, spawn.Z()}),
			entity.NewText(text.Colourf("<grey>moyai.club</grey>"), spawn),
		} {
			w.AddEntity(e)
		}

		for _, e := range []*entity.Ent{
			entity.NewText(text.Colourf("<b><red>Crates & Keys</red></b>"), mgl64.Vec3{crate.X(), crate.Y() + 2.5, crate.Z()}),
			entity.NewText(text.Colourf("<yellow>Crate Keys can be obtained easily with /reclaim (once).</yellow>"), mgl64.Vec3{crate.X(), crate.Y() + 2, crate.Z()}),
			entity.NewText(text.Colourf("<grey>To the left is the Partner Crate and to the right is the KOTH Crate.</grey>"), mgl64.Vec3{crate.X(), crate.Y() + 1.6, crate.Z()}),
			entity.NewText(text.Colourf("<grey>In front are the three donor crates.</grey>"), mgl64.Vec3{crate.X(), crate.Y() + 1.1, crate.Z()}),
			entity.NewText(text.Colourf("<red>Buy Keys</red>: store.moyai.club"), mgl64.Vec3{crate.X(), crate.Y() + 0.5, crate.Z()}),
		} {
			w.AddEntity(e)
		}

		for _, e := range []*entity.Ent{
			entity.NewText(text.Colourf("<b><red>Kits</red></b>"), mgl64.Vec3{kit.X(), kit.Y() + 2.5, kit.Z()}),
			entity.NewText(text.Colourf("<yellow>Kits can be equipped using /kit</yellow>"), mgl64.Vec3{kit.X(), kit.Y() + 2, kit.Z()}),
			entity.NewText(text.Colourf("<grey>All kits have a 4-hour cooldown (2-hours for donors).</grey>"), mgl64.Vec3{kit.X(), kit.Y() + 1.6, kit.Z()}),
			entity.NewText(text.Colourf("<grey>Certain kits must be purchased through the store.</grey>"), mgl64.Vec3{kit.X(), kit.Y() + 1.1, kit.Z()}),
			entity.NewText(text.Colourf("<red>Buy Kits</red>: store.moyai.club"), mgl64.Vec3{kit.X(), kit.Y() + 0.5, kit.Z()}),
		} {
			w.AddEntity(e)
		}

		for _, e := range []*entity.Ent{
			entity.NewText(text.Colourf("<b><red>Shop</red></b>"), mgl64.Vec3{shop.X(), shop.Y() + 2.5, shop.Z()}),
			entity.NewText(text.Colourf("<yellow>You can buy and sell items at the shop for credits!</yellow>"), mgl64.Vec3{shop.X(), shop.Y() + 2, shop.Z()}),
			entity.NewText(text.Colourf("<grey>You can deposit credits into your team with /f deposit.</grey>"), mgl64.Vec3{shop.X(), shop.Y() + 1.6, shop.Z()}),
			entity.NewText(text.Colourf("<grey>Credits can also be used in certain hours during the Black Market...</grey>"), mgl64.Vec3{shop.X(), shop.Y() + 1.1, shop.Z()}),
			entity.NewText(text.Colourf("<red>Buy Credits</red>: store.moyai.club"), mgl64.Vec3{shop.X(), shop.Y() + 0.5, shop.Z()}),
		} {
			w.AddEntity(e)
		}

		for _, e := range []*entity.Ent{
			entity.NewText(text.Colourf("<b><red>Chasing Glory!</red></b>"), mgl64.Vec3{adv.X(), adv.Y() + 2.5, adv.Z()}),
			entity.NewText(text.Colourf("<yellow>Create a team via /t create and a claim via /t claim to get started!</yellow>"), mgl64.Vec3{adv.X(), adv.Y() + 2, adv.Z()}),
			entity.NewText(text.Colourf("<grey>Score team points by killing users, capturing KOTHs, and other events.</grey>"), mgl64.Vec3{adv.X(), adv.Y() + 1.6, adv.Z()}),
			entity.NewText(text.Colourf("<grey>The top three factiosn will receive prizes every map. March down South Road and get started!</grey>"), mgl64.Vec3{adv.X(), adv.Y() + 1.1, adv.Z()}),
			entity.NewText(text.Colourf("<red>Buy Ranks</red>: store.moyai.club"), mgl64.Vec3{adv.X(), adv.Y() + 0.5, adv.Z()}),
		} {
			w.AddEntity(e)
		}
	}

}

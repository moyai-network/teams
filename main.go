package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unsafe"

	"github.com/bedrock-gophers/intercept"
	"github.com/bedrock-gophers/spawner/spawner"

	"github.com/bedrock-gophers/console/console"
	"github.com/bedrock-gophers/inv/inv"
	"github.com/bedrock-gophers/tebex/tebex"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/df-mc/npc"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	cr "github.com/moyai-network/teams/moyai/crate"
	"github.com/moyai-network/teams/moyai/data"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	ent "github.com/moyai-network/teams/moyai/entity"
	"github.com/moyai-network/teams/moyai/menu"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft"

	"github.com/go-gl/mathgl/mgl64"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/sotw"

	//proxypacket "github.com/paroxity/portal/socket/packet"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/command"
	"github.com/moyai-network/teams/moyai/user"

	"github.com/moyai-network/teams/moyai/crate"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

var (
	shopSigns   = []shopSign{}
	cowSpawners = []cube.Pos{
		{-17, 63, -14},
		{-15, 63, -18},
		{-13, 63, -22},
		{-11, 63, -26},
	}
	endermanSpawners = []cube.Pos{
		{28, 63, -33},
		{27, 63, -36},
		{26, 63, -39},
		{25, 63, -42},
	}
)

func main() {
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	lang.Register(language.English)
	lang.Register(language.Spanish)
	lang.Register(language.French)

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	console.Enable(log)
	conf, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	config := configure(conf, log)
	/*ac := oomph.New(oomph.OomphSettings{
		LocalAddress:  ":19133",
		RemoteAddress: ":19132",
		RequirePacks:  true,
	})

	ac.Listen(&config, text.Colourf("<dark-red>balls</dark-red>"), []minecraft.Protocol{}, true, false)
	go func() {
		for {
			_, err := ac.Accept()
			if err != nil {
				return
			}

			log.Println("LOL BRO I CONNECTED VIA OOMPGH")

		}
	}()*/
	pk := intercept.NewPacketListener()
	pk.Listen(&config, ":19132", []minecraft.Protocol{})

	go func() {
		for {
			p, err := pk.Accept()
			if err != nil {
				return
			}
			p.Handle(user.NewPacketHandler(p))
		}
	}()

	srv := moyai.NewServer(config)
	handleServerClose(srv)

	nether, end := moyai.ConfigureDimensions(config.Entities)
	w := srv.World()
	configureWorld(w)
	configureWorld(nether)
	configureWorld(end)

	placeSpawners(w)
	clearEntities(w)
	placeText(w, conf)
	placeSlapper(w)
	placeCrates(w)
	placeShopSigns(w)

	go tickBlackMarket(srv)
	go tickClearLag(srv)

	registerCommands(srv)
	srv.Listen()

	store := loadStore(conf.Moyai.Tebex, log)

	for srv.Accept(acceptFunc(store, conf.Proxy.Enabled, srv)) {
		// Do nothing.
	}
}

func tickClearLag(srv *server.Server) {
	t := time.NewTicker(time.Minute / 2)
	defer t.Stop()

	for range t.C {
		for _, e := range srv.World().Entities() {
			if et, ok := e.(*entity.Ent); ok && et.Type() == (entity.ItemType{}) {
				age := fetchPrivateField[time.Duration](et, "age")
				if age > (time.Minute*5)/2 {
					srv.World().RemoveEntity(e)
				}
			}
		}
	}

}

// fetchPrivateField fetches a private field of a session.
func fetchPrivateField[T any](v any, name string) T {
	reflectedValue := reflect.ValueOf(v).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)
	privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

	return privateFieldValue.Interface().(T)
}

func tickBlackMarket(srv *server.Server) {
	t := time.NewTicker(time.Minute * 15)
	defer t.Stop()

	for range t.C {
		if time.Since(moyai.LastBlackMarket()) < time.Hour {
			continue
		}

		if rand.Intn(4) == 0 {
			moyai.SetLastBlackMarket(time.Now())
			for _, p := range srv.Players() {
				p.PlaySound(sound.BarrelOpen{})
				p.PlaySound(sound.FireworkHugeBlast{})
				p.PlaySound(sound.FireworkLaunch{})
				p.PlaySound(sound.Note{})
				user.Broadcastf("blackmarket.opened")
			}
		}
	}
}

func placeText(w *world.World, c moyai.Config) {
	spawn := mgl64.Vec3{0.5, 68, 5}
	crate := mgl64.Vec3{5.5, 68, 15.5}
	kit := mgl64.Vec3{-7.5, 67, 39.5}
	shop := mgl64.Vec3{13, 62, 60}
	adv := mgl64.Vec3{0, 62, 81}
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

func placeSlapper(w *world.World) {
	_ = npc.Create(npc.Settings{
		Name:       text.Colourf("<green>Click to use kits</green>"),
		Skin:       skin.Skin{},
		Scale:      1,
		Yaw:        215,
		MainHand:   item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, 1)),
		Helmet:     item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(ench.Protection{}, 1)),
		Chestplate: item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(ench.Protection{}, 1)),
		Leggings:   item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(ench.Protection{}, 1)),
		Boots:      item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(ench.Protection{}, 1)),

		Position: mgl64.Vec3{-7, 65, 38.5},
	}, w, func(p *player.Player) {
		if men, ok := menu.NewKitsMenu(p); ok {
			inv.SendMenu(p, men)
		}
	})
}

func placeSpawners(w *world.World) {
	for _, pos := range cowSpawners {
		sp := spawner.New(ent.NewCow, pos.Vec3Centre(), w, time.Second*30, 25, true)
		w.SetBlock(pos, sp, nil)
	}
	for _, pos := range endermanSpawners {
		sp := spawner.New(ent.NewEnderman, pos.Vec3Centre(), w, time.Second*5, 25, true)
		w.SetBlock(pos, sp, nil)
	}
}

func clearEntities(w *world.World) {
	for _, e := range w.Entities() {
		if _, ok := e.Type().(entity.TextType); ok {
			w.RemoveEntity(e)
		}

		if _, ok := e.Type().(entity.ItemType); ok {
			w.RemoveEntity(e)
		}
	}
}

func placeCrates(w *world.World) {
	for _, c := range crate.All() {
		b := block.NewChest()
		b.Facing = c.Facing().Direction()
		b.CustomName = text.Colourf("%s <grey>Crate</grey>", c.Name())

		*b.Inventory() = *inventory.New(27, nil)

		var items [27]item.Stack
		for i, r := range c.Rewards() {
			if r.Stack().Empty() {
				continue // Ignore this, ill fix it later
			}
			st := ench.AddEnchantmentLore(r.Stack())
			st = st.WithLore(append(st.Lore(), text.Colourf("<yellow>Chance: %d%%</yellow>", r.Chance()))...)
			items[i] = st
		}
		for i, s := range items {
			if s.Empty() {
				items[i] = item.NewStack(block.StainedGlass{Colour: item.ColourRed()}, 1)
			}
		}

		for s, i := range items {
			_ = b.Inventory().SetItem(s, i)
		}

		b.Inventory().Handle(cr.Handler{})

		w.SetBlock(cube.PosFromVec3(c.Position()), b, nil)
		t := entity.NewText(text.Colourf("%s <grey>Crate</grey>\n<yellow>Left click to open crate</yellow>\n<grey>Right click to see rewards</grey>", c.Name()), c.Position().Add(mgl64.Vec3{0, 2, 0}))
		w.AddEntity(t)
	}
}

func configure(conf moyai.Config, log *logrus.Logger) server.Config {
	c, err := conf.Config(log)
	if err != nil {
		panic(err)
	}
	c.Entities = ent.Registry

	c.Name = text.Colourf("<bold><redstone>MOYAI</redstone></bold>") + "§8"
	c.Allower = moyai.NewAllower(conf.Moyai.Whitelisted)
	return c
}

func configureWorld(w *world.World) {
	w.SetDifficulty(world.DifficultyHard)
	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeSurvival)
	w.SetTime(6000)
	w.StopTime()
	w.SetTickRange(0)
	w.StopThundering()
	w.StopRaining()
	w.SetSpawn(cube.Pos{0, 80, 0})

	l := world.NewLoader(8, w, world.NopViewer{})
	l.Move(w.Spawn().Vec3Middle())
	l.Load(math.MaxInt)

	inv.PlaceFakeContainer(w, cube.Pos{0, 255, 0})
}

func recoverFunc(srv *server.Server) func() {
	return func() {
		time.Sleep(time.Millisecond * 500)
		data.FlushCache()

		sotw.Save()
		_ = srv.Close()
		os.Exit(1)
	}
}

func acceptFunc(store *tebex.Client, proxy bool, srv *server.Server) func(*player.Player) {
	return func(p *player.Player) {
		inv.RedirectPlayerPackets(p, recoverFunc(srv))
		store.ExecuteCommands(p)

		u, _ := data.LoadUserOrCreate(p.Name(), p.XUID())
		if !u.Roles.Contains(role.Default{}) {
			u.Roles.Add(role.Default{})
		}
		u.Roles.Add(role.Donor1{})

		p.Message(lang.Translatef(u.Language, "discord.message"))
		if proxy {
			//info := moyai.SearchInfo(p.UUID())
			//p.Handle(user.NewHandler(p, info.XUID))
		} else {
			p.Handle(user.NewHandler(p, p.XUID()))
			p.Armour().Handle(user.NewArmourHandler(p))
		}
		p.RemoveScoreboard()
		for _, ef := range p.Effects() {
			p.RemoveEffect(ef.Type())
		}
		p.ShowCoordinates()
		p.SetFood(20)

		data.SaveUser(u)
		inv := p.Inventory()
		for slot, i := range inv.Slots() {
			for _, sp := range it.SpecialItems() {
				if _, ok := i.Value(sp.Key()); ok {
					_ = inv.SetItem(slot, it.NewSpecialItem(sp, i.Count()))
				}
			}
		}

		// This is here in case we have any chunk
		// glitches with the proxy
		p.SetImmobile()
		p.SetAttackImmunity(time.Millisecond * 500)
		time.AfterFunc(time.Millisecond*500, func() {
			if p != nil {
				p.SetMobile()
			}
		})

		w := p.World()
		for _, e := range w.Entities() {
			if !area.Spawn(w).Vec3WithinOrEqualFloorXZ(e.Position()) {
				continue
			}
			if _, ok := e.(*player.Player); ok {
				continue
			}
			if e.Type() == (entity.TextType{}) {
				continue
			}

			p.HideEntity(e)
		}
		// u.Roles.Add(role.Pharaoh{})
	}

}

// shopSign is a sign that can be placed in the world to create a shop. It can be used to buy or sell items.
type shopSign struct {
	buy      bool
	it       world.Item
	quantity int
	price    int
	pos      cube.Pos
}

func placeShopSigns(w *world.World) {
	for _, s := range shopSigns {
		var txt string
		state := text.Colourf("<green>[Buy]</green>")
		if !s.buy {
			state = text.Colourf("<red>[Sell]</red>")
		}

		name, _ := s.it.EncodeItem()
		txt = fmt.Sprintf("%s\n%s\n%d\n$%d", state, formatItemName(name), s.quantity, s.price)

		b := block.Sign{Front: block.SignText{
			Text: txt,
		}}
		w.SetBlock(s.pos, b, nil)
	}
}

func formatItemName(s string) string {
	split := strings.Split(s, "_")
	for i, str := range split {
		upperCasesPrefix := unicode.ToUpper(rune(str[0]))
		split[i] = string(upperCasesPrefix) + str[1:]
	}
	return strings.Join(split, " ")
}

// loadStore initializes the Tebex store connection.
func loadStore(key string, log *logrus.Logger) *tebex.Client {
	store := tebex.NewClient(log, time.Second*5, key)
	name, domain, err := store.Information()
	if err != nil {
		log.Fatalf("tebex: %v", err)
		return nil
	}
	log.Infof("Connected to Tebex under %v (%v).", name, domain)
	return store
}

func registerCommands(srv *server.Server) {
	for _, c := range []cmd.Command{
		cmd.New("staff", text.Colourf("Staff management commands."), nil, command.StaffMode{}),
		cmd.New("rename", text.Colourf("Rename your items."), nil, command.Rename{}),
		cmd.New("stop", text.Colourf("Stop the server."), nil, command.NewStop(srv)),
		cmd.New("pots", text.Colourf("Place potion chests."), nil, command.Pots{}),
		cmd.New("fix", text.Colourf("Fix your inventory."), nil, command.Fix{}, command.FixAll{}),
		cmd.New("chat", text.Colourf("Chat management commands."), nil, command.ChatMute{}, command.ChatUnMute{}, command.ChatCoolDown{}),
		cmd.New("t", text.Colourf("The main team management command."), []string{"f"},
			command.TeamCreate{},
			command.TeamRename{},
			command.TeamInformation{},
			command.TeamDisband{},
			command.TeamInvite{},
			command.TeamJoin{},
			command.TeamWho{},
			command.TeamLeave{},
			command.TeamKick{},
			command.TeamPromote{},
			command.TeamDemote{},
			command.TeamTop{},
			command.TeamClaim{},
			command.TeamUnClaim{},
			command.TeamSetHome{},
			command.TeamHome{},
			command.TeamList{},
			command.TeamUnFocus{},
			command.TeamFocusPlayer{},
			command.TeamFocusTeam{},
			command.TeamChat{},
			command.TeamWithdraw{},
			command.TeamDeposit{},
			command.TeamWithdrawAll{},
			command.TeamDepositAll{},
			command.TeamStuck{},
			command.TeamRally{},
			command.TeamUnRally{},
			command.TeamMap{},
			command.TeamSetDTR{},
			command.TeamDelete{},
		), cmd.New("whitelist", text.Colourf("Whitelist commands."), []string{"wl"}, command.WhiteListAdd{}, command.WhiteListRemove{}),
		cmd.New("tl", text.Colourf("Send your location to teammates"), nil, command.TL{}),
		cmd.New("balance", text.Colourf("Manage your balance."), []string{"bal"}, command.Balance{}, command.BalancePayOnline{}, command.BalancePayOffline{}, command.BalanceAdd{}, command.BalanceAddOffline{}),
		cmd.New("colour", text.Colourf("Customize the colour of your archer."), nil, command.Colour{}),
		cmd.New("clear", text.Colourf("Clear your Inventory."), nil, command.Clear{}),
		cmd.New("clearlag", text.Colourf("Clears all ground entitys."), nil, command.ClearLag{}),
		cmd.New("logout", text.Colourf("Safely logout of PVP."), nil, command.Logout{}),
		cmd.New("pvp", text.Colourf("Manage PVP timer."), nil, command.PvpEnable{}),
		cmd.New("role", text.Colourf("Manage user roles."), nil, command.RoleAdd{}, command.RoleRemove{}, command.RoleAddOffline{}, command.RoleRemoveOffline{}, command.RoleList{}),
		cmd.New("teleport", text.Colourf("Teleport yourself or another player to a position."), []string{"tp"}, command.TeleportToPos{}, command.TeleportTargetsToPos{}, command.TeleportTargetsToTarget{}, command.TeleportToTarget{}),
		cmd.New("reclaim", text.Colourf("Manage a reclaim."), nil, command.Reclaim{}, command.ReclaimReset{}),
		cmd.New("kit", text.Colourf("Choose a kit."), nil, command.Kit{}, command.KitReset{}),
		cmd.New("ban", text.Colourf("Unleash the ban hammer."), nil, command.Ban{}, command.BanOffline{}, command.BanList{}, command.BanLiftOffline{} /*command.BanInfoOffline{},*/, command.BanForm{}),
		//cmd.New("blacklist", text.Colourf("Blacklist little evaders."), nil, command.Blacklist{}, command.BlacklistOffline{}, command.BlacklistList{}, command.BlacklistLiftOffline{}, command.BlacklistInfoOffline{}, command.BlacklistForm{}),
		cmd.New("kick", text.Colourf("Kick a user."), nil, command.Kick{}),
		cmd.New("mute", text.Colourf("Mute a user."), nil, command.MuteList{}, command.MuteInfo{}, command.MuteInfoOffline{}, command.MuteLift{}, command.MuteLiftOffline{}, command.MuteForm{}, command.Mute{}, command.MuteOffline{}),
		cmd.New("whisper", text.Colourf("Send a private message to a player."), []string{"w", "tell", "msg"}, command.Whisper{}),
		cmd.New("reply", text.Colourf("Reply to the last whispered player."), []string{"r"}, command.Reply{}),
		cmd.New("fly", text.Colourf("Toggle flight."), nil, command.Fly{}),
		cmd.New("sotw", text.Colourf("SOTW management commands."), nil, command.SOTWStart{}, command.SOTWEnd{}, command.SOTWDisable{}),
		cmd.New("freeze", text.Colourf("Freeze possible cheaters."), nil, command.Freeze{}),
		cmd.New("gamemode", text.Colourf("Manage gamemodes."), []string{"gm"}, command.GameMode{}),
		cmd.New("key", text.Colourf("Manage keys"), nil, command.Key{}, command.KeyAll{}),
		cmd.New("koth", text.Colourf("Manage KOTHs."), nil, command.KothStart{}, command.KothStop{}, command.KothList{}),
		cmd.New("pp", text.Colourf("Manage partner packages."), nil, command.PartnerPackageAll{}, command.PartnerPackage{}),
		cmd.New("ping", text.Colourf("Check your ping."), nil, command.Ping{}),
		//cmd.New("data", text.Colourf("Clear data."), nil, command.DataReset{}),
		cmd.New("nick", text.Colourf("Change your nickname."), nil, command.Nick{}, command.NickReset{}),
		cmd.New("vanish", text.Colourf("Vanish as staff."), []string{"v"}, command.Vanish{}),
		cmd.New("lang", text.Colourf("Change your language."), nil, lang.Lang{}),
		cmd.New("blockshop", text.Colourf("Access the blockshop to buy items."), nil, command.BlockShop{}),
		cmd.New("enderchest", text.Colourf("Access your enderchest."), []string{"ec"}, command.Enderchest{}),
		cmd.New("blackmarket", text.Colourf("Access the secret items of the black market"), nil, command.BlackMarket{}),
		cmd.New("trim", text.Colourf("Add trims to your armor"), nil, command.Trim{}, command.TrimClear{}),
		cmd.New("end", text.Colourf("End your adventure."), nil, command.End{}),
		cmd.New("nether", text.Colourf("End your adventure."), nil, command.Nether{}),
		cmd.New("settings", text.Colourf("Access your settings."), nil, command.Settings{}),
		cmd.New("tag", text.Colourf("Manage your tags."), nil, command.TagAddOnline{}, command.TagAddOffline{}, command.TagRemoveOnline{}, command.TagRemoveOffline{}, command.TagSet{}),
		cmd.New("cape", text.Colourf("Manage your capes."), nil, command.Cape{}),
	} {
		cmd.Register(c)
	}

	//cmd.Register(cmd.New("hub", text.Colourf("Return to the Moyai Hub."), []string{"lobby"}, command.Hub{}))
}

// pix converts an image.Image into a []uint8.
func pix(i image.Image) (p []uint8) {
	for y := 0; y <= i.Bounds().Max.Y-1; y++ {
		for x := 0; x <= i.Bounds().Max.X-1; x++ {
			color := i.At(x, y)
			r, g, b, a := color.RGBA()
			p = append(p, uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}
	return
}

func handleServerClose(srv *server.Server) {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		for _, p := range srv.Players() {
			p.Message(text.Colourf("<green>Travelling to <black>The</black> <gold>Hub</gold>...</green>"))
			/*sock, ok := moyai.Socket()
			if ok {
				go func() {
					_ = sock.WritePacket(&proxypacket.TransferRequest{
						PlayerUUID: p.UUID(),
						Server:     "syn.lobby",
					})
				}()
			}*/
		}
		time.Sleep(time.Millisecond * 500)
		data.FlushCache()

		sotw.Save()
		if err := srv.Close(); err != nil {
			logrus.Fatalln("close server: %v", err)
		}
	}()
}

func readConfig() (moyai.Config, error) {
	c := moyai.DefaultConfig()
	g := gophig.NewGophig("./config", "toml", 0777)

	err := g.GetConf(&c)
	if os.IsNotExist(err) {
		err = g.SetConf(c)
	}
	return c, err
}

package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moyai-network/carrot/tebex"
	"github.com/moyai-network/teams/moyai/koth"

	"github.com/moyai-network/moose/worlds"
	"github.com/oomph-ac/oomph"
	"github.com/oomph-ac/oomph/utils"

	"github.com/bedrock-gophers/packethandler"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose/data"
	ent "github.com/moyai-network/teams/moyai/entity"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/kit"
	"github.com/moyai-network/teams/moyai/sotw"

	proxypacket "github.com/paroxity/portal/socket/packet"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/item/recipe"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/playerdb"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/df-mc/npc"
	_ "github.com/moyai-network/moose/console"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/command"
	"github.com/moyai-network/teams/moyai/user"

	"github.com/moyai-network/moose/crate"
	cr "github.com/moyai-network/teams/moyai/crate"
	ench "github.com/moyai-network/teams/moyai/enchantment"

	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

var playerProvider *playerdb.Provider
var worldProvider *mcdb.DB

func main() {
	data.LoadTeams()

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	lang.Register(language.English)
	lang.Register(language.Spanish)
	//lang.Register(language.EuropeanSpanish)

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if config.Proxy.Enabled {
		moyai.NewProxySocket()
	}

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	c, err := config.Config(log)
	if err != nil {
		//panic(err)
	}

	c.Entities = ent.Registry

	pProv, err := playerdb.NewProvider("assets/players")
	if err != nil {
		panic(err)
	}
	playerProvider = pProv
	c.PlayerProvider = playerProvider

	wProv, err := mcdb.Config{
		Log: log,
	}.Open("assets/hcfworld")
	if err != nil {
		panic(err)
	}
	worldProvider = wProv
	c.WorldProvider = worldProvider

	c.Name = text.Colourf("<bold><redstone>MOYAI</redstone></bold>") + "§8"
	c.Generator = func(dim world.Dimension) world.Generator { return nil }
	c.Allower = moyai.NewAllower(config.Moyai.Whitelisted)

	if config.Oomph.Enabled {
		o := oomph.New(log, ":19132")
		o.Listen(&c, text.Colourf("<red>Moyai</red>"), []minecraft.Protocol{}, true, config.Proxy.Enabled)
		go func() {
			for {
				p, err := o.Accept()
				if err != nil {
					return
				}

				p.SetMovementMode(utils.ModeSemiAuthoritative)
				p.SetCombatMode(utils.ModeSemiAuthoritative)

				p.Handle(user.NewOomphHandler(p))
			}
		}()
	} else {
		pk := packethandler.NewPacketListener()
		pk.Listen(&c, ":19136", []minecraft.Protocol{})
		go func() {
			for {
				p, err := pk.Accept()
				if err != nil {
					return
				}
				p.Handle(user.NewPacketHandler(p))
			}
		}()
	}

	srv := c.New()

	clg := time.NewTicker(time.Minute * 5)
	go func() {
		for range clg.C {
			var itemCount int64

			for _, e := range srv.World().Entities() {
				if _, ok := e.Type().(entity.ItemType); ok {
					err := e.Close()
					if err != nil {
						continue
					}
					srv.World().RemoveEntity(e)
					itemCount++
				}
			}
		}
	}()

	t := time.NewTicker(time.Minute * 1)
	go func() {
		for range t.C {
			for _, u := range user.All() {
				_ = playerProvider.Save(u.Player().UUID(), u.Player().Data())
				_ = worldProvider.SavePlayerSpawnPosition(u.Player().UUID(), cube.PosFromVec3(u.Player().Position()))
				d, ok := data.LoadUser(u.Player().Name())
				if ok {
					_ = data.SaveUser(d)
				}
			}
			log.Println("Saving all users/world data.")
			for _, t := range data.Teams() {
				data.SaveTeam(t)
			}
			log.Println("Saving all team/data.")
		}
	}()

	handleServerClose(srv)

	w := srv.World()
	w.Handle(&worlds.Handler{})
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

	// LOL
	for x := -50; x < -30; x++ {
		for z := 190; z < 210; z++ {
			b := w.Block(cube.Pos{x, 28, z})
			if _, ok := b.(block.Hopper); ok {
				w.SetBlock(cube.Pos{x, 28, z}, block.Air{}, nil)
			}
		}
	}

	// Doing this twice so that w.Entities isn't empty, even when it shouldn't be
	for _, c := range crate.All() {
		b := block.NewChest()
		b.CustomName = text.Colourf("%s <grey>Crate</grey>", c.Name())

		*b.Inventory() = *inventory.New(27, nil)

		var items [27]item.Stack
		for i, r := range c.Rewards() {
			if r == nil {
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

	}

	for _, e := range w.Entities() {
		if _, ok := e.Type().(entity.TextType); ok {
			w.RemoveEntity(e)
		}

		if _, ok := e.Type().(entity.ItemType); ok {
			w.RemoveEntity(e)
		}
	}

	for _, c := range crate.All() {
		t := entity.NewText(text.Colourf("%s <grey>Crate</grey>\n<yellow>Right click to open crate</yellow>\n<grey>Left click to see rewards</grey>", c.Name()), c.Position().Add(mgl64.Vec3{0, 2, 0}))
		w.AddEntity(t)
	}

	koth.Broadcast = user.Broadcast

	registerCommands(srv)
	registerRecipes()
	//registerSlappers(w)

	srv.Listen()

	store := loadStore(config.Moyai.Tebex, log)

	for srv.Accept(acceptFunc(store, config.Proxy.Enabled)) {
		// Do nothing.
	}

}

func acceptFunc(store *tebex.Client, proxy bool) func(*player.Player) {
	return func(p *player.Player) {
		store.ExecuteCommands(p)
		if proxy {
			info := moyai.SearchInfo(p.UUID())
			p.Handle(user.NewHandler(p, info.XUID))
		} else {
			p.Handle(user.NewHandler(p, p.XUID()))
		}
		p.SetGameMode(world.GameModeSurvival)
		p.RemoveScoreboard()
		for _, ef := range p.Effects() {
			p.RemoveEffect(ef.Type())
		}
		p.ShowCoordinates()
		p.SetFood(20)

		p.Message(text.Colourf("<green>Make sure to join our discord</green><grey>:</grey> <yellow>discord.gg/moyai</yellow>"))
		inv := p.Inventory()
		for slot, i := range inv.Slots() {
			for _, sp := range it.SpecialItems() {
				if _, ok := i.Value(sp.Key()); ok {
					_ = inv.SetItem(slot, it.NewSpecialItem(sp, i.Count()))
				}
			}
		}

		u, _ := data.LoadUserOrCreate(p.Name())
		if u.Language() == (language.Tag{}) {
			p.SendForm(lang.NewPromptForm())
		}
		// u.Roles.Add(role.Pharaoh{})
	}

} // loadStore initializes the Tebex store connection.
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
		cmd.New("t", text.Colourf("<aqua>The main team management command.</aqua>"), []string{"f"},
			command.TeamCreate{},
			command.NewTeamInformation(srv),
			command.TeamDisband{},
			command.TeamInvite{},
			command.TeamJoin{},
			command.NewTeamWho(srv),
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
		), cmd.New("whitelist", text.Colourf("<aqua>Whitelist commands.</aqua>"), []string{"wl"}, command.WhiteListAdd{}, command.WhiteListRemove{}),
		cmd.New("tl", text.Colourf("<aqua>Send your location to teammates</aqua>"), nil, command.TL{}),
		cmd.New("balance", text.Colourf("<aqua>Manage your balance.</aqua>"), []string{"bal"}, command.Balance{}, command.BalancePayOnline{}, command.BalancePayOffline{}),
		cmd.New("colour", text.Colourf("<aqua>Customize the colour of your archer.</aqua>"), nil, command.Colour{}),
		cmd.New("clear", text.Colourf("<aqua>Clear your Inventory.</aqua>"), nil, command.Clear{}),
		cmd.New("clearlag", text.Colourf("<aqua>Clears all ground entitys.</aqua>"), nil, command.ClearLag{}),
		cmd.New("logout", text.Colourf("<aqua>Safely logout of PVP.</aqua>"), nil, command.Logout{}),
		cmd.New("pvp", text.Colourf("<aqua>Manage PVP timer.</aqua>"), nil, command.PvpEnable{}),
		cmd.New("role", text.Colourf("<aqua>Manage user roles.</aqua>"), nil, command.RoleAdd{}, command.RoleRemove{}, command.RoleAddOffline{}, command.RoleRemoveOffline{}, command.RoleList{}),
		cmd.New("teleport", text.Colourf("<aqua>Teleport yourself or another player to a position.</aqua>"), []string{"tp"}, command.TeleportToPos{}, command.TeleportTargetsToPos{}, command.TeleportTargetsToTarget{}, command.TeleportToTarget{}),
		cmd.New("reclaim", text.Colourf("<aqua>Manage a reclaim.</aqua>"), nil, command.Reclaim{}, command.ReclaimReset{}),
		cmd.New("kit", text.Colourf("<aqua>Choose a kit.</aqua>"), nil, command.Kit{}),
		cmd.New("ban", text.Colourf("<aqua>Unleash the ban hammer.</aqua>"), nil, command.Ban{}, command.BanOffline{}, command.BanList{}, command.BanLiftOffline{}, command.BanInfoOffline{}, command.BanForm{}),
		cmd.New("blacklist", text.Colourf("<aqua>Blacklist little evaders.</aqua>"), nil, command.Blacklist{}, command.BlacklistOffline{}, command.BlacklistList{}, command.BlacklistLiftOffline{}, command.BlacklistInfoOffline{}, command.BlacklistForm{}),
		cmd.New("kick", text.Colourf("<aqua>Kick a user.</aqua>"), nil, command.Kick{}),
		cmd.New("mute", text.Colourf("<aqua>Mute a user.</aqua>"), nil, command.MuteList{}, command.MuteInfo{}, command.MuteInfoOffline{}, command.MuteLift{}, command.MuteLiftOffline{}, command.MuteForm{}, command.Mute{}, command.MuteOffline{}),
		cmd.New("whisper", text.Colourf("<aqua>Send a private message to a player.</aqua>"), []string{"w", "tell", "msg"}, command.Whisper{}),
		cmd.New("reply", text.Colourf("<aqua>Reply to the last whispered player.</aqua>"), []string{"r"}, command.Reply{}),
		cmd.New("fly", text.Colourf("<aqua>Toggle flight.</aqua>"), nil, command.Fly{}),
		cmd.New("sotw", text.Colourf("<aqua>SOTW management commands.</aqua>"), nil, command.SOTWStart{}, command.SOTWEnd{}, command.SOTWDisable{}),
		cmd.New("freeze", text.Colourf("<aqua>Freeze possible cheaters.</aqua>"), nil, command.Freeze{}),
		cmd.New("gamemode", text.Colourf("<aqua>Manage gamemodes.</aqua>"), []string{"gm"}, command.GameMode{}),
		cmd.New("key", text.Colourf("<aqua>Manage keys</aqua>"), nil, command.Key{}, command.KeyAll{}),
		cmd.New("koth", text.Colourf("<aqua>Manage KOTHs.</aqua>"), nil, command.KothStart{}, command.KothStop{}, command.KothList{}),
		cmd.New("pp", text.Colourf("<aqua>Manage partner packages.</aqua>"), nil, command.PartnerPackageAll{}, command.PartnerPackage{}),
		cmd.New("ping", text.Colourf("<aqua>Check your ping.</aqua>"), nil, command.Ping{}),
		//cmd.New("data", text.Colourf("<aqua>Clear data.</aqua>"), nil, command.DataReset{}),
		cmd.New("vanish", text.Colourf("<aqua>Vanish as staff.</aqua>"), []string{"v"}, command.Vanish{}),
	} {
		cmd.Register(c)
	}
	if sock, ok := moyai.Socket(); ok {
		cmd.Register(cmd.New("hub", text.Colourf("<aqua>Return to the Moyai Hub.</aqua>"), []string{"lobby"}, command.NewHub(sock)))
	}
}

func registerRecipes() {
	recipe.Register(recipe.NewShapeless([]item.Stack{
		item.NewStack(block.Wood{}, 1),
	}, item.NewStack(block.Planks{}, 4), "crafting_table"))

	recipe.Register(recipe.NewShapeless([]item.Stack{
		item.NewStack(block.Log{}, 1),
	}, item.NewStack(block.Planks{}, 4), "crafting_table"))

	recipe.Register(recipe.NewShaped([]item.Stack{
		item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0),
		item.NewStack(block.Air{}, 0), item.NewStack(block.Planks{}, 1), item.NewStack(block.Air{}, 0),
		item.NewStack(block.Air{}, 0), item.NewStack(block.Planks{}, 1), item.NewStack(block.Air{}, 0),
	}, item.NewStack(item.Stick{}, 4), recipe.NewShape(3, 3), "crafting_table"))

	recipe.Register(recipe.NewShaped([]item.Stack{
		item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0),
		item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0),
		item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1),
	}, item.NewStack(block.Slab{
		Block: block.Planks{},
	}, 3), recipe.NewShape(3, 3), "crafting_table"))

	recipe.Register(recipe.NewShaped([]item.Stack{
		item.NewStack(block.Planks{}, 1), item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0),
		item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1), item.NewStack(block.Air{}, 0),
		item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1),
	}, item.NewStack(block.Stairs{
		Block: block.Planks{},
	}, 3), recipe.NewShape(3, 3), "crafting_table"))

	recipe.Register(recipe.NewShaped([]item.Stack{
		item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1),
		item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1),
		item.NewStack(block.Air{}, 0), item.NewStack(item.Stick{}, 1), item.NewStack(block.Air{}, 0),
	}, item.NewStack(block.Sign{
		Wood: block.OakWood(),
	}, 3), recipe.NewShape(3, 3), "crafting_table"))

	recipe.Register(recipe.NewShaped([]item.Stack{
		item.NewStack(block.Air{}, 0), item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1),
		item.NewStack(block.Air{}, 0), item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1),
		item.NewStack(block.Air{}, 0), item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1),
	}, item.NewStack(block.WoodDoor{
		Wood: block.OakWood(),
	}, 3), recipe.NewShape(3, 3), "crafting_table"))

	recipe.Register(recipe.NewShaped([]item.Stack{
		item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1),
		item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1), item.NewStack(block.Planks{}, 1),
		item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0),
	}, item.NewStack(block.WoodTrapdoor{
		Wood: block.OakWood(),
	}, 3), recipe.NewShape(3, 3), "crafting_table"))

	recipe.Register(recipe.NewShaped([]item.Stack{
		item.NewStack(item.Stick{}, 1), item.NewStack(block.Planks{}, 1), item.NewStack(item.Stick{}, 1),
		item.NewStack(item.Stick{}, 1), item.NewStack(block.Planks{}, 1), item.NewStack(item.Stick{}, 1),
		item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0),
	}, item.NewStack(block.WoodFenceGate{
		Wood: block.OakWood(),
	}, 3), recipe.NewShape(3, 3), "crafting_table"))

	recipe.Register(recipe.NewShaped([]item.Stack{
		item.NewStack(block.Planks{}, 1), item.NewStack(item.Stick{}, 1), item.NewStack(block.Planks{}, 1),
		item.NewStack(block.Planks{}, 1), item.NewStack(item.Stick{}, 1), item.NewStack(block.Planks{}, 1),
		item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0), item.NewStack(block.Air{}, 0),
	}, item.NewStack(block.WoodFence{
		Wood: block.OakWood(),
	}, 3), recipe.NewShape(3, 3), "crafting_table"))
}

func registerSlappers(w *world.World) {
	s, err := decodeSkin("./assets/models/steve.png", "geometry.humanoid.custom", "")

	if err != nil {
		panic(err)
	}
	npc.Create(npc.Settings{
		Name:       text.Colourf("<aqua>Diamond</aqua>"),
		Position:   mgl64.Vec3{9, 80, 21},
		Skin:       s,
		Immobile:   true,
		Vulnerable: false,
		Helmet:     item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1),
		Chestplate: item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1),
		Leggings:   item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1),
		Boots:      item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1),
		MainHand:   item.NewStack(item.Sword{item.ToolTierDiamond}, 1),
	}, w, func(p *player.Player) {
		u, err := data.LoadUserOrCreate(p.Name())
		if err != nil {
			return
		}
		cd := u.GameMode.Teams.Kits.Key("diamond")
		if cd.Active() {
			p.Handler().(*user.Handler).Message("command.kit.cooldown", cd.Remaining().Round(time.Second))
			return
		} else {
			cd.Set(5 * time.Minute)
		}
		kit.Apply(kit.Diamond{}, p)
	})
}

func decodeSkin(path, geometry, geometryPath string) (skin.Skin, error) {
	var f *os.File
	var err error

	if f, err = os.OpenFile(path, os.O_RDWR, 0777); os.IsNotExist(err) {
		return skin.Skin{}, err
	}

	img, err := png.Decode(f)
	if err != nil {
		return skin.Skin{}, err
	}
	s := skin.New(img.Bounds().Max.X, img.Bounds().Max.Y)
	s.Pix = pix(img)
	s.ModelConfig = skin.ModelConfig{Default: geometry}
	if geometryPath != "" {
		b, err := ioutil.ReadFile(geometryPath)
		if err != nil {
			panic(err)
		}
		s.Model = b
	}
	return s, nil
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
			sock, ok := moyai.Socket()
			if ok {
				go func() {
					_ = sock.WritePacket(&proxypacket.TransferRequest{
						PlayerUUID: p.UUID(),
						Server:     "syn.lobby",
					})
				}()
			}
		}
		time.Sleep(time.Millisecond * 500)
		_ = data.CloseUsers()
		data.CloseTeams()
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

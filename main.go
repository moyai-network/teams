package main

import (
	"github.com/bedrock-gophers/console/console"
	"github.com/bedrock-gophers/intercept"
	"github.com/bedrock-gophers/inv/inv"
	"github.com/bedrock-gophers/tebex/tebex"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	cr "github.com/moyai-network/teams/moyai/crate"
	"github.com/moyai-network/teams/moyai/data"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	ent "github.com/moyai-network/teams/moyai/entity"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft"
	"image"
	"math"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func main() {
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	lang.Register(language.English)

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	console.Enable(log)
	conf, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	config := configure(conf, log)
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

	w := srv.World()
	configureWorld(w)
	clearEntities(w)
	placeCrates(w)
	inv.PlaceFakeContainer(w, cube.Pos{0, 255, 0})

	registerCommands(srv)

	srv.Listen()

	store := loadStore(conf.Moyai.Tebex, log)

	for srv.Accept(acceptFunc(store, conf.Proxy.Enabled)) {
		// Do nothing.
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
		t := entity.NewText(text.Colourf("%s <grey>Crate</grey>\n<yellow>Right click to open crate</yellow>\n<grey>Left click to see rewards</grey>", c.Name()), c.Position().Add(mgl64.Vec3{0, 2, 0}))
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
}

func acceptFunc(store *tebex.Client, proxy bool) func(*player.Player) {
	return func(p *player.Player) {
		inv.RedirectPlayerPackets(p)
		store.ExecuteCommands(p)
		if proxy {
			//info := moyai.SearchInfo(p.UUID())
			//p.Handle(user.NewHandler(p, info.XUID))
		} else {
			p.Handle(user.NewHandler(p, p.XUID()))
			p.Armour().Handle(user.NewArmourHandler(p))
		}
		p.SetGameMode(world.GameModeSurvival)
		p.RemoveScoreboard()
		for _, ef := range p.Effects() {
			p.RemoveEffect(ef.Type())
		}
		p.ShowCoordinates()
		p.SetFood(20)

		u, _ := data.LoadUserOrCreate(p.Name(), p.XUID())
		if !u.Roles.Contains(role.Default{}) {
			u.Roles.Add(role.Default{})
		}
		p.Message(lang.Translatef(u.Language, "discord.message"))
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
		cmd.New("t", text.Colourf("<dark-red>The main team management command.</dark-red>"), []string{"f"},
			command.TeamCreate{},
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
			//command.TeamRally{},
			//command.TeamUnRally{},
			//command.TeamMap{},
			command.TeamSetDTR{},
			command.TeamDelete{},
		), cmd.New("whitelist", text.Colourf("<dark-red>Whitelist commands.</dark-red>"), []string{"wl"}, command.WhiteListAdd{}, command.WhiteListRemove{}),
		cmd.New("tl", text.Colourf("<dark-red>Send your location to teammates</dark-red>"), nil, command.TL{}),
		cmd.New("balance", text.Colourf("<dark-red>Manage your balance.</dark-red>"), []string{"bal"}, command.Balance{}, command.BalancePayOnline{}, command.BalancePayOffline{}, command.BalanceAdd{}, command.BalanceAddOffline{}),
		cmd.New("colour", text.Colourf("<dark-red>Customize the colour of your archer.</dark-red>"), nil, command.Colour{}),
		cmd.New("clear", text.Colourf("<dark-red>Clear your Inventory.</dark-red>"), nil, command.Clear{}),
		cmd.New("clearlag", text.Colourf("<dark-red>Clears all ground entitys.</dark-red>"), nil, command.ClearLag{}),
		cmd.New("logout", text.Colourf("<dark-red>Safely logout of PVP.</dark-red>"), nil, command.Logout{}),
		cmd.New("pvp", text.Colourf("<dark-red>Manage PVP timer.</dark-red>"), nil, command.PvpEnable{}),
		cmd.New("role", text.Colourf("<dark-red>Manage user roles.</dark-red>"), nil, command.RoleAdd{}, command.RoleRemove{}, command.RoleAddOffline{}, command.RoleRemoveOffline{}, command.RoleList{}),
		cmd.New("teleport", text.Colourf("<dark-red>Teleport yourself or another player to a position.</dark-red>"), []string{"tp"}, command.TeleportToPos{}, command.TeleportTargetsToPos{}, command.TeleportTargetsToTarget{}, command.TeleportToTarget{}),
		cmd.New("reclaim", text.Colourf("<dark-red>Manage a reclaim.</dark-red>"), nil, command.Reclaim{}, command.ReclaimReset{}),
		cmd.New("kit", text.Colourf("<dark-red>Choose a kit.</dark-red>"), nil, command.Kit{}),
		//cmd.New("ban", text.Colourf("<dark-red>Unleash the ban hammer.</dark-red>"), nil, command.Ban{}, command.BanOffline{}, command.BanList{}, command.BanLiftOffline{}, command.BanInfoOffline{}, command.BanForm{}),
		//cmd.New("blacklist", text.Colourf("<dark-red>Blacklist little evaders.</dark-red>"), nil, command.Blacklist{}, command.BlacklistOffline{}, command.BlacklistList{}, command.BlacklistLiftOffline{}, command.BlacklistInfoOffline{}, command.BlacklistForm{}),
		cmd.New("kick", text.Colourf("<dark-red>Kick a user.</dark-red>"), nil, command.Kick{}),
		//cmd.New("mute", text.Colourf("<dark-red>Mute a user.</dark-red>"), nil, command.MuteList{}, command.MuteInfo{}, command.MuteInfoOffline{}, command.MuteLift{}, command.MuteLiftOffline{}, command.MuteForm{}, command.Mute{}, command.MuteOffline{}),
		cmd.New("whisper", text.Colourf("<dark-red>Send a private message to a player.</dark-red>"), []string{"w", "tell", "msg"}, command.Whisper{}),
		cmd.New("reply", text.Colourf("<dark-red>Reply to the last whispered player.</dark-red>"), []string{"r"}, command.Reply{}),
		cmd.New("fly", text.Colourf("<dark-red>Toggle flight.</dark-red>"), nil, command.Fly{}),
		cmd.New("sotw", text.Colourf("<dark-red>SOTW management commands.</dark-red>"), nil, command.SOTWStart{}, command.SOTWEnd{}, command.SOTWDisable{}),
		cmd.New("freeze", text.Colourf("<dark-red>Freeze possible cheaters.</dark-red>"), nil, command.Freeze{}),
		cmd.New("gamemode", text.Colourf("<dark-red>Manage gamemodes.</dark-red>"), []string{"gm"}, command.GameMode{}),
		cmd.New("key", text.Colourf("<dark-red>Manage keys</dark-red>"), nil, command.Key{}, command.KeyAll{}),
		cmd.New("koth", text.Colourf("<dark-red>Manage KOTHs.</dark-red>"), nil, command.KothStart{}, command.KothStop{}, command.KothList{}),
		cmd.New("pp", text.Colourf("<dark-red>Manage partner packages.</dark-red>"), nil, command.PartnerPackageAll{}, command.PartnerPackage{}),
		cmd.New("ping", text.Colourf("<dark-red>Check your ping.</dark-red>"), nil, command.Ping{}),
		//cmd.New("data", text.Colourf("<dark-red>Clear data.</dark-red>"), nil, command.DataReset{}),
		cmd.New("vanish", text.Colourf("<dark-red>Vanish as staff.</dark-red>"), []string{"v"}, command.Vanish{}),
	} {
		cmd.Register(c)
	}

	//cmd.Register(cmd.New("hub", text.Colourf("<dark-red>Return to the Moyai Hub.</dark-red>"), []string{"lobby"}, command.Hub{}))
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

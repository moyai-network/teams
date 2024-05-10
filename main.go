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
	log.Level = logrus.InfoLevel

	console.Enable(log)
	conf, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	config := configure(conf, log)
	srv := moyai.NewServer(config)

	handleServerClose(srv)

	w := srv.World()
	configureWorld(w)
	clearEntities(w)
	placeCrates(w)
	inv.PlaceFakeContainer(w, cube.Pos{0, 255, 0})

	registerCommands(srv)

	pk := intercept.NewPacketListener()
	pk.Listen(&config, ":19133", []minecraft.Protocol{})
	go func() {
		for {
			p, err := pk.Accept()
			if err != nil {
				return
			}
			p.Handle(user.NewPacketHandler(p))
		}
	}()
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
		store.ExecuteCommands(p)
		if proxy {
			//info := moyai.SearchInfo(p.UUID())
			//p.Handle(user.NewHandler(p, info.XUID))
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
			//command.TeamMap{},
			command.TeamSetDTR{},
			command.TeamDelete{},
		), cmd.New("whitelist", text.Colourf("<aqua>Whitelist commands.</aqua>"), []string{"wl"}, command.WhiteListAdd{}, command.WhiteListRemove{}),
		cmd.New("tl", text.Colourf("<aqua>Send your location to teammates</aqua>"), nil, command.TL{}),
		cmd.New("balance", text.Colourf("<aqua>Manage your balance.</aqua>"), []string{"bal"}, command.Balance{}, command.BalancePayOnline{}, command.BalancePayOffline{}, command.BalanceAdd{}, command.BalanceAddOffline{}),
		cmd.New("colour", text.Colourf("<aqua>Customize the colour of your archer.</aqua>"), nil, command.Colour{}),
		cmd.New("clear", text.Colourf("<aqua>Clear your Inventory.</aqua>"), nil, command.Clear{}),
		cmd.New("clearlag", text.Colourf("<aqua>Clears all ground entitys.</aqua>"), nil, command.ClearLag{}),
		cmd.New("logout", text.Colourf("<aqua>Safely logout of PVP.</aqua>"), nil, command.Logout{}),
		cmd.New("pvp", text.Colourf("<aqua>Manage PVP timer.</aqua>"), nil, command.PvpEnable{}),
		cmd.New("role", text.Colourf("<aqua>Manage user roles.</aqua>"), nil, command.RoleAdd{}, command.RoleRemove{}, command.RoleAddOffline{}, command.RoleRemoveOffline{}, command.RoleList{}),
		cmd.New("teleport", text.Colourf("<aqua>Teleport yourself or another player to a position.</aqua>"), []string{"tp"}, command.TeleportToPos{}, command.TeleportTargetsToPos{}, command.TeleportTargetsToTarget{}, command.TeleportToTarget{}),
		cmd.New("reclaim", text.Colourf("<aqua>Manage a reclaim.</aqua>"), nil, command.Reclaim{}, command.ReclaimReset{}),
		cmd.New("kit", text.Colourf("<aqua>Choose a kit.</aqua>"), nil, command.Kit{}),
		//cmd.New("ban", text.Colourf("<aqua>Unleash the ban hammer.</aqua>"), nil, command.Ban{}, command.BanOffline{}, command.BanList{}, command.BanLiftOffline{}, command.BanInfoOffline{}, command.BanForm{}),
		//cmd.New("blacklist", text.Colourf("<aqua>Blacklist little evaders.</aqua>"), nil, command.Blacklist{}, command.BlacklistOffline{}, command.BlacklistList{}, command.BlacklistLiftOffline{}, command.BlacklistInfoOffline{}, command.BlacklistForm{}),
		cmd.New("kick", text.Colourf("<aqua>Kick a user.</aqua>"), nil, command.Kick{}),
		//cmd.New("mute", text.Colourf("<aqua>Mute a user.</aqua>"), nil, command.MuteList{}, command.MuteInfo{}, command.MuteInfoOffline{}, command.MuteLift{}, command.MuteLiftOffline{}, command.MuteForm{}, command.Mute{}, command.MuteOffline{}),
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

	cmd.Register(cmd.New("hub", text.Colourf("<aqua>Return to the Moyai Hub.</aqua>"), []string{"lobby"}, command.Hub{}))
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

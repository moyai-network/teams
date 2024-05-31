package minecraft

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bedrock-gophers/console/console"
	"github.com/bedrock-gophers/intercept"
	"github.com/bedrock-gophers/inv/inv"
	"github.com/bedrock-gophers/tebex/tebex"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	ent "github.com/moyai-network/teams/moyai/entity"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/oomph-ac/oomph"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

// Run runs the Minecraft server.
func Run() error {
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	registerLanguages()

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	console.Enable(log)
	conf, err := readConfig()
	if err != nil {
		return err
	}
	config := configure(conf, log)

	srv := moyai.NewServer(config)
	handleServerClose(srv)

	registerCommands(srv)
	srv.Listen()

	moyai.ConfigureDimensions(config.Entities, conf.Nether.Folder, conf.End.Folder)
	moyai.ConfigureDeathban(config.Entities, "./assets/deathban2")
	configureWorlds()

	placeSpawners()
	placeText(conf)
	placeSlapper()
	placeCrates()
	placeShopSigns()

	go tickBlackMarket(srv)
	go tickClearLag()

	store := loadStore(conf.Moyai.Tebex, log)
	for srv.Accept(acceptFunc(store, conf.Proxy.Enabled, srv)) {
		// Do nothing.
	}

	return nil
}

// registerLanguages registers all languages that are available in the server.
func registerLanguages() {
	lang.Register(language.English)
	lang.Register(language.Spanish)
	lang.Register(language.French)
}

// tickBlackMarket runs a ticker that checks every 15 minutes if the black market should be opened. The black
// market is opened randomly with a 25% chance every 15 minutes, but only if the last black market was opened
// more than an hour ago.
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

// configure initializes the server configuration.
func configure(conf moyai.Config, log *logrus.Logger) server.Config {
	c, err := conf.Config(log)
	if err != nil {
		panic(err)
	}
	c.Entities = ent.Registry

	c.Name = text.Colourf("<bold><redstone>MOYAI</redstone></bold>") + "§8"
	c.Allower = moyai.NewAllower(conf.Moyai.Whitelisted)

	configurePacketListener(&c, conf.Oomph.Enabled)
	return c
}

// configurePacketListener configures the packet listener for the server.
func configurePacketListener(conf *server.Config, oomphEnabled bool) {
	if oomphEnabled {
		ac := oomph.New(oomph.OomphSettings{
			LocalAddress:  ":19133",
			RemoteAddress: ":19132",
			RequirePacks:  true,
		})

		ac.Listen(conf, text.Colourf(conf.Name), []minecraft.Protocol{}, true, false)
		go func() {
			for {
				_, err := ac.Accept()
				if err != nil {
					return
				}

			}
		}()
		return
	}
	pk := intercept.NewPacketListener()
	pk.Listen(conf, ":19132", []minecraft.Protocol{})

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

// handleServerClose handles the closing of the server.
func handleServerClose(srv *server.Server) {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		time.Sleep(time.Millisecond * 500)
		data.FlushCache()

		sotw.Save()
		err := moyai.Nether().Close()
		if err != nil {
			logrus.Fatalln("close nether: %v", err)
		}
		moyai.End().Close()
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

// acceptFunc returns a function that is called when a player joins the server.
func acceptFunc(store *tebex.Client, proxy bool, srv *server.Server) func(*player.Player) {
	return func(p *player.Player) {
		inv.RedirectPlayerPackets(p, func() {
			time.Sleep(time.Millisecond * 500)
			data.FlushCache()

			sotw.Save()
			_ = srv.Close()
			os.Exit(1)
		})

		store.ExecuteCommands(p)

		u, _ := data.LoadUserOrCreate(p.Name(), p.XUID())
		if !u.Roles.Contains(role.Default{}) {
			u.Roles.Add(role.Default{})
		}
		u.Roles.Add(role.Donor1{})

		p.Message(lang.Translatef(u.Language, "discord.message"))
		p.Handle(user.NewHandler(p, p.XUID()))
		p.Armour().Handle(user.NewArmourHandler(p))
		p.RemoveScoreboard()
		for _, ef := range p.Effects() {
			p.RemoveEffect(ef.Type())
		}
		p.ShowCoordinates()
		p.SetFood(20)

		data.SaveUser(u)
		in := p.Inventory()
		for slot, i := range in.Slots() {
			for _, sp := range it.SpecialItems() {
				if _, ok := i.Value(sp.Key()); ok {
					_ = in.SetItem(slot, it.NewSpecialItem(sp, i.Count()))
				}
			}
		}

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
	}
}

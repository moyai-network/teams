package minecraft

import (
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
	v486 "github.com/flonja/multiversion/protocols/v486"

	_ "github.com/flonja/multiversion/protocols" // VERY IMPORTANT
	// v486 "github.com/flonja/multiversion/protocols/v486"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	ent "github.com/moyai-network/teams/moyai/entity"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

// Run runs the Minecraft server.
func Run() error {
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
	handleServerClose()

	registerCommands()
	srv.Listen()

	moyai.ConfigureDimensions(config.Entities, conf.Nether.Folder, conf.End.Folder)
	moyai.ConfigureDeathban(config.Entities, conf.DeathBan.Folder)
	configureWorlds()

	placeSpawners()
	placeText(conf)
	placeSlapper()
	placeCrates()
	placeShopSigns()

	go tickVotes()
	go tickBlackMarket(srv)
	go tickClearLag()

	store := loadStore(conf.Moyai.Tebex, log)
	for srv.Accept(acceptFunc(store)) {
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

// tickVotes ticks new votes every 5 minutes of all users.
func tickVotes() {
	t := time.NewTicker(time.Second * 5)
	go func() {
		for range t.C {
			usrs := data.NewVoters()
			for _, u := range usrs {
				u.Roles.Add(role.Voter{})
				u.Roles.Expire(role.Voter{}, time.Now().Add(time.Hour*24))
				moyai.Broadcastf("vote.broadcast", u.DisplayName)
				data.SaveUser(u)
			}
		}
	}()
}

// configure initializes the server configuration.
func configure(conf moyai.Config, log *logrus.Logger) server.Config {
	c, err := conf.Config(log)
	if err != nil {
		panic(err)
	}
	c.Entities = ent.Registry

	c.Name = text.Colourf("<bold><redstone>MOYAI</redstone></bold>") + "§8"
	c.ShutdownMessage = text.Colourf("<red>MoyaiHCF has restarted; please join back shortly or join discord.gg/moyai for more info!</red>") + "§8"
	c.JoinMessage = "<green>[+] %s</green>"
	c.QuitMessage = "<red>[-] %s</red>"
	c.Allower = moyai.NewAllower(conf.Moyai.Whitelisted)

	configurePacketListener(&c, conf.Oomph.Enabled)
	return c
}

// configurePacketListener configures the packet listener for the server.
func configurePacketListener(conf *server.Config, oomphEnabled bool) {
	/*if oomphEnabled {
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
	}*/
	pk := intercept.NewPacketListener()
	pk.Listen(conf, ":19132", []minecraft.Protocol{
		v486.New(),
		// v582.New(),
		// v589.New(),
		// v594.New(),
		// v618.New(),
		// v622.New(),
		// v630.New(),
		// v649.New(),
		// v662.New(),
	})

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
func handleServerClose() {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		moyai.Close()
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
func acceptFunc(store *tebex.Client) func(*player.Player) {
	return func(p *player.Player) {
		inv.RedirectPlayerPackets(p, func() {
			moyai.Close()
			os.Exit(1)
		})
		store.ExecuteCommands(p)

		p.Handle(user.NewHandler(p, p.XUID()))
		p.Armour().Handle(user.NewArmourHandler(p))
		p.RemoveScoreboard()
		for _, ef := range p.Effects() {
			p.RemoveEffect(ef.Type())
		}
		p.ShowCoordinates()
		p.SetFood(20)

		in := p.Inventory()
		for slot, i := range in.Slots() {
			for _, sp := range append(it.SpecialItems(), it.PartnerItems()...) {
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

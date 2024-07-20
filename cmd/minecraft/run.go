package minecraft

import (
	"fmt"
	"github.com/bedrock-gophers/intercept/intercept"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bedrock-gophers/knockback/knockback"

	"github.com/bedrock-gophers/tag/tag"

	"github.com/bedrock-gophers/role/role"
	"github.com/moyai-network/teams/moyai/roles"
	"github.com/oomph-ac/oomph"
	"github.com/oomph-ac/oomph/handler"

	"github.com/bedrock-gophers/console/console"
	"github.com/bedrock-gophers/tebex/tebex"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"

	_ "github.com/flonja/multiversion/protocols" // VERY IMPORTANT

	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	ent "github.com/moyai-network/teams/moyai/entity"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

// Run runs the Minecraft server.
func Run() error {
	err := knockback.Load("assets/knockback.json")
	if err != nil {
		return err
	}
	err = role.Load("assets/roles")
	if err != nil {
		return err
	}
	err = tag.Load("assets/tags")
	if err != nil {
		return err
	}
	registerLanguages()

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	// chat.Global.Subscribe(chat.StdoutSubscriber{})

	console.Enable(log)
	conf, err := readConfig()
	if err != nil {
		return err
	}
	config := configure(conf, log)

	srv := moyai.NewServer(config)
	handleServerClose()

	registerCommands()
	store := loadStore(conf.Moyai.Tebex, log)
	srv.Listen()
	go executeTebexCommands(store, srv)

	moyai.ConfigureDimensions(config.Entities, conf.Nether.Folder, conf.End.Folder)
	moyai.ConfigureDeathban(config.Entities, conf.DeathBan.Folder)
	configureWorlds()

	//placeSpawners()
	placeText(conf)
	placeSlapper()
	placeCrates()
	placeShopSigns()

	go tickVotes()
	//go tickBlackMarket(srv)
	go tickClearLag()

	go startBroadcats()
	go startPlayerBroadcasts()

	go startLeaderboard()

	go startChatGame()
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
	for range t.C {
		usrs := data.NewVoters()
		for _, u := range usrs {
			u.Roles.Add(roles.Voter())
			u.Roles.Expire(roles.Voter(), time.Now().Add(time.Hour*24))
			moyai.Broadcastf("vote.broadcast", u.DisplayName)
			data.SaveUser(u)
		}
	}
}

// startBroadcats starts the broadcasts.
func startBroadcats() {
	broadcasts := [...]string{
		"moyai.broadcast.discord",
		"moyai.broadcast.store",
		"moyai.broadcast.emojis",
		"moyai.broadcast.settings",
		"moyai.broadcast.feedback",
		//"moyai.broadcast.report",
		"moyai.broadcast.rules",
		"moyai.broadcast.vote",
	}

	var cursor int
	t := time.NewTicker(time.Minute * 4)
	defer t.Stop()
	for range t.C {
		message := broadcasts[cursor]
		for _, p := range moyai.Players() {
			u, err := data.LoadUserFromName(p.Name())
			if err != nil {
				continue
			}
			p.Message(lang.Translatef(*u.Language, "moyai.broadcast.notice", lang.Translate(*u.Language, message)))
		}
		if cursor++; cursor == len(broadcasts) {
			cursor = 0
		}
	}
}

// startPlayerBroadcasts starts the player broadcasts.
func startPlayerBroadcasts() {
	t := time.NewTicker(time.Minute * 5)
	for range t.C {
		players := moyai.Players()
		var plus []string
		for _, p := range players {
			u, err := data.LoadUserFromName(p.Name())
			if err != nil {
				continue
			}
			if roles.Premium(u.Roles.Highest()) {
				plus = append(plus, p.Name())
			}
		}

		for _, p := range players {
			u, err := data.LoadUserFromName(p.Name())
			if err != nil {
				continue
			}
			p.Message(lang.Translatef(*u.Language, "moyai.broadcast.plus", len(plus), strings.Join(plus, ", ")))
		}
	}
}

var chatGameWords = []string{
	"moyai",
	"beacon",
	"diamond",
	"nether",
	"ender",
	"dragon",
	"blaze",
	"wither",
	"creeper",
	"spider",
	"zombie",
	"skeleton",
	"pig",
	"sheep",
	"cow",
	"chicken",
	"squid",
	"wolf",
	"villager",
	"rabbit",
	"guardian",
	"shulker",
	"sword",
	"shield",
	"bow",
	"carrot",
}

// startChatGame starts the chat game.
func startChatGame() {
	t := time.NewTicker(time.Minute * 10)
	for range t.C {
		cursor := rand.Intn(len(chatGameWords))
		word := chatGameWords[cursor]

		scramble := func(word string) string {
			runes := []rune(word)
			for i := 0; i < len(runes); i++ {
				j := rand.Intn(len(runes))
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		}

		scrambled := scramble(word)
		moyai.SetChatGameWord(word)
		for _, p := range moyai.Players() {
			u, err := data.LoadUserFromName(p.Name())
			if err != nil {
				continue
			}
			p.Message(lang.Translatef(*u.Language, "moyai.broadcast.chatgame", scrambled))
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
	c.ShutdownMessage = text.Colourf("<red>MoyaiHCF has restarted; please join back shortly or join discord.moyai.club for more info!</red>") + "§8"
	c.JoinMessage = "<green>[+] %s</green>"
	c.QuitMessage = "<red>[-] %s</red>"
	c.Allower = moyai.NewAllower(conf.Moyai.Whitelisted)

	configurePacketListener(&c, conf.Oomph.Enabled)
	return c
}

// configurePacketListener configures the packet listener for the server.
func configurePacketListener(conf *server.Config, oomphEnabled bool) {
	if oomphEnabled {
		ac := oomph.New(oomph.OomphSettings{
			LocalAddress:  ":19132",
			RemoteAddress: ":19133",
			RequirePacks:  true,
		})

		ac.Listen(conf, text.Colourf(conf.Name), []minecraft.Protocol{}, true, false)
		go func() {
			for {
				p, err := ac.Accept()
				if err != nil {
					return
				}

				p.Player.SetLog(logrus.New())
				p.Player.MovementMode = 0
				p.Player.Handler(handler.HandlerIDMovement).(*handler.MovementHandler).CorrectionThreshold = 100000000

				// TODO: Handle events
			}
		}()
		return
	}
}

func executeTebexCommands(store *tebex.Client, srv *server.Server) {
	for {
		for _, p := range srv.Players() {
			store.ExecuteCommands(p)
			<-time.After(time.Second / 4)
		}
		<-time.After(time.Second * 20)
	}
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
		intercept.Intercept(p)
		p.Handle(temporaryHandler{})
		go store.ExecuteCommands(p)

		h, err := user.NewHandler(p, p.XUID())
		if err != nil {
			fmt.Printf("new handler: %v\n", err)
			p.Disconnect(text.Colourf("<red>Unknown Error. Please contact developers at discord.moyai.club</red>"))
			return
		}
		p.Handle(h)
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

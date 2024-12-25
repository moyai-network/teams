package minecraft

import (
	"fmt"
	data2 "github.com/moyai-network/teams/internal/core/data"
	ent "github.com/moyai-network/teams/internal/core/entity"
	it "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/core/roles"
	user2 "github.com/moyai-network/teams/internal/core/user"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bedrock-gophers/intercept/intercept"

	"github.com/bedrock-gophers/knockback/knockback"

	"github.com/bedrock-gophers/tag/tag"

	"github.com/bedrock-gophers/console/console"
	"github.com/bedrock-gophers/role/role"
	"github.com/bedrock-gophers/tebex/tebex"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"

	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/text/language"
)

// Run runs the Minecraft server.
func Run() error {
	internal.Assemble()

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

	// chat.Global.Subscribe(chat.StdoutSubscriber{})
	conf, err := readConfig()
	if err != nil {
		return err
	}
	config := configure(conf, slog.Default())

	srv := internal.NewServer(config)
	console.Enable(srv, slog.Default())
	handleServerClose()

	registerCommands()
	store := loadStore(conf.Moyai.Tebex, slog.Default())
	srv.Listen()
	//go executeTebexCommands(store, srv)

	internal.ConfigureDimensions(config.Entities, conf.Nether.Folder, conf.End.Folder)
	internal.ConfigureDeathban(config.Entities, conf.DeathBan.Folder)
	configureWorlds()

	//placeSpawners()
	//placeText(conf)
	//placeSlapper()
	placeCrates()
	placeShopSigns()

	go tickVotes()
	//go tickBlackMarket(srv)
	go tickClearLag()

	go startBroadcats()
	go startPlayerBroadcasts()

	go startLeaderboard()

	go startChatGame()
	for p := range srv.Accept() {
		acceptFunc(store)(p)
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
		usrs := data2.NewVoters()
		for _, u := range usrs {
			u.Roles.Add(roles.Voter())
			u.Roles.Expire(roles.Voter(), time.Now().Add(time.Hour*24))
			internal.Broadcastf(nil, "vote.broadcast", u.DisplayName)
			data2.SaveUser(u)
		}
	}
}

// startBroadcats starts the broadcasts.
func startBroadcats() {
	broadcasts := [...]string{
		"internal.broadcast.discord",
		"internal.broadcast.store",
		"internal.broadcast.emojis",
		"internal.broadcast.settings",
		"internal.broadcast.feedback",
		//"internal.broadcast.report",
		"internal.broadcast.rules",
		"internal.broadcast.vote",
	}

	var cursor int
	t := time.NewTicker(time.Minute * 4)
	defer t.Stop()
	for range t.C {
		message := broadcasts[cursor]
		for p := range internal.Players(nil) {
			u, err := data2.LoadUserFromName(p.Name())
			if err != nil {
				continue
			}
			p.Message(lang.Translatef(*u.Language, "internal.broadcast.notice", lang.Translate(*u.Language, message)))
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
		players := internal.Players(nil)
		var plus []string
		for p := range players {
			u, err := data2.LoadUserFromName(p.Name())
			if err != nil {
				continue
			}
			if roles.Premium(u.Roles.Highest()) {
				plus = append(plus, p.Name())
			}
		}

		for p := range players {
			u, err := data2.LoadUserFromName(p.Name())
			if err != nil {
				continue
			}
			p.Message(lang.Translatef(*u.Language, "internal.broadcast.plus", len(plus), strings.Join(plus, ", ")))
		}
	}
}

var chatGameWords = []string{
	"internal",
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
		internal.SetChatGameWord(word)
		for p := range internal.Players(nil) {
			u, err := data2.LoadUserFromName(p.Name())
			if err != nil {
				continue
			}
			p.Message(lang.Translatef(*u.Language, "internal.broadcast.chatgame", scrambled))
		}
	}
}

// configure initializes the server configuration.
func configure(conf internal.Config, log *slog.Logger) server.Config {
	c, err := conf.Config(log)
	if err != nil {
		panic(err)
	}
	c.Entities = ent.Registry

	c.Name = text.Colourf("<bold><redstone>MOYAI</redstone></bold>") + "§8"
	//c.ShutdownMessage = text.Colourf("<red>MoyaiHCF has restarted; please join back shortly or join discord.internal.club for more info!</red>") + "§8"
	//c.JoinMessage = "<green>[+] %s</green>"
	//c.QuitMessage = "<red>[-] %s</red>"
	c.Allower = internal.NewAllower(conf.Moyai.Whitelisted)

	//configurePacketListener(&c, conf.Oomph.Enabled)
	return c
}

// configurePacketListener configures the packet listener for the server.
func configurePacketListener(conf *server.Config, oomphEnabled bool) {
	/*if oomphEnabled {
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
	*/
}

/*func executeTebexCommands(store *tebex.Client, srv *server.Server) {
	for {
		for p := range srv.Players() {
			store.ExecuteCommands(p)
			<-time.After(time.Second / 4)
		}
		<-time.After(time.Second * 20)
	}
}*/

// loadStore initializes the Tebex store connection.
func loadStore(key string, log *slog.Logger) *tebex.Client {
	store := tebex.NewClient(log, time.Second*5, key)
	name, domain, err := store.Information()
	if err != nil {
		log.Error("tebex: %v", err)
		return nil
	}
	log.Info("Connected to Tebex under %v (%v).", name, domain)
	return store
}

// handleServerClose handles the closing of the server.
func handleServerClose() {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		internal.Close()
	}()
}

func readConfig() (internal.Config, error) {
	c := internal.DefaultConfig()
	g := gophig.NewGophig("./configs/config", "toml", 0777)

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
		//p.Handle(temporaryHandler{})
		go store.ExecuteCommands(p)

		h, err := user2.NewHandler(p, p.XUID())
		if err != nil {
			fmt.Printf("new handler: %v\n", err)
			p.Disconnect(text.Colourf("<red>Unknown Error. Please contact developers at discord.internal.club</red>"))
			return
		}
		p.Handle(h)
		p.Armour().Handle(user2.NewArmourHandler(p))
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

		//p.SetImmobile()
		/*p.SetAttackImmunity(time.Millisecond * 500)
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
		}*/
	}
}

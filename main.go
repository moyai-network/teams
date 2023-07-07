package main

import (
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bedrock-gophers/packethandler"
	"github.com/moyai-network/moose/worlds"
	"github.com/moyai-network/teams/moyai/data"
	ent "github.com/moyai-network/teams/moyai/entity"
	"github.com/moyai-network/teams/moyai/sotw"
	proxypacket "github.com/paroxity/portal/socket/packet"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	_ "github.com/moyai-network/moose/console"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/command"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

func main() {
	lang.Register(language.English)

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	c, err := config.Config(log)
	if err != nil {
		panic(err)
	}

	c.Entities = ent.Registry

	c.Name = text.Colourf("<bold><redstone>SYN</redstone></bold>") + "§8"
	c.Generator = func(dim world.Dimension) world.Generator { return nil }
	c.Allower = moyai.NewAllower(config.Moyai.Whitelisted)

	pk := packethandler.NewPacketListener()
	pk.Listen(":19134", &c, true)
	go func() {
		for {
			p, err := pk.Accept()
			if err != nil {
				return
			}
			p.Handle(user.NewPacketHandler(p))
		}
	}()

	srv := c.New()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		sotw.Save()
		for _, p := range srv.Players() {
			p.Message(text.Colourf("<green>Travelling to <black>The</black> <gold>Hub</gold>...</green>"))
			_ = moyai.Socket().WritePacket(&proxypacket.TransferRequest{
				PlayerUUID: p.UUID(),
				Server:     "syn.lobby",
			})
		}
		time.Sleep(time.Second)
		_ = data.Close()
		if err := srv.Close(); err != nil {
			log.Errorf("close server: %v", err)
		}
	}()

	w := srv.World()
	w.Handle(&worlds.Handler{})
	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeSurvival)
	w.SetTime(6000)
	w.StopTime()
	w.SetTickRange(0)
	w.StopThundering()

	l := world.NewLoader(8, w, world.NopViewer{})
	l.Move(w.Spawn().Vec3Middle())
	l.Load(math.MaxInt)

	for _, e := range w.Entities() {
		w.RemoveEntity(e)
	}
	registerCommands(srv)

	srv.Listen()
	for srv.Accept(acceptFunc(config.Proxy.Enabled)) {
		// Do nothing.
	}
}

func acceptFunc(proxy bool) func(*player.Player) {
	return func(p *player.Player) {
		if proxy {
			info := moyai.SearchInfo(p.UUID())
			p.Handle(user.NewHandler(p, info.XUID))
		} else {
			p.Handle(user.NewHandler(p, p.XUID()))
		}
		p.SetGameMode(world.GameModeSurvival)
	}
}

func registerCommands(srv *server.Server) {
	for _, c := range []cmd.Command{
		cmd.New("t", text.Colourf("<aqua>The main team management command.</aqua>"), []string{"f"},
			command.TeamCreate{},
			command.NewTeamInformation(srv),
			//command.TeamDisband{},
			command.TeamInvite{},
			command.TeamJoin{},
			command.NewTeamWho(srv),
			command.TeamLeave{},
			command.TeamKick{},
			// command.TeamPromote{},
			// command.TeamDemote{},
			// command.TeamTop{},
			// command.TeamClaim{},
			// command.TeamUnClaim{},
			// command.TeamSetHome{},
			// command.TeamHome{},
			// command.TeamList{},
			// command.TeamUnFocus{},
			// command.TeamFocusPlayer{},
			// command.TeamFocusTeam{},
			// command.TeamChat{},
			// command.TeamWithdraw{},
			// command.TeamDeposit{},
			// command.TeamWithdrawAll{},
			// command.TeamDepositAll{},
			// command.TeamStuck{},
		), cmd.New("whitelist", text.Colourf("<aqua>Whitelist commands.</aqua>"), []string{"wl"}, command.WhiteListAdd{}, command.WhiteListRemove{}),
		cmd.New("balance", text.Colourf("<aqua>Manage your balance.</aqua>"), []string{"bal"}, command.Balance{}, command.BalancePayOnline{}, command.BalancePayOffline{}),
		cmd.New("colour", text.Colourf("<aqua>Customize the colour of your archer.</aqua>"), nil, command.Colour{}),
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
		cmd.New("sotw", text.Colourf("SOTW management commands.</aqua>"), nil, command.SOTWStart{}, command.SOTWEnd{}, command.SOTWDisable{}),
	} {
		cmd.Register(c)
	}
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

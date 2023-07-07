package main

import (
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/bedrock-gophers/packethandler"
	"github.com/moyai-network/teams/moyai/data"
	ent "github.com/moyai-network/teams/moyai/entity"

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

	proxypacket "github.com/paroxity/portal/socket/packet"
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
		_ = data.Close()
		if err := srv.Close(); err != nil {
			log.Errorf("close server: %v", err)
		}
	}()

	w := srv.World()
	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeAdventure)
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
	registerCommands()

	srv.Listen()
	for srv.Accept(accept) {
		// Do nothing.
	}
}

func accept(p *player.Player) {
	// config, err := readConfig()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	if true {
		moyai.Socket().WritePacket(&proxypacket.PlayerInfoRequest{
			PlayerUUID: p.UUID(),
		})
		moyai.SocketUnlock()
		logrus.Info("Allah1")
		ch := moyai.PlayerInfo()
		info := <-ch
		logrus.Info(info)
		if info.PlayerUUID != p.UUID() {
			logrus.Info("Allah3")
			inf, ok := moyai.SearchInfo(p.UUID())
			if !ok {
				logrus.Info("Allah4")
				p.Disconnect("fuckshit")
			}
			logrus.Info("Allah4")
			info = inf
		}
		logrus.Info("Allah5")
		p.Handle(user.NewHandler(p, info.XUID))
	} else {
		p.Handle(user.NewHandler(p, p.XUID()))
	}
	p.SetGameMode(world.GameModeSurvival)
}

func registerCommands() {
	for _, c := range []cmd.Command{
		cmd.New("team", text.Colourf("<aqua>Team management commands.</aqua>"), []string{"t", "f"}, command.TeamCreate{}, command.TeamInvite{}, command.TeamJoin{}),
		cmd.New("whitelist", text.Colourf("<aqua>Whitelist commands.</aqua>"), []string{"wl"}, command.WhiteListAdd{}, command.WhiteListRemove{}),
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

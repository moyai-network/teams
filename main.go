package main

import (
	"github.com/bedrock-gophers/packethandler"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
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

	c.Name = text.Colourf("<bold><redstone>SYN</redstone></bold>") + "§8"
	c.Generator = func(dim world.Dimension) world.Generator { return nil }
	c.Allower = moyai.NewAllower(config.Moyai.Whitelisted)

	pk := packethandler.NewPacketListener()
	pk.Listen(&c, ":19132", []minecraft.Protocol{})
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
	p.Inventory().Handle(user.NewArmourHandler(p))
	p.Handle(user.NewHandler(p))
	p.SetGameMode(world.GameModeSurvival)
	p.Teleport(mgl64.Vec3{0, 74, 0})
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

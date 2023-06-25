package main

import (
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"math"
	"os"
)

func main() {
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

	c.Name = text.Colourf("<bold><gold>Vasar</gold></bold>") + "§8"
	//c.Allower = &allower{config.Vasar.Whitelisted}
	//c.Resources = append(c.Resources, resource.MustCompile(config.Pack.Path).WithContentKey(config.Pack.Key))
	c.Generator = func(dim world.Dimension) world.Generator { return nil }

	srv := c.New()
	srv.CloseOnProgramEnd()

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
		_ = e.Close()
	}

	for srv.Accept(nil) {
		// Do nothing.
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

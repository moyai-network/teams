package moyai

import (
	"github.com/bedrock-gophers/packethandler"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai/entity"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/oomph-ac/oomph"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
)

var srv *server.Server

func Server() *server.Server {
	return srv
}

func NewServer(conf Config, log *logrus.Logger) *server.Server {
	c := configure(conf, log)
	listen(conf.Oomph.Enabled, c, log)
	return c.New()
}

func configure(conf Config, log *logrus.Logger) server.Config {
	c, err := conf.Config(log)
	if err != nil {
		panic(err)
	}

	c.Entities = entity.Registry

	c.Name = text.Colourf("<bold><redstone>MOYAI</redstone></bold>") + "§8"
	c.Generator = func(dim world.Dimension) world.Generator { return nil }
	c.Allower = NewAllower(conf.Moyai.Whitelisted)
	return c
}

func listen(oomphEnabled bool, conf server.Config, log *logrus.Logger) {
	if oomphEnabled {
		o := oomph.New(log, oomph.OomphSettings{
			RemoteAddress:  "127.0.0.1:19133",
			Authentication: true,
		})
		o.Listen(&conf, text.Colourf("<red>Moyai</red>"), []minecraft.Protocol{}, true, false)
		go func() {
			for {
				p, err := o.Accept()
				if err != nil {
					return
				}
				_ = p
			}
		}()
		return
	}

	pk := packethandler.NewPacketListener()
	pk.Listen(&conf, ":19133", []minecraft.Protocol{})
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

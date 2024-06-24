package discord

import (
	"context"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/moyai-network/teams/cmd/discord/command"
	"github.com/moyai-network/teams/moyai"
	"log"
	"os"
)

func Run() {
	r := cmdroute.NewRouter()
	s := state.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	r.Use(cmdroute.Deferrable(s, cmdroute.DeferOpts{}))

	g, err := s.Guild(1111055709300342826)
	if err != nil {
		panic(err)
	}

	h := command.NewHandler(r, s, g.ID)
	h.RegisterCommands()

	s.AddInteractionHandler(r)
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentGuildVoiceStates)
	s.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := s.Me()
		log.Println("Bot Connected as ", me.Tag())
	})

	go func() {
		if err := s.Connect(context.TODO()); err != nil {
			log.Println("cannot connect:", err)
		}
	}()

	moyai.SetDiscordState(s)
}

package main

import (
	"github.com/moyai-network/teams/cmd/discord"
	"github.com/moyai-network/teams/cmd/konsole"
	_ "net/http/pprof"

	"github.com/joho/godotenv"
	"github.com/moyai-network/teams/cmd/minecraft"
)

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	discord.Run()
	konsole.Run()
	err := minecraft.Run()
	if err != nil {
		panic(err)
	}
}

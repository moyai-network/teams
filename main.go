package main

import (
	"github.com/joho/godotenv"
	"github.com/moyai-network/teams/cmd/discord"
	"github.com/moyai-network/teams/cmd/konsole"
	"github.com/moyai-network/teams/cmd/minecraft"
	_ "net/http/pprof"
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

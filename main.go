package main

import (
	"github.com/moyai-network/teams/cmd/konsole"
	_ "net/http/pprof"

	"github.com/moyai-network/teams/cmd/minecraft"
)

func main() {
	konsole.Run()
	err := minecraft.Run()
	if err != nil {
		panic(err)
	}
}

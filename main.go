package main

import (
	_ "net/http/pprof"

	"github.com/moyai-network/teams/cmd/minecraft"
)

func main() {
	err := minecraft.Run()
	if err != nil {
		panic(err)
	}
}

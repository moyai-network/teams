package main

import (
	"github.com/moyai-network/teams/cmd/minecraft"
	_ "net/http/pprof"
)

func main() {
	err := minecraft.Run()
	if err != nil {
		panic(err)
	}
}

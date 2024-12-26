package main

import (
	"github.com/moyai-network/teams/cmd/minecraft"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		_ = http.ListenAndServe(":8080", nil)
	}()

	err := minecraft.Run()
	if err != nil {
		panic(err)
	}
}

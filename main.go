package main

import (
	"github.com/joho/godotenv"
	"github.com/moyai-network/teams/cmd/minecraft"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	go func() {
		_ = http.ListenAndServe(":8080", nil)
	}()

	err := minecraft.Run()
	if err != nil {
		panic(err)
	}
}

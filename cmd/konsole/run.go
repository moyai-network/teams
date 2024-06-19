package konsole

import (
	"github.com/bedrock-gophers/konsole/konsole"
	"github.com/bedrock-gophers/konsole/konsole/app"
	"github.com/df-mc/dragonfly/server/player/chat"
)

func Run() {
	a := app.New("/897h6sad9ads9a8d6sa8d6as987d6216t12t6hw1278wy12876tw21t86gw1287tgw8712f8terg9eg98122REF8765WAF65GWAFR6536754DVEZ3265F7E6F7253ZE6R5723EZRF657G23SAFTYFTSFF647SAF6D4ASF67")
	go func() {
		err := a.ListenAndServe(":6969")
		if err != nil {
			panic(err)
		}
	}()

	ws := konsole.NewWebSocketServer(chat.StdoutSubscriber{}, "Mybigbootynigga69420creampieinthegayass", formatter{})
	chat.Global.Subscribe(ws)

	go func() {
		err := ws.ListenAndServe(":8080")
		if err != nil {
			panic(err)
		}
	}()
}

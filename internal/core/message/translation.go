package message

import (
	"github.com/df-mc/dragonfly/server/player/chat"
)

func Translate(key string) chat.Translation {
	return chat.Translate(NewResolver(key), 0, "msg not found")
}

package message

import "github.com/df-mc/dragonfly/server/player/chat"

type teamMessages struct{}

func (teamMessages) ErrAlreadyInTeam() chat.Translation {
	return Translate("test.test")
}

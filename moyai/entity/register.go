package entity

import (
	"github.com/bedrock-gophers/spawner/spawner"
)

func init() {
	spawner.RegisterEntityType(BlazeType{}, NewBlaze)
	spawner.RegisterEntityType(CowType{}, NewCow)
	spawner.RegisterEntityType(EndermanType{}, NewEnderman)
}

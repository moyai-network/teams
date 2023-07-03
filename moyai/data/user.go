package data

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/role"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

var (
	userCollection *mongo.Collection
)

// User is a structure containing the data of an offline user. It also contains useful functions that can be used
// externally to modify offline user data, such as roles.
type User struct {
	// xuid is the xuid of the user.
	XUID string
	// displayName is the display name of the user.
	DisplayName string
	// name is the name of the user.
	Name string
	// deviceID is the device ID of the user.
	DeviceID string
	// selfSignedID is the self-signed ID of the user.
	SelfSignedID string
	// address is the hashed IP address of the user.
	Address string
	// firstLogin is the time the user first logged in.
	FirstLogin time.Time
	// playTime is the duration the user has played for on the server.
	PlayTime time.Duration

	// Roles is the roles manager of the User.
	Roles *role.Roles
	// Mute is the mute information of the User.
	Mute moose.Punishment
	// Ban is the ban information of the User.
	Ban moose.Punishment
	// Stats contains the stats of the user.
	Stats Stats
	// Whitelisted is true if the user is whitelisted.
	Whitelisted bool
	// Frozen is the frozen state of the user.
	Frozen bool

	// Balance is the balance in the user's bank.
	Balance float64
	// Invitations is a map of the teams that invited the user, with the invitation's expiry.
	Invitations moose.MappedCoolDown[string]
	// Kits represents the kits cool-downs
	Kits moose.MappedCoolDown[string]
	// Lives is the amount of lives the user has left.
	Lives int
	// DeathBan is the death-ban cool-down
	DeathBan *moose.CoolDown
}

// Stats contains all the stats of a user.
type Stats struct {
	// Kills is the amount of players the user has killed.
	Kills uint32
	// Deaths is the amount of times the user has died.
	Deaths uint32

	// KillStreak is the current streak of kills the user has without dying.
	KillStreak uint32
	// BestKillStreak is the highest kill-streak the user has ever gotten.
	BestKillStreak uint32
}

// DefaultUser creates a default user.
func DefaultUser(xuid, name string) User {
	return User{
		XUID:        xuid,
		FirstLogin:  time.Now(),
		DisplayName: name,
		Name:        strings.ToLower(name),
		Roles:       role.NewRoles([]moose.Role{role.Default{}}, map[moose.Role]time.Time{}),
		Kits:        moose.NewMappedCoolDown[string](),
		Invitations: moose.NewMappedCoolDown[string](),
		DeathBan:    moose.NewCoolDown(),
		Balance:     250,
	}
}

// LoadUser loads a user using the given player.
func LoadUser(p *player.Player) (User, error) {
	result := userCollection.FindOne(ctx(), bson.M{"$or": []bson.M{{"name": strings.ToLower(p.Name())}, {"xuid": p.XUID()}}})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return DefaultUser(p.XUID(), p.Name()), nil
		}
		return User{}, err
	}
	var data User
	data.DeathBan = moose.NewCoolDown()
	data.Invitations = moose.NewMappedCoolDown[string]()
	data.Kits = moose.NewMappedCoolDown[string]()

	err := result.Decode(&data)
	if err != nil {
		return User{}, err
	}

	for key, value := range data.Invitations {
		if !value.Active() {
			delete(data.Invitations, key)
		}
	}
	for key, value := range data.Kits {
		if !value.Active() {
			delete(data.Kits, key)
		}
	}

	return data, nil
}

// SaveUser saves the provided user into the database.
func SaveUser(u User) error {
	filter := bson.M{"name": bson.M{"$eq": u.Name}}
	update := bson.M{"$set": u}

	res, err := userCollection.UpdateOne(ctx(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		_, err = userCollection.InsertOne(ctx(), u)
		return err
	}
	return nil
}

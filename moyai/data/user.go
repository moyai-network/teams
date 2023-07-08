package data

import (
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/role"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"sync"
	"time"
)

var (
	userCollection *mongo.Collection

	usersMu sync.Mutex
	users   = map[string]User{}
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

	// LastMessageFrom is the name of the player that sent the user a message.
	LastMessageFrom string
	// Balance is the balance in the user's bank.
	Balance float64
	// Invitations is a map of the teams that invited the user, with the invitation's expiry.
	Invitations moose.MappedCoolDown[string]
	// Kits represents the kits cool-downs.
	Kits moose.MappedCoolDown[string]
	// Lives is the amount of lives the user has left.
	Lives int
	// DeathBan is the death-ban cool-down.
	DeathBan *moose.CoolDown
	// Report is the report cool-down.
	Report *moose.CoolDown
	// SOTW is whether the user their SOTW timer enabled, or not.
	SOTW bool
	// Reclaimed is whether the user has already used their reclaim perk.
	Reclaimed bool
	// PVP is the PVP timer of the user.
	PVP *moose.CoolDown
}

func (u User) Team() (Team, bool) {
	for _, t := range Teams() {
		if t.Member(u.Name) {
			return t, true
		}
	}
	return Team{}, false
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
func DefaultUser(name string) User {
	return User{
		FirstLogin:  time.Now(),
		DisplayName: name,
		Name:        strings.ToLower(name),
		Roles:       role.NewRoles([]moose.Role{role.Default{}}, map[moose.Role]time.Time{}),
		Kits:        moose.NewMappedCoolDown[string](),
		Invitations: moose.NewMappedCoolDown[string](),
		DeathBan:    moose.NewCoolDown(),
		PVP:         moose.NewCoolDown(),
		Report:      moose.NewCoolDown(),
		Balance:     250,
		SOTW:        true,
	}
}

// LoadUser loads a user using the given name or xuid.
func LoadUser(name string) (User, error) {
	usersMu.Lock()
	defer usersMu.Unlock()

	if u, ok := users[strings.ToLower(name)]; ok {
		return u, nil
	}
	filter := bson.M{"$eq": bson.M{"name": strings.ToLower(name)}}

	result := userCollection.FindOne(ctx(), filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return DefaultUser(name), nil
		}
		return User{}, err
	}
	var u User

	u.DeathBan = moose.NewCoolDown()
	u.PVP = moose.NewCoolDown()
	u.Invitations = moose.NewMappedCoolDown[string]()
	u.Kits = moose.NewMappedCoolDown[string]()
	u.Report = moose.NewCoolDown()

	err := result.Decode(&u)
	if err != nil {
		return User{}, err
	}

	for key, value := range u.Invitations {
		if !value.Active() {
			delete(u.Invitations, key)
		}
	}
	for key, value := range u.Kits {
		if !value.Active() {
			delete(u.Kits, key)
		}
	}
	users[u.Name] = u

	return u, nil
}

// LoadUsersCond loads users using the given filter.
func LoadUsersCond(cond any) ([]User, error) {
	collection := db.Collection("users")
	count, err := collection.EstimatedDocumentCount(ctx())
	if err != nil {
		return nil, err
	}

	var data = make([]User, count)

	for d := 0; d < int(count); d++ {
		data[d] = User{
			DeathBan:    moose.NewCoolDown(),
			PVP:         moose.NewCoolDown(),
			Invitations: moose.NewMappedCoolDown[string](),
			Kits:        moose.NewMappedCoolDown[string](),
			Report:      moose.NewCoolDown(),
		}
	}

	cursor, err := collection.Find(ctx(), cond)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx(), &data); err != nil {
		return nil, err
	}

	usersMu.Lock()
	for i, u := range data {
		if d, ok := users[u.Name]; ok {
			data[i] = d
		}
	}
	usersMu.Unlock()

	return data, nil
}

// SaveUser saves the provided user into the database.
func SaveUser(u User) error {
	usersMu.Lock()
	users[u.Name] = u
	usersMu.Unlock()
	return nil
}

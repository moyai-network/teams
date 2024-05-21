package data

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/moyai-network/teams/internal/cooldown"
	"github.com/moyai-network/teams/internal/punishment"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/moyai-network/teams/moyai/tag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userCollection *mongo.Collection
	userMu         sync.Mutex
	users          = map[string]User{}
)

func userCached(f func(User) bool) (User, bool) {
	userMu.Lock()
	defer userMu.Unlock()
	for _, u := range users {
		if f(u) {
			return u, true
		}
	}
	return User{}, false
}

func saveUserData(u User) error {
	filter := bson.M{"xuid": bson.M{"$eq": u.XUID}}
	update := bson.M{"$set": u}

	res, err := userCollection.UpdateOne(ctx(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		_, err = userCollection.InsertOne(ctx(), u)
	}
	return err
}

type Stats struct {
	Kills          int `bson:"kills"`
	Deaths         int `bson:"deaths"`
	Assists        int `bson:"Assists"`
	KillStreak     int `bson:"streak"`
	BestKillStreak int `bson:"best_streak"`
}

type User struct {
	XUID        string `bson:"xuid"`
	Name        string `bson:"name"`
	DisplayName string `bson:"display_name"`
	Whitelisted bool   `bson:"whitelisted"`

	Address      string `bson:"address"`
	DeviceID     string `bson:"device_id"`
	SelfSignedID string `bson:"self_signed_id"`

	Roles    *role.Roles `bson:"roles"`
	Tags     *tag.Tags   `bson:"tags"`
	Language Language    `bson:"language"`

	Frozen bool `bson:"frozen"`
	// LastMessageFrom is the name of the player that sent the user a message.
	LastMessageFrom string

	Teams struct {
		// ChatType is the type of chat the user is in.
		ChatType int
		// Balance is the balance in the user's bank.
		Balance float64
		// Invitations is a map of the teams that invited the user, with the invitation's expiry.
		Invitations cooldown.MappedCoolDown[string]
		// Kits represents the kits cool-downs.
		Kits cooldown.MappedCoolDown[string]
		// Lives is the amount of lives the user has left.
		Lives int
		// DeathBan is the death-ban cool-down.
		DeathBan *cooldown.CoolDown
		// Report is the report cool-down.
		Report *cooldown.CoolDown
		// SOTW is whether the user their SOTW timer enabled, or not.
		SOTW bool
		// Reclaimed is whether the user has already used their reclaim perk.
		Reclaimed bool
		// PVP is the PVP timer of the user.
		PVP *cooldown.CoolDown
		// Create is the team create cooldown of the user.
		Create *cooldown.CoolDown
		// Ban is the ban of the user.
		Ban punishment.Punishment
		// Mute is the mute of the user.
		Mute punishment.Punishment
		// Dead is the live status of the logger.
		// If true, the user should be cleared.
		Dead bool
		// Stats contains the stats of the user.
		Stats Stats `bson:"stats"`
	} `bson:"nocta"`
}

func DefaultUser(name, xuid string) User {
	u := User{
		XUID:        xuid,
		Name:        strings.ToLower(name),
		DisplayName: name,
		Whitelisted: false,
	}
	u.Roles = role.NewRoles([]role.Role{}, map[role.Role]time.Time{})
	u.Tags = tag.NewTags([]tag.Tag{})
	u.Teams.Invitations = cooldown.NewMappedCoolDown[string]()
	u.Teams.Kits = cooldown.NewMappedCoolDown[string]()
	u.Teams.DeathBan = cooldown.NewCoolDown()
	u.Teams.Report = cooldown.NewCoolDown()
	u.Teams.PVP = cooldown.NewCoolDown()
	u.Teams.Create = cooldown.NewCoolDown()
	u.Teams.Stats = Stats{}

	return u
}

func LoadUserOrCreate(name, xuid string) (User, error) {
	u, err := LoadUserFromXUID(xuid)
	if errors.Is(err, mongo.ErrNoDocuments) {
		u = DefaultUser(name, xuid)
		userMu.Lock()
		users[u.XUID] = u
		userMu.Unlock()
		return u, nil
	}
	return u, err
}

func LoadUserFromXUID(xuid string) (User, error) {
	if u, ok := userCached(func(u User) bool {
		return u.XUID == xuid
	}); ok {
		return u, nil
	}
	return decodeSingleUserFromFilter(bson.M{"xuid": bson.M{"$eq": xuid}})
}

func LoadUserFromName(name string) (User, error) {
	name = strings.ToLower(name)

	if u, ok := userCached(func(u User) bool {
		return u.Name == name
	}); ok {
		return u, nil
	}

	return decodeSingleUserFromFilter(bson.M{"name": bson.M{"$eq": name}})
}

func LoadAllUsers() ([]User, error) {
	return loadUsersFromFilter(bson.M{})
}

func LoadUsersFromAddress(address string) ([]User, error) {
	filter := bson.M{"address": bson.M{"$eq": address}}
	return loadUsersFromFilter(filter)
}

func LoadUsersFromDeviceID(did string) ([]User, error) {
	filter := bson.M{"device_id": bson.M{"$eq": did}}
	return loadUsersFromFilter(filter)
}

func LoadUsersFromSelfSignedID(ssid string) ([]User, error) {
	filter := bson.M{"self_signed_id": bson.M{"$eq": ssid}}
	return loadUsersFromFilter(filter)
}

func LoadUsersFromRole(r role.Role) ([]User, error) {
	filter := bson.M{"roles": bson.M{"$elemMatch": bson.M{"$eq": r.Name()}}}
	return loadUsersFromFilter(filter)
}

func loadUsersFromFilter(filter any) ([]User, error) {
	cursor, err := userCollection.Find(ctx(), filter)
	if err != nil {
		return nil, err
	}

	var data []User
	if err = cursor.All(ctx(), &data); err != nil {
		return nil, err
	}

	for i, u := range data {
		userMu.Lock()
		if _, ok := users[u.XUID]; ok {
			data[i] = users[u.XUID]
		} else {
			users[u.XUID] = u
		}
		userMu.Unlock()
	}

	return data, nil
}

func SaveUser(u User) {
	userMu.Lock()
	users[u.XUID] = u
	userMu.Unlock()
}

func decodeSingleUserFromFilter(filter any) (User, error) {
	return decodeSingleUserResult(userCollection.FindOne(ctx(), filter))
}

func decodeSingleUserResult(result *mongo.SingleResult) (User, error) {
	var u User
	u.Roles = role.NewRoles([]role.Role{}, map[role.Role]time.Time{})
	u.Tags = tag.NewTags([]tag.Tag{})
	u.Teams.Invitations = cooldown.NewMappedCoolDown[string]()
	u.Teams.Kits = cooldown.NewMappedCoolDown[string]()
	u.Teams.DeathBan = cooldown.NewCoolDown()
	u.Teams.Report = cooldown.NewCoolDown()
	u.Teams.PVP = cooldown.NewCoolDown()
	u.Teams.Create = cooldown.NewCoolDown()
	u.Teams.Stats = Stats{}

	err := result.Decode(&u)
	if err != nil {
		return User{}, err
	}

	userMu.Lock()
	users[u.XUID] = u
	userMu.Unlock()

	return u, nil
}

func init() {
	t := time.NewTicker(60 * time.Minute)
	go func() {
		for range t.C {
			FlushCache()
		}
	}()
}

package data

import (
	"errors"
	"fmt"
	"github.com/bedrock-gophers/tag/tag"
	"github.com/moyai-network/teams/internal/core/roles"
	model2 "github.com/moyai-network/teams/internal/model"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"

	"github.com/restartfu/sets"

	"github.com/bedrock-gophers/cooldown/cooldown"
	"github.com/bedrock-gophers/role/role"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userCollection *mongo.Collection
	userMu         sync.Mutex
	users          = map[string]model2.User{}
)

func userCached(f func(model2.User) bool) (model2.User, bool) {
	userMu.Lock()
	defer userMu.Unlock()
	for _, u := range users {
		if f(u) {
			return u, true
		}
	}
	return model2.User{}, false
}

func saveUserData(u model2.User) error {
	filter := bson.M{"name": bson.M{"$eq": u.Name}}
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

func saveBatchUserData(users []model2.User) error {
	var models []mongo.WriteModel
	for _, u := range users {
		filter := bson.M{"name": bson.M{"$eq": u.Name}}
		update := bson.M{"$set": u}

		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}

	_, err := userCollection.BulkWrite(ctx(), models)
	return err
}

func LoadUserOrCreate(name, xuid string) (model2.User, error) {
	u, err := LoadUserFromName(name)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Println("LoadUserOrCreate: no user found")
		u = model2.DefaultUser(name, xuid)
		userMu.Lock()
		users[u.Name] = u
		userMu.Unlock()
		return u, nil
	}
	return u, err
}

func LoadUserFromCode(code string) (model2.User, error) {
	if u, ok := userCached(func(u model2.User) bool {
		return u.LinkCode == code
	}); ok {
		return u, nil
	}
	return decodeSingleUserFromFilter(bson.M{"link_code": bson.M{"$eq": code}})
}

func LinkUser(code string, sender *discord.User) (model2.User, error) {
	id := sender.ID.String()
	if _, err := LoadUserFromDiscordID(id); err == nil {
		return model2.User{}, errors.New("already linked")
	}
	u, err := LoadUserFromCode(code)
	if err != nil || len(u.LinkCode) == 0 || len(u.DiscordID) > 0 {
		return model2.User{}, errors.New("invalid code")
	}

	u.DiscordID = id
	u.LinkCode = ""

	u.Roles.Add(roles.Nitro())
	SaveUser(u)
	return u, nil
}

func UnlinkUser(u model2.User, s *state.State, gID discord.GuildID) error {
	if len(u.DiscordID) == 0 {
		return errors.New("not linked")
	}
	discordID, _ := strconv.Atoi(u.DiscordID)
	u.DiscordID = ""

	SaveUser(u)
	_ = s.ModifyMember(gID, discord.UserID(discordID), api.ModifyMemberData{
		Nick: option.NewString(""),
	})

	_ = s.RemoveRole(gID, discord.UserID(discordID), discord.RoleID(1255290630922436698), "Unlinking")
	return nil
}

func LoadUserFromDiscordID(did string) (model2.User, error) {
	if u, ok := userCached(func(u model2.User) bool {
		return u.DiscordID == did
	}); ok {
		return u, nil
	}
	return decodeSingleUserFromFilter(bson.M{"discord_id": bson.M{"$eq": did}})

}

func LoadUserFromName(name string) (model2.User, error) {
	name = strings.ToLower(name)

	if u, ok := userCached(func(u model2.User) bool {
		return u.Name == name
	}); ok {
		return u, nil
	}

	return decodeSingleUserFromFilter(bson.M{"name": bson.M{"$eq": name}})
}

func LoadAllUsers() ([]model2.User, error) {
	return loadUsersFromFilter(bson.M{})
}

func LoadUsersFromAddress(address string) ([]model2.User, error) {
	filter := bson.M{"address": bson.M{"$eq": address}}
	return loadUsersFromFilter(filter)
}

func LoadUsersFromDeviceID(did string) ([]model2.User, error) {
	filter := bson.M{"device_id": bson.M{"$eq": did}}
	return loadUsersFromFilter(filter)
}

func LoadUsersFromSelfSignedID(ssid string) ([]model2.User, error) {
	filter := bson.M{"self_signed_id": bson.M{"$eq": ssid}}
	return loadUsersFromFilter(filter)
}

func LoadUsersFromRole(r role.Role) ([]model2.User, error) {
	filter := bson.M{"roles.roles": bson.M{"$elemMatch": bson.M{"$eq": r.Name()}}}
	return loadUsersFromFilter(filter)
}

func loadUsersFromFilter(filter any) ([]model2.User, error) {
	cursor, err := userCollection.Find(ctx(), filter)
	if err != nil {
		return nil, err
	}

	n, err := userCollection.CountDocuments(ctx(), filter)
	if err != nil {
		return nil, err
	}
	data := make([]model2.User, n)
	for i := range data {
		data[i] = model2.DefaultUser("loadUsersFromFilter", "")
	}

	if err = cursor.All(ctx(), &data); err != nil {
		return nil, err
	}

	userMu.Lock()
	for i, u := range data {
		if _, ok := users[u.Name]; ok {
			data[i] = users[u.Name]
		} else {
			users[u.Name] = u
		}
	}
	userMu.Unlock()

	return data, nil
}

func SaveUser(u model2.User) {
	u.LastSaved = time.Now()
	userMu.Lock()
	users[u.Name] = u
	userMu.Unlock()
}

func FlushUser(u model2.User) {
	userMu.Lock()
	delete(users, u.Name)
	userMu.Unlock()

	err := saveUserData(u)
	if err != nil {
		fmt.Println(err)
	}
}

func decodeSingleUserFromFilter(filter any) (model2.User, error) {
	return decodeSingleUserResult(userCollection.FindOne(ctx(), filter))
}

func decodeSingleUserResult(result *mongo.SingleResult) (model2.User, error) {
	var u model2.User
	u.Roles = role.NewRoles([]role.Role{}, map[role.Role]time.Time{})
	u.Tags = tag.NewTags([]tag.Tag{}, tag.Tag{})
	u.Teams.Invitations = cooldown.NewMappedCoolDown[string]()
	u.Teams.Kits = cooldown.NewMappedCoolDown[string]()
	u.Teams.KOTHStart = cooldown.NewCoolDown()
	u.Teams.DeathBan = cooldown.NewCoolDown()
	u.Teams.Report = cooldown.NewCoolDown()
	u.Teams.Refill = cooldown.NewCoolDown()
	u.Teams.PVP = cooldown.NewCoolDown()
	u.Teams.Create = cooldown.NewCoolDown()
	u.Teams.Stats = model2.Stats{}
	u.Teams.ClaimedRewards = sets.New[int]()
	u.Teams.DeathInventory = &model2.Inventory{}
	u.Language = &model2.Language{}

	err := result.Decode(&u)
	if err != nil {
		return model2.User{}, err
	}

	userMu.Lock()
	users[u.Name] = u
	userMu.Unlock()

	return u, nil
}

func init() {
	t := time.NewTicker(1 * time.Minute)
	go func() {
		for range t.C {
			FlushCache()
		}
	}()
}

package data

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/moyai-network/teams/internal/sets"

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

// Settings is a structure containing the settings of a user.
type Settings struct {
	// Language is the language of the user.
	Language string
	// Display is the display settings of the user.
	Display struct {
		// ScoreboardDisabled is whether the user wants to see the scoreboard.
		ScoreboardDisabled bool
		// BossBar is whether the user wants to see their Bossbar.
		Bossbar bool
		// ActiveTag is the active tag of the user.
		ActiveTag string
	}
	Visual struct {
		// Lightning is true if lightning deaths should be enabled.
		Lightning bool
		// Splashes is true if potion splashes should be enabled.
		Splashes bool
		// PearlAnimation is true if players should appear to zoom instead of instantly teleport.
		PearlAnimation bool
	}
	// Privacy is the privacy settings of the user.
	Privacy struct {
		// PrivateMessages is whether the user wants to receive private messages.
		PrivateMessages bool
		// PublicStatistics is true if the user's statistics should be public.
		PublicStatistics bool
		// DuelRequests is true if duel requests should be allowed.
		DuelRequests bool
	}
	// Gameplay is the gameplay settings of the user.
	Gameplay struct {
		// ToggleSprint is true if the user should automatically toggle sprinting.
		ToggleSprint bool
		// AutoReapplyKit is true if the user should automatically reapply the kit.
		AutoReapplyKit bool
		// PreventInterference is true if the user should prevent interference with other players.
		PreventInterference bool
		// PreventClutter is true if clutter should be prevented.
		PreventClutter bool
		// InstantRespawn is true if the user should respawn instantly.
		InstantRespawn bool
	}
	// Advanced is a section of settings related to advanced features, such as capes or splash colours.
	Advanced struct {
		// Cape is the name of the user's cape.
		Cape string
		// ParticleMultiplier is the multiplier of combat particles.
		ParticleMultiplier int
		// PotionSplashColor is the colour of the potion splash particles.
		PotionSplashColor string
	}
}

// DefaultSettings returns the default settings.
func DefaultSettings() Settings {
	s := Settings{}

	s.Language = "en"

	s.Display.Bossbar = true
	s.Gameplay.AutoReapplyKit = true

	s.Privacy.PrivateMessages = true
	s.Privacy.DuelRequests = true
	s.Privacy.PublicStatistics = true

	s.Visual.Lightning = true
	s.Visual.Splashes = true

	return s
}

type User struct {
	XUID        string `bson:"xuid"`
	Name        string `bson:"name"`
	DisplayName string `bson:"display_name"`
	Whitelisted bool   `bson:"whitelisted"`

	DiscordID string `bson:"discord_id"`
	LinkCode  string `bson:"link_code"`

	StaffMode bool `bson:"staff_mode"`
	Vanished  bool `bson:"vanished"`

	Address      string `bson:"address"`
	DeviceID     string `bson:"device_id"`
	SelfSignedID string `bson:"self_signed_id"`

	Roles    *role.Roles `bson:"roles"`
	Tags     *tag.Tags   `bson:"tags"`
	Language *Language   `bson:"language"`

	// PlayTime is the total playtime of the user.
	PlayTime time.Duration `bson:"playtime"`
	// Frozen is whether the user is frozen.
	Frozen bool `bson:"frozen"`
	// LastMessageFrom is the name of the player that sent the user a message.
	LastMessageFrom string
	// LastVote is the last time the user voted.
	LastVote time.Time

	Teams struct {
		// Position is the position of the user.
		Position mgl64.Vec3
		// DeathInventory is the inventory of the user when they died.
		DeathInventory *Inventory `bson:"death_inventory"`
		// ChatType is the type of chat the user is in.
		ChatType int
		// Balance is the balance in the user's bank.
		Balance float64
		// Invitations is a map of the teams that invited the user, with the invitation's expiry.
		Invitations cooldown.MappedCoolDown[string]
		// Kits represents the kits cool-downs.
		Kits cooldown.MappedCoolDown[string]
		// KOTHStart is the KOTH start cool-down.
		KOTHStart *cooldown.CoolDown
		// Lives is the amount of lives the user has left.
		Lives int
		// DeathBan is the death-ban cool-down.
		DeathBan time.Time
		// Report is the report cool-down.
		Report *cooldown.CoolDown
		// Refill is the refill cool-down.
		Refill *cooldown.CoolDown
		// SOTW is whether the user their SOTW timer enabled, or not.
		SOTW bool
		// ClaimedRewards is a set of all the claimed rewards
		ClaimedRewards sets.Set[int]
		// Reclaimed is whether the user has already used their reclaim perk.
		Reclaimed bool
		// PVP is the PVP timer of the user.
		PVP *cooldown.CoolDown
		// Create is the team create cooldown of the user.
		Create *cooldown.CoolDown
		// GodApple is cooldown for god apples.
		GodApple *cooldown.CoolDown
		// Ban is the ban of the user.
		Ban punishment.Punishment
		// Mute is the mute of the user.
		Mute punishment.Punishment
		// Dead is the live status of the logger.
		// If true, the user should be cleared.
		Dead bool
		// Stats contains the stats of the user.
		Stats Stats `bson:"stats"`

		Settings Settings
	} `bson:"nocta"`
	lastSaved time.Time
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
	u.Teams.KOTHStart = cooldown.NewCoolDown()
	u.Teams.Report = cooldown.NewCoolDown()
	u.Teams.Refill = cooldown.NewCoolDown()
	u.Teams.GodApple = cooldown.NewCoolDown()
	u.Teams.PVP = cooldown.NewCoolDown()
	u.Teams.PVP.Set(time.Hour + time.Second)
	u.Teams.Create = cooldown.NewCoolDown()
	u.Teams.Stats = Stats{}
	u.Teams.ClaimedRewards = sets.New[int]()
	u.Teams.DeathInventory = &Inventory{}
	u.Language = &Language{}

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

func LoadUserFromCode(code string) (User, error) {
	if u, ok := userCached(func(u User) bool {
		return u.LinkCode == code
	}); ok {
		return u, nil
	}
	return decodeSingleUserFromFilter(bson.M{"link_code": bson.M{"$eq": code}})
}

func LinkUser(code string, sender *discord.User) (User, error) {
	id := sender.ID.String()
	if _, err := LoadUserFromDiscordID(id); err == nil {
		return User{}, errors.New("already linked")
	}
	u, err := LoadUserFromCode(code)
	if err != nil || len(u.LinkCode) == 0 || len(u.DiscordID) > 0 {
		return User{}, errors.New("invalid code")
	}

	u.DiscordID = id
	u.LinkCode = ""

	u.Roles.Add(role.Nitro{})
	SaveUser(u)
	return u, nil
}

func UnlinkUser(u User, s *state.State, gID discord.GuildID) error {
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

func LoadUserFromDiscordID(did string) (User, error) {
	if u, ok := userCached(func(u User) bool {
		return u.DiscordID == did
	}); ok {
		return u, nil
	}
	return decodeSingleUserFromFilter(bson.M{"discord_id": bson.M{"$eq": did}})

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
	filter := bson.M{"roles.roles": bson.M{"$elemMatch": bson.M{"$eq": r.Name()}}}
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
	u.lastSaved = time.Now()
	userMu.Lock()
	users[u.XUID] = u
	userMu.Unlock()

	go func() {
		err := saveUserData(u)
		if err != nil {
			log.Println("Error saving user data:", err)
		}
	}()
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
	u.Teams.KOTHStart = cooldown.NewCoolDown()
	u.Teams.Report = cooldown.NewCoolDown()
	u.Teams.Refill = cooldown.NewCoolDown()
	u.Teams.PVP = cooldown.NewCoolDown()
	u.Teams.Create = cooldown.NewCoolDown()
	u.Teams.Stats = Stats{}
	u.Teams.ClaimedRewards = sets.New[int]()
	u.Teams.DeathInventory = &Inventory{}
	u.Language = &Language{}

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

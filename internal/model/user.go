package model

import (
	"github.com/bedrock-gophers/cooldown/cooldown"
	"github.com/bedrock-gophers/role/role"
	"github.com/bedrock-gophers/tag/tag"
	"github.com/restartfu/sets"
	"strings"
	"time"
)

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
	Tags     *tag.Tags `bson:"tags"`
	Language *Language `bson:"language"`

	// PlayTime is the total playtime of the user.
	PlayTime time.Duration `bson:"playtime"`
	// Frozen is whether the user is frozen.
	Frozen bool `bson:"frozen"`
	// LastMessageFrom is the name of the player that sent the user a message.
	LastMessageFrom string
	// LastVote is the last time the user voted.
	LastVote time.Time

	Teams struct {
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
		DeathBan *cooldown.CoolDown
		// DeathBanned is whether the user is/was death-banned.
		DeathBanned bool
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
		Ban Punishment
		// Mute is the mute of the user.
		Mute Punishment
		// Dead is the live status of the logger.
		// If true, the user should be cleared.
		//Dead bool
		// Stats contains the stats of the user.
		Stats Stats `bson:"stats"`

		Settings Settings
	} `bson:"nocta"`
	LastSaved time.Time
}

func DefaultUser(name, xuid string) User {
	u := User{
		XUID:        xuid,
		Name:        strings.ToLower(name),
		DisplayName: name,
		Whitelisted: false,
	}
	u.Roles = role.NewRoles([]role.Role{}, map[role.Role]time.Time{})
	u.Tags = tag.NewTags([]tag.Tag{}, tag.Tag{})
	u.Teams.Invitations = cooldown.NewMappedCoolDown[string]()
	u.Teams.Kits = cooldown.NewMappedCoolDown[string]()
	u.Teams.KOTHStart = cooldown.NewCoolDown()
	u.Teams.Report = cooldown.NewCoolDown()
	u.Teams.Refill = cooldown.NewCoolDown()
	u.Teams.GodApple = cooldown.NewCoolDown()
	u.Teams.DeathBan = cooldown.NewCoolDown()
	u.Teams.PVP = cooldown.NewCoolDown()
	u.Teams.PVP.Set(time.Hour + time.Second)
	if !u.Teams.PVP.Paused() {
		u.Teams.PVP.TogglePause()
	}
	u.Teams.Create = cooldown.NewCoolDown()
	u.Teams.Stats = Stats{}
	u.Teams.ClaimedRewards = sets.New[int]()
	u.Teams.DeathInventory = &Inventory{}
	u.Language = &Language{}

	return u
}

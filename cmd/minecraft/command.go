package minecraft

import (
	"github.com/bedrock-gophers/knockback/knockback"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/adapter/command"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/roles"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// registerCommands registers all commands that are available in the server.
func registerCommands() {
	for _, c := range []cmd.Command{
		cmd.New("knockback", "", []string{"kb"}, knockback.Menu{Allower: operatorAllower{}}),
		cmd.New("unlink", text.Colourf("Unlink your discord account."), nil, command.Unlink{}),
		cmd.New("link", text.Colourf("Link your discord account."), nil, command.Link{}),
		cmd.New("revive", text.Colourf("Revive a player."), nil, command.Revive{}),
		cmd.New("prizes", "play time rewards", nil, command.Prizes{}),
		cmd.New("lives", "lives management commands", nil, command.Lives{}, command.LivesGiveOnline{}, command.LivesGiveOffline{}),
		cmd.New("staff", text.Colourf("Staff management commands."), nil, command.StaffMode{}),
		cmd.New("rename", text.Colourf("Rename your items."), nil, command.Rename{}),
		cmd.New("stop", text.Colourf("Stop the server."), nil, command.Stop{}),
		cmd.New("pots", text.Colourf("Place potion chests."), nil, command.Pots{}),
		cmd.New("fix", text.Colourf("Fix your inventory."), nil, command.Fix{}, command.FixAll{}),
		cmd.New("chat", text.Colourf("Chat management commands."), nil, command.ChatMute{}, command.ChatUnMute{}, command.ChatCoolDown{}),
		cmd.New("t", text.Colourf("The main team management command."), []string{"f"},
			command.TeamCreate{},
			command.TeamRename{},
			command.TeamInformation{},
			command.TeamDisband{},
			command.TeamInvite{},
			command.TeamJoin{},
			command.TeamWho{},
			command.TeamLeave{},
			command.TeamKick{},
			command.TeamPromote{},
			command.TeamDemote{},
			command.TeamTop{},
			command.TeamClaim{},
			command.TeamUnClaim{},
			command.TeamSetHome{},
			command.TeamHome{},
			command.TeamList{},
			command.TeamUnFocus{},
			command.TeamFocusPlayer{},
			command.TeamFocusTeam{},
			command.TeamChat{},
			command.TeamWithdraw{},
			command.TeamW{},
			command.TeamDeposit{},
			command.TeamD{},
			command.TeamWithdrawAll{},
			command.TeamWAll{},
			command.TeamDepositAll{},
			command.TeamDAll{},
			command.TeamStuck{},
			command.TeamRally{},
			command.TeamUnRally{},
			command.TeamMap{},
			command.TeamClearMap{},
			command.TeamSetDTR{},
			command.TeamDelete{},
			command.TeamCamp{},
			command.TeamIncrementDTR{},
			command.TeamDecrementDTR{},
			command.TeamResetRegen{},
			command.TeamSetPoints{},
		), cmd.New("whitelist", text.Colourf("Whitelist commands."), []string{"wl"}, command.WhiteListAdd{}, command.WhiteListRemove{}),
		cmd.New("tl", text.Colourf("Send your location to teammates."), nil, command.TL{}),
		cmd.New("balance", text.Colourf("Manage your balance."), []string{"bal"}, command.Balance{}, command.BalancePayOnline{}, command.BalancePayOffline{}, command.BalanceAdd{}, command.BalanceAddOffline{}),
		cmd.New("colour", text.Colourf("Customize the colour of your archer."), []string{"color", "dye"}, command.Colour{}),
		cmd.New("clear", text.Colourf("Clear your Inventory."), nil, command.Clear{}),
		cmd.New("clearlag", text.Colourf("Clears all ground entitys."), nil, command.ClearLag{}),
		cmd.New("logout", text.Colourf("Safely logout of PVP."), nil, command.Logout{}),
		cmd.New("pvp", text.Colourf("Manage PVP timer."), nil, command.PvpEnable{}),
		cmd.New("role", text.Colourf("Manage user roles."), nil, command.RoleAdd{}, command.RoleRemove{}, command.RoleAddOffline{}, command.RoleRemoveOffline{}, command.RoleList{}),
		cmd.New("teleport", text.Colourf("Teleport yourself or another player to a position."), []string{"tp"}, command.TeleportToPos{}, command.TeleportTargetsToPos{}, command.TeleportTargetsToTarget{}, command.TeleportToTarget{}),
		cmd.New("reclaim", text.Colourf("Manage a reclaim."), nil, command.Reclaim{}, command.ReclaimReset{}),
		cmd.New("kit", text.Colourf("Choose a kit."), nil, command.Kit{}, command.KitReset{}),
		cmd.New("ban", text.Colourf("Unleash the ban hammer."), nil, command.Ban{}, command.BanOffline{}, command.BanList{}, command.BanLiftOffline{} /*command.BanInfoOffline{},*/, command.BanForm{}),
		cmd.New("kick", text.Colourf("Kick a user."), nil, command.Kick{}),
		cmd.New("mute", text.Colourf("Mute a user."), nil, command.MuteList{}, command.MuteInfo{}, command.MuteInfoOffline{}, command.MuteLift{}, command.MuteLiftOffline{}, command.MuteForm{}, command.Mute{}, command.MuteOffline{}),
		cmd.New("whisper", text.Colourf("Send a private message to a player."), []string{"w", "tell", "msg"}, command.Whisper{}),
		cmd.New("reply", text.Colourf("Reply to the last whispered player."), []string{"r"}, command.Reply{}),
		cmd.New("fly", text.Colourf("Toggle flight."), nil, command.Fly{}),
		cmd.New("sotw", text.Colourf("SOTW management commands."), nil, command.SOTWStart{}, command.SOTWEnd{}, command.SOTWDisable{}),
		cmd.New("freeze", text.Colourf("Freeze possible cheaters."), nil, command.Freeze{}),
		cmd.New("gamemode", text.Colourf("Manage gamemodes."), []string{"gm"}, command.GameMode{}),
		cmd.New("key", text.Colourf("Manage keys"), nil, command.Key{}, command.KeyAll{}),
		cmd.New("koth", text.Colourf("Manage KOTHs."), nil, command.KothStart{}, command.KothStop{}, command.KothList{}),
		cmd.New("pp", text.Colourf("Manage partner packages."), nil, command.PartnerPackageAll{}, command.PartnerPackage{}, command.PartnerItemsRefresh{}),
		cmd.New("ping", text.Colourf("Check your ping."), nil, command.Ping{}),
		cmd.New("data", text.Colourf("Clear data."), nil, command.DataReset{}),
		cmd.New("vanish", text.Colourf("Vanish as staff."), []string{"v"}, command.Vanish{}),
		cmd.New("lang", text.Colourf("Change your language."), nil, lang.Lang{}),
		cmd.New("blockshop", text.Colourf("Access the blockshop to buy items."), nil, command.BlockShop{}),
		cmd.New("enderchest", text.Colourf("Access your enderchest."), []string{"ec"}, command.Enderchest{}),
		cmd.New("blackmarket", text.Colourf("Access the secret items of the black market"), nil, command.BlackMarket{}),
		cmd.New("trim", text.Colourf("Add trims to your armor"), nil, command.Trim{}, command.TrimClear{}),
		cmd.New("end", text.Colourf("End your adventure."), nil, command.End{}),
		cmd.New("nether", text.Colourf("End your adventure."), nil, command.Nether{}),
		cmd.New("settings", text.Colourf("Access your settings."), nil, command.Settings{}),
		cmd.New("tag", text.Colourf("Manage your tags."), nil, command.TagAddOnline{}, command.TagAddOffline{}, command.TagRemoveOnline{}, command.TagRemoveOffline{}, command.TagSet{}),
		cmd.New("cape", text.Colourf("Manage your capes."), nil, command.Cape{}),
		cmd.New("conquest", text.Colourf("Manage Conquest."), nil, command.ConquestStart{}, command.ConquestStop{}),
		cmd.New("eotw", text.Colourf("EOTW management commands."), nil, command.EOTWStart{}, command.EOTWEnd{}),
		cmd.New("lastinv", text.Colourf("Access last inventory of players."), nil, command.LastInv{}),
		cmd.New("leaderboards", text.Colourf("See the leaderboards for kills, deaths, killstreaks, and KDR"), []string{"lb"}, command.LeaderboardKills{}, command.LeaderboardDeaths{}, command.LeaderboardKillStreaks{}, command.LeaderboardKDR{}),
		cmd.New("stats", text.Colourf("View your or other player's stats."), nil, command.StatsOnlineCommand{}, command.StatsOfflineCommand{}),
		cmd.New("report", text.Colourf("Report a player."), nil, command.Report{}),
	} {
		cmd.Register(c)
	}
}

type operatorAllower struct{}

func (operatorAllower) Allow(src cmd.Source) bool {
	p, ok := src.(*player.Player)
	if !ok {
		return false
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return false
	}

	return u.Roles.Contains(roles.Operator())
}

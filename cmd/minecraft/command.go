package minecraft

import (
	"github.com/bedrock-gophers/role/role"
	"github.com/df-mc/dragonfly/server/cmd"
	command2 "github.com/moyai-network/teams/internal/adapter/command"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// registerCommands registers all commands that are available in the server.
func registerCommands() {
	for _, c := range []cmd.Command{
		//cmd.New("knockback", text.Colourf("Manage server KB"), []string{"kb"}, knockback.Menu{Allower: operatorAllower{}}),
		cmd.New("alias", text.Colourf("Find aliases of a player."), nil, command2.AliasOffline{}, command2.AliasOnline{}),
		cmd.New("unlink", text.Colourf("Unlink your discord account."), nil, command2.Unlink{}),
		cmd.New("link", text.Colourf("Link your discord account."), nil, command2.Link{}),
		cmd.New("revive", text.Colourf("Revive a player."), nil, command2.Revive{}),
		cmd.New("prizes", "play time rewards", nil, command2.Prizes{}),
		cmd.New("lives", "lives management commands", nil, command2.Lives{}, command2.LivesGiveOnline{}, command2.LivesGiveOffline{}),
		cmd.New("staff", text.Colourf("Staff management commands."), nil, command2.StaffMode{}),
		cmd.New("rename", text.Colourf("Rename your items."), nil, command2.Rename{}),
		cmd.New("stop", text.Colourf("Stop the server."), nil, command2.Stop{}),
		cmd.New("pots", text.Colourf("Place potion chests."), nil, command2.Pots{}),
		cmd.New("fix", text.Colourf("Fix your inventory."), nil, command2.Fix{}, command2.FixAll{}),
		cmd.New("chat", text.Colourf("Chat management commands."), nil, command2.ChatMute{}, command2.ChatUnMute{}, command2.ChatCoolDown{}),
		cmd.New("t", text.Colourf("The main team management command."), []string{"f"},
			command2.TeamCreate{},
			command2.TeamRename{},
			command2.TeamInformation{},
			command2.TeamDisband{},
			command2.TeamInvite{},
			command2.TeamJoin{},
			command2.TeamWho{},
			command2.TeamLeave{},
			command2.TeamKick{},
			command2.TeamPromote{},
			command2.TeamDemote{},
			command2.TeamTop{},
			command2.TeamClaim{},
			command2.TeamUnClaim{},
			command2.TeamSetHome{},
			command2.TeamHome{},
			command2.TeamList{},
			command2.TeamUnFocus{},
			command2.TeamFocusPlayer{},
			command2.TeamFocusTeam{},
			command2.TeamChat{},
			command2.TeamWithdraw{},
			command2.TeamW{},
			command2.TeamDeposit{},
			command2.TeamD{},
			command2.TeamWithdrawAll{},
			command2.TeamWAll{},
			command2.TeamDepositAll{},
			command2.TeamDAll{},
			command2.TeamStuck{},
			command2.TeamRally{},
			command2.TeamUnRally{},
			command2.TeamMap{},
			command2.TeamClearMap{},
			command2.TeamSetDTR{},
			command2.TeamDelete{},
			command2.TeamCamp{},
			command2.TeamIncrementDTR{},
			command2.TeamDecrementDTR{},
			command2.TeamResetRegen{},
			command2.TeamSetPoints{},
		), cmd.New("whitelist", text.Colourf("Whitelist commands."), []string{"wl"}, command2.WhiteListAdd{}, command2.WhiteListRemove{}),
		cmd.New("tl", text.Colourf("Send your location to teammates."), nil, command2.TL{}),
		cmd.New("balance", text.Colourf("Manage your balance."), []string{"bal"}, command2.Balance{}, command2.BalancePayOnline{}, command2.BalancePayOffline{}, command2.BalanceAdd{}, command2.BalanceAddOffline{}),
		cmd.New("colour", text.Colourf("Customize the colour of your archer."), []string{"color", "dye"}, command2.Colour{}),
		cmd.New("clear", text.Colourf("Clear your Inventory."), nil, command2.Clear{}),
		cmd.New("clearlag", text.Colourf("Clears all ground entitys."), nil, command2.ClearLag{}),
		cmd.New("logout", text.Colourf("Safely logout of PVP."), nil, command2.Logout{}),
		cmd.New("pvp", text.Colourf("Manage PVP timer."), nil, command2.PvpEnable{}),
		cmd.New("role", text.Colourf("Manage user roles."), nil, command2.RoleAdd{}, command2.RoleRemove{}, command2.RoleAddOffline{}, command2.RoleRemoveOffline{}, command2.RoleList{}),
		cmd.New("teleport", text.Colourf("Teleport yourself or another player to a position."), []string{"tp"}, command2.TeleportToPos{}, command2.TeleportTargetsToPos{}, command2.TeleportTargetsToTarget{}, command2.TeleportToTarget{}),
		cmd.New("reclaim", text.Colourf("Manage a reclaim."), nil, command2.Reclaim{}, command2.ReclaimReset{}),
		cmd.New("kit", text.Colourf("Choose a kit."), nil, command2.Kit{}, command2.KitReset{}),
		cmd.New("ban", text.Colourf("Unleash the ban hammer."), nil, command2.Ban{}, command2.BanOffline{}, command2.BanList{}, command2.BanLiftOffline{} /*command.BanInfoOffline{},*/, command2.BanForm{}),
		//cmd.New("blacklist", text.Colourf("Blacklist little evaders."), nil, command.Blacklist{}, command.BlacklistOffline{}, command.BlacklistList{}, command.BlacklistLiftOffline{}, command.BlacklistInfoOffline{}, command.BlacklistForm{}),
		cmd.New("kick", text.Colourf("Kick a user."), nil, command2.Kick{}),
		cmd.New("mute", text.Colourf("Mute a user."), nil, command2.MuteList{}, command2.MuteInfo{}, command2.MuteInfoOffline{}, command2.MuteLift{}, command2.MuteLiftOffline{}, command2.MuteForm{}, command2.Mute{}, command2.MuteOffline{}),
		cmd.New("whisper", text.Colourf("Send a private message to a player."), []string{"w", "tell", "msg"}, command2.Whisper{}),
		cmd.New("reply", text.Colourf("Reply to the last whispered player."), []string{"r"}, command2.Reply{}),
		cmd.New("fly", text.Colourf("Toggle flight."), nil, command2.Fly{}),
		cmd.New("sotw", text.Colourf("SOTW management commands."), nil, command2.SOTWStart{}, command2.SOTWEnd{}, command2.SOTWDisable{}),
		cmd.New("freeze", text.Colourf("Freeze possible cheaters."), nil, command2.Freeze{}),
		cmd.New("gamemode", text.Colourf("Manage gamemodes."), []string{"gm"}, command2.GameMode{}),
		cmd.New("key", text.Colourf("Manage keys"), nil, command2.Key{}, command2.KeyAll{}),
		cmd.New("koth", text.Colourf("Manage KOTHs."), nil, command2.KothStart{}, command2.KothStop{}, command2.KothList{}),
		cmd.New("pp", text.Colourf("Manage partner packages."), nil, command2.PartnerPackageAll{}, command2.PartnerPackage{}, command2.PartnerItemsRefresh{}),
		cmd.New("ping", text.Colourf("Check your ping."), nil, command2.Ping{}),
		cmd.New("data", text.Colourf("Clear data."), nil, command2.DataReset{}),
		//cmd.New("nick", text.Colourf("Change your nickname."), nil, command.NickReset{}, command.Nick{}),
		cmd.New("vanish", text.Colourf("Vanish as staff."), []string{"v"}, command2.Vanish{}),
		cmd.New("lang", text.Colourf("Change your language."), nil, lang.Lang{}),
		cmd.New("blockshop", text.Colourf("Access the blockshop to buy items."), nil, command2.BlockShop{}),
		cmd.New("enderchest", text.Colourf("Access your enderchest."), []string{"ec"}, command2.Enderchest{}),
		cmd.New("blackmarket", text.Colourf("Access the secret items of the black market"), nil, command2.BlackMarket{}),
		cmd.New("trim", text.Colourf("Add trims to your armor"), nil, command2.Trim{}, command2.TrimClear{}),
		cmd.New("end", text.Colourf("End your adventure."), nil, command2.End{}),
		cmd.New("nether", text.Colourf("End your adventure."), nil, command2.Nether{}),
		cmd.New("settings", text.Colourf("Access your settings."), nil, command2.Settings{}),
		cmd.New("tag", text.Colourf("Manage your tags."), nil, command2.TagAddOnline{}, command2.TagAddOffline{}, command2.TagRemoveOnline{}, command2.TagRemoveOffline{}, command2.TagSet{}),
		cmd.New("cape", text.Colourf("Manage your capes."), nil, command2.Cape{}),
		cmd.New("conquest", text.Colourf("Manage Conquest."), nil, command2.ConquestStart{}, command2.ConquestStop{}),
		cmd.New("eotw", text.Colourf("EOTW management commands."), nil, command2.EOTWStart{}, command2.EOTWEnd{}),
		cmd.New("lastinv", text.Colourf("Access last inventory of players."), nil, command2.LastInv{}),
		cmd.New("leaderboards", text.Colourf("See the leaderboards for kills, deaths, killstreaks, and KDR"), []string{"lb"}, command2.LeaderboardKills{}, command2.LeaderboardDeaths{}, command2.LeaderboardKillStreaks{}, command2.LeaderboardKDR{}),
		cmd.New("stats", text.Colourf("View your or other player's stats."), nil, command2.StatsOnlineCommand{}, command2.StatsOfflineCommand{}),
		cmd.New("report", text.Colourf("Report a player."), nil, command2.Report{}),
	} {
		cmd.Register(c)
	}

	//cmd.Register(cmd.New("hub", text.Colourf("Return to the Moyai Hub."), []string{"lobby"}, command.Hub{}))
}

// operatorAllower is an allower that allows all users with the operator role to execute a command.
type operatorAllower struct{}

// Allow ...
func (operatorAllower) Allow(s cmd.Source) bool {
	return command2.Allow(s, true, []role.Role{}...)
}

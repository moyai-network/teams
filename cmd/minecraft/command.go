package minecraft

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/command"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// registerCommands registers all commands that are available in the server.
func registerCommands(srv *server.Server) {
	for _, c := range []cmd.Command{
		cmd.New("staff", text.Colourf("Staff management commands."), nil, command.StaffMode{}),
		cmd.New("rename", text.Colourf("Rename your items."), nil, command.Rename{}),
		cmd.New("stop", text.Colourf("Stop the server."), nil, command.NewStop(srv)),
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
			command.TeamDeposit{},
			command.TeamWithdrawAll{},
			command.TeamDepositAll{},
			command.TeamStuck{},
			command.TeamRally{},
			command.TeamUnRally{},
			command.TeamMap{},
			command.TeamSetDTR{},
			command.TeamDelete{},
		), cmd.New("whitelist", text.Colourf("Whitelist commands."), []string{"wl"}, command.WhiteListAdd{}, command.WhiteListRemove{}),
		cmd.New("tl", text.Colourf("Send your location to teammates"), nil, command.TL{}),
		cmd.New("balance", text.Colourf("Manage your balance."), []string{"bal"}, command.Balance{}, command.BalancePayOnline{}, command.BalancePayOffline{}, command.BalanceAdd{}, command.BalanceAddOffline{}),
		cmd.New("colour", text.Colourf("Customize the colour of your archer."), nil, command.Colour{}),
		cmd.New("clear", text.Colourf("Clear your Inventory."), nil, command.Clear{}),
		cmd.New("clearlag", text.Colourf("Clears all ground entitys."), nil, command.ClearLag{}),
		cmd.New("logout", text.Colourf("Safely logout of PVP."), nil, command.Logout{}),
		cmd.New("pvp", text.Colourf("Manage PVP timer."), nil, command.PvpEnable{}),
		cmd.New("role", text.Colourf("Manage user roles."), nil, command.RoleAdd{}, command.RoleRemove{}, command.RoleAddOffline{}, command.RoleRemoveOffline{}, command.RoleList{}),
		cmd.New("teleport", text.Colourf("Teleport yourself or another player to a position."), []string{"tp"}, command.TeleportToPos{}, command.TeleportTargetsToPos{}, command.TeleportTargetsToTarget{}, command.TeleportToTarget{}),
		cmd.New("reclaim", text.Colourf("Manage a reclaim."), nil, command.Reclaim{}, command.ReclaimReset{}),
		cmd.New("kit", text.Colourf("Choose a kit."), nil, command.Kit{}, command.KitReset{}),
		cmd.New("ban", text.Colourf("Unleash the ban hammer."), nil, command.Ban{}, command.BanOffline{}, command.BanList{}, command.BanLiftOffline{} /*command.BanInfoOffline{},*/, command.BanForm{}),
		//cmd.New("blacklist", text.Colourf("Blacklist little evaders."), nil, command.Blacklist{}, command.BlacklistOffline{}, command.BlacklistList{}, command.BlacklistLiftOffline{}, command.BlacklistInfoOffline{}, command.BlacklistForm{}),
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
		cmd.New("pp", text.Colourf("Manage partner packages."), nil, command.PartnerPackageAll{}, command.PartnerPackage{}),
		cmd.New("ping", text.Colourf("Check your ping."), nil, command.Ping{}),
		//cmd.New("data", text.Colourf("Clear data."), nil, command.DataReset{}),
		cmd.New("nick", text.Colourf("Change your nickname."), nil, command.Nick{}, command.NickReset{}),
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
	} {
		cmd.Register(c)
	}

	//cmd.Register(cmd.New("hub", text.Colourf("Return to the Moyai Hub."), []string{"lobby"}, command.Hub{}))
}

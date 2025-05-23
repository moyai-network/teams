package command

import (
	"github.com/bedrock-gophers/tag/tag"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/core"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// TagAddOnline is a command that can be used to add a tag to a player online.
type TagAddOnline struct {
	adminAllower
	Sub    cmd.SubCommand `cmd:"add"`
	Target []cmd.Target
	Tag    tagList
}

// TagAddOffline is a command that can be used to add a tag to a player offline.
type TagAddOffline struct {
	adminAllower
	Sub    cmd.SubCommand `cmd:"add"`
	Target string
	Tag    tagList
}

// TagRemoveOnline is a command that can be used to remove a tag from a player online.
type TagRemoveOnline struct {
	adminAllower
	Sub    cmd.SubCommand `cmd:"remove"`
	Target []cmd.Target
	Tag    tagList
}

// TagRemoveOffline is a command that can be used to remove a tag from a player offline.
type TagRemoveOffline struct {
	adminAllower
	Sub    cmd.SubCommand `cmd:"remove"`
	Target string
	Tag    tagList
}

// TagSet is a command that can be used to set your active tag.
type TagSet struct {
	Sub cmd.SubCommand `cmd:"set"`
	Tag ownedTagList
}

// Run ...
func (t TagAddOnline) Run(_ cmd.Source, out *cmd.Output, tx *world.Tx) {
	target, ok := t.Target[0].(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(target.Name())
	if !ok {
		return
	}

	tg, _ := tag.ByName(string(t.Tag))
	if u.Tags.Contains(tg) {
		out.Errorf("The player already owns this tag.")
		return
	}
	u.Tags.Add(tg)

	core.UserRepository.Save(u)
	out.Print(text.Colourf("<green>The tag has been added to the player.</green>"))
}

// Run ...
func (t TagAddOffline) Run(_ cmd.Source, out *cmd.Output, tx *world.Tx) {
	u, ok := core.UserRepository.FindByName(t.Target)
	if !ok {
		return
	}

	tg, _ := tag.ByName(string(t.Tag))
	if u.Tags.Contains(tg) {
		out.Errorf("The player already owns this tag.")
		return
	}
	u.Tags.Add(tg)

	core.UserRepository.Save(u)
	out.Print(text.Colourf("<green>The tag has been added to the player.</green>"))
}

// Run ...
func (t TagRemoveOnline) Run(_ cmd.Source, out *cmd.Output, tx *world.Tx) {
	target, ok := t.Target[0].(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(target.Name())
	if !ok {
		return
	}

	tg, _ := tag.ByName(string(t.Tag))
	if !u.Tags.Contains(tg) {
		out.Errorf("The player does not own this tag.")
		return
	}
	u.Tags.Remove(tg)

	core.UserRepository.Save(u)
	out.Print(text.Colourf("<green>The tag has been removed from the player.</green>"))
}

// Run ...
func (t TagRemoveOffline) Run(_ cmd.Source, out *cmd.Output, tx *world.Tx) {
	u, ok := core.UserRepository.FindByName(t.Target)
	if !ok {
		out.Errorf("The player does not exist.")
		return
	}

	tg, _ := tag.ByName(string(t.Tag))
	if !u.Tags.Contains(tg) {
		out.Errorf("The player does not own this tag.")
		return
	}
	u.Tags.Remove(tg)

	core.UserRepository.Save(u)
	out.Print(text.Colourf("<green>The tag has been removed from the player.</green>"))
}

// Run ...
func (t TagSet) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	if string(t.Tag) == "none" {
		out.Print(text.Colourf("<green>Your active tag has been removed.</green>"))
		u.Teams.Settings.Display.ActiveTag = ""
		core.UserRepository.Save(u)
		return
	}

	tg, ok := tag.ByName(string(t.Tag))
	if !ok {
		out.Errorf("The tag does not exist.")
		return
	}

	if !u.Tags.Contains(tg) {
		out.Errorf("You do not own this tag.")
		return
	}

	if u.Teams.Settings.Display.ActiveTag == string(t.Tag) {
		out.Errorf("The tag is already your active tag.")
		return
	}

	u.Teams.Settings.Display.ActiveTag = tg.Name()
	core.UserRepository.Save(u)
	out.Print(text.Colourf("<green>Your active tag has been set to </green>%s<green>.</green>", tg.Format()))
}

// ownedTagList is a type that implements the cmd.Enum interface for the tag command.
type ownedTagList string

// Type ...
func (ownedTagList) Type() string {
	return "ownedTag"
}

// Options ...
func (ownedTagList) Options(src cmd.Source) (list []string) {
	list = append(list, "none")
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}

	for _, t := range u.Tags.All() {
		list = append(list, t.Name())
	}

	return list
}

// tagList is a type that implements the cmd.Enum interface for the tag command.
type tagList string

// Type ...
func (tagList) Type() string {
	return "tag"
}

// Options ...
func (tagList) Options(cmd.Source) (list []string) {
	for _, t := range tag.All() {
		list = append(list, t.Name())
	}
	return list
}
